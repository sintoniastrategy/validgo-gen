package astbuilder

import (
	"go/ast"
	"go/token"
	"testing"
)

func TestIntegration_CompleteFunctionGeneration(t *testing.T) {
	// Test creating a complete function using all builders
	config := BuilderConfig{
		PackageName:  "test",
		ImportPrefix: "github.com/test",
		UsePointers:  true,
	}

	builder := NewBuilder(config)
	exprBuilder := NewExpressionBuilder(builder)
	stmtBuilder := NewStatementBuilder(builder)
	funcBuilder := NewFunctionBuilder(builder)
	var typeBuilder *TypeBuilder

	// Add imports
	builder.AddImport("fmt")
	builder.AddImport("net/http")

	// Create function parameters
	params := []*ast.Field{
		funcBuilder.Param("r", "*http.Request"),
		funcBuilder.Param("w", "http.ResponseWriter"),
	}

	// Create type builder for later use
	typeBuilder = NewTypeBuilder(builder)

	// Create function results
	results := []*ast.Field{
		funcBuilder.ResultAnonymous("error"),
	}

	// Use typeBuilder to create a simple type
	_ = typeBuilder.Ident("string")

	// Build function body
	body := []ast.Stmt{
		// Declare variables
		stmtBuilder.DeclareVar("name", "string", nil),
		stmtBuilder.DeclareVar("age", "int", nil),

		// Extract query parameters
		stmtBuilder.AssignDefine(
			exprBuilder.Ident("name"),
			exprBuilder.MethodCall(
				exprBuilder.MethodCall(exprBuilder.Ident("r"), "URL"),
				"Query",
			),
		),

		// Validate required parameter
		stmtBuilder.If(
			exprBuilder.Equal(exprBuilder.Ident("name"), exprBuilder.String("")),
			[]ast.Stmt{
				stmtBuilder.CallStmt(
					exprBuilder.Select(exprBuilder.Ident("http"), "Error"),
					exprBuilder.Ident("w"),
					exprBuilder.String("name is required"),
					exprBuilder.Select(exprBuilder.Ident("http"), "StatusBadRequest"),
				),
				stmtBuilder.ReturnEmpty(),
			},
		),

		// Set age
		stmtBuilder.Assign(exprBuilder.Ident("age"), exprBuilder.Int(25)),

		// Write response
		stmtBuilder.CallStmt(
			exprBuilder.MethodCall(exprBuilder.Ident("w"), "WriteHeader"),
			exprBuilder.Select(exprBuilder.Ident("http"), "StatusOK"),
		),
		stmtBuilder.CallStmt(
			exprBuilder.MethodCall(exprBuilder.Ident("w"), "Write"),
			exprBuilder.Call(
				exprBuilder.Select(exprBuilder.Ident("fmt"), "Sprintf"),
				exprBuilder.String("Hello %s, age %d"),
				exprBuilder.Ident("name"),
				exprBuilder.Ident("age"),
			),
		),

		// Return success
		stmtBuilder.Return(exprBuilder.Nil()),
	}

	// Create function declaration
	funcDecl := funcBuilder.Func("HandleRequest", nil, params, results, body)

	// Add function to builder
	builder.AddDeclaration(funcDecl)

	// Build the complete file
	file := builder.BuildFile()

	// Verify the file structure
	if file == nil {
		t.Fatal("BuildFile returned nil")
	}

	if file.Name.Name != "test" {
		t.Errorf("Expected package name 'test', got '%s'", file.Name.Name)
	}

	if len(file.Imports) != 2 {
		t.Errorf("Expected 2 imports, got %d", len(file.Imports))
	}

	if len(file.Decls) != 2 {
		t.Errorf("Expected 2 declarations (imports + function), got %d", len(file.Decls))
		for i, decl := range file.Decls {
			t.Logf("Declaration %d: %T", i, decl)
		}
	}

	// Verify function declaration (should be the second one)
	if funcDecl, ok := file.Decls[1].(*ast.FuncDecl); !ok {
		t.Fatal("Expected function declaration")
	} else {
		if funcDecl.Name.Name != "HandleRequest" {
			t.Errorf("Expected function name 'HandleRequest', got '%s'", funcDecl.Name.Name)
		}

		if len(funcDecl.Type.Params.List) != 2 {
			t.Errorf("Expected 2 parameters, got %d", len(funcDecl.Type.Params.List))
		}

		if len(funcDecl.Type.Results.List) != 1 {
			t.Errorf("Expected 1 result, got %d", len(funcDecl.Type.Results.List))
		}

		if len(funcDecl.Body.List) != 8 {
			t.Errorf("Expected 8 statements in body, got %d", len(funcDecl.Body.List))
		}
	}
}

