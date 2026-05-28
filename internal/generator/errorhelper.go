package generator

import (
	"go/ast"
	"go/parser"
	"go/token"
)

// standardErrorHelperSrc is the source for the per-package error-handler
// hook every generated handlers.go file embeds. It declares:
//
//   - ErrorHandler — function-type alias matching every codegen error site.
//   - Option / WithErrorHandler — variadic option for NewHandler.
//   - (*Handler).SetErrorHandler — setter for the aggregator-loop pattern
//     (one Handler.SetErrorHandler call per generated package, instead of
//     threading WithErrorHandler through every per-package constructor).
//   - DefaultErrorHandler — the legacy {"error":"<msg>"} body emitted via
//     net/http.Error. Exposed so consumers can wrap it from a custom
//     ErrorHandler.
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

// WithErrorHandler replaces DefaultErrorHandler with the supplied function
// for this Handler. Pass it to NewHandler.
func WithErrorHandler(eh ErrorHandler) Option {
	return func(h *Handler) { h.errorHandler = eh }
}

// SetErrorHandler is the post-construction equivalent of WithErrorHandler.
// Useful when Handlers are constructed by an injector (fx, wire, ...) and
// the error handler is configured later.
func (h *Handler) SetErrorHandler(eh ErrorHandler) { h.errorHandler = eh }

// DefaultErrorHandler writes a text/plain body shaped as the JSON literal
// {"error":"<msg>"} via net/http.Error. This is the same body the
// generator emitted before the ErrorHandler hook was introduced — consumers
// that don't override it see no behaviour change.
var DefaultErrorHandler ErrorHandler = func(w http.ResponseWriter, r *http.Request, status int, msg string) {
	http.Error(w, fmt.Sprintf("{\"error\":%s}", strconv.Quote(msg)), status)
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
	g.AddHandlersImport("fmt")
	g.AddHandlersImport("net/http")
	g.AddHandlersImport("strconv")
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
