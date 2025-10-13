package astbuilder

import (
	"go/ast"
	"go/token"
	"testing"
)

func TestNewInterfaceBuilder(t *testing.T) {
	builder := NewInterfaceBuilder()

	if builder == nil {
		t.Fatal("NewInterfaceBuilder returned nil")
	}

	if builder.name != "" {
		t.Errorf("Expected empty name initially, got %s", builder.name)
	}

	if builder.methods == nil {
		t.Fatal("methods slice is nil")
	}

	if len(builder.methods) != 0 {
		t.Errorf("Expected empty methods initially, got %d methods", len(builder.methods))
	}
}

func TestNewInterfaceMethodBuilder(t *testing.T) {
	methodBuilder := NewInterfaceMethodBuilder()

	if methodBuilder == nil {
		t.Fatal("NewInterfaceMethodBuilder returned nil")
	}

	if methodBuilder.name != "" {
		t.Errorf("Expected empty name initially, got %s", methodBuilder.name)
	}

	if methodBuilder.interfaceBuilder != nil {
		t.Error("Expected nil interface builder reference initially")
	}

	if methodBuilder.params == nil {
		t.Fatal("params slice is nil")
	}

	if methodBuilder.results == nil {
		t.Fatal("results slice is nil")
	}

	if len(methodBuilder.params) != 0 {
		t.Errorf("Expected empty params initially, got %d params", len(methodBuilder.params))
	}

	if len(methodBuilder.results) != 0 {
		t.Errorf("Expected empty results initially, got %d results", len(methodBuilder.results))
	}
}

func TestInterfaceBuilder_WithName(t *testing.T) {
	builder := NewInterfaceBuilder()

	result := builder.WithName("TestInterface")
	if result != builder {
		t.Error("WithName should return the builder for chaining")
	}

	if builder.name != "TestInterface" {
		t.Errorf("Expected name 'TestInterface', got %s", builder.name)
	}
}

func TestInterfaceBuilder_WithMethod(t *testing.T) {
	builder := NewInterfaceBuilder()

	// Create a method builder independently
	methodBuilder := NewInterfaceMethodBuilder().
		WithName("TestMethod").
		AddArgField(StringField("param")).
		AddRetvalField(ErrorField())

	// Add the method to the interface
	result := builder.WithMethod(methodBuilder)
	if result != builder {
		t.Error("WithMethod should return the builder for chaining")
	}

	// Check that the method was added
	if builder.MethodCount() != 1 {
		t.Errorf("Expected 1 method, got %d", builder.MethodCount())
	}

	// Check that the method builder reference was set
	if methodBuilder.interfaceBuilder != builder {
		t.Error("method builder should reference the interface builder")
	}

	// Check the method details
	addedMethod := builder.GetMethod(0)
	if addedMethod != methodBuilder {
		t.Error("Added method should be the same instance")
	}

	if addedMethod.name != "TestMethod" {
		t.Errorf("Expected method name 'TestMethod', got %s", addedMethod.name)
	}
}

func TestInterfaceBuilder_WithMethodNil(t *testing.T) {
	builder := NewInterfaceBuilder()

	defer func() {
		if r := recover(); r == nil {
			t.Error("WithMethod should panic when method builder is nil")
		}
	}()

	builder.WithMethod(nil)
}

func TestInterfaceBuilder_WithMethodExistingReference(t *testing.T) {
	builder1 := NewInterfaceBuilder()
	builder2 := NewInterfaceBuilder()

	// Create a method builder and add to first builder
	methodBuilder := NewInterfaceMethodBuilder().WithName("TestMethod")
	builder1.WithMethod(methodBuilder)

	// Add to second builder - should work and update reference
	result := builder2.WithMethod(methodBuilder)
	if result != builder2 {
		t.Error("WithMethod should return the builder for chaining")
	}

	// The method builder should now reference builder2
	if methodBuilder.interfaceBuilder != builder2 {
		t.Error("method builder should reference builder2")
	}

	if builder2.MethodCount() != 1 {
		t.Errorf("Expected 1 method in builder2, got %d", builder2.MethodCount())
	}
}

