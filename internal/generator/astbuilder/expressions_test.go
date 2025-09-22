package astbuilder

import (
	"go/ast"
	"go/token"
	"testing"
)

func TestExpressionBuilder_Ident(t *testing.T) {
	builder := NewBuilder(BuilderConfig{PackageName: "test"})
	exprBuilder := NewExpressionBuilder(builder)

	expr := exprBuilder.Ident("test")

	if expr == nil {
		t.Fatal("Ident returned nil")
	}

	if ident, ok := expr.(*ast.Ident); !ok {
		t.Fatal("Ident should return *ast.Ident")
	} else if ident.Name != "test" {
		t.Errorf("Expected identifier name 'test', got '%s'", ident.Name)
	}
}

func TestExpressionBuilder_String(t *testing.T) {
	builder := NewBuilder(BuilderConfig{PackageName: "test"})
	exprBuilder := NewExpressionBuilder(builder)

	expr := exprBuilder.String("test")

	if expr == nil {
		t.Fatal("String returned nil")
	}

	if lit, ok := expr.(*ast.BasicLit); !ok {
		t.Fatal("String should return *ast.BasicLit")
	} else if lit.Kind != token.STRING {
		t.Error("String should return a string literal")
	} else if lit.Value != `"test"` {
		t.Errorf("Expected string value '\"test\"', got '%s'", lit.Value)
	}
}

func TestExpressionBuilder_Int(t *testing.T) {
	builder := NewBuilder(BuilderConfig{PackageName: "test"})
	exprBuilder := NewExpressionBuilder(builder)

	expr := exprBuilder.Int(42)

	if expr == nil {
		t.Fatal("Int returned nil")
	}

	if lit, ok := expr.(*ast.BasicLit); !ok {
		t.Fatal("Int should return *ast.BasicLit")
	} else if lit.Kind != token.INT {
		t.Error("Int should return an integer literal")
	} else if lit.Value != "42" {
		t.Errorf("Expected integer value '42', got '%s'", lit.Value)
	}
}

func TestExpressionBuilder_Bool(t *testing.T) {
	builder := NewBuilder(BuilderConfig{PackageName: "test"})
	exprBuilder := NewExpressionBuilder(builder)

	expr := exprBuilder.Bool(true)

	if expr == nil {
		t.Fatal("Bool returned nil")
	}

	if ident, ok := expr.(*ast.Ident); !ok {
		t.Fatal("Bool should return *ast.Ident")
	} else if ident.Name != "true" {
		t.Errorf("Expected boolean value 'true', got '%s'", ident.Name)
	}
}

func TestExpressionBuilder_Nil(t *testing.T) {
	builder := NewBuilder(BuilderConfig{PackageName: "test"})
	exprBuilder := NewExpressionBuilder(builder)

	expr := exprBuilder.Nil()

	if expr == nil {
		t.Fatal("Nil returned nil")
	}

	if ident, ok := expr.(*ast.Ident); !ok {
		t.Fatal("Nil should return *ast.Ident")
	} else if ident.Name != "nil" {
		t.Errorf("Expected nil value 'nil', got '%s'", ident.Name)
	}
}

func TestExpressionBuilder_Select(t *testing.T) {
	builder := NewBuilder(BuilderConfig{PackageName: "test"})
	exprBuilder := NewExpressionBuilder(builder)

	receiver := exprBuilder.Ident("obj")
	expr := exprBuilder.Select(receiver, "field")

	if expr == nil {
		t.Fatal("Select returned nil")
	}

	if sel, ok := expr.(*ast.SelectorExpr); !ok {
		t.Fatal("Select should return *ast.SelectorExpr")
	} else if sel.Sel.Name != "field" {
		t.Errorf("Expected selector name 'field', got '%s'", sel.Sel.Name)
	}
}

func TestExpressionBuilder_Call(t *testing.T) {
	builder := NewBuilder(BuilderConfig{PackageName: "test"})
	exprBuilder := NewExpressionBuilder(builder)

	fun := exprBuilder.Ident("test")
	expr := exprBuilder.Call(fun, exprBuilder.String("arg"))

	if expr == nil {
		t.Fatal("Call returned nil")
	}

	if call, ok := expr.(*ast.CallExpr); !ok {
		t.Fatal("Call should return *ast.CallExpr")
	} else if len(call.Args) != 1 {
		t.Errorf("Expected 1 argument, got %d", len(call.Args))
	}
}

