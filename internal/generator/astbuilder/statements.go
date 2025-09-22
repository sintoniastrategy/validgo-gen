package astbuilder

import (
	"go/ast"
	"go/token"
)

// StatementBuilder provides methods for building AST statements
type StatementBuilder struct {
	builder *Builder
}

// NewStatementBuilder creates a new statement builder
func NewStatementBuilder(builder *Builder) *StatementBuilder {
	return &StatementBuilder{builder: builder}
}

// DeclareVar creates a variable declaration statement
func (s *StatementBuilder) DeclareVar(name, typeName string, value ast.Expr) ast.Stmt {
	spec := &ast.ValueSpec{
		Names: []*ast.Ident{s.builder.ident(name)},
		Type:  s.builder.ident(typeName),
	}

	if value != nil {
		spec.Values = []ast.Expr{value}
	}

	return &ast.DeclStmt{
		Decl: &ast.GenDecl{
			Tok:   token.VAR,
			Specs: []ast.Spec{spec},
		},
	}
}

// DeclareVarWithType creates a variable declaration statement with a type expression
func (s *StatementBuilder) DeclareVarWithType(name string, typeExpr ast.Expr, value ast.Expr) ast.Stmt {
	spec := &ast.ValueSpec{
		Names: []*ast.Ident{s.builder.ident(name)},
		Type:  typeExpr,
	}

	if value != nil {
		spec.Values = []ast.Expr{value}
	}

	return &ast.DeclStmt{
		Decl: &ast.GenDecl{
			Tok:   token.VAR,
			Specs: []ast.Spec{spec},
		},
	}
}

// DeclareConst creates a constant declaration statement
func (s *StatementBuilder) DeclareConst(name, typeName string, value ast.Expr) ast.Stmt {
	spec := &ast.ValueSpec{
		Names:  []*ast.Ident{s.builder.ident(name)},
		Type:   s.builder.ident(typeName),
		Values: []ast.Expr{value},
	}

	return &ast.DeclStmt{
		Decl: &ast.GenDecl{
			Tok:   token.CONST,
			Specs: []ast.Spec{spec},
		},
	}
}

// Assign creates an assignment statement
func (s *StatementBuilder) Assign(lhs, rhs ast.Expr) ast.Stmt {
	return &ast.AssignStmt{
		Lhs: []ast.Expr{lhs},
		Rhs: []ast.Expr{rhs},
		Tok: token.ASSIGN,
	}
}

// AssignDefine creates a short variable declaration statement (:=)
func (s *StatementBuilder) AssignDefine(lhs, rhs ast.Expr) ast.Stmt {
	return &ast.AssignStmt{
		Lhs: []ast.Expr{lhs},
		Rhs: []ast.Expr{rhs},
		Tok: token.DEFINE,
	}
}

// AssignMultiple creates a multiple assignment statement
func (s *StatementBuilder) AssignMultiple(lhs, rhs []ast.Expr) ast.Stmt {
	return &ast.AssignStmt{
		Lhs: lhs,
		Rhs: rhs,
		Tok: token.ASSIGN,
	}
}

// AssignDefineMultiple creates a multiple short variable declaration statement
func (s *StatementBuilder) AssignDefineMultiple(lhs, rhs []ast.Expr) ast.Stmt {
	return &ast.AssignStmt{
		Lhs: lhs,
		Rhs: rhs,
		Tok: token.DEFINE,
	}
}

// If creates an if statement
func (s *StatementBuilder) If(cond ast.Expr, body []ast.Stmt) ast.Stmt {
	return &ast.IfStmt{
		Cond: cond,
		Body: &ast.BlockStmt{List: body},
	}
}

// IfElse creates an if-else statement
func (s *StatementBuilder) IfElse(cond ast.Expr, ifBody, elseBody []ast.Stmt) ast.Stmt {
	return &ast.IfStmt{
		Cond: cond,
		Body: &ast.BlockStmt{List: ifBody},
		Else: &ast.BlockStmt{List: elseBody},
	}
}

// IfElseIf creates an if-else if-else statement
func (s *StatementBuilder) IfElseIf(cond ast.Expr, ifBody []ast.Stmt, elseIf *ast.IfStmt) ast.Stmt {
	return &ast.IfStmt{
		Cond: cond,
		Body: &ast.BlockStmt{List: ifBody},
		Else: elseIf,
	}
}

// For creates a for loop statement
func (s *StatementBuilder) For(init ast.Stmt, cond ast.Expr, post ast.Stmt, body []ast.Stmt) ast.Stmt {
	return &ast.ForStmt{
		Init: init,
		Cond: cond,
		Post: post,
		Body: &ast.BlockStmt{List: body},
	}
}

// ForRange creates a for-range loop statement
func (s *StatementBuilder) ForRange(key, value ast.Expr, x ast.Expr, body []ast.Stmt) ast.Stmt {
	var keyExpr, valueExpr ast.Expr
	if key != nil {
		keyExpr = key
	}
	if value != nil {
		valueExpr = value
	}

	return &ast.RangeStmt{
		Key:   keyExpr,
		Value: valueExpr,
		X:     x,
		Body:  &ast.BlockStmt{List: body},
		Tok:   token.DEFINE,
	}
}