func TestInterfaceMethodBuilder_WithName(t *testing.T) {
	methodBuilder := NewInterfaceMethodBuilder()

	result := methodBuilder.WithName("TestMethod")
	if result != methodBuilder {
		t.Error("WithName should return the method builder for chaining")
	}

	if methodBuilder.name != "TestMethod" {
		t.Errorf("Expected method name 'TestMethod', got %s", methodBuilder.name)
	}
}

func TestInterfaceMethodBuilder_AddArgField(t *testing.T) {
	methodBuilder := NewInterfaceMethodBuilder()

	// Test adding named parameter
	result := methodBuilder.AddArgField(StringField("param1"))
	if result != methodBuilder {
		t.Error("AddArgField should return the method builder for chaining")
	}

	if len(methodBuilder.params) != 1 {
		t.Errorf("Expected 1 parameter, got %d", len(methodBuilder.params))
	}

	param := methodBuilder.params[0]
	if !param.HasName() {
		t.Error("Expected parameter to have a name")
	}

	if param.GetName() != "param1" {
		t.Errorf("Expected parameter name 'param1', got %s", param.GetName())
	}

	// Test adding unnamed parameter
	methodBuilder.AddArgField(IntField(""))
	if len(methodBuilder.params) != 2 {
		t.Errorf("Expected 2 parameters, got %d", len(methodBuilder.params))
	}

	unnamedParam := methodBuilder.params[1]
	if unnamedParam.HasName() {
		t.Error("Expected unnamed parameter to not have a name")
	}
}

func TestInterfaceMethodBuilder_AddRetvalField(t *testing.T) {
	methodBuilder := NewInterfaceMethodBuilder()

	// Test adding named return value
	result := methodBuilder.AddRetvalField(StringField("result1"))
	if result != methodBuilder {
		t.Error("AddRetvalField should return the method builder for chaining")
	}

	if len(methodBuilder.results) != 1 {
		t.Errorf("Expected 1 return value, got %d", len(methodBuilder.results))
	}

	retval := methodBuilder.results[0]
	if !retval.HasName() {
		t.Error("Expected return value to have a name")
	}

	if retval.GetName() != "result1" {
		t.Errorf("Expected return value name 'result1', got %s", retval.GetName())
	}

	// Test adding unnamed return value
	methodBuilder.AddRetvalField(ErrorField())
	if len(methodBuilder.results) != 2 {
		t.Errorf("Expected 2 return values, got %d", len(methodBuilder.results))
	}

	unnamedRetval := methodBuilder.results[1]
	if unnamedRetval.HasName() {
		t.Error("Expected unnamed return value to not have a name")
	}
}

func TestInterfaceMethodBuilder_WithMethod(t *testing.T) {
	builder := NewInterfaceBuilder()
	methodBuilder := NewInterfaceMethodBuilder().
		WithName("TestMethod").
		AddArgField(ContextField("ctx")).
		AddArgField(IdentField("req", "Request")).
		AddRetvalField(IdentField("", "Response")).
		AddRetvalField(ErrorField())

	result := builder.WithMethod(methodBuilder)
	if result != builder {
		t.Error("WithMethod should return the interface builder")
	}

	if len(builder.methods) != 1 {
		t.Errorf("Expected 1 method in interface, got %d", len(builder.methods))
	}

	method := builder.methods[0]
	if method.name != "TestMethod" {
		t.Errorf("Expected method name 'TestMethod', got %s", method.name)
	}

	if len(method.params) != 2 {
		t.Errorf("Expected 2 parameters, got %d", len(method.params))
	}

	if len(method.results) != 2 {
		t.Errorf("Expected 2 return values, got %d", len(method.results))
	}
}