func TestExpressionBuilder_MethodCall(t *testing.T) {
	builder := NewBuilder(BuilderConfig{PackageName: "test"})
	exprBuilder := NewExpressionBuilder(builder)

	receiver := exprBuilder.Ident("obj")
	expr := exprBuilder.MethodCall(receiver, "method", exprBuilder.String("arg"))

	if expr == nil {
		t.Fatal("MethodCall returned nil")
	}

	if call, ok := expr.(*ast.CallExpr); !ok {
		t.Fatal("MethodCall should return *ast.CallExpr")
	} else if len(call.Args) != 1 {
		t.Errorf("Expected 1 argument, got %d", len(call.Args))
	}
}

func TestExpressionBuilder_Equal(t *testing.T) {
	builder := NewBuilder(BuilderConfig{PackageName: "test"})
	exprBuilder := NewExpressionBuilder(builder)

	left := exprBuilder.Ident("a")
	right := exprBuilder.Ident("b")
	expr := exprBuilder.Equal(left, right)

	if expr == nil {
		t.Fatal("Equal returned nil")
	}

	if binary, ok := expr.(*ast.BinaryExpr); !ok {
		t.Fatal("Equal should return *ast.BinaryExpr")
	} else if binary.Op != token.EQL {
		t.Error("Equal should use EQL token")
	}
}

func TestExpressionBuilder_NotEqual(t *testing.T) {
	builder := NewBuilder(BuilderConfig{PackageName: "test"})
	exprBuilder := NewExpressionBuilder(builder)

	left := exprBuilder.Ident("a")
	right := exprBuilder.Ident("b")
	expr := exprBuilder.NotEqual(left, right)

	if expr == nil {
		t.Fatal("NotEqual returned nil")
	}

	if binary, ok := expr.(*ast.BinaryExpr); !ok {
		t.Fatal("NotEqual should return *ast.BinaryExpr")
	} else if binary.Op != token.NEQ {
		t.Error("NotEqual should use NEQ token")
	}
}

func TestExpressionBuilder_AddressOf(t *testing.T) {
	builder := NewBuilder(BuilderConfig{PackageName: "test"})
	exprBuilder := NewExpressionBuilder(builder)

	expr := exprBuilder.Ident("test")
	addr := exprBuilder.AddressOf(expr)

	if addr == nil {
		t.Fatal("AddressOf returned nil")
	}

	if unary, ok := addr.(*ast.UnaryExpr); !ok {
		t.Fatal("AddressOf should return *ast.UnaryExpr")
	} else if unary.Op != token.AND {
		t.Error("AddressOf should use AND token")
	}
}

func TestExpressionBuilder_Deref(t *testing.T) {
	builder := NewBuilder(BuilderConfig{PackageName: "test"})
	exprBuilder := NewExpressionBuilder(builder)

	expr := exprBuilder.Ident("test")
	deref := exprBuilder.Deref(expr)

	if deref == nil {
		t.Fatal("Deref returned nil")
	}

	if unary, ok := deref.(*ast.UnaryExpr); !ok {
		t.Fatal("Deref should return *ast.UnaryExpr")
	} else if unary.Op != token.MUL {
		t.Error("Deref should use MUL token")
	}
}

func TestExpressionBuilder_CompositeLit(t *testing.T) {
	builder := NewBuilder(BuilderConfig{PackageName: "test"})
	exprBuilder := NewExpressionBuilder(builder)

	expr := exprBuilder.CompositeLit("Test", exprBuilder.String("value"))

	if expr == nil {
		t.Fatal("CompositeLit returned nil")
	}

	if comp, ok := expr.(*ast.CompositeLit); !ok {
		t.Fatal("CompositeLit should return *ast.CompositeLit")
	} else if comp.Type.(*ast.Ident).Name != "Test" {
		t.Errorf("Expected type name 'Test', got '%s'", comp.Type.(*ast.Ident).Name)
	} else if len(comp.Elts) != 1 {
		t.Errorf("Expected 1 element, got %d", len(comp.Elts))
	}
}

