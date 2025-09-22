package astbuilder

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

func TestNewValidationBuilder(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	validationConfig := ValidationConfig{
		PackageName:   "api",
		UsePointers:   true,
		ValidatorName: "validator",
	}

	validationBuilder := NewValidationBuilder(builder, validationConfig)

	if validationBuilder == nil {
		t.Fatal("NewValidationBuilder returned nil")
	}

	if validationBuilder.config.PackageName != "api" {
		t.Errorf("Expected package name 'api', got '%s'", validationBuilder.config.PackageName)
	}

	if !validationBuilder.config.UsePointers {
		t.Error("Expected UsePointers to be true")
	}

	if validationBuilder.config.ValidatorName != "validator" {
		t.Errorf("Expected validator name 'validator', got '%s'", validationBuilder.config.ValidatorName)
	}

	if validationBuilder.builder != builder {
		t.Error("Expected validation builder to use the provided builder")
	}
}

func TestValidationBuilder_BuildObjectValidation(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	validationConfig := ValidationConfig{
		PackageName:   "api",
		UsePointers:   true,
		ValidatorName: "validator",
	}

	validationBuilder := NewValidationBuilder(builder, validationConfig)

	// Create test schema
	schema := &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type: &openapi3.Types{"object"},
			Properties: openapi3.Schemas{
				"name": &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type:      &openapi3.Types{"string"},
						MinLength: 1,
						MaxLength: func() *uint64 { v := uint64(100); return &v }(),
					},
				},
				"age": &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: &openapi3.Types{"integer"},
						Min:  func() *float64 { v := 0.0; return &v }(),
						Max:  func() *float64 { v := 120.0; return &v }(),
					},
				},
			},
			Required: []string{"name"},
		},
	}

	err := validationBuilder.BuildObjectValidation("User", schema)
	if err != nil {
		t.Fatalf("BuildObjectValidation failed: %v", err)
	}

	// Check that declarations were added
	decls := builder.decls
	if len(decls) == 0 {
		t.Error("Expected validation declarations to be added")
	}

	// Check that validator import was added
	imports := builder.imports
	if !imports["github.com/go-playground/validator/v10"] {
		t.Error("Expected validator import to be added")
	}
}

func TestValidationBuilder_BuildArrayValidation(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	validationConfig := ValidationConfig{
		PackageName:   "api",
		UsePointers:   true,
		ValidatorName: "validator",
	}

	validationBuilder := NewValidationBuilder(builder, validationConfig)

	// Create test array schema
	schema := &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type: &openapi3.Types{"array"},
			Items: &openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: &openapi3.Types{"string"},
				},
			},
			MinItems: 1,
			MaxItems: func() *uint64 { v := uint64(10); return &v }(),
		},
	}

	err := validationBuilder.BuildArrayValidation("UserList", schema)
	if err != nil {
		t.Fatalf("BuildArrayValidation failed: %v", err)
	}

	// Check that declarations were added
	decls := builder.decls
	if len(decls) == 0 {
		t.Error("Expected validation declarations to be added")
	}

	// Check that validator import was added
	imports := builder.imports
	if !imports["github.com/go-playground/validator/v10"] {
		t.Error("Expected validator import to be added")
	}
}

func TestValidationBuilder_BuildFieldValidation(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	validationConfig := ValidationConfig{
		PackageName:   "api",
		UsePointers:   true,
		ValidatorName: "validator",
	}

	validationBuilder := NewValidationBuilder(builder, validationConfig)

	// Test string field validation
	stringSchema := &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type:      &openapi3.Types{"string"},
			MinLength: 1,
			MaxLength: func() *uint64 { v := uint64(100); return &v }(),
			Pattern:   "^[a-zA-Z]+$",
		},
	}

	rules := validationBuilder.BuildFieldValidation("name", stringSchema)
	expectedRules := []string{"min=1", "max=100", "regexp=^[a-zA-Z]+$"}

	if len(rules) != len(expectedRules) {
		t.Errorf("Expected %d rules, got %d", len(expectedRules), len(rules))
	}

	for _, expectedRule := range expectedRules {
		found := false
		for _, rule := range rules {
			if rule == expectedRule {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected rule '%s' not found", expectedRule)
		}
	}
}