func TestInterfaceBuilder_Build(t *testing.T) {
	method1 := NewInterfaceMethodBuilder().
		WithName("Method1").
		AddArgField(StringField("param")).
		AddRetvalField(ErrorField())

	method2 := NewInterfaceMethodBuilder().
		WithName("Method2").
		AddRetvalField(IntField("result"))

	builder := NewInterfaceBuilder().
		WithName("TestInterface").
		WithMethod(method1).
		WithMethod(method2)

	decl := builder.Build()

	if decl == nil {
		t.Fatal("Build returned nil")
	}

	if decl.Tok != token.TYPE {
		t.Errorf("Expected token.TYPE, got %v", decl.Tok)
	}

	if len(decl.Specs) != 1 {
		t.Errorf("Expected 1 spec, got %d", len(decl.Specs))
	}

	typeSpec := decl.Specs[0].(*ast.TypeSpec)
	if typeSpec.Name.Name != "TestInterface" {
		t.Errorf("Expected type name 'TestInterface', got %s", typeSpec.Name.Name)
	}

	interfaceType := typeSpec.Type.(*ast.InterfaceType)
	if len(interfaceType.Methods.List) != 2 {
		t.Errorf("Expected 2 methods, got %d", len(interfaceType.Methods.List))
	}
}

func TestInterfaceBuilder_BuildAsDeclaration(t *testing.T) {
	methodBuilder := NewInterfaceMethodBuilder().
		WithName("TestMethod")

	builder := NewInterfaceBuilder().
		WithName("TestInterface").
		WithMethod(methodBuilder)

	decl := builder.BuildAsDeclaration()

	if decl == nil {
		t.Fatal("BuildAsDeclaration returned nil")
	}

	if _, ok := decl.(*ast.GenDecl); !ok {
		t.Error("BuildAsDeclaration should return *ast.GenDecl")
	}
}

func TestInterfaceBuilder_BuildWithoutName(t *testing.T) {
	builder := NewInterfaceBuilder()

	defer func() {
		if r := recover(); r == nil {
			t.Error("Build should panic when interface name is not set")
		}
	}()

	builder.Build()
}

func TestInterfaceBuilder_WithMethodWithoutName(t *testing.T) {
	builder := NewInterfaceBuilder()
	methodBuilder := NewInterfaceMethodBuilder()
	// Don't set the method name

	// This should not panic since we're not calling BuildMethod anymore
	// The method will be added but won't have a name
	builder.WithMethod(methodBuilder)

	// The method should still be added, but with empty name
	if builder.MethodCount() != 1 {
		t.Errorf("Expected 1 method, got %d", builder.MethodCount())
	}

	method := builder.GetMethod(0)
	if method.name != "" {
		t.Errorf("Expected empty method name, got %s", method.name)
	}
}

func TestInterfaceMethodBuilder_HelperMethods(t *testing.T) {
	builder := NewInterfaceBuilder()
	methodBuilder := NewInterfaceMethodBuilder().
		WithName("TestMethod").
		AddContextArg().
		AddErrorRetval()

	builder.WithMethod(methodBuilder)

	method := builder.methods[0]

	// Check context parameter
	if len(method.params) != 1 {
		t.Errorf("Expected 1 parameter, got %d", len(method.params))
	}

	contextParam := method.params[0]
	if !contextParam.HasName() {
		t.Error("Expected context parameter to have a name")
	}

	if contextParam.GetName() != "ctx" {
		t.Errorf("Expected context parameter name 'ctx', got %s", contextParam.GetName())
	}

	// Check error return value
	if len(method.results) != 1 {
		t.Errorf("Expected 1 return value, got %d", len(method.results))
	}

	errorResult := method.results[0]
	if errorResult.HasName() {
		t.Error("Expected error return value to not have a name")
	}
}

