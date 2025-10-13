package astbuilder

import (
	"go/ast"
	"go/token"
	"strings"
)

// FieldBuilder provides a fluent interface for building ast.Field
type FieldBuilder struct {
	name         string
	typeBuilder  *SimpleTypeBuilder
	jsonTags     []string
	validateTags []string
}

// NewFieldBuilder creates a new FieldBuilder
func NewFieldBuilder() *FieldBuilder {
	return &FieldBuilder{
		typeBuilder:  NewSimpleTypeBuilder(),
		jsonTags:     make([]string, 0),
		validateTags: make([]string, 0),
	}
}

// WithName sets the field name
// Returns the builder for method chaining
func (fb *FieldBuilder) WithName(name string) *FieldBuilder {
	fb.name = name
	return fb
}

// WithType sets the field type using a SimpleTypeBuilder
// Returns the builder for method chaining
func (fb *FieldBuilder) WithType(typeBuilder *SimpleTypeBuilder) *FieldBuilder {
	if typeBuilder == nil {
		panic("type builder cannot be nil")
	}
	fb.typeBuilder = typeBuilder.Clone()
	return fb
}

// AddJSONTag adds a JSON tag to the field
// Returns the builder for method chaining
func (fb *FieldBuilder) AddJSONTag(tag string) *FieldBuilder {
	if tag != "" {
		fb.jsonTags = append(fb.jsonTags, tag)
	}
	return fb
}

// AddJSONTags adds multiple JSON tags to the field
// Returns the builder for method chaining
func (fb *FieldBuilder) AddJSONTags(tags ...string) *FieldBuilder {
	for _, tag := range tags {
		if tag != "" {
			fb.jsonTags = append(fb.jsonTags, tag)
		}
	}
	return fb
}

// SetJSONTags replaces all JSON tags with the provided tags
// Returns the builder for method chaining
func (fb *FieldBuilder) SetJSONTags(tags ...string) *FieldBuilder {
	fb.jsonTags = make([]string, 0)
	for _, tag := range tags {
		if tag != "" {
			fb.jsonTags = append(fb.jsonTags, tag)
		}
	}
	return fb
}

// AddValidateTag adds a validate tag to the field
// Returns the builder for method chaining
func (fb *FieldBuilder) AddValidateTag(tag string) *FieldBuilder {
	if tag != "" {
		fb.validateTags = append(fb.validateTags, tag)
	}
	return fb
}

// AddValidateTags adds multiple validate tags to the field
// Returns the builder for method chaining
func (fb *FieldBuilder) AddValidateTags(tags ...string) *FieldBuilder {
	for _, tag := range tags {
		if tag != "" {
			fb.validateTags = append(fb.validateTags, tag)
		}
	}
	return fb
}

// SetValidateTags replaces all validate tags with the provided tags
// Returns the builder for method chaining
func (fb *FieldBuilder) SetValidateTags(tags ...string) *FieldBuilder {
	fb.validateTags = make([]string, 0)
	for _, tag := range tags {
		if tag != "" {
			fb.validateTags = append(fb.validateTags, tag)
		}
	}
	return fb
}

// Build creates the ast.Field
func (fb *FieldBuilder) Build() *ast.Field {
	if fb.typeBuilder == nil {
		panic("field must have a type builder")
	}

	field := &ast.Field{
		Type: fb.typeBuilder.Build(),
	}

	// Set name if provided
	if fb.name != "" {
		field.Names = []*ast.Ident{ast.NewIdent(fb.name)}
	}

	// Generate tag from JSON and validate tags
	tagParts := make([]string, 0)

	// Add JSON tags
	if len(fb.jsonTags) > 0 {
		jsonTag := strings.Join(fb.jsonTags, ",")
		tagParts = append(tagParts, "json:\""+jsonTag+"\"")
	}

	// Add validate tags
	if len(fb.validateTags) > 0 {
		validateTag := strings.Join(fb.validateTags, ",")
		tagParts = append(tagParts, "validate:\""+validateTag+"\"")
	}

	// Set tag if any tags were provided
	if len(tagParts) > 0 {
		fullTag := strings.Join(tagParts, " ")
		field.Tag = &ast.BasicLit{
			Kind:  token.STRING,
			Value: "`" + fullTag + "`",
		}
	}

	return field
}

