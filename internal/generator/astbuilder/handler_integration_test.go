package astbuilder

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

func TestHandlerBuilder_CompleteWorkflow(t *testing.T) {
	// Test complete handler building workflow
	config := BuilderConfig{
		PackageName:  "api",
		ImportPrefix: "github.com/example/api",
		UsePointers:  true,
	}

	builder := NewBuilder(config)

	handlerConfig := HandlerConfig{
		PackageName: "api",
		UsePointers: true,
	}

	handlerBuilder := NewHandlerBuilder(builder, handlerConfig)

	// Test 1: Build handler struct
	err := handlerBuilder.BuildHandlerStruct()
	if err != nil {
		t.Fatalf("BuildHandlerStruct failed: %v", err)
	}

	// Test 2: Build constructor
	err = handlerBuilder.BuildConstructor()
	if err != nil {
		t.Fatalf("BuildConstructor failed: %v", err)
	}

	// Test 3: Build interface
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

	err = handlerBuilder.BuildInterface("UserHandler", methods)
	if err != nil {
		t.Fatalf("BuildInterface failed: %v", err)
	}

	// Test 4: Build routes function
	err = handlerBuilder.BuildRoutesFunction()
	if err != nil {
		t.Fatalf("BuildRoutesFunction failed: %v", err)
	}

	// Verify results
	decls := builder.decls
	if len(decls) < 4 {
		t.Errorf("Expected at least 4 declarations, got %d", len(decls))
	}

	// Check that all required imports are present
	imports := builder.imports
	requiredImports := []string{
		"github.com/go-playground/validator/v10",
		"github.com/go-chi/chi/v5",
	}

	for _, requiredImport := range requiredImports {
		if !imports[requiredImport] {
			t.Errorf("Expected import '%s' to be present", requiredImport)
		}
	}
}

func TestHandlerBuilder_ChainedWorkflow(t *testing.T) {
	// Test chained handler building workflow
	config := BuilderConfig{
		PackageName:  "api",
		ImportPrefix: "github.com/example/api",
		UsePointers:  true,
	}

	builder := NewBuilder(config)

	handlerConfig := HandlerConfig{
		PackageName: "api",
		UsePointers: true,
	}

	handlerBuilder := NewHandlerBuilder(builder, handlerConfig)

	// Test chained workflow
	handlerBuilder.
		AddHandlerField("logger", "*log.Logger").
		AddInterfaceMethod("GetUser", []FieldSpec{{Name: "id", Type: "string"}}, []FieldSpec{{Name: "user", Type: "User"}, {Name: "error", Type: "error"}}).
		AddRoute("GET", "/users/{id}", "GetUser").
		AddResponseWriter("Write", []string{"200", "404"})

	// Verify that all operations were performed
	decls := builder.decls
	if len(decls) < 4 {
		t.Errorf("Expected at least 4 declarations from chained workflow, got %d", len(decls))
	}
}