func TestIntegration_StructGeneration(t *testing.T) {
	// Test creating a complete struct using all builders
	config := BuilderConfig{
		PackageName:  "test",
		ImportPrefix: "github.com/test",
		UsePointers:  true,
	}

	builder := NewBuilder(config)
	typeBuilder := NewTypeBuilder(builder)
	exprBuilder := NewExpressionBuilder(builder)

	// Add imports
	builder.AddImport("time")
	builder.AddImport("github.com/go-playground/validator/v10")

	// Create struct fields using type builder
	typeBuilder = NewTypeBuilder(builder)
	fields := []*ast.Field{
		typeBuilder.Field("ID", exprBuilder.Ident("int"), `json:"id" validate:"required"`),
		typeBuilder.Field("Name", exprBuilder.Ident("string"), `json:"name" validate:"required,min=1,max=100"`),
		typeBuilder.Field("Email", exprBuilder.Ident("string"), `json:"email" validate:"required,email"`),
		typeBuilder.Field("Age", exprBuilder.Star(exprBuilder.Ident("int")), `json:"age,omitempty" validate:"omitempty,min=0,max=150"`),
		typeBuilder.Field("CreatedAt", exprBuilder.Star(exprBuilder.Select(exprBuilder.Ident("time"), "Time")), `json:"created_at,omitempty"`),
		typeBuilder.Field("Tags", exprBuilder.SliceType(exprBuilder.Ident("string")), `json:"tags,omitempty" validate:"omitempty,dive,min=1"`),
	}

	// Create struct declaration
	structDecl := typeBuilder.StructAlias("User", fields)

	// Add struct to builder
	builder.AddDeclaration(structDecl)

	// Build the complete file
	file := builder.BuildFile()

	// Verify the file structure
	if file == nil {
		t.Fatal("BuildFile returned nil")
	}

	if file.Name.Name != "test" {
		t.Errorf("Expected package name 'test', got '%s'", file.Name.Name)
	}

	if len(file.Imports) != 2 {
		t.Errorf("Expected 2 imports, got %d", len(file.Imports))
	}

	if len(file.Decls) != 2 {
		t.Errorf("Expected 2 declarations (imports + struct), got %d", len(file.Decls))
	}

	// Verify struct declaration (should be the second one)
	if genDecl, ok := file.Decls[1].(*ast.GenDecl); !ok {
		t.Fatal("Expected type declaration")
	} else if genDecl.Tok != token.TYPE {
		t.Error("Expected TYPE token")
	} else if len(genDecl.Specs) != 1 {
		t.Errorf("Expected 1 spec, got %d", len(genDecl.Specs))
	} else if typeSpec, ok := genDecl.Specs[0].(*ast.TypeSpec); !ok {
		t.Fatal("Expected type spec")
	} else if typeSpec.Name.Name != "User" {
		t.Errorf("Expected type name 'User', got '%s'", typeSpec.Name.Name)
	} else if structType, ok := typeSpec.Type.(*ast.StructType); !ok {
		t.Fatal("Expected struct type")
	} else if len(structType.Fields.List) != 6 {
		t.Errorf("Expected 6 fields, got %d", len(structType.Fields.List))
	}
}

