package generator

import (
	"fmt"
	"go/ast"
	"go/token"
	"sort"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-faster/errors"
)

func (g *Generator) AddParseQueryParamsMethod(baseName string, params openapi3.Parameters) error {
	bodyList := []ast.Stmt{
		&ast.DeclStmt{
			Decl: &ast.GenDecl{
				Tok: token.VAR,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: []*ast.Ident{I("queryParams")},
						Type:  Sel(I(g.GetCurrentModelsPackage()), baseName+"QueryParams"),
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
					Fun: Sel(&ast.CallExpr{
						Fun:  Sel(Sel(I("r"), "URL"), "Query"),
						Args: []ast.Expr{},
					}, "Get"),
					Args: []ast.Expr{Str(param.Value.Name)},
				},
			},
		})
		if param.Value.Required {
			bodyList = append(bodyList, &ast.IfStmt{
				Cond: Eq(I(varName), Str("")),
				Body: &ast.BlockStmt{
					List: []ast.Stmt{Ret2(I("nil"),
						&ast.CallExpr{
							Fun: Sel(I("errors"), "New"),
							Args: []ast.Expr{
								Str(param.Value.Name + " query param is required"),
							},
						},
					)},
				},
			})
			g.AddHandlersImport("github.com/go-faster/errors")
			switch {
			case param.Value.Schema.Value.Type.Permits("string"):
				bodyList = append(bodyList,
					g.AssignStringField("queryParams", varName, FormatGoLikeIdentifier(param.Value.Name), param.Value.Schema, param.Value.Required)...,
				)
			default:
				return errors.New(fmt.Sprintf("unsupported path parameter type: %v", param.Value.Schema.Value.Type)) //nolint:revive
			}
		} else {
			bodyList = append(bodyList, &ast.IfStmt{
				Cond: Ne(I(varName), Str("")),
				Body: &ast.BlockStmt{
					List: g.AssignStringField("queryParams", varName, FormatGoLikeIdentifier(param.Value.Name), param.Value.Schema, param.Value.Required),
				},
			})
		}
	}
	bodyList = append(bodyList, &ast.AssignStmt{
		Lhs: []ast.Expr{I("err")},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: Sel(Sel(I("h"), "validator"), "Struct"),
				Args: []ast.Expr{
					I("queryParams"),
				},
			},
		},
	})
	bodyList = append(bodyList, &ast.IfStmt{
		Cond: Ne(I("err"), I("nil")),
		Body: &ast.BlockStmt{List: []ast.Stmt{Ret2(I("nil"), I("err"))}},
	})

	bodyList = append(bodyList, Ret2(Amp(I("queryParams")), I("nil")))

	g.HandlersFile.restDecls = append(g.HandlersFile.restDecls, Func("parse"+baseName+"QueryParams",
		Field("h", Star(I("Handler")), ""),
		[]*ast.Field{
			Field("r", Star(Sel(I("http"), "Request")), ""),
		},
		[]*ast.Field{
			Field("", Star(Sel(I(g.GetCurrentModelsPackage()), baseName+"QueryParams")), ""),
			Field("", I("error"), ""),
		},
		bodyList,
	))

	return nil
}

func (g *Generator) AssignStringField(paramsName string, varName string, fieldName string, param *openapi3.SchemaRef, required bool) []ast.Stmt {
	if param.Value.Format == "date-time" {
		g.AddHandlersImport("time")
		var result []ast.Stmt
		result = append(result, &ast.AssignStmt{
			Lhs: []ast.Expr{
				I("parsed" + fieldName),
				I("err"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: Sel(I("time"), "Parse"),
					Args: []ast.Expr{
						Sel(I("time"), "RFC3339"),
						I(varName),
					},
				},
			},
		})
		result = append(result, &ast.IfStmt{
			Cond: Ne(I("err"), I("nil")),
			Body: &ast.BlockStmt{
				List: []ast.Stmt{Ret2(
					I("nil"),
					&ast.CallExpr{
						Fun: Sel(I("errors"), "Wrap"),
						Args: []ast.Expr{
							I("err"),
							Str(fieldName + " is not a valid date-time format"),
						},
					},
				)},
			},
		})
		var rhs ast.Expr
		if required && !g.HandlersFile.requiredFieldsArePointers {
			rhs = I("parsed" + fieldName)
		} else {
			rhs = Amp(I("parsed" + fieldName))
		}

		return append(result, &ast.AssignStmt{
			Lhs: []ast.Expr{Sel(I(paramsName), fieldName)},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{rhs},
		})
	}
	var rhs ast.Expr
	if required && !g.HandlersFile.requiredFieldsArePointers {
		rhs = I(varName)
	} else {
		rhs = Amp(I(varName))
	}

	return []ast.Stmt{&ast.AssignStmt{
		Lhs: []ast.Expr{Sel(I(paramsName), fieldName)},
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{rhs},
	}}
}

