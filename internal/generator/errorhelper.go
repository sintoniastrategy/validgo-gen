package generator

import (
	"go/ast"
	"go/parser"
	"go/token"
)

// standardErrorHelperSrc is the source for the per-package error-envelope
// helper that every generated handlers.go file embeds. The helper writes a
// JSON body matching the standard error envelope: {"code","error","req_id"}.
//
// The "package _" prefix is replaced with the generated package's name when
// these decls are appended to the file by GenerateHandlersFile.
const standardErrorHelperSrc = `package _

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

func writeStandardError(w http.ResponseWriter, r *http.Request, status int, msg string) {
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

// AddStandardErrorDecls appends the package-level statusToCode map and
// writeStandardError function to the generated handlers file. Safe to call
// more than once per file — only the first call adds the decls.
func (g *Generator) AddStandardErrorDecls() {
	for _, d := range g.HandlersFile.extraDecls {
		if fn, ok := d.(*ast.FuncDecl); ok && fn.Name != nil && fn.Name.Name == "writeStandardError" {
			return
		}
	}
	g.AddHandlersImport("encoding/json")
	g.AddHandlersImport("net/http")
	g.AddHandlersImportWithAlias("chimw", "github.com/go-chi/chi/v5/middleware")
	g.HandlersFile.extraDecls = append(g.HandlersFile.extraDecls, parseStandardErrorHelperDecls()...)
}

// writeStandardErrorCall builds the AST for
//
//	writeStandardError(w, r, http.<StatusConst>, <msgExpr>)
func writeStandardErrorCall(statusConst string, msgExpr ast.Expr) *ast.ExprStmt {
	return &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: I("writeStandardError"),
			Args: []ast.Expr{
				I("w"),
				I("r"),
				Sel(I("http"), statusConst),
				msgExpr,
			},
		},
	}
}
