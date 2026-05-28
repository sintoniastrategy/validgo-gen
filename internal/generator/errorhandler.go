package generator

import (
	"go/ast"
	"go/parser"
	"go/token"
)

const standardErrorHandlerSrc = `package _

type ErrorHandler = func(w http.ResponseWriter, r *http.Request, status int, msg string)

type Option func(*Handler)

func WithErrorHandler(eh ErrorHandler) Option {
	return func(h *Handler) { h.errorHandler = eh }
}

func (h *Handler) SetErrorHandler(eh ErrorHandler) { h.errorHandler = eh }

var DefaultErrorHandler ErrorHandler = func(w http.ResponseWriter, r *http.Request, status int, msg string) {
	http.Error(w, fmt.Sprintf("{\"error\":%s}", strconv.Quote(msg)), status)
}
`

func parseStandardErrorHandlerDecls() []ast.Decl {
	file, err := parser.ParseFile(token.NewFileSet(), "", standardErrorHandlerSrc, 0)
	if err != nil {
		panic(err)
	}
	return file.Decls
}

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
	g.HandlersFile.extraDecls = append(g.HandlersFile.extraDecls, parseStandardErrorHandlerDecls()...)
}

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