func (g *Generator) AddParseHeadersMethod(baseName string, params openapi3.Parameters) error {
	bodyList := []ast.Stmt{
		&ast.DeclStmt{
			Decl: &ast.GenDecl{
				Tok: token.VAR,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: []*ast.Ident{I("headers")},
						Type:  Sel(I(g.GetCurrentModelsPackage()), baseName+"Headers"),
					},
				},
			},
		},
	}
	for _, param := range params {
		if param.Value.Schema == nil || param.Value.Schema.Value == nil {
			continue
		}
		if g.Opts.AllowRemoteAddrParam && param.Value.Name == "Remote-Addr" && param.Value.Schema.Value.Format == "remote-addr" {
			bodyList = append(bodyList, &ast.AssignStmt{
				Lhs: []ast.Expr{Sel(I("headers"), FormatGoLikeIdentifier(param.Value.Name))},
				Tok: token.ASSIGN,
				Rhs: []ast.Expr{Sel(I("r"), "RemoteAddr")},
			})
			continue
		}
		varName := GoIdentLowercase(FormatGoLikeIdentifier(param.Value.Name))
		bodyList = append(bodyList, &ast.AssignStmt{
			Lhs: []ast.Expr{I(varName)},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun:  Sel(Sel(I("r"), "Header"), "Get"),
					Args: []ast.Expr{Str(param.Value.Name)},
				},
			},
		})
		if param.Value.Required {
			bodyList = append(bodyList, &ast.IfStmt{
				Cond: Eq(I(varName), Str("")),
				Body: &ast.BlockStmt{
					List: []ast.Stmt{Ret2(I("nil"),
						&ast.CallExpr{
							Fun: Sel(I("errors"), "New"),
							Args: []ast.Expr{
								Str(param.Value.Name + " header is required"),
							},
						},
					)},
				},
			})
			g.AddHandlersImport("github.com/go-faster/errors")
			switch {
			case param.Value.Schema.Value.Type.Permits("string"):
				bodyList = append(bodyList,
					g.AssignStringField("headers", varName, FormatGoLikeIdentifier(param.Value.Name),
						param.Value.Schema, param.Value.Required,
					)...,
				)
			default:
				return errors.New("unsupported path parameter type: " + fmt.Sprint(param.Value.Schema.Value.Type))
			}
		} else {
			bodyList = append(bodyList, &ast.IfStmt{
				Cond: Ne(I(varName), Str("")),
				Body: &ast.BlockStmt{
					List: g.AssignStringField("headers", varName, FormatGoLikeIdentifier(param.Value.Name),
						param.Value.Schema, param.Value.Required,
					),
				},
			})
		}
	}
	bodyList = append(bodyList, &ast.AssignStmt{
		Lhs: []ast.Expr{I("err")},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: Sel(Sel(I("h"), "validator"), "Struct"),
				Args: []ast.Expr{
					I("headers"),
				},
			},
		},
	})
	bodyList = append(bodyList, &ast.IfStmt{
		Cond: Ne(I("err"), I("nil")),
		Body: &ast.BlockStmt{List: []ast.Stmt{Ret2(I("nil"), I("err"))}},
	})
	bodyList = append(bodyList, Ret2(Amp(I("headers")), I("nil")))
	g.HandlersFile.restDecls = append(g.HandlersFile.restDecls, Func("parse"+baseName+"Headers",
		Field("h", Star(I("Handler")), ""),
		[]*ast.Field{
			Field("r", Star(Sel(I("http"), "Request")), ""),
		},
		[]*ast.Field{
			Field("", Star(Sel(I(g.GetCurrentModelsPackage()), baseName+"Headers")), ""),
			Field("", I("error"), ""),
		},
		bodyList,
	))

	return nil
}

