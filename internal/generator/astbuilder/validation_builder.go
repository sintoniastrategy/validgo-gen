package astbuilder

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// ValidationConfig holds configuration for validation building
type ValidationConfig struct {
	PackageName   string
	UsePointers   bool
	ImportPrefix  string
	ValidatorName string
	ErrorHandler  string
}

// ValidationRule represents a validation rule
type ValidationRule struct {
	Field     string
	Validator string
	Message   string
	Value     interface{}
}

// ValidationBuilder provides high-level methods for building Go validation code
type ValidationBuilder struct {
	builder *Builder
	config  ValidationConfig
}

// NewValidationBuilder creates a new validation builder
func NewValidationBuilder(builder *Builder, config ValidationConfig) *ValidationBuilder {
	return &ValidationBuilder{
		builder: builder,
		config:  config,
	}
}

// BuildObjectValidation builds validation for an object schema
func (v *ValidationBuilder) BuildObjectValidation(modelName string, schema *openapi3.SchemaRef) error {
	if schema == nil || schema.Value == nil {
		return fmt.Errorf("schema cannot be nil")
	}

	// Add validator import
	v.builder.AddImport("github.com/go-playground/validator/v10")

	// Build validation for each property
	for propName, propSchema := range schema.Value.Properties {
		if propSchema == nil || propSchema.Value == nil {
			continue
		}

		// Generate validation rules for this property
		rules := v.BuildFieldValidation(propName, propSchema)
		if len(rules) > 0 {
			// Create validation method for this field
			if err := v.buildFieldValidationMethod(modelName, propName, rules); err != nil {
				return fmt.Errorf("failed to build validation for field %s: %w", propName, err)
			}
		}
	}

	// Build main validation method
	if err := v.buildMainValidationMethod(modelName, schema); err != nil {
		return fmt.Errorf("failed to build main validation method: %w", err)
	}

	return nil
}

// BuildArrayValidation builds validation for an array schema
func (v *ValidationBuilder) BuildArrayValidation(modelName string, schema *openapi3.SchemaRef) error {
	if schema == nil || schema.Value == nil {
		return fmt.Errorf("schema cannot be nil")
	}

	// Add validator import
	v.builder.AddImport("github.com/go-playground/validator/v10")

	// Build array validation
	if err := v.buildArrayValidationMethod(modelName, schema); err != nil {
		return fmt.Errorf("failed to build array validation method: %w", err)
	}

	return nil
}

// BuildFieldValidation builds validation rules for a field
func (v *ValidationBuilder) BuildFieldValidation(fieldName string, schema *openapi3.SchemaRef) []string {
	if schema == nil || schema.Value == nil {
		return nil
	}

	var rules []string
	schemaValue := schema.Value

	// Required validation
	if schemaValue.Required != nil {
		for _, requiredField := range schemaValue.Required {
			if requiredField == fieldName {
				rules = append(rules, "required")
				break
			}
		}
	}

	// Type-specific validations
	if schemaValue.Type != nil && len(*schemaValue.Type) > 0 {
		switch (*schemaValue.Type)[0] {
		case "string":
			rules = append(rules, v.buildStringValidations(schemaValue)...)
		case "integer", "number":
			rules = append(rules, v.buildNumericValidations(schemaValue)...)
		case "array":
			rules = append(rules, v.buildArrayValidations(schemaValue)...)
		}
	}

	// Format validations
	if schemaValue.Format != "" {
		rules = append(rules, v.buildFormatValidations(schemaValue.Format)...)
	}

	// Enum validations
	if len(schemaValue.Enum) > 0 {
		rules = append(rules, v.buildEnumValidations(schemaValue.Enum)...)
	}

	return rules
}

// High-level validation methods

