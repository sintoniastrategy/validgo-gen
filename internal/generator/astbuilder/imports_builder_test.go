package astbuilder

import (
	"go/token"
	"testing"
)

func TestNewImportsBuilder(t *testing.T) {
	packagePrefix := "github.com/test"
	builder := NewImportsBuilder(packagePrefix)

	if builder == nil {
		t.Fatal("NewImportsBuilder returned nil")
	}

	if builder.packagePrefix != packagePrefix {
		t.Errorf("Expected packagePrefix %s, got %s", packagePrefix, builder.packagePrefix)
	}

	if builder.imports == nil {
		t.Fatal("imports map is nil")
	}

	if len(builder.imports) != 0 {
		t.Errorf("Expected empty imports, got %d imports", len(builder.imports))
	}
}

func TestImportsBuilder_AddImport(t *testing.T) {
	builder := NewImportsBuilder("github.com/test")

	// Test adding a single import
	result := builder.AddImport("fmt")
	if result != builder {
		t.Error("AddImport should return the builder for chaining")
	}

	if !builder.HasImport("fmt") {
		t.Error("Import 'fmt' should be present")
	}

	if builder.Count() != 1 {
		t.Errorf("Expected 1 import, got %d", builder.Count())
	}
}

func TestImportsBuilder_AddImports(t *testing.T) {
	builder := NewImportsBuilder("github.com/test")

	// Test adding multiple imports
	result := builder.AddImports("fmt", "strings", "context")
	if result != builder {
		t.Error("AddImports should return the builder for chaining")
	}

	expectedImports := []string{"fmt", "strings", "context"}
	for _, imp := range expectedImports {
		if !builder.HasImport(imp) {
			t.Errorf("Import '%s' should be present", imp)
		}
	}

	if builder.Count() != 3 {
		t.Errorf("Expected 3 imports, got %d", builder.Count())
	}
}

func TestImportsBuilder_Build(t *testing.T) {
	builder := NewImportsBuilder("github.com/test")

	// Add imports in different categories
	builder.AddImports(
		"fmt",                             // system import
		"strings",                         // system import
		"github.com/gin-gonic/gin",        // library import
		"github.com/test/package",         // user import
		"github.com/test/another-package", // user import
	)

	specs, declSpecs := builder.Build()

	// Check that we have the right number of specs
	expectedCount := 5
	if len(specs) != expectedCount {
		t.Errorf("Expected %d import specs, got %d", expectedCount, len(specs))
	}

	if len(declSpecs) != expectedCount {
		t.Errorf("Expected %d decl specs, got %d", expectedCount, len(declSpecs))
	}

	// Verify the specs are of the correct type
	for i, spec := range specs {
		if spec == nil {
			t.Errorf("Import spec %d is nil", i)
			continue
		}

		if spec.Path == nil {
			t.Errorf("Import spec %d has nil path", i)
			continue
		}

		// Check that it's a basic literal with a string value
		if spec.Path.Kind != token.STRING {
			t.Errorf("Import spec %d should be a string literal", i)
		}
	}

	// Verify declSpecs are the same as specs
	for i, declSpec := range declSpecs {
		if declSpec != specs[i] {
			t.Errorf("Decl spec %d should match import spec %d", i, i)
		}
	}
}

func TestImportsBuilder_BuildWithAlias(t *testing.T) {
	builder := NewImportsBuilder("github.com/test")
	builder.AddImports("fmt", "github.com/gin-gonic/gin")

	aliases := map[string]string{
		"fmt":                      "f",
		"github.com/gin-gonic/gin": "gin",
	}

	specs, declSpecs := builder.BuildWithAlias(aliases)

	if len(specs) != 2 {
		t.Errorf("Expected 2 import specs, got %d", len(specs))
	}

	if len(declSpecs) != 2 {
		t.Errorf("Expected 2 decl specs, got %d", len(declSpecs))
	}

	// Check that aliases are applied
	for _, spec := range specs {
		if spec.Path == nil {
			continue
		}

		pathValue := spec.Path.Value[1 : len(spec.Path.Value)-1] // Remove quotes

		if pathValue == "fmt" && (spec.Name == nil || spec.Name.Name != "f") {
			t.Error("fmt import should have alias 'f'")
		}

		if pathValue == "github.com/gin-gonic/gin" && (spec.Name == nil || spec.Name.Name != "gin") {
			t.Error("gin import should have alias 'gin'")
		}
	}
}

func TestImportsBuilder_HasImport(t *testing.T) {
	builder := NewImportsBuilder("github.com/test")

	// Test with no imports
	if builder.HasImport("fmt") {
		t.Error("Should not have 'fmt' import when none added")
	}

	// Add an import and test
	builder.AddImport("fmt")
	if !builder.HasImport("fmt") {
		t.Error("Should have 'fmt' import after adding it")
	}

	// Test with different import
	if builder.HasImport("strings") {
		t.Error("Should not have 'strings' import when not added")
	}
}

