package astbuilder

import (
	"go/ast"
	"go/token"
)

// SimpleTypeBuilder provides a fluent interface for building simple type expressions
// It can build expressions like "a", "a.B", "package.Type", etc.
type SimpleTypeBuilder struct {
	elements  []string
	asPointer bool
}

// NewSimpleTypeBuilder creates a new SimpleTypeBuilder
func NewSimpleTypeBuilder() *SimpleTypeBuilder {
	return &SimpleTypeBuilder{
		elements: make([]string, 0),
	}
}

// AddElement adds an element to the type expression
// Returns the builder for method chaining
func (stb *SimpleTypeBuilder) AddElement(element string) *SimpleTypeBuilder {
	if element == "" {
		return stb
	}
	stb.elements = append(stb.elements, element)
	return stb
}

// AddElements adds multiple elements to the type expression
// Returns the builder for method chaining
func (stb *SimpleTypeBuilder) AddElements(elements ...string) *SimpleTypeBuilder {
	for _, element := range elements {
		if element != "" {
			stb.elements = append(stb.elements, element)
		}
	}
	return stb
}

// AsPointer sets whether the type should be built as a pointer (*Type)
// Returns the builder for method chaining
func (stb *SimpleTypeBuilder) AsPointer(isPointer bool) *SimpleTypeBuilder {
	stb.asPointer = isPointer
	return stb
}

// Build creates the ast.Expr for the simple type
func (stb *SimpleTypeBuilder) Build() ast.Expr {
	if len(stb.elements) == 0 {
		panic("simple type must have at least one element")
	}

	var expr ast.Expr
	if len(stb.elements) == 1 {
		// Single element - just an identifier
		expr = ast.NewIdent(stb.elements[0])
	} else {
		// Multiple elements - build a selector expression
		expr = ast.NewIdent(stb.elements[0])
		for i := 1; i < len(stb.elements); i++ {
			expr = &ast.SelectorExpr{
				X:   expr,
				Sel: ast.NewIdent(stb.elements[i]),
			}
		}
	}

	// If asPointer is true, wrap the expression in a StarExpr
	if stb.asPointer {
		return &ast.StarExpr{X: expr}
	}
	return expr
}

// BuildAsIdent creates an ast.Ident if there's only one element, panics otherwise
func (stb *SimpleTypeBuilder) BuildAsIdent() *ast.Ident {
	if len(stb.elements) != 1 {
		panic("BuildAsIdent requires exactly one element")
	}
	return ast.NewIdent(stb.elements[0])
}

// BuildAsSelector creates an ast.SelectorExpr if there are multiple elements, panics otherwise
func (stb *SimpleTypeBuilder) BuildAsSelector() *ast.SelectorExpr {
	if len(stb.elements) < 2 {
		panic("BuildAsSelector requires at least two elements")
	}

	var expr ast.Expr = ast.NewIdent(stb.elements[0])
	for i := 1; i < len(stb.elements); i++ {
		expr = &ast.SelectorExpr{
			X:   expr,
			Sel: ast.NewIdent(stb.elements[i]),
		}
	}
	return expr.(*ast.SelectorExpr)
}

// ElementCount returns the number of elements in the type
func (stb *SimpleTypeBuilder) ElementCount() int {
	return len(stb.elements)
}

// HasElements returns true if the type has any elements
func (stb *SimpleTypeBuilder) HasElements() bool {
	return len(stb.elements) > 0
}

// GetElements returns a copy of the elements slice
func (stb *SimpleTypeBuilder) GetElements() []string {
	elements := make([]string, len(stb.elements))
	copy(elements, stb.elements)
	return elements
}

// Clear removes all elements from the type
func (stb *SimpleTypeBuilder) Clear() *SimpleTypeBuilder {
	stb.elements = make([]string, 0)
	return stb
}