// AddRequiredFieldsValidation adds validation for required fields
func (v *ValidationBuilder) AddRequiredFieldsValidation(fields []string) *ValidationBuilder {
	stmtBuilder := NewStatementBuilder(v.builder)
	exprBuilder := NewExpressionBuilder(v.builder)

	// Create validation for each required field
	for _, field := range fields {
		validationStmt := stmtBuilder.If(
			exprBuilder.Equal(exprBuilder.Ident(field), exprBuilder.Nil()),
			[]ast.Stmt{
				stmtBuilder.Return(
					exprBuilder.Nil(),
					exprBuilder.Call(
						exprBuilder.Select(exprBuilder.Ident("errors"), "New"),
						exprBuilder.String(fmt.Sprintf("field %s is required", field)),
					),
				),
			},
		)
		v.builder.AddStatement(validationStmt)
	}

	// Add errors import
	v.builder.AddImport("errors")

	return v
}

// AddFieldValidation adds validation for a specific field
func (v *ValidationBuilder) AddFieldValidation(fieldName string, validators []string) *ValidationBuilder {
	stmtBuilder := NewStatementBuilder(v.builder)
	exprBuilder := NewExpressionBuilder(v.builder)

	// Create validation for each validator
	for _, validator := range validators {
		validationStmt := stmtBuilder.If(
			exprBuilder.NotEqual(exprBuilder.Ident(fieldName), exprBuilder.Nil()),
			[]ast.Stmt{
				stmtBuilder.AssignDefine(
					exprBuilder.Ident("err"),
					exprBuilder.Call(
						exprBuilder.Select(exprBuilder.Ident("validator"), "Var"),
						exprBuilder.Ident(fieldName),
						exprBuilder.String(validator),
					),
				),
				stmtBuilder.If(
					exprBuilder.NotEqual(exprBuilder.Ident("err"), exprBuilder.Nil()),
					[]ast.Stmt{
						stmtBuilder.Return(exprBuilder.Nil(), exprBuilder.Ident("err")),
					},
				),
			},
		)
		v.builder.AddStatement(validationStmt)
	}

	// Add validator import
	v.builder.AddImport("github.com/go-playground/validator/v10")

	return v
}

// AddJSONUnmarshal adds JSON unmarshaling validation
func (v *ValidationBuilder) AddJSONUnmarshal() *ValidationBuilder {
	stmtBuilder := NewStatementBuilder(v.builder)
	exprBuilder := NewExpressionBuilder(v.builder)

	// Create JSON unmarshal validation
	unmarshalStmt := stmtBuilder.AssignDefine(
		exprBuilder.Ident("err"),
		exprBuilder.Call(
			exprBuilder.Select(exprBuilder.Ident("json"), "Unmarshal"),
			exprBuilder.Ident("data"),
			exprBuilder.Ident("target"),
		),
	)
	v.builder.AddStatement(unmarshalStmt)

	// Add error handling
	errorStmt := stmtBuilder.If(
		exprBuilder.NotEqual(exprBuilder.Ident("err"), exprBuilder.Nil()),
		[]ast.Stmt{
			stmtBuilder.Return(exprBuilder.Nil(), exprBuilder.Ident("err")),
		},
	)
	v.builder.AddStatement(errorStmt)

	// Add json import
	v.builder.AddImport("encoding/json")

	return v
}

// AddErrorHandling adds error handling validation
func (v *ValidationBuilder) AddErrorHandling() *ValidationBuilder {
	stmtBuilder := NewStatementBuilder(v.builder)
	exprBuilder := NewExpressionBuilder(v.builder)

	// Create error handling pattern
	errorStmt := stmtBuilder.If(
		exprBuilder.NotEqual(exprBuilder.Ident("err"), exprBuilder.Nil()),
		[]ast.Stmt{
			stmtBuilder.Return(exprBuilder.Nil(), exprBuilder.Ident("err")),
		},
	)
	v.builder.AddStatement(errorStmt)

	return v
}

// Helper methods

