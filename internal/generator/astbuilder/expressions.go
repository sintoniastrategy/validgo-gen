package astbuilder

import (
	"go/ast"
	"go/token"
	"strconv"
)

// ExpressionBuilder provides methods for building AST expressions
type ExpressionBuilder struct {
	builder *Builder
}

// NewExpressionBuilder creates a new expression builder
func NewExpressionBuilder(builder *Builder) *ExpressionBuilder {
	return &ExpressionBuilder{builder: builder}
}

// Ident creates an identifier expression
func (e *ExpressionBuilder) Ident(name string) ast.Expr {
	return e.builder.ident(name)
}

// String creates a string literal expression
func (e *ExpressionBuilder) String(value string) ast.Expr {
	return e.builder.str(value)
}

// Int creates an integer literal expression
func (e *ExpressionBuilder) Int(value int) ast.Expr {
	return &ast.BasicLit{
		Kind:  token.INT,
		Value: strconv.Itoa(value),
	}
}

// Int64 creates an int64 literal expression
func (e *ExpressionBuilder) Int64(value int64) ast.Expr {
	return &ast.BasicLit{
		Kind:  token.INT,
		Value: strconv.FormatInt(value, 10),
	}
}

// Float creates a float literal expression
func (e *ExpressionBuilder) Float(value float64) ast.Expr {
	return &ast.BasicLit{
		Kind:  token.FLOAT,
		Value: strconv.FormatFloat(value, 'f', -1, 64),
	}
}

// Bool creates a boolean literal expression
func (e *ExpressionBuilder) Bool(value bool) ast.Expr {
	return &ast.Ident{
		Name: strconv.FormatBool(value),
	}
}

// Nil creates a nil literal expression
func (e *ExpressionBuilder) Nil() ast.Expr {
	return e.builder.ident("nil")
}

// True creates a true boolean literal expression
func (e *ExpressionBuilder) True() ast.Expr {
	return e.builder.ident("true")
}

// False creates a false boolean literal expression
func (e *ExpressionBuilder) False() ast.Expr {
	return e.builder.ident("false")
}

// Select creates a selector expression (e.g., obj.field)
func (e *ExpressionBuilder) Select(receiver ast.Expr, field string) ast.Expr {
	return e.builder.selector(receiver, field)
}

// Call creates a call expression (e.g., func(args...))
func (e *ExpressionBuilder) Call(fun ast.Expr, args ...ast.Expr) ast.Expr {
	return e.builder.call(fun, args...)
}

// MethodCall creates a method call expression (e.g., obj.method(args...))
func (e *ExpressionBuilder) MethodCall(receiver ast.Expr, method string, args ...ast.Expr) ast.Expr {
	return e.builder.call(e.builder.selector(receiver, method), args...)
}

// Index creates an index expression (e.g., array[index])
func (e *ExpressionBuilder) Index(array ast.Expr, index ast.Expr) ast.Expr {
	return &ast.IndexExpr{
		X:     array,
		Index: index,
	}
}

// Slice creates a slice expression (e.g., array[low:high])
func (e *ExpressionBuilder) Slice(array ast.Expr, low, high ast.Expr) ast.Expr {
	return &ast.SliceExpr{
		X:    array,
		Low:  low,
		High: high,
	}
}

// AddressOf creates an address-of expression (e.g., &value)
func (e *ExpressionBuilder) AddressOf(expr ast.Expr) ast.Expr {
	return e.builder.unary(token.AND, expr)
}

// Deref creates a dereference expression (e.g., *ptr)
func (e *ExpressionBuilder) Deref(expr ast.Expr) ast.Expr {
	return e.builder.unary(token.MUL, expr)
}

// Not creates a logical not expression (e.g., !expr)
func (e *ExpressionBuilder) Not(expr ast.Expr) ast.Expr {
	return e.builder.unary(token.NOT, expr)
}

// Plus creates a unary plus expression (e.g., +expr)
func (e *ExpressionBuilder) Plus(expr ast.Expr) ast.Expr {
	return e.builder.unary(token.ADD, expr)
}