func TestInterfaceBuilder_UtilityMethods(t *testing.T) {
	builder := NewInterfaceBuilder().
		WithName("TestInterface")

	// Test MethodCount
	if builder.MethodCount() != 0 {
		t.Errorf("Expected 0 methods initially, got %d", builder.MethodCount())
	}

	// Test HasMethods
	if builder.HasMethods() {
		t.Error("Expected HasMethods to return false initially")
	}

	// Add a method
	methodBuilder := NewInterfaceMethodBuilder().
		WithName("TestMethod")
	builder.WithMethod(methodBuilder)

	// Test MethodCount after adding method
	if builder.MethodCount() != 1 {
		t.Errorf("Expected 1 method, got %d", builder.MethodCount())
	}

	// Test HasMethods after adding method
	if !builder.HasMethods() {
		t.Error("Expected HasMethods to return true after adding method")
	}

	// Test GetMethodNames
	names := builder.GetMethodNames()
	if len(names) != 1 {
		t.Errorf("Expected 1 method name, got %d", len(names))
	}

	if names[0] != "TestMethod" {
		t.Errorf("Expected method name 'TestMethod', got %s", names[0])
	}

	// Test ClearMethods
	result := builder.ClearMethods()
	if result != builder {
		t.Error("ClearMethods should return the builder for chaining")
	}

	if builder.MethodCount() != 0 {
		t.Errorf("Expected 0 methods after clear, got %d", builder.MethodCount())
	}

	if builder.HasMethods() {
		t.Error("Expected HasMethods to return false after clear")
	}
}

func TestInterfaceBuilder_MethodChaining(t *testing.T) {
	method1 := NewInterfaceMethodBuilder().
		WithName("Method1").
		AddArgField(StringField("param1")).
		AddRetvalField(ErrorField())

	method2 := NewInterfaceMethodBuilder().
		WithName("Method2").
		AddContextArg().
		AddRetvalField(IntField("result"))

	builder := NewInterfaceBuilder().
		WithName("TestInterface").
		WithMethod(method1).
		WithMethod(method2)

	if builder.name != "TestInterface" {
		t.Errorf("Expected interface name 'TestInterface', got %s", builder.name)
	}

	if builder.MethodCount() != 2 {
		t.Errorf("Expected 2 methods, got %d", builder.MethodCount())
	}

	names := builder.GetMethodNames()
	expectedNames := []string{"Method1", "Method2"}
	for i, expected := range expectedNames {
		if i >= len(names) || names[i] != expected {
			t.Errorf("Expected method name %d to be %s, got %s", i, expected, names[i])
		}
	}
}

func TestInterfaceBuilder_ComplexExample(t *testing.T) {
	// Create an interface similar to the example in the generated code
	methodBuilder := NewInterfaceMethodBuilder().
		WithName("HandleCreate").
		AddArgField(ContextField("ctx")).
		AddArgField(SelectorField("r", "apimodels", "CreateRequest")).
		AddRetvalField(SelectorField("", "apimodels", "CreateResponse")).
		AddErrorRetval()

	builder := NewInterfaceBuilder().
		WithName("CreateHandler").
		WithMethod(methodBuilder)

	decl := builder.Build()

	// Verify the structure
	typeSpec := decl.Specs[0].(*ast.TypeSpec)
	if typeSpec.Name.Name != "CreateHandler" {
		t.Errorf("Expected interface name 'CreateHandler', got %s", typeSpec.Name.Name)
	}

	interfaceType := typeSpec.Type.(*ast.InterfaceType)
	if len(interfaceType.Methods.List) != 1 {
		t.Errorf("Expected 1 method, got %d", len(interfaceType.Methods.List))
	}

	method := interfaceType.Methods.List[0]
	if method.Names[0].Name != "HandleCreate" {
		t.Errorf("Expected method name 'HandleCreate', got %s", method.Names[0].Name)
	}

	funcType := method.Type.(*ast.FuncType)
	if len(funcType.Params.List) != 2 {
		t.Errorf("Expected 2 parameters, got %d", len(funcType.Params.List))
	}

	if len(funcType.Results.List) != 2 {
		t.Errorf("Expected 2 return values, got %d", len(funcType.Results.List))
	}
}

