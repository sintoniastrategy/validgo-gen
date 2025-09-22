package astbuilder

import (
	"go/ast"
	"go/token"
	"testing"
)

func TestNewBuilder(t *testing.T) {
	config := BuilderConfig{
		PackageName:  "test",
		ImportPrefix: "github.com/test",
		UsePointers:  true,
	}

	builder := NewBuilder(config)

	if builder == nil {
		t.Fatal("NewBuilder returned nil")
	}

	if builder.config.PackageName != "test" {
		t.Errorf("Expected package name 'test', got '%s'", builder.config.PackageName)
	}

	if builder.config.ImportPrefix != "github.com/test" {
		t.Errorf("Expected import prefix 'github.com/test', got '%s'", builder.config.ImportPrefix)
	}

	if !builder.config.UsePointers {
		t.Error("Expected UsePointers to be true")
	}

	if builder.imports == nil {
		t.Error("Expected imports map to be initialized")
	}

	if builder.stmts == nil {
		t.Error("Expected statements slice to be initialized")
	}

	if builder.decls == nil {
		t.Error("Expected declarations slice to be initialized")
	}
}

func TestBuilder_AddImport(t *testing.T) {
	builder := NewBuilder(BuilderConfig{PackageName: "test"})

	builder.AddImport("fmt")
	builder.AddImport("net/http")
	builder.AddImport("fmt") // Duplicate

	if builder.ImportCount() != 2 {
		t.Errorf("Expected 2 imports, got %d", builder.ImportCount())
	}

	if !builder.HasImports() {
		t.Error("Expected builder to have imports")
	}
}

func TestBuilder_AddStatement(t *testing.T) {
	builder := NewBuilder(BuilderConfig{PackageName: "test"})

	stmt := &ast.ExprStmt{
		X: &ast.Ident{Name: "test"},
	}

	builder.AddStatement(stmt)

	if builder.StatementCount() != 1 {
		t.Errorf("Expected 1 statement, got %d", builder.StatementCount())
	}

	if !builder.HasStatements() {
		t.Error("Expected builder to have statements")
	}
}

func TestBuilder_AddDeclaration(t *testing.T) {
	builder := NewBuilder(BuilderConfig{PackageName: "test"})

	decl := &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: &ast.Ident{Name: "Test"},
				Type: &ast.Ident{Name: "string"},
			},
		},
	}

	builder.AddDeclaration(decl)

	if builder.DeclarationCount() != 1 {
		t.Errorf("Expected 1 declaration, got %d", builder.DeclarationCount())
	}

	if !builder.HasDeclarations() {
		t.Error("Expected builder to have declarations")
	}
}

func TestBuilder_AddStatements(t *testing.T) {
	builder := NewBuilder(BuilderConfig{PackageName: "test"})

	stmts := []ast.Stmt{
		&ast.ExprStmt{X: &ast.Ident{Name: "stmt1"}},
		&ast.ExprStmt{X: &ast.Ident{Name: "stmt2"}},
	}

	builder.AddStatements(stmts)

	if builder.StatementCount() != 2 {
		t.Errorf("Expected 2 statements, got %d", builder.StatementCount())
	}
}

func TestBuilder_AddDeclarations(t *testing.T) {
	builder := NewBuilder(BuilderConfig{PackageName: "test"})

	decls := []ast.Decl{
		&ast.GenDecl{
			Tok: token.TYPE,
			Specs: []ast.Spec{
				&ast.TypeSpec{
					Name: &ast.Ident{Name: "Type1"},
					Type: &ast.Ident{Name: "string"},
				},
			},
		},
		&ast.GenDecl{
			Tok: token.TYPE,
			Specs: []ast.Spec{
				&ast.TypeSpec{
					Name: &ast.Ident{Name: "Type2"},
					Type: &ast.Ident{Name: "int"},
				},
			},
		},
	}

	builder.AddDeclarations(decls)

	if builder.DeclarationCount() != 2 {
		t.Errorf("Expected 2 declarations, got %d", builder.DeclarationCount())
	}
}

