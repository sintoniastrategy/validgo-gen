package astbuilder

import (
	"go/ast"
	"go/token"
	"testing"
)

func TestNewStructBuilder(t *testing.T) {
	builder := NewStructBuilder()

	if builder == nil {
		t.Fatal("NewStructBuilder returned nil")
	}

	if builder.name != "" {
		t.Errorf("Expected empty name initially, got %s", builder.name)
	}

	if len(builder.fields) != 0 {
		t.Errorf("Expected empty fields initially, got %d fields", len(builder.fields))
	}
}

func TestStructBuilder_WithName(t *testing.T) {
	builder := NewStructBuilder()

	result := builder.WithName("TestStruct")
	if result != builder {
		t.Error("WithName should return the builder for chaining")
	}

	if builder.name != "TestStruct" {
		t.Errorf("Expected name 'TestStruct', got %s", builder.name)
	}

	// Test chaining
	builder.WithName("AnotherStruct").WithName("FinalStruct")
	if builder.name != "FinalStruct" {
		t.Errorf("Expected chained name 'FinalStruct', got %s", builder.name)
	}
}

func TestStructBuilder_AddField(t *testing.T) {
	builder := NewStructBuilder()

	fieldBuilder := StringField("name")
	result := builder.AddField(fieldBuilder)

	if result != builder {
		t.Error("AddField should return the builder for chaining")
	}

	if len(builder.fields) != 1 {
		t.Errorf("Expected 1 field, got %d", len(builder.fields))
	}

	// Test that the field was cloned (not the same reference)
	if builder.fields[0] == fieldBuilder {
		t.Error("Field should be cloned, not the same reference")
	}

	// Test adding another field
	builder.AddField(IntField("age"))
	if len(builder.fields) != 2 {
		t.Errorf("Expected 2 fields, got %d", len(builder.fields))
	}
}

func TestStructBuilder_AddFieldNil(t *testing.T) {
	builder := NewStructBuilder()

	defer func() {
		if r := recover(); r == nil {
			t.Error("AddField should panic when field builder is nil")
		}
	}()

	builder.AddField(nil)
}

func TestStructBuilder_AddFields(t *testing.T) {
	builder := NewStructBuilder()

	field1 := StringField("name")
	field2 := IntField("age")
	field3 := BoolField("active")

	result := builder.AddFields(field1, field2, field3)

	if result != builder {
		t.Error("AddFields should return the builder for chaining")
	}

	if len(builder.fields) != 3 {
		t.Errorf("Expected 3 fields, got %d", len(builder.fields))
	}

	// Test that fields were cloned
	if builder.fields[0] == field1 {
		t.Error("First field should be cloned")
	}
	if builder.fields[1] == field2 {
		t.Error("Second field should be cloned")
	}
	if builder.fields[2] == field3 {
		t.Error("Third field should be cloned")
	}
}

func TestStructBuilder_AddFieldsNil(t *testing.T) {
	builder := NewStructBuilder()

	defer func() {
		if r := recover(); r == nil {
			t.Error("AddFields should panic when any field builder is nil")
		}
	}()

	builder.AddFields(StringField("name"), nil, IntField("age"))
}

func TestStructBuilder_Build(t *testing.T) {
	builder := NewStructBuilder().
		WithName("Person").
		AddField(StringField("name")).
		AddField(IntField("age"))

	typeSpec := builder.Build()

	if typeSpec.Name.Name != "Person" {
		t.Errorf("Expected type name 'Person', got %s", typeSpec.Name.Name)
	}

	structType, ok := typeSpec.Type.(*ast.StructType)
	if !ok {
		t.Fatal("Type should be *ast.StructType")
	}

	if len(structType.Fields.List) != 2 {
		t.Errorf("Expected 2 fields, got %d", len(structType.Fields.List))
	}

	// Check first field
	field1 := structType.Fields.List[0]
	if len(field1.Names) != 1 || field1.Names[0].Name != "name" {
		t.Error("First field should be named 'name'")
	}
	if ident, ok := field1.Type.(*ast.Ident); !ok || ident.Name != "string" {
		t.Error("First field should be of type 'string'")
	}

	// Check second field
	field2 := structType.Fields.List[1]
	if len(field2.Names) != 1 || field2.Names[0].Name != "age" {
		t.Error("Second field should be named 'age'")
	}
	if ident, ok := field2.Type.(*ast.Ident); !ok || ident.Name != "int" {
		t.Error("Second field should be of type 'int'")
	}
}

