package astbuilder

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// ParameterConfig holds configuration for parameter parsing
type ParameterConfig struct {
	BaseName      string
	PackageName   string
	UsePointers   bool
	ParameterType string // "Query", "Header", "Cookie", "Path"
	StructName    string
	ImportPrefix  string
}

// ParameterParser provides high-level methods for parsing OpenAPI parameters
type ParameterParser struct {
	builder *Builder
	config  ParameterConfig
}

// NewParameterParser creates a new parameter parser
func NewParameterParser(builder *Builder, config ParameterConfig) *ParameterParser {
	return &ParameterParser{
		builder: builder,
		config:  config,
	}
}

// ParseParameters parses all parameters of a specific type
func (p *ParameterParser) ParseParameters(params openapi3.Parameters) error {
	if len(params) == 0 {
		return nil
	}

	// Filter parameters by type
	filteredParams := p.filterParametersByType(params)
	if len(filteredParams) == 0 {
		return nil
	}

	// Generate struct name if not provided
	if p.config.StructName == "" {
		p.config.StructName = p.generateStructName()
	}

	// Build the parameter parsing method
	return p.buildParameterParsingMethod(filteredParams)
}

// ParseQueryParams parses query parameters
func (p *ParameterParser) ParseQueryParams(params openapi3.Parameters) error {
	config := p.config
	config.ParameterType = "Query"
	config.StructName = p.config.BaseName + "QueryParams"

	parser := &ParameterParser{
		builder: p.builder,
		config:  config,
	}

	return parser.ParseParameters(params)
}

// ParseHeaders parses header parameters
func (p *ParameterParser) ParseHeaders(params openapi3.Parameters) error {
	config := p.config
	config.ParameterType = "Header"
	config.StructName = p.config.BaseName + "Headers"

	parser := &ParameterParser{
		builder: p.builder,
		config:  config,
	}

	return parser.ParseParameters(params)
}

// ParseCookies parses cookie parameters
func (p *ParameterParser) ParseCookies(params openapi3.Parameters) error {
	config := p.config
	config.ParameterType = "Cookie"
	config.StructName = p.config.BaseName + "Cookies"

	parser := &ParameterParser{
		builder: p.builder,
		config:  config,
	}

	return parser.ParseParameters(params)
}

// ParsePathParams parses path parameters
func (p *ParameterParser) ParsePathParams(params openapi3.Parameters) error {
	config := p.config
	config.ParameterType = "Path"
	config.StructName = p.config.BaseName + "PathParams"

	// Add chi import for path parameters
	p.builder.AddImport("github.com/go-chi/chi/v5")

	parser := &ParameterParser{
		builder: p.builder,
		config:  config,
	}

	return parser.ParseParameters(params)
}

// DeclareStruct declares the parameter struct
func (p *ParameterParser) DeclareStruct() *ParameterParser {
	typeBuilder := NewTypeBuilder(p.builder)

	// Create struct fields from parameters
	fields := p.createStructFields()

	// Create struct declaration
	structDecl := typeBuilder.StructAlias(p.config.StructName, fields)
	p.builder.AddDeclaration(structDecl)

	return p
}

// ExtractParameter extracts a single parameter
func (p *ParameterParser) ExtractParameter(param *openapi3.Parameter) *ParameterParser {
	stmtBuilder := NewStatementBuilder(p.builder)

	// Get parameter value based on type
	paramValue := p.getParameterValue(param)

	// Declare variable for the parameter
	declStmt := stmtBuilder.DeclareVar(
		p.getParameterVarName(param),
		p.getGoType(param),
		paramValue,
	)

	p.builder.AddStatement(declStmt)

	return p
}

// ValidateRequired validates required parameters
func (p *ParameterParser) ValidateRequired(param *openapi3.Parameter) *ParameterParser {
	if !param.Required {
		return p
	}

	stmtBuilder := NewStatementBuilder(p.builder)
	exprBuilder := NewExpressionBuilder(p.builder)

	// Add validation for required parameters
	validationStmt := stmtBuilder.If(
		exprBuilder.Equal(
			exprBuilder.Ident(p.getParameterVarName(param)),
			exprBuilder.String(""),
		),
		[]ast.Stmt{
			stmtBuilder.Return(
				exprBuilder.Nil(),
				exprBuilder.Call(
					exprBuilder.Select(exprBuilder.Ident("errors"), "New"),
					exprBuilder.String(fmt.Sprintf("%s is required", param.Name)),
				),
			),
		},
	)

	p.builder.AddStatement(validationStmt)
	p.builder.AddImport("github.com/go-faster/errors")

	return p
}

