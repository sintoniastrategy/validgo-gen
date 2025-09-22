package astbuilder

import (
	"go/ast"
)

// PatternBuilder provides methods for building common AST patterns
type PatternBuilder struct {
	builder *Builder
}

// NewPatternBuilder creates a new pattern builder
func NewPatternBuilder(builder *Builder) *PatternBuilder {
	return &PatternBuilder{builder: builder}
}

// ErrorHandlingPattern creates a common error handling pattern
func (p *PatternBuilder) ErrorHandlingPattern(varName string, callExpr ast.Expr, errorMessage string) []ast.Stmt {
	stmtBuilder := NewStatementBuilder(p.builder)
	exprBuilder := NewExpressionBuilder(p.builder)

	return []ast.Stmt{
		stmtBuilder.AssignDefine(
			exprBuilder.Ident(varName),
			callExpr,
		),
		stmtBuilder.If(
			exprBuilder.NotEqual(exprBuilder.Ident(varName), exprBuilder.Nil()),
			[]ast.Stmt{
				stmtBuilder.Return(exprBuilder.Nil(), exprBuilder.Ident(varName)),
			},
		),
	}
}

// ValidationPattern creates a common validation pattern
func (p *PatternBuilder) ValidationPattern(structName string, validatorName string) []ast.Stmt {
	stmtBuilder := NewStatementBuilder(p.builder)
	exprBuilder := NewExpressionBuilder(p.builder)

	return []ast.Stmt{
		stmtBuilder.AssignDefine(
			exprBuilder.Ident("err"),
			exprBuilder.MethodCall(
				exprBuilder.Select(exprBuilder.Ident("h"), validatorName),
				"Struct",
				exprBuilder.Ident(structName),
			),
		),
		stmtBuilder.If(
			exprBuilder.NotEqual(exprBuilder.Ident("err"), exprBuilder.Nil()),
			[]ast.Stmt{
				stmtBuilder.Return(exprBuilder.Nil(), exprBuilder.Ident("err")),
			},
		),
	}
}

// JSONUnmarshalPattern creates a common JSON unmarshaling pattern
func (p *PatternBuilder) JSONUnmarshalPattern(varName string, jsonData ast.Expr, target ast.Expr) []ast.Stmt {
	stmtBuilder := NewStatementBuilder(p.builder)
	exprBuilder := NewExpressionBuilder(p.builder)

	p.builder.AddImport("encoding/json")

	return []ast.Stmt{
		stmtBuilder.AssignDefine(
			exprBuilder.Ident("err"),
			exprBuilder.Call(
				exprBuilder.Select(exprBuilder.Ident("json"), "Unmarshal"),
				jsonData,
				exprBuilder.AddressOf(target),
			),
		),
		stmtBuilder.If(
			exprBuilder.NotEqual(exprBuilder.Ident("err"), exprBuilder.Nil()),
			[]ast.Stmt{
				stmtBuilder.Return(exprBuilder.Nil(), exprBuilder.Ident("err")),
			},
		),
	}
}

// JSONMarshalPattern creates a common JSON marshaling pattern
func (p *PatternBuilder) JSONMarshalPattern(varName string, data ast.Expr) []ast.Stmt {
	stmtBuilder := NewStatementBuilder(p.builder)
	exprBuilder := NewExpressionBuilder(p.builder)

	p.builder.AddImport("encoding/json")

	return []ast.Stmt{
		stmtBuilder.AssignDefine(
			exprBuilder.Ident(varName),
			exprBuilder.Call(
				exprBuilder.Select(exprBuilder.Ident("json"), "Marshal"),
				data,
			),
		),
		stmtBuilder.If(
			exprBuilder.NotEqual(exprBuilder.Ident("err"), exprBuilder.Nil()),
			[]ast.Stmt{
				stmtBuilder.Return(exprBuilder.Nil(), exprBuilder.Ident("err")),
			},
		),
	}
}

// HTTPErrorPattern creates a common HTTP error response pattern
func (p *PatternBuilder) HTTPErrorPattern(writer, message ast.Expr, statusCode ast.Expr) []ast.Stmt {
	stmtBuilder := NewStatementBuilder(p.builder)
	exprBuilder := NewExpressionBuilder(p.builder)

	p.builder.AddImport("net/http")

	return []ast.Stmt{
		stmtBuilder.CallStmt(
			exprBuilder.Select(exprBuilder.Ident("http"), "Error"),
			writer,
			message,
			statusCode,
		),
		stmtBuilder.ReturnEmpty(),
	}
}

