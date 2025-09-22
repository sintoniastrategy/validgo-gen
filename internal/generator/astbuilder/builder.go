package astbuilder

import (
	"go/ast"
	"go/token"
	"sort"
	"strconv"
	"strings"
)

// BuilderConfig holds configuration for the AST builder
type BuilderConfig struct {
	PackageName  string
	ImportPrefix string
	UsePointers  bool
}

// Builder is the main AST builder that provides a fluent interface for building Go AST
type Builder struct {
	config  BuilderConfig
	imports map[string]bool
	stmts   []ast.Stmt
	decls   []ast.Decl
}

// NewBuilder creates a new AST builder with the given configuration
func NewBuilder(config BuilderConfig) *Builder {
	return &Builder{
		config:  config,
		imports: make(map[string]bool),
		stmts:   make([]ast.Stmt, 0),
		decls:   make([]ast.Decl, 0),
	}
}

// AddImport adds an import to the builder
func (b *Builder) AddImport(path string) *Builder {
	b.imports[path] = true
	return b
}

// AddStatement adds a statement to the builder
func (b *Builder) AddStatement(stmt ast.Stmt) *Builder {
	b.stmts = append(b.stmts, stmt)
	return b
}

// AddDeclaration adds a declaration to the builder
func (b *Builder) AddDeclaration(decl ast.Decl) *Builder {
	b.decls = append(b.decls, decl)
	return b
}

// AddStatements adds multiple statements to the builder
func (b *Builder) AddStatements(stmts []ast.Stmt) *Builder {
	b.stmts = append(b.stmts, stmts...)
	return b
}

// AddDeclarations adds multiple declarations to the builder
func (b *Builder) AddDeclarations(decls []ast.Decl) *Builder {
	b.decls = append(b.decls, decls...)
	return b
}

// Build returns the built AST components
func (b *Builder) Build() ([]ast.Stmt, []ast.Decl, []string) {
	// Sort imports for consistent output
	imports := make([]string, 0, len(b.imports))
	for imp := range b.imports {
		imports = append(imports, imp)
	}
	sort.Strings(imports)

	return b.stmts, b.decls, imports
}

// GetConfig returns the builder configuration
func (b *Builder) GetConfig() BuilderConfig {
	return b.config
}

// Clear clears all statements and declarations from the builder
func (b *Builder) Clear() *Builder {
	b.stmts = b.stmts[:0]
	b.decls = b.decls[:0]
	b.imports = make(map[string]bool)
	return b
}

// ClearStatements clears only the statements from the builder
func (b *Builder) ClearStatements() *Builder {
	b.stmts = b.stmts[:0]
	return b
}

// ClearDeclarations clears only the declarations from the builder
func (b *Builder) ClearDeclarations() *Builder {
	b.decls = b.decls[:0]
	return b
}

// HasImports returns true if the builder has any imports
func (b *Builder) HasImports() bool {
	return len(b.imports) > 0
}

// HasStatements returns true if the builder has any statements
func (b *Builder) HasStatements() bool {
	return len(b.stmts) > 0
}

// HasDeclarations returns true if the builder has any declarations
func (b *Builder) HasDeclarations() bool {
	return len(b.decls) > 0
}

// StatementCount returns the number of statements in the builder
func (b *Builder) StatementCount() int {
	return len(b.stmts)
}

// DeclarationCount returns the number of declarations in the builder
func (b *Builder) DeclarationCount() int {
	return len(b.decls)
}

// ImportCount returns the number of imports in the builder
func (b *Builder) ImportCount() int {
	return len(b.imports)
}

// Clone creates a copy of the builder with the same configuration
func (b *Builder) Clone() *Builder {
	clone := &Builder{
		config:  b.config,
		imports: make(map[string]bool),
		stmts:   make([]ast.Stmt, len(b.stmts)),
		decls:   make([]ast.Decl, len(b.decls)),
	}

	// Copy imports
	for imp := range b.imports {
		clone.imports[imp] = true
	}

	// Copy statements and declarations (shallow copy)
	copy(clone.stmts, b.stmts)
	copy(clone.decls, b.decls)

	return clone
}

// Merge merges another builder into this one
func (b *Builder) Merge(other *Builder) *Builder {
	// Merge imports
	for imp := range other.imports {
		b.imports[imp] = true
	}

	// Merge statements
	b.stmts = append(b.stmts, other.stmts...)

	// Merge declarations
	b.decls = append(b.decls, other.decls...)

	return b
}