// Clone creates a copy of the SimpleTypeBuilder
func (stb *SimpleTypeBuilder) Clone() *SimpleTypeBuilder {
	clone := &SimpleTypeBuilder{
		elements:  make([]string, len(stb.elements)),
		asPointer: stb.asPointer,
	}
	copy(clone.elements, stb.elements)
	return clone
}

// Helper methods for common types

// String creates a simple type builder for "string"
func String() *SimpleTypeBuilder {
	return NewSimpleTypeBuilder().AddElement("string")
}

// Int creates a simple type builder for "int"
func Int() *SimpleTypeBuilder {
	return NewSimpleTypeBuilder().AddElement("int")
}

// Bool creates a simple type builder for "bool"
func Bool() *SimpleTypeBuilder {
	return NewSimpleTypeBuilder().AddElement("bool")
}

// Error creates a simple type builder for "error"
func Error() *SimpleTypeBuilder {
	return NewSimpleTypeBuilder().AddElement("error")
}

// Context creates a simple type builder for "context.Context"
func Context() *SimpleTypeBuilder {
	return NewSimpleTypeBuilder().AddElements("context", "Context")
}

// Ident creates a simple type builder for a single identifier
func Ident(name string) *SimpleTypeBuilder {
	return NewSimpleTypeBuilder().AddElement(name)
}

// Selector creates a simple type builder for a selector expression like "package.Type"
func Selector(packageName, typeName string) *SimpleTypeBuilder {
	return NewSimpleTypeBuilder().AddElements(packageName, typeName)
}

// Pointer creates a pointer to the type built by this builder
func (stb *SimpleTypeBuilder) Pointer() ast.Expr {
	return &ast.StarExpr{X: stb.Build()}
}

// Slice creates a slice of the type built by this builder
func (stb *SimpleTypeBuilder) Slice() ast.Expr {
	return &ast.ArrayType{Elt: stb.Build()}
}

// Array creates an array of the type built by this builder
func (stb *SimpleTypeBuilder) Array(length int) ast.Expr {
	return &ast.ArrayType{
		Len: &ast.BasicLit{Kind: 0, Value: string(rune(length + '0'))}, // This is simplified
		Elt: stb.Build(),
	}
}

// TypeExpressionBuilder is an interface that can build ast.Expr types
type TypeExpressionBuilder interface {
	Build() ast.Expr
}

// ArrayTypeBuilder provides a fluent interface for building slice types
// It can build expressions like "[]string", "[][]int", etc.
type ArrayTypeBuilder struct {
	element TypeExpressionBuilder
}

// NewArrayTypeBuilder creates a new ArrayTypeBuilder
func NewArrayTypeBuilder() *ArrayTypeBuilder {
	return &ArrayTypeBuilder{}
}

// WithElement sets the element type using a TypeExpressionBuilder (SimpleTypeBuilder or ArrayTypeBuilder)
// Returns the builder for method chaining
func (atb *ArrayTypeBuilder) WithElement(element TypeExpressionBuilder) *ArrayTypeBuilder {
	if element == nil {
		panic("element cannot be nil")
	}
	atb.element = element
	return atb
}

// Build creates the ast.Expr for the array type
func (atb *ArrayTypeBuilder) Build() ast.Expr {
	if atb.element == nil {
		panic("array type must have an element type")
	}
	return &ast.ArrayType{
		Elt: atb.element.Build(),
	}
}

// Utility methods for ArrayTypeBuilder

// HasElement returns true if the array has an element type
func (atb *ArrayTypeBuilder) HasElement() bool {
	return atb.element != nil
}

// GetElement returns the element TypeExpressionBuilder
func (atb *ArrayTypeBuilder) GetElement() TypeExpressionBuilder {
	return atb.element
}

