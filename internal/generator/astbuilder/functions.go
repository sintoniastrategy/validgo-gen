package astbuilder

import (
	"go/ast"
)

// FunctionBuilder provides methods for building AST function declarations
type FunctionBuilder struct {
	builder *Builder
}

// NewFunctionBuilder creates a new function builder
func NewFunctionBuilder(builder *Builder) *FunctionBuilder {
	return &FunctionBuilder{builder: builder}
}

// Func creates a function declaration
func (f *FunctionBuilder) Func(name string, receiver *ast.Field, params, results []*ast.Field, body []ast.Stmt) *ast.FuncDecl {
	funcDecl := &ast.FuncDecl{
		Name: f.builder.ident(name),
		Type: &ast.FuncType{
			Params:  &ast.FieldList{List: params},
			Results: &ast.FieldList{List: results},
		},
		Body: &ast.BlockStmt{List: body},
	}

	if receiver != nil {
		funcDecl.Recv = &ast.FieldList{List: []*ast.Field{receiver}}
	}

	return funcDecl
}

// Method creates a method declaration
func (f *FunctionBuilder) Method(receiver *ast.Field, name string, params, results []*ast.Field, body []ast.Stmt) *ast.FuncDecl {
	return f.Func(name, receiver, params, results, body)
}

// Function creates a function declaration without receiver
func (f *FunctionBuilder) Function(name string, params, results []*ast.Field, body []ast.Stmt) *ast.FuncDecl {
	return f.Func(name, nil, params, results, body)
}

// Param creates a function parameter
func (f *FunctionBuilder) Param(name, typeName string) *ast.Field {
	return &ast.Field{
		Names: []*ast.Ident{f.builder.ident(name)},
		Type:  f.builder.ident(typeName),
	}
}

// ParamWithType creates a function parameter with a type expression
func (f *FunctionBuilder) ParamWithType(name string, typeExpr ast.Expr) *ast.Field {
	return &ast.Field{
		Names: []*ast.Ident{f.builder.ident(name)},
		Type:  typeExpr,
	}
}

// ParamMultiple creates a function parameter with multiple names
func (f *FunctionBuilder) ParamMultiple(names []string, typeName string) *ast.Field {
	field := &ast.Field{
		Type: f.builder.ident(typeName),
	}

	if len(names) > 0 {
		field.Names = make([]*ast.Ident, len(names))
		for i, name := range names {
			field.Names[i] = f.builder.ident(name)
		}
	}

	return field
}

// ParamMultipleWithType creates a function parameter with multiple names and a type expression
func (f *FunctionBuilder) ParamMultipleWithType(names []string, typeExpr ast.Expr) *ast.Field {
	field := &ast.Field{
		Type: typeExpr,
	}

	if len(names) > 0 {
		field.Names = make([]*ast.Ident, len(names))
		for i, name := range names {
			field.Names[i] = f.builder.ident(name)
		}
	}

	return field
}

// ParamVariadic creates a variadic function parameter
func (f *FunctionBuilder) ParamVariadic(name, typeName string) *ast.Field {
	return &ast.Field{
		Names: []*ast.Ident{f.builder.ident(name)},
		Type:  &ast.Ellipsis{Elt: f.builder.ident(typeName)},
	}
}

// ParamVariadicWithType creates a variadic function parameter with a type expression
func (f *FunctionBuilder) ParamVariadicWithType(name string, typeExpr ast.Expr) *ast.Field {
	return &ast.Field{
		Names: []*ast.Ident{f.builder.ident(name)},
		Type:  &ast.Ellipsis{Elt: typeExpr},
	}
}

// Result creates a function result
func (f *FunctionBuilder) Result(name, typeName string) *ast.Field {
	return &ast.Field{
		Names: []*ast.Ident{f.builder.ident(name)},
		Type:  f.builder.ident(typeName),
	}
}

// ResultWithType creates a function result with a type expression
func (f *FunctionBuilder) ResultWithType(name string, typeExpr ast.Expr) *ast.Field {
	return &ast.Field{
		Names: []*ast.Ident{f.builder.ident(name)},
		Type:  typeExpr,
	}
}

// ResultMultiple creates a function result with multiple names
func (f *FunctionBuilder) ResultMultiple(names []string, typeName string) *ast.Field {
	field := &ast.Field{
		Type: f.builder.ident(typeName),
	}

	if len(names) > 0 {
		field.Names = make([]*ast.Ident, len(names))
		for i, name := range names {
			field.Names[i] = f.builder.ident(name)
		}
	}

	return field
}

// ResultMultipleWithType creates a function result with multiple names and a type expression
func (f *FunctionBuilder) ResultMultipleWithType(names []string, typeExpr ast.Expr) *ast.Field {
	field := &ast.Field{
		Type: typeExpr,
	}

	if len(names) > 0 {
		field.Names = make([]*ast.Ident, len(names))
		for i, name := range names {
			field.Names[i] = f.builder.ident(name)
		}
	}

	return field
}

