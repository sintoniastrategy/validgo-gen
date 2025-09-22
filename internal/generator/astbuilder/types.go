package astbuilder

import (
	"go/ast"
	"go/token"
	"strings"
)

// TypeBuilder provides methods for building AST type expressions and declarations
type TypeBuilder struct {
	builder *Builder
}

// NewTypeBuilder creates a new type builder
func NewTypeBuilder(builder *Builder) *TypeBuilder {
	return &TypeBuilder{builder: builder}
}

// Ident creates a type identifier
func (t *TypeBuilder) Ident(name string) ast.Expr {
	return t.builder.ident(name)
}

// Select creates a selector type expression (e.g., pkg.Type)
func (t *TypeBuilder) Select(pkg, typeName string) ast.Expr {
	return t.builder.selector(t.builder.ident(pkg), typeName)
}

// Star creates a pointer type (e.g., *T)
func (t *TypeBuilder) Star(typeExpr ast.Expr) ast.Expr {
	return &ast.StarExpr{
		X: typeExpr,
	}
}

// Array creates an array type (e.g., [n]T)
func (t *TypeBuilder) Array(length ast.Expr, elementType ast.Expr) ast.Expr {
	return &ast.ArrayType{
		Len: length,
		Elt: elementType,
	}
}

// Slice creates a slice type (e.g., []T)
func (t *TypeBuilder) Slice(elementType ast.Expr) ast.Expr {
	return &ast.ArrayType{
		Elt: elementType,
	}
}

// Map creates a map type (e.g., map[K]V)
func (t *TypeBuilder) Map(keyType, valueType ast.Expr) ast.Expr {
	return &ast.MapType{
		Key:   keyType,
		Value: valueType,
	}
}

// Chan creates a channel type (e.g., chan T)
func (t *TypeBuilder) Chan(valueType ast.Expr, dir ast.ChanDir) ast.Expr {
	return &ast.ChanType{
		Value: valueType,
		Dir:   dir,
	}
}

// SendChan creates a send-only channel type (e.g., chan<- T)
func (t *TypeBuilder) SendChan(valueType ast.Expr) ast.Expr {
	return t.Chan(valueType, ast.SEND)
}

// RecvChan creates a receive-only channel type (e.g., <-chan T)
func (t *TypeBuilder) RecvChan(valueType ast.Expr) ast.Expr {
	return t.Chan(valueType, ast.RECV)
}

// Func creates a function type
func (t *TypeBuilder) Func(params, results []*ast.Field) ast.Expr {
	return &ast.FuncType{
		Params:  &ast.FieldList{List: params},
		Results: &ast.FieldList{List: results},
	}
}

// Interface creates an interface type
func (t *TypeBuilder) Interface(methods []*ast.Field) ast.Expr {
	return &ast.InterfaceType{
		Methods: &ast.FieldList{List: methods},
	}
}

// Struct creates a struct type
func (t *TypeBuilder) Struct(fields []*ast.Field) ast.Expr {
	return &ast.StructType{
		Fields: &ast.FieldList{List: fields},
	}
}

// Field creates a struct field
func (t *TypeBuilder) Field(name string, typeExpr ast.Expr, tag string) *ast.Field {
	field := &ast.Field{
		Type: typeExpr,
	}

	if name != "" {
		field.Names = []*ast.Ident{t.builder.ident(name)}
	}

	if tag != "" {
		field.Tag = &ast.BasicLit{
			Kind:  token.STRING,
			Value: "`" + tag + "`",
		}
	}

	return field
}

// FieldWithNames creates a struct field with multiple names
func (t *TypeBuilder) FieldWithNames(names []string, typeExpr ast.Expr, tag string) *ast.Field {
	field := &ast.Field{
		Type: typeExpr,
	}

	if len(names) > 0 {
		field.Names = make([]*ast.Ident, len(names))
		for i, name := range names {
			field.Names[i] = t.builder.ident(name)
		}
	}

	if tag != "" {
		field.Tag = &ast.BasicLit{
			Kind:  token.STRING,
			Value: "`" + tag + "`",
		}
	}

	return field
}

// FieldList creates a field list
func (t *TypeBuilder) FieldList(fields []*ast.Field) *ast.FieldList {
	return &ast.FieldList{
		List: fields,
	}
}

// Method creates a method signature for interfaces
func (t *TypeBuilder) Method(name string, params, results []*ast.Field) *ast.Field {
	return &ast.Field{
		Names: []*ast.Ident{t.builder.ident(name)},
		Type:  t.Func(params, results),
	}
}

// Ellipsis creates an ellipsis type (e.g., ...T)
func (t *TypeBuilder) Ellipsis(typeExpr ast.Expr) ast.Expr {
	return &ast.Ellipsis{
		Elt: typeExpr,
	}
}