// BuildFile creates a complete Go file AST
func (b *Builder) BuildFile() *ast.File {
	_, decls, imports := b.Build()

	file := &ast.File{
		Name:  ast.NewIdent(b.config.PackageName),
		Decls: make([]ast.Decl, 0),
	}

	// Add imports if any
	if len(imports) > 0 {
		importSpecs := make([]*ast.ImportSpec, 0, len(imports))
		for _, imp := range imports {
			importSpecs = append(importSpecs, &ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"` + imp + `"`,
				},
			})
		}

		file.Imports = importSpecs
		file.Decls = append(file.Decls, &ast.GenDecl{
			Tok:   token.IMPORT,
			Specs: make([]ast.Spec, 0, len(importSpecs)),
		})

		// Add import specs to the import declaration
		if len(file.Decls) > 0 {
			if importDecl, ok := file.Decls[0].(*ast.GenDecl); ok {
				for _, spec := range importSpecs {
					importDecl.Specs = append(importDecl.Specs, spec)
				}
			}
		}
	}

	// Add other declarations
	file.Decls = append(file.Decls, decls...)

	return file
}

// BuildFunction creates a function declaration with the given name and body
func (b *Builder) BuildFunction(name string, receiver *ast.Field, params, results []*ast.Field) *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: ast.NewIdent(name),
		Type: &ast.FuncType{
			Params:  &ast.FieldList{List: params},
			Results: &ast.FieldList{List: results},
		},
		Body: &ast.BlockStmt{List: b.stmts},
		Recv: func() *ast.FieldList {
			if receiver != nil {
				return &ast.FieldList{List: []*ast.Field{receiver}}
			}
			return nil
		}(),
	}
}

// BuildBlock creates a block statement with the current statements
func (b *Builder) BuildBlock() *ast.BlockStmt {
	return &ast.BlockStmt{List: b.stmts}
}

// Helper function to create a string literal
func (b *Builder) str(value string) *ast.BasicLit {
	return &ast.BasicLit{
		Kind:  token.STRING,
		Value: `"` + strings.ReplaceAll(value, `"`, `\"`) + `"`,
	}
}

// Helper function to create an integer literal
func (b *Builder) int(value int) *ast.BasicLit {
	return &ast.BasicLit{
		Kind:  token.INT,
		Value: strconv.Itoa(value),
	}
}

// Helper function to create a boolean literal
func (b *Builder) bool(value bool) *ast.Ident {
	if value {
		return ast.NewIdent("true")
	}
	return ast.NewIdent("false")
}

// Helper function to create a nil literal
func (b *Builder) nil() *ast.Ident {
	return ast.NewIdent("nil")
}

// Helper function to create an identifier
func (b *Builder) ident(name string) *ast.Ident {
	return ast.NewIdent(name)
}

// Helper function to create a selector expression
func (b *Builder) selector(x ast.Expr, sel string) *ast.SelectorExpr {
	return &ast.SelectorExpr{
		X:   x,
		Sel: ast.NewIdent(sel),
	}
}

// Helper function to create a call expression
func (b *Builder) call(fun ast.Expr, args ...ast.Expr) *ast.CallExpr {
	return &ast.CallExpr{
		Fun:  fun,
		Args: args,
	}
}

// Helper function to create a binary expression
func (b *Builder) binary(x ast.Expr, op token.Token, y ast.Expr) *ast.BinaryExpr {
	return &ast.BinaryExpr{
		X:  x,
		Op: op,
		Y:  y,
	}
}

// Helper function to create a unary expression
func (b *Builder) unary(op token.Token, x ast.Expr) *ast.UnaryExpr {
	return &ast.UnaryExpr{
		Op: op,
		X:  x,
	}
}

// Public Expression Methods

// Call creates a function call expression
func (b *Builder) Call(receiver, method string, args ...ast.Expr) ast.Expr {
	var fun ast.Expr
	if receiver != "" {
		fun = b.selector(b.ident(receiver), method)
	} else {
		fun = b.ident(method)
	}
	return b.call(fun, args...)
}

// Select creates a selector expression (receiver.field)
func (b *Builder) Select(receiver, field string) ast.Expr {
	return b.selector(b.ident(receiver), field)
}

// Ident creates an identifier expression
func (b *Builder) Ident(name string) ast.Expr {
	return b.ident(name)
}

// String creates a string literal expression
func (b *Builder) String(value string) ast.Expr {
	return b.str(value)
}

// Int creates an integer literal expression
func (b *Builder) Int(value int) ast.Expr {
	return b.int(value)
}

// Bool creates a boolean literal expression
func (b *Builder) Bool(value bool) ast.Expr {
	return b.bool(value)
}

// Nil creates a nil literal expression
func (b *Builder) Nil() ast.Expr {
	return b.nil()
}

// AddressOf creates an address-of expression (&expr)
func (b *Builder) AddressOf(expr ast.Expr) ast.Expr {
	return b.unary(token.AND, expr)
}

// Deref creates a dereference expression (*expr)
func (b *Builder) Deref(expr ast.Expr) ast.Expr {
	return b.unary(token.MUL, expr)
}

// Public Statement Methods

// DeclareVar creates a variable declaration statement
func (b *Builder) DeclareVar(name, typeName string, value ast.Expr) ast.Stmt {
	spec := &ast.ValueSpec{
		Names: []*ast.Ident{b.ident(name)},
		Type:  b.ident(typeName),
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

// Assign creates an assignment statement
func (b *Builder) Assign(lhs, rhs ast.Expr) ast.Stmt {
	return &ast.AssignStmt{
		Lhs: []ast.Expr{lhs},
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{rhs},
	}
}

// If creates an if statement
func (b *Builder) If(cond ast.Expr, body []ast.Stmt) ast.Stmt {
	return &ast.IfStmt{
		Cond: cond,
		Body: &ast.BlockStmt{List: body},
	}
}

// IfElse creates an if-else statement
func (b *Builder) IfElse(cond ast.Expr, ifBody, elseBody []ast.Stmt) ast.Stmt {
	stmt := &ast.IfStmt{
		Cond: cond,
		Body: &ast.BlockStmt{List: ifBody},
	}

	if len(elseBody) > 0 {
		stmt.Else = &ast.BlockStmt{List: elseBody}
	}

	return stmt
}

// Return creates a return statement
func (b *Builder) Return(values ...ast.Expr) ast.Stmt {
	return &ast.ReturnStmt{
		Results: values,
	}
}

// CallStmt creates a call statement
func (b *Builder) CallStmt(receiver, method string, args ...ast.Expr) ast.Stmt {
	var fun ast.Expr
	if receiver != "" {
		fun = b.selector(b.ident(receiver), method)
	} else {
		fun = b.ident(method)
	}

	return &ast.ExprStmt{
		X: b.call(fun, args...),
	}
}
