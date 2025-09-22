package astbuilder

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/jolfzverb/codegen/internal/generator"
)

func TestNewSchemaBuilder(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	schemaConfig := SchemaConfig{
		PackageName: "api",
		UsePointers: true,
	}

	schemaBuilder := NewSchemaBuilder(builder, schemaConfig)

	if schemaBuilder == nil {
		t.Fatal("NewSchemaBuilder returned nil")
	}

	if schemaBuilder.config.PackageName != "api" {
		t.Errorf("Expected package name 'api', got '%s'", schemaBuilder.config.PackageName)
	}

	if !schemaBuilder.config.UsePointers {
		t.Error("Expected UsePointers to be true")
	}

	if schemaBuilder.builder != builder {
		t.Error("Expected schema builder to use the provided builder")
	}
}

func TestSchemaBuilder_BuildStruct(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	schemaConfig := SchemaConfig{
		PackageName: "api",
		UsePointers: true,
	}

	schemaBuilder := NewSchemaBuilder(builder, schemaConfig)

	// Create test struct
	model := generator.SchemaStruct{
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

	err := schemaBuilder.BuildStruct(model)
	if err != nil {
		t.Fatalf("BuildStruct failed: %v", err)
	}

	// Check that declaration was added
	decls := builder.decls
	if len(decls) == 0 {
		t.Error("Expected struct declaration to be added")
	}
}

func TestSchemaBuilder_BuildTypeAlias(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	schemaConfig := SchemaConfig{
		PackageName: "api",
		UsePointers: true,
	}

	schemaBuilder := NewSchemaBuilder(builder, schemaConfig)

	err := schemaBuilder.BuildTypeAlias("UserID", "int")
	if err != nil {
		t.Fatalf("BuildTypeAlias failed: %v", err)
	}

	// Check that declaration was added
	decls := builder.decls
	if len(decls) == 0 {
		t.Error("Expected type alias declaration to be added")
	}
}

func TestSchemaBuilder_BuildSliceAlias(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	schemaConfig := SchemaConfig{
		PackageName: "api",
		UsePointers: true,
	}

	schemaBuilder := NewSchemaBuilder(builder, schemaConfig)

	err := schemaBuilder.BuildSliceAlias("UserIDs", "int")
	if err != nil {
		t.Fatalf("BuildSliceAlias failed: %v", err)
	}

	// Check that declaration was added
	decls := builder.decls
	if len(decls) == 0 {
		t.Error("Expected slice alias declaration to be added")
	}
}

func TestSchemaBuilder_BuildField(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	schemaConfig := SchemaConfig{
		PackageName: "api",
		UsePointers: true,
	}

	schemaBuilder := NewSchemaBuilder(builder, schemaConfig)

	field := generator.SchemaField{
		Name:        "Name",
		Type:        "string",
		TagJSON:     []string{"name"},
		TagValidate: []string{"required"},
		Required:    true,
	}

	astField := schemaBuilder.BuildField(field)

	if astField == nil {
		t.Fatal("BuildField returned nil")
	}

	if len(astField.Names) != 1 {
		t.Errorf("Expected 1 field name, got %d", len(astField.Names))
	}

	if astField.Names[0].Name != "Name" {
		t.Errorf("Expected field name 'Name', got '%s'", astField.Names[0].Name)
	}
}

func TestSchemaBuilder_AddStruct(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	schemaConfig := SchemaConfig{
		PackageName: "api",
		UsePointers: true,
	}

	schemaBuilder := NewSchemaBuilder(builder, schemaConfig)

	fields := []generator.SchemaField{
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
	}

	schemaBuilder.AddStruct("User", fields)

	// Check that declaration was added
	decls := builder.decls
	if len(decls) == 0 {
		t.Error("Expected struct declaration to be added")
	}
}

func TestSchemaBuilder_AddField(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	schemaConfig := SchemaConfig{
		PackageName: "api",
		UsePointers: true,
	}

	schemaBuilder := NewSchemaBuilder(builder, schemaConfig)

	tags := map[string]string{
		"json":     "name",
		"validate": "required",
		"required": "true",
	}

	schemaBuilder.AddField("Name", "string", tags)

	// Check that declaration was added
	decls := builder.decls
	if len(decls) == 0 {
		t.Error("Expected struct declaration to be added")
	}
}

func TestSchemaBuilder_AddTypeAlias(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	schemaConfig := SchemaConfig{
		PackageName: "api",
		UsePointers: true,
	}

	schemaBuilder := NewSchemaBuilder(builder, schemaConfig)

	schemaBuilder.AddTypeAlias("UserID", "int")

	// Check that declaration was added
	decls := builder.decls
	if len(decls) == 0 {
		t.Error("Expected type alias declaration to be added")
	}
}

func TestSchemaBuilder_AddSliceType(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	schemaConfig := SchemaConfig{
		PackageName: "api",
		UsePointers: true,
	}

	schemaBuilder := NewSchemaBuilder(builder, schemaConfig)

	schemaBuilder.AddSliceType("UserIDs", "int")

	// Check that declaration was added
	decls := builder.decls
	if len(decls) == 0 {
		t.Error("Expected slice type declaration to be added")
	}
}

func TestSchemaBuilder_BuildFromOpenAPISchema_Object(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	schemaConfig := SchemaConfig{
		PackageName: "api",
		UsePointers: true,
	}

	schemaBuilder := NewSchemaBuilder(builder, schemaConfig)

	// Create OpenAPI schema
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
						Type: &openapi3.Types{openapi3.TypeString},
					},
				},
			},
			Required: []string{"id", "name"},
		},
	}

	err := schemaBuilder.BuildFromOpenAPISchema("User", schema)
	if err != nil {
		t.Fatalf("BuildFromOpenAPISchema failed: %v", err)
	}

	// Check that declaration was added
	decls := builder.decls
	if len(decls) == 0 {
		t.Error("Expected struct declaration to be added")
	}
}

