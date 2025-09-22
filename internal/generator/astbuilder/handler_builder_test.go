package astbuilder

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

func TestNewHandlerBuilder(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	handlerConfig := HandlerConfig{
		PackageName: "api",
		UsePointers: true,
	}

	handlerBuilder := NewHandlerBuilder(builder, handlerConfig)

	if handlerBuilder == nil {
		t.Fatal("NewHandlerBuilder returned nil")
	}

	if handlerBuilder.config.PackageName != "api" {
		t.Errorf("Expected package name 'api', got '%s'", handlerBuilder.config.PackageName)
	}

	if !handlerBuilder.config.UsePointers {
		t.Error("Expected UsePointers to be true")
	}

	if handlerBuilder.builder != builder {
		t.Error("Expected handler builder to use the provided builder")
	}
}

func TestHandlerBuilder_BuildHandlerStruct(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	handlerConfig := HandlerConfig{
		PackageName: "api",
		UsePointers: true,
	}

	handlerBuilder := NewHandlerBuilder(builder, handlerConfig)

	err := handlerBuilder.BuildHandlerStruct()
	if err != nil {
		t.Fatalf("BuildHandlerStruct failed: %v", err)
	}

	// Check that declaration was added
	decls := builder.decls
	if len(decls) == 0 {
		t.Error("Expected handler struct declaration to be added")
	}

	// Check that validator import was added
	imports := builder.imports
	if !imports["github.com/go-playground/validator/v10"] {
		t.Error("Expected validator import to be added")
	}
}

func TestHandlerBuilder_BuildConstructor(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	handlerConfig := HandlerConfig{
		PackageName: "api",
		UsePointers: true,
	}

	handlerBuilder := NewHandlerBuilder(builder, handlerConfig)

	err := handlerBuilder.BuildConstructor()
	if err != nil {
		t.Fatalf("BuildConstructor failed: %v", err)
	}

	// Check that declaration was added
	decls := builder.decls
	if len(decls) == 0 {
		t.Error("Expected constructor declaration to be added")
	}

	// Check that validator import was added
	imports := builder.imports
	if !imports["github.com/go-playground/validator/v10"] {
		t.Error("Expected validator import to be added")
	}
}

func TestHandlerBuilder_BuildInterface(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	handlerConfig := HandlerConfig{
		PackageName: "api",
		UsePointers: true,
	}

	handlerBuilder := NewHandlerBuilder(builder, handlerConfig)

	// Create test methods
	methods := []MethodSpec{
		{
			Name: "GetUser",
			Params: []FieldSpec{
				{Name: "id", Type: "string"},
			},
			Returns: []FieldSpec{
				{Name: "user", Type: "User"},
				{Name: "error", Type: "error"},
			},
		},
		{
			Name: "CreateUser",
			Params: []FieldSpec{
				{Name: "user", Type: "User"},
			},
			Returns: []FieldSpec{
				{Name: "error", Type: "error"},
			},
		},
	}

	err := handlerBuilder.BuildInterface("UserHandler", methods)
	if err != nil {
		t.Fatalf("BuildInterface failed: %v", err)
	}

	// Check that declaration was added
	decls := builder.decls
	if len(decls) == 0 {
		t.Error("Expected interface declaration to be added")
	}
}

func TestHandlerBuilder_BuildRoutesFunction(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	handlerConfig := HandlerConfig{
		PackageName: "api",
		UsePointers: true,
	}

	handlerBuilder := NewHandlerBuilder(builder, handlerConfig)

	err := handlerBuilder.BuildRoutesFunction()
	if err != nil {
		t.Fatalf("BuildRoutesFunction failed: %v", err)
	}

	// Check that declaration was added
	decls := builder.decls
	if len(decls) == 0 {
		t.Error("Expected routes function declaration to be added")
	}

	// Check that chi import was added
	imports := builder.imports
	if !imports["github.com/go-chi/chi/v5"] {
		t.Error("Expected chi import to be added")
	}
}

func TestHandlerBuilder_AddHandlerField(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	handlerConfig := HandlerConfig{
		PackageName: "api",
		UsePointers: true,
	}

	handlerBuilder := NewHandlerBuilder(builder, handlerConfig)

	handlerBuilder.AddHandlerField("logger", "*log.Logger")

	// Check that declaration was added
	decls := builder.decls
	if len(decls) == 0 {
		t.Error("Expected handler struct declaration to be added")
	}
}

func TestHandlerBuilder_AddInterfaceMethod(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	handlerConfig := HandlerConfig{
		PackageName: "api",
		UsePointers: true,
	}

	handlerBuilder := NewHandlerBuilder(builder, handlerConfig)

	params := []FieldSpec{
		{Name: "id", Type: "string"},
	}
	returns := []FieldSpec{
		{Name: "user", Type: "User"},
		{Name: "error", Type: "error"},
	}

	handlerBuilder.AddInterfaceMethod("GetUser", params, returns)

	// Check that declaration was added
	decls := builder.decls
	if len(decls) == 0 {
		t.Error("Expected interface declaration to be added")
	}
}

func TestHandlerBuilder_AddRoute(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	handlerConfig := HandlerConfig{
		PackageName: "api",
		UsePointers: true,
	}

	handlerBuilder := NewHandlerBuilder(builder, handlerConfig)

	handlerBuilder.AddRoute("GET", "/users/{id}", "GetUser")

	// Check that statement was added
	stmts := builder.stmts
	if len(stmts) == 0 {
		t.Error("Expected route statement to be added")
	}

	// Check that chi import was added
	imports := builder.imports
	if !imports["github.com/go-chi/chi/v5"] {
		t.Error("Expected chi import to be added")
	}
}