func TestInterfaceBuilder_NewUtilityMethods(t *testing.T) {
	method1 := NewInterfaceMethodBuilder().
		WithName("Method1").
		AddArgField(StringField("param"))

	method2 := NewInterfaceMethodBuilder().
		WithName("Method2").
		AddRetvalField(IntField("result"))

	builder := NewInterfaceBuilder().
		WithName("TestInterface").
		WithMethod(method1).
		WithMethod(method2)

	// Test GetMethod
	retrievedMethod1 := builder.GetMethod(0)
	if retrievedMethod1 == nil {
		t.Fatal("GetMethod(0) returned nil")
	}
	if retrievedMethod1.name != "Method1" {
		t.Errorf("Expected method name 'Method1', got %s", retrievedMethod1.name)
	}

	retrievedMethod2 := builder.GetMethod(1)
	if retrievedMethod2 == nil {
		t.Fatal("GetMethod(1) returned nil")
	}
	if retrievedMethod2.name != "Method2" {
		t.Errorf("Expected method name 'Method2', got %s", retrievedMethod2.name)
	}

	// Test out of bounds
	if builder.GetMethod(2) != nil {
		t.Error("GetMethod(2) should return nil for out of bounds")
	}

	if builder.GetMethod(-1) != nil {
		t.Error("GetMethod(-1) should return nil for negative index")
	}

	// Test GetMethodByName
	foundMethod1 := builder.GetMethodByName("Method1")
	if foundMethod1 == nil {
		t.Fatal("GetMethodByName('Method1') returned nil")
	}
	if foundMethod1 != retrievedMethod1 {
		t.Error("GetMethodByName should return the same method instance")
	}

	foundMethod2 := builder.GetMethodByName("Method2")
	if foundMethod2 == nil {
		t.Fatal("GetMethodByName('Method2') returned nil")
	}
	if foundMethod2 != retrievedMethod2 {
		t.Error("GetMethodByName should return the same method instance")
	}

	// Test non-existent method
	if builder.GetMethodByName("NonExistent") != nil {
		t.Error("GetMethodByName('NonExistent') should return nil")
	}

	// Test RemoveMethod
	result := builder.RemoveMethod(0)
	if result != builder {
		t.Error("RemoveMethod should return the builder for chaining")
	}

	if builder.MethodCount() != 1 {
		t.Errorf("Expected 1 method after removal, got %d", builder.MethodCount())
	}

	// The remaining method should be Method2
	remainingMethod := builder.GetMethod(0)
	if remainingMethod.name != "Method2" {
		t.Errorf("Expected remaining method to be 'Method2', got %s", remainingMethod.name)
	}

	// Test RemoveMethod out of bounds (should not panic)
	builder.RemoveMethod(5)
	if builder.MethodCount() != 1 {
		t.Errorf("RemoveMethod out of bounds should not affect count, got %d", builder.MethodCount())
	}

	// Test RemoveMethodByName
	builder.RemoveMethodByName("Method2")
	if builder.MethodCount() != 0 {
		t.Errorf("Expected 0 methods after removing by name, got %d", builder.MethodCount())
	}

	// Test RemoveMethodByName on non-existent method (should not panic)
	builder.RemoveMethodByName("NonExistent")
	if builder.MethodCount() != 0 {
		t.Errorf("RemoveMethodByName on non-existent should not affect count, got %d", builder.MethodCount())
	}
}

// Integration tests moved from interface_method_builder_integration_test.go

func TestInterfaceMethodBuilder_AddArgFieldIntegration(t *testing.T) {
	fieldBuilder := NewFieldBuilder().
		WithName("param").
		WithType(String())

	methodBuilder := NewInterfaceMethodBuilder().
		WithName("TestMethod")

	result := methodBuilder.AddArgField(fieldBuilder)
	if result != methodBuilder {
		t.Error("AddArgField should return the method builder for chaining")
	}

	if len(methodBuilder.params) != 1 {
		t.Errorf("Expected 1 parameter, got %d", len(methodBuilder.params))
	}

	param := methodBuilder.params[0]
	if !param.HasName() {
		t.Error("Expected parameter to have a name")
	}

	if param.GetName() != "param" {
		t.Errorf("Expected parameter name 'param', got %s", param.GetName())
	}
}