func (g *Generator) AddParseCookiesMethod(baseName string, params openapi3.Parameters) error {
	bodyList := []ast.Stmt{
		&ast.DeclStmt{
			Decl: &ast.GenDecl{
				Tok: token.VAR,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: []*ast.Ident{I("cookies")},
						Type:  Sel(I(g.GetCurrentModelsPackage()), baseName+"Cookies"),
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
			Lhs: []ast.Expr{I(varName), I("err")},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun:  Sel(I("r"), "Cookie"),
					Args: []ast.Expr{Str(param.Value.Name)},
				},
			},
		})

		if param.Value.Required {
			bodyList = append(bodyList, &ast.IfStmt{
				Cond: Ne(I("err"), I("nil")),
				Body: &ast.BlockStmt{List: []ast.Stmt{Ret2(I("nil"), I("err"))}},
			})
		} else {
			bodyList = append(bodyList, &ast.IfStmt{
				Cond: &ast.BinaryExpr{
					X:  Ne(I("err"), I("nil")),
					Op: token.LAND,
					Y: &ast.UnaryExpr{
						Op: token.NOT,
						X: &ast.CallExpr{
							Fun: Sel(I("errors"), "Is"),
							Args: []ast.Expr{
								I("err"),
								Sel(I("http"), "ErrNoCookie"),
							},
						},
					},
				},
				Body: &ast.BlockStmt{List: []ast.Stmt{Ret2(I("nil"), I("err"))}},
			})
			g.AddHandlersImport("github.com/go-faster/errors")
		}

		if param.Value.Required {
			bodyList = append(bodyList, &ast.AssignStmt{
				Lhs: []ast.Expr{I(varName + "Value")},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{Sel(I(varName), "Value")},
			})

			switch {
			case param.Value.Schema.Value.Type.Permits("string"):
				bodyList = append(bodyList,
					g.AssignStringField("cookies", varName+"Value", FormatGoLikeIdentifier(param.Value.Name),
						param.Value.Schema, param.Value.Required,
					)...,
				)
			default:
				return errors.New("unsupported path parameter type: " + fmt.Sprint(param.Value.Schema.Value.Type))
			}
		} else {
			ifBody := []ast.Stmt{&ast.AssignStmt{
				Lhs: []ast.Expr{I(varName + "Value")},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{Sel(I(varName), "Value")},
			}}
			ifBody = append(ifBody,
				g.AssignStringField("cookies", varName+"Value", FormatGoLikeIdentifier(param.Value.Name),
					param.Value.Schema, param.Value.Required,
				)...,
			)
			bodyList = append(bodyList, &ast.IfStmt{
				Cond: Eq(I("err"), I("nil")),
				Body: &ast.BlockStmt{
					List: ifBody,
				},
			})
		}
	}
	bodyList = append(bodyList, &ast.AssignStmt{
		Lhs: []ast.Expr{I("err")},
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun:  Sel(Sel(I("h"), "validator"), "Struct"),
				Args: []ast.Expr{I("cookies")},
			},
		},
	})
	bodyList = append(bodyList, &ast.IfStmt{
		Cond: Ne(I("err"), I("nil")),
		Body: &ast.BlockStmt{List: []ast.Stmt{Ret2(I("nil"), I("err"))}},
	})
	bodyList = append(bodyList, Ret2(Amp(I("cookies")), I("nil")))
	g.HandlersFile.restDecls = append(g.HandlersFile.restDecls, Func("parse"+baseName+"Cookies",
		Field("h", Star(I("Handler")), ""),
		[]*ast.Field{
			Field("r", Star(Sel(I("http"), "Request")), ""),
		},
		[]*ast.Field{
			Field("", Star(Sel(I(g.GetCurrentModelsPackage()), baseName+"Cookies")), ""),
			Field("", I("error"), ""),
		},
		bodyList,
	))

	return nil
}

func (g *Generator) GetValidateFuncStmt(typeName string, ref string) ast.Expr {
	validateFuncName := "Validate" + typeName + "JSON"

	if ref == "" || !refIsExternal(ref) {
		return I(validateFuncName)
	}

	filename := parseFilenameFromRef(ref)
	if filename == "" {
		return I(validateFuncName)
	}

	parts := strings.Split(ref, "/")
	if len(parts) == 0 {
		return I(validateFuncName)
	}

	validateFuncName = "Validate" + parts[len(parts)-1] + "JSON"

	g.YAMLFilesToProcess = append(g.YAMLFilesToProcess, g.GetYAMLFilePath(filename))
	g.AddHandlersImport(g.GetHandlersImportForFile(filename))
	modelName := g.GetModelName(filename)
	return Sel(I(modelName), validateFuncName)
}

