package astbuilder

import (
	"go/ast"
	"go/token"
	"testing"
)

func TestStatementBuilder_DeclareVar(t *testing.T) {
	builder := NewBuilder(BuilderConfig{PackageName: "test"})
	stmtBuilder := NewStatementBuilder(builder)

	stmt := stmtBuilder.DeclareVar("test", "string", nil)

	if stmt == nil {
		t.Fatal("DeclareVar returned nil")
	}

	if declStmt, ok := stmt.(*ast.DeclStmt); !ok {
		t.Fatal("DeclareVar should return *ast.DeclStmt")
	} else if genDecl, ok := declStmt.Decl.(*ast.GenDecl); !ok {
		t.Fatal("DeclareVar should contain *ast.GenDecl")
	} else if genDecl.Tok != token.VAR {
		t.Error("DeclareVar should use VAR token")
	} else if len(genDecl.Specs) != 1 {
		t.Errorf("Expected 1 spec, got %d", len(genDecl.Specs))
	}
}

func TestStatementBuilder_DeclareVarWithType(t *testing.T) {
	builder := NewBuilder(BuilderConfig{PackageName: "test"})
	stmtBuilder := NewStatementBuilder(builder)
	exprBuilder := NewExpressionBuilder(builder)

	typeExpr := exprBuilder.Ident("string")
	stmt := stmtBuilder.DeclareVarWithType("test", typeExpr, nil)

	if stmt == nil {
		t.Fatal("DeclareVarWithType returned nil")
	}

	if _, ok := stmt.(*ast.DeclStmt); !ok {
		t.Fatal("DeclareVarWithType should return *ast.DeclStmt")
	}
}

func TestStatementBuilder_Assign(t *testing.T) {
	builder := NewBuilder(BuilderConfig{PackageName: "test"})
	stmtBuilder := NewStatementBuilder(builder)
	exprBuilder := NewExpressionBuilder(builder)

	lhs := exprBuilder.Ident("a")
	rhs := exprBuilder.Ident("b")
	stmt := stmtBuilder.Assign(lhs, rhs)

	if stmt == nil {
		t.Fatal("Assign returned nil")
	}

	if assign, ok := stmt.(*ast.AssignStmt); !ok {
		t.Fatal("Assign should return *ast.AssignStmt")
	} else if assign.Tok != token.ASSIGN {
		t.Error("Assign should use ASSIGN token")
	} else if len(assign.Lhs) != 1 {
		t.Errorf("Expected 1 LHS, got %d", len(assign.Lhs))
	} else if len(assign.Rhs) != 1 {
		t.Errorf("Expected 1 RHS, got %d", len(assign.Rhs))
	}
}

func TestStatementBuilder_AssignDefine(t *testing.T) {
	builder := NewBuilder(BuilderConfig{PackageName: "test"})
	stmtBuilder := NewStatementBuilder(builder)
	exprBuilder := NewExpressionBuilder(builder)

	lhs := exprBuilder.Ident("a")
	rhs := exprBuilder.Ident("b")
	stmt := stmtBuilder.AssignDefine(lhs, rhs)

	if stmt == nil {
		t.Fatal("AssignDefine returned nil")
	}

	if assign, ok := stmt.(*ast.AssignStmt); !ok {
		t.Fatal("AssignDefine should return *ast.AssignStmt")
	} else if assign.Tok != token.DEFINE {
		t.Error("AssignDefine should use DEFINE token")
	}
}

func TestStatementBuilder_If(t *testing.T) {
	builder := NewBuilder(BuilderConfig{PackageName: "test"})
	stmtBuilder := NewStatementBuilder(builder)
	exprBuilder := NewExpressionBuilder(builder)

	cond := exprBuilder.Ident("test")
	body := []ast.Stmt{stmtBuilder.ReturnEmpty()}
	stmt := stmtBuilder.If(cond, body)

	if stmt == nil {
		t.Fatal("If returned nil")
	}

	if ifStmt, ok := stmt.(*ast.IfStmt); !ok {
		t.Fatal("If should return *ast.IfStmt")
	} else if ifStmt.Cond != cond {
		t.Error("If should preserve the condition")
	} else if len(ifStmt.Body.List) != 1 {
		t.Errorf("Expected 1 statement in body, got %d", len(ifStmt.Body.List))
	}
}