// Clone creates a copy of the ArrayTypeBuilder
func (atb *ArrayTypeBuilder) Clone() *ArrayTypeBuilder {
	clone := &ArrayTypeBuilder{}
	if atb.element != nil {
		// Clone the element if it's a SimpleTypeBuilder
		if stb, ok := atb.element.(*SimpleTypeBuilder); ok {
			clone.element = stb.Clone()
		} else if atb, ok := atb.element.(*ArrayTypeBuilder); ok {
			clone.element = atb.Clone()
		} else {
			// For other types, we can't clone, so we'll panic
			panic("cannot clone unknown TypeExpressionBuilder type")
		}
	}
	return clone
}

// Helper functions for creating arrays

// StringSlice creates an ArrayTypeBuilder for []string
func StringSlice() *ArrayTypeBuilder {
	return NewArrayTypeBuilder().WithElement(String())
}

// IntSlice creates an ArrayTypeBuilder for []int
func IntSlice() *ArrayTypeBuilder {
	return NewArrayTypeBuilder().WithElement(Int())
}

// BoolSlice creates an ArrayTypeBuilder for []bool
func BoolSlice() *ArrayTypeBuilder {
	return NewArrayTypeBuilder().WithElement(Bool())
}

// ErrorSlice creates an ArrayTypeBuilder for []error
func ErrorSlice() *ArrayTypeBuilder {
	return NewArrayTypeBuilder().WithElement(Error())
}

// ContextSlice creates an ArrayTypeBuilder for []context.Context
func ContextSlice() *ArrayTypeBuilder {
	return NewArrayTypeBuilder().WithElement(Context())
}

// IdentSlice creates an ArrayTypeBuilder for []Identifier
func IdentSlice(identifier string) *ArrayTypeBuilder {
	return NewArrayTypeBuilder().WithElement(Ident(identifier))
}

// SelectorSlice creates an ArrayTypeBuilder for []package.Type
func SelectorSlice(packageName, typeName string) *ArrayTypeBuilder {
	return NewArrayTypeBuilder().WithElement(Selector(packageName, typeName))
}

// SliceOf creates an ArrayTypeBuilder for []TypeExpressionBuilder
func SliceOf(element TypeExpressionBuilder) *ArrayTypeBuilder {
	return NewArrayTypeBuilder().WithElement(element)
}

// TypeAliasBuilder provides a fluent interface for building type aliases
// It can build declarations like "type NameOfAlias []string"
type TypeAliasBuilder struct {
	name        string
	typeBuilder TypeExpressionBuilder
}

// NewTypeAliasBuilder creates a new TypeAliasBuilder
func NewTypeAliasBuilder() *TypeAliasBuilder {
	return &TypeAliasBuilder{}
}

// WithName sets the alias name
// Returns the builder for method chaining
func (tab *TypeAliasBuilder) WithName(name string) *TypeAliasBuilder {
	if name == "" {
		panic("alias name cannot be empty")
	}
	tab.name = name
	return tab
}

// WithType sets the underlying type using a TypeExpressionBuilder
// Returns the builder for method chaining
func (tab *TypeAliasBuilder) WithType(typeBuilder TypeExpressionBuilder) *TypeAliasBuilder {
	if typeBuilder == nil {
		panic("type builder cannot be nil")
	}
	tab.typeBuilder = typeBuilder
	return tab
}

// Build creates the ast.TypeSpec for the type alias
func (tab *TypeAliasBuilder) Build() *ast.TypeSpec {
	if tab.name == "" {
		panic("type alias must have a name")
	}
	if tab.typeBuilder == nil {
		panic("type alias must have a type")
	}

	return &ast.TypeSpec{
		Name: ast.NewIdent(tab.name),
		Type: tab.typeBuilder.Build(),
	}
}

// BuildAsDeclaration creates the ast.GenDecl for the type alias
func (tab *TypeAliasBuilder) BuildAsDeclaration() *ast.GenDecl {
	return &ast.GenDecl{
		Tok:   token.TYPE,
		Specs: []ast.Spec{tab.Build()},
	}
}

// Utility methods for TypeAliasBuilder

// HasName returns true if the alias has a name
func (tab *TypeAliasBuilder) HasName() bool {
	return tab.name != ""
}

