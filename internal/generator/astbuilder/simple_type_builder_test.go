package astbuilder

import (
	"go/ast"
	"testing"
)

func TestNewSimpleTypeBuilder(t *testing.T) {
	builder := NewSimpleTypeBuilder()

	if builder == nil {
		t.Fatal("NewSimpleTypeBuilder returned nil")
	}

	if builder.elements == nil {
		t.Fatal("elements slice is nil")
	}

	if len(builder.elements) != 0 {
		t.Errorf("Expected empty elements initially, got %d", len(builder.elements))
	}
}

func TestSimpleTypeBuilder_AddElement(t *testing.T) {
	builder := NewSimpleTypeBuilder()

	// Test adding a single element
	result := builder.AddElement("string")
	if result != builder {
		t.Error("AddElement should return the builder for chaining")
	}

	if len(builder.elements) != 1 {
		t.Errorf("Expected 1 element, got %d", len(builder.elements))
	}

	if builder.elements[0] != "string" {
		t.Errorf("Expected element 'string', got %s", builder.elements[0])
	}

	// Test adding another element
	builder.AddElement("Context")
	if len(builder.elements) != 2 {
		t.Errorf("Expected 2 elements, got %d", len(builder.elements))
	}

	if builder.elements[1] != "Context" {
		t.Errorf("Expected element 'Context', got %s", builder.elements[1])
	}

	// Test adding empty element (should be ignored)
	builder.AddElement("")
	if len(builder.elements) != 2 {
		t.Errorf("Expected 2 elements after adding empty, got %d", len(builder.elements))
	}
}

func TestSimpleTypeBuilder_AddElements(t *testing.T) {
	builder := NewSimpleTypeBuilder()

	// Test adding multiple elements
	result := builder.AddElements("context", "Context")
	if result != builder {
		t.Error("AddElements should return the builder for chaining")
	}

	if len(builder.elements) != 2 {
		t.Errorf("Expected 2 elements, got %d", len(builder.elements))
	}

	if builder.elements[0] != "context" {
		t.Errorf("Expected first element 'context', got %s", builder.elements[0])
	}

	if builder.elements[1] != "Context" {
		t.Errorf("Expected second element 'Context', got %s", builder.elements[1])
	}

	// Test adding elements with empty strings (should be ignored)
	builder.AddElements("", "Type", "")
	if len(builder.elements) != 3 {
		t.Errorf("Expected 3 elements after adding with empty strings, got %d", len(builder.elements))
	}

	if builder.elements[2] != "Type" {
		t.Errorf("Expected third element 'Type', got %s", builder.elements[2])
	}
}

func TestSimpleTypeBuilder_Build(t *testing.T) {
	// Test single element (should create ast.Ident)
	builder := NewSimpleTypeBuilder().AddElement("string")
	expr := builder.Build()

	if ident, ok := expr.(*ast.Ident); ok {
		if ident.Name != "string" {
			t.Errorf("Expected ident name 'string', got %s", ident.Name)
		}
	} else {
		t.Error("Single element should create ast.Ident")
	}

	// Test multiple elements (should create ast.SelectorExpr)
	builder = NewSimpleTypeBuilder().AddElements("context", "Context")
	expr = builder.Build()

	if selector, ok := expr.(*ast.SelectorExpr); ok {
		if ident, ok := selector.X.(*ast.Ident); ok {
			if ident.Name != "context" {
				t.Errorf("Expected selector X name 'context', got %s", ident.Name)
			}
		} else {
			t.Error("Selector X should be ast.Ident")
		}

		if selector.Sel.Name != "Context" {
			t.Errorf("Expected selector Sel name 'Context', got %s", selector.Sel.Name)
		}
	} else {
		t.Error("Multiple elements should create ast.SelectorExpr")
	}

	// Test three elements (should create nested selector)
	builder = NewSimpleTypeBuilder().AddElements("package", "subpackage", "Type")
	expr = builder.Build()

	if selector, ok := expr.(*ast.SelectorExpr); ok {
		if selector.Sel.Name != "Type" {
			t.Errorf("Expected outermost selector name 'Type', got %s", selector.Sel.Name)
		}

		if innerSelector, ok := selector.X.(*ast.SelectorExpr); ok {
			if innerSelector.Sel.Name != "subpackage" {
				t.Errorf("Expected inner selector name 'subpackage', got %s", innerSelector.Sel.Name)
			}
		} else {
			t.Error("Inner expression should be ast.SelectorExpr")
		}
	} else {
		t.Error("Three elements should create ast.SelectorExpr")
	}
}