func (v *ValidationBuilder) buildFieldValidationMethod(modelName, fieldName string, rules []string) error {
	funcBuilder := NewFunctionBuilder(v.builder)
	stmtBuilder := NewStatementBuilder(v.builder)
	exprBuilder := NewExpressionBuilder(v.builder)

	// Create method name
	methodName := fmt.Sprintf("Validate%s%s", v.toPascalCase(fieldName), v.toPascalCase(modelName))

	// Create method parameters
	params := []*ast.Field{
		funcBuilder.Param("v", "*validator.Validate"),
		funcBuilder.Param("value", "interface{}"),
	}

	// Create method results
	results := []*ast.Field{
		funcBuilder.ResultAnonymous("error"),
	}

	// Build validation body
	var body []ast.Stmt
	for _, rule := range rules {
		validationStmt := stmtBuilder.AssignDefine(
			exprBuilder.Ident("err"),
			exprBuilder.Call(
				exprBuilder.Select(exprBuilder.Ident("v"), "Var"),
				exprBuilder.Ident("value"),
				exprBuilder.String(rule),
			),
		)
		body = append(body, validationStmt)

		errorStmt := stmtBuilder.If(
			exprBuilder.NotEqual(exprBuilder.Ident("err"), exprBuilder.Nil()),
			[]ast.Stmt{
				stmtBuilder.Return(exprBuilder.Nil(), exprBuilder.Ident("err")),
			},
		)
		body = append(body, errorStmt)
	}

	// Add success return
	body = append(body, stmtBuilder.Return(exprBuilder.Nil(), exprBuilder.Nil()))

	// Create method declaration
	methodDecl := funcBuilder.Function(methodName, params, results, body)
	v.builder.AddDeclaration(methodDecl)

	return nil
}

func (v *ValidationBuilder) buildMainValidationMethod(modelName string, schema *openapi3.SchemaRef) error {
	funcBuilder := NewFunctionBuilder(v.builder)
	stmtBuilder := NewStatementBuilder(v.builder)
	exprBuilder := NewExpressionBuilder(v.builder)

	// Create method name
	methodName := fmt.Sprintf("Validate%s", v.toPascalCase(modelName))

	// Create method parameters
	params := []*ast.Field{
		funcBuilder.Param("v", "*validator.Validate"),
		funcBuilder.Param("model", fmt.Sprintf("*%s", v.toPascalCase(modelName))),
	}

	// Create method results
	results := []*ast.Field{
		funcBuilder.ResultAnonymous("error"),
	}

	// Build validation body
	body := []ast.Stmt{
		stmtBuilder.AssignDefine(
			exprBuilder.Ident("err"),
			exprBuilder.Call(
				exprBuilder.Select(exprBuilder.Ident("v"), "Struct"),
				exprBuilder.Ident("model"),
			),
		),
		stmtBuilder.If(
			exprBuilder.NotEqual(exprBuilder.Ident("err"), exprBuilder.Nil()),
			[]ast.Stmt{
				stmtBuilder.Return(exprBuilder.Nil(), exprBuilder.Ident("err")),
			},
		),
		stmtBuilder.Return(exprBuilder.Nil(), exprBuilder.Nil()),
	}

	// Create method declaration
	methodDecl := funcBuilder.Function(methodName, params, results, body)
	v.builder.AddDeclaration(methodDecl)

	return nil
}

func (v *ValidationBuilder) buildArrayValidationMethod(modelName string, schema *openapi3.SchemaRef) error {
	funcBuilder := NewFunctionBuilder(v.builder)
	stmtBuilder := NewStatementBuilder(v.builder)
	exprBuilder := NewExpressionBuilder(v.builder)

	// Create method name
	methodName := fmt.Sprintf("Validate%sArray", v.toPascalCase(modelName))

	// Create method parameters
	params := []*ast.Field{
		funcBuilder.Param("v", "*validator.Validate"),
		funcBuilder.Param("array", "[]interface{}"),
	}

	// Create method results
	results := []*ast.Field{
		funcBuilder.ResultAnonymous("error"),
	}

	// Build validation body
	body := []ast.Stmt{
		stmtBuilder.AssignDefine(
			exprBuilder.Ident("err"),
			exprBuilder.Call(
				exprBuilder.Select(exprBuilder.Ident("v"), "Var"),
				exprBuilder.Ident("array"),
				exprBuilder.String("dive"),
			),
		),
		stmtBuilder.If(
			exprBuilder.NotEqual(exprBuilder.Ident("err"), exprBuilder.Nil()),
			[]ast.Stmt{
				stmtBuilder.Return(exprBuilder.Nil(), exprBuilder.Ident("err")),
			},
		),
		stmtBuilder.Return(exprBuilder.Nil(), exprBuilder.Nil()),
	}

	// Create method declaration
	methodDecl := funcBuilder.Function(methodName, params, results, body)
	v.builder.AddDeclaration(methodDecl)

	return nil
}