func TestValidationBuilder_AddRequiredFieldsValidation(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	validationConfig := ValidationConfig{
		PackageName:   "api",
		UsePointers:   true,
		ValidatorName: "validator",
	}

	validationBuilder := NewValidationBuilder(builder, validationConfig)

	fields := []string{"name", "email", "age"}
	validationBuilder.AddRequiredFieldsValidation(fields)

	// Check that statements were added
	stmts := builder.stmts
	if len(stmts) == 0 {
		t.Error("Expected validation statements to be added")
	}

	// Check that errors import was added
	imports := builder.imports
	if !imports["errors"] {
		t.Error("Expected errors import to be added")
	}
}

func TestValidationBuilder_AddFieldValidation(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	validationConfig := ValidationConfig{
		PackageName:   "api",
		UsePointers:   true,
		ValidatorName: "validator",
	}

	validationBuilder := NewValidationBuilder(builder, validationConfig)

	validators := []string{"required", "email", "min=1", "max=100"}
	validationBuilder.AddFieldValidation("email", validators)

	// Check that statements were added
	stmts := builder.stmts
	if len(stmts) == 0 {
		t.Error("Expected validation statements to be added")
	}

	// Check that validator import was added
	imports := builder.imports
	if !imports["github.com/go-playground/validator/v10"] {
		t.Error("Expected validator import to be added")
	}
}

func TestValidationBuilder_AddJSONUnmarshal(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	validationConfig := ValidationConfig{
		PackageName:   "api",
		UsePointers:   true,
		ValidatorName: "validator",
	}

	validationBuilder := NewValidationBuilder(builder, validationConfig)

	validationBuilder.AddJSONUnmarshal()

	// Check that statements were added
	stmts := builder.stmts
	if len(stmts) == 0 {
		t.Error("Expected JSON unmarshal statements to be added")
	}

	// Check that json import was added
	imports := builder.imports
	if !imports["encoding/json"] {
		t.Error("Expected json import to be added")
	}
}

func TestValidationBuilder_AddErrorHandling(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	validationConfig := ValidationConfig{
		PackageName:   "api",
		UsePointers:   true,
		ValidatorName: "validator",
	}

	validationBuilder := NewValidationBuilder(builder, validationConfig)

	validationBuilder.AddErrorHandling()

	// Check that statements were added
	stmts := builder.stmts
	if len(stmts) == 0 {
		t.Error("Expected error handling statements to be added")
	}
}

func TestValidationBuilder_FluentInterface(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	validationConfig := ValidationConfig{
		PackageName:   "api",
		UsePointers:   true,
		ValidatorName: "validator",
	}

	validationBuilder := NewValidationBuilder(builder, validationConfig)

	// Test fluent interface
	validationBuilder = validationBuilder.WithPackageName("custom").
		WithUsePointers(false).
		WithImportPrefix("github.com/custom").
		WithValidatorName("customValidator").
		WithErrorHandler("customErrorHandler")

	if validationBuilder.config.PackageName != "custom" {
		t.Errorf("Expected package name 'custom', got '%s'", validationBuilder.config.PackageName)
	}

	if validationBuilder.config.UsePointers {
		t.Error("Expected UsePointers to be false")
	}

	if validationBuilder.config.ImportPrefix != "github.com/custom" {
		t.Errorf("Expected import prefix 'github.com/custom', got '%s'", validationBuilder.config.ImportPrefix)
	}

	if validationBuilder.config.ValidatorName != "customValidator" {
		t.Errorf("Expected validator name 'customValidator', got '%s'", validationBuilder.config.ValidatorName)
	}

	if validationBuilder.config.ErrorHandler != "customErrorHandler" {
		t.Errorf("Expected error handler 'customErrorHandler', got '%s'", validationBuilder.config.ErrorHandler)
	}
}

