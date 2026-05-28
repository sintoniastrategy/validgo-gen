package generator

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"io"
	"slices"
	"sort"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-faster/errors"
)

type HandlersFile struct {
	requiredFieldsArePointers bool
	packageName               *ast.Ident
	packageImports            []string
	packageAliasedImports     map[string]string // alias -> path
	interfaceDecls            []*ast.GenDecl

	handlerDecl            *ast.GenDecl
	handlerDeclQAFieldList *ast.FieldList // quick access to handler struct field list

	handlerConstructorDecl                       *ast.FuncDecl
	handlerConstructorDeclQAArgs                 *ast.FieldList    // quick access to handler constructor args
	handlerConstructorDeclQAConstructorComposite *ast.CompositeLit // quick access to handler struct initializer

	addRoutesDecl         *ast.FuncDecl
	handleDeclQASwitches  map[string]*ast.BlockStmt
	restDecls             []*ast.FuncDecl
	extraDecls            []ast.Decl
	hasContainsNullMethod bool
}

func (g *Generator) InitHandlerImports() {
	g.AddHandlersImport("github.com/go-playground/validator/v10")
	g.AddHandlersImport("github.com/go-chi/chi/v5")
}

func (g *Generator) InitHandlerStruct() {
	fieldList := &ast.FieldList{
		List: []*ast.Field{Field("validator", Star(Sel(I("validator"), "Validate")), "")},
	}
	handlerDecl := &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: I("Handler"),
				Type: &ast.StructType{
					Fields: fieldList,
				},
			},
		},
	}
	g.HandlersFile.handlerDecl = handlerDecl
	g.HandlersFile.handlerDeclQAFieldList = fieldList
}

func (g *Generator) InitHandlerConstructor() {
	initializerComposite := &ast.CompositeLit{
		Type: I("Handler"),
		Elts: []ast.Expr{
			&ast.KeyValueExpr{
				Key: I("validator"),
				Value: &ast.CallExpr{
					Fun: Sel(I("validator"), "New"),
					Args: []ast.Expr{
						&ast.CallExpr{
							Fun: Sel(I("validator"), "WithRequiredStructEnabled"),
						},
					},
				},
			},
		},
	}

	g.HandlersFile.handlerConstructorDecl = Func(
		"NewHandler",
		nil,
		nil,
		FieldA(Field("", Star(I("Handler")), "")),
		[]ast.Stmt{Ret1(Amp(initializerComposite))},
	)

	g.HandlersFile.handlerConstructorDeclQAArgs = g.HandlersFile.handlerConstructorDecl.Type.Params
	g.HandlersFile.handlerConstructorDeclQAConstructorComposite = initializerComposite
}

// FinalizeHandlerConstructor extends the Handler struct, constructor and
// initializer with the Option-3 error-handler plumbing: an `errorHandler`
// field, an `opts ...Option` variadic param, and a body that applies each
// option to the constructed *Handler. Called from GenerateHandlersFile only
// when the file has at least one route (no plumbing on schema-only packages
// like the def fixture).
func (g *Generator) FinalizeHandlerConstructor() {
	// 1. Append `errorHandler ErrorHandler` to the Handler struct.
	g.HandlersFile.handlerDeclQAFieldList.List = append(
		g.HandlersFile.handlerDeclQAFieldList.List,
		Field("errorHandler", I("ErrorHandler"), ""),
	)

	// 2. Append `errorHandler: DefaultErrorHandler` to the composite literal.
	g.HandlersFile.handlerConstructorDeclQAConstructorComposite.Elts = append(
		g.HandlersFile.handlerConstructorDeclQAConstructorComposite.Elts,
		&ast.KeyValueExpr{
			Key:   I("errorHandler"),
			Value: I("DefaultErrorHandler"),
		},
	)

	// 3. Append `opts ...Option` to the constructor params.
	g.HandlersFile.handlerConstructorDeclQAArgs.List = append(
		g.HandlersFile.handlerConstructorDeclQAArgs.List,
		Field("opts", &ast.Ellipsis{Elt: I("Option")}, ""),
	)

	// 4. Rewrite the body from `return &Handler{...}` to
	//        h := &Handler{...}
	//        for _, opt := range opts { opt(h) }
	//        return h
	initializer := g.HandlersFile.handlerConstructorDeclQAConstructorComposite
	g.HandlersFile.handlerConstructorDecl.Body.List = []ast.Stmt{
		&ast.AssignStmt{
			Lhs: []ast.Expr{I("h")},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{Amp(initializer)},
		},
		&ast.RangeStmt{
			Key:   I("_"),
			Value: I("opt"),
			Tok:   token.DEFINE,
			X:     I("opts"),
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ExprStmt{
						X: &ast.CallExpr{
							Fun:  I("opt"),
							Args: []ast.Expr{I("h")},
						},
					},
				},
			},
		},
		Ret1(I("h")),
	}
}