// ForRangeAssign creates a for-range loop statement with assignment
func (s *StatementBuilder) ForRangeAssign(key, value ast.Expr, x ast.Expr, body []ast.Stmt) ast.Stmt {
	var keyExpr, valueExpr ast.Expr
	if key != nil {
		keyExpr = key
	}
	if value != nil {
		valueExpr = value
	}

	return &ast.RangeStmt{
		Key:   keyExpr,
		Value: valueExpr,
		X:     x,
		Body:  &ast.BlockStmt{List: body},
		Tok:   token.ASSIGN,
	}
}

// Switch creates a switch statement
func (s *StatementBuilder) Switch(tag ast.Expr, cases []ast.Stmt) ast.Stmt {
	return &ast.SwitchStmt{
		Tag:  tag,
		Body: &ast.BlockStmt{List: cases},
	}
}

// TypeSwitch creates a type switch statement
func (s *StatementBuilder) TypeSwitch(assign ast.Stmt, cases []ast.Stmt) ast.Stmt {
	return &ast.TypeSwitchStmt{
		Assign: assign,
		Body:   &ast.BlockStmt{List: cases},
	}
}

// Case creates a case clause
func (s *StatementBuilder) Case(values []ast.Expr, body []ast.Stmt) ast.Stmt {
	return &ast.CaseClause{
		List: values,
		Body: body,
	}
}

// Default creates a default case clause
func (s *StatementBuilder) Default(body []ast.Stmt) ast.Stmt {
	return &ast.CaseClause{
		Body: body,
	}
}

// Select creates a select statement
func (s *StatementBuilder) Select(cases []ast.Stmt) ast.Stmt {
	return &ast.SelectStmt{
		Body: &ast.BlockStmt{List: cases},
	}
}

// CommCase creates a communication case clause
func (s *StatementBuilder) CommCase(comm ast.Stmt, body []ast.Stmt) ast.Stmt {
	return &ast.CommClause{
		Comm: comm,
		Body: body,
	}
}

// CommDefault creates a default communication case clause
func (s *StatementBuilder) CommDefault(body []ast.Stmt) ast.Stmt {
	return &ast.CommClause{
		Body: body,
	}
}

// Return creates a return statement
func (s *StatementBuilder) Return(values ...ast.Expr) ast.Stmt {
	return &ast.ReturnStmt{
		Results: values,
	}
}

// ReturnEmpty creates an empty return statement
func (s *StatementBuilder) ReturnEmpty() ast.Stmt {
	return &ast.ReturnStmt{
		Results: []ast.Expr{},
	}
}

// CallStmt creates a call statement (expression statement)
func (s *StatementBuilder) CallStmt(fun ast.Expr, args ...ast.Expr) ast.Stmt {
	return &ast.ExprStmt{
		X: s.builder.call(fun, args...),
	}
}

// MethodCallStmt creates a method call statement
func (s *StatementBuilder) MethodCallStmt(receiver ast.Expr, method string, args ...ast.Expr) ast.Stmt {
	return &ast.ExprStmt{
		X: s.builder.call(s.builder.selector(receiver, method), args...),
	}
}

// Send creates a send statement (e.g., ch <- value)
func (s *StatementBuilder) Send(ch, value ast.Expr) ast.Stmt {
	return &ast.SendStmt{
		Chan:  ch,
		Value: value,
	}
}

// Go creates a go statement (e.g., go func())
func (s *StatementBuilder) Go(call ast.Expr) ast.Stmt {
	return &ast.GoStmt{
		Call: call.(*ast.CallExpr),
	}
}

// Defer creates a defer statement (e.g., defer func())
func (s *StatementBuilder) Defer(call ast.Expr) ast.Stmt {
	return &ast.DeferStmt{
		Call: call.(*ast.CallExpr),
	}
}

// Label creates a labeled statement
func (s *StatementBuilder) Label(name string, stmt ast.Stmt) ast.Stmt {
	return &ast.LabeledStmt{
		Label: s.builder.ident(name),
		Stmt:  stmt,
	}
}

// Goto creates a goto statement
func (s *StatementBuilder) Goto(label string) ast.Stmt {
	return &ast.BranchStmt{
		Tok:   token.GOTO,
		Label: s.builder.ident(label),
	}
}

// Break creates a break statement
func (s *StatementBuilder) Break(label string) ast.Stmt {
	stmt := &ast.BranchStmt{
		Tok: token.BREAK,
	}
	if label != "" {
		stmt.Label = s.builder.ident(label)
	}
	return stmt
}

// Continue creates a continue statement
func (s *StatementBuilder) Continue(label string) ast.Stmt {
	stmt := &ast.BranchStmt{
		Tok: token.CONTINUE,
	}
	if label != "" {
		stmt.Label = s.builder.ident(label)
	}
	return stmt
}

// Fallthrough creates a fallthrough statement
func (s *StatementBuilder) Fallthrough() ast.Stmt {
	return &ast.BranchStmt{
		Tok: token.FALLTHROUGH,
	}
}