// ResultAnonymous creates an anonymous function result
func (f *FunctionBuilder) ResultAnonymous(typeName string) *ast.Field {
	return &ast.Field{
		Type: f.builder.ident(typeName),
	}
}

// ResultAnonymousWithType creates an anonymous function result with a type expression
func (f *FunctionBuilder) ResultAnonymousWithType(typeExpr ast.Expr) *ast.Field {
	return &ast.Field{
		Type: typeExpr,
	}
}

// Receiver creates a method receiver
func (f *FunctionBuilder) Receiver(name, typeName string) *ast.Field {
	return &ast.Field{
		Names: []*ast.Ident{f.builder.ident(name)},
		Type:  f.builder.ident(typeName),
	}
}

// ReceiverWithType creates a method receiver with a type expression
func (f *FunctionBuilder) ReceiverWithType(name string, typeExpr ast.Expr) *ast.Field {
	return &ast.Field{
		Names: []*ast.Ident{f.builder.ident(name)},
		Type:  typeExpr,
	}
}

// ReceiverPointer creates a pointer method receiver
func (f *FunctionBuilder) ReceiverPointer(name, typeName string) *ast.Field {
	return &ast.Field{
		Names: []*ast.Ident{f.builder.ident(name)},
		Type:  &ast.StarExpr{X: f.builder.ident(typeName)},
	}
}

// ReceiverPointerWithType creates a pointer method receiver with a type expression
func (f *FunctionBuilder) ReceiverPointerWithType(name string, typeExpr ast.Expr) *ast.Field {
	return &ast.Field{
		Names: []*ast.Ident{f.builder.ident(name)},
		Type:  &ast.StarExpr{X: typeExpr},
	}
}

// InterfaceMethod creates a method signature for interfaces
func (f *FunctionBuilder) InterfaceMethod(name string, params, results []*ast.Field) *ast.Field {
	return &ast.Field{
		Names: []*ast.Ident{f.builder.ident(name)},
		Type: &ast.FuncType{
			Params:  &ast.FieldList{List: params},
			Results: &ast.FieldList{List: results},
		},
	}
}

// InterfaceMethodWithType creates a method signature for interfaces with a function type
func (f *FunctionBuilder) InterfaceMethodWithType(name string, funcType ast.Expr) *ast.Field {
	return &ast.Field{
		Names: []*ast.Ident{f.builder.ident(name)},
		Type:  funcType,
	}
}

// FieldList creates a field list
func (f *FunctionBuilder) FieldList(fields []*ast.Field) *ast.FieldList {
	return &ast.FieldList{
		List: fields,
	}
}

// EmptyFieldList creates an empty field list
func (f *FunctionBuilder) EmptyFieldList() *ast.FieldList {
	return &ast.FieldList{
		List: []*ast.Field{},
	}
}

// BuildFunction creates a function declaration with the current builder statements
func (f *FunctionBuilder) BuildFunction(name string, receiver *ast.Field, params, results []*ast.Field) *ast.FuncDecl {
	return f.Func(name, receiver, params, results, f.builder.stmts)
}

// BuildMethod creates a method declaration with the current builder statements
func (f *FunctionBuilder) BuildMethod(receiver *ast.Field, name string, params, results []*ast.Field) *ast.FuncDecl {
	return f.Method(receiver, name, params, results, f.builder.stmts)
}

// BuildFunctionWithoutReceiver creates a function declaration without receiver
func (f *FunctionBuilder) BuildFunctionWithoutReceiver(name string, params, results []*ast.Field) *ast.FuncDecl {
	return f.Function(name, params, results, f.builder.stmts)
}

// AddToBuilder adds the function declaration to the builder
func (f *FunctionBuilder) AddToBuilder(name string, receiver *ast.Field, params, results []*ast.Field, body []ast.Stmt) *FunctionBuilder {
	funcDecl := f.Func(name, receiver, params, results, body)
	f.builder.AddDeclaration(funcDecl)
	return f
}

// AddMethodToBuilder adds the method declaration to the builder
func (f *FunctionBuilder) AddMethodToBuilder(receiver *ast.Field, name string, params, results []*ast.Field, body []ast.Stmt) *FunctionBuilder {
	return f.AddToBuilder(name, receiver, params, results, body)
}

// AddFunctionToBuilder adds the function declaration to the builder
func (f *FunctionBuilder) AddFunctionToBuilder(name string, params, results []*ast.Field, body []ast.Stmt) *FunctionBuilder {
	return f.AddToBuilder(name, nil, params, results, body)
}

// Helper method to get the underlying builder
func (f *FunctionBuilder) Builder() *Builder {
	return f.builder
}