func (g *Generator) InitRoutesFunc() {
	g.HandlersFile.addRoutesDecl = Func(
		"AddRoutes",
		Field("h", Star(I("Handler")), ""),
		FieldA(Field("router", Sel(I("chi"), "Router"), "")),
		nil,
		[]ast.Stmt{},
	)
}

func (g *Generator) InitHandlerFields(packageName string) {
	g.HandlersFile.packageName = I(packageName)

	g.InitHandlerImports()

	g.InitHandlerStruct()

	g.InitHandlerConstructor()

	g.InitRoutesFunc()
}

func (g *Generator) NewHandlersFile() {
	g.HandlersFile = &HandlersFile{
		requiredFieldsArePointers: g.Opts.RequiredFieldsArePointers,
	}
}

func (g *Generator) WriteHandlersToOutput(output io.Writer) error {
	const op = "generator.HandlersFile.WriteToOutput"
	// go/ast package is great!
	_, err := output.Write([]byte("// Code generated by github.com/sintoniastrategy/validgo-gen; DO NOT EDIT.\n\n"))
	if err != nil {
		return errors.Wrap(err, op)
	}

	file := g.GenerateHandlersFile()
	err = format.Node(output, token.NewFileSet(), file)
	if err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}

func (g *Generator) AddHandlersInterface(name string, methodName string, requestName string, responseName string) {
	var methodParams []*ast.Field
	methodParams = append(methodParams, Field("ctx", Sel(I("context"), "Context"), ""))
	methodParams = append(methodParams, Field("r", Sel(I(g.GetCurrentModelsPackage()), requestName), ""))
	var methodResults []*ast.Field
	methodResults = append(methodResults, Field("", Star(Sel(I(g.GetCurrentModelsPackage()), responseName)), ""))
	methodResults = append(methodResults, Field("", I("error"), ""))
	g.HandlersFile.interfaceDecls = append(g.HandlersFile.interfaceDecls, &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: I(name),
				Type: &ast.InterfaceType{
					Methods: &ast.FieldList{
						List: []*ast.Field{{
							Names: []*ast.Ident{I(methodName)},
							Type: &ast.FuncType{
								Params: &ast.FieldList{
									List: methodParams,
								},
								Results: &ast.FieldList{
									List: methodResults,
								},
							},
						}},
					},
				},
			},
		},
	})
}

func (g *Generator) AddDependencyToHandlers(baseName string) {
	fieldName := GoIdentLowercase(baseName)

	g.HandlersFile.handlerDeclQAFieldList.List = append(g.HandlersFile.handlerDeclQAFieldList.List,
		Field(fieldName, I(baseName+"Handler"), ""))

	g.HandlersFile.handlerConstructorDeclQAArgs.List = append(g.HandlersFile.handlerConstructorDeclQAArgs.List,
		Field(fieldName, I(baseName+"Handler"), ""))

	g.HandlersFile.handlerConstructorDeclQAConstructorComposite.Elts = append(
		g.HandlersFile.handlerConstructorDeclQAConstructorComposite.Elts, &ast.KeyValueExpr{
			Key:   I(fieldName),
			Value: I(fieldName),
		},
	)
}

func (g *Generator) AddHandlersImport(path string) {
	if slices.Contains(g.HandlersFile.packageImports, path) {
		return
	}
	g.HandlersFile.packageImports = append(g.HandlersFile.packageImports, path)
}

func (g *Generator) AddHandlersImportWithAlias(alias, path string) {
	if g.HandlersFile.packageAliasedImports == nil {
		g.HandlersFile.packageAliasedImports = make(map[string]string)
	}
	if existing, ok := g.HandlersFile.packageAliasedImports[alias]; ok && existing == path {
		return
	}
	g.HandlersFile.packageAliasedImports[alias] = path
}