func TestHandlerBuilder_OpenAPIIntegration(t *testing.T) {
	// Test OpenAPI integration
	config := BuilderConfig{
		PackageName:  "api",
		ImportPrefix: "github.com/example/api",
		UsePointers:  true,
	}

	builder := NewBuilder(config)

	handlerConfig := HandlerConfig{
		PackageName: "api",
		UsePointers: true,
	}

	handlerBuilder := NewHandlerBuilder(builder, handlerConfig)

	// Create comprehensive OpenAPI specification
	paths := openapi3.Paths{}
	paths.Set("/users", &openapi3.PathItem{
		Get: &openapi3.Operation{
			OperationID: "getUsers",
			Summary:     "Get all users",
			Tags:        []string{"users"},
		},
		Post: &openapi3.Operation{
			OperationID: "createUser",
			Summary:     "Create a new user",
			Tags:        []string{"users"},
		},
	})
	paths.Set("/users/{id}", &openapi3.PathItem{
		Get: &openapi3.Operation{
			OperationID: "getUser",
			Summary:     "Get user by ID",
			Tags:        []string{"users"},
		},
		Put: &openapi3.Operation{
			OperationID: "updateUser",
			Summary:     "Update user by ID",
			Tags:        []string{"users"},
		},
		Delete: &openapi3.Operation{
			OperationID: "deleteUser",
			Summary:     "Delete user by ID",
			Tags:        []string{"users"},
		},
	})
	paths.Set("/users/{id}/posts", &openapi3.PathItem{
		Get: &openapi3.Operation{
			OperationID: "getUserPosts",
			Summary:     "Get user posts",
			Tags:        []string{"users", "posts"},
		},
	})

	spec := &openapi3.T{
		Info: &openapi3.Info{
			Title:   "User API",
			Version: "1.0.0",
		},
		Paths: &paths,
	}

	err := handlerBuilder.BuildFromOpenAPI(spec)
	if err != nil {
		t.Fatalf("BuildFromOpenAPI failed: %v", err)
	}

	// Verify that declarations were added
	decls := builder.decls
	if len(decls) == 0 {
		t.Error("Expected declarations to be added")
	}

	// Check that all required imports are present
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

func TestHandlerBuilder_ResponseWriterGeneration(t *testing.T) {
	// Test response writer generation
	config := BuilderConfig{
		PackageName:  "api",
		ImportPrefix: "github.com/example/api",
		UsePointers:  true,
	}

	builder := NewBuilder(config)

	handlerConfig := HandlerConfig{
		PackageName: "api",
		UsePointers: true,
	}

	handlerBuilder := NewHandlerBuilder(builder, handlerConfig)

	// Test different status codes
	codes := []string{"200", "201", "400", "401", "403", "404", "422", "500", "502", "503"}
	handlerBuilder.AddResponseWriter("Write", codes)

	// Verify that declarations were added
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

func TestHandlerBuilder_InterfaceGeneration(t *testing.T) {
	// Test interface generation with various method signatures
	config := BuilderConfig{
		PackageName:  "api",
		ImportPrefix: "github.com/example/api",
		UsePointers:  true,
	}

	builder := NewBuilder(config)

	handlerConfig := HandlerConfig{
		PackageName: "api",
		UsePointers: true,
	}

	handlerBuilder := NewHandlerBuilder(builder, handlerConfig)

	// Test different method signatures
	methods := []MethodSpec{
		{
			Name:    "GetUser",
			Params:  []FieldSpec{{Name: "id", Type: "string"}},
			Returns: []FieldSpec{{Name: "user", Type: "User"}, {Name: "error", Type: "error"}},
		},
		{
			Name:    "CreateUser",
			Params:  []FieldSpec{{Name: "user", Type: "User"}},
			Returns: []FieldSpec{{Name: "error", Type: "error"}},
		},
		{
			Name:    "UpdateUser",
			Params:  []FieldSpec{{Name: "id", Type: "string"}, {Name: "user", Type: "User"}},
			Returns: []FieldSpec{{Name: "error", Type: "error"}},
		},
		{
			Name:    "DeleteUser",
			Params:  []FieldSpec{{Name: "id", Type: "string"}},
			Returns: []FieldSpec{{Name: "error", Type: "error"}},
		},
		{
			Name:    "ListUsers",
			Params:  []FieldSpec{{Name: "limit", Type: "int"}, {Name: "offset", Type: "int"}},
			Returns: []FieldSpec{{Name: "users", Type: "[]User"}, {Name: "error", Type: "error"}},
		},
	}

	err := handlerBuilder.BuildInterface("UserService", methods)
	if err != nil {
		t.Fatalf("BuildInterface failed: %v", err)
	}

	// Verify that declaration was added
	decls := builder.decls
	if len(decls) == 0 {
		t.Error("Expected interface declaration to be added")
	}
}

func TestHandlerBuilder_RouteGeneration(t *testing.T) {
	// Test route generation
	config := BuilderConfig{
		PackageName:  "api",
		ImportPrefix: "github.com/example/api",
		UsePointers:  true,
	}

	builder := NewBuilder(config)

	handlerConfig := HandlerConfig{
		PackageName: "api",
		UsePointers: true,
	}

	handlerBuilder := NewHandlerBuilder(builder, handlerConfig)

	// Test different HTTP methods and paths
	routes := []struct {
		method      string
		path        string
		handlerName string
	}{
		{"GET", "/users", "GetUsers"},
		{"POST", "/users", "CreateUser"},
		{"GET", "/users/{id}", "GetUser"},
		{"PUT", "/users/{id}", "UpdateUser"},
		{"DELETE", "/users/{id}", "DeleteUser"},
		{"GET", "/users/{id}/posts", "GetUserPosts"},
		{"POST", "/users/{id}/posts", "CreateUserPost"},
	}

	for _, route := range routes {
		handlerBuilder.AddRoute(route.method, route.path, route.handlerName)
	}

	// Verify that declarations were added
	decls := builder.decls
	if len(decls) == 0 {
		t.Error("Expected route declarations to be added")
	}

	// Check that chi import was added
	imports := builder.imports
	if !imports["github.com/go-chi/chi/v5"] {
		t.Error("Expected chi import to be added")
	}
}

func TestHandlerBuilder_HandlerNameGeneration(t *testing.T) {
	// Test handler name generation
	config := BuilderConfig{
		PackageName:  "api",
		ImportPrefix: "github.com/example/api",
		UsePointers:  true,
	}

	builder := NewBuilder(config)

	handlerConfig := HandlerConfig{
		PackageName: "api",
		UsePointers: true,
	}

	handlerBuilder := NewHandlerBuilder(builder, handlerConfig)

	// Test cases for handler name generation
	testCases := []struct {
		operationID string
		method      string
		path        string
		expected    string
	}{
		{"getUser", "GET", "/users/{id}", "GetUser"},
		{"", "GET", "/users", "GETUsers"},
		{"", "POST", "/users", "POSTUsers"},
		{"", "PUT", "/users/{id}", "PUTUsers"},
		{"", "DELETE", "/users/{id}", "DELETEUsers"},
		{"", "GET", "/users/{id}/posts", "GETUsersPosts"},
		{"", "POST", "/users/{id}/posts", "POSTUsersPosts"},
	}

	for _, tc := range testCases {
		handlerName := handlerBuilder.generateHandlerName(tc.operationID, tc.method, tc.path)
		if handlerName != tc.expected {
			t.Errorf("Expected handler name '%s' for operationID='%s', method='%s', path='%s', got '%s'",
				tc.expected, tc.operationID, tc.method, tc.path, handlerName)
		}
	}
}

func TestHandlerBuilder_EmptySpecification(t *testing.T) {
	// Test handling of empty OpenAPI specification
	config := BuilderConfig{
		PackageName:  "api",
		ImportPrefix: "github.com/example/api",
		UsePointers:  true,
	}

	builder := NewBuilder(config)

	handlerConfig := HandlerConfig{
		PackageName: "api",
		UsePointers: true,
	}

	handlerBuilder := NewHandlerBuilder(builder, handlerConfig)

	// Test with empty specification
	emptyPaths := openapi3.Paths{}
	emptySpec := &openapi3.T{
		Paths: &emptyPaths,
	}

	err := handlerBuilder.BuildFromOpenAPI(emptySpec)
	if err != nil {
		t.Fatalf("BuildFromOpenAPI with empty spec failed: %v", err)
	}

	// Should still add basic handler components
	decls := builder.decls
	if len(decls) < 3 {
		t.Errorf("Expected at least 3 declarations for empty spec, got %d", len(decls))
	}
}
