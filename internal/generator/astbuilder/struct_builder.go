package astbuilder

import (
	"go/ast"
	"go/token"
)

// StructBuilder provides a fluent interface for building Go structs
type StructBuilder struct {
	name   string
	fields []*FieldBuilder
}

// NewStructBuilder creates a new StructBuilder
func NewStructBuilder() *StructBuilder {
	return &StructBuilder{
		fields: make([]*FieldBuilder, 0),
	}
}

// WithName sets the struct name
// Returns the builder for method chaining
func (sb *StructBuilder) WithName(name string) *StructBuilder {
	sb.name = name
	return sb
}

// AddField adds a field to the struct using a FieldBuilder
// Returns the builder for method chaining
func (sb *StructBuilder) AddField(fieldBuilder *FieldBuilder) *StructBuilder {
	if fieldBuilder == nil {
		panic("field builder cannot be nil")
	}
	sb.fields = append(sb.fields, fieldBuilder.Clone())
	return sb
}

// AddFields adds multiple fields to the struct
// Returns the builder for method chaining
func (sb *StructBuilder) AddFields(fieldBuilders ...*FieldBuilder) *StructBuilder {
	for _, fieldBuilder := range fieldBuilders {
		if fieldBuilder == nil {
			panic("field builder cannot be nil")
		}
		sb.fields = append(sb.fields, fieldBuilder.Clone())
	}
	return sb
}

// Build creates the ast.TypeSpec for the struct
func (sb *StructBuilder) Build() *ast.TypeSpec {
	if sb.name == "" {
		panic("struct must have a name")
	}

	// Convert FieldBuilder array to ast.Field array
	astFields := make([]*ast.Field, len(sb.fields))
	for i, fieldBuilder := range sb.fields {
		astFields[i] = fieldBuilder.Build()
	}

	return &ast.TypeSpec{
		Name: ast.NewIdent(sb.name),
		Type: &ast.StructType{
			Fields: &ast.FieldList{
				List: astFields,
			},
		},
	}
}

// BuildAsDeclaration creates the ast.GenDecl for the struct
func (sb *StructBuilder) BuildAsDeclaration() *ast.GenDecl {
	return &ast.GenDecl{
		Tok:   token.TYPE,
		Specs: []ast.Spec{sb.Build()},
	}
}

// Utility methods

// HasName returns true if the struct has a name
func (sb *StructBuilder) HasName() bool {
	return sb.name != ""
}

// GetName returns the struct name
func (sb *StructBuilder) GetName() string {
	return sb.name
}

// FieldCount returns the number of fields in the struct
func (sb *StructBuilder) FieldCount() int {
	return len(sb.fields)
}

// HasFields returns true if the struct has any fields
func (sb *StructBuilder) HasFields() bool {
	return len(sb.fields) > 0
}

// GetFields returns a copy of the fields slice
func (sb *StructBuilder) GetFields() []*FieldBuilder {
	fields := make([]*FieldBuilder, len(sb.fields))
	for i, field := range sb.fields {
		fields[i] = field.Clone()
	}
	return fields
}

// GetField returns the field at the specified index
// Returns nil if index is out of bounds
func (sb *StructBuilder) GetField(index int) *FieldBuilder {
	if index < 0 || index >= len(sb.fields) {
		return nil
	}
	return sb.fields[index].Clone()
}

// GetFieldByName returns the first field with the specified name
// Returns nil if no field with that name is found
func (sb *StructBuilder) GetFieldByName(name string) *FieldBuilder {
	for _, field := range sb.fields {
		if field.GetName() == name {
			return field.Clone()
		}
	}
	return nil
}

// RemoveField removes the field at the specified index
// Returns the builder for method chaining
func (sb *StructBuilder) RemoveField(index int) *StructBuilder {
	if index < 0 || index >= len(sb.fields) {
		return sb
	}
	sb.fields = append(sb.fields[:index], sb.fields[index+1:]...)
	return sb
}

// RemoveFieldByName removes the first field with the specified name
// Returns the builder for method chaining
func (sb *StructBuilder) RemoveFieldByName(name string) *StructBuilder {
	for i, field := range sb.fields {
		if field.GetName() == name {
			sb.fields = append(sb.fields[:i], sb.fields[i+1:]...)
			return sb
		}
	}
	return sb
}

// Clear removes all fields from the struct
// Returns the builder for method chaining
func (sb *StructBuilder) Clear() *StructBuilder {
	sb.fields = make([]*FieldBuilder, 0)
	return sb
}

// Clone creates a copy of the StructBuilder
func (sb *StructBuilder) Clone() *StructBuilder {
	clone := &StructBuilder{
		name:   sb.name,
		fields: make([]*FieldBuilder, len(sb.fields)),
	}
	for i, field := range sb.fields {
		clone.fields[i] = field.Clone()
	}
	return clone
}

// Helper methods for common field types

// AddStringField adds a string field to the struct
func (sb *StructBuilder) AddStringField(name string) *StructBuilder {
	return sb.AddField(StringField(name))
}

// AddIntField adds an int field to the struct
func (sb *StructBuilder) AddIntField(name string) *StructBuilder {
	return sb.AddField(IntField(name))
}

// AddBoolField adds a bool field to the struct
func (sb *StructBuilder) AddBoolField(name string) *StructBuilder {
	return sb.AddField(BoolField(name))
}

// AddContextField adds a context.Context field to the struct
func (sb *StructBuilder) AddContextField(name string) *StructBuilder {
	return sb.AddField(ContextField(name))
}

// AddIdentField adds a field with a single identifier type
func (sb *StructBuilder) AddIdentField(name, typeName string) *StructBuilder {
	return sb.AddField(IdentField(name, typeName))
}

// AddSelectorField adds a field with a selector expression type (e.g., "pkg.Type")
func (sb *StructBuilder) AddSelectorField(name, pkg, typeName string) *StructBuilder {
	return sb.AddField(SelectorField(name, pkg, typeName))
}