func (g *Generator) GenerateImportsSpecs(imp []string, aliased map[string]string) ([]*ast.ImportSpec, []ast.Spec) {
	type importEntry struct {
		alias string
		path  string
	}
	classify := func(path string) (system, lib, mine bool) {
		if strings.HasPrefix(path, g.Opts.PackagePrefix) {
			return false, false, true
		}
		prefix := strings.SplitN(path, "/", 2)[0] //nolint:mnd
		if strings.Contains(prefix, ".") {
			return false, true, false
		}
		return true, false, false
	}

	var systemImports []importEntry
	var libImports []importEntry
	var myImports []importEntry
	for _, path := range imp {
		sys, lib, mine := classify(path)
		switch {
		case sys:
			systemImports = append(systemImports, importEntry{path: path})
		case lib:
			libImports = append(libImports, importEntry{path: path})
		case mine:
			myImports = append(myImports, importEntry{path: path})
		}
	}
	aliases := make([]string, 0, len(aliased))
	for alias := range aliased {
		aliases = append(aliases, alias)
	}
	sort.Strings(aliases)
	for _, alias := range aliases {
		path := aliased[alias]
		sys, lib, mine := classify(path)
		switch {
		case sys:
			systemImports = append(systemImports, importEntry{alias: alias, path: path})
		case lib:
			libImports = append(libImports, importEntry{alias: alias, path: path})
		case mine:
			myImports = append(myImports, importEntry{alias: alias, path: path})
		}
	}

	sortEntries := func(s []importEntry) {
		sort.Slice(s, func(i, j int) bool { return s[i].path < s[j].path })
	}
	sortEntries(systemImports)
	sortEntries(libImports)
	sortEntries(myImports)

	makeSpec := func(e importEntry) *ast.ImportSpec {
		spec := &ast.ImportSpec{Path: Str(e.path)}
		if e.alias != "" {
			spec.Name = I(e.alias)
		}
		return spec
	}

	specs := make([]*ast.ImportSpec, 0, len(imp)+len(aliased))
	for _, e := range systemImports {
		specs = append(specs, makeSpec(e))
	}
	// Add a space to separate system and library imports
	// but go/ast is too great for that
	for _, e := range libImports {
		specs = append(specs, makeSpec(e))
	}
	// Add a space to separate library and user imports
	// but go/ast is too great for that
	for _, e := range myImports {
		specs = append(specs, makeSpec(e))
	}

	declSpecs := make([]ast.Spec, 0, len(specs))
	for _, spec := range specs {
		declSpecs = append(declSpecs, spec)
	}

	return specs, declSpecs
}

func (g *Generator) GenerateHandlersFile() *ast.File {
	if len(g.HandlersFile.addRoutesDecl.Body.List) > 0 {
		g.FinalizeHandlerConstructor()
		g.AddStandardErrorDecls()
	}

	importSpecs, declSpecs := g.GenerateImportsSpecs(g.HandlersFile.packageImports, g.HandlersFile.packageAliasedImports)

	g.FinalizeHandlerSwitches()

	file := &ast.File{
		Name:    g.HandlersFile.packageName,
		Decls:   []ast.Decl{},
		Imports: importSpecs,
	}

	file.Decls = append(file.Decls, &ast.GenDecl{
		Tok:   token.IMPORT,
		Specs: declSpecs,
	})
	for _, d := range g.HandlersFile.interfaceDecls {
		file.Decls = append(file.Decls, d)
	}

	file.Decls = append(file.Decls, g.HandlersFile.handlerDecl)
	file.Decls = append(file.Decls, g.HandlersFile.handlerConstructorDecl)
	file.Decls = append(file.Decls, g.HandlersFile.addRoutesDecl)
	for _, d := range g.HandlersFile.restDecls {
		file.Decls = append(file.Decls, d)
	}
	file.Decls = append(file.Decls, g.HandlersFile.extraDecls...)

	return file
}

func (g *Generator) AddRouteToRouter(baseName string, method string, pathName string) {
	g.HandlersFile.addRoutesDecl.Body.List = append(g.HandlersFile.addRoutesDecl.Body.List, &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: Sel(I("router"), method),
			Args: []ast.Expr{
				Str(pathName),
				Sel(I("h"), "handle"+baseName),
			},
		},
	})
}