// HTTPWriteHeaderPattern creates a common HTTP header writing pattern
func (p *PatternBuilder) HTTPWriteHeaderPattern(writer, headerName, headerValue ast.Expr) []ast.Stmt {
	stmtBuilder := NewStatementBuilder(p.builder)
	exprBuilder := NewExpressionBuilder(p.builder)

	return []ast.Stmt{
		stmtBuilder.CallStmt(
			exprBuilder.MethodCall(
				exprBuilder.MethodCall(writer, "Header"),
				"Set",
				headerName,
				headerValue,
			),
		),
	}
}

// HTTPWritePattern creates a common HTTP response writing pattern
func (p *PatternBuilder) HTTPWritePattern(writer, data ast.Expr) []ast.Stmt {
	stmtBuilder := NewStatementBuilder(p.builder)
	exprBuilder := NewExpressionBuilder(p.builder)

	return []ast.Stmt{
		stmtBuilder.CallStmt(
			exprBuilder.MethodCall(writer, "Write"),
			data,
		),
	}
}

// HTTPWriteHeaderAndDataPattern creates a pattern for writing HTTP headers and data
func (p *PatternBuilder) HTTPWriteHeaderAndDataPattern(writer, statusCode, data ast.Expr) []ast.Stmt {
	stmtBuilder := NewStatementBuilder(p.builder)
	exprBuilder := NewExpressionBuilder(p.builder)

	return []ast.Stmt{
		stmtBuilder.CallStmt(
			exprBuilder.MethodCall(writer, "WriteHeader"),
			statusCode,
		),
		stmtBuilder.CallStmt(
			exprBuilder.MethodCall(writer, "Write"),
			data,
		),
	}
}

// ParameterExtractionPattern creates a common parameter extraction pattern
func (p *PatternBuilder) ParameterExtractionPattern(paramName string, source ast.Expr, method string, required bool) []ast.Stmt {
	stmtBuilder := NewStatementBuilder(p.builder)
	exprBuilder := NewExpressionBuilder(p.builder)

	stmts := []ast.Stmt{
		stmtBuilder.AssignDefine(
			exprBuilder.Ident(paramName),
			exprBuilder.MethodCall(source, method),
		),
	}

	if required {
		stmts = append(stmts, stmtBuilder.If(
			exprBuilder.Equal(exprBuilder.Ident(paramName), exprBuilder.String("")),
			[]ast.Stmt{
				stmtBuilder.Return(
					exprBuilder.Nil(),
					exprBuilder.Call(
						exprBuilder.Select(exprBuilder.Ident("errors"), "New"),
						exprBuilder.String(paramName+" is required"),
					),
				),
			},
		))
		p.builder.AddImport("github.com/go-faster/errors")
	}

	return stmts
}

// QueryParamExtractionPattern creates a query parameter extraction pattern
func (p *PatternBuilder) QueryParamExtractionPattern(paramName, paramKey string, required bool) []ast.Stmt {
	exprBuilder := NewExpressionBuilder(p.builder)

	return p.ParameterExtractionPattern(
		paramName,
		exprBuilder.MethodCall(exprBuilder.Ident("r"), "URL"),
		"Query",
		required,
	)
}

// HeaderExtractionPattern creates a header extraction pattern
func (p *PatternBuilder) HeaderExtractionPattern(paramName, headerName string, required bool) []ast.Stmt {
	exprBuilder := NewExpressionBuilder(p.builder)

	return p.ParameterExtractionPattern(
		paramName,
		exprBuilder.MethodCall(exprBuilder.Ident("r"), "Header"),
		"Get",
		required,
	)
}

// CookieExtractionPattern creates a cookie extraction pattern
func (p *PatternBuilder) CookieExtractionPattern(paramName, cookieName string, required bool) []ast.Stmt {
	exprBuilder := NewExpressionBuilder(p.builder)

	return p.ParameterExtractionPattern(
		paramName,
		exprBuilder.MethodCall(exprBuilder.Ident("r"), "Header"),
		"Get",
		required,
	)
}

