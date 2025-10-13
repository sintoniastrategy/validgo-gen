package astbuilder

/*
import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// SchemaConfig holds configuration for schema building
type SchemaConfig struct {
	PackageName  string
	UsePointers  bool
	ImportPrefix string
}

// SchemaBuilder provides high-level methods for building Go schemas from OpenAPI schemas
type SchemaBuilder struct {
	builder *Builder
	config  SchemaConfig
}

// NewSchemaBuilder creates a new schema builder
func NewSchemaBuilder(builder *Builder, config SchemaConfig) *SchemaBuilder {
	return &SchemaBuilder{
		builder: builder,
		config:  config,
	}
}

// BuildStruct builds a Go struct from a SchemaStruct
func (s *SchemaBuilder) BuildStruct(model generator.SchemaStruct) error {
	if model.Name == "" {
		return fmt.Errorf("struct name cannot be empty")
	}

	typeBuilder := NewTypeBuilder(s.builder)

	// Create struct fields
	fields := make([]*ast.Field, 0, len(model.Fields))
	for _, field := range model.Fields {
		astField := s.BuildField(field)
		fields = append(fields, astField)
	}

	// Create struct declaration
	structDecl := typeBuilder.StructAlias(model.Name, fields)
	s.builder.AddDeclaration(structDecl)

	return nil
}

// BuildTypeAlias builds a type alias
func (s *SchemaBuilder) BuildTypeAlias(name, typeName string) error {
	if name == "" || typeName == "" {
		return fmt.Errorf("name and typeName cannot be empty")
	}

	typeBuilder := NewTypeBuilder(s.builder)

	// Create type alias
	typeDecl := typeBuilder.TypeAlias(name, typeName)
	s.builder.AddDeclaration(typeDecl)

	return nil
}

// BuildSliceAlias builds a slice type alias
func (s *SchemaBuilder) BuildSliceAlias(name, elementType string) error {
	if name == "" || elementType == "" {
		return fmt.Errorf("name and elementType cannot be empty")
	}

	typeBuilder := NewTypeBuilder(s.builder)
	exprBuilder := NewExpressionBuilder(s.builder)

	// Create slice type
	sliceType := exprBuilder.SliceType(exprBuilder.Ident(elementType))

	// Create slice alias declaration
	sliceDecl := typeBuilder.TypeAliasWithType(name, sliceType)
	s.builder.AddDeclaration(sliceDecl)

	return nil
}

// BuildField builds an AST field from a SchemaField
func (s *SchemaBuilder) BuildField(field generator.SchemaField) *ast.Field {
	typeBuilder := NewTypeBuilder(s.builder)
	exprBuilder := NewExpressionBuilder(s.builder)

	// Determine Go type
	goType := s.getGoType(field.Type)
	var typeExpr ast.Expr

	// Add time import if using time.Time
	if goType == "time.Time" {
		s.builder.AddImport("time")
	}

	// Handle pointer types for optional fields
	if !field.Required && s.config.UsePointers {
		typeExpr = exprBuilder.Star(exprBuilder.Ident(goType))
	} else {
		typeExpr = exprBuilder.Ident(goType)
	}

	// Build tags
	tags := s.buildTags(field)

	// Create field with valid Go identifier
	goFieldName := generator.FormatGoLikeIdentifier(field.Name)
	astField := typeBuilder.Field(goFieldName, typeExpr, tags)

	return astField
}

// High-level schema methods

// AddStruct adds a struct with the given name and fields
func (s *SchemaBuilder) AddStruct(name string, fields []generator.SchemaField) *SchemaBuilder {
	model := generator.SchemaStruct{
		Name:   name,
		Fields: fields,
	}

	if err := s.BuildStruct(model); err != nil {
		// In a real implementation, you might want to handle this error
		// For now, we'll just continue
	}

	return s
}

// AddField adds a field to the current struct
func (s *SchemaBuilder) AddField(name, typeName string, tags map[string]string) *SchemaBuilder {
	field := generator.SchemaField{
		Name:        name,
		Type:        typeName,
		TagJSON:     []string{name},
		TagValidate: []string{},
		Required:    true,
	}

	// Convert tags map to field tags
	if jsonTag, ok := tags["json"]; ok {
		field.TagJSON = []string{jsonTag}
	}
	if validateTag, ok := tags["validate"]; ok {
		field.TagValidate = []string{validateTag}
	}
	if required, ok := tags["required"]; ok && required == "false" {
		field.Required = false
	}

	// This is a simplified version - in practice, you'd need to track the current struct
	// For now, we'll just create a single-field struct
	model := generator.SchemaStruct{
		Name:   s.generateStructName(name),
		Fields: []generator.SchemaField{field},
	}

	s.BuildStruct(model)

	return s
}

// AddTypeAlias adds a type alias
func (s *SchemaBuilder) AddTypeAlias(name, typeName string) *SchemaBuilder {
	s.BuildTypeAlias(name, typeName)
	return s
}

// AddSliceType adds a slice type alias
func (s *SchemaBuilder) AddSliceType(name, elementType string) *SchemaBuilder {
	s.BuildSliceAlias(name, elementType)
	return s
}

// BuildFromOpenAPISchema builds a Go struct from an OpenAPI schema
func (s *SchemaBuilder) BuildFromOpenAPISchema(name string, schema *openapi3.SchemaRef) error {
	if schema == nil || schema.Value == nil {
		return fmt.Errorf("schema cannot be nil")
	}

	// Handle different schema types
	switch {
	case schema.Value.Type.Permits(openapi3.TypeObject):
		return s.buildObjectSchema(name, schema)
	case schema.Value.Type.Permits(openapi3.TypeArray):
		return s.buildArraySchema(name, schema)
	case schema.Value.Type.Permits(openapi3.TypeString):
		return s.buildStringSchema(name, schema)
	case schema.Value.Type.Permits(openapi3.TypeInteger):
		return s.buildIntegerSchema(name, schema)
	case schema.Value.Type.Permits(openapi3.TypeNumber):
		return s.buildNumberSchema(name, schema)
	case schema.Value.Type.Permits(openapi3.TypeBoolean):
		return s.buildBooleanSchema(name, schema)
	default:
		return s.buildStringSchema(name, schema) // Default to string
	}
}

// Helper methods

func (s *SchemaBuilder) buildObjectSchema(name string, schema *openapi3.SchemaRef) error {
	fields := make([]generator.SchemaField, 0)

	// Process properties
	for propName, propSchema := range schema.Value.Properties {
		field := s.createFieldFromOpenAPI(propName, propSchema, schema.Value.Required)
		fields = append(fields, field)
	}

	// Create struct
	model := generator.SchemaStruct{
		Name:   name,
		Fields: fields,
	}

	return s.BuildStruct(model)
}

func (s *SchemaBuilder) buildArraySchema(name string, schema *openapi3.SchemaRef) error {
	// Get element type
	elementType := "string" // Default
	if schema.Value.Items != nil && schema.Value.Items.Value != nil {
		elementType = s.getOpenAPITypeAsGo(schema.Value.Items.Value.Type, schema.Value.Items.Value.Format)
	}

	// Create slice alias
	return s.BuildSliceAlias(name, elementType)
}

func (s *SchemaBuilder) buildStringSchema(name string, schema *openapi3.SchemaRef) error {
	// Create type alias for string
	return s.BuildTypeAlias(name, "string")
}

func (s *SchemaBuilder) buildIntegerSchema(name string, schema *openapi3.SchemaRef) error {
	// Create type alias for integer
	return s.BuildTypeAlias(name, "int")
}

func (s *SchemaBuilder) buildNumberSchema(name string, schema *openapi3.SchemaRef) error {
	// Create type alias for number
	return s.BuildTypeAlias(name, "float64")
}

func (s *SchemaBuilder) buildBooleanSchema(name string, schema *openapi3.SchemaRef) error {
	// Create type alias for boolean
	return s.BuildTypeAlias(name, "bool")
}

func (s *SchemaBuilder) createFieldFromOpenAPI(name string, schema *openapi3.SchemaRef, required []string) generator.SchemaField {
	field := generator.SchemaField{
		Name:     s.toPascalCase(name),
		Type:     s.getOpenAPITypeAsGo(schema.Value.Type, schema.Value.Format),
		TagJSON:  []string{name},
		Required: s.isFieldRequired(name, required),
	}

	// Add validation tags based on schema constraints
	validationTags := s.getValidationTags(schema.Value)
	field.TagValidate = validationTags

	return field
}

func (s *SchemaBuilder) getOpenAPITypeAsGo(openapiType *openapi3.Types, format string) string {
	if openapiType == nil {
		return "string"
	}

	switch {
	case openapiType.Permits(openapi3.TypeString):
		// Check for date-time format
		if format == "date-time" {
			return "time.Time"
		}
		return "string"
	case openapiType.Permits(openapi3.TypeInteger):
		return "int"
	case openapiType.Permits(openapi3.TypeNumber):
		return "float64"
	case openapiType.Permits(openapi3.TypeBoolean):
		return "bool"
	case openapiType.Permits(openapi3.TypeArray):
		return "[]string" // Default array type
	default:
		return "string"
	}
}

func (s *SchemaBuilder) isFieldRequired(fieldName string, required []string) bool {
	for _, req := range required {
		if req == fieldName {
			return true
		}
	}
	return false
}

func (s *SchemaBuilder) getValidationTags(schema *openapi3.Schema) []string {
	var tags []string

	// Required validation
	if schema.Type.Permits(openapi3.TypeString) {
		if schema.MinLength > 0 {
			tags = append(tags, fmt.Sprintf("min=%d", schema.MinLength))
		}
		if schema.MaxLength != nil {
			tags = append(tags, fmt.Sprintf("max=%d", *schema.MaxLength))
		}
		if schema.Pattern != "" {
			tags = append(tags, fmt.Sprintf("pattern=%s", schema.Pattern))
		}
	}

	if schema.Type.Permits(openapi3.TypeInteger) || schema.Type.Permits(openapi3.TypeNumber) {
		if schema.Min != nil {
			tags = append(tags, fmt.Sprintf("min=%f", *schema.Min))
		}
		if schema.Max != nil {
			tags = append(tags, fmt.Sprintf("max=%f", *schema.Max))
		}
	}

	// Email validation
	if schema.Format == "email" {
		tags = append(tags, "email")
	}

	// URL validation
	if schema.Format == "uri" {
		tags = append(tags, "url")
	}

	return tags
}

func (s *SchemaBuilder) getGoType(typeName string) string {
	// Handle common Go types
	switch typeName {
	case "string":
		return "string"
	case "int":
		return "int"
	case "int64":
		return "int64"
	case "float64":
		return "float64"
	case "bool":
		return "bool"
	case "time.Time":
		return "time.Time"
	default:
		// Check if it's a slice type
		if strings.HasPrefix(typeName, "[]") {
			return typeName
		}
		// Check if it's a map type
		if strings.HasPrefix(typeName, "map[") {
			return typeName
		}
		// Default to string
		return "string"
	}
}

func (s *SchemaBuilder) buildTags(field generator.SchemaField) string {
	tags := make(map[string]string)

	// JSON tag
	if len(field.TagJSON) > 0 {
		jsonTag := field.TagJSON[0]
		if !field.Required {
			jsonTag += ",omitempty"
		}
		tags["json"] = jsonTag
	}

	// Validation tag
	if len(field.TagValidate) > 0 {
		tags["validate"] = strings.Join(field.TagValidate, ",")
	}

	// Build tag string
	var tagParts []string
	for key, value := range tags {
		tagParts = append(tagParts, fmt.Sprintf("%s:\"%s\"", key, value))
	}

	return strings.Join(tagParts, " ")
}

func (s *SchemaBuilder) toPascalCase(str string) string {
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

func (s *SchemaBuilder) generateStructName(fieldName string) string {
	return s.toPascalCase(fieldName) + "Struct"
}

// Fluent interface methods

// WithPackageName sets the package name
func (s *SchemaBuilder) WithPackageName(name string) *SchemaBuilder {
	s.config.PackageName = name
	return s
}

// WithUsePointers sets whether to use pointers
func (s *SchemaBuilder) WithUsePointers(use bool) *SchemaBuilder {
	s.config.UsePointers = use
	return s
}

// WithImportPrefix sets the import prefix
func (s *SchemaBuilder) WithImportPrefix(prefix string) *SchemaBuilder {
	s.config.ImportPrefix = prefix
	return s
}

// GetConfig returns the current configuration
func (s *SchemaBuilder) GetConfig() SchemaConfig {
	return s.config
}

// GetBuilder returns the underlying builder
func (s *SchemaBuilder) GetBuilder() *Builder {
	return s.builder
}
*/