func (g *Generator) GetHandler(baseName string) *ast.BlockStmt {
	if g.HandlersFile.handleDeclQASwitches == nil {
		return nil
	}
	if blockStmt, ok := g.HandlersFile.handleDeclQASwitches[baseName]; ok {
		return blockStmt
	}

	return nil
}

func (g *Generator) CreateHandler(baseName string) {
	g.AddHandlersImport("mime")

	switchBody := &ast.BlockStmt{
		List: []ast.Stmt{},
	}

	handleFunc := Func(
		"handle"+baseName,
		Field("h", Star(I("Handler")), ""),
		[]*ast.Field{
			Field("w", Sel(I("http"), "ResponseWriter"), ""),
			Field("r", Star(Sel(I("http"), "Request")), ""),
		},
		nil,
		[]ast.Stmt{
			// contentType, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
			&ast.AssignStmt{
				Lhs: []ast.Expr{
					I("contentType"),
					I("_"),
					I("_"),
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: Sel(I("mime"), "ParseMediaType"),
						Args: []ast.Expr{
							&ast.CallExpr{
								Fun:  Sel(I("r.Header"), "Get"),
								Args: []ast.Expr{Str("Content-Type")},
							},
						},
					},
				},
			},
			&ast.SwitchStmt{
				Tag:  I("contentType"),
				Body: switchBody,
			},
		},
	)

	g.HandlersFile.restDecls = append(g.HandlersFile.restDecls, handleFunc)

	if g.HandlersFile.handleDeclQASwitches == nil {
		g.HandlersFile.handleDeclQASwitches = make(map[string]*ast.BlockStmt)
	}
	g.HandlersFile.handleDeclQASwitches[baseName] = switchBody
}

// CreateDirectHandler generates a handler that directly delegates to the request handler
// without checking Content-Type. Used for operations with no request body (e.g. GET, DELETE)
// where Content-Type is irrelevant.
func (g *Generator) CreateDirectHandler(baseName string) {
	handleFunc := Func(
		"handle"+baseName,
		Field("h", Star(I("Handler")), ""),
		[]*ast.Field{
			Field("w", Sel(I("http"), "ResponseWriter"), ""),
			Field("r", Star(Sel(I("http"), "Request")), ""),
		},
		nil,
		[]ast.Stmt{
			&ast.ExprStmt{
				X: &ast.CallExpr{
					Fun:  Sel(I("h"), "handle"+baseName+"Request"),
					Args: []ast.Expr{I("w"), I("r")},
				},
			},
		},
	)

	g.HandlersFile.restDecls = append(g.HandlersFile.restDecls, handleFunc)
}

func (g *Generator) FinalizeHandlerSwitches() {
	if g.HandlersFile.handleDeclQASwitches == nil {
		return
	}
	for _, blockStmt := range g.HandlersFile.handleDeclQASwitches {
		blockStmt.List = append(blockStmt.List, &ast.CaseClause{
			List: nil,
			Body: []ast.Stmt{
				writeStandardErrorCall("StatusUnsupportedMediaType", Str("Unsupported Content-Type")),
				Ret(),
			},
		})
	}
}

func (g *Generator) AddContentTypeHandler(baseName string, rawContentType string) {
	if g.HandlersFile.handleDeclQASwitches == nil {
		return
	}
	if blockStmt, ok := g.HandlersFile.handleDeclQASwitches[baseName]; ok {
		stmts := []ast.Stmt{
			&ast.ExprStmt{
				X: &ast.CallExpr{
					Fun: Sel(I("h"), "handle"+baseName+"Request"),
					Args: []ast.Expr{
						I("w"),
						I("r"),
					},
				},
			},
			Ret(),
		}

		blockStmt.List = append(blockStmt.List, &ast.CaseClause{
			List: []ast.Expr{Str(rawContentType)},
			Body: stmts,
		},
		)

		if rawContentType == applicationJSONCT {
			blockStmt.List = append(blockStmt.List, &ast.CaseClause{
				List: []ast.Expr{Str("")},
				Body: stmts,
			})
		}
	}
}

