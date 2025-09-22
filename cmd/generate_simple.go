package main

import (
	"go/ast"
	"go/format"
	"go/token"
	"log"
	"os"
	"strings"

	"github.com/jolfzverb/codegen/internal/generator/astbuilder"
	"github.com/jolfzverb/codegen/internal/generator/options"
)

func main() {
	opts, err := options.GetOptions()
	if err != nil {
		log.Fatal("Failed to get options:", err)
	}

	log.Println("🚀 Generating code using new AST builder abstractions...")

	// Create AST builder
	config := astbuilder.BuilderConfig{
		PackageName:  "generated",
		ImportPrefix: opts.PackagePrefix,
		UsePointers:  opts.RequiredFieldsArePointers,
	}
	builder := astbuilder.NewBuilder(config)

	// Add imports
	builder.AddImport("net/http")
	builder.AddImport("github.com/go-chi/chi/v5")
	builder.AddImport("github.com/go-playground/validator/v10")

	// Add a simple User struct
	userStruct := &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent("User"),
				Type: &ast.StructType{
					Fields: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{ast.NewIdent("ID")},
								Type:  ast.NewIdent("int"),
								Tag:   &ast.BasicLit{Kind: token.STRING, Value: "`json:\"id\"`"},
							},
							{
								Names: []*ast.Ident{ast.NewIdent("Name")},
								Type:  ast.NewIdent("string"),
								Tag:   &ast.BasicLit{Kind: token.STRING, Value: "`json:\"name\"`"},
							},
							{
								Names: []*ast.Ident{ast.NewIdent("Email")},
								Type:  ast.NewIdent("string"),
								Tag:   &ast.BasicLit{Kind: token.STRING, Value: "`json:\"email\" validate:\"email\"`"},
							},
						},
					},
				},
			},
		},
	}
	builder.AddDeclaration(userStruct)

	// Add Handler struct
	handlerStruct := &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent("Handler"),
				Type: &ast.StructType{
					Fields: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{ast.NewIdent("validator")},
								Type:  &ast.StarExpr{X: ast.NewIdent("validator.Validate")},
							},
						},
					},
				},
			},
		},
	}
	builder.AddDeclaration(handlerStruct)

	// Add NewHandler function
	newHandlerFunc := &ast.FuncDecl{
		Name: ast.NewIdent("NewHandler"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("validator")},
						Type:  &ast.StarExpr{X: ast.NewIdent("validator.Validate")},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{X: ast.NewIdent("Handler")},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.UnaryExpr{
							Op: token.AND,
							X: &ast.CompositeLit{
								Type: ast.NewIdent("Handler"),
								Elts: []ast.Expr{
									&ast.KeyValueExpr{
										Key:   ast.NewIdent("validator"),
										Value: ast.NewIdent("validator"),
									},
								},
							},
						},
					},
				},
			},
		},
	}
	builder.AddDeclaration(newHandlerFunc)

	// Add GetUsers handler
	getUsersFunc := &ast.FuncDecl{
		Name: ast.NewIdent("GetUsers"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("w")},
						Type:  ast.NewIdent("http.ResponseWriter"),
					},
					{
						Names: []*ast.Ident{ast.NewIdent("r")},
						Type:  &ast.StarExpr{X: ast.NewIdent("http.Request")},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("w"),
							Sel: ast.NewIdent("WriteHeader"),
						},
						Args: []ast.Expr{
							&ast.BasicLit{Kind: token.INT, Value: "200"},
						},
					},
				},
			},
		},
	}
	builder.AddDeclaration(getUsersFunc)

	// Add AddRoutes function
	addRoutesFunc := &ast.FuncDecl{
		Name: ast.NewIdent("AddRoutes"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("h")},
						Type:  &ast.StarExpr{X: ast.NewIdent("Handler")},
					},
					{
						Names: []*ast.Ident{ast.NewIdent("r")},
						Type:  &ast.StarExpr{X: ast.NewIdent("chi.Mux")},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("r"),
							Sel: ast.NewIdent("GET"),
						},
						Args: []ast.Expr{
							&ast.BasicLit{Kind: token.STRING, Value: `"/users"`},
							ast.NewIdent("GetUsers"),
						},
					},
				},
			},
		},
	}
	builder.AddDeclaration(addRoutesFunc)

	// Build the AST file
	file := builder.BuildFile()
	if file == nil {
		log.Fatal("Failed to build AST file")
	}

	// Format the code
	var buf strings.Builder
	fset := token.NewFileSet()
	err = format.Node(&buf, fset, file)
	if err != nil {
		log.Fatal("Failed to format code:", err)
	}

	// Create output directory
	outputDir := "internal/generated/simple"
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		log.Fatal("Failed to create output directory:", err)
	}

	// Write the generated code
	outputFile := outputDir + "/generated.go"
	err = os.WriteFile(outputFile, []byte(buf.String()), 0644)
	if err != nil {
		log.Fatal("Failed to write output file:", err)
	}

	log.Printf("✅ Code generated successfully: %s", outputFile)
	log.Println("📊 Generated code using new AST builder abstractions!")
}
