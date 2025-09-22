package astbuilder

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// HandlerConfig holds configuration for handler building
type HandlerConfig struct {
	PackageName  string
	UsePointers  bool
	ImportPrefix string
}

// MethodSpec represents a method specification
type MethodSpec struct {
	Name    string
	Params  []FieldSpec
	Returns []FieldSpec
}

// FieldSpec represents a field specification
type FieldSpec struct {
	Name string
	Type string
}

// RouteSpec represents a route specification
type RouteSpec struct {
	Method      string
	Path        string
	HandlerName string
}

// HandlerBuilder provides high-level methods for building Go handlers
type HandlerBuilder struct {
	builder *Builder
	config  HandlerConfig
}

// NewHandlerBuilder creates a new handler builder
func NewHandlerBuilder(builder *Builder, config HandlerConfig) *HandlerBuilder {
	return &HandlerBuilder{
		builder: builder,
		config:  config,
	}
}

// BuildHandlerStruct builds the main handler struct
func (h *HandlerBuilder) BuildHandlerStruct() error {
	typeBuilder := NewTypeBuilder(h.builder)
	exprBuilder := NewExpressionBuilder(h.builder)

	// Create handler struct fields
	fields := []*ast.Field{
		typeBuilder.Field("validator", exprBuilder.Star(exprBuilder.Select(exprBuilder.Ident("validator"), "Validate")), ""),
	}

	// Create handler struct declaration
	handlerDecl := typeBuilder.StructAlias("Handler", fields)
	h.builder.AddDeclaration(handlerDecl)

	// Add validator import
	h.builder.AddImport("github.com/go-playground/validator/v10")

	return nil
}

// BuildConstructor builds the handler constructor
func (h *HandlerBuilder) BuildConstructor() error {
	funcBuilder := NewFunctionBuilder(h.builder)
	exprBuilder := NewExpressionBuilder(h.builder)
	stmtBuilder := NewStatementBuilder(h.builder)

	// Create constructor parameters
	params := []*ast.Field{
		funcBuilder.Param("validator", "*validator.Validate"),
	}

	// Create constructor results
	results := []*ast.Field{
		funcBuilder.ResultAnonymous("*Handler"),
	}

	// Create constructor body
	body := []ast.Stmt{
		stmtBuilder.Return(
			exprBuilder.AddressOf(
				exprBuilder.CompositeLitWithType(
					exprBuilder.Ident("Handler"),
					exprBuilder.KeyValue(exprBuilder.Ident("validator"), exprBuilder.Ident("validator")),
				),
			),
		),
	}

	// Create constructor declaration
	constructorDecl := funcBuilder.Function("NewHandler", params, results, body)
	h.builder.AddDeclaration(constructorDecl)

	// Add validator import
	h.builder.AddImport("github.com/go-playground/validator/v10")

	return nil
}

// BuildInterface builds a handler interface
func (h *HandlerBuilder) BuildInterface(name string, methods []MethodSpec) error {
	if name == "" {
		return fmt.Errorf("interface name cannot be empty")
	}

	typeBuilder := NewTypeBuilder(h.builder)
	exprBuilder := NewExpressionBuilder(h.builder)

	// Create interface methods
	interfaceMethods := make([]*ast.Field, 0, len(methods))
	for _, method := range methods {
		// Create method parameters
		params := make([]*ast.Field, 0, len(method.Params))
		for _, param := range method.Params {
			params = append(params, typeBuilder.Field(param.Name, exprBuilder.Ident(param.Type), ""))
		}

		// Create method results
		results := make([]*ast.Field, 0, len(method.Returns))
		for _, result := range method.Returns {
			results = append(results, typeBuilder.Field(result.Name, exprBuilder.Ident(result.Type), ""))
		}

		// Create method field
		methodField := typeBuilder.Field(method.Name, exprBuilder.FuncType(params, results), "")
		interfaceMethods = append(interfaceMethods, methodField)
	}

	// Create interface declaration
	interfaceDecl := typeBuilder.InterfaceAlias(name, interfaceMethods)
	h.builder.AddDeclaration(interfaceDecl)

	return nil
}

