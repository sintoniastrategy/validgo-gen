package astbuilder

import (
	"go/ast"
	"go/token"
)

// FieldBuilder provides a fluent interface for building ast.Field
type FieldBuilder struct {
	name        string
	typeBuilder *SimpleTypeBuilder
	tag         string
}

// NewFieldBuilder creates a new FieldBuilder
func NewFieldBuilder() *FieldBuilder {
	return &FieldBuilder{
		typeBuilder: NewSimpleTypeBuilder(),
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

// WithTag sets the field tag
// Returns the builder for method chaining
func (fb *FieldBuilder) WithTag(tag string) *FieldBuilder {
	fb.tag = tag
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

	// Set tag if provided
	if fb.tag != "" {
		field.Tag = &ast.BasicLit{
			Kind:  token.STRING,
			Value: fb.tag,
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

// HasTag returns true if the field has a tag
func (fb *FieldBuilder) HasTag() bool {
	return fb.tag != ""
}

// GetTag returns the field tag
func (fb *FieldBuilder) GetTag() string {
	return fb.tag
}

// Clone creates a copy of the FieldBuilder
func (fb *FieldBuilder) Clone() *FieldBuilder {
	clone := &FieldBuilder{
		name: fb.name,
		tag:  fb.tag,
	}
	if fb.typeBuilder != nil {
		clone.typeBuilder = fb.typeBuilder.Clone()
	}
	return clone
}