func (g *Generator) AddHandleOperationMethodHandlers(baseName string) {
	g.HandlersFile.restDecls = append(g.HandlersFile.restDecls, Func(
		"handle"+baseName+"Request",
		Field("h", Star(I("Handler")), ""),
		[]*ast.Field{
			Field("w", Sel(I("http"), "ResponseWriter"), ""),
			Field("r", Star(Sel(I("http"), "Request")), ""),
		},
		nil,
		[]ast.Stmt{
			&ast.AssignStmt{
				Lhs: []ast.Expr{
					I("request"),
					I("err"),
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: Sel(I("h"), "parse"+baseName+"Request"),
						Args: []ast.Expr{
							I("r"),
						},
					},
				},
			},
			&ast.IfStmt{
				Cond: Ne(I("err"), I("nil")),
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						writeStandardErrorCall("StatusBadRequest", &ast.CallExpr{
							Fun: Sel(I("err"), "Error"),
						}),
						Ret(),
					},
				},
			},
			&ast.AssignStmt{
				Lhs: []ast.Expr{
					I("ctx"),
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun:  Sel(I("r"), "Context"),
						Args: []ast.Expr{},
					},
				},
			},
			&ast.AssignStmt{
				Lhs: []ast.Expr{
					I("response"),
					I("err"),
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: Sel(Sel(I("h"), GoIdentLowercase(baseName)), "Handle"+baseName),
						Args: []ast.Expr{
							I("ctx"),
							Star(I("request")),
						},
					},
				},
			},
			&ast.IfStmt{
				Cond: &ast.BinaryExpr{
					X:  Ne(I("err"), I("nil")),
					Op: token.LOR,
					Y:  Eq(I("response"), I("nil")),
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						writeStandardErrorCall("StatusInternalServerError", Str("Internal server error")),
						Ret(),
					},
				},
			},
			&ast.ExprStmt{
				X: &ast.CallExpr{
					Fun: Sel(I("h"), "write"+baseName+"Response"),
					Args: []ast.Expr{
						I("w"),
						I("r"),
						I("response"),
					},
				},
			},
			Ret(),
		},
	))
}

func (g *Generator) AddWriteResponseMethodHandlers(baseName string, codes []string, operation *openapi3.Operation) error {
	switchBody := &ast.BlockStmt{
		List: []ast.Stmt{},
	}
	for _, code := range codes {
		response := operation.Responses.Value(code)

		caseBody := []ast.Stmt{}
		caseBody = append(caseBody, &ast.IfStmt{
			Cond: Eq(Sel(I("response"), "Response"+code), I("nil")),
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					writeStandardErrorCall("StatusInternalServerError", Str("Internal server error")),
					Ret(),
				},
			},
		})

		if len(response.Value.Headers) > 0 {
			caseBody = append(caseBody,
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: Sel(I("h"), "write"+baseName+code+"ResponseHeaders"),
						Args: []ast.Expr{
							I("w"),
							I("r"),
							Sel(I("response"), "Response"+code),
						},
					},
				})
		}

		if len(response.Value.Content) > 0 {
			if len(response.Value.Content) > 1 {
				return errors.New("multiple content types are not supported for response code " + code)
			}
			var contentType string
			for key := range response.Value.Content {
				contentType = key
				break
			}
			caseBody = append(caseBody,
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: Sel(&ast.CallExpr{
							Fun:  Sel(I("w"), "Header"),
							Args: []ast.Expr{},
						}, "Set"),
						Args: []ast.Expr{
							Str("Content-Type"),
							Str(g.getContentTypeHeadeValue(contentType)),
						},
					},
				},
			)

		}

		caseBody = append(caseBody, &ast.ExprStmt{
			X: &ast.CallExpr{
				Fun:  Sel(I("w"), "WriteHeader"),
				Args: []ast.Expr{Sel(I("response"), "StatusCode")},
			},
		})
		caseBody = append(caseBody, &ast.ExprStmt{
			X: &ast.CallExpr{
				Fun: Sel(I("h"), "write"+baseName+code+"Response"),
				Args: []ast.Expr{
					I("w"),
					I("r"),
					Sel(I("response"), "Response"+code),
				},
			},
		})
		caseBody = append(caseBody, &ast.ReturnStmt{})
		switchBody.List = append(switchBody.List, &ast.CaseClause{
			List: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.INT,
					Value: code,
				},
			},
			Body: caseBody,
		})
	}

	writeResponseFunc := Func(
		"write"+baseName+"Response",
		Field("h", Star(I("Handler")), ""),
		[]*ast.Field{
			Field("w", Sel(I("http"), "ResponseWriter"), ""),
			Field("r", Star(Sel(I("http"), "Request")), ""),
			Field("response", Star(Sel(I(g.GetCurrentModelsPackage()), baseName+"Response")), ""),
		},
		nil,
		[]ast.Stmt{
			&ast.SwitchStmt{
				Tag:  Sel(I("response"), "StatusCode"),
				Body: switchBody,
			},
			writeStandardErrorCall("StatusInternalServerError", Str("Internal server error")),
		},
	)

	g.HandlersFile.restDecls = append(g.HandlersFile.restDecls, writeResponseFunc)
	return nil
}