func TestImportsBuilder_RemoveImport(t *testing.T) {
	builder := NewImportsBuilder("github.com/test")
	builder.AddImports("fmt", "strings", "context")

	// Remove an import
	result := builder.RemoveImport("strings")
	if result != builder {
		t.Error("RemoveImport should return the builder for chaining")
	}

	if builder.HasImport("strings") {
		t.Error("Should not have 'strings' import after removing it")
	}

	// Verify other imports are still there
	if !builder.HasImport("fmt") {
		t.Error("Should still have 'fmt' import")
	}

	if !builder.HasImport("context") {
		t.Error("Should still have 'context' import")
	}

	if builder.Count() != 2 {
		t.Errorf("Expected 2 imports, got %d", builder.Count())
	}
}

func TestImportsBuilder_Clear(t *testing.T) {
	builder := NewImportsBuilder("github.com/test")
	builder.AddImports("fmt", "strings", "context")

	if builder.Count() != 3 {
		t.Errorf("Expected 3 imports before clear, got %d", builder.Count())
	}

	result := builder.Clear()
	if result != builder {
		t.Error("Clear should return the builder for chaining")
	}

	if builder.Count() != 0 {
		t.Errorf("Expected 0 imports after clear, got %d", builder.Count())
	}

	if builder.HasImport("fmt") {
		t.Error("Should not have any imports after clear")
	}
}

func TestImportsBuilder_Count(t *testing.T) {
	builder := NewImportsBuilder("github.com/test")

	if builder.Count() != 0 {
		t.Errorf("Expected 0 imports initially, got %d", builder.Count())
	}

	builder.AddImport("fmt")
	if builder.Count() != 1 {
		t.Errorf("Expected 1 import, got %d", builder.Count())
	}

	builder.AddImports("strings", "context")
	if builder.Count() != 3 {
		t.Errorf("Expected 3 imports, got %d", builder.Count())
	}

	builder.RemoveImport("fmt")
	if builder.Count() != 2 {
		t.Errorf("Expected 2 imports after removal, got %d", builder.Count())
	}
}

func TestImportsBuilder_GetImports(t *testing.T) {
	builder := NewImportsBuilder("github.com/test")
	builder.AddImports("context", "fmt", "strings")

	imports := builder.GetImports()

	if len(imports) != 3 {
		t.Errorf("Expected 3 imports, got %d", len(imports))
	}

	// Check that imports are sorted
	expected := []string{"context", "fmt", "strings"}
	for i, imp := range imports {
		if imp != expected[i] {
			t.Errorf("Expected import %d to be %s, got %s", i, expected[i], imp)
		}
	}
}

func TestImportsBuilder_ImportOrdering(t *testing.T) {
	builder := NewImportsBuilder("github.com/test")

	// Add imports in different categories
	builder.AddImports(
		"github.com/test/user-package",           // user import
		"fmt",                                    // system import
		"github.com/gin-gonic/gin",               // library import
		"strings",                                // system import
		"github.com/test/another-user-package",   // user import
		"github.com/go-playground/validator/v10", // library import
	)

	specs, _ := builder.Build()

	if len(specs) != 6 {
		t.Errorf("Expected 6 import specs, got %d", len(specs))
	}

	// Extract import paths
	var paths []string
	for _, spec := range specs {
		path := spec.Path.Value[1 : len(spec.Path.Value)-1] // Remove quotes
		paths = append(paths, path)
	}

	// System imports should come first
	systemImports := []string{"fmt", "strings"}
	for i, expected := range systemImports {
		if paths[i] != expected {
			t.Errorf("Expected system import %d to be %s, got %s", i, expected, paths[i])
		}
	}

	// Library imports should come second
	libStart := 2
	libraryImports := []string{"github.com/gin-gonic/gin", "github.com/go-playground/validator/v10"}
	for i, expected := range libraryImports {
		if paths[libStart+i] != expected {
			t.Errorf("Expected library import %d to be %s, got %s", i, expected, paths[libStart+i])
		}
	}

	// User imports should come last
	userStart := 4
	userImports := []string{"github.com/test/another-user-package", "github.com/test/user-package"}
	for i, expected := range userImports {
		if paths[userStart+i] != expected {
			t.Errorf("Expected user import %d to be %s, got %s", i, expected, paths[userStart+i])
		}
	}
}

func TestImportsBuilder_MethodChaining(t *testing.T) {
	builder := NewImportsBuilder("github.com/test")

	// Test method chaining
	result := builder.
		AddImport("fmt").
		AddImports("strings", "context").
		AddImport("io").
		RemoveImport("context").
		AddImport("os")

	if result != builder {
		t.Error("Method chaining should return the same builder instance")
	}

	expectedImports := []string{"fmt", "io", "os", "strings"}
	actualImports := builder.GetImports()

	if len(actualImports) != len(expectedImports) {
		t.Errorf("Expected %d imports, got %d", len(expectedImports), len(actualImports))
	}

	for i, expected := range expectedImports {
		if actualImports[i] != expected {
			t.Errorf("Expected import %d to be %s, got %s", i, expected, actualImports[i])
		}
	}
}