func TestStatementBuilder_IfElse(t *testing.T) {
	builder := NewBuilder(BuilderConfig{PackageName: "test"})
	stmtBuilder := NewStatementBuilder(builder)
	exprBuilder := NewExpressionBuilder(builder)

	cond := exprBuilder.Ident("test")
	ifBody := []ast.Stmt{stmtBuilder.ReturnEmpty()}
	elseBody := []ast.Stmt{stmtBuilder.ReturnEmpty()}
	stmt := stmtBuilder.IfElse(cond, ifBody, elseBody)

	if stmt == nil {
		t.Fatal("IfElse returned nil")
	}

	if ifStmt, ok := stmt.(*ast.IfStmt); !ok {
		t.Fatal("IfElse should return *ast.IfStmt")
	} else if ifStmt.Cond != cond {
		t.Error("IfElse should preserve the condition")
	} else if len(ifStmt.Body.List) != 1 {
		t.Errorf("Expected 1 statement in if body, got %d", len(ifStmt.Body.List))
	} else if ifStmt.Else == nil {
		t.Error("IfElse should have an else clause")
	}
}

func TestStatementBuilder_For(t *testing.T) {
	builder := NewBuilder(BuilderConfig{PackageName: "test"})
	stmtBuilder := NewStatementBuilder(builder)
	exprBuilder := NewExpressionBuilder(builder)

	init := stmtBuilder.DeclareVar("i", "int", exprBuilder.Int(0))
	cond := exprBuilder.Less(exprBuilder.Ident("i"), exprBuilder.Int(10))
	post := stmtBuilder.Inc(exprBuilder.Ident("i"))
	body := []ast.Stmt{stmtBuilder.ReturnEmpty()}
	stmt := stmtBuilder.For(init, cond, post, body)

	if stmt == nil {
		t.Fatal("For returned nil")
	}

	if forStmt, ok := stmt.(*ast.ForStmt); !ok {
		t.Fatal("For should return *ast.ForStmt")
	} else if forStmt.Init != init {
		t.Error("For should preserve the init statement")
	} else if forStmt.Cond != cond {
		t.Error("For should preserve the condition")
	} else if forStmt.Post != post {
		t.Error("For should preserve the post statement")
	} else if len(forStmt.Body.List) != 1 {
		t.Errorf("Expected 1 statement in body, got %d", len(forStmt.Body.List))
	}
}

func TestStatementBuilder_ForRange(t *testing.T) {
	builder := NewBuilder(BuilderConfig{PackageName: "test"})
	stmtBuilder := NewStatementBuilder(builder)
	exprBuilder := NewExpressionBuilder(builder)

	key := exprBuilder.Ident("k")
	value := exprBuilder.Ident("v")
	x := exprBuilder.Ident("m")
	body := []ast.Stmt{stmtBuilder.ReturnEmpty()}
	stmt := stmtBuilder.ForRange(key, value, x, body)

	if stmt == nil {
		t.Fatal("ForRange returned nil")
	}

	if rangeStmt, ok := stmt.(*ast.RangeStmt); !ok {
		t.Fatal("ForRange should return *ast.RangeStmt")
	} else if rangeStmt.Key != key {
		t.Error("ForRange should preserve the key")
	} else if rangeStmt.Value != value {
		t.Error("ForRange should preserve the value")
	} else if rangeStmt.X != x {
		t.Error("ForRange should preserve the range expression")
	} else if rangeStmt.Tok != token.DEFINE {
		t.Error("ForRange should use DEFINE token")
	} else if len(rangeStmt.Body.List) != 1 {
		t.Errorf("Expected 1 statement in body, got %d", len(rangeStmt.Body.List))
	}
}

func TestStatementBuilder_Return(t *testing.T) {
	builder := NewBuilder(BuilderConfig{PackageName: "test"})
	stmtBuilder := NewStatementBuilder(builder)
	exprBuilder := NewExpressionBuilder(builder)

	expr := exprBuilder.Ident("test")
	stmt := stmtBuilder.Return(expr)

	if stmt == nil {
		t.Fatal("Return returned nil")
	}

	if ret, ok := stmt.(*ast.ReturnStmt); !ok {
		t.Fatal("Return should return *ast.ReturnStmt")
	} else if len(ret.Results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(ret.Results))
	} else if ret.Results[0] != expr {
		t.Error("Return should preserve the expression")
	}
}

func TestStatementBuilder_ReturnEmpty(t *testing.T) {
	builder := NewBuilder(BuilderConfig{PackageName: "test"})
	stmtBuilder := NewStatementBuilder(builder)

	stmt := stmtBuilder.ReturnEmpty()

	if stmt == nil {
		t.Fatal("ReturnEmpty returned nil")
	}

	if ret, ok := stmt.(*ast.ReturnStmt); !ok {
		t.Fatal("ReturnEmpty should return *ast.ReturnStmt")
	} else if len(ret.Results) != 0 {
		t.Errorf("Expected 0 results, got %d", len(ret.Results))
	}
}