func TestSimpleTypeBuilder_BuildWithoutElements(t *testing.T) {
	builder := NewSimpleTypeBuilder()

	defer func() {
		if r := recover(); r == nil {
			t.Error("Build should panic when no elements are present")
		}
	}()

	builder.Build()
}

func TestSimpleTypeBuilder_BuildAsIdent(t *testing.T) {
	builder := NewSimpleTypeBuilder().AddElement("string")
	ident := builder.BuildAsIdent()

	if ident.Name != "string" {
		t.Errorf("Expected ident name 'string', got %s", ident.Name)
	}

	// Test with multiple elements (should panic)
	builder = NewSimpleTypeBuilder().AddElements("context", "Context")

	defer func() {
		if r := recover(); r == nil {
			t.Error("BuildAsIdent should panic with multiple elements")
		}
	}()

	builder.BuildAsIdent()
}

func TestSimpleTypeBuilder_BuildAsSelector(t *testing.T) {
	builder := NewSimpleTypeBuilder().AddElements("context", "Context")
	selector := builder.BuildAsSelector()

	if selector.Sel.Name != "Context" {
		t.Errorf("Expected selector name 'Context', got %s", selector.Sel.Name)
	}

	// Test with single element (should panic)
	builder = NewSimpleTypeBuilder().AddElement("string")

	defer func() {
		if r := recover(); r == nil {
			t.Error("BuildAsSelector should panic with single element")
		}
	}()

	builder.BuildAsSelector()
}

func TestSimpleTypeBuilder_UtilityMethods(t *testing.T) {
	builder := NewSimpleTypeBuilder()

	// Test ElementCount
	if builder.ElementCount() != 0 {
		t.Errorf("Expected 0 elements initially, got %d", builder.ElementCount())
	}

	// Test HasElements
	if builder.HasElements() {
		t.Error("Expected HasElements to return false initially")
	}

	// Add elements
	builder.AddElements("context", "Context")

	// Test ElementCount after adding elements
	if builder.ElementCount() != 2 {
		t.Errorf("Expected 2 elements, got %d", builder.ElementCount())
	}

	// Test HasElements after adding elements
	if !builder.HasElements() {
		t.Error("Expected HasElements to return true after adding elements")
	}

	// Test GetElements
	elements := builder.GetElements()
	if len(elements) != 2 {
		t.Errorf("Expected 2 elements, got %d", len(elements))
	}

	if elements[0] != "context" {
		t.Errorf("Expected first element 'context', got %s", elements[0])
	}

	if elements[1] != "Context" {
		t.Errorf("Expected second element 'Context', got %s", elements[1])
	}

	// Test Clear
	result := builder.Clear()
	if result != builder {
		t.Error("Clear should return the builder for chaining")
	}

	if builder.ElementCount() != 0 {
		t.Errorf("Expected 0 elements after clear, got %d", builder.ElementCount())
	}

	if builder.HasElements() {
		t.Error("Expected HasElements to return false after clear")
	}
}

