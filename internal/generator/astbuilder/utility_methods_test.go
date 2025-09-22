package astbuilder

import (
	"go/ast"
	"go/token"
	"testing"
)

func TestBuilder_ExpressionMethods(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	// Test Call method
	callExpr := builder.Call("fmt", "Printf", builder.String("Hello %s"), builder.Ident("name"))
	if callExpr == nil {
		t.Fatal("Call returned nil")
	}
	if call, ok := callExpr.(*ast.CallExpr); !ok {
		t.Fatal("Call did not return CallExpr")
	} else {
		if sel, ok := call.Fun.(*ast.SelectorExpr); !ok {
			t.Fatal("Call.Fun is not SelectorExpr")
		} else {
			if sel.Sel.Name != "Printf" {
				t.Errorf("Expected method name 'Printf', got '%s'", sel.Sel.Name)
			}
		}
		if len(call.Args) != 2 {
			t.Errorf("Expected 2 arguments, got %d", len(call.Args))
		}
	}

	// Test Call without receiver
	callExpr2 := builder.Call("", "someFunction", builder.String("arg"))
	if callExpr2 == nil {
		t.Fatal("Call without receiver returned nil")
	}
	if call, ok := callExpr2.(*ast.CallExpr); !ok {
		t.Fatal("Call without receiver did not return CallExpr")
	} else {
		if ident, ok := call.Fun.(*ast.Ident); !ok {
			t.Fatal("Call.Fun is not Ident")
		} else {
			if ident.Name != "someFunction" {
				t.Errorf("Expected function name 'someFunction', got '%s'", ident.Name)
			}
		}
	}

	// Test Select method
	selectExpr := builder.Select("user", "Name")
	if selectExpr == nil {
		t.Fatal("Select returned nil")
	}
	if sel, ok := selectExpr.(*ast.SelectorExpr); !ok {
		t.Fatal("Select did not return SelectorExpr")
	} else {
		if ident, ok := sel.X.(*ast.Ident); !ok {
			t.Fatal("Select.X is not Ident")
		} else {
			if ident.Name != "user" {
				t.Errorf("Expected receiver 'user', got '%s'", ident.Name)
			}
		}
		if sel.Sel.Name != "Name" {
			t.Errorf("Expected field name 'Name', got '%s'", sel.Sel.Name)
		}
	}

	// Test Ident method
	identExpr := builder.Ident("variable")
	if identExpr == nil {
		t.Fatal("Ident returned nil")
	}
	if ident, ok := identExpr.(*ast.Ident); !ok {
		t.Fatal("Ident did not return Ident")
	} else {
		if ident.Name != "variable" {
			t.Errorf("Expected identifier 'variable', got '%s'", ident.Name)
		}
	}

	// Test String method
	strExpr := builder.String("hello world")
	if strExpr == nil {
		t.Fatal("String returned nil")
	}
	if lit, ok := strExpr.(*ast.BasicLit); !ok {
		t.Fatal("String did not return BasicLit")
	} else {
		if lit.Kind != token.STRING {
			t.Error("String did not return STRING token")
		}
		if lit.Value != `"hello world"` {
			t.Errorf("Expected string value '\"hello world\"', got '%s'", lit.Value)
		}
	}

	// Test Int method
	intExpr := builder.Int(42)
	if intExpr == nil {
		t.Fatal("Int returned nil")
	}
	if lit, ok := intExpr.(*ast.BasicLit); !ok {
		t.Fatal("Int did not return BasicLit")
	} else {
		if lit.Kind != token.INT {
			t.Error("Int did not return INT token")
		}
		if lit.Value != "42" {
			t.Errorf("Expected int value '42', got '%s'", lit.Value)
		}
	}

	// Test Bool method
	boolExpr := builder.Bool(true)
	if boolExpr == nil {
		t.Fatal("Bool returned nil")
	}
	if ident, ok := boolExpr.(*ast.Ident); !ok {
		t.Fatal("Bool did not return Ident")
	} else {
		if ident.Name != "true" {
			t.Errorf("Expected bool value 'true', got '%s'", ident.Name)
		}
	}

	// Test Nil method
	nilExpr := builder.Nil()
	if nilExpr == nil {
		t.Fatal("Nil returned nil")
	}
	if ident, ok := nilExpr.(*ast.Ident); !ok {
		t.Fatal("Nil did not return Ident")
	} else {
		if ident.Name != "nil" {
			t.Errorf("Expected nil value 'nil', got '%s'", ident.Name)
		}
	}

	// Test AddressOf method
	addrExpr := builder.AddressOf(builder.Ident("variable"))
	if addrExpr == nil {
		t.Fatal("AddressOf returned nil")
	}
	if unary, ok := addrExpr.(*ast.UnaryExpr); !ok {
		t.Fatal("AddressOf did not return UnaryExpr")
	} else {
		if unary.Op != token.AND {
			t.Error("AddressOf did not return AND token")
		}
	}

	// Test Deref method
	derefExpr := builder.Deref(builder.Ident("pointer"))
	if derefExpr == nil {
		t.Fatal("Deref returned nil")
	}
	if unary, ok := derefExpr.(*ast.UnaryExpr); !ok {
		t.Fatal("Deref did not return UnaryExpr")
	} else {
		if unary.Op != token.MUL {
			t.Error("Deref did not return MUL token")
		}
	}
}