func TestSchemaBuilder_BuildFromOpenAPISchema_Array(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	schemaConfig := SchemaConfig{
		PackageName: "api",
		UsePointers: true,
	}

	schemaBuilder := NewSchemaBuilder(builder, schemaConfig)

	// Create OpenAPI array schema
	schema := &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type: &openapi3.Types{openapi3.TypeArray},
			Items: &openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: &openapi3.Types{openapi3.TypeString},
				},
			},
		},
	}

	err := schemaBuilder.BuildFromOpenAPISchema("UserNames", schema)
	if err != nil {
		t.Fatalf("BuildFromOpenAPISchema failed: %v", err)
	}

	// Check that declaration was added
	decls := builder.decls
	if len(decls) == 0 {
		t.Error("Expected slice type declaration to be added")
	}
}

func TestSchemaBuilder_BuildFromOpenAPISchema_String(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	schemaConfig := SchemaConfig{
		PackageName: "api",
		UsePointers: true,
	}

	schemaBuilder := NewSchemaBuilder(builder, schemaConfig)

	// Create OpenAPI string schema
	schema := &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type: &openapi3.Types{openapi3.TypeString},
		},
	}

	err := schemaBuilder.BuildFromOpenAPISchema("UserName", schema)
	if err != nil {
		t.Fatalf("BuildFromOpenAPISchema failed: %v", err)
	}

	// Check that declaration was added
	decls := builder.decls
	if len(decls) == 0 {
		t.Error("Expected type alias declaration to be added")
	}
}

func TestSchemaBuilder_FluentInterface(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	schemaConfig := SchemaConfig{
		PackageName: "api",
		UsePointers: true,
	}

	schemaBuilder := NewSchemaBuilder(builder, schemaConfig)

	// Test fluent interface
	schemaBuilder = schemaBuilder.WithPackageName("custom").
		WithUsePointers(false).
		WithImportPrefix("github.com/custom")

	if schemaBuilder.config.PackageName != "custom" {
		t.Errorf("Expected package name 'custom', got '%s'", schemaBuilder.config.PackageName)
	}

	if schemaBuilder.config.UsePointers {
		t.Error("Expected UsePointers to be false")
	}

	if schemaBuilder.config.ImportPrefix != "github.com/custom" {
		t.Errorf("Expected import prefix 'github.com/custom', got '%s'", schemaBuilder.config.ImportPrefix)
	}
}

func TestSchemaBuilder_GetConfig(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	schemaConfig := SchemaConfig{
		PackageName: "api",
		UsePointers: true,
	}

	schemaBuilder := NewSchemaBuilder(builder, schemaConfig)

	retrievedConfig := schemaBuilder.GetConfig()

	if retrievedConfig.PackageName != schemaConfig.PackageName {
		t.Error("GetConfig should return the same package name")
	}

	if retrievedConfig.UsePointers != schemaConfig.UsePointers {
		t.Error("GetConfig should return the same UsePointers value")
	}
}

func TestSchemaBuilder_GetBuilder(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	schemaConfig := SchemaConfig{
		PackageName: "api",
		UsePointers: true,
	}

	schemaBuilder := NewSchemaBuilder(builder, schemaConfig)

	retrievedBuilder := schemaBuilder.GetBuilder()

	if retrievedBuilder != builder {
		t.Error("GetBuilder should return the same builder instance")
	}
}

func TestSchemaBuilder_HelperMethods(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	schemaConfig := SchemaConfig{
		PackageName: "api",
		UsePointers: true,
	}

	schemaBuilder := NewSchemaBuilder(builder, schemaConfig)

	// Test PascalCase conversion
	pascalCase := schemaBuilder.toPascalCase("user_name")
	expected := "UserName"
	if pascalCase != expected {
		t.Errorf("Expected PascalCase 'UserName', got '%s'", pascalCase)
	}

	// Test Go type mapping
	goType := schemaBuilder.getGoType("string")
	if goType != "string" {
		t.Errorf("Expected Go type 'string', got '%s'", goType)
	}

	goType = schemaBuilder.getGoType("[]string")
	if goType != "[]string" {
		t.Errorf("Expected Go type '[]string', got '%s'", goType)
	}
}

func TestSchemaBuilder_ErrorHandling(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	schemaConfig := SchemaConfig{
		PackageName: "api",
		UsePointers: true,
	}

	schemaBuilder := NewSchemaBuilder(builder, schemaConfig)

	// Test empty struct name
	err := schemaBuilder.BuildStruct(generator.SchemaStruct{Name: ""})
	if err == nil {
		t.Error("Expected error for empty struct name")
	}

	// Test empty type alias name
	err = schemaBuilder.BuildTypeAlias("", "int")
	if err == nil {
		t.Error("Expected error for empty type alias name")
	}

	// Test empty slice alias name
	err = schemaBuilder.BuildSliceAlias("", "int")
	if err == nil {
		t.Error("Expected error for empty slice alias name")
	}

	// Test nil OpenAPI schema
	err = schemaBuilder.BuildFromOpenAPISchema("Test", nil)
	if err == nil {
		t.Error("Expected error for nil OpenAPI schema")
	}
}
