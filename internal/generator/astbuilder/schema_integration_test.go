package astbuilder

/*
import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/jolfzverb/codegen/internal/generator"
)

func TestSchemaBuilder_CompleteWorkflow(t *testing.T) {
	// Test complete schema building workflow
	config := BuilderConfig{
		PackageName:  "api",
		ImportPrefix: "github.com/example/api",
		UsePointers:  true,
	}

	builder := NewBuilder(config)

	schemaConfig := SchemaConfig{
		PackageName: "api",
		UsePointers: true,
	}

	schemaBuilder := NewSchemaBuilder(builder, schemaConfig)

	// Test 1: Build struct from generator.SchemaStruct
	userStruct := generator.SchemaStruct{
		Name: "User",
		Fields: []generator.SchemaField{
			{
				Name:        "ID",
				Type:        "int",
				TagJSON:     []string{"id"},
				TagValidate: []string{"required"},
				Required:    true,
			},
			{
				Name:        "Name",
				Type:        "string",
				TagJSON:     []string{"name"},
				TagValidate: []string{"required", "min=1"},
				Required:    true,
			},
			{
				Name:        "Email",
				Type:        "string",
				TagJSON:     []string{"email"},
				TagValidate: []string{"required", "email"},
				Required:    true,
			},
			{
				Name:        "Age",
				Type:        "int",
				TagJSON:     []string{"age"},
				TagValidate: []string{},
				Required:    false,
			},
		},
	}

	err := schemaBuilder.BuildStruct(userStruct)
	if err != nil {
		t.Fatalf("BuildStruct failed: %v", err)
	}

	// Test 2: Build type aliases
	err = schemaBuilder.BuildTypeAlias("UserID", "int")
	if err != nil {
		t.Fatalf("BuildTypeAlias failed: %v", err)
	}

	err = schemaBuilder.BuildSliceAlias("UserIDs", "int")
	if err != nil {
		t.Fatalf("BuildSliceAlias failed: %v", err)
	}

	// Test 3: Build from OpenAPI schema
	openapiSchema := &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type: &openapi3.Types{openapi3.TypeObject},
			Properties: openapi3.Schemas{
				"id": &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: &openapi3.Types{openapi3.TypeInteger},
					},
				},
				"name": &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: &openapi3.Types{openapi3.TypeString},
					},
				},
				"email": &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type:   &openapi3.Types{openapi3.TypeString},
						Format: "email",
					},
				},
			},
			Required: []string{"id", "name", "email"},
		},
	}

	err = schemaBuilder.BuildFromOpenAPISchema("OpenAPIUser", openapiSchema)
	if err != nil {
		t.Fatalf("BuildFromOpenAPISchema failed: %v", err)
	}

	// Verify results
	decls := builder.decls
	if len(decls) < 4 {
		t.Errorf("Expected at least 4 declarations, got %d", len(decls))
	}
}

func TestSchemaBuilder_ChainedWorkflow(t *testing.T) {
	// Test chained schema building workflow
	config := BuilderConfig{
		PackageName:  "api",
		ImportPrefix: "github.com/example/api",
		UsePointers:  true,
	}

	builder := NewBuilder(config)

	schemaConfig := SchemaConfig{
		PackageName: "api",
		UsePointers: true,
	}

	schemaBuilder := NewSchemaBuilder(builder, schemaConfig)

	// Test chained workflow
	schemaBuilder.
		AddTypeAlias("UserID", "int").
		AddSliceType("UserIDs", "int").
		AddStruct("User", []generator.SchemaField{
			{
				Name:        "ID",
				Type:        "int",
				TagJSON:     []string{"id"},
				TagValidate: []string{"required"},
				Required:    true,
			},
			{
				Name:        "Name",
				Type:        "string",
				TagJSON:     []string{"name"},
				TagValidate: []string{"required"},
				Required:    true,
			},
		}).
		AddField("Email", "string", map[string]string{
			"json":     "email",
			"validate": "required,email",
		})

	// Verify that all operations were performed
	decls := builder.decls
	if len(decls) < 4 {
		t.Errorf("Expected at least 4 declarations from chained workflow, got %d", len(decls))
	}
}

func TestSchemaBuilder_OpenAPITypeMapping(t *testing.T) {
	// Test OpenAPI type to Go type mapping
	config := BuilderConfig{
		PackageName:  "api",
		ImportPrefix: "github.com/example/api",
		UsePointers:  true,
	}

	builder := NewBuilder(config)

	schemaConfig := SchemaConfig{
		PackageName: "api",
		UsePointers: true,
	}

	schemaBuilder := NewSchemaBuilder(builder, schemaConfig)

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
		schema := &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Type: &tc.openapiType,
			},
		}

		err := schemaBuilder.BuildFromOpenAPISchema("TestType", schema)
		if err != nil {
			t.Fatalf("BuildFromOpenAPISchema failed for type %v: %v", tc.openapiType, err)
		}
	}

	// Verify that all type declarations were added
	decls := builder.decls
	if len(decls) != len(testCases) {
		t.Errorf("Expected %d declarations, got %d", len(testCases), len(decls))
	}
}

func TestSchemaBuilder_ValidationTags(t *testing.T) {
	// Test validation tag generation
	config := BuilderConfig{
		PackageName:  "api",
		ImportPrefix: "github.com/example/api",
		UsePointers:  true,
	}

	builder := NewBuilder(config)

	schemaConfig := SchemaConfig{
		PackageName: "api",
		UsePointers: true,
	}

	schemaBuilder := NewSchemaBuilder(builder, schemaConfig)

	// Test schema with validation constraints
	schema := &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type:      &openapi3.Types{openapi3.TypeString},
			MinLength: 1,
			MaxLength: func() *uint64 { v := uint64(100); return &v }(),
			Pattern:   "^[a-zA-Z0-9]+$",
			Format:    "email",
		},
	}

	err := schemaBuilder.BuildFromOpenAPISchema("ValidatedString", schema)
	if err != nil {
		t.Fatalf("BuildFromOpenAPISchema failed: %v", err)
	}

	// Verify that declaration was added
	decls := builder.decls
	if len(decls) == 0 {
		t.Error("Expected declaration to be added")
	}
}

func TestSchemaBuilder_ComplexObject(t *testing.T) {
	// Test building a complex object with nested properties
	config := BuilderConfig{
		PackageName:  "api",
		ImportPrefix: "github.com/example/api",
		UsePointers:  true,
	}

	builder := NewBuilder(config)

	schemaConfig := SchemaConfig{
		PackageName: "api",
		UsePointers: true,
	}

	schemaBuilder := NewSchemaBuilder(builder, schemaConfig)

	// Create complex OpenAPI schema
	schema := &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type: &openapi3.Types{openapi3.TypeObject},
			Properties: openapi3.Schemas{
				"id": &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: &openapi3.Types{openapi3.TypeInteger},
					},
				},
				"name": &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type:      &openapi3.Types{openapi3.TypeString},
						MinLength: 1,
						MaxLength: func() *uint64 { v := uint64(100); return &v }(),
					},
				},
				"email": &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type:   &openapi3.Types{openapi3.TypeString},
						Format: "email",
					},
				},
				"age": &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: &openapi3.Types{openapi3.TypeInteger},
						Min:  func() *float64 { v := float64(0); return &v }(),
						Max:  func() *float64 { v := float64(150); return &v }(),
					},
				},
				"tags": &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: &openapi3.Types{openapi3.TypeArray},
						Items: &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type: &openapi3.Types{openapi3.TypeString},
							},
						},
					},
				},
			},
			Required: []string{"id", "name", "email"},
		},
	}

	err := schemaBuilder.BuildFromOpenAPISchema("ComplexUser", schema)
	if err != nil {
		t.Fatalf("BuildFromOpenAPISchema failed: %v", err)
	}

	// Verify that declaration was added
	decls := builder.decls
	if len(decls) == 0 {
		t.Error("Expected declaration to be added")
	}
}

func TestSchemaBuilder_PointerHandling(t *testing.T) {
	// Test pointer handling for optional fields
	config := BuilderConfig{
		PackageName:  "api",
		ImportPrefix: "github.com/example/api",
		UsePointers:  true,
	}

	builder := NewBuilder(config)

	schemaConfig := SchemaConfig{
		PackageName: "api",
		UsePointers: true,
	}

	schemaBuilder := NewSchemaBuilder(builder, schemaConfig)

	// Create struct with optional fields
	model := generator.SchemaStruct{
		Name: "UserWithOptionals",
		Fields: []generator.SchemaField{
			{
				Name:        "ID",
				Type:        "int",
				TagJSON:     []string{"id"},
				TagValidate: []string{"required"},
				Required:    true,
			},
			{
				Name:        "OptionalName",
				Type:        "string",
				TagJSON:     []string{"optional_name"},
				TagValidate: []string{},
				Required:    false,
			},
		},
	}

	err := schemaBuilder.BuildStruct(model)
	if err != nil {
		t.Fatalf("BuildStruct failed: %v", err)
	}

	// Verify that declaration was added
	decls := builder.decls
	if len(decls) == 0 {
		t.Error("Expected declaration to be added")
	}
}

func TestSchemaBuilder_EmptyStruct(t *testing.T) {
	// Test building an empty struct
	config := BuilderConfig{
		PackageName:  "api",
		ImportPrefix: "github.com/example/api",
		UsePointers:  true,
	}

	builder := NewBuilder(config)

	schemaConfig := SchemaConfig{
		PackageName: "api",
		UsePointers: true,
	}

	schemaBuilder := NewSchemaBuilder(builder, schemaConfig)

	// Create empty struct
	model := generator.SchemaStruct{
		Name:   "EmptyStruct",
		Fields: []generator.SchemaField{},
	}

	err := schemaBuilder.BuildStruct(model)
	if err != nil {
		t.Fatalf("BuildStruct failed: %v", err)
	}

	// Verify that declaration was added
	decls := builder.decls
	if len(decls) == 0 {
		t.Error("Expected declaration to be added")
	}
}
*/