func TestBuilder_StatementMethods(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	// Test DeclareVar method
	declStmt := builder.DeclareVar("name", "string", builder.String("John"))
	if declStmt == nil {
		t.Fatal("DeclareVar returned nil")
	}
	if decl, ok := declStmt.(*ast.DeclStmt); !ok {
		t.Fatal("DeclareVar did not return DeclStmt")
	} else {
		if genDecl, ok := decl.Decl.(*ast.GenDecl); !ok {
			t.Fatal("DeclareVar.Decl is not GenDecl")
		} else {
			if genDecl.Tok != token.VAR {
				t.Error("DeclareVar did not return VAR token")
			}
			if len(genDecl.Specs) != 1 {
				t.Errorf("Expected 1 spec, got %d", len(genDecl.Specs))
			}
			if spec, ok := genDecl.Specs[0].(*ast.ValueSpec); !ok {
				t.Fatal("Spec is not ValueSpec")
			} else {
				if len(spec.Names) != 1 {
					t.Errorf("Expected 1 name, got %d", len(spec.Names))
				}
				if spec.Names[0].Name != "name" {
					t.Errorf("Expected name 'name', got '%s'", spec.Names[0].Name)
				}
				if spec.Type.(*ast.Ident).Name != "string" {
					t.Errorf("Expected type 'string', got '%s'", spec.Type.(*ast.Ident).Name)
				}
				if len(spec.Values) != 1 {
					t.Errorf("Expected 1 value, got %d", len(spec.Values))
				}
			}
		}
	}

	// Test DeclareVar without value
	declStmt2 := builder.DeclareVar("count", "int", nil)
	if declStmt2 == nil {
		t.Fatal("DeclareVar without value returned nil")
	}
	if decl, ok := declStmt2.(*ast.DeclStmt); !ok {
		t.Fatal("DeclareVar without value did not return DeclStmt")
	} else {
		if genDecl, ok := decl.Decl.(*ast.GenDecl); !ok {
			t.Fatal("DeclareVar without value.Decl is not GenDecl")
		} else {
			if spec, ok := genDecl.Specs[0].(*ast.ValueSpec); !ok {
				t.Fatal("Spec is not ValueSpec")
			} else {
				if len(spec.Values) != 0 {
					t.Errorf("Expected 0 values, got %d", len(spec.Values))
				}
			}
		}
	}

	// Test Assign method
	assignStmt := builder.Assign(builder.Ident("x"), builder.Int(10))
	if assignStmt == nil {
		t.Fatal("Assign returned nil")
	}
	if assign, ok := assignStmt.(*ast.AssignStmt); !ok {
		t.Fatal("Assign did not return AssignStmt")
	} else {
		if assign.Tok != token.ASSIGN {
			t.Error("Assign did not return ASSIGN token")
		}
		if len(assign.Lhs) != 1 {
			t.Errorf("Expected 1 LHS, got %d", len(assign.Lhs))
		}
		if len(assign.Rhs) != 1 {
			t.Errorf("Expected 1 RHS, got %d", len(assign.Rhs))
		}
	}

	// Test If method
	ifStmt := builder.If(builder.Ident("condition"), []ast.Stmt{
		builder.CallStmt("", "doSomething"),
	})
	if ifStmt == nil {
		t.Fatal("If returned nil")
	}
	if ifSt, ok := ifStmt.(*ast.IfStmt); !ok {
		t.Fatal("If did not return IfStmt")
	} else {
		if ifSt.Cond == nil {
			t.Error("If.Cond is nil")
		}
		if ifSt.Body == nil {
			t.Error("If.Body is nil")
		} else {
			if len(ifSt.Body.List) != 1 {
				t.Errorf("Expected 1 statement in body, got %d", len(ifSt.Body.List))
			}
		}
	}

	// Test IfElse method
	ifElseStmt := builder.IfElse(
		builder.Ident("condition"),
		[]ast.Stmt{builder.CallStmt("", "doIf")},
		[]ast.Stmt{builder.CallStmt("", "doElse")},
	)
	if ifElseStmt == nil {
		t.Fatal("IfElse returned nil")
	}
	if ifSt, ok := ifElseStmt.(*ast.IfStmt); !ok {
		t.Fatal("IfElse did not return IfStmt")
	} else {
		if ifSt.Else == nil {
			t.Error("IfElse.Else is nil")
		} else {
			if block, ok := ifSt.Else.(*ast.BlockStmt); !ok {
				t.Fatal("IfElse.Else is not BlockStmt")
			} else {
				if len(block.List) != 1 {
					t.Errorf("Expected 1 statement in else body, got %d", len(block.List))
				}
			}
		}
	}

	// Test IfElse without else body
	ifElseStmt2 := builder.IfElse(
		builder.Ident("condition"),
		[]ast.Stmt{builder.CallStmt("", "doIf")},
		[]ast.Stmt{},
	)
	if ifElseStmt2 == nil {
		t.Fatal("IfElse without else body returned nil")
	}
	if ifSt, ok := ifElseStmt2.(*ast.IfStmt); !ok {
		t.Fatal("IfElse without else body did not return IfStmt")
	} else {
		if ifSt.Else != nil {
			t.Error("IfElse without else body should have nil Else")
		}
	}

	// Test Return method
	returnStmt := builder.Return(builder.String("success"))
	if returnStmt == nil {
		t.Fatal("Return returned nil")
	}
	if ret, ok := returnStmt.(*ast.ReturnStmt); !ok {
		t.Fatal("Return did not return ReturnStmt")
	} else {
		if len(ret.Results) != 1 {
			t.Errorf("Expected 1 result, got %d", len(ret.Results))
		}
	}

	// Test Return without values
	returnStmt2 := builder.Return()
	if returnStmt2 == nil {
		t.Fatal("Return without values returned nil")
	}
	if ret, ok := returnStmt2.(*ast.ReturnStmt); !ok {
		t.Fatal("Return without values did not return ReturnStmt")
	} else {
		if len(ret.Results) != 0 {
			t.Errorf("Expected 0 results, got %d", len(ret.Results))
		}
	}

	// Test CallStmt method
	callStmt := builder.CallStmt("fmt", "Println", builder.String("Hello"))
	if callStmt == nil {
		t.Fatal("CallStmt returned nil")
	}
	if exprStmt, ok := callStmt.(*ast.ExprStmt); !ok {
		t.Fatal("CallStmt did not return ExprStmt")
	} else {
		if call, ok := exprStmt.X.(*ast.CallExpr); !ok {
			t.Fatal("CallStmt.X is not CallExpr")
		} else {
			if sel, ok := call.Fun.(*ast.SelectorExpr); !ok {
				t.Fatal("CallStmt.Fun is not SelectorExpr")
			} else {
				if sel.Sel.Name != "Println" {
					t.Errorf("Expected method name 'Println', got '%s'", sel.Sel.Name)
				}
			}
		}
	}

	// Test CallStmt without receiver
	callStmt2 := builder.CallStmt("", "someFunction", builder.Int(42))
	if callStmt2 == nil {
		t.Fatal("CallStmt without receiver returned nil")
	}
	if exprStmt, ok := callStmt2.(*ast.ExprStmt); !ok {
		t.Fatal("CallStmt without receiver did not return ExprStmt")
	} else {
		if call, ok := exprStmt.X.(*ast.CallExpr); !ok {
			t.Fatal("CallStmt without receiver.X is not CallExpr")
		} else {
			if ident, ok := call.Fun.(*ast.Ident); !ok {
				t.Fatal("CallStmt without receiver.Fun is not Ident")
			} else {
				if ident.Name != "someFunction" {
					t.Errorf("Expected function name 'someFunction', got '%s'", ident.Name)
				}
			}
		}
	}
}