func TestValidationBuilder_GetConfig(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	validationConfig := ValidationConfig{
		PackageName:   "api",
		UsePointers:   true,
		ValidatorName: "validator",
	}

	validationBuilder := NewValidationBuilder(builder, validationConfig)

	retrievedConfig := validationBuilder.GetConfig()

	if retrievedConfig.PackageName != validationConfig.PackageName {
		t.Error("GetConfig should return the same package name")
	}

	if retrievedConfig.UsePointers != validationConfig.UsePointers {
		t.Error("GetConfig should return the same UsePointers value")
	}

	if retrievedConfig.ValidatorName != validationConfig.ValidatorName {
		t.Error("GetConfig should return the same validator name")
	}
}

func TestValidationBuilder_GetBuilder(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	validationConfig := ValidationConfig{
		PackageName:   "api",
		UsePointers:   true,
		ValidatorName: "validator",
	}

	validationBuilder := NewValidationBuilder(builder, validationConfig)

	retrievedBuilder := validationBuilder.GetBuilder()

	if retrievedBuilder != builder {
		t.Error("GetBuilder should return the same builder instance")
	}
}

func TestValidationBuilder_HelperMethods(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	validationConfig := ValidationConfig{
		PackageName:   "api",
		UsePointers:   true,
		ValidatorName: "validator",
	}

	validationBuilder := NewValidationBuilder(builder, validationConfig)

	// Test PascalCase conversion
	pascalCase := validationBuilder.toPascalCase("user_name_field")
	expected := "UserNameField"
	if pascalCase != expected {
		t.Errorf("Expected PascalCase 'UserNameField', got '%s'", pascalCase)
	}

	// Test empty string
	pascalCase = validationBuilder.toPascalCase("")
	expected = ""
	if pascalCase != expected {
		t.Errorf("Expected empty string, got '%s'", pascalCase)
	}
}

func TestValidationBuilder_ErrorHandling(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	validationConfig := ValidationConfig{
		PackageName:   "api",
		UsePointers:   true,
		ValidatorName: "validator",
	}

	validationBuilder := NewValidationBuilder(builder, validationConfig)

	// Test nil schema
	err := validationBuilder.BuildObjectValidation("User", nil)
	if err == nil {
		t.Error("Expected error for nil schema")
	}

	// Test nil schema value
	schema := &openapi3.SchemaRef{}
	err = validationBuilder.BuildObjectValidation("User", schema)
	if err == nil {
		t.Error("Expected error for nil schema value")
	}
}

