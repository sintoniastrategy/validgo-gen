package astbuilder

import (
	"go/ast"
	"go/token"
)

// InterfaceBuilder provides a fluent interface for building Go interface declarations
type InterfaceBuilder struct {
	name    string
	methods []*InterfaceMethodBuilder
}

// InterfaceMethodBuilder provides a fluent interface for building interface methods
type InterfaceMethodBuilder struct {
	interfaceBuilder *InterfaceBuilder
	name             string
	params           []*FieldBuilder
	results          []*FieldBuilder
}

// NewInterfaceMethodBuilder creates a new InterfaceMethodBuilder
func NewInterfaceMethodBuilder() *InterfaceMethodBuilder {
	return &InterfaceMethodBuilder{
		params:  make([]*FieldBuilder, 0),
		results: make([]*FieldBuilder, 0),
	}
}

// CreateInterfaceBuilder creates a new InterfaceBuilder
func NewInterfaceBuilder() *InterfaceBuilder {
	return &InterfaceBuilder{
		methods: make([]*InterfaceMethodBuilder, 0),
	}
}

// WithName sets the interface name
// Returns the builder for method chaining
func (ib *InterfaceBuilder) WithName(name string) *InterfaceBuilder {
	ib.name = name
	return ib
}

// WithMethod adds an existing method builder to this interface
// Returns the interface builder for method chaining
func (ib *InterfaceBuilder) WithMethod(methodBuilder *InterfaceMethodBuilder) *InterfaceBuilder {
	if methodBuilder == nil {
		panic("method builder cannot be nil")
	}

	// Always update the interface builder reference to this interface
	methodBuilder.interfaceBuilder = ib

	ib.methods = append(ib.methods, methodBuilder)
	return ib
}

// Build creates the interface declaration AST
func (ib *InterfaceBuilder) Build() *ast.GenDecl {
	if ib.name == "" {
		panic("interface name is required")
	}

	// Convert InterfaceMethodBuilder instances to ast.Field instances
	fields := make([]*ast.Field, 0, len(ib.methods))
	for _, methodBuilder := range ib.methods {
		fields = append(fields, methodBuilder.toField())
	}

	interfaceType := &ast.InterfaceType{
		Methods: &ast.FieldList{
			List: fields,
		},
	}

	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent(ib.name),
				Type: interfaceType,
			},
		},
	}
}

// BuildAsDeclaration creates the interface declaration AST
// This is an alias for Build() for consistency with other builders
func (ib *InterfaceBuilder) BuildAsDeclaration() ast.Decl {
	return ib.Build()
}

// WithName sets the method name
// Returns the method builder for method chaining
func (imb *InterfaceMethodBuilder) WithName(name string) *InterfaceMethodBuilder {
	imb.name = name
	return imb
}

// AddArgField adds a parameter to the method using a FieldBuilder
// Returns the method builder for method chaining
func (imb *InterfaceMethodBuilder) AddArgField(fieldBuilder *FieldBuilder) *InterfaceMethodBuilder {
	if fieldBuilder == nil {
		panic("field builder cannot be nil")
	}
	imb.params = append(imb.params, fieldBuilder.Clone())
	return imb
}

// AddRetvalField adds a return value to the method using a FieldBuilder
// Returns the method builder for method chaining
func (imb *InterfaceMethodBuilder) AddRetvalField(fieldBuilder *FieldBuilder) *InterfaceMethodBuilder {
	if fieldBuilder == nil {
		panic("field builder cannot be nil")
	}
	imb.results = append(imb.results, fieldBuilder.Clone())
	return imb
}

// toField converts the InterfaceMethodBuilder to an ast.Field
func (imb *InterfaceMethodBuilder) toField() *ast.Field {
	// Convert FieldBuilder arrays to ast.Field arrays
	params := make([]*ast.Field, len(imb.params))
	for i, paramBuilder := range imb.params {
		params[i] = paramBuilder.Build()
	}

	results := make([]*ast.Field, len(imb.results))
	for i, resultBuilder := range imb.results {
		results[i] = resultBuilder.Build()
	}

	return &ast.Field{
		Names: []*ast.Ident{ast.NewIdent(imb.name)},
		Type: &ast.FuncType{
			Params:  &ast.FieldList{List: params},
			Results: &ast.FieldList{List: results},
		},
	}
}

// Helper methods for common parameter types

// AddContextArg adds a context.Context parameter
func (imb *InterfaceMethodBuilder) AddContextArg() *InterfaceMethodBuilder {
	return imb.AddArgField(ContextField("ctx"))
}

// AddErrorRetval adds an error return value
func (imb *InterfaceMethodBuilder) AddErrorRetval() *InterfaceMethodBuilder {
	return imb.AddRetvalField(ErrorField())
}

// Utility methods for the InterfaceBuilder

// MethodCount returns the number of methods in the interface
func (ib *InterfaceBuilder) MethodCount() int {
	return len(ib.methods)
}

// HasMethods returns true if the interface has any methods
func (ib *InterfaceBuilder) HasMethods() bool {
	return len(ib.methods) > 0
}

// ClearMethods removes all methods from the interface
func (ib *InterfaceBuilder) ClearMethods() *InterfaceBuilder {
	ib.methods = make([]*InterfaceMethodBuilder, 0)
	return ib
}

// GetMethodNames returns a slice of all method names in the interface
func (ib *InterfaceBuilder) GetMethodNames() []string {
	names := make([]string, 0, len(ib.methods))
	for _, method := range ib.methods {
		if method.name != "" {
			names = append(names, method.name)
		}
	}
	return names
}

// GetMethod returns the InterfaceMethodBuilder at the specified index
func (ib *InterfaceBuilder) GetMethod(index int) *InterfaceMethodBuilder {
	if index < 0 || index >= len(ib.methods) {
		return nil
	}
	return ib.methods[index]
}

// GetMethodByName returns the InterfaceMethodBuilder with the specified name
func (ib *InterfaceBuilder) GetMethodByName(name string) *InterfaceMethodBuilder {
	for _, method := range ib.methods {
		if method.name == name {
			return method
		}
	}
	return nil
}

// RemoveMethod removes the method at the specified index
func (ib *InterfaceBuilder) RemoveMethod(index int) *InterfaceBuilder {
	if index >= 0 && index < len(ib.methods) {
		ib.methods = append(ib.methods[:index], ib.methods[index+1:]...)
	}
	return ib
}

// RemoveMethodByName removes the method with the specified name
func (ib *InterfaceBuilder) RemoveMethodByName(name string) *InterfaceBuilder {
	for i, method := range ib.methods {
		if method.name == name {
			ib.methods = append(ib.methods[:i], ib.methods[i+1:]...)
			break
		}
	}
	return ib
}