// BuildRoutesFunction builds the routes function
func (h *HandlerBuilder) BuildRoutesFunction() error {
	funcBuilder := NewFunctionBuilder(h.builder)

	// Create routes function parameters
	params := []*ast.Field{
		funcBuilder.Param("h", "*Handler"),
		funcBuilder.Param("r", "*chi.Mux"),
	}

	// Create routes function body - include any statements that were added
	body := h.builder.GetStatements()
	if len(body) == 0 {
		// If no statements were added, create an empty body
		body = []ast.Stmt{}
	}

	// Create routes function declaration
	routesDecl := funcBuilder.Function("AddRoutes", params, nil, body)
	h.builder.AddDeclaration(routesDecl)

	// Add chi import
	h.builder.AddImport("github.com/go-chi/chi/v5")

	return nil
}

// High-level handler methods

// AddHandlerField adds a field to the handler struct
func (h *HandlerBuilder) AddHandlerField(name, typeName string) *HandlerBuilder {
	// This is a simplified version - in practice, you'd need to track the current struct
	// For now, we'll just create a new handler struct with the additional field
	typeBuilder := NewTypeBuilder(h.builder)
	exprBuilder := NewExpressionBuilder(h.builder)

	// Create handler struct with additional field
	fields := []*ast.Field{
		typeBuilder.Field("validator", exprBuilder.Star(exprBuilder.Select(exprBuilder.Ident("validator"), "Validate")), ""),
		typeBuilder.Field(name, exprBuilder.Ident(typeName), ""),
	}

	// Create handler struct declaration
	handlerDecl := typeBuilder.StructAlias("Handler", fields)
	h.builder.AddDeclaration(handlerDecl)

	return h
}

// AddInterfaceMethod adds a method to the current interface
func (h *HandlerBuilder) AddInterfaceMethod(name string, params, returns []FieldSpec) *HandlerBuilder {
	// This is a simplified version - in practice, you'd need to track the current interface
	// For now, we'll just create a new interface with the method
	method := MethodSpec{
		Name:    name,
		Params:  params,
		Returns: returns,
	}

	h.BuildInterface("HandlerInterface", []MethodSpec{method})
	return h
}

// AddRoute adds a route to the routes function
func (h *HandlerBuilder) AddRoute(method, path, handlerName string) *HandlerBuilder {
	// For now, we'll just add the route as a statement to be included in the routes function
	// In a more sophisticated implementation, we'd track and modify the existing routes function
	exprBuilder := NewExpressionBuilder(h.builder)
	stmtBuilder := NewStatementBuilder(h.builder)

	// Create route statement with proper handler function
	routeStmt := stmtBuilder.MethodCallStmt(
		exprBuilder.Ident("r"),
		strings.Title(strings.ToLower(method)),
		exprBuilder.String(path),
		exprBuilder.Call(
			exprBuilder.Ident("http.HandlerFunc"),
			exprBuilder.Ident(handlerName),
		),
	)

	// Add the route statement to the builder
	h.builder.AddStatement(routeStmt)

	// Add chi import
	h.builder.AddImport("github.com/go-chi/chi/v5")

	return h
}

// AddResponseWriter adds response writer methods
func (h *HandlerBuilder) AddResponseWriter(baseName string, codes []string) *HandlerBuilder {
	funcBuilder := NewFunctionBuilder(h.builder)
	exprBuilder := NewExpressionBuilder(h.builder)
	stmtBuilder := NewStatementBuilder(h.builder)

	// Create response writer methods for each status code
	for _, code := range codes {
		methodName := fmt.Sprintf("%s%s", baseName, code)

		// Create method parameters
		params := []*ast.Field{
			funcBuilder.Param("w", "http.ResponseWriter"),
			funcBuilder.Param("data", "interface{}"),
		}

		// Create method body
		body := []ast.Stmt{
			stmtBuilder.Assign(
				exprBuilder.MethodCall(exprBuilder.Ident("w"), "Header"),
				exprBuilder.Call(
					exprBuilder.Select(exprBuilder.Ident("http"), "StatusText"),
					exprBuilder.Ident(code),
				),
			),
			stmtBuilder.Assign(
				exprBuilder.MethodCall(exprBuilder.Ident("w"), "WriteHeader", exprBuilder.Ident(code)),
				exprBuilder.Nil(),
			),
			// Add JSON encoding here
			// JSON encoding will be added here
		}

		// Create method declaration
		methodDecl := funcBuilder.Function(methodName, params, nil, body)
		h.builder.AddDeclaration(methodDecl)
	}

	// Add http import
	h.builder.AddImport("net/http")

	return h
}