func TestInterfaceMethodBuilder_AddArgFieldNil(t *testing.T) {
	methodBuilder := NewInterfaceMethodBuilder().WithName("TestMethod")

	defer func() {
		if r := recover(); r == nil {
			t.Error("AddArgField should panic when field builder is nil")
		}
	}()

	methodBuilder.AddArgField(nil)
}

func TestInterfaceMethodBuilder_AddArgFieldWithContext(t *testing.T) {
	methodBuilder := NewInterfaceMethodBuilder().WithName("TestMethod")
	fieldBuilder := ContextField("ctx")

	result := methodBuilder.AddArgField(fieldBuilder)
	if result != methodBuilder {
		t.Error("AddArgField should return the method builder for chaining")
	}

	if len(methodBuilder.params) != 1 {
		t.Errorf("Expected 1 parameter, got %d", len(methodBuilder.params))
	}

	param := methodBuilder.params[0]
	if !param.HasName() {
		t.Error("Expected parameter to have a name")
	}

	if param.GetName() != "ctx" {
		t.Errorf("Expected parameter name 'ctx', got %s", param.GetName())
	}
}

func TestInterfaceMethodBuilder_AddArgFieldNilType(t *testing.T) {
	methodBuilder := NewInterfaceMethodBuilder().WithName("TestMethod")

	defer func() {
		if r := recover(); r == nil {
			t.Error("AddArgField should panic when field builder is nil")
		}
	}()

	methodBuilder.AddArgField(nil)
}

func TestInterfaceMethodBuilder_AddRetvalFieldIntegration(t *testing.T) {
	fieldBuilder := NewFieldBuilder().
		WithName("result").
		WithType(Int())

	methodBuilder := NewInterfaceMethodBuilder().
		WithName("TestMethod")

	result := methodBuilder.AddRetvalField(fieldBuilder)
	if result != methodBuilder {
		t.Error("AddRetvalField should return the method builder for chaining")
	}

	if len(methodBuilder.results) != 1 {
		t.Errorf("Expected 1 return value, got %d", len(methodBuilder.results))
	}

	retval := methodBuilder.results[0]
	if !retval.HasName() {
		t.Error("Expected return value to have a name")
	}

	if retval.GetName() != "result" {
		t.Errorf("Expected return value name 'result', got %s", retval.GetName())
	}
}

func TestInterfaceMethodBuilder_AddRetvalFieldNil(t *testing.T) {
	methodBuilder := NewInterfaceMethodBuilder().WithName("TestMethod")

	defer func() {
		if r := recover(); r == nil {
			t.Error("AddRetvalField should panic when field builder is nil")
		}
	}()

	methodBuilder.AddRetvalField(nil)
}

func TestInterfaceMethodBuilder_AddRetvalFieldWithError(t *testing.T) {
	methodBuilder := NewInterfaceMethodBuilder().WithName("TestMethod")
	fieldBuilder := ErrorField()

	result := methodBuilder.AddRetvalField(fieldBuilder)
	if result != methodBuilder {
		t.Error("AddRetvalField should return the method builder for chaining")
	}

	if len(methodBuilder.results) != 1 {
		t.Errorf("Expected 1 return value, got %d", len(methodBuilder.results))
	}

	retval := methodBuilder.results[0]
	if retval.HasName() {
		t.Error("Expected error return value to not have a name")
	}
}

func TestInterfaceMethodBuilder_AddRetvalFieldNilType(t *testing.T) {
	methodBuilder := NewInterfaceMethodBuilder().WithName("TestMethod")

	defer func() {
		if r := recover(); r == nil {
			t.Error("AddRetvalField should panic when field builder is nil")
		}
	}()

	methodBuilder.AddRetvalField(nil)
}

func TestInterfaceMethodBuilder_AddPointerRetvalField(t *testing.T) {
	methodBuilder := NewInterfaceMethodBuilder().WithName("TestMethod")
	fieldBuilder := StringField("ptr")

	result := methodBuilder.AddRetvalField(fieldBuilder)
	if result != methodBuilder {
		t.Error("AddRetvalField should return the method builder for chaining")
	}

	if len(methodBuilder.results) != 1 {
		t.Errorf("Expected 1 return value, got %d", len(methodBuilder.results))
	}

	retval := methodBuilder.results[0]
	if !retval.HasName() {
		t.Error("Expected return value to have a name")
	}

	if retval.GetName() != "ptr" {
		t.Errorf("Expected return value name 'ptr', got %s", retval.GetName())
	}
}