func TestStructBuilder_BuildWithoutName(t *testing.T) {
	builder := NewStructBuilder().AddField(StringField("name"))

	defer func() {
		if r := recover(); r == nil {
			t.Error("Build should panic when struct has no name")
		}
	}()

	builder.Build()
}

func TestStructBuilder_BuildAsDeclaration(t *testing.T) {
	builder := NewStructBuilder().
		WithName("Person").
		AddField(StringField("name"))

	decl := builder.BuildAsDeclaration()

	if decl.Tok != token.TYPE {
		t.Error("Declaration should have token.TYPE")
	}

	if len(decl.Specs) != 1 {
		t.Errorf("Expected 1 spec, got %d", len(decl.Specs))
	}

	typeSpec, ok := decl.Specs[0].(*ast.TypeSpec)
	if !ok {
		t.Fatal("Spec should be *ast.TypeSpec")
	}

	if typeSpec.Name.Name != "Person" {
		t.Errorf("Expected type name 'Person', got %s", typeSpec.Name.Name)
	}
}

func TestStructBuilder_UtilityMethods(t *testing.T) {
	builder := NewStructBuilder()

	// Test HasName
	if builder.HasName() {
		t.Error("HasName should return false initially")
	}

	builder.WithName("TestStruct")
	if !builder.HasName() {
		t.Error("HasName should return true after setting name")
	}

	// Test GetName
	if builder.GetName() != "TestStruct" {
		t.Errorf("Expected name 'TestStruct', got %s", builder.GetName())
	}

	// Test FieldCount
	if builder.FieldCount() != 0 {
		t.Errorf("Expected 0 fields initially, got %d", builder.FieldCount())
	}

	// Test HasFields
	if builder.HasFields() {
		t.Error("HasFields should return false initially")
	}

	builder.AddField(StringField("name"))
	if !builder.HasFields() {
		t.Error("HasFields should return true after adding field")
	}

	if builder.FieldCount() != 1 {
		t.Errorf("Expected 1 field, got %d", builder.FieldCount())
	}
}

func TestStructBuilder_GetFields(t *testing.T) {
	builder := NewStructBuilder().
		AddField(StringField("name")).
		AddField(IntField("age"))

	fields := builder.GetFields()

	if len(fields) != 2 {
		t.Errorf("Expected 2 fields, got %d", len(fields))
	}

	// Test that returned fields are clones
	if fields[0] == builder.fields[0] {
		t.Error("Returned fields should be clones")
	}

	// Test that modifying returned fields doesn't affect builder
	fields[0].WithName("modified")
	if builder.fields[0].GetName() == "modified" {
		t.Error("Modifying returned fields should not affect builder")
	}
}

func TestStructBuilder_GetField(t *testing.T) {
	builder := NewStructBuilder().
		AddField(StringField("name")).
		AddField(IntField("age"))

	// Test valid index
	field := builder.GetField(0)
	if field == nil {
		t.Fatal("GetField(0) should not return nil")
	}
	if field.GetName() != "name" {
		t.Errorf("Expected field name 'name', got %s", field.GetName())
	}

	// Test invalid index
	field = builder.GetField(-1)
	if field != nil {
		t.Error("GetField(-1) should return nil")
	}

	field = builder.GetField(2)
	if field != nil {
		t.Error("GetField(2) should return nil")
	}
}