func TestBuilder_UtilityMethodsIntegration(t *testing.T) {
	// Test integration of utility methods
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	// Create a complex function using utility methods
	body := []ast.Stmt{
		builder.DeclareVar("name", "string", builder.String("John")),
		builder.DeclareVar("age", "int", builder.Int(30)),
		builder.If(
			builder.Call("", "isValid", builder.Ident("name")),
			[]ast.Stmt{
				builder.CallStmt("fmt", "Printf", builder.String("Name: %s, Age: %d"), builder.Ident("name"), builder.Ident("age")),
			},
		),
		builder.Return(builder.Nil()),
	}

	// Create function declaration
	funcDecl := &ast.FuncDecl{
		Name: builder.Ident("processUser").(*ast.Ident),
		Type: &ast.FuncType{
			Params: &ast.FieldList{},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: builder.Ident("error").(*ast.Ident),
					},
				},
			},
		},
		Body: &ast.BlockStmt{List: body},
	}

	// Add function to builder
	builder.AddDeclaration(funcDecl)

	// Build the file
	file := builder.BuildFile()

	// Verify the file structure
	if file == nil {
		t.Fatal("BuildFile returned nil")
	}

	if len(file.Decls) != 1 {
		t.Errorf("Expected 1 declaration, got %d", len(file.Decls))
	}

	if funcDecl, ok := file.Decls[0].(*ast.FuncDecl); !ok {
		t.Fatal("Expected function declaration")
	} else {
		if funcDecl.Name.Name != "processUser" {
			t.Errorf("Expected function name 'processUser', got '%s'", funcDecl.Name.Name)
		}

		if len(funcDecl.Body.List) != 4 {
			t.Errorf("Expected 4 statements in body, got %d", len(funcDecl.Body.List))
		}
	}
}