func TestSimpleTypeBuilder_Clone(t *testing.T) {
	builder := NewSimpleTypeBuilder().AddElements("context", "Context")
	clone := builder.Clone()

	if clone == builder {
		t.Error("Clone should return a different instance")
	}

	if clone.ElementCount() != 2 {
		t.Errorf("Expected clone to have 2 elements, got %d", clone.ElementCount())
	}

	elements := clone.GetElements()
	if elements[0] != "context" || elements[1] != "Context" {
		t.Error("Clone should have the same elements")
	}

	// Modify original
	builder.AddElement("Type")

	// Clone should be unaffected
	if clone.ElementCount() != 2 {
		t.Errorf("Expected clone to still have 2 elements, got %d", clone.ElementCount())
	}
}

func TestSimpleTypeBuilder_HelperMethods(t *testing.T) {
	// Test String()
	builder := String()
	if builder.ElementCount() != 1 {
		t.Errorf("Expected String() to have 1 element, got %d", builder.ElementCount())
	}
	if builder.elements[0] != "string" {
		t.Errorf("Expected String() element to be 'string', got %s", builder.elements[0])
	}

	// Test Int()
	builder = Int()
	if builder.ElementCount() != 1 {
		t.Errorf("Expected Int() to have 1 element, got %d", builder.ElementCount())
	}
	if builder.elements[0] != "int" {
		t.Errorf("Expected Int() element to be 'int', got %s", builder.elements[0])
	}

	// Test Error()
	builder = Error()
	if builder.ElementCount() != 1 {
		t.Errorf("Expected Error() to have 1 element, got %d", builder.ElementCount())
	}
	if builder.elements[0] != "error" {
		t.Errorf("Expected Error() element to be 'error', got %s", builder.elements[0])
	}

	// Test Context()
	builder = Context()
	if builder.ElementCount() != 2 {
		t.Errorf("Expected Context() to have 2 elements, got %d", builder.ElementCount())
	}
	if builder.elements[0] != "context" || builder.elements[1] != "Context" {
		t.Errorf("Expected Context() elements to be 'context' and 'Context', got %s and %s", builder.elements[0], builder.elements[1])
	}

	// Test Ident()
	builder = Ident("MyType")
	if builder.ElementCount() != 1 {
		t.Errorf("Expected Ident() to have 1 element, got %d", builder.ElementCount())
	}
	if builder.elements[0] != "MyType" {
		t.Errorf("Expected Ident() element to be 'MyType', got %s", builder.elements[0])
	}

	// Test Selector()
	builder = Selector("package", "Type")
	if builder.ElementCount() != 2 {
		t.Errorf("Expected Selector() to have 2 elements, got %d", builder.ElementCount())
	}
	if builder.elements[0] != "package" || builder.elements[1] != "Type" {
		t.Errorf("Expected Selector() elements to be 'package' and 'Type', got %s and %s", builder.elements[0], builder.elements[1])
	}
}

func TestSimpleTypeBuilder_Pointer(t *testing.T) {
	builder := NewSimpleTypeBuilder().AddElement("string")
	ptrExpr := builder.Pointer()

	if starExpr, ok := ptrExpr.(*ast.StarExpr); ok {
		if ident, ok := starExpr.X.(*ast.Ident); ok {
			if ident.Name != "string" {
				t.Errorf("Expected pointer target to be 'string', got %s", ident.Name)
			}
		} else {
			t.Error("Pointer target should be ast.Ident")
		}
	} else {
		t.Error("Pointer should create ast.StarExpr")
	}
}

func TestSimpleTypeBuilder_Slice(t *testing.T) {
	builder := NewSimpleTypeBuilder().AddElement("string")
	sliceExpr := builder.Slice()

	if arrayType, ok := sliceExpr.(*ast.ArrayType); ok {
		if ident, ok := arrayType.Elt.(*ast.Ident); ok {
			if ident.Name != "string" {
				t.Errorf("Expected slice element type to be 'string', got %s", ident.Name)
			}
		} else {
			t.Error("Slice element type should be ast.Ident")
		}

		if arrayType.Len != nil {
			t.Error("Slice should have nil length")
		}
	} else {
		t.Error("Slice should create ast.ArrayType")
	}
}