// Block creates a block statement
func (s *StatementBuilder) Block(stmts []ast.Stmt) ast.Stmt {
	return &ast.BlockStmt{
		List: stmts,
	}
}

// Empty creates an empty statement
func (s *StatementBuilder) Empty() ast.Stmt {
	return &ast.EmptyStmt{}
}

// Inc creates an increment statement (e.g., x++)
func (s *StatementBuilder) Inc(expr ast.Expr) ast.Stmt {
	return &ast.IncDecStmt{
		X:   expr,
		Tok: token.INC,
	}
}

// Dec creates a decrement statement (e.g., x--)
func (s *StatementBuilder) Dec(expr ast.Expr) ast.Stmt {
	return &ast.IncDecStmt{
		X:   expr,
		Tok: token.DEC,
	}
}

// AddAssign creates an add-assign statement (e.g., x += y)
func (s *StatementBuilder) AddAssign(lhs, rhs ast.Expr) ast.Stmt {
	return &ast.AssignStmt{
		Lhs: []ast.Expr{lhs},
		Rhs: []ast.Expr{rhs},
		Tok: token.ADD_ASSIGN,
	}
}

// SubAssign creates a subtract-assign statement (e.g., x -= y)
func (s *StatementBuilder) SubAssign(lhs, rhs ast.Expr) ast.Stmt {
	return &ast.AssignStmt{
		Lhs: []ast.Expr{lhs},
		Rhs: []ast.Expr{rhs},
		Tok: token.SUB_ASSIGN,
	}
}

// MulAssign creates a multiply-assign statement (e.g., x *= y)
func (s *StatementBuilder) MulAssign(lhs, rhs ast.Expr) ast.Stmt {
	return &ast.AssignStmt{
		Lhs: []ast.Expr{lhs},
		Rhs: []ast.Expr{rhs},
		Tok: token.MUL_ASSIGN,
	}
}

// DivAssign creates a divide-assign statement (e.g., x /= y)
func (s *StatementBuilder) DivAssign(lhs, rhs ast.Expr) ast.Stmt {
	return &ast.AssignStmt{
		Lhs: []ast.Expr{lhs},
		Rhs: []ast.Expr{rhs},
		Tok: token.QUO_ASSIGN,
	}
}

// ModAssign creates a modulo-assign statement (e.g., x %= y)
func (s *StatementBuilder) ModAssign(lhs, rhs ast.Expr) ast.Stmt {
	return &ast.AssignStmt{
		Lhs: []ast.Expr{lhs},
		Rhs: []ast.Expr{rhs},
		Tok: token.REM_ASSIGN,
	}
}

// AndAssign creates a bitwise AND-assign statement (e.g., x &= y)
func (s *StatementBuilder) AndAssign(lhs, rhs ast.Expr) ast.Stmt {
	return &ast.AssignStmt{
		Lhs: []ast.Expr{lhs},
		Rhs: []ast.Expr{rhs},
		Tok: token.AND_ASSIGN,
	}
}

// OrAssign creates a bitwise OR-assign statement (e.g., x |= y)
func (s *StatementBuilder) OrAssign(lhs, rhs ast.Expr) ast.Stmt {
	return &ast.AssignStmt{
		Lhs: []ast.Expr{lhs},
		Rhs: []ast.Expr{rhs},
		Tok: token.OR_ASSIGN,
	}
}

// XorAssign creates a bitwise XOR-assign statement (e.g., x ^= y)
func (s *StatementBuilder) XorAssign(lhs, rhs ast.Expr) ast.Stmt {
	return &ast.AssignStmt{
		Lhs: []ast.Expr{lhs},
		Rhs: []ast.Expr{rhs},
		Tok: token.XOR_ASSIGN,
	}
}

// ShlAssign creates a left shift-assign statement (e.g., x <<= y)
func (s *StatementBuilder) ShlAssign(lhs, rhs ast.Expr) ast.Stmt {
	return &ast.AssignStmt{
		Lhs: []ast.Expr{lhs},
		Rhs: []ast.Expr{rhs},
		Tok: token.SHL_ASSIGN,
	}
}

// ShrAssign creates a right shift-assign statement (e.g., x >>= y)
func (s *StatementBuilder) ShrAssign(lhs, rhs ast.Expr) ast.Stmt {
	return &ast.AssignStmt{
		Lhs: []ast.Expr{lhs},
		Rhs: []ast.Expr{rhs},
		Tok: token.SHR_ASSIGN,
	}
}

// AndNotAssign creates a bitwise AND NOT-assign statement (e.g., x &^= y)
func (s *StatementBuilder) AndNotAssign(lhs, rhs ast.Expr) ast.Stmt {
	return &ast.AssignStmt{
		Lhs: []ast.Expr{lhs},
		Rhs: []ast.Expr{rhs},
		Tok: token.AND_NOT_ASSIGN,
	}
}

// Helper method to get the underlying builder
func (s *StatementBuilder) Builder() *Builder {
	return s.builder
}