func TestBuilder_UtilityMethodsChaining(t *testing.T) {
	// Test method chaining and complex expressions
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	// Create complex expression using chaining
	complexExpr := builder.Call(
		"json",
		"Marshal",
		builder.Call(
			"",
			"createUser",
			builder.String("John"),
			builder.Int(30),
			builder.Bool(true),
		),
	)

	if complexExpr == nil {
		t.Fatal("Complex expression returned nil")
	}

	// Create statement using the complex expression
	stmt := builder.Assign(
		builder.Ident("data"),
		complexExpr,
	)

	if stmt == nil {
		t.Fatal("Statement with complex expression returned nil")
	}

	// Verify the structure
	if assign, ok := stmt.(*ast.AssignStmt); !ok {
		t.Fatal("Statement is not AssignStmt")
	} else {
		if len(assign.Lhs) != 1 {
			t.Errorf("Expected 1 LHS, got %d", len(assign.Lhs))
		}
		if len(assign.Rhs) != 1 {
			t.Errorf("Expected 1 RHS, got %d", len(assign.Rhs))
		}

		// Verify the complex expression structure
		if call, ok := assign.Rhs[0].(*ast.CallExpr); !ok {
			t.Fatal("RHS is not CallExpr")
		} else {
			if sel, ok := call.Fun.(*ast.SelectorExpr); !ok {
				t.Fatal("Call.Fun is not SelectorExpr")
			} else {
				if sel.Sel.Name != "Marshal" {
					t.Errorf("Expected method name 'Marshal', got '%s'", sel.Sel.Name)
				}
			}

			if len(call.Args) != 1 {
				t.Errorf("Expected 1 argument, got %d", len(call.Args))
			}

			// Verify nested call
			if nestedCall, ok := call.Args[0].(*ast.CallExpr); !ok {
				t.Fatal("Argument is not CallExpr")
			} else {
				if ident, ok := nestedCall.Fun.(*ast.Ident); !ok {
					t.Fatal("Nested call.Fun is not Ident")
				} else {
					if ident.Name != "createUser" {
						t.Errorf("Expected function name 'createUser', got '%s'", ident.Name)
					}
				}

				if len(nestedCall.Args) != 3 {
					t.Errorf("Expected 3 arguments in nested call, got %d", len(nestedCall.Args))
				}
			}
		}
	}
}