func TestStructBuilder_GetFieldByName(t *testing.T) {
	builder := NewStructBuilder().
		AddField(StringField("name")).
		AddField(IntField("age"))

	// Test existing field
	field := builder.GetFieldByName("name")
	if field == nil {
		t.Fatal("GetFieldByName('name') should not return nil")
	}
	if field.GetName() != "name" {
		t.Errorf("Expected field name 'name', got %s", field.GetName())
	}

	// Test non-existing field
	field = builder.GetFieldByName("nonexistent")
	if field != nil {
		t.Error("GetFieldByName('nonexistent') should return nil")
	}
}

func TestStructBuilder_RemoveField(t *testing.T) {
	builder := NewStructBuilder().
		AddField(StringField("name")).
		AddField(IntField("age")).
		AddField(BoolField("active"))

	// Test removing middle field
	result := builder.RemoveField(1)
	if result != builder {
		t.Error("RemoveField should return the builder for chaining")
	}

	if builder.FieldCount() != 2 {
		t.Errorf("Expected 2 fields after removal, got %d", builder.FieldCount())
	}

	// Check that the correct field was removed
	if builder.fields[0].GetName() != "name" {
		t.Error("First field should still be 'name'")
	}
	if builder.fields[1].GetName() != "active" {
		t.Error("Second field should now be 'active'")
	}

	// Test invalid index
	builder.RemoveField(-1)
	if builder.FieldCount() != 2 {
		t.Error("RemoveField with invalid index should not affect field count")
	}

	builder.RemoveField(10)
	if builder.FieldCount() != 2 {
		t.Error("RemoveField with invalid index should not affect field count")
	}
}

func TestStructBuilder_RemoveFieldByName(t *testing.T) {
	builder := NewStructBuilder().
		AddField(StringField("name")).
		AddField(IntField("age")).
		AddField(BoolField("active"))

	// Test removing existing field
	result := builder.RemoveFieldByName("age")
	if result != builder {
		t.Error("RemoveFieldByName should return the builder for chaining")
	}

	if builder.FieldCount() != 2 {
		t.Errorf("Expected 2 fields after removal, got %d", builder.FieldCount())
	}

	// Check that the correct field was removed
	if builder.fields[0].GetName() != "name" {
		t.Error("First field should still be 'name'")
	}
	if builder.fields[1].GetName() != "active" {
		t.Error("Second field should now be 'active'")
	}

	// Test removing non-existing field
	builder.RemoveFieldByName("nonexistent")
	if builder.FieldCount() != 2 {
		t.Error("RemoveFieldByName with non-existing field should not affect field count")
	}
}

func TestStructBuilder_Clear(t *testing.T) {
	builder := NewStructBuilder().
		AddField(StringField("name")).
		AddField(IntField("age"))

	result := builder.Clear()
	if result != builder {
		t.Error("Clear should return the builder for chaining")
	}

	if builder.FieldCount() != 0 {
		t.Errorf("Expected 0 fields after clear, got %d", builder.FieldCount())
	}

	if builder.HasFields() {
		t.Error("HasFields should return false after clear")
	}
}

func TestStructBuilder_Clone(t *testing.T) {
	original := NewStructBuilder().
		WithName("Person").
		AddField(StringField("name")).
		AddField(IntField("age"))

	clone := original.Clone()

	// Test that clone has the same values
	if clone.GetName() != original.GetName() {
		t.Error("Clone should have the same name")
	}

	if clone.FieldCount() != original.FieldCount() {
		t.Error("Clone should have the same number of fields")
	}

	// Test that clone has different field references
	for i := range original.fields {
		if clone.fields[i] == original.fields[i] {
			t.Error("Clone should have different field references")
		}
	}

	// Test that modifying clone doesn't affect original
	clone.WithName("ModifiedPerson")
	clone.AddField(BoolField("active"))

	if original.GetName() == "ModifiedPerson" {
		t.Error("Modifying clone name should not affect original")
	}

	if original.FieldCount() == 3 {
		t.Error("Modifying clone fields should not affect original")
	}
}