func (g *Generator) AddParseRequestBodyMethod(baseName string, contentType string, body *openapi3.RequestBodyRef) error {
	bodyList := []ast.Stmt{}
	if !body.Value.Required {
		bodyList = append(bodyList, &ast.IfStmt{
			Cond: Eq(Sel(I("r"), "Body"), I("nil")),
			Body: &ast.BlockStmt{List: []ast.Stmt{Ret2(I("nil"), I("nil"))}},
		})
	}

	typeName := baseName + "RequestBody"
	var bodyType ast.Expr
	content, ok := body.Value.Content[contentType]
	bodyType = Sel(I(g.GetCurrentModelsPackage()), typeName)
	if ok && content.Schema != nil {
		if content.Schema.Ref != "" {
			var importPath string
			typeName, importPath = g.ParseRefTypeName(content.Schema.Ref)
			bodyType = Sel(I(g.GetCurrentModelsPackage()), typeName)
			if importPath != "" {
				g.AddHandlersImport(importPath)
			}
			if refIsExternal(content.Schema.Ref) {
				bodyType = I(typeName)
			}
		}
	}
	bodyList = append(bodyList, &ast.DeclStmt{
		Decl: &ast.GenDecl{
			Tok: token.VAR,
			Specs: []ast.Spec{
				&ast.ValueSpec{
					Names: []*ast.Ident{I("bodyJSON")},
					Type:  Sel(I("json"), "RawMessage"),
				},
			},
		},
	})
	g.AddHandlersImport("encoding/json")
	bodyList = append(bodyList, &ast.AssignStmt{
		Lhs: []ast.Expr{I("err")},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: Sel(
					&ast.CallExpr{
						Fun:  Sel(I("json"), "NewDecoder"),
						Args: []ast.Expr{Sel(I("r"), "Body")},
					},
					"Decode",
				),
				Args: []ast.Expr{
					Amp(I("bodyJSON")),
				},
			},
		},
	})
	bodyList = append(bodyList, &ast.IfStmt{
		Cond: Ne(I("err"), I("nil")),
		Body: &ast.BlockStmt{List: []ast.Stmt{Ret2(I("nil"), I("err"))}},
	})
	bodyList = append(bodyList, &ast.AssignStmt{
		Lhs: []ast.Expr{I("err")},
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: g.GetValidateFuncStmt(typeName, content.Schema.Ref),
				Args: []ast.Expr{
					I("bodyJSON"),
				},
			},
		},
	})
	bodyList = append(bodyList, &ast.IfStmt{
		Cond: Ne(I("err"), I("nil")),
		Body: &ast.BlockStmt{List: []ast.Stmt{Ret2(I("nil"), I("err"))}},
	})

	bodyList = append(bodyList, &ast.DeclStmt{
		Decl: &ast.GenDecl{
			Tok: token.VAR,
			Specs: []ast.Spec{
				&ast.ValueSpec{
					Names: []*ast.Ident{I("body")},
					Type:  bodyType,
				},
			},
		},
	})

	bodyList = append(bodyList, &ast.AssignStmt{
		Lhs: []ast.Expr{I("err")},
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun:  Sel(I("json"), "Unmarshal"),
				Args: []ast.Expr{I("bodyJSON"), Amp(I("body"))},
			},
		},
	})
	bodyList = append(bodyList, &ast.IfStmt{
		Cond: Ne(I("err"), I("nil")),
		Body: &ast.BlockStmt{List: []ast.Stmt{Ret2(I("nil"), I("err"))}},
	})

	bodyList = append(bodyList, &ast.AssignStmt{
		Lhs: []ast.Expr{I("err")},
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: Sel(Sel(I("h"), "validator"), "Struct"),
				Args: []ast.Expr{
					I("body"),
				},
			},
		},
	})
	bodyList = append(bodyList, &ast.IfStmt{
		Cond: Ne(I("err"), I("nil")),
		Body: &ast.BlockStmt{List: []ast.Stmt{Ret2(I("nil"), I("err"))}},
	})
	bodyList = append(bodyList, Ret2(Amp(I("body")), I("nil")))

	g.HandlersFile.restDecls = append(g.HandlersFile.restDecls, Func(
		"parse"+baseName+"RequestBody",
		Field("h", Star(I("Handler")), ""),
		[]*ast.Field{
			Field("r", Star(Sel(I("http"), "Request")), ""),
		},
		[]*ast.Field{
			Field("", Star(bodyType), ""),
			Field("", I("error"), ""),
		},
		bodyList,
	))

	return nil
}