// Minus creates a unary minus expression (e.g., -expr)
func (e *ExpressionBuilder) Minus(expr ast.Expr) ast.Expr {
	return e.builder.unary(token.SUB, expr)
}

// Equal creates an equality expression (e.g., x == y)
func (e *ExpressionBuilder) Equal(left, right ast.Expr) ast.Expr {
	return e.builder.binary(left, token.EQL, right)
}

// NotEqual creates an inequality expression (e.g., x != y)
func (e *ExpressionBuilder) NotEqual(left, right ast.Expr) ast.Expr {
	return e.builder.binary(left, token.NEQ, right)
}

// Less creates a less than expression (e.g., x < y)
func (e *ExpressionBuilder) Less(left, right ast.Expr) ast.Expr {
	return e.builder.binary(left, token.LSS, right)
}

// LessEqual creates a less than or equal expression (e.g., x <= y)
func (e *ExpressionBuilder) LessEqual(left, right ast.Expr) ast.Expr {
	return e.builder.binary(left, token.LEQ, right)
}

// Greater creates a greater than expression (e.g., x > y)
func (e *ExpressionBuilder) Greater(left, right ast.Expr) ast.Expr {
	return e.builder.binary(left, token.GTR, right)
}

// GreaterEqual creates a greater than or equal expression (e.g., x >= y)
func (e *ExpressionBuilder) GreaterEqual(left, right ast.Expr) ast.Expr {
	return e.builder.binary(left, token.GEQ, right)
}

// Add creates an addition expression (e.g., x + y)
func (e *ExpressionBuilder) Add(left, right ast.Expr) ast.Expr {
	return e.builder.binary(left, token.ADD, right)
}

// Sub creates a subtraction expression (e.g., x - y)
func (e *ExpressionBuilder) Sub(left, right ast.Expr) ast.Expr {
	return e.builder.binary(left, token.SUB, right)
}

// Mul creates a multiplication expression (e.g., x * y)
func (e *ExpressionBuilder) Mul(left, right ast.Expr) ast.Expr {
	return e.builder.binary(left, token.MUL, right)
}

// Div creates a division expression (e.g., x / y)
func (e *ExpressionBuilder) Div(left, right ast.Expr) ast.Expr {
	return e.builder.binary(left, token.QUO, right)
}

// Mod creates a modulo expression (e.g., x % y)
func (e *ExpressionBuilder) Mod(left, right ast.Expr) ast.Expr {
	return e.builder.binary(left, token.REM, right)
}

// And creates a logical AND expression (e.g., x && y)
func (e *ExpressionBuilder) And(left, right ast.Expr) ast.Expr {
	return e.builder.binary(left, token.LAND, right)
}

// Or creates a logical OR expression (e.g., x || y)
func (e *ExpressionBuilder) Or(left, right ast.Expr) ast.Expr {
	return e.builder.binary(left, token.LOR, right)
}

// BitwiseAnd creates a bitwise AND expression (e.g., x & y)
func (e *ExpressionBuilder) BitwiseAnd(left, right ast.Expr) ast.Expr {
	return e.builder.binary(left, token.AND, right)
}

// BitwiseOr creates a bitwise OR expression (e.g., x | y)
func (e *ExpressionBuilder) BitwiseOr(left, right ast.Expr) ast.Expr {
	return e.builder.binary(left, token.OR, right)
}

// BitwiseXor creates a bitwise XOR expression (e.g., x ^ y)
func (e *ExpressionBuilder) BitwiseXor(left, right ast.Expr) ast.Expr {
	return e.builder.binary(left, token.XOR, right)
}

// LeftShift creates a left shift expression (e.g., x << y)
func (e *ExpressionBuilder) LeftShift(left, right ast.Expr) ast.Expr {
	return e.builder.binary(left, token.SHL, right)
}

// RightShift creates a right shift expression (e.g., x >> y)
func (e *ExpressionBuilder) RightShift(left, right ast.Expr) ast.Expr {
	return e.builder.binary(left, token.SHR, right)
}