func (g *Generator) getContentTypeHeadeValue(contentType string) string {
	textualContentType := map[string]struct{}{
		"text/plain":             {},
		"text/html":              {},
		"text/css":               {},
		"application/javascript": {},
		"application/xml":        {},
		"application/json":       {},
	}
	if _, ok := textualContentType[contentType]; ok {
		return fmt.Sprintf("%s; charset=utf-8", contentType)
	}
	return contentType
}

func (g *Generator) AddWriteHeadersForResponseCode(baseName string, code string, response *openapi3.ResponseRef) error {
	var body []ast.Stmt

	g.AddHandlersImport("encoding/json")
	body = append(body, &ast.AssignStmt{
		Lhs: []ast.Expr{
			I("headersJSON"),
			I("err"),
		},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun:  Sel(I("json"), "Marshal"),
				Args: []ast.Expr{Sel(I("resp"), "Headers")},
			},
		},
	})
	body = append(body, &ast.IfStmt{
		Cond: Ne(I("err"), I("nil")),
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				writeStandardErrorCall("StatusInternalServerError", Str("Internal server error")),
				Ret(),
			},
		},
	})
	body = append(body, &ast.DeclStmt{
		Decl: &ast.GenDecl{
			Tok: token.VAR,
			Specs: []ast.Spec{
				&ast.ValueSpec{
					Names: []*ast.Ident{I("headers")},
					Type: &ast.MapType{
						Key:   I("string"),
						Value: I("string"),
					},
				},
			},
		},
	})
	body = append(body, &ast.AssignStmt{
		Lhs: []ast.Expr{I("err")},
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: Sel(I("json"), "Unmarshal"),
				Args: []ast.Expr{
					I("headersJSON"),
					Amp(I("headers")),
				},
			},
		},
	})
	body = append(body, &ast.IfStmt{
		Cond: Ne(I("err"), I("nil")),
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				writeStandardErrorCall("StatusInternalServerError", Str("Internal server error")),
				Ret(),
			},
		},
	})
	body = append(body, &ast.RangeStmt{
		Key:   I("key"),
		Value: I("value"),
		Tok:   token.DEFINE,
		X:     I("headers"),
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: Sel(&ast.CallExpr{
							Fun:  Sel(I("w"), "Header"),
							Args: []ast.Expr{},
						}, "Set"),
						Args: []ast.Expr{
							I("key"),
							I("value"),
						},
					},
				},
			},
		},
	})

	writeResponseFunc := Func(
		"write"+baseName+code+"ResponseHeaders",
		Field("h", Star(I("Handler")), ""),
		[]*ast.Field{
			Field("w", Sel(I("http"), "ResponseWriter"), ""),
			Field("r", Star(Sel(I("http"), "Request")), ""),
			Field("resp", Star(Sel(I(g.GetCurrentModelsPackage()), baseName+"Response"+code)), ""),
		},
		nil,
		body,
	)

	g.HandlersFile.restDecls = append(g.HandlersFile.restDecls, writeResponseFunc)

	return nil
}