func (g *Generator) AddParseRequestMethod(baseName string, contentType string, pathParams openapi3.Parameters,
	queryParams openapi3.Parameters, headers openapi3.Parameters, cookieParams openapi3.Parameters,
	body *openapi3.RequestBodyRef,
) {
	bodyList := []ast.Stmt{}
	elts := []ast.Expr{}
	if len(pathParams) > 0 {
		elts = append(elts, &ast.KeyValueExpr{
			Key:   I("Path"),
			Value: Star(I("pathParams")),
		})
		bodyList = append(bodyList, &ast.AssignStmt{
			Lhs: []ast.Expr{
				I("pathParams"),
				I("err"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: Sel(I("h"), "parse"+baseName+"PathParams"),
					Args: []ast.Expr{
						I("r"),
					},
				},
			},
		})
		bodyList = append(bodyList, &ast.IfStmt{
			Cond: Ne(I("err"), I("nil")),
			Body: &ast.BlockStmt{List: []ast.Stmt{Ret2(I("nil"), I("err"))}},
		})
	}
	if len(queryParams) > 0 {
		elts = append(elts, &ast.KeyValueExpr{
			Key:   I("Query"),
			Value: Star(I("queryParams")),
		})
		bodyList = append(bodyList, &ast.AssignStmt{
			Lhs: []ast.Expr{
				I("queryParams"),
				I("err"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: Sel(I("h"), "parse"+baseName+"QueryParams"),
					Args: []ast.Expr{
						I("r"),
					},
				},
			},
		})
		bodyList = append(bodyList, &ast.IfStmt{
			Cond: Ne(I("err"), I("nil")),
			Body: &ast.BlockStmt{List: []ast.Stmt{Ret2(I("nil"), I("err"))}},
		})
	}
	if len(headers) > 0 {
		elts = append(elts, &ast.KeyValueExpr{
			Key:   I("Headers"),
			Value: Star(I("headers")),
		})
		bodyList = append(bodyList, &ast.AssignStmt{
			Lhs: []ast.Expr{
				I("headers"),
				I("err"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: Sel(I("h"), "parse"+baseName+"Headers"),
					Args: []ast.Expr{
						I("r"),
					},
				},
			},
		})
		bodyList = append(bodyList, &ast.IfStmt{
			Cond: Ne(I("err"), I("nil")),
			Body: &ast.BlockStmt{List: []ast.Stmt{Ret2(I("nil"), I("err"))}},
		})
	}
	if len(cookieParams) > 0 {
		elts = append(elts, &ast.KeyValueExpr{
			Key:   I("Cookies"),
			Value: Star(I("cookieParams")),
		})
		bodyList = append(bodyList, &ast.AssignStmt{
			Lhs: []ast.Expr{
				I("cookieParams"),
				I("err"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: Sel(I("h"), "parse"+baseName+"Cookies"),
					Args: []ast.Expr{
						I("r"),
					},
				},
			},
		})
		bodyList = append(bodyList, &ast.IfStmt{
			Cond: Ne(I("err"), I("nil")),
			Body: &ast.BlockStmt{List: []ast.Stmt{Ret2(I("nil"), I("err"))}},
		})
	}
	if body != nil && body.Value != nil {
		content, ok := body.Value.Content[contentType]
		if ok && content.Schema != nil {
			if body.Value.Required {
				elts = append(elts, &ast.KeyValueExpr{
					Key:   I("Body"),
					Value: Star(I("body")),
				})
			} else {
				elts = append(elts, &ast.KeyValueExpr{
					Key:   I("Body"),
					Value: I("body"),
				})
			}
			bodyList = append(bodyList, &ast.AssignStmt{
				Lhs: []ast.Expr{
					I("body"),
					I("err"),
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: Sel(I("h"), "parse"+baseName+"RequestBody"),
						Args: []ast.Expr{
							I("r"),
						},
					},
				},
			})
			bodyList = append(bodyList, &ast.IfStmt{
				Cond: Ne(I("err"), I("nil")),
				Body: &ast.BlockStmt{List: []ast.Stmt{Ret2(I("nil"), I("err"))}},
			})
		}
	}

	bodyList = append(bodyList,
		Ret2(Amp(&ast.CompositeLit{
			Type: Sel(I(g.GetCurrentModelsPackage()), baseName+"Request"),
			Elts: elts,
		}),
			I("nil"),
		),
	)

	g.HandlersFile.restDecls = append(g.HandlersFile.restDecls, Func(
		"parse"+baseName+"Request",
		Field("h", Star(I("Handler")), ""),
		[]*ast.Field{
			Field("r", Star(Sel(I("http"), "Request")), ""),
		},
		[]*ast.Field{
			Field("", Star(Sel(I(g.GetCurrentModelsPackage()), baseName+"Request")), ""),
			Field("", I("error"), ""),
		},
		bodyList,
	))
}

func (g *Generator) AddCreateResponseModel(baseName string, code string, response *openapi3.ResponseRef) error {
	arglist := []*ast.Field{}
	constructorArgs := []ast.Expr{}

	if len(response.Value.Content) > 0 {
		// assume there is a json body
		json, ok := response.Value.Content["application/json"]
		if !ok {
			return errors.New("response content type 'application/json' not found")
		}
		if json.Schema != nil {
			typeName := baseName + "Response" + code + "Body"
			var astType ast.Expr
			astType = Sel(I(g.GetCurrentModelsPackage()), typeName)
			if json.Schema.Ref != "" {
				schemaRef := resolveSchemaRefAgainstResponse(response.Ref, json.Schema.Ref)
				var importPath string
				typeName, importPath = g.ParseRefTypeName(schemaRef)
				if refIsExternal(schemaRef) {
					astType = I(typeName)
				} else {
					astType = Sel(I(g.GetCurrentModelsPackage()), typeName)
				}
				if importPath != "" {
					g.AddHandlersImport(importPath)
				}
			}
			arglist = append(arglist, &ast.Field{
				Names: []*ast.Ident{I("body")},
				Type:  astType,
			})
			constructorArgs = append(constructorArgs, &ast.KeyValueExpr{
				Key:   I("Body"),
				Value: I("body"),
			})
		}
	}

	if len(response.Value.Headers) > 0 {
		arglist = append(arglist, &ast.Field{
			Names: []*ast.Ident{I("headers")},
			Type:  Sel(I(g.GetCurrentModelsPackage()), baseName+"Response"+code+"Headers"),
		})
		constructorArgs = append(constructorArgs, &ast.KeyValueExpr{
			Key:   I("Headers"),
			Value: I("headers"),
		})
	}

	g.HandlersFile.restDecls = append(g.HandlersFile.restDecls, Func(baseName+code+"Response",
		nil,
		arglist,
		[]*ast.Field{
			Field("", Star(Sel(I(g.GetCurrentModelsPackage()), baseName+"Response")), ""),
		},
		[]ast.Stmt{Ret1(
			Amp(&ast.CompositeLit{
				Type: Sel(I(g.GetCurrentModelsPackage()), baseName+"Response"),
				Elts: []ast.Expr{
					&ast.KeyValueExpr{
						Key: I("StatusCode"),
						Value: &ast.BasicLit{
							Kind:  token.INT,
							Value: code,
						},
					},
					&ast.KeyValueExpr{
						Key: I("Response" + code),
						Value: Amp(&ast.CompositeLit{
							Type: Sel(I(g.GetCurrentModelsPackage()), baseName+"Response"+code),
							Elts: constructorArgs,
						}),
					},
				},
			}),
		)},
	))

	return nil
}