func TestInterfaceMethodBuilder_ComplexExample(t *testing.T) {
	// Create a method similar to the generated code example
	methodBuilder := NewInterfaceMethodBuilder().
		WithName("HandleCreate").
		AddArgField(ContextField("ctx")).
		AddArgField(SelectorField("r", "apimodels", "CreateRequest")).
		AddRetvalField(SelectorField("", "apimodels", "CreateResponse")).
		AddRetvalField(ErrorField())

	builder := NewInterfaceBuilder().
		WithName("CreateHandler").
		WithMethod(methodBuilder)

	decl := builder.Build()

	// Verify the structure
	typeSpec := decl.Specs[0].(*ast.TypeSpec)
	if typeSpec.Name.Name != "CreateHandler" {
		t.Errorf("Expected interface name 'CreateHandler', got %s", typeSpec.Name.Name)
	}

	interfaceType := typeSpec.Type.(*ast.InterfaceType)
	if len(interfaceType.Methods.List) != 1 {
		t.Errorf("Expected 1 method, got %d", len(interfaceType.Methods.List))
	}

	method := interfaceType.Methods.List[0]
	if method.Names[0].Name != "HandleCreate" {
		t.Errorf("Expected method name 'HandleCreate', got %s", method.Names[0].Name)
	}

	funcType := method.Type.(*ast.FuncType)
	if len(funcType.Params.List) != 2 {
		t.Errorf("Expected 2 parameters, got %d", len(funcType.Params.List))
	}

	if len(funcType.Results.List) != 2 {
		t.Errorf("Expected 2 return values, got %d", len(funcType.Results.List))
	}

	// Check first parameter (context.Context)
	ctxParam := funcType.Params.List[0]
	if ctxParam.Names[0].Name != "ctx" {
		t.Error("First parameter should be named 'ctx'")
	}

	// Check second parameter (apimodels.CreateRequest)
	reqParam := funcType.Params.List[1]
	if reqParam.Names[0].Name != "r" {
		t.Error("Second parameter should be named 'r'")
	}

	// Check first return value (pointer to apimodels.CreateResponse)
	responseRet := funcType.Results.List[0]
	if len(responseRet.Names) != 0 {
		t.Error("First return value should be unnamed")
	}

	// Check second return value (error)
	errorRet := funcType.Results.List[1]
	if len(errorRet.Names) != 0 {
		t.Error("Second return value should be unnamed")
	}
}

func TestInterfaceMethodBuilder_MethodChainingWithNewBuilders(t *testing.T) {
	methodBuilder := NewInterfaceMethodBuilder().
		WithName("TestMethod").
		AddArgField(StringField("name")).
		AddArgField(ContextField("ctx")).
		AddRetvalField(IntField("count")).
		AddRetvalField(ErrorField())

	builder := NewInterfaceBuilder().
		WithName("TestInterface").
		WithMethod(methodBuilder)

	if builder.MethodCount() != 1 {
		t.Errorf("Expected 1 method, got %d", builder.MethodCount())
	}

	method := builder.GetMethod(0)
	if len(method.params) != 2 {
		t.Errorf("Expected 2 parameters, got %d", len(method.params))
	}

	if len(method.results) != 2 {
		t.Errorf("Expected 2 return values, got %d", len(method.results))
	}

	// Check parameters
	if method.params[0].GetName() != "name" {
		t.Error("First parameter should be named 'name'")
	}

	if method.params[1].GetName() != "ctx" {
		t.Error("Second parameter should be named 'ctx'")
	}

	// Check return values
	if method.results[0].GetName() != "count" {
		t.Error("First return value should be named 'count'")
	}

	if method.results[1].HasName() {
		t.Error("Second return value should be unnamed")
	}
}