func TestHandlerBuilder_AddResponseWriter(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	handlerConfig := HandlerConfig{
		PackageName: "api",
		UsePointers: true,
	}

	handlerBuilder := NewHandlerBuilder(builder, handlerConfig)

	codes := []string{"200", "400", "404", "500"}
	handlerBuilder.AddResponseWriter("Write", codes)

	// Check that declarations were added
	decls := builder.decls
	if len(decls) == 0 {
		t.Error("Expected response writer declarations to be added")
	}

	// Check that http import was added
	imports := builder.imports
	if !imports["net/http"] {
		t.Error("Expected http import to be added")
	}
}

func TestHandlerBuilder_BuildFromOpenAPI(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	handlerConfig := HandlerConfig{
		PackageName: "api",
		UsePointers: true,
	}

	handlerBuilder := NewHandlerBuilder(builder, handlerConfig)

	// Create test OpenAPI specification
	paths := openapi3.Paths{}
	paths.Set("/users", &openapi3.PathItem{
		Get: &openapi3.Operation{
			OperationID: "getUsers",
			Summary:     "Get all users",
		},
		Post: &openapi3.Operation{
			OperationID: "createUser",
			Summary:     "Create a new user",
		},
	})
	paths.Set("/users/{id}", &openapi3.PathItem{
		Get: &openapi3.Operation{
			OperationID: "getUser",
			Summary:     "Get user by ID",
		},
	})

	spec := &openapi3.T{
		Paths: &paths,
	}

	err := handlerBuilder.BuildFromOpenAPI(spec)
	if err != nil {
		t.Fatalf("BuildFromOpenAPI failed: %v", err)
	}

	// Check that declarations were added
	decls := builder.decls
	if len(decls) == 0 {
		t.Error("Expected declarations to be added")
	}

	// Check that required imports were added
	imports := builder.imports
	requiredImports := []string{
		"github.com/go-playground/validator/v10",
		"github.com/go-chi/chi/v5",
		"net/http",
	}

	for _, requiredImport := range requiredImports {
		if !imports[requiredImport] {
			t.Errorf("Expected import '%s' to be present", requiredImport)
		}
	}
}

func TestHandlerBuilder_FluentInterface(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	handlerConfig := HandlerConfig{
		PackageName: "api",
		UsePointers: true,
	}

	handlerBuilder := NewHandlerBuilder(builder, handlerConfig)

	// Test fluent interface
	handlerBuilder = handlerBuilder.WithPackageName("custom").
		WithUsePointers(false).
		WithImportPrefix("github.com/custom")

	if handlerBuilder.config.PackageName != "custom" {
		t.Errorf("Expected package name 'custom', got '%s'", handlerBuilder.config.PackageName)
	}

	if handlerBuilder.config.UsePointers {
		t.Error("Expected UsePointers to be false")
	}

	if handlerBuilder.config.ImportPrefix != "github.com/custom" {
		t.Errorf("Expected import prefix 'github.com/custom', got '%s'", handlerBuilder.config.ImportPrefix)
	}
}

func TestHandlerBuilder_GetConfig(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	handlerConfig := HandlerConfig{
		PackageName: "api",
		UsePointers: true,
	}

	handlerBuilder := NewHandlerBuilder(builder, handlerConfig)

	retrievedConfig := handlerBuilder.GetConfig()

	if retrievedConfig.PackageName != handlerConfig.PackageName {
		t.Error("GetConfig should return the same package name")
	}

	if retrievedConfig.UsePointers != handlerConfig.UsePointers {
		t.Error("GetConfig should return the same UsePointers value")
	}
}

func TestHandlerBuilder_GetBuilder(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	handlerConfig := HandlerConfig{
		PackageName: "api",
		UsePointers: true,
	}

	handlerBuilder := NewHandlerBuilder(builder, handlerConfig)

	retrievedBuilder := handlerBuilder.GetBuilder()

	if retrievedBuilder != builder {
		t.Error("GetBuilder should return the same builder instance")
	}
}

func TestHandlerBuilder_HelperMethods(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	handlerConfig := HandlerConfig{
		PackageName: "api",
		UsePointers: true,
	}

	handlerBuilder := NewHandlerBuilder(builder, handlerConfig)

	// Test PascalCase conversion
	pascalCase := handlerBuilder.toPascalCase("get_user_by_id")
	expected := "GetUserById"
	if pascalCase != expected {
		t.Errorf("Expected PascalCase 'GetUserById', got '%s'", pascalCase)
	}

	// Test handler name generation
	handlerName := handlerBuilder.generateHandlerName("getUser", "GET", "/users/{id}")
	expected = "GetUser"
	if handlerName != expected {
		t.Errorf("Expected handler name 'GetUser', got '%s'", handlerName)
	}

	// Test handler name generation without operation ID
	handlerName = handlerBuilder.generateHandlerName("", "POST", "/users")
	expected = "POSTUsers"
	if handlerName != expected {
		t.Errorf("Expected handler name 'POSTUsers', got '%s'", handlerName)
	}
}

func TestHandlerBuilder_ErrorHandling(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	handlerConfig := HandlerConfig{
		PackageName: "api",
		UsePointers: true,
	}

	handlerBuilder := NewHandlerBuilder(builder, handlerConfig)

	// Test empty interface name
	err := handlerBuilder.BuildInterface("", []MethodSpec{})
	if err == nil {
		t.Error("Expected error for empty interface name")
	}

	// Test nil OpenAPI specification
	err = handlerBuilder.BuildFromOpenAPI(nil)
	if err == nil {
		t.Error("Expected error for nil OpenAPI specification")
	}
}