// TypeAssert creates a type assertion expression (e.g., x.(T))
func (e *ExpressionBuilder) TypeAssert(expr ast.Expr, typeName string) ast.Expr {
	return &ast.TypeAssertExpr{
		X:    expr,
		Type: e.builder.ident(typeName),
	}
}

// TypeAssertWithType creates a type assertion expression with a type (e.g., x.(T))
func (e *ExpressionBuilder) TypeAssertWithType(expr ast.Expr, typeExpr ast.Expr) ast.Expr {
	return &ast.TypeAssertExpr{
		X:    expr,
		Type: typeExpr,
	}
}

// CompositeLit creates a composite literal expression (e.g., Type{field: value})
func (e *ExpressionBuilder) CompositeLit(typeName string, elements ...ast.Expr) ast.Expr {
	return &ast.CompositeLit{
		Type: e.builder.ident(typeName),
		Elts: elements,
	}
}

// CompositeLitWithType creates a composite literal expression with a type
func (e *ExpressionBuilder) CompositeLitWithType(typeExpr ast.Expr, elements ...ast.Expr) ast.Expr {
	return &ast.CompositeLit{
		Type: typeExpr,
		Elts: elements,
	}
}

// KeyValue creates a key-value expression for composite literals
func (e *ExpressionBuilder) KeyValue(key, value ast.Expr) ast.Expr {
	return &ast.KeyValueExpr{
		Key:   key,
		Value: value,
	}
}

// ArrayType creates an array type expression (e.g., [n]T)
func (e *ExpressionBuilder) ArrayType(length ast.Expr, elementType ast.Expr) ast.Expr {
	return &ast.ArrayType{
		Len: length,
		Elt: elementType,
	}
}

// SliceType creates a slice type expression (e.g., []T)
func (e *ExpressionBuilder) SliceType(elementType ast.Expr) ast.Expr {
	return &ast.ArrayType{
		Elt: elementType,
	}
}

// MapType creates a map type expression (e.g., map[K]V)
func (e *ExpressionBuilder) MapType(keyType, valueType ast.Expr) ast.Expr {
	return &ast.MapType{
		Key:   keyType,
		Value: valueType,
	}
}

// ChanType creates a channel type expression (e.g., chan T)
func (e *ExpressionBuilder) ChanType(valueType ast.Expr, dir ast.ChanDir) ast.Expr {
	return &ast.ChanType{
		Value: valueType,
		Dir:   dir,
	}
}

// Star creates a pointer type expression (e.g., *T)
func (e *ExpressionBuilder) Star(typeExpr ast.Expr) ast.Expr {
	return &ast.StarExpr{
		X: typeExpr,
	}
}

// FuncType creates a function type expression
func (e *ExpressionBuilder) FuncType(params, results []*ast.Field) ast.Expr {
	return &ast.FuncType{
		Params:  &ast.FieldList{List: params},
		Results: &ast.FieldList{List: results},
	}
}

// InterfaceType creates an interface type expression
func (e *ExpressionBuilder) InterfaceType(methods []*ast.Field) ast.Expr {
	return &ast.InterfaceType{
		Methods: &ast.FieldList{List: methods},
	}
}

// StructType creates a struct type expression
func (e *ExpressionBuilder) StructType(fields []*ast.Field) ast.Expr {
	return &ast.StructType{
		Fields: &ast.FieldList{List: fields},
	}
}

// Paren creates a parenthesized expression (e.g., (expr))
func (e *ExpressionBuilder) Paren(expr ast.Expr) ast.Expr {
	return &ast.ParenExpr{
		X: expr,
	}
}

// Ellipsis creates an ellipsis expression (e.g., ...T)
func (e *ExpressionBuilder) Ellipsis(typeExpr ast.Expr) ast.Expr {
	return &ast.Ellipsis{
		Elt: typeExpr,
	}
}

// Chan creates a channel type expression (e.g., chan T)
func (e *ExpressionBuilder) Chan(valueType ast.Expr, dir ast.ChanDir) ast.Expr {
	return &ast.ChanType{
		Value: valueType,
		Dir:   dir,
	}
}

// Helper method to get the underlying builder
func (e *ExpressionBuilder) Builder() *Builder {
	return e.builder
}
