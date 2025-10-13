package astbuilder

import (
	"go/ast"
	"go/token"
	"testing"
)

func TestNewFieldBuilder(t *testing.T) {
	builder := NewFieldBuilder()

	if builder == nil {
		t.Fatal("NewFieldBuilder returned nil")
	}

	if builder.name != "" {
		t.Errorf("Expected empty name initially, got %s", builder.name)
	}

	if builder.typeBuilder == nil {
		t.Fatal("typeBuilder should not be nil")
	}

	if builder.tag != "" {
		t.Errorf("Expected empty tag initially, got %s", builder.tag)
	}
}

func TestFieldBuilder_WithName(t *testing.T) {
	builder := NewFieldBuilder()

	result := builder.WithName("fieldName")
	if result != builder {
		t.Error("WithName should return the builder for chaining")
	}

	if builder.name != "fieldName" {
		t.Errorf("Expected name 'fieldName', got %s", builder.name)
	}
}

func TestFieldBuilder_WithType(t *testing.T) {
	builder := NewFieldBuilder()
	typeBuilder := String()

	result := builder.WithType(typeBuilder)
	if result != builder {
		t.Error("WithType should return the builder for chaining")
	}

	if builder.typeBuilder == nil {
		t.Fatal("typeBuilder should not be nil")
	}

	// Verify it's a clone, not the same instance
	if builder.typeBuilder == typeBuilder {
		t.Error("typeBuilder should be a clone, not the same instance")
	}

	// Test with nil type builder
	defer func() {
		if r := recover(); r == nil {
			t.Error("WithType should panic when type builder is nil")
		}
	}()

	builder.WithType(nil)
}

func TestFieldBuilder_WithTag(t *testing.T) {
	builder := NewFieldBuilder()

	result := builder.WithTag(`json:"name"`)
	if result != builder {
		t.Error("WithTag should return the builder for chaining")
	}

	if builder.tag != `json:"name"` {
		t.Errorf("Expected tag 'json:\"name\"', got %s", builder.tag)
	}
}

func TestFieldBuilder_Build(t *testing.T) {
	// Test with SimpleTypeBuilder
	builder := NewFieldBuilder().
		WithName("fieldName").
		WithType(String()).
		WithTag(`json:"name"`)

	field := builder.Build()

	if field == nil {
		t.Fatal("Build returned nil")
	}

	// Check name
	if len(field.Names) != 1 {
		t.Errorf("Expected 1 name, got %d", len(field.Names))
	}

	if field.Names[0].Name != "fieldName" {
		t.Errorf("Expected name 'fieldName', got %s", field.Names[0].Name)
	}

	// Check type
	if ident, ok := field.Type.(*ast.Ident); ok {
		if ident.Name != "string" {
			t.Errorf("Expected type 'string', got %s", ident.Name)
		}
	} else {
		t.Error("Type should be ast.Ident")
	}

	// Check tag
	if field.Tag == nil {
		t.Fatal("Tag should not be nil")
	}

	if field.Tag.Value != `json:"name"` {
		t.Errorf("Expected tag 'json:\"name\"', got %s", field.Tag.Value)
	}

	if field.Tag.Kind != token.STRING {
		t.Error("Tag should be a string literal")
	}
}

func TestFieldBuilder_BuildWithoutName(t *testing.T) {
	builder := NewFieldBuilder().WithType(String())
	field := builder.Build()

	if len(field.Names) != 0 {
		t.Errorf("Expected 0 names for unnamed field, got %d", len(field.Names))
	}
}

func TestFieldBuilder_BuildWithoutTag(t *testing.T) {
	builder := NewFieldBuilder().WithName("field").WithType(String())
	field := builder.Build()

	if field.Tag != nil {
		t.Error("Tag should be nil when not set")
	}
}

func TestFieldBuilder_BuildWithoutType(t *testing.T) {
	builder := NewFieldBuilder().WithName("field")

	defer func() {
		if r := recover(); r == nil {
			t.Error("Build should panic when no type is set")
		}
	}()

	builder.Build()
}