// PathParamExtractionPattern creates a path parameter extraction pattern
func (p *PatternBuilder) PathParamExtractionPattern(paramName, paramKey string, required bool) []ast.Stmt {
	exprBuilder := NewExpressionBuilder(p.builder)

	p.builder.AddImport("github.com/go-chi/chi/v5")

	return p.ParameterExtractionPattern(
		paramName,
		exprBuilder.Call(exprBuilder.Select(exprBuilder.Ident("chi"), "URLParam"), exprBuilder.Ident("r"), exprBuilder.String(paramKey)),
		"",
		required,
	)
}

// StructInitializationPattern creates a struct initialization pattern
func (p *PatternBuilder) StructInitializationPattern(structName string, fields map[string]ast.Expr) []ast.Stmt {
	stmtBuilder := NewStatementBuilder(p.builder)
	exprBuilder := NewExpressionBuilder(p.builder)

	// Create composite literal elements
	var elements []ast.Expr
	for fieldName, value := range fields {
		elements = append(elements, exprBuilder.KeyValue(
			exprBuilder.String(fieldName),
			value,
		))
	}

	return []ast.Stmt{
		stmtBuilder.AssignDefine(
			exprBuilder.Ident(structName),
			exprBuilder.CompositeLit(structName, elements...),
		),
	}
}

// MapInitializationPattern creates a map initialization pattern
func (p *PatternBuilder) MapInitializationPattern(mapName, keyType, valueType string, entries map[string]ast.Expr) []ast.Stmt {
	stmtBuilder := NewStatementBuilder(p.builder)
	exprBuilder := NewExpressionBuilder(p.builder)

	// Create map type
	mapType := exprBuilder.MapType(
		exprBuilder.Ident(keyType),
		exprBuilder.Ident(valueType),
	)

	// Create composite literal elements
	var elements []ast.Expr
	for key, value := range entries {
		elements = append(elements, exprBuilder.KeyValue(
			exprBuilder.String(key),
			value,
		))
	}

	return []ast.Stmt{
		stmtBuilder.AssignDefine(
			exprBuilder.Ident(mapName),
			exprBuilder.CompositeLitWithType(mapType, elements...),
		),
	}
}

// SliceInitializationPattern creates a slice initialization pattern
func (p *PatternBuilder) SliceInitializationPattern(sliceName, elementType string, elements []ast.Expr) []ast.Stmt {
	stmtBuilder := NewStatementBuilder(p.builder)
	exprBuilder := NewExpressionBuilder(p.builder)

	// Create slice type
	sliceType := exprBuilder.SliceType(exprBuilder.Ident(elementType))

	return []ast.Stmt{
		stmtBuilder.AssignDefine(
			exprBuilder.Ident(sliceName),
			exprBuilder.CompositeLitWithType(sliceType, elements...),
		),
	}
}

// SwitchPattern creates a switch statement pattern
func (p *PatternBuilder) SwitchPattern(tag ast.Expr, cases map[ast.Expr][]ast.Stmt, defaultCase []ast.Stmt) ast.Stmt {
	stmtBuilder := NewStatementBuilder(p.builder)

	var caseStmts []ast.Stmt
	for value, body := range cases {
		caseStmts = append(caseStmts, stmtBuilder.Case([]ast.Expr{value}, body))
	}

	if len(defaultCase) > 0 {
		caseStmts = append(caseStmts, stmtBuilder.Default(defaultCase))
	}

	return stmtBuilder.Switch(tag, caseStmts)
}

// TypeSwitchPattern creates a type switch statement pattern
func (p *PatternBuilder) TypeSwitchPattern(assign ast.Stmt, cases map[string][]ast.Stmt, defaultCase []ast.Stmt) ast.Stmt {
	stmtBuilder := NewStatementBuilder(p.builder)
	exprBuilder := NewExpressionBuilder(p.builder)

	var caseStmts []ast.Stmt
	for typeName, body := range cases {
		caseStmts = append(caseStmts, stmtBuilder.Case([]ast.Expr{exprBuilder.Ident(typeName)}, body))
	}

	if len(defaultCase) > 0 {
		caseStmts = append(caseStmts, stmtBuilder.Default(defaultCase))
	}

	return stmtBuilder.TypeSwitch(assign, caseStmts)
}

// Helper method to get the underlying builder
func (p *PatternBuilder) Builder() *Builder {
	return p.builder
}
