package astbuilder

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

func TestParameterParser_CompleteWorkflow(t *testing.T) {
	// Test complete parameter parsing workflow
	config := BuilderConfig{
		PackageName:  "api",
		ImportPrefix: "github.com/example/api",
		UsePointers:  true,
	}

	builder := NewBuilder(config)

	// Create parameter parser
	paramConfig := ParameterConfig{
		BaseName:    "User",
		PackageName: "api",
		UsePointers: true,
	}

	parser := NewParameterParser(builder, paramConfig)

	// Create comprehensive test parameters
	params := openapi3.Parameters{
		// Query parameters
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
		{
			Value: &openapi3.Parameter{
				Name:     "active",
				In:       "query",
				Required: false,
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: &openapi3.Types{openapi3.TypeBoolean},
					},
				},
			},
		},
		// Header parameters
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
		{
			Value: &openapi3.Parameter{
				Name:     "Content-Type",
				In:       "header",
				Required: false,
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: &openapi3.Types{openapi3.TypeString},
					},
				},
			},
		},
		// Path parameters
		{
			Value: &openapi3.Parameter{
				Name:     "user_id",
				In:       "path",
				Required: true,
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: &openapi3.Types{openapi3.TypeString},
					},
				},
			},
		},
		// Cookie parameters
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

	// Test query parameter parsing
	err := parser.ParseQueryParams(params)
	if err != nil {
		t.Fatalf("ParseQueryParams failed: %v", err)
	}

	// Test header parameter parsing
	err = parser.ParseHeaders(params)
	if err != nil {
		t.Fatalf("ParseHeaders failed: %v", err)
	}

	// Test path parameter parsing
	err = parser.ParsePathParams(params)
	if err != nil {
		t.Fatalf("ParsePathParams failed: %v", err)
	}

	// Test cookie parameter parsing
	err = parser.ParseCookies(params)
	if err != nil {
		t.Fatalf("ParseCookies failed: %v", err)
	}

	// Verify results
	decls := builder.decls
	if len(decls) < 4 {
		t.Errorf("Expected at least 4 declarations, got %d", len(decls))
	}

	// Check that all required imports are present
	imports := builder.imports
	requiredImports := []string{
		"net/http",
		"github.com/go-chi/chi/v5",
	}

	for _, requiredImport := range requiredImports {
		if !imports[requiredImport] {
			t.Errorf("Expected import '%s' to be present", requiredImport)
		}
	}
}

func TestParameterParser_ChainedWorkflow(t *testing.T) {
	// Test chained parameter parsing workflow
	config := BuilderConfig{
		PackageName:  "api",
		ImportPrefix: "github.com/example/api",
		UsePointers:  true,
	}

	builder := NewBuilder(config)

	// Create parameter parser
	paramConfig := ParameterConfig{
		BaseName:    "User",
		PackageName: "api",
		UsePointers: true,
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

	// Test chained workflow
	parser.DeclareStruct().
		ExtractParameter(param).
		ValidateRequired(param).
		AssignToField(param).
		ValidateStruct().
		ReturnResult()

	// Verify that all operations were performed
	stmts := builder.stmts
	if len(stmts) == 0 {
		t.Error("Expected statements to be added from chained workflow")
	}

	decls := builder.decls
	if len(decls) == 0 {
		t.Error("Expected declarations to be added from chained workflow")
	}
}

func TestParameterParser_EmptyParameters(t *testing.T) {
	// Test handling of empty parameters
	config := BuilderConfig{
		PackageName:  "api",
		ImportPrefix: "github.com/example/api",
		UsePointers:  true,
	}

	builder := NewBuilder(config)

	paramConfig := ParameterConfig{
		BaseName:    "User",
		PackageName: "api",
		UsePointers: true,
	}

	parser := NewParameterParser(builder, paramConfig)

	// Test with empty parameters
	emptyParams := openapi3.Parameters{}

	err := parser.ParseQueryParams(emptyParams)
	if err != nil {
		t.Fatalf("ParseQueryParams with empty params failed: %v", err)
	}

	// Should not add any declarations or statements
	decls := builder.decls
	if len(decls) != 0 {
		t.Errorf("Expected no declarations for empty params, got %d", len(decls))
	}

	stmts := builder.stmts
	if len(stmts) != 0 {
		t.Errorf("Expected no statements for empty params, got %d", len(stmts))
	}
}

func TestParameterParser_MixedParameterTypes(t *testing.T) {
	// Test parsing mixed parameter types
	config := BuilderConfig{
		PackageName:  "api",
		ImportPrefix: "github.com/example/api",
		UsePointers:  true,
	}

	builder := NewBuilder(config)

	paramConfig := ParameterConfig{
		BaseName:    "User",
		PackageName: "api",
		UsePointers: true,
	}

	parser := NewParameterParser(builder, paramConfig)

	// Create mixed parameter types
	params := openapi3.Parameters{
		// Query parameter
		{
			Value: &openapi3.Parameter{
				Name:     "query_param",
				In:       "query",
				Required: true,
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: &openapi3.Types{openapi3.TypeString},
					},
				},
			},
		},
		// Header parameter
		{
			Value: &openapi3.Parameter{
				Name:     "header_param",
				In:       "header",
				Required: false,
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: &openapi3.Types{openapi3.TypeString},
					},
				},
			},
		},
		// Path parameter
		{
			Value: &openapi3.Parameter{
				Name:     "path_param",
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

	// Parse query parameters (should only process query types)
	err := parser.ParseQueryParams(params)
	if err != nil {
		t.Fatalf("ParseQueryParams failed: %v", err)
	}

	// Should only process the query parameter
	decls := builder.decls
	if len(decls) != 1 {
		t.Errorf("Expected 1 declaration for query params, got %d", len(decls))
	}

	// Parse header parameters (should only process header types)
	err = parser.ParseHeaders(params)
	if err != nil {
		t.Fatalf("ParseHeaders failed: %v", err)
	}

	// Should now have 2 declarations (query + header)
	decls = builder.decls
	if len(decls) != 2 {
		t.Errorf("Expected 2 declarations after parsing headers, got %d", len(decls))
	}
}

func TestParameterParser_TypeMapping(t *testing.T) {
	// Test OpenAPI type to Go type mapping
	config := BuilderConfig{
		PackageName:  "api",
		ImportPrefix: "github.com/example/api",
		UsePointers:  true,
	}

	builder := NewBuilder(config)

	paramConfig := ParameterConfig{
		BaseName:    "User",
		PackageName: "api",
		UsePointers: true,
	}

	parser := NewParameterParser(builder, paramConfig)

	// Test different OpenAPI types
	testCases := []struct {
		openapiType    openapi3.Types
		expectedGoType string
	}{
		{openapi3.Types{openapi3.TypeString}, "string"},
		{openapi3.Types{openapi3.TypeInteger}, "int"},
		{openapi3.Types{openapi3.TypeNumber}, "float64"},
		{openapi3.Types{openapi3.TypeBoolean}, "bool"},
		{openapi3.Types{openapi3.TypeArray}, "[]string"},
	}

	for _, tc := range testCases {
		param := &openapi3.Parameter{
			Name: "test_param",
			In:   "query",
			Schema: &openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: &tc.openapiType,
				},
			},
		}

		goType := parser.getGoType(param)
		if goType != tc.expectedGoType {
			t.Errorf("Expected Go type '%s' for OpenAPI type %v, got '%s'",
				tc.expectedGoType, tc.openapiType, goType)
		}
	}
}
