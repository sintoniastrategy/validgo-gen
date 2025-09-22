package astbuilder

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

func TestNewParameterParser(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	paramConfig := ParameterConfig{
		BaseName:      "User",
		PackageName:   "api",
		UsePointers:   true,
		ParameterType: "Query",
	}

	parser := NewParameterParser(builder, paramConfig)

	if parser == nil {
		t.Fatal("NewParameterParser returned nil")
	}

	if parser.config.BaseName != "User" {
		t.Errorf("Expected BaseName 'User', got '%s'", parser.config.BaseName)
	}

	if parser.config.ParameterType != "Query" {
		t.Errorf("Expected ParameterType 'Query', got '%s'", parser.config.ParameterType)
	}

	if parser.builder != builder {
		t.Error("Expected parser to use the provided builder")
	}
}

func TestParameterParser_ParseQueryParams(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	paramConfig := ParameterConfig{
		BaseName:    "User",
		PackageName: "api",
		UsePointers: true,
	}

	parser := NewParameterParser(builder, paramConfig)

	// Create test parameters
	params := openapi3.Parameters{
		{
			Value: &openapi3.Parameter{
				Name:     "id",
				In:       "query",
				Required: true,
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: &openapi3.Types{openapi3.TypeString},
					},
				},
			},
		},
		{
			Value: &openapi3.Parameter{
				Name:     "limit",
				In:       "query",
				Required: false,
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: &openapi3.Types{openapi3.TypeInteger},
					},
				},
			},
		},
	}

	err := parser.ParseQueryParams(params)
	if err != nil {
		t.Fatalf("ParseQueryParams failed: %v", err)
	}

	// Check that declarations were added
	decls := builder.decls
	if len(decls) == 0 {
		t.Error("Expected declarations to be added")
	}

	// Check that imports were added
	if !builder.HasImports() {
		t.Error("Expected imports to be added")
	}
}

func TestParameterParser_ParseHeaders(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	paramConfig := ParameterConfig{
		BaseName:    "User",
		PackageName: "api",
		UsePointers: true,
	}

	parser := NewParameterParser(builder, paramConfig)

	// Create test parameters
	params := openapi3.Parameters{
		{
			Value: &openapi3.Parameter{
				Name:     "Authorization",
				In:       "header",
				Required: true,
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: &openapi3.Types{openapi3.TypeString},
					},
				},
			},
		},
	}

	err := parser.ParseHeaders(params)
	if err != nil {
		t.Fatalf("ParseHeaders failed: %v", err)
	}

	// Check that declarations were added
	decls := builder.decls
	if len(decls) == 0 {
		t.Error("Expected declarations to be added")
	}
}

func TestParameterParser_ParseCookies(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	paramConfig := ParameterConfig{
		BaseName:    "User",
		PackageName: "api",
		UsePointers: true,
	}

	parser := NewParameterParser(builder, paramConfig)

	// Create test parameters
	params := openapi3.Parameters{
		{
			Value: &openapi3.Parameter{
				Name:     "session_id",
				In:       "cookie",
				Required: true,
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: &openapi3.Types{openapi3.TypeString},
					},
				},
			},
		},
	}

	err := parser.ParseCookies(params)
	if err != nil {
		t.Fatalf("ParseCookies failed: %v", err)
	}

	// Check that declarations were added
	decls := builder.decls
	if len(decls) == 0 {
		t.Error("Expected declarations to be added")
	}
}

func TestParameterParser_ParsePathParams(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	paramConfig := ParameterConfig{
		BaseName:    "User",
		PackageName: "api",
		UsePointers: true,
	}

	parser := NewParameterParser(builder, paramConfig)

	// Create test parameters
	params := openapi3.Parameters{
		{
			Value: &openapi3.Parameter{
				Name:     "id",
				In:       "path",
				Required: true,
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: &openapi3.Types{openapi3.TypeString},
					},
				},
			},
		},
	}

	err := parser.ParsePathParams(params)
	if err != nil {
		t.Fatalf("ParsePathParams failed: %v", err)
	}

	// Check that declarations were added
	decls := builder.decls
	if len(decls) == 0 {
		t.Error("Expected declarations to be added")
	}

	// Check that chi import was added
	imports := builder.imports
	if !imports["github.com/go-chi/chi/v5"] {
		t.Error("Expected chi import to be added for path parameters")
	}
}

func TestParameterParser_DeclareStruct(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	paramConfig := ParameterConfig{
		BaseName:    "User",
		PackageName: "api",
		UsePointers: true,
		StructName:  "UserQueryParams",
	}

	parser := NewParameterParser(builder, paramConfig)

	// Declare struct
	parser.DeclareStruct()

	// Check that declaration was added
	decls := builder.decls
	if len(decls) == 0 {
		t.Error("Expected struct declaration to be added")
	}
}

func TestParameterParser_ExtractParameter(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	paramConfig := ParameterConfig{
		BaseName:      "User",
		PackageName:   "api",
		UsePointers:   true,
		ParameterType: "Query",
	}

	parser := NewParameterParser(builder, paramConfig)

	// Create test parameter
	param := &openapi3.Parameter{
		Name:     "id",
		In:       "query",
		Required: true,
		Schema: &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Type: &openapi3.Types{openapi3.TypeString},
			},
		},
	}

	// Extract parameter
	parser.ExtractParameter(param)

	// Check that statement was added
	stmts := builder.stmts
	if len(stmts) == 0 {
		t.Error("Expected parameter extraction statement to be added")
	}
}