func TestFieldBuilder_HelperMethods(t *testing.T) {
	// Test StringField
	builder := StringField("name")
	field := builder.Build()

	if len(field.Names) != 1 || field.Names[0].Name != "name" {
		t.Error("StringField should set the name correctly")
	}

	if ident, ok := field.Type.(*ast.Ident); !ok || ident.Name != "string" {
		t.Error("StringField should set type to string")
	}

	// Test IntField
	builder = IntField("count")
	field = builder.Build()

	if len(field.Names) != 1 || field.Names[0].Name != "count" {
		t.Error("IntField should set the name correctly")
	}

	if ident, ok := field.Type.(*ast.Ident); !ok || ident.Name != "int" {
		t.Error("IntField should set type to int")
	}

	// Test ErrorField
	builder = ErrorField()
	field = builder.Build()

	if len(field.Names) != 0 {
		t.Error("ErrorField should not have a name")
	}

	if ident, ok := field.Type.(*ast.Ident); !ok || ident.Name != "error" {
		t.Error("ErrorField should set type to error")
	}

	// Test ContextField
	builder = ContextField("ctx")
	field = builder.Build()

	if len(field.Names) != 1 || field.Names[0].Name != "ctx" {
		t.Error("ContextField should set the name correctly")
	}

	if selector, ok := field.Type.(*ast.SelectorExpr); !ok || selector.Sel.Name != "Context" {
		t.Error("ContextField should set type to context.Context")
	}

	// Test IdentField
	builder = IdentField("id", "UserID")
	field = builder.Build()

	if len(field.Names) != 1 || field.Names[0].Name != "id" {
		t.Error("IdentField should set the name correctly")
	}

	if ident, ok := field.Type.(*ast.Ident); !ok || ident.Name != "UserID" {
		t.Error("IdentField should set type to UserID")
	}

	// Test SelectorField
	builder = SelectorField("req", "models", "Request")
	field = builder.Build()

	if len(field.Names) != 1 || field.Names[0].Name != "req" {
		t.Error("SelectorField should set the name correctly")
	}

	if selector, ok := field.Type.(*ast.SelectorExpr); !ok || selector.Sel.Name != "Request" {
		t.Error("SelectorField should set type to models.Request")
	}

	// Test with a more complex type using SimpleTypeBuilder
	builder = NewFieldBuilder().
		WithName("complex").
		WithType(Selector("package", "ComplexType"))
	field = builder.Build()

	if len(field.Names) != 1 || field.Names[0].Name != "complex" {
		t.Error("Complex field should set the name correctly")
	}

	if selector, ok := field.Type.(*ast.SelectorExpr); !ok || selector.Sel.Name != "ComplexType" {
		t.Error("Complex field should create a selector type")
	}

}

func TestFieldBuilder_UtilityMethods(t *testing.T) {
	builder := NewFieldBuilder()

	// Test HasName
	if builder.HasName() {
		t.Error("Expected HasName to return false initially")
	}

	builder.WithName("test")
	if !builder.HasName() {
		t.Error("Expected HasName to return true after setting name")
	}

	// Test GetName
	if builder.GetName() != "test" {
		t.Errorf("Expected GetName to return 'test', got %s", builder.GetName())
	}

	// Test HasTag
	if builder.HasTag() {
		t.Error("Expected HasTag to return false initially")
	}

	builder.WithTag("test")
	if !builder.HasTag() {
		t.Error("Expected HasTag to return true after setting tag")
	}

	// Test GetTag
	if builder.GetTag() != "test" {
		t.Errorf("Expected GetTag to return 'test', got %s", builder.GetTag())
	}
}

func TestFieldBuilder_Clone(t *testing.T) {
	builder := NewFieldBuilder().
		WithName("test").
		WithType(String()).
		WithTag("test")

	clone := builder.Clone()

	if clone == builder {
		t.Error("Clone should return a different instance")
	}

	if clone.GetName() != "test" {
		t.Error("Clone should have the same name")
	}

	if clone.GetTag() != "test" {
		t.Error("Clone should have the same tag")
	}

	// Verify type builder is cloned
	if clone.typeBuilder == builder.typeBuilder {
		t.Error("Clone should have a different type builder instance")
	}

	// Modify original
	builder.WithName("modified")

	// Clone should be unaffected
	if clone.GetName() != "test" {
		t.Error("Clone should be unaffected by original modifications")
	}
}