func TestSimpleTypeBuilder_AsPointer(t *testing.T) {
	// Test with single element
	builder := NewSimpleTypeBuilder().AddElement("string").AsPointer(true)
	expr := builder.Build()

	if starExpr, ok := expr.(*ast.StarExpr); ok {
		if ident, ok := starExpr.X.(*ast.Ident); ok {
			if ident.Name != "string" {
				t.Errorf("Expected pointer target to be 'string', got %s", ident.Name)
			}
		} else {
			t.Error("Pointer target should be ast.Ident")
		}
	} else {
		t.Error("AsPointer(true) should create ast.StarExpr")
	}

	// Test with multiple elements (selector expression)
	builder = NewSimpleTypeBuilder().AddElements("context", "Context").AsPointer(true)
	expr = builder.Build()

	if starExpr, ok := expr.(*ast.StarExpr); ok {
		if selector, ok := starExpr.X.(*ast.SelectorExpr); ok {
			if selector.Sel.Name != "Context" {
				t.Errorf("Expected selector name 'Context', got %s", selector.Sel.Name)
			}
		} else {
			t.Error("Pointer target should be ast.SelectorExpr")
		}
	} else {
		t.Error("AsPointer(true) with selector should create ast.StarExpr")
	}

	// Test with AsPointer(false) - should not create pointer
	builder = NewSimpleTypeBuilder().AddElement("string").AsPointer(false)
	expr = builder.Build()

	if _, ok := expr.(*ast.StarExpr); ok {
		t.Error("AsPointer(false) should not create ast.StarExpr")
	}

	if ident, ok := expr.(*ast.Ident); ok {
		if ident.Name != "string" {
			t.Errorf("Expected identifier 'string', got %s", ident.Name)
		}
	} else {
		t.Error("AsPointer(false) should create ast.Ident")
	}

	// Test default behavior (should not create pointer)
	builder = NewSimpleTypeBuilder().AddElement("string")
	expr = builder.Build()

	if _, ok := expr.(*ast.StarExpr); ok {
		t.Error("Default behavior should not create ast.StarExpr")
	}

	if ident, ok := expr.(*ast.Ident); ok {
		if ident.Name != "string" {
			t.Errorf("Expected identifier 'string', got %s", ident.Name)
		}
	} else {
		t.Error("Default behavior should create ast.Ident")
	}
}

func TestSimpleTypeBuilder_AsPointerMethodChaining(t *testing.T) {
	builder := NewSimpleTypeBuilder().
		AddElement("string").
		AsPointer(true)

	result := builder.AsPointer(false)
	if result != builder {
		t.Error("AsPointer should return the builder for chaining")
	}

	expr := builder.Build()
	if _, ok := expr.(*ast.StarExpr); ok {
		t.Error("AsPointer(false) should not create ast.StarExpr")
	}

	// Test chaining with true
	builder.AsPointer(true)
	expr = builder.Build()
	if _, ok := expr.(*ast.StarExpr); !ok {
		t.Error("AsPointer(true) should create ast.StarExpr")
	}
}

func TestSimpleTypeBuilder_AsPointerClone(t *testing.T) {
	original := NewSimpleTypeBuilder().
		AddElement("string").
		AsPointer(true)

	clone := original.Clone()

	// Test that clone has the same asPointer setting
	originalExpr := original.Build()
	cloneExpr := clone.Build()

	if _, ok := originalExpr.(*ast.StarExpr); !ok {
		t.Error("Original should create ast.StarExpr")
	}

	if _, ok := cloneExpr.(*ast.StarExpr); !ok {
		t.Error("Clone should create ast.StarExpr")
	}

	// Test that modifying clone doesn't affect original
	clone.AsPointer(false)
	originalExpr = original.Build()
	cloneExpr = clone.Build()

	if _, ok := originalExpr.(*ast.StarExpr); !ok {
		t.Error("Original should still create ast.StarExpr after clone modification")
	}

	if _, ok := cloneExpr.(*ast.StarExpr); ok {
		t.Error("Clone should not create ast.StarExpr after modification")
	}
}