func (g *Generator) AddContainsNullIfNeeded() {
	if g.HandlersFile.hasContainsNullMethod {
		return
	}

	g.HandlersFile.hasContainsNullMethod = true
	g.HandlersFile.restDecls = append(g.HandlersFile.restDecls, Func("containsNull",
		nil,
		[]*ast.Field{
			Field("data", Sel(I("json"), "RawMessage"), ""),
		},
		[]*ast.Field{
			Field("", I("bool"), ""),
		},
		[]ast.Stmt{
			&ast.DeclStmt{
				Decl: &ast.GenDecl{
					Tok: token.VAR,
					Specs: []ast.Spec{
						&ast.ValueSpec{
							Names: []*ast.Ident{I("temp")},
							Type:  I("any"),
						},
					},
				},
			},
			&ast.AssignStmt{
				Lhs: []ast.Expr{I("err")},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: Sel(I("json"), "Unmarshal"),
						Args: []ast.Expr{
							I("data"),
							Amp(I("temp")),
						},
					},
				},
			},
			&ast.IfStmt{
				Cond: Ne(I("err"), I("nil")),
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						Ret1(I("false")),
					},
				},
			},
			&ast.ReturnStmt{
				Results: []ast.Expr{
					&ast.BinaryExpr{
						X:  I("temp"),
						Op: token.EQL,
						Y:  I("nil"),
					},
				},
			},
		},
	))
	g.AddHandlersImport("encoding/json")
}

