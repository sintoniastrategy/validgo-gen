package astbuilder

import (
	"go/ast"
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
