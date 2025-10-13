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

	if len(builder.jsonTags) != 0 {
		t.Errorf("Expected empty jsonTags initially, got %v", builder.jsonTags)
	}

	if len(builder.validateTags) != 0 {
		t.Errorf("Expected empty validateTags initially, got %v", builder.validateTags)
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

func TestFieldBuilder_AddJSONTag(t *testing.T) {
	builder := NewFieldBuilder()

	result := builder.AddJSONTag("name")
	if result != builder {
		t.Error("AddJSONTag should return the builder for chaining")
	}

	if len(builder.jsonTags) != 1 || builder.jsonTags[0] != "name" {
		t.Errorf("Expected jsonTags ['name'], got %v", builder.jsonTags)
	}

	// Test adding another tag
	builder.AddJSONTag("omitempty")
	if len(builder.jsonTags) != 2 || builder.jsonTags[1] != "omitempty" {
		t.Errorf("Expected jsonTags ['name', 'omitempty'], got %v", builder.jsonTags)
	}
}

func TestFieldBuilder_AddValidateTag(t *testing.T) {
	builder := NewFieldBuilder()

	result := builder.AddValidateTag("required")
	if result != builder {
		t.Error("AddValidateTag should return the builder for chaining")
	}

	if len(builder.validateTags) != 1 || builder.validateTags[0] != "required" {
		t.Errorf("Expected validateTags ['required'], got %v", builder.validateTags)
	}

	// Test adding another tag
	builder.AddValidateTag("min=1")
	if len(builder.validateTags) != 2 || builder.validateTags[1] != "min=1" {
		t.Errorf("Expected validateTags ['required', 'min=1'], got %v", builder.validateTags)
	}
}

func TestFieldBuilder_Build(t *testing.T) {
	// Test with SimpleTypeBuilder
	builder := NewFieldBuilder().
		WithName("fieldName").
		WithType(String()).
		AddJSONTag("name")

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

	expectedTag := "`json:\"name\"`"
	if field.Tag.Value != expectedTag {
		t.Errorf("Expected tag '%s', got %s", expectedTag, field.Tag.Value)
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

	// Test HasTags
	if builder.HasTags() {
		t.Error("Expected HasTags to return false initially")
	}

	builder.AddJSONTag("test")
	if !builder.HasTags() {
		t.Error("Expected HasTags to return true after adding tag")
	}

	// Test HasJSONTags
	if !builder.HasJSONTags() {
		t.Error("Expected HasJSONTags to return true after adding JSON tag")
	}

	// Test HasValidateTags
	if builder.HasValidateTags() {
		t.Error("Expected HasValidateTags to return false initially")
	}

	builder.AddValidateTag("required")
	if !builder.HasValidateTags() {
		t.Error("Expected HasValidateTags to return true after adding validate tag")
	}

	// Test GetJSONTags
	jsonTags := builder.GetJSONTags()
	if len(jsonTags) != 1 || jsonTags[0] != "test" {
		t.Errorf("Expected GetJSONTags to return ['test'], got %v", jsonTags)
	}

	// Test GetValidateTags
	validateTags := builder.GetValidateTags()
	if len(validateTags) != 1 || validateTags[0] != "required" {
		t.Errorf("Expected GetValidateTags to return ['required'], got %v", validateTags)
	}
}

func TestFieldBuilder_Clone(t *testing.T) {
	builder := NewFieldBuilder().
		WithName("test").
		WithType(String()).
		AddJSONTag("test").
		AddValidateTag("required")

	clone := builder.Clone()

	if clone == builder {
		t.Error("Clone should return a different instance")
	}

	if clone.GetName() != "test" {
		t.Error("Clone should have the same name")
	}

	// Verify tags are cloned
	jsonTags := clone.GetJSONTags()
	if len(jsonTags) != 1 || jsonTags[0] != "test" {
		t.Error("Clone should have the same JSON tags")
	}

	validateTags := clone.GetValidateTags()
	if len(validateTags) != 1 || validateTags[0] != "required" {
		t.Error("Clone should have the same validate tags")
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

func TestFieldBuilder_TagMethods(t *testing.T) {
	builder := NewFieldBuilder()

	// Test AddJSONTags
	builder.AddJSONTags("name", "omitempty")
	if len(builder.jsonTags) != 2 {
		t.Errorf("Expected 2 JSON tags, got %d", len(builder.jsonTags))
	}

	// Test SetJSONTags
	builder.SetJSONTags("id", "required")
	if len(builder.jsonTags) != 2 || builder.jsonTags[0] != "id" || builder.jsonTags[1] != "required" {
		t.Errorf("Expected JSON tags ['id', 'required'], got %v", builder.jsonTags)
	}

	// Test AddValidateTags
	builder.AddValidateTags("required", "min=1")
	if len(builder.validateTags) != 2 {
		t.Errorf("Expected 2 validate tags, got %d", len(builder.validateTags))
	}

	// Test SetValidateTags
	builder.SetValidateTags("max=100")
	if len(builder.validateTags) != 1 || builder.validateTags[0] != "max=100" {
		t.Errorf("Expected validate tags ['max=100'], got %v", builder.validateTags)
	}

	// Test ClearJSONTags
	builder.ClearJSONTags()
	if len(builder.jsonTags) != 0 {
		t.Error("JSON tags should be cleared")
	}

	// Test ClearValidateTags
	builder.ClearValidateTags()
	if len(builder.validateTags) != 0 {
		t.Error("Validate tags should be cleared")
	}

	// Test ClearAllTags
	builder.AddJSONTag("test").AddValidateTag("required")
	builder.ClearAllTags()
	if len(builder.jsonTags) != 0 || len(builder.validateTags) != 0 {
		t.Error("All tags should be cleared")
	}
}

func TestFieldBuilder_BuildWithMultipleTags(t *testing.T) {
	builder := NewFieldBuilder().
		WithName("field").
		WithType(String()).
		AddJSONTags("name", "omitempty").
		AddValidateTags("required", "min=1")

	field := builder.Build()

	if field.Tag == nil {
		t.Fatal("Tag should not be nil")
	}

	expectedTag := "`json:\"name,omitempty\" validate:\"required,min=1\"`"
	if field.Tag.Value != expectedTag {
		t.Errorf("Expected tag '%s', got %s", expectedTag, field.Tag.Value)
	}
}

func TestFieldBuilder_BuildWithOnlyJSONTags(t *testing.T) {
	builder := NewFieldBuilder().
		WithName("field").
		WithType(String()).
		AddJSONTag("name")

	field := builder.Build()

	if field.Tag == nil {
		t.Fatal("Tag should not be nil")
	}

	expectedTag := "`json:\"name\"`"
	if field.Tag.Value != expectedTag {
		t.Errorf("Expected tag '%s', got %s", expectedTag, field.Tag.Value)
	}
}

func TestFieldBuilder_BuildWithOnlyValidateTags(t *testing.T) {
	builder := NewFieldBuilder().
		WithName("field").
		WithType(String()).
		AddValidateTag("required")

	field := builder.Build()

	if field.Tag == nil {
		t.Fatal("Tag should not be nil")
	}

	expectedTag := "`validate:\"required\"`"
	if field.Tag.Value != expectedTag {
		t.Errorf("Expected tag '%s', got %s", expectedTag, field.Tag.Value)
	}
}