func (g *Generator) AddObjectValidate(modelName string, schema *openapi3.SchemaRef) error {
	const op = "generator.AddObjectValidate"
	requiredFieldsMap := make(map[string]bool, 0)
	nullableFields := make([]string, 0)
	objectFields := make(map[string]ast.Expr, 0)

	for _, requiredField := range schema.Value.Required {
		requiredFieldsMap[requiredField] = true
	}

	for fieldName, fieldSchema := range schema.Value.Properties {
		if fieldSchema.Value == nil {
			continue
		}
		if fieldSchema.Value.Nullable && requiredFieldsMap[fieldName] {
			nullableFields = append(nullableFields, fieldName)
		}
		if fieldSchema.Value.Type.Permits(openapi3.TypeObject) {
			fieldType, err := g.GetFieldTypeFromSchema(modelName, fieldName, fieldSchema)
			if err != nil {
				return errors.Wrap(err, op)
			}
			objectFields[fieldName] = g.GetValidateFuncStmt(fieldType, fieldSchema.Ref)
		}
		if fieldSchema.Value.Type.Permits(openapi3.TypeArray) {
			if fieldSchema.Value.Items != nil {
				itemsType := g.getMostNestedArrayItemType(fieldSchema.Value.Items)
				if itemsType != nil && itemsType.Permits(openapi3.TypeObject) {
					fieldType, err := g.GetFieldTypeFromSchema(modelName, fieldName, fieldSchema)
					if err != nil {
						return errors.Wrap(err, op)
					}
					objectFields[fieldName] = g.GetValidateFuncStmt(fieldType, fieldSchema.Ref)
				}
			}
		}
	}
	requiredFields := make([]string, 0, len(requiredFieldsMap))
	for fieldName := range requiredFieldsMap {
		requiredFields = append(requiredFields, fieldName)
	}
	sort.Strings(requiredFields)
	sort.Strings(nullableFields)

	funcBody := make([]ast.Stmt, 0, len(objectFields))

	if len(requiredFields) > 0 {
		requiredFieldsElts := make([]ast.Expr, 0, len(requiredFields))
		for _, fieldName := range requiredFields {
			requiredFieldsElts = append(requiredFieldsElts, &ast.KeyValueExpr{
				Key:   Str(fieldName),
				Value: I("true"),
			})
		}
		funcBody = append(funcBody, &ast.AssignStmt{
			Lhs: []ast.Expr{I("requiredFields")},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CompositeLit{
					Type: &ast.MapType{
						Key:   I("string"),
						Value: I("bool"),
					},
					Elts: requiredFieldsElts,
				},
			},
		})

		nullableFieldsElts := make([]ast.Expr, 0, len(nullableFields))
		for _, fieldName := range nullableFields {
			nullableFieldsElts = append(nullableFieldsElts, &ast.KeyValueExpr{
				Key:   Str(fieldName),
				Value: I("true"),
			})
		}
		funcBody = append(funcBody, &ast.AssignStmt{
			Lhs: []ast.Expr{I("nullableFields")},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CompositeLit{
					Type: &ast.MapType{
						Key:   I("string"),
						Value: I("bool"),
					},
					Elts: nullableFieldsElts,
				},
			},
		})
	}

	if len(requiredFields) > 0 || len(objectFields) > 0 {
		funcBody = append(funcBody, &ast.DeclStmt{
			Decl: &ast.GenDecl{
				Tok: token.VAR,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: []*ast.Ident{I("obj")},
						Type:  &ast.MapType{Key: I("string"), Value: Sel(I("json"), "RawMessage")},
					},
				},
			},
		})
		funcBody = append(funcBody, &ast.AssignStmt{
			Lhs: []ast.Expr{I("err")},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: Sel(I("json"), "Unmarshal"),
					Args: []ast.Expr{
						I("jsonData"),
						Amp(I("obj")),
					},
				},
			},
		})
		funcBody = append(funcBody, &ast.IfStmt{
			Cond: Ne(I("err"), I("nil")),
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					Ret1(I("err")),
				},
			},
		})
	}
	if len(requiredFields) > 0 || len(objectFields) > 0 {
		funcBody = append(funcBody, &ast.DeclStmt{
			Decl: &ast.GenDecl{
				Tok: token.VAR,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: []*ast.Ident{I("val")},
						Type:  Sel(I("json"), "RawMessage"),
					},
				},
			},
		})
		funcBody = append(funcBody, &ast.DeclStmt{
			Decl: &ast.GenDecl{
				Tok: token.VAR,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: []*ast.Ident{I("exists")},
						Type:  I("bool"),
					},
				},
			},
		})
	}
	if len(requiredFields) > 0 {
		funcBody = append(funcBody, &ast.RangeStmt{
			Key: I("field"),
			Tok: token.DEFINE,
			X:   I("requiredFields"),
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.AssignStmt{
						Lhs: []ast.Expr{I("val"), I("exists")},
						Tok: token.ASSIGN,
						Rhs: []ast.Expr{
							&ast.IndexExpr{
								X:     I("obj"),
								Index: I("field"),
							},
						},
					},
					&ast.IfStmt{
						Cond: &ast.UnaryExpr{
							Op: token.NOT,
							X:  I("exists"),
						},
						Body: &ast.BlockStmt{
							List: []ast.Stmt{
								Ret1(&ast.CallExpr{
									Fun: Sel(I("errors"), "New"),
									Args: []ast.Expr{
										&ast.BinaryExpr{
											X: &ast.BinaryExpr{
												X:  Str("field "),
												Op: token.ADD,
												Y:  I("field"),
											},
											Op: token.ADD,
											Y:  Str(" is required"),
										},
									},
								}),
							},
						},
					},
					&ast.IfStmt{
						Cond: &ast.BinaryExpr{
							X: &ast.UnaryExpr{
								Op: token.NOT,
								X: &ast.IndexExpr{
									X:     I("nullableFields"),
									Index: I("field"),
								},
							},
							Op: token.LAND,
							Y: &ast.CallExpr{
								Fun:  I("containsNull"),
								Args: []ast.Expr{I("val")},
							},
						},
						Body: &ast.BlockStmt{
							List: []ast.Stmt{
								Ret1(&ast.CallExpr{
									Fun: Sel(I("errors"), "New"),
									Args: []ast.Expr{
										&ast.BinaryExpr{
											X: &ast.BinaryExpr{
												X:  Str("field "),
												Op: token.ADD,
												Y:  I("field"),
											},
											Op: token.ADD,
											Y:  Str(" cannot be null"),
										},
									},
								}),
							},
						},
					},
				},
			},
		})
		g.AddContainsNullIfNeeded()
		g.AddHandlersImport("github.com/go-faster/errors")
	}

	objectFieldsNames := make([]string, 0, len(objectFields))
	for fieldName := range objectFields {
		objectFieldsNames = append(objectFieldsNames, fieldName)
	}
	sort.Strings(objectFieldsNames)

	for _, fieldName := range objectFieldsNames {
		fieldValidationFunc := objectFields[fieldName]
		funcBody = append(funcBody, &ast.AssignStmt{
			Lhs: []ast.Expr{I("val"), I("exists")},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				&ast.IndexExpr{
					X:     I("obj"),
					Index: Str(fieldName),
				},
			},
		})
		funcBody = append(funcBody, &ast.IfStmt{
			Cond: &ast.BinaryExpr{
				X:  I("exists"),
				Op: token.LAND,
				Y: &ast.UnaryExpr{
					Op: token.NOT,
					X: &ast.CallExpr{
						Fun:  I("containsNull"),
						Args: []ast.Expr{I("val")},
					},
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.AssignStmt{
						Lhs: []ast.Expr{I("err")},
						Tok: token.ASSIGN,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun:  fieldValidationFunc,
								Args: []ast.Expr{I("val")},
							},
						},
					},
					&ast.IfStmt{
						Cond: Ne(I("err"), I("nil")),
						Body: &ast.BlockStmt{
							List: []ast.Stmt{
								Ret1(&ast.CallExpr{
									Fun: Sel(I("errors"), "Wrap"),
									Args: []ast.Expr{
										I("err"),
										Str("field " + fieldName + " is not valid"),
									},
								}),
							},
						},
					},
				},
			},
		})

		g.AddContainsNullIfNeeded()
		g.AddHandlersImport("github.com/go-faster/errors")
	}

	funcBody = append(funcBody, &ast.ReturnStmt{
		Results: []ast.Expr{I("nil")},
	})

	fieldName := "jsonData"
	if len(requiredFields) == 0 && len(objectFields) == 0 {
		fieldName = "_"
	}

	g.HandlersFile.restDecls = append(g.HandlersFile.restDecls, Func("Validate"+modelName+"JSON",
		nil,
		[]*ast.Field{
			Field(fieldName, Sel(I("json"), "RawMessage"), ""),
		},
		[]*ast.Field{
			Field("", I("error"), ""),
		},
		funcBody,
	))
	return nil
}

