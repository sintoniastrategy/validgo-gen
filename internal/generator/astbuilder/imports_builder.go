package astbuilder

import (
	"go/ast"
	"go/token"
	"sort"
	"strings"
)

// ImportsBuilder provides a fluent interface for building import specifications
type ImportsBuilder struct {
	imports       map[string]bool
	packagePrefix string
}

// NewImportsBuilder creates a new ImportsBuilder
func NewImportsBuilder(packagePrefix string) *ImportsBuilder {
	return &ImportsBuilder{
		imports:       make(map[string]bool),
		packagePrefix: packagePrefix,
	}
}

// AddImport adds an import path to the builder
// Returns the builder for method chaining
func (ib *ImportsBuilder) AddImport(path string) *ImportsBuilder {
	ib.imports[path] = true
	return ib
}

// AddImports adds multiple import paths to the builder
// Returns the builder for method chaining
func (ib *ImportsBuilder) AddImports(paths ...string) *ImportsBuilder {
	for _, path := range paths {
		ib.imports[path] = true
	}
	return ib
}

// Build compiles the list of imports into importSpec and importDecl
// Returns ([]*ast.ImportSpec, []ast.Spec) similar to GenerateImportsSpecs method
func (ib *ImportsBuilder) Build() ([]*ast.ImportSpec, []ast.Spec) {
	// Convert map to slice
	var systemImports []string
	var libImports []string
	var myImports []string

	for path := range ib.imports {
		if strings.HasPrefix(path, ib.packagePrefix) {
			myImports = append(myImports, path)
			continue
		}

		prefix := strings.SplitN(path, "/", 2)[0]
		if strings.Contains(prefix, ".") {
			libImports = append(libImports, path)
			continue
		}

		systemImports = append(systemImports, path)
	}

	// Sort imports for consistent output
	sort.Strings(systemImports)
	sort.Strings(libImports)
	sort.Strings(myImports)

	// Build import specs in the correct order
	specs := make([]*ast.ImportSpec, 0, len(ib.imports))

	// Add system imports first
	for _, path := range systemImports {
		specs = append(specs, &ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: `"` + path + `"`}})
	}

	// Add library imports second
	for _, path := range libImports {
		specs = append(specs, &ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: `"` + path + `"`}})
	}

	// Add user imports last
	for _, path := range myImports {
		specs = append(specs, &ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: `"` + path + `"`}})
	}

	// Convert to decl specs
	declSpecs := make([]ast.Spec, 0, len(specs))
	for _, spec := range specs {
		declSpecs = append(declSpecs, spec)
	}

	return specs, declSpecs
}

// BuildWithAlias compiles the list of imports with optional aliases
// Returns ([]*ast.ImportSpec, []ast.Spec) with alias support
func (ib *ImportsBuilder) BuildWithAlias(aliases map[string]string) ([]*ast.ImportSpec, []ast.Spec) {
	// Convert map to slice
	var systemImports []string
	var libImports []string
	var myImports []string

	for path := range ib.imports {
		if strings.HasPrefix(path, ib.packagePrefix) {
			myImports = append(myImports, path)
			continue
		}

		prefix := strings.SplitN(path, "/", 2)[0]
		if strings.Contains(prefix, ".") {
			libImports = append(libImports, path)
			continue
		}

		systemImports = append(systemImports, path)
	}

	// Sort imports for consistent output
	sort.Strings(systemImports)
	sort.Strings(libImports)
	sort.Strings(myImports)

	// Build import specs in the correct order
	specs := make([]*ast.ImportSpec, 0, len(ib.imports))

	// Add system imports first
	for _, path := range systemImports {
		spec := &ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: `"` + path + `"`}}
		if alias, exists := aliases[path]; exists {
			spec.Name = ast.NewIdent(alias)
		}
		specs = append(specs, spec)
	}

	// Add library imports second
	for _, path := range libImports {
		spec := &ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: `"` + path + `"`}}
		if alias, exists := aliases[path]; exists {
			spec.Name = ast.NewIdent(alias)
		}
		specs = append(specs, spec)
	}

	// Add user imports last
	for _, path := range myImports {
		spec := &ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: `"` + path + `"`}}
		if alias, exists := aliases[path]; exists {
			spec.Name = ast.NewIdent(alias)
		}
		specs = append(specs, spec)
	}

	// Convert to decl specs
	declSpecs := make([]ast.Spec, 0, len(specs))
	for _, spec := range specs {
		declSpecs = append(declSpecs, spec)
	}

	return specs, declSpecs
}

// HasImport checks if an import path is already added
func (ib *ImportsBuilder) HasImport(path string) bool {
	return ib.imports[path]
}

// RemoveImport removes an import path from the builder
func (ib *ImportsBuilder) RemoveImport(path string) *ImportsBuilder {
	delete(ib.imports, path)
	return ib
}

// Clear removes all imports from the builder
func (ib *ImportsBuilder) Clear() *ImportsBuilder {
	ib.imports = make(map[string]bool)
	return ib
}

// Count returns the number of imports
func (ib *ImportsBuilder) Count() int {
	return len(ib.imports)
}

// GetImports returns a slice of all import paths
func (ib *ImportsBuilder) GetImports() []string {
	imports := make([]string, 0, len(ib.imports))
	for path := range ib.imports {
		imports = append(imports, path)
	}
	sort.Strings(imports)
	return imports
}
