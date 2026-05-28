package generator

import (
	"go/ast"
	"go/parser"
	"go/token"
)

// standardErrorHelperSrc is the source for the per-package error-envelope
// helpers that every generated handlers.go file embeds. It declares:
//
//   - ErrorHandler — function-type alias matching every codegen error site.
//   - Option / WithErrorHandler — variadic option for NewHandler.
//   - (*Handler).SetErrorHandler — setter for the aggregator-loop pattern
//     (one Handler.SetErrorHandler call per generated package, instead of
//     threading WithErrorHandler through every per-package constructor).
//   - DefaultErrorHandler — the {code,error,req_id} JSON envelope. Exposed
//     so consumers can wrap it from a custom ErrorHandler.
//
// The "package _" prefix is replaced with the generated package's name when
// these decls are appended to the file by GenerateHandlersFile.
const standardErrorHelperSrc = `package _

// ErrorHandler is invoked by generated handlers at every internal failure
// (request parse error, unsupported Content-Type, handler returning nil,
// JSON encode failure). It is a type alias rather than a named type so a
// single setter signature can drive Handlers from multiple generated
// packages without per-package conversions.
type ErrorHandler = func(w http.ResponseWriter, r *http.Request, status int, msg string)

type Option func(*Handler)

// WithErrorHandler replaces the default {code,error,req_id} envelope with
// the supplied function for this Handler. Pass it to NewHandler.
func WithErrorHandler(eh ErrorHandler) Option {
	return func(h *Handler) { h.errorHandler = eh }
}

// SetErrorHandler is the post-construction equivalent of WithErrorHandler.
// Useful when Handlers are constructed by an injector (fx, wire, ...) and
// the error handler is configured later.
func (h *Handler) SetErrorHandler(eh ErrorHandler) { h.errorHandler = eh }

var statusToCode = map[int]string{
	400: "BadRequest",
	401: "Unauthorized",
	403: "Forbidden",
	404: "NotFound",
	409: "Conflict",
	415: "UnsupportedMediaType",
	429: "TooManyRequests",
	500: "InternalServerError",
}

// DefaultErrorHandler writes the standard {code,error,req_id} JSON envelope.
// Status code -> "code" mapping is the canonical Go HTTP name (400 ->
// "BadRequest", 415 -> "UnsupportedMediaType", ...) and falls back to
// "Error" for unmapped statuses. "req_id" is read via chi's RequestID
// middleware; mount it to populate the field, otherwise it stays empty.
var DefaultErrorHandler ErrorHandler = func(w http.ResponseWriter, r *http.Request, status int, msg string) {
	code, ok := statusToCode[status]
	if !ok {
		code = "Error"
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"code":   code,
		"error":  msg,
		"req_id": chimw.GetReqID(r.Context()),
	})
}
`

func parseStandardErrorHelperDecls() []ast.Decl {
	file, err := parser.ParseFile(token.NewFileSet(), "", standardErrorHelperSrc, 0)
	if err != nil {
		panic(err)
	}
	return file.Decls
}

// AddStandardErrorDecls appends the ErrorHandler/Option/SetErrorHandler/
// DefaultErrorHandler decls to the generated handlers file. Safe to call
// more than once per file — only the first call adds the decls.
func (g *Generator) AddStandardErrorDecls() {
	for _, d := range g.HandlersFile.extraDecls {
		if v, ok := d.(*ast.GenDecl); ok {
			for _, spec := range v.Specs {
				if vs, ok := spec.(*ast.ValueSpec); ok {
					for _, name := range vs.Names {
						if name.Name == "DefaultErrorHandler" {
							return
						}
					}
				}
			}
		}
	}
	g.AddHandlersImport("encoding/json")
	g.AddHandlersImport("net/http")
	g.AddHandlersImportWithAlias("chimw", "github.com/go-chi/chi/v5/middleware")
	g.HandlersFile.extraDecls = append(g.HandlersFile.extraDecls, parseStandardErrorHelperDecls()...)
}

// writeStandardErrorCall builds the AST for the call site
//
//	h.errorHandler(w, r, http.<StatusConst>, <msgExpr>)
//
// every generated error site invokes. The receiver routes the call through
// the user-supplied ErrorHandler (set via WithErrorHandler/SetErrorHandler)
// or DefaultErrorHandler when no override was supplied.
func writeStandardErrorCall(statusConst string, msgExpr ast.Expr) *ast.ExprStmt {
	return &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: Sel(I("h"), "errorHandler"),
			Args: []ast.Expr{
				I("w"),
				I("r"),
				Sel(I("http"), statusConst),
				msgExpr,
			},
		},
	}
}