func (g *Generator) AddArrayValidate(modelName string, schema *openapi3.SchemaRef) error {
	const op = "generator.AddArrayValidate"

	elemType, err := g.GetFieldTypeFromSchema(modelName, "Item", schema.Value.Items)
	if err != nil {
		return errors.Wrap(err, op)
	}
	validateFunc := g.GetValidateFuncStmt(elemType, schema.Value.Items.Ref)
	g.HandlersFile.restDecls = append(g.HandlersFile.restDecls, Func("Validate"+modelName+"JSON",
		nil,
		[]*ast.Field{
			Field("jsonData", Sel(I("json"), "RawMessage"), ""),
		},
		[]*ast.Field{
			Field("", I("error"), ""),
		},
		[]ast.Stmt{
			&ast.DeclStmt{
				Decl: &ast.GenDecl{
					Tok: token.VAR,
					Specs: []ast.Spec{
						&ast.ValueSpec{
							Names: []*ast.Ident{I("arr")},
							Type:  &ast.ArrayType{Elt: Sel(I("json"), "RawMessage")},
						},
					},
				},
			},
			&ast.AssignStmt{
				Lhs: []ast.Expr{I("err")},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: Sel(I("json"), "Unmarshal"),
						Args: []ast.Expr{
							I("jsonData"),
							Amp(I("arr")),
						},
					},
				},
			},
			&ast.IfStmt{
				Cond: Ne(I("err"), I("nil")),
				Body: &ast.BlockStmt{
					List: []ast.Stmt{Ret1(I("err"))},
				},
			},
			&ast.RangeStmt{
				Key:   I("index"),
				Value: I("obj"),
				Tok:   token.DEFINE,
				X:     I("arr"),
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.IfStmt{
							Cond: &ast.UnaryExpr{
								Op: token.NOT,
								X: &ast.CallExpr{
									Fun:  I("containsNull"),
									Args: []ast.Expr{I("obj")},
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.AssignStmt{
										Lhs: []ast.Expr{I("err")},
										Tok: token.ASSIGN,
										Rhs: []ast.Expr{
											&ast.CallExpr{
												Fun:  validateFunc,
												Args: []ast.Expr{I("obj")},
											},
										},
									},
									&ast.IfStmt{
										Cond: Ne(I("err"), I("nil")),
										Body: &ast.BlockStmt{
											List: []ast.Stmt{Ret1(&ast.CallExpr{
												Fun: Sel(I("errors"), "Wrapf"),
												Args: []ast.Expr{
													I("err"),
													Str("error validating object at index %d"),
													I("index"),
												},
											})},
										},
									},
								},
							},
						},
					},
				},
			},
			&ast.ReturnStmt{
				Results: []ast.Expr{I("nil")},
			},
		},
	))
	g.AddHandlersImport("github.com/go-faster/errors")

	return nil
}