// BuildFromOpenAPI builds handlers from OpenAPI specification
func (h *HandlerBuilder) BuildFromOpenAPI(spec *openapi3.T) error {
	if spec == nil {
		return fmt.Errorf("OpenAPI specification cannot be nil")
	}

	// Build handler struct
	if err := h.BuildHandlerStruct(); err != nil {
		return fmt.Errorf("failed to build handler struct: %w", err)
	}

	// Build constructor
	if err := h.BuildConstructor(); err != nil {
		return fmt.Errorf("failed to build constructor: %w", err)
	}

	// Process paths and operations first to collect routes
	if spec.Paths != nil {
		paths := spec.Paths.Map()
		for path, pathItem := range paths {
			if pathItem == nil {
				continue
			}

			// Process each HTTP method
			for method, operation := range pathItem.Operations() {
				if operation == nil {
					continue
				}

				// Generate handler method name
				handlerName := h.generateHandlerName(operation.OperationID, method, path)

				// Add route
				h.AddRoute(method, path, handlerName)

				// Generate handler method
				if err := h.buildHandlerMethod(operation, handlerName); err != nil {
					return fmt.Errorf("failed to build handler method %s: %w", handlerName, err)
				}
			}
		}
	}

	// Build routes function after all routes have been added
	if err := h.BuildRoutesFunction(); err != nil {
		return fmt.Errorf("failed to build routes function: %w", err)
	}

	return nil
}

// Helper methods

func (h *HandlerBuilder) buildHandlerMethod(operation *openapi3.Operation, methodName string) error {
	funcBuilder := NewFunctionBuilder(h.builder)
	exprBuilder := NewExpressionBuilder(h.builder)
	stmtBuilder := NewStatementBuilder(h.builder)

	// Create method parameters - use standard http.HandlerFunc signature
	params := []*ast.Field{
		funcBuilder.Param("w", "http.ResponseWriter"),
		funcBuilder.Param("r", "*http.Request"),
	}

	// Create method body
	body := []ast.Stmt{
		// Add method implementation here
		// Implementation for methodName
		stmtBuilder.MethodCallStmt(
			exprBuilder.Ident("w"),
			"WriteHeader",
			exprBuilder.Int(200),
		),
	}

	// Create method declaration
	methodDecl := funcBuilder.Function(methodName, params, nil, body)
	h.builder.AddDeclaration(methodDecl)

	// Add http import
	h.builder.AddImport("net/http")

	return nil
}

func (h *HandlerBuilder) generateHandlerName(operationID, method, path string) string {
	if operationID != "" {
		return h.toPascalCase(operationID)
	}

	// Generate name from method and path
	pathParts := strings.Split(strings.Trim(path, "/"), "/")
	var nameParts []string

	// Add method
	nameParts = append(nameParts, strings.ToUpper(method))

	// Add path parts
	for _, part := range pathParts {
		if part != "" && !strings.HasPrefix(part, "{") {
			nameParts = append(nameParts, h.toPascalCase(part))
		}
	}

	return strings.Join(nameParts, "")
}

func (h *HandlerBuilder) toPascalCase(str string) string {
	if str == "" {
		return ""
	}

	// Simple PascalCase conversion
	words := strings.Split(str, "_")
	for i, word := range words {
		if word != "" {
			words[i] = strings.Title(word)
		}
	}

	return strings.Join(words, "")
}

// Fluent interface methods

// WithPackageName sets the package name
func (h *HandlerBuilder) WithPackageName(name string) *HandlerBuilder {
	h.config.PackageName = name
	return h
}

// WithUsePointers sets whether to use pointers
func (h *HandlerBuilder) WithUsePointers(use bool) *HandlerBuilder {
	h.config.UsePointers = use
	return h
}

// WithImportPrefix sets the import prefix
func (h *HandlerBuilder) WithImportPrefix(prefix string) *HandlerBuilder {
	h.config.ImportPrefix = prefix
	return h
}

// GetConfig returns the current configuration
func (h *HandlerBuilder) GetConfig() HandlerConfig {
	return h.config
}

// GetBuilder returns the underlying builder
func (h *HandlerBuilder) GetBuilder() *Builder {
	return h.builder
}