func TestExpressionBuilder_KeyValue(t *testing.T) {
	builder := NewBuilder(BuilderConfig{PackageName: "test"})
	exprBuilder := NewExpressionBuilder(builder)

	key := exprBuilder.String("key")
	value := exprBuilder.String("value")
	expr := exprBuilder.KeyValue(key, value)

	if expr == nil {
		t.Fatal("KeyValue returned nil")
	}

	if kv, ok := expr.(*ast.KeyValueExpr); !ok {
		t.Fatal("KeyValue should return *ast.KeyValueExpr")
	} else if kv.Key != key {
		t.Error("KeyValue should preserve the key")
	} else if kv.Value != value {
		t.Error("KeyValue should preserve the value")
	}
}

func TestExpressionBuilder_SliceType(t *testing.T) {
	builder := NewBuilder(BuilderConfig{PackageName: "test"})
	exprBuilder := NewExpressionBuilder(builder)

	elementType := exprBuilder.Ident("string")
	expr := exprBuilder.SliceType(elementType)

	if expr == nil {
		t.Fatal("SliceType returned nil")
	}

	if array, ok := expr.(*ast.ArrayType); !ok {
		t.Fatal("SliceType should return *ast.ArrayType")
	} else if array.Len != nil {
		t.Error("SliceType should not have a length")
	} else if array.Elt != elementType {
		t.Error("SliceType should preserve the element type")
	}
}

func TestExpressionBuilder_MapType(t *testing.T) {
	builder := NewBuilder(BuilderConfig{PackageName: "test"})
	exprBuilder := NewExpressionBuilder(builder)

	keyType := exprBuilder.Ident("string")
	valueType := exprBuilder.Ident("int")
	expr := exprBuilder.MapType(keyType, valueType)

	if expr == nil {
		t.Fatal("MapType returned nil")
	}

	if mapType, ok := expr.(*ast.MapType); !ok {
		t.Fatal("MapType should return *ast.MapType")
	} else if mapType.Key != keyType {
		t.Error("MapType should preserve the key type")
	} else if mapType.Value != valueType {
		t.Error("MapType should preserve the value type")
	}
}

func TestExpressionBuilder_ChanType(t *testing.T) {
	builder := NewBuilder(BuilderConfig{PackageName: "test"})
	exprBuilder := NewExpressionBuilder(builder)

	valueType := exprBuilder.Ident("int")
	expr := exprBuilder.Chan(valueType, ast.SEND)

	if expr == nil {
		t.Fatal("Chan returned nil")
	}

	if chanType, ok := expr.(*ast.ChanType); !ok {
		t.Fatal("Chan should return *ast.ChanType")
	} else if chanType.Value != valueType {
		t.Error("Chan should preserve the value type")
	} else if chanType.Dir != ast.SEND {
		t.Error("Chan should preserve the direction")
	}
}

func TestExpressionBuilder_Star(t *testing.T) {
	builder := NewBuilder(BuilderConfig{PackageName: "test"})
	exprBuilder := NewExpressionBuilder(builder)

	typeExpr := exprBuilder.Ident("int")
	expr := exprBuilder.Star(typeExpr)

	if expr == nil {
		t.Fatal("Star returned nil")
	}

	if star, ok := expr.(*ast.StarExpr); !ok {
		t.Fatal("Star should return *ast.StarExpr")
	} else if star.X != typeExpr {
		t.Error("Star should preserve the type expression")
	}
}

func TestExpressionBuilder_Paren(t *testing.T) {
	builder := NewBuilder(BuilderConfig{PackageName: "test"})
	exprBuilder := NewExpressionBuilder(builder)

	inner := exprBuilder.Ident("test")
	expr := exprBuilder.Paren(inner)

	if expr == nil {
		t.Fatal("Paren returned nil")
	}

	if paren, ok := expr.(*ast.ParenExpr); !ok {
		t.Fatal("Paren should return *ast.ParenExpr")
	} else if paren.X != inner {
		t.Error("Paren should preserve the inner expression")
	}
}