func TestBuilder_Build(t *testing.T) {
	builder := NewBuilder(BuilderConfig{PackageName: "test"})

	builder.AddImport("fmt")
	builder.AddImport("net/http")

	stmt := &ast.ExprStmt{X: &ast.Ident{Name: "test"}}
	builder.AddStatement(stmt)

	decl := &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: &ast.Ident{Name: "Test"},
				Type: &ast.Ident{Name: "string"},
			},
		},
	}
	builder.AddDeclaration(decl)

	stmts, decls, imports := builder.Build()

	if len(stmts) != 1 {
		t.Errorf("Expected 1 statement, got %d", len(stmts))
	}

	if len(decls) != 1 {
		t.Errorf("Expected 1 declaration, got %d", len(decls))
	}

	if len(imports) != 2 {
		t.Errorf("Expected 2 imports, got %d", len(imports))
	}

	// Check import order
	if imports[0] != "fmt" {
		t.Errorf("Expected first import to be 'fmt', got '%s'", imports[0])
	}

	if imports[1] != "net/http" {
		t.Errorf("Expected second import to be 'net/http', got '%s'", imports[1])
	}
}

func TestBuilder_Clear(t *testing.T) {
	builder := NewBuilder(BuilderConfig{PackageName: "test"})

	builder.AddImport("fmt")
	builder.AddStatement(&ast.ExprStmt{X: &ast.Ident{Name: "test"}})
	builder.AddDeclaration(&ast.GenDecl{Tok: token.TYPE})

	builder.Clear()

	if builder.HasImports() {
		t.Error("Expected builder to not have imports after clear")
	}

	if builder.HasStatements() {
		t.Error("Expected builder to not have statements after clear")
	}

	if builder.HasDeclarations() {
		t.Error("Expected builder to not have declarations after clear")
	}
}

func TestBuilder_ClearStatements(t *testing.T) {
	builder := NewBuilder(BuilderConfig{PackageName: "test"})

	builder.AddImport("fmt")
	builder.AddStatement(&ast.ExprStmt{X: &ast.Ident{Name: "test"}})
	builder.AddDeclaration(&ast.GenDecl{Tok: token.TYPE})

	builder.ClearStatements()

	if !builder.HasImports() {
		t.Error("Expected builder to still have imports after clear statements")
	}

	if builder.HasStatements() {
		t.Error("Expected builder to not have statements after clear statements")
	}

	if !builder.HasDeclarations() {
		t.Error("Expected builder to still have declarations after clear statements")
	}
}

func TestBuilder_ClearDeclarations(t *testing.T) {
	builder := NewBuilder(BuilderConfig{PackageName: "test"})

	builder.AddImport("fmt")
	builder.AddStatement(&ast.ExprStmt{X: &ast.Ident{Name: "test"}})
	builder.AddDeclaration(&ast.GenDecl{Tok: token.TYPE})

	builder.ClearDeclarations()

	if !builder.HasImports() {
		t.Error("Expected builder to still have imports after clear declarations")
	}

	if !builder.HasStatements() {
		t.Error("Expected builder to still have statements after clear declarations")
	}

	if builder.HasDeclarations() {
		t.Error("Expected builder to not have declarations after clear declarations")
	}
}

func TestBuilder_Clone(t *testing.T) {
	builder := NewBuilder(BuilderConfig{PackageName: "test"})

	builder.AddImport("fmt")
	builder.AddStatement(&ast.ExprStmt{X: &ast.Ident{Name: "test"}})
	builder.AddDeclaration(&ast.GenDecl{Tok: token.TYPE})

	clone := builder.Clone()

	if clone == nil {
		t.Fatal("Clone returned nil")
	}

	if clone.config.PackageName != builder.config.PackageName {
		t.Error("Clone should have same package name")
	}

	if clone.ImportCount() != builder.ImportCount() {
		t.Error("Clone should have same number of imports")
	}

	if clone.StatementCount() != builder.StatementCount() {
		t.Error("Clone should have same number of statements")
	}

	if clone.DeclarationCount() != builder.DeclarationCount() {
		t.Error("Clone should have same number of declarations")
	}

	// Modify original
	builder.AddImport("net/http")

	// Clone should not be affected
	if clone.ImportCount() != 1 {
		t.Error("Clone should not be affected by changes to original")
	}
}