func (v *ValidationBuilder) buildStringValidations(schema *openapi3.Schema) []string {
	var rules []string

	if schema.MinLength > 0 {
		rules = append(rules, fmt.Sprintf("min=%d", schema.MinLength))
	}
	if schema.MaxLength != nil {
		rules = append(rules, fmt.Sprintf("max=%d", *schema.MaxLength))
	}
	if schema.Pattern != "" {
		rules = append(rules, fmt.Sprintf("regexp=%s", schema.Pattern))
	}

	return rules
}

func (v *ValidationBuilder) buildNumericValidations(schema *openapi3.Schema) []string {
	var rules []string

	if schema.Min != nil {
		rules = append(rules, fmt.Sprintf("min=%v", *schema.Min))
	}
	if schema.Max != nil {
		rules = append(rules, fmt.Sprintf("max=%v", *schema.Max))
	}
	if schema.MultipleOf != nil {
		rules = append(rules, fmt.Sprintf("multipleof=%v", *schema.MultipleOf))
	}

	return rules
}

func (v *ValidationBuilder) buildArrayValidations(schema *openapi3.Schema) []string {
	var rules []string

	if schema.MinItems > 0 {
		rules = append(rules, fmt.Sprintf("min=%d", schema.MinItems))
	}
	if schema.MaxItems != nil {
		rules = append(rules, fmt.Sprintf("max=%d", *schema.MaxItems))
	}
	if schema.UniqueItems {
		rules = append(rules, "unique")
	}

	return rules
}

func (v *ValidationBuilder) buildFormatValidations(format string) []string {
	switch format {
	case "email":
		return []string{"email"}
	case "uri":
		return []string{"uri"}
	case "date":
		return []string{"datetime=2006-01-02"}
	case "date-time":
		return []string{"datetime=2006-01-02T15:04:05Z07:00"}
	case "uuid":
		return []string{"uuid"}
	default:
		return nil
	}
}

func (v *ValidationBuilder) buildEnumValidations(enum []interface{}) []string {
	if len(enum) == 0 {
		return nil
	}

	var values []string
	for _, val := range enum {
		values = append(values, fmt.Sprintf("%v", val))
	}

	return []string{fmt.Sprintf("oneof=%s", strings.Join(values, " "))}
}

func (v *ValidationBuilder) toPascalCase(str string) string {
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
func (v *ValidationBuilder) WithPackageName(name string) *ValidationBuilder {
	v.config.PackageName = name
	return v
}

// WithUsePointers sets whether to use pointers
func (v *ValidationBuilder) WithUsePointers(use bool) *ValidationBuilder {
	v.config.UsePointers = use
	return v
}

// WithImportPrefix sets the import prefix
func (v *ValidationBuilder) WithImportPrefix(prefix string) *ValidationBuilder {
	v.config.ImportPrefix = prefix
	return v
}

// WithValidatorName sets the validator variable name
func (v *ValidationBuilder) WithValidatorName(name string) *ValidationBuilder {
	v.config.ValidatorName = name
	return v
}

// WithErrorHandler sets the error handler name
func (v *ValidationBuilder) WithErrorHandler(name string) *ValidationBuilder {
	v.config.ErrorHandler = name
	return v
}

// GetConfig returns the current configuration
func (v *ValidationBuilder) GetConfig() ValidationConfig {
	return v.config
}

// GetBuilder returns the underlying builder
func (v *ValidationBuilder) GetBuilder() *Builder {
	return v.builder
}