// Paren creates a parenthesized type (e.g., (T))
func (t *TypeBuilder) Paren(typeExpr ast.Expr) ast.Expr {
	return &ast.ParenExpr{
		X: typeExpr,
	}
}

// TypeDecl creates a type declaration
func (t *TypeBuilder) TypeDecl(name string, typeExpr ast.Expr) ast.Decl {
	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: t.builder.ident(name),
				Type: typeExpr,
			},
		},
	}
}

// TypeAlias creates a type alias declaration
func (t *TypeBuilder) TypeAlias(name, alias string) ast.Decl {
	return t.TypeDecl(name, t.Ident(alias))
}

// TypeAliasWithType creates a type alias declaration with a type expression
func (t *TypeBuilder) TypeAliasWithType(name string, typeExpr ast.Expr) ast.Decl {
	return t.TypeDecl(name, typeExpr)
}

// SliceAlias creates a slice type alias
func (t *TypeBuilder) SliceAlias(name, elementType string) ast.Decl {
	return t.TypeDecl(name, t.Slice(t.Ident(elementType)))
}

// SliceAliasWithType creates a slice type alias with a type expression
func (t *TypeBuilder) SliceAliasWithType(name string, elementType ast.Expr) ast.Decl {
	return t.TypeDecl(name, t.Slice(elementType))
}

// ArrayAlias creates an array type alias
func (t *TypeBuilder) ArrayAlias(name string, length ast.Expr, elementType string) ast.Decl {
	return t.TypeDecl(name, t.Array(length, t.Ident(elementType)))
}

// ArrayAliasWithType creates an array type alias with a type expression
func (t *TypeBuilder) ArrayAliasWithType(name string, length ast.Expr, elementType ast.Expr) ast.Decl {
	return t.TypeDecl(name, t.Array(length, elementType))
}

// MapAlias creates a map type alias
func (t *TypeBuilder) MapAlias(name, keyType, valueType string) ast.Decl {
	return t.TypeDecl(name, t.Map(t.Ident(keyType), t.Ident(valueType)))
}

// MapAliasWithType creates a map type alias with type expressions
func (t *TypeBuilder) MapAliasWithType(name string, keyType, valueType ast.Expr) ast.Decl {
	return t.TypeDecl(name, t.Map(keyType, valueType))
}

// ChanAlias creates a channel type alias
func (t *TypeBuilder) ChanAlias(name, valueType string, dir ast.ChanDir) ast.Decl {
	return t.TypeDecl(name, t.Chan(t.Ident(valueType), dir))
}

// ChanAliasWithType creates a channel type alias with a type expression
func (t *TypeBuilder) ChanAliasWithType(name string, valueType ast.Expr, dir ast.ChanDir) ast.Decl {
	return t.TypeDecl(name, t.Chan(valueType, dir))
}

// StructAlias creates a struct type alias
func (t *TypeBuilder) StructAlias(name string, fields []*ast.Field) ast.Decl {
	return t.TypeDecl(name, t.Struct(fields))
}

// InterfaceAlias creates an interface type alias
func (t *TypeBuilder) InterfaceAlias(name string, methods []*ast.Field) ast.Decl {
	return t.TypeDecl(name, t.Interface(methods))
}

// FuncAlias creates a function type alias
func (t *TypeBuilder) FuncAlias(name string, params, results []*ast.Field) ast.Decl {
	return t.TypeDecl(name, t.Func(params, results))
}

// BuildTags creates struct tags from a map
func (t *TypeBuilder) BuildTags(tags map[string]string) string {
	if len(tags) == 0 {
		return ""
	}

	var parts []string
	for key, value := range tags {
		parts = append(parts, key+`:"`+value+`"`)
	}

	return strings.Join(parts, " ")
}

// BuildJSONTag creates a JSON tag
func (t *TypeBuilder) BuildJSONTag(name string, omitempty bool) string {
	tag := `json:"` + name + `"`
	if omitempty {
		tag += `,omitempty`
	}
	return tag
}

// BuildValidateTag creates a validation tag
func (t *TypeBuilder) BuildValidateTag(rules ...string) string {
	if len(rules) == 0 {
		return ""
	}
	return `validate:"` + strings.Join(rules, ",") + `"`
}

// BuildTagsFromFields creates struct tags from field specifications
func (t *TypeBuilder) BuildTagsFromFields(jsonName string, omitempty bool, validateRules ...string) string {
	tags := make(map[string]string)

	if jsonName != "" {
		jsonTag := jsonName
		if omitempty {
			jsonTag += ",omitempty"
		}
		tags["json"] = jsonTag
	}

	if len(validateRules) > 0 {
		tags["validate"] = strings.Join(validateRules, ",")
	}

	return t.BuildTags(tags)
}

// Helper method to get the underlying builder
func (t *TypeBuilder) Builder() *Builder {
	return t.builder
}