// Helper methods for common field types

// StringField creates a field builder for a string field
func StringField(name string) *FieldBuilder {
	return NewFieldBuilder().
		WithName(name).
		WithType(String())
}

// IntField creates a field builder for an int field
func IntField(name string) *FieldBuilder {
	return NewFieldBuilder().
		WithName(name).
		WithType(Int())
}

// BoolField creates a field builder for a bool field
func BoolField(name string) *FieldBuilder {
	return NewFieldBuilder().
		WithName(name).
		WithType(Bool())
}

// ErrorField creates a field builder for an error field (unnamed)
func ErrorField() *FieldBuilder {
	return NewFieldBuilder().WithType(Error())
}

// ContextField creates a field builder for a context.Context field
func ContextField(name string) *FieldBuilder {
	return NewFieldBuilder().
		WithName(name).
		WithType(Context())
}

// IdentField creates a field builder for a field with a single identifier type
func IdentField(name, typeName string) *FieldBuilder {
	return NewFieldBuilder().
		WithName(name).
		WithType(Ident(typeName))
}

// SelectorField creates a field builder for a field with a selector type
func SelectorField(name, packageName, typeName string) *FieldBuilder {
	return NewFieldBuilder().
		WithName(name).
		WithType(Selector(packageName, typeName))
}

// Utility methods

// HasName returns true if the field has a name
func (fb *FieldBuilder) HasName() bool {
	return fb.name != ""
}

// GetName returns the field name
func (fb *FieldBuilder) GetName() string {
	return fb.name
}

// HasTags returns true if the field has any tags (JSON or validate)
func (fb *FieldBuilder) HasTags() bool {
	return len(fb.jsonTags) > 0 || len(fb.validateTags) > 0
}

// HasJSONTags returns true if the field has JSON tags
func (fb *FieldBuilder) HasJSONTags() bool {
	return len(fb.jsonTags) > 0
}

// HasValidateTags returns true if the field has validate tags
func (fb *FieldBuilder) HasValidateTags() bool {
	return len(fb.validateTags) > 0
}

// GetJSONTags returns a copy of the JSON tags
func (fb *FieldBuilder) GetJSONTags() []string {
	tags := make([]string, len(fb.jsonTags))
	copy(tags, fb.jsonTags)
	return tags
}

// GetValidateTags returns a copy of the validate tags
func (fb *FieldBuilder) GetValidateTags() []string {
	tags := make([]string, len(fb.validateTags))
	copy(tags, fb.validateTags)
	return tags
}

// ClearJSONTags removes all JSON tags
func (fb *FieldBuilder) ClearJSONTags() *FieldBuilder {
	fb.jsonTags = make([]string, 0)
	return fb
}

// ClearValidateTags removes all validate tags
func (fb *FieldBuilder) ClearValidateTags() *FieldBuilder {
	fb.validateTags = make([]string, 0)
	return fb
}

// ClearAllTags removes all tags (JSON and validate)
func (fb *FieldBuilder) ClearAllTags() *FieldBuilder {
	fb.jsonTags = make([]string, 0)
	fb.validateTags = make([]string, 0)
	return fb
}

// Clone creates a copy of the FieldBuilder
func (fb *FieldBuilder) Clone() *FieldBuilder {
	clone := &FieldBuilder{
		name:         fb.name,
		jsonTags:     make([]string, len(fb.jsonTags)),
		validateTags: make([]string, len(fb.validateTags)),
	}
	copy(clone.jsonTags, fb.jsonTags)
	copy(clone.validateTags, fb.validateTags)
	if fb.typeBuilder != nil {
		clone.typeBuilder = fb.typeBuilder.Clone()
	}
	return clone
}