func TestStructBuilder_HelperMethods(t *testing.T) {
	builder := NewStructBuilder()

	// Test AddStringField
	result := builder.AddStringField("name")
	if result != builder {
		t.Error("AddStringField should return the builder for chaining")
	}

	// Test AddIntField
	builder.AddIntField("age")

	// Test AddBoolField
	builder.AddBoolField("active")

	// Test AddContextField
	builder.AddContextField("ctx")

	// Test AddIdentField
	builder.AddIdentField("customType", "CustomType")

	// Test AddSelectorField
	builder.AddSelectorField("response", "apimodels", "Response")

	if builder.FieldCount() != 6 {
		t.Errorf("Expected 6 fields, got %d", builder.FieldCount())
	}

	// Verify field types
	fields := builder.GetFields()
	if fields[0].GetName() != "name" {
		t.Error("First field should be named 'name'")
	}
	if fields[1].GetName() != "age" {
		t.Error("Second field should be named 'age'")
	}
	if fields[2].GetName() != "active" {
		t.Error("Third field should be named 'active'")
	}
	if fields[3].GetName() != "ctx" {
		t.Error("Fourth field should be named 'ctx'")
	}
	if fields[4].GetName() != "customType" {
		t.Error("Fifth field should be named 'customType'")
	}
	if fields[5].GetName() != "response" {
		t.Error("Sixth field should be named 'response'")
	}
}

func TestStructBuilder_MethodChaining(t *testing.T) {
	builder := NewStructBuilder().
		WithName("Person").
		AddField(StringField("name")).
		AddField(IntField("age")).
		AddFields(BoolField("active"), ContextField("ctx")).
		AddStringField("email").
		AddIntField("score")

	if builder.GetName() != "Person" {
		t.Errorf("Expected name 'Person', got %s", builder.GetName())
	}

	if builder.FieldCount() != 6 {
		t.Errorf("Expected 6 fields, got %d", builder.FieldCount())
	}

	// Test that all operations return the same builder
	operations := []func() *StructBuilder{
		func() *StructBuilder { return builder.WithName("Test") },
		func() *StructBuilder { return builder.AddField(StringField("test")) },
		func() *StructBuilder { return builder.AddFields(IntField("test2")) },
		func() *StructBuilder { return builder.AddStringField("test3") },
		func() *StructBuilder { return builder.Clear() },
	}

	for i, op := range operations {
		if op() != builder {
			t.Errorf("Operation %d should return the same builder", i)
		}
	}
}

func TestStructBuilder_ComplexExample(t *testing.T) {
	// Create a complex struct similar to what might be generated
	builder := NewStructBuilder().
		WithName("CreateRequest").
		AddStringField("name").
		AddStringField("description").
		AddIntField("priority").
		AddBoolField("active").
		AddContextField("ctx").
		AddSelectorField("metadata", "apimodels", "Metadata").
		AddField(NewFieldBuilder().
			WithName("tags").
			WithType(String()).
			AddJSONTags("tags", "omitempty"))

	typeSpec := builder.Build()

	if typeSpec.Name.Name != "CreateRequest" {
		t.Errorf("Expected type name 'CreateRequest', got %s", typeSpec.Name.Name)
	}

	structType := typeSpec.Type.(*ast.StructType)
	if len(structType.Fields.List) != 7 {
		t.Errorf("Expected 7 fields, got %d", len(structType.Fields.List))
	}

	// Verify specific fields
	fields := structType.Fields.List
	if fields[0].Names[0].Name != "name" {
		t.Error("First field should be named 'name'")
	}
	if fields[5].Names[0].Name != "metadata" {
		t.Error("Sixth field should be named 'metadata'")
	}
	if fields[6].Names[0].Name != "tags" {
		t.Error("Seventh field should be named 'tags'")
	}
}