func (g *Generator) AddWriteResponseCode(baseName string, code string, response *openapi3.ResponseRef) error {
	var body []ast.Stmt

	if len(response.Value.Content) > 1 {
		return errors.New("multiple responses are not supported")
	}
	for key, value := range response.Value.Content {
		if key != applicationJSONCT {
			return errors.New("only application/json content type is supported")
		}
		if value.Schema != nil {
			g.AddHandlersImport("encoding/json")
			body = append(body, &ast.AssignStmt{
				Lhs: []ast.Expr{I("err")},
				Tok: token.ASSIGN,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: Sel(&ast.CallExpr{
							Fun:  Sel(I("json"), "NewEncoder"),
							Args: []ast.Expr{I("w")},
						}, "Encode"),

						Args: []ast.Expr{Sel(I("resp"), "Body")},
					},
				},
			})
			body = append(body, &ast.IfStmt{
				Cond: Ne(I("err"), I("nil")),
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						writeStandardErrorCall("StatusInternalServerError", Str("Internal server error")),
						Ret(),
					},
				},
			})
		}
	}

	if len(body) > 0 {
		body = append([]ast.Stmt{&ast.DeclStmt{
			Decl: &ast.GenDecl{
				Tok: token.VAR,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: []*ast.Ident{I("err")},
						Type:  I("error"),
					},
				},
			},
		}}, body...)
	}

	writeResponseFunc := Func(
		"write"+baseName+code+"Response",
		Field("h", Star(I("Handler")), ""),
		[]*ast.Field{
			Field("w", Sel(I("http"), "ResponseWriter"), ""),
			Field("r", Star(Sel(I("http"), "Request")), ""),
			Field("resp", Star(Sel(I(g.GetCurrentModelsPackage()), baseName+"Response"+code)), ""),
		},
		nil,
		body,
	)

	g.HandlersFile.restDecls = append(g.HandlersFile.restDecls, writeResponseFunc)

	return nil
}

func (g *Generator) AddParsePathParamsMethod(baseName string, params openapi3.Parameters) error {
	bodyList := []ast.Stmt{
		&ast.DeclStmt{
			Decl: &ast.GenDecl{
				Tok: token.VAR,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: []*ast.Ident{I("pathParams")},
						Type:  Sel(I(g.GetCurrentModelsPackage()), baseName+"PathParams"),
					},
				},
			},
		},
	}

	for _, param := range params {
		if param.Value.Schema == nil || param.Value.Schema.Value == nil {
			continue
		}

		varName := GoIdentLowercase(FormatGoLikeIdentifier(param.Value.Name))
		bodyList = append(bodyList, &ast.AssignStmt{
			Lhs: []ast.Expr{I(varName)},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun:  Sel(I("chi"), "URLParam"),
					Args: []ast.Expr{I("r"), Str(param.Value.Name)},
				},
			},
		})
		bodyList = append(bodyList, &ast.IfStmt{
			Cond: Eq(I(varName), Str("")),
			Body: &ast.BlockStmt{
				List: []ast.Stmt{Ret2(I("nil"),
					&ast.CallExpr{
						Fun:  Sel(I("errors"), "New"),
						Args: []ast.Expr{Str(param.Value.Name + " path param is required")},
					},
				)},
			},
		})
		g.AddHandlersImport("github.com/go-faster/errors")
		switch {
		case param.Value.Schema.Value.Type.Permits("string"):
			bodyList = append(bodyList, &ast.AssignStmt{
				Lhs: []ast.Expr{Sel(I("pathParams"), FormatGoLikeIdentifier(param.Value.Name))},
				Tok: token.ASSIGN,
				Rhs: []ast.Expr{
					I(varName),
				},
			})
		default:
			return errors.New(fmt.Sprintf("unsupported path parameter type: %v", param.Value.Schema.Value.Type)) //nolint:revive
		}
	}

	bodyList = append(bodyList, &ast.AssignStmt{
		Lhs: []ast.Expr{I("err")},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: Sel(Sel(I("h"), "validator"), "Struct"),
				Args: []ast.Expr{
					I("pathParams"),
				},
			},
		},
	})
	bodyList = append(bodyList, &ast.IfStmt{
		Cond: Ne(I("err"), I("nil")),
		Body: &ast.BlockStmt{List: []ast.Stmt{Ret2(I("nil"), I("err"))}},
	})
	bodyList = append(bodyList, Ret2(Amp(I("pathParams")), I("nil")))

	parsePathParamsFunc := Func(
		"parse"+baseName+"PathParams",
		Field("h", Star(I("Handler")), ""),
		[]*ast.Field{
			Field("r", Star(Sel(I("http"), "Request")), ""),
		},
		[]*ast.Field{
			Field("", Star(Sel(I(g.GetCurrentModelsPackage()), baseName+"PathParams")), ""),
			Field("", I("error"), ""),
		},
		bodyList,
	)

	g.HandlersFile.restDecls = append(g.HandlersFile.restDecls, parsePathParamsFunc)

	return nil
}