func TestParameterParser_ValidateRequired(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	paramConfig := ParameterConfig{
		BaseName:      "User",
		PackageName:   "api",
		UsePointers:   true,
		ParameterType: "Query",
	}

	parser := NewParameterParser(builder, paramConfig)

	// Create required parameter
	param := &openapi3.Parameter{
		Name:     "id",
		In:       "query",
		Required: true,
		Schema: &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Type: &openapi3.Types{openapi3.TypeString},
			},
		},
	}

	// Validate required
	parser.ValidateRequired(param)

	// Check that validation statement was added
	stmts := builder.stmts
	if len(stmts) == 0 {
		t.Error("Expected validation statement to be added")
	}

	// Check that errors import was added
	imports := builder.imports
	if !imports["github.com/go-faster/errors"] {
		t.Error("Expected errors import to be added for validation")
	}
}

func TestParameterParser_AssignToField(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	paramConfig := ParameterConfig{
		BaseName:      "User",
		PackageName:   "api",
		UsePointers:   true,
		ParameterType: "Query",
		StructName:    "UserQueryParams",
	}

	parser := NewParameterParser(builder, paramConfig)

	// Create test parameter
	param := &openapi3.Parameter{
		Name:     "id",
		In:       "query",
		Required: true,
		Schema: &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Type: &openapi3.Types{openapi3.TypeString},
			},
		},
	}

	// Assign to field
	parser.AssignToField(param)

	// Check that assignment statement was added
	stmts := builder.stmts
	if len(stmts) == 0 {
		t.Error("Expected assignment statement to be added")
	}
}

func TestParameterParser_ValidateStruct(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	paramConfig := ParameterConfig{
		BaseName:      "User",
		PackageName:   "api",
		UsePointers:   true,
		ParameterType: "Query",
		StructName:    "UserQueryParams",
	}

	parser := NewParameterParser(builder, paramConfig)

	// Validate struct
	parser.ValidateStruct()

	// Check that validation statements were added
	stmts := builder.stmts
	if len(stmts) == 0 {
		t.Error("Expected validation statements to be added")
	}

	// Check that validator import was added
	imports := builder.imports
	if !imports["github.com/go-playground/validator/v10"] {
		t.Error("Expected validator import to be added")
	}
}

func TestParameterParser_ReturnResult(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	paramConfig := ParameterConfig{
		BaseName:      "User",
		PackageName:   "api",
		UsePointers:   true,
		ParameterType: "Query",
		StructName:    "UserQueryParams",
	}

	parser := NewParameterParser(builder, paramConfig)

	// Return result
	parser.ReturnResult()

	// Check that return statement was added
	stmts := builder.stmts
	if len(stmts) == 0 {
		t.Error("Expected return statement to be added")
	}
}

func TestParameterParser_FluentInterface(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	paramConfig := ParameterConfig{
		BaseName:    "User",
		PackageName: "api",
		UsePointers: true,
	}

	parser := NewParameterParser(builder, paramConfig)

	// Test fluent interface
	parser = parser.WithStructName("CustomStruct").
		WithPackageName("custom").
		WithUsePointers(false)

	if parser.config.StructName != "CustomStruct" {
		t.Errorf("Expected StructName 'CustomStruct', got '%s'", parser.config.StructName)
	}

	if parser.config.PackageName != "custom" {
		t.Errorf("Expected PackageName 'custom', got '%s'", parser.config.PackageName)
	}

	if parser.config.UsePointers {
		t.Error("Expected UsePointers to be false")
	}
}

func TestParameterParser_GetConfig(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	paramConfig := ParameterConfig{
		BaseName:    "User",
		PackageName: "api",
		UsePointers: true,
	}

	parser := NewParameterParser(builder, paramConfig)

	retrievedConfig := parser.GetConfig()

	if retrievedConfig.BaseName != paramConfig.BaseName {
		t.Error("GetConfig should return the same BaseName")
	}

	if retrievedConfig.PackageName != paramConfig.PackageName {
		t.Error("GetConfig should return the same PackageName")
	}

	if retrievedConfig.UsePointers != paramConfig.UsePointers {
		t.Error("GetConfig should return the same UsePointers value")
	}
}

func TestParameterParser_GetBuilder(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	paramConfig := ParameterConfig{
		BaseName:    "User",
		PackageName: "api",
		UsePointers: true,
	}

	parser := NewParameterParser(builder, paramConfig)

	retrievedBuilder := parser.GetBuilder()

	if retrievedBuilder != builder {
		t.Error("GetBuilder should return the same builder instance")
	}
}

func TestParameterParser_HelperMethods(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	paramConfig := ParameterConfig{
		BaseName:      "User",
		PackageName:   "api",
		UsePointers:   true,
		ParameterType: "Query",
	}

	parser := NewParameterParser(builder, paramConfig)

	// Test parameter var name generation
	param := &openapi3.Parameter{Name: "user_id"}
	varName := parser.getParameterVarName(param)
	expectedVarName := "user_idparam"
	if varName != expectedVarName {
		t.Errorf("Expected var name '%s', got '%s'", expectedVarName, varName)
	}

	// Test field name generation
	fieldName := parser.getFieldName(param)
	expectedFieldName := "User_id"
	if fieldName != expectedFieldName {
		t.Errorf("Expected field name '%s', got '%s'", expectedFieldName, fieldName)
	}

	// Test Go type mapping
	param.Schema = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type: &openapi3.Types{openapi3.TypeString},
		},
	}
	goType := parser.getGoType(param)
	if goType != "string" {
		t.Errorf("Expected Go type 'string', got '%s'", goType)
	}
}