// AssignToField assigns parameter to struct field
func (p *ParameterParser) AssignToField(param *openapi3.Parameter) *ParameterParser {
	stmtBuilder := NewStatementBuilder(p.builder)
	exprBuilder := NewExpressionBuilder(p.builder)

	// Assign parameter to struct field
	assignStmt := stmtBuilder.Assign(
		exprBuilder.Select(exprBuilder.Ident(p.config.StructName), p.getFieldName(param)),
		exprBuilder.Ident(p.getParameterVarName(param)),
	)

	p.builder.AddStatement(assignStmt)

	return p
}

// ValidateStruct validates the entire struct
func (p *ParameterParser) ValidateStruct() *ParameterParser {
	stmtBuilder := NewStatementBuilder(p.builder)
	exprBuilder := NewExpressionBuilder(p.builder)

	// Add struct validation
	validationStmt := stmtBuilder.AssignDefine(
		exprBuilder.Ident("err"),
		exprBuilder.MethodCall(
			exprBuilder.Select(exprBuilder.Ident("h"), "validator"),
			"Struct",
			exprBuilder.Ident(p.config.StructName),
		),
	)

	p.builder.AddStatement(validationStmt)

	// Add error handling
	errorStmt := stmtBuilder.If(
		exprBuilder.NotEqual(exprBuilder.Ident("err"), exprBuilder.Nil()),
		[]ast.Stmt{
			stmtBuilder.Return(exprBuilder.Nil(), exprBuilder.Ident("err")),
		},
	)

	p.builder.AddStatement(errorStmt)
	p.builder.AddImport("github.com/go-playground/validator/v10")

	return p
}

// ReturnResult returns the parsed parameters
func (p *ParameterParser) ReturnResult() *ParameterParser {
	stmtBuilder := NewStatementBuilder(p.builder)
	exprBuilder := NewExpressionBuilder(p.builder)

	// Return the struct and nil error
	returnStmt := stmtBuilder.Return(
		exprBuilder.Ident(p.config.StructName),
		exprBuilder.Nil(),
	)

	p.builder.AddStatement(returnStmt)

	return p
}

// Helper methods

func (p *ParameterParser) filterParametersByType(params openapi3.Parameters) openapi3.Parameters {
	var filtered []*openapi3.ParameterRef
	for _, param := range params {
		if param.Value != nil && strings.ToLower(param.Value.In) == strings.ToLower(p.config.ParameterType) {
			filtered = append(filtered, param)
		}
	}
	return filtered
}

func (p *ParameterParser) generateStructName() string {
	return p.config.BaseName + p.config.ParameterType + "Params"
}

func (p *ParameterParser) buildParameterParsingMethod(params openapi3.Parameters) error {
	// Create method name
	methodName := fmt.Sprintf("Parse%sParams", p.config.ParameterType)

	// Create method parameters
	funcBuilder := NewFunctionBuilder(p.builder)
	paramsList := []*ast.Field{
		funcBuilder.Param("r", "*http.Request"),
	}

	// Create method results
	resultsList := []*ast.Field{
		funcBuilder.ResultAnonymous(p.config.StructName),
		funcBuilder.ResultAnonymous("error"),
	}

	// Add imports
	p.builder.AddImport("net/http")

	// Build method body
	body := p.buildMethodBody(params)

	// Create method declaration
	methodDecl := funcBuilder.Method(
		funcBuilder.Receiver("h", "*Handler"),
		methodName,
		paramsList,
		resultsList,
		body,
	)

	p.builder.AddDeclaration(methodDecl)

	return nil
}