func TestValidationBuilder_StringValidations(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	validationConfig := ValidationConfig{
		PackageName:   "api",
		UsePointers:   true,
		ValidatorName: "validator",
	}

	validationBuilder := NewValidationBuilder(builder, validationConfig)

	// Test string validations
	schema := &openapi3.Schema{
		Type:      &openapi3.Types{"string"},
		MinLength: 5,
		MaxLength: func() *uint64 { v := uint64(50); return &v }(),
		Pattern:   "^[a-zA-Z0-9]+$",
	}

	rules := validationBuilder.buildStringValidations(schema)
	expectedRules := []string{"min=5", "max=50", "regexp=^[a-zA-Z0-9]+$"}

	if len(rules) != len(expectedRules) {
		t.Errorf("Expected %d rules, got %d", len(expectedRules), len(rules))
	}

	for _, expectedRule := range expectedRules {
		found := false
		for _, rule := range rules {
			if rule == expectedRule {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected rule '%s' not found", expectedRule)
		}
	}
}

func TestValidationBuilder_NumericValidations(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	validationConfig := ValidationConfig{
		PackageName:   "api",
		UsePointers:   true,
		ValidatorName: "validator",
	}

	validationBuilder := NewValidationBuilder(builder, validationConfig)

	// Test numeric validations
	schema := &openapi3.Schema{
		Type:       &openapi3.Types{"integer"},
		Min:        func() *float64 { v := 0.0; return &v }(),
		Max:        func() *float64 { v := 100.0; return &v }(),
		MultipleOf: func() *float64 { v := 2.0; return &v }(),
	}

	rules := validationBuilder.buildNumericValidations(schema)
	expectedRules := []string{"min=0", "max=100", "multipleof=2"}

	if len(rules) != len(expectedRules) {
		t.Errorf("Expected %d rules, got %d", len(expectedRules), len(rules))
	}

	for _, expectedRule := range expectedRules {
		found := false
		for _, rule := range rules {
			if rule == expectedRule {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected rule '%s' not found", expectedRule)
		}
	}
}

func TestValidationBuilder_ArrayValidations(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	validationConfig := ValidationConfig{
		PackageName:   "api",
		UsePointers:   true,
		ValidatorName: "validator",
	}

	validationBuilder := NewValidationBuilder(builder, validationConfig)

	// Test array validations
	schema := &openapi3.Schema{
		Type:        &openapi3.Types{"array"},
		MinItems:    1,
		MaxItems:    func() *uint64 { v := uint64(10); return &v }(),
		UniqueItems: true,
	}

	rules := validationBuilder.buildArrayValidations(schema)
	expectedRules := []string{"min=1", "max=10", "unique"}

	if len(rules) != len(expectedRules) {
		t.Errorf("Expected %d rules, got %d", len(expectedRules), len(rules))
	}

	for _, expectedRule := range expectedRules {
		found := false
		for _, rule := range rules {
			if rule == expectedRule {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected rule '%s' not found", expectedRule)
		}
	}
}

func TestValidationBuilder_FormatValidations(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	validationConfig := ValidationConfig{
		PackageName:   "api",
		UsePointers:   true,
		ValidatorName: "validator",
	}

	validationBuilder := NewValidationBuilder(builder, validationConfig)

	// Test format validations
	testCases := []struct {
		format        string
		expectedRules []string
	}{
		{"email", []string{"email"}},
		{"uri", []string{"uri"}},
		{"date", []string{"datetime=2006-01-02"}},
		{"date-time", []string{"datetime=2006-01-02T15:04:05Z07:00"}},
		{"uuid", []string{"uuid"}},
		{"unknown", nil},
	}

	for _, tc := range testCases {
		rules := validationBuilder.buildFormatValidations(tc.format)
		if len(rules) != len(tc.expectedRules) {
			t.Errorf("For format '%s', expected %d rules, got %d", tc.format, len(tc.expectedRules), len(rules))
		}

		for i, expectedRule := range tc.expectedRules {
			if i < len(rules) && rules[i] != expectedRule {
				t.Errorf("For format '%s', expected rule '%s', got '%s'", tc.format, expectedRule, rules[i])
			}
		}
	}
}

func TestValidationBuilder_EnumValidations(t *testing.T) {
	config := BuilderConfig{PackageName: "test"}
	builder := NewBuilder(config)

	validationConfig := ValidationConfig{
		PackageName:   "api",
		UsePointers:   true,
		ValidatorName: "validator",
	}

	validationBuilder := NewValidationBuilder(builder, validationConfig)

	// Test enum validations
	enum := []interface{}{"red", "green", "blue"}
	rules := validationBuilder.buildEnumValidations(enum)
	expectedRule := "oneof=red green blue"

	if len(rules) != 1 {
		t.Errorf("Expected 1 rule, got %d", len(rules))
	}

	if len(rules) > 0 && rules[0] != expectedRule {
		t.Errorf("Expected rule '%s', got '%s'", expectedRule, rules[0])
	}

	// Test empty enum
	emptyEnum := []interface{}{}
	rules = validationBuilder.buildEnumValidations(emptyEnum)
	if len(rules) != 0 {
		t.Errorf("Expected 0 rules for empty enum, got %d", len(rules))
	}
}