func TestBuilder_Merge(t *testing.T) {
	builder1 := NewBuilder(BuilderConfig{PackageName: "test1"})
	builder1.AddImport("fmt")
	builder1.AddStatement(&ast.ExprStmt{X: &ast.Ident{Name: "stmt1"}})

	builder2 := NewBuilder(BuilderConfig{PackageName: "test2"})
	builder2.AddImport("net/http")
	builder2.AddStatement(&ast.ExprStmt{X: &ast.Ident{Name: "stmt2"}})
	builder2.AddDeclaration(&ast.GenDecl{Tok: token.TYPE})

	builder1.Merge(builder2)

	if builder1.ImportCount() != 2 {
		t.Errorf("Expected 2 imports after merge, got %d", builder1.ImportCount())
	}

	if builder1.StatementCount() != 2 {
		t.Errorf("Expected 2 statements after merge, got %d", builder1.StatementCount())
	}

	if builder1.DeclarationCount() != 1 {
		t.Errorf("Expected 1 declaration after merge, got %d", builder1.DeclarationCount())
	}
}

func TestBuilder_BuildFile(t *testing.T) {
	builder := NewBuilder(BuilderConfig{PackageName: "test"})

	builder.AddImport("fmt")
	builder.AddStatement(&ast.ExprStmt{X: &ast.Ident{Name: "test"}})

	file := builder.BuildFile()

	if file == nil {
		t.Fatal("BuildFile returned nil")
	}

	if file.Name.Name != "test" {
		t.Errorf("Expected package name 'test', got '%s'", file.Name.Name)
	}

	if len(file.Imports) != 1 {
		t.Errorf("Expected 1 import, got %d", len(file.Imports))
	}

	if len(file.Decls) != 1 {
		t.Errorf("Expected 1 declaration, got %d", len(file.Decls))
	}
}

func TestBuilder_BuildFunction(t *testing.T) {
	builder := NewBuilder(BuilderConfig{PackageName: "test"})

	builder.AddStatement(&ast.ExprStmt{X: &ast.Ident{Name: "test"}})

	funcDecl := builder.BuildFunction("TestFunc", nil, []*ast.Field{}, []*ast.Field{})

	if funcDecl == nil {
		t.Fatal("BuildFunction returned nil")
	}

	if funcDecl.Name.Name != "TestFunc" {
		t.Errorf("Expected function name 'TestFunc', got '%s'", funcDecl.Name.Name)
	}

	if len(funcDecl.Body.List) != 1 {
		t.Errorf("Expected 1 statement in function body, got %d", len(funcDecl.Body.List))
	}
}

func TestBuilder_BuildBlock(t *testing.T) {
	builder := NewBuilder(BuilderConfig{PackageName: "test"})

	builder.AddStatement(&ast.ExprStmt{X: &ast.Ident{Name: "stmt1"}})
	builder.AddStatement(&ast.ExprStmt{X: &ast.Ident{Name: "stmt2"}})

	block := builder.BuildBlock()

	if block == nil {
		t.Fatal("BuildBlock returned nil")
	}

	if len(block.List) != 2 {
		t.Errorf("Expected 2 statements in block, got %d", len(block.List))
	}
}

func TestBuilder_GetConfig(t *testing.T) {
	config := BuilderConfig{
		PackageName:  "test",
		ImportPrefix: "github.com/test",
		UsePointers:  true,
	}

	builder := NewBuilder(config)
	retrievedConfig := builder.GetConfig()

	if retrievedConfig.PackageName != config.PackageName {
		t.Error("GetConfig should return the same package name")
	}

	if retrievedConfig.ImportPrefix != config.ImportPrefix {
		t.Error("GetConfig should return the same import prefix")
	}

	if retrievedConfig.UsePointers != config.UsePointers {
		t.Error("GetConfig should return the same UsePointers value")
	}
}