func TestStatementBuilder_CallStmt(t *testing.T) {
	builder := NewBuilder(BuilderConfig{PackageName: "test"})
	stmtBuilder := NewStatementBuilder(builder)
	exprBuilder := NewExpressionBuilder(builder)

	fun := exprBuilder.Ident("test")
	arg := exprBuilder.String("arg")
	stmt := stmtBuilder.CallStmt(fun, arg)

	if stmt == nil {
		t.Fatal("CallStmt returned nil")
	}

	if exprStmt, ok := stmt.(*ast.ExprStmt); !ok {
		t.Fatal("CallStmt should return *ast.ExprStmt")
	} else if call, ok := exprStmt.X.(*ast.CallExpr); !ok {
		t.Fatal("CallStmt should contain *ast.CallExpr")
	} else if call.Fun != fun {
		t.Error("CallStmt should preserve the function")
	} else if len(call.Args) != 1 {
		t.Errorf("Expected 1 argument, got %d", len(call.Args))
	}
}

func TestStatementBuilder_MethodCallStmt(t *testing.T) {
	builder := NewBuilder(BuilderConfig{PackageName: "test"})
	stmtBuilder := NewStatementBuilder(builder)
	exprBuilder := NewExpressionBuilder(builder)

	receiver := exprBuilder.Ident("obj")
	arg := exprBuilder.String("arg")
	stmt := stmtBuilder.MethodCallStmt(receiver, "method", arg)

	if stmt == nil {
		t.Fatal("MethodCallStmt returned nil")
	}

	if exprStmt, ok := stmt.(*ast.ExprStmt); !ok {
		t.Fatal("MethodCallStmt should return *ast.ExprStmt")
	} else if call, ok := exprStmt.X.(*ast.CallExpr); !ok {
		t.Fatal("MethodCallStmt should contain *ast.CallExpr")
	} else if len(call.Args) != 1 {
		t.Errorf("Expected 1 argument, got %d", len(call.Args))
	}
}

func TestStatementBuilder_Inc(t *testing.T) {
	builder := NewBuilder(BuilderConfig{PackageName: "test"})
	stmtBuilder := NewStatementBuilder(builder)
	exprBuilder := NewExpressionBuilder(builder)

	expr := exprBuilder.Ident("i")
	stmt := stmtBuilder.Inc(expr)

	if stmt == nil {
		t.Fatal("Inc returned nil")
	}

	if incDec, ok := stmt.(*ast.IncDecStmt); !ok {
		t.Fatal("Inc should return *ast.IncDecStmt")
	} else if incDec.X != expr {
		t.Error("Inc should preserve the expression")
	} else if incDec.Tok != token.INC {
		t.Error("Inc should use INC token")
	}
}

func TestStatementBuilder_Dec(t *testing.T) {
	builder := NewBuilder(BuilderConfig{PackageName: "test"})
	stmtBuilder := NewStatementBuilder(builder)
	exprBuilder := NewExpressionBuilder(builder)

	expr := exprBuilder.Ident("i")
	stmt := stmtBuilder.Dec(expr)

	if stmt == nil {
		t.Fatal("Dec returned nil")
	}

	if incDec, ok := stmt.(*ast.IncDecStmt); !ok {
		t.Fatal("Dec should return *ast.IncDecStmt")
	} else if incDec.X != expr {
		t.Error("Dec should preserve the expression")
	} else if incDec.Tok != token.DEC {
		t.Error("Dec should use DEC token")
	}
}

func TestStatementBuilder_Block(t *testing.T) {
	builder := NewBuilder(BuilderConfig{PackageName: "test"})
	stmtBuilder := NewStatementBuilder(builder)

	stmts := []ast.Stmt{
		stmtBuilder.ReturnEmpty(),
		stmtBuilder.ReturnEmpty(),
	}
	stmt := stmtBuilder.Block(stmts)

	if stmt == nil {
		t.Fatal("Block returned nil")
	}

	if block, ok := stmt.(*ast.BlockStmt); !ok {
		t.Fatal("Block should return *ast.BlockStmt")
	} else if len(block.List) != 2 {
		t.Errorf("Expected 2 statements, got %d", len(block.List))
	}
}

func TestStatementBuilder_Empty(t *testing.T) {
	builder := NewBuilder(BuilderConfig{PackageName: "test"})
	stmtBuilder := NewStatementBuilder(builder)

	stmt := stmtBuilder.Empty()

	if stmt == nil {
		t.Fatal("Empty returned nil")
	}

	if _, ok := stmt.(*ast.EmptyStmt); !ok {
		t.Fatal("Empty should return *ast.EmptyStmt")
	}
}