// GetName returns the alias name
func (tab *TypeAliasBuilder) GetName() string {
	return tab.name
}

// HasType returns true if the alias has a type
func (tab *TypeAliasBuilder) HasType() bool {
	return tab.typeBuilder != nil
}

// GetType returns the underlying type builder
func (tab *TypeAliasBuilder) GetType() TypeExpressionBuilder {
	return tab.typeBuilder
}

// Clone creates a copy of the TypeAliasBuilder
func (tab *TypeAliasBuilder) Clone() *TypeAliasBuilder {
	clone := &TypeAliasBuilder{
		name: tab.name,
	}

	if tab.typeBuilder != nil {
		// Clone the type builder based on its type
		if stb, ok := tab.typeBuilder.(*SimpleTypeBuilder); ok {
			clone.typeBuilder = stb.Clone()
		} else if atb, ok := tab.typeBuilder.(*ArrayTypeBuilder); ok {
			clone.typeBuilder = atb.Clone()
		} else {
			panic("cannot clone unknown TypeExpressionBuilder type")
		}
	}

	return clone
}

// Helper functions for creating type aliases

// StringSliceAlias creates a TypeAliasBuilder for "type AliasName []string"
func StringSliceAlias(name string) *TypeAliasBuilder {
	return NewTypeAliasBuilder().
		WithName(name).
		WithType(StringSlice())
}

// IntSliceAlias creates a TypeAliasBuilder for "type AliasName []int"
func IntSliceAlias(name string) *TypeAliasBuilder {
	return NewTypeAliasBuilder().
		WithName(name).
		WithType(IntSlice())
}

// BoolSliceAlias creates a TypeAliasBuilder for "type AliasName []bool"
func BoolSliceAlias(name string) *TypeAliasBuilder {
	return NewTypeAliasBuilder().
		WithName(name).
		WithType(BoolSlice())
}

// ErrorSliceAlias creates a TypeAliasBuilder for "type AliasName []error"
func ErrorSliceAlias(name string) *TypeAliasBuilder {
	return NewTypeAliasBuilder().
		WithName(name).
		WithType(ErrorSlice())
}

// ContextSliceAlias creates a TypeAliasBuilder for "type AliasName []context.Context"
func ContextSliceAlias(name string) *TypeAliasBuilder {
	return NewTypeAliasBuilder().
		WithName(name).
		WithType(ContextSlice())
}

// IdentAlias creates a TypeAliasBuilder for "type AliasName IdentType"
func IdentAlias(name, typeName string) *TypeAliasBuilder {
	return NewTypeAliasBuilder().
		WithName(name).
		WithType(Ident(typeName))
}

// SelectorAlias creates a TypeAliasBuilder for "type AliasName package.Type"
func SelectorAlias(name, packageName, typeName string) *TypeAliasBuilder {
	return NewTypeAliasBuilder().
		WithName(name).
		WithType(Selector(packageName, typeName))
}

// PointerAlias creates a TypeAliasBuilder for "type AliasName *Type"
func PointerAlias(name, typeName string) *TypeAliasBuilder {
	return NewTypeAliasBuilder().
		WithName(name).
		WithType(Ident(typeName).AsPointer(true))
}

// SliceAlias creates a TypeAliasBuilder for "type AliasName []Type"
func SliceAlias(name, typeName string) *TypeAliasBuilder {
	return NewTypeAliasBuilder().
		WithName(name).
		WithType(SliceOf(Ident(typeName)))
}

// ArrayAlias creates a TypeAliasBuilder for "type AliasName []Type" (same as SliceAlias)
func ArrayAlias(name, typeName string) *TypeAliasBuilder {
	return SliceAlias(name, typeName)
}

// CustomAlias creates a TypeAliasBuilder for "type AliasName CustomType"
func CustomAlias(name string, typeBuilder TypeExpressionBuilder) *TypeAliasBuilder {
	return NewTypeAliasBuilder().
		WithName(name).
		WithType(typeBuilder)
}