func (p *ParameterParser) buildMethodBody(params openapi3.Parameters) []ast.Stmt {
	// Create struct instance
	exprBuilder := NewExpressionBuilder(p.builder)
	stmtBuilder := NewStatementBuilder(p.builder)

	// Declare struct variable
	structDecl := stmtBuilder.DeclareVar(
		p.config.StructName,
		p.config.StructName,
		nil,
	)

	body := []ast.Stmt{structDecl}

	// Process each parameter
	for _, param := range params {
		if param.Value == nil {
			continue
		}

		// Extract parameter
		paramValue := p.getParameterValue(param.Value)
		paramVar := p.getParameterVarName(param.Value)
		goType := p.getGoType(param.Value)

		// Declare parameter variable
		paramDecl := stmtBuilder.DeclareVar(paramVar, goType, paramValue)
		body = append(body, paramDecl)

		// Validate required
		if param.Value.Required {
			validation := stmtBuilder.If(
				exprBuilder.Equal(exprBuilder.Ident(paramVar), exprBuilder.String("")),
				[]ast.Stmt{
					stmtBuilder.Return(
						exprBuilder.Ident(p.config.StructName),
						exprBuilder.Call(
							exprBuilder.Select(exprBuilder.Ident("errors"), "New"),
							exprBuilder.String(fmt.Sprintf("%s is required", param.Value.Name)),
						),
					),
				},
			)
			body = append(body, validation)
		}

		// Assign to struct field
		assign := stmtBuilder.Assign(
			exprBuilder.Select(exprBuilder.Ident(p.config.StructName), p.getFieldName(param.Value)),
			exprBuilder.Ident(paramVar),
		)
		body = append(body, assign)
	}

	// Add struct validation
	validation := stmtBuilder.AssignDefine(
		exprBuilder.Ident("err"),
		exprBuilder.MethodCall(
			exprBuilder.Select(exprBuilder.Ident("h"), "validator"),
			"Struct",
			exprBuilder.Ident(p.config.StructName),
		),
	)
	body = append(body, validation)

	// Add error handling
	errorCheck := stmtBuilder.If(
		exprBuilder.NotEqual(exprBuilder.Ident("err"), exprBuilder.Nil()),
		[]ast.Stmt{
			stmtBuilder.Return(exprBuilder.Ident(p.config.StructName), exprBuilder.Ident("err")),
		},
	)
	body = append(body, errorCheck)

	// Return result
	returnStmt := stmtBuilder.Return(exprBuilder.Ident(p.config.StructName), exprBuilder.Nil())
	body = append(body, returnStmt)

	return body
}

func (p *ParameterParser) createStructFields() []*ast.Field {
	// This would be called from DeclareStruct to create the struct fields
	// For now, return empty - this will be implemented when we have the parameters
	return []*ast.Field{}
}

func (p *ParameterParser) getParameterValue(param *openapi3.Parameter) ast.Expr {
	exprBuilder := NewExpressionBuilder(p.builder)

	switch p.config.ParameterType {
	case "Query":
		return exprBuilder.MethodCall(
			exprBuilder.MethodCall(exprBuilder.Ident("r"), "URL"),
			"Query",
		)
	case "Header":
		return exprBuilder.MethodCall(
			exprBuilder.MethodCall(exprBuilder.Ident("r"), "Header"),
			"Get",
		)
	case "Cookie":
		return exprBuilder.MethodCall(
			exprBuilder.MethodCall(exprBuilder.Ident("r"), "Header"),
			"Get",
		)
	case "Path":
		return exprBuilder.Call(
			exprBuilder.Select(exprBuilder.Ident("chi"), "URLParam"),
			exprBuilder.Ident("r"),
			exprBuilder.String(param.Name),
		)
	default:
		return exprBuilder.String("")
	}
}

func (p *ParameterParser) getParameterVarName(param *openapi3.Parameter) string {
	return strings.ToLower(param.Name) + "param"
}

func (p *ParameterParser) getFieldName(param *openapi3.Parameter) string {
	// Convert parameter name to Go field name (PascalCase)
	return strings.Title(param.Name)
}

func (p *ParameterParser) getGoType(param *openapi3.Parameter) string {
	if param.Schema == nil || param.Schema.Value == nil {
		return "string"
	}

	schema := param.Schema.Value
	switch {
	case schema.Type.Permits(openapi3.TypeString):
		return "string"
	case schema.Type.Permits(openapi3.TypeInteger):
		return "int"
	case schema.Type.Permits(openapi3.TypeNumber):
		return "float64"
	case schema.Type.Permits(openapi3.TypeBoolean):
		return "bool"
	case schema.Type.Permits(openapi3.TypeArray):
		return "[]string"
	default:
		return "string"
	}
}

// Fluent interface methods for chaining

// WithStructName sets the struct name
func (p *ParameterParser) WithStructName(name string) *ParameterParser {
	p.config.StructName = name
	return p
}

// WithPackageName sets the package name
func (p *ParameterParser) WithPackageName(name string) *ParameterParser {
	p.config.PackageName = name
	return p
}

// WithUsePointers sets whether to use pointers
func (p *ParameterParser) WithUsePointers(use bool) *ParameterParser {
	p.config.UsePointers = use
	return p
}

// GetConfig returns the current configuration
func (p *ParameterParser) GetConfig() ParameterConfig {
	return p.config
}

// GetBuilder returns the underlying builder
func (p *ParameterParser) GetBuilder() *Builder {
	return p.builder
}
