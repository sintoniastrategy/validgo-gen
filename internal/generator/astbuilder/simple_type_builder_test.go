package astbuilder

import (
	"go/ast"
	"go/token"
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

func TestNewArrayTypeBuilder(t *testing.T) {
	builder := NewArrayTypeBuilder()

	if builder == nil {
		t.Fatal("NewArrayTypeBuilder returned nil")
	}

	if builder.HasElement() {
		t.Error("Expected no element initially")
	}
}

func TestArrayTypeBuilder_WithElement(t *testing.T) {
	builder := NewArrayTypeBuilder()
	stringBuilder := String()

	result := builder.WithElement(stringBuilder)
	if result != builder {
		t.Error("WithElement should return the builder for chaining")
	}

	if !builder.HasElement() {
		t.Error("Expected element to be set")
	}

	if builder.GetElement() != stringBuilder {
		t.Error("Expected element to be the same reference")
	}
}

func TestArrayTypeBuilder_WithElementNil(t *testing.T) {
	builder := NewArrayTypeBuilder()

	defer func() {
		if r := recover(); r == nil {
			t.Error("WithElement should panic when element is nil")
		}
	}()

	builder.WithElement(nil)
}

func TestArrayTypeBuilder_Build(t *testing.T) {
	// Test with SimpleTypeBuilder
	builder := NewArrayTypeBuilder().WithElement(String())
	expr := builder.Build()

	if arrayType, ok := expr.(*ast.ArrayType); ok {
		if ident, ok := arrayType.Elt.(*ast.Ident); ok {
			if ident.Name != "string" {
				t.Errorf("Expected element type 'string', got %s", ident.Name)
			}
		} else {
			t.Error("Element type should be ast.Ident")
		}

		if arrayType.Len != nil {
			t.Error("Array type should have nil length for slices")
		}
	} else {
		t.Error("Build should create ast.ArrayType")
	}
}

func TestArrayTypeBuilder_BuildWithoutElement(t *testing.T) {
	builder := NewArrayTypeBuilder()

	defer func() {
		if r := recover(); r == nil {
			t.Error("Build should panic when no element is set")
		}
	}()

	builder.Build()
}

func TestArrayTypeBuilder_NestedArrays(t *testing.T) {
	// Test nested arrays: [][]string
	innerBuilder := NewArrayTypeBuilder().WithElement(String())
	outerBuilder := NewArrayTypeBuilder().WithElement(innerBuilder)

	expr := outerBuilder.Build()

	if outerArrayType, ok := expr.(*ast.ArrayType); ok {
		if innerArrayType, ok := outerArrayType.Elt.(*ast.ArrayType); ok {
			if ident, ok := innerArrayType.Elt.(*ast.Ident); ok {
				if ident.Name != "string" {
					t.Errorf("Expected inner element type 'string', got %s", ident.Name)
				}
			} else {
				t.Error("Inner element type should be ast.Ident")
			}
		} else {
			t.Error("Inner element should be ast.ArrayType")
		}
	} else {
		t.Error("Outer expression should be ast.ArrayType")
	}
}

func TestArrayTypeBuilder_UtilityMethods(t *testing.T) {
	builder := NewArrayTypeBuilder()

	// Test HasElement
	if builder.HasElement() {
		t.Error("Expected HasElement to return false initially")
	}

	builder.WithElement(String())
	if !builder.HasElement() {
		t.Error("Expected HasElement to return true after setting element")
	}

	// Test GetElement
	element := builder.GetElement()
	if element == nil {
		t.Error("GetElement should not return nil")
	}

	// Verify it's the same element reference
	if element != builder.element {
		t.Error("GetElement should return the same element reference")
	}
}

func TestArrayTypeBuilder_Clone(t *testing.T) {
	stringBuilder := String()
	original := NewArrayTypeBuilder().WithElement(stringBuilder)
	clone := original.Clone()

	// Test that clone has the same element type
	if !clone.HasElement() {
		t.Error("Clone should have an element")
	}

	// Test that clone has different element reference
	if clone.GetElement() == original.GetElement() {
		t.Error("Clone should have a different element reference")
	}

	// Test that modifying clone doesn't affect original
	clone.WithElement(Int())
	if original.GetElement() != stringBuilder {
		t.Error("Modifying clone should not affect original")
	}
}

func TestArrayTypeBuilder_CloneNested(t *testing.T) {
	// Test cloning nested arrays
	inner := NewArrayTypeBuilder().WithElement(String())
	outer := NewArrayTypeBuilder().WithElement(inner)

	clone := outer.Clone()

	// Test that clone has nested structure
	if !clone.HasElement() {
		t.Error("Clone should have an element")
	}

	if innerClone, ok := clone.GetElement().(*ArrayTypeBuilder); ok {
		if !innerClone.HasElement() {
			t.Error("Inner clone should have an element")
		}
	} else {
		t.Error("Clone element should be ArrayTypeBuilder")
	}
}

func TestArrayTypeBuilder_HelperFunctions(t *testing.T) {
	// Test StringSlice
	stringSlice := StringSlice()
	if !stringSlice.HasElement() {
		t.Error("StringSlice should have an element")
	}

	// Test IntSlice
	intSlice := IntSlice()
	if !intSlice.HasElement() {
		t.Error("IntSlice should have an element")
	}

	// Test BoolSlice
	boolSlice := BoolSlice()
	if !boolSlice.HasElement() {
		t.Error("BoolSlice should have an element")
	}

	// Test ErrorSlice
	errorSlice := ErrorSlice()
	if !errorSlice.HasElement() {
		t.Error("ErrorSlice should have an element")
	}

	// Test ContextSlice
	contextSlice := ContextSlice()
	if !contextSlice.HasElement() {
		t.Error("ContextSlice should have an element")
	}

	// Test IdentSlice
	identSlice := IdentSlice("CustomType")
	if !identSlice.HasElement() {
		t.Error("IdentSlice should have an element")
	}

	// Test SelectorSlice
	selectorSlice := SelectorSlice("pkg", "Type")
	if !selectorSlice.HasElement() {
		t.Error("SelectorSlice should have an element")
	}

	// Test SliceOf
	sliceOf := SliceOf(String())
	if !sliceOf.HasElement() {
		t.Error("SliceOf should have an element")
	}
}

func TestArrayTypeBuilder_ComplexNesting(t *testing.T) {
	// Test complex nesting: [][]context.Context
	contextSlice := ContextSlice()
	nestedSlice := SliceOf(contextSlice)

	expr := nestedSlice.Build()

	if outerArrayType, ok := expr.(*ast.ArrayType); ok {
		if innerArrayType, ok := outerArrayType.Elt.(*ast.ArrayType); ok {
			if selector, ok := innerArrayType.Elt.(*ast.SelectorExpr); ok {
				if selector.Sel.Name != "Context" {
					t.Errorf("Expected selector name 'Context', got %s", selector.Sel.Name)
				}
			} else {
				t.Error("Inner element should be ast.SelectorExpr")
			}
		} else {
			t.Error("Inner element should be ast.ArrayType")
		}
	} else {
		t.Error("Outer expression should be ast.ArrayType")
	}
}

func TestArrayTypeBuilder_MethodChaining(t *testing.T) {
	builder := NewArrayTypeBuilder().
		WithElement(String())

	if !builder.HasElement() {
		t.Error("Method chaining should work correctly")
	}

	// Test that all operations return the same builder
	operations := []func() *ArrayTypeBuilder{
		func() *ArrayTypeBuilder { return builder.WithElement(Int()) },
	}

	for i, op := range operations {
		if op() != builder {
			t.Errorf("Operation %d should return the same builder", i)
		}
	}
}

func TestNewTypeAliasBuilder(t *testing.T) {
	builder := NewTypeAliasBuilder()

	if builder == nil {
		t.Fatal("NewTypeAliasBuilder returned nil")
	}

	if builder.HasName() {
		t.Error("Expected no name initially")
	}

	if builder.HasType() {
		t.Error("Expected no type initially")
	}
}

func TestTypeAliasBuilder_WithName(t *testing.T) {
	builder := NewTypeAliasBuilder()

	result := builder.WithName("MyAlias")
	if result != builder {
		t.Error("WithName should return the builder for chaining")
	}

	if !builder.HasName() {
		t.Error("Expected name to be set")
	}

	if builder.GetName() != "MyAlias" {
		t.Errorf("Expected name 'MyAlias', got %s", builder.GetName())
	}
}

func TestTypeAliasBuilder_WithNameEmpty(t *testing.T) {
	builder := NewTypeAliasBuilder()

	defer func() {
		if r := recover(); r == nil {
			t.Error("WithName should panic when name is empty")
		}
	}()

	builder.WithName("")
}

func TestTypeAliasBuilder_WithType(t *testing.T) {
	builder := NewTypeAliasBuilder()
	stringBuilder := String()

	result := builder.WithType(stringBuilder)
	if result != builder {
		t.Error("WithType should return the builder for chaining")
	}

	if !builder.HasType() {
		t.Error("Expected type to be set")
	}

	if builder.GetType() != stringBuilder {
		t.Error("Expected type to be the same reference")
	}
}

func TestTypeAliasBuilder_WithTypeNil(t *testing.T) {
	builder := NewTypeAliasBuilder()

	defer func() {
		if r := recover(); r == nil {
			t.Error("WithType should panic when type is nil")
		}
	}()

	builder.WithType(nil)
}

func TestTypeAliasBuilder_Build(t *testing.T) {
	// Test with SimpleTypeBuilder
	builder := NewTypeAliasBuilder().
		WithName("StringAlias").
		WithType(String())

	spec := builder.Build()

	if spec.Name.Name != "StringAlias" {
		t.Errorf("Expected name 'StringAlias', got %s", spec.Name.Name)
	}

	if ident, ok := spec.Type.(*ast.Ident); ok {
		if ident.Name != "string" {
			t.Errorf("Expected type 'string', got %s", ident.Name)
		}
	} else {
		t.Error("Type should be ast.Ident")
	}
}

func TestTypeAliasBuilder_BuildWithArrayType(t *testing.T) {
	// Test with ArrayTypeBuilder
	builder := NewTypeAliasBuilder().
		WithName("StringSliceAlias").
		WithType(StringSlice())

	spec := builder.Build()

	if spec.Name.Name != "StringSliceAlias" {
		t.Errorf("Expected name 'StringSliceAlias', got %s", spec.Name.Name)
	}

	if arrayType, ok := spec.Type.(*ast.ArrayType); ok {
		if ident, ok := arrayType.Elt.(*ast.Ident); ok {
			if ident.Name != "string" {
				t.Errorf("Expected array element type 'string', got %s", ident.Name)
			}
		} else {
			t.Error("Array element should be ast.Ident")
		}
	} else {
		t.Error("Type should be ast.ArrayType")
	}
}

func TestTypeAliasBuilder_BuildWithoutName(t *testing.T) {
	builder := NewTypeAliasBuilder().WithType(String())

	defer func() {
		if r := recover(); r == nil {
			t.Error("Build should panic when no name is set")
		}
	}()

	builder.Build()
}

func TestTypeAliasBuilder_BuildWithoutType(t *testing.T) {
	builder := NewTypeAliasBuilder().WithName("MyAlias")

	defer func() {
		if r := recover(); r == nil {
			t.Error("Build should panic when no type is set")
		}
	}()

	builder.Build()
}

func TestTypeAliasBuilder_BuildAsDeclaration(t *testing.T) {
	builder := NewTypeAliasBuilder().
		WithName("MyAlias").
		WithType(String())

	decl := builder.BuildAsDeclaration()

	if decl.Tok != token.TYPE {
		t.Error("Declaration should have TYPE token")
	}

	if len(decl.Specs) != 1 {
		t.Error("Declaration should have exactly one spec")
	}

	if typeSpec, ok := decl.Specs[0].(*ast.TypeSpec); ok {
		if typeSpec.Name.Name != "MyAlias" {
			t.Errorf("Expected spec name 'MyAlias', got %s", typeSpec.Name.Name)
		}
	} else {
		t.Error("Spec should be ast.TypeSpec")
	}
}

func TestTypeAliasBuilder_UtilityMethods(t *testing.T) {
	builder := NewTypeAliasBuilder()

	// Test HasName
	if builder.HasName() {
		t.Error("Expected HasName to return false initially")
	}

	builder.WithName("TestAlias")
	if !builder.HasName() {
		t.Error("Expected HasName to return true after setting name")
	}

	// Test GetName
	name := builder.GetName()
	if name != "TestAlias" {
		t.Errorf("Expected name 'TestAlias', got %s", name)
	}

	// Test HasType
	if builder.HasType() {
		t.Error("Expected HasType to return false initially")
	}

	builder.WithType(String())
	if !builder.HasType() {
		t.Error("Expected HasType to return true after setting type")
	}

	// Test GetType
	typeBuilder := builder.GetType()
	if typeBuilder == nil {
		t.Error("GetType should not return nil")
	}
}

func TestTypeAliasBuilder_Clone(t *testing.T) {
	stringBuilder := String()
	original := NewTypeAliasBuilder().
		WithName("OriginalAlias").
		WithType(stringBuilder)

	clone := original.Clone()

	// Test that clone has the same name
	if clone.GetName() != "OriginalAlias" {
		t.Error("Clone should have the same name")
	}

	// Test that clone has different type reference
	if clone.GetType() == original.GetType() {
		t.Error("Clone should have a different type reference")
	}

	// Test that modifying clone doesn't affect original
	clone.WithName("CloneAlias")
	if original.GetName() != "OriginalAlias" {
		t.Error("Modifying clone should not affect original")
	}
}

func TestTypeAliasBuilder_CloneWithArrayType(t *testing.T) {
	// Test cloning with ArrayTypeBuilder
	original := NewTypeAliasBuilder().
		WithName("ArrayAlias").
		WithType(StringSlice())

	clone := original.Clone()

	// Test that clone has the same name
	if clone.GetName() != "ArrayAlias" {
		t.Error("Clone should have the same name")
	}

	// Test that clone has different type reference
	if clone.GetType() == original.GetType() {
		t.Error("Clone should have a different type reference")
	}

	// Test that the cloned type works
	spec := clone.Build()
	if arrayType, ok := spec.Type.(*ast.ArrayType); ok {
		if ident, ok := arrayType.Elt.(*ast.Ident); ok {
			if ident.Name != "string" {
				t.Errorf("Expected array element type 'string', got %s", ident.Name)
			}
		}
	}
}

func TestTypeAliasBuilder_HelperFunctions(t *testing.T) {
	// Test StringSliceAlias
	stringSliceAlias := StringSliceAlias("StringList")
	spec := stringSliceAlias.Build()

	if spec.Name.Name != "StringList" {
		t.Errorf("Expected name 'StringList', got %s", spec.Name.Name)
	}

	if arrayType, ok := spec.Type.(*ast.ArrayType); ok {
		if ident, ok := arrayType.Elt.(*ast.Ident); ok {
			if ident.Name != "string" {
				t.Errorf("Expected array element type 'string', got %s", ident.Name)
			}
		}
	}

	// Test IntSliceAlias
	intSliceAlias := IntSliceAlias("IntList")
	spec = intSliceAlias.Build()

	if arrayType, ok := spec.Type.(*ast.ArrayType); ok {
		if ident, ok := arrayType.Elt.(*ast.Ident); ok {
			if ident.Name != "int" {
				t.Errorf("Expected array element type 'int', got %s", ident.Name)
			}
		}
	}

	// Test BoolSliceAlias
	boolSliceAlias := BoolSliceAlias("BoolList")
	spec = boolSliceAlias.Build()

	if arrayType, ok := spec.Type.(*ast.ArrayType); ok {
		if ident, ok := arrayType.Elt.(*ast.Ident); ok {
			if ident.Name != "bool" {
				t.Errorf("Expected array element type 'bool', got %s", ident.Name)
			}
		}
	}

	// Test IdentAlias
	identAlias := IdentAlias("MyInt", "int")
	spec = identAlias.Build()

	if spec.Name.Name != "MyInt" {
		t.Errorf("Expected name 'MyInt', got %s", spec.Name.Name)
	}

	if ident, ok := spec.Type.(*ast.Ident); ok {
		if ident.Name != "int" {
			t.Errorf("Expected type 'int', got %s", ident.Name)
		}
	}

	// Test SelectorAlias
	selectorAlias := SelectorAlias("MyContext", "context", "Context")
	spec = selectorAlias.Build()

	if spec.Name.Name != "MyContext" {
		t.Errorf("Expected name 'MyContext', got %s", spec.Name.Name)
	}

	if selector, ok := spec.Type.(*ast.SelectorExpr); ok {
		if selector.Sel.Name != "Context" {
			t.Errorf("Expected selector name 'Context', got %s", selector.Sel.Name)
		}
	}

	// Test PointerAlias
	pointerAlias := PointerAlias("MyStringPtr", "string")
	spec = pointerAlias.Build()

	if spec.Name.Name != "MyStringPtr" {
		t.Errorf("Expected name 'MyStringPtr', got %s", spec.Name.Name)
	}

	if starExpr, ok := spec.Type.(*ast.StarExpr); ok {
		if ident, ok := starExpr.X.(*ast.Ident); ok {
			if ident.Name != "string" {
				t.Errorf("Expected pointer target type 'string', got %s", ident.Name)
			}
		}
	}

	// Test SliceAlias
	sliceAlias := SliceAlias("MyStringSlice", "string")
	spec = sliceAlias.Build()

	if spec.Name.Name != "MyStringSlice" {
		t.Errorf("Expected name 'MyStringSlice', got %s", spec.Name.Name)
	}

	if arrayType, ok := spec.Type.(*ast.ArrayType); ok {
		if ident, ok := arrayType.Elt.(*ast.Ident); ok {
			if ident.Name != "string" {
				t.Errorf("Expected slice element type 'string', got %s", ident.Name)
			}
		}
	}

	// Test CustomAlias
	customAlias := CustomAlias("MyCustom", StringSlice())
	spec = customAlias.Build()

	if spec.Name.Name != "MyCustom" {
		t.Errorf("Expected name 'MyCustom', got %s", spec.Name.Name)
	}

	if arrayType, ok := spec.Type.(*ast.ArrayType); ok {
		if ident, ok := arrayType.Elt.(*ast.Ident); ok {
			if ident.Name != "string" {
				t.Errorf("Expected custom type element 'string', got %s", ident.Name)
			}
		}
	}
}

func TestTypeAliasBuilder_ComplexTypes(t *testing.T) {
	// Test nested arrays: [][]string
	nestedArray := SliceOf(StringSlice())
	alias := CustomAlias("StringMatrix", nestedArray)

	spec := alias.Build()

	if outerArrayType, ok := spec.Type.(*ast.ArrayType); ok {
		if innerArrayType, ok := outerArrayType.Elt.(*ast.ArrayType); ok {
			if ident, ok := innerArrayType.Elt.(*ast.Ident); ok {
				if ident.Name != "string" {
					t.Errorf("Expected inner array element type 'string', got %s", ident.Name)
				}
			} else {
				t.Error("Inner array element should be ast.Ident")
			}
		} else {
			t.Error("Inner element should be ast.ArrayType")
		}
	} else {
		t.Error("Outer type should be ast.ArrayType")
	}
}

func TestTypeAliasBuilder_MethodChaining(t *testing.T) {
	builder := NewTypeAliasBuilder().
		WithName("ChainedAlias").
		WithType(String())

	if !builder.HasName() || !builder.HasType() {
		t.Error("Method chaining should work correctly")
	}

	// Test that all operations return the same builder
	operations := []func() *TypeAliasBuilder{
		func() *TypeAliasBuilder { return builder.WithName("NewName") },
		func() *TypeAliasBuilder { return builder.WithType(Int()) },
	}

	for i, op := range operations {
		if op() != builder {
			t.Errorf("Operation %d should return the same builder", i)
		}
	}
}