func TestIntegration_InterfaceGeneration(t *testing.T) {
	// Test creating a complete interface using all builders
	config := BuilderConfig{
		PackageName:  "test",
		ImportPrefix: "github.com/test",
		UsePointers:  true,
	}

	builder := NewBuilder(config)
	typeBuilder := NewTypeBuilder(builder)
	funcBuilder := NewFunctionBuilder(builder)

	// Add imports
	builder.AddImport("context")
	builder.AddImport("net/http")

	// Create interface methods
	methods := []*ast.Field{
		funcBuilder.InterfaceMethod(
			"HandleRequest",
			[]*ast.Field{
				funcBuilder.Param("ctx", "context.Context"),
				funcBuilder.Param("r", "*http.Request"),
			},
			[]*ast.Field{
				funcBuilder.ResultAnonymous("error"),
			},
		),
		funcBuilder.InterfaceMethod(
			"Validate",
			[]*ast.Field{
				funcBuilder.Param("data", "interface{}"),
			},
			[]*ast.Field{
				funcBuilder.ResultAnonymous("error"),
			},
		),
	}

	// Create interface declaration
	interfaceDecl := typeBuilder.InterfaceAlias("Handler", methods)

	// Add interface to builder
	builder.AddDeclaration(interfaceDecl)

	// Build the complete file
	file := builder.BuildFile()

	// Verify the file structure
	if file == nil {
		t.Fatal("BuildFile returned nil")
	}

	if file.Name.Name != "test" {
		t.Errorf("Expected package name 'test', got '%s'", file.Name.Name)
	}

	if len(file.Imports) != 2 {
		t.Errorf("Expected 2 imports, got %d", len(file.Imports))
	}

	if len(file.Decls) != 2 {
		t.Errorf("Expected 2 declarations (imports + interface), got %d", len(file.Decls))
	}

	// Verify interface declaration (should be the second one)
	if genDecl, ok := file.Decls[1].(*ast.GenDecl); !ok {
		t.Fatal("Expected type declaration")
	} else if genDecl.Tok != token.TYPE {
		t.Error("Expected TYPE token")
	} else if len(genDecl.Specs) != 1 {
		t.Errorf("Expected 1 spec, got %d", len(genDecl.Specs))
	} else if typeSpec, ok := genDecl.Specs[0].(*ast.TypeSpec); !ok {
		t.Fatal("Expected type spec")
	} else if typeSpec.Name.Name != "Handler" {
		t.Errorf("Expected type name 'Handler', got '%s'", typeSpec.Name.Name)
	} else if interfaceType, ok := typeSpec.Type.(*ast.InterfaceType); !ok {
		t.Fatal("Expected interface type")
	} else if len(interfaceType.Methods.List) != 2 {
		t.Errorf("Expected 2 methods, got %d", len(interfaceType.Methods.List))
	}
}

func TestIntegration_PatternUsage(t *testing.T) {
	// Test using common patterns
	config := BuilderConfig{
		PackageName:  "test",
		ImportPrefix: "github.com/test",
		UsePointers:  true,
	}

	builder := NewBuilder(config)
	patternBuilder := NewPatternBuilder(builder)
	exprBuilder := NewExpressionBuilder(builder)

	// Add imports
	builder.AddImport("encoding/json")
	builder.AddImport("net/http")

	// Use error handling pattern
	errorPattern := patternBuilder.ErrorHandlingPattern(
		"err",
		exprBuilder.Call(exprBuilder.Ident("someFunction")),
		"failed to call function",
	)

	// Use validation pattern
	validationPattern := patternBuilder.ValidationPattern("user", "validator")

	// Use JSON unmarshal pattern
	jsonPattern := patternBuilder.JSONUnmarshalPattern(
		"err",
		exprBuilder.Ident("jsonData"),
		exprBuilder.Ident("user"),
	)

	// Add all patterns to builder
	builder.AddStatements(errorPattern)
	builder.AddStatements(validationPattern)
	builder.AddStatements(jsonPattern)

	// Build the complete file
	file := builder.BuildFile()

	// Verify the file structure
	if file == nil {
		t.Fatal("BuildFile returned nil")
	}

	if file.Name.Name != "test" {
		t.Errorf("Expected package name 'test', got '%s'", file.Name.Name)
	}

	if len(file.Imports) != 2 {
		t.Errorf("Expected 2 imports, got %d", len(file.Imports))
	}

	// Verify statements were added
	stmts, _, _ := builder.Build()
	if len(stmts) != 6 { // 3 patterns * 2 statements each
		t.Errorf("Expected 6 statements, got %d", len(stmts))
	}
}
