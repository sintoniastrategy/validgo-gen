package astbuilder

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

func TestValidationBuilder_CompleteWorkflow(t *testing.T) {
	// Test complete validation building workflow
	config := BuilderConfig{
		PackageName:  "api",
		ImportPrefix: "github.com/example/api",
		UsePointers:  true,
	}

	builder := NewBuilder(config)

	validationConfig := ValidationConfig{
		PackageName:   "api",
		UsePointers:   true,
		ValidatorName: "validator",
	}

	validationBuilder := NewValidationBuilder(builder, validationConfig)

	// Test 1: Add required fields validation
	validationBuilder.AddRequiredFieldsValidation([]string{"name", "email"})

	// Test 2: Add field validation
	validationBuilder.AddFieldValidation("email", []string{"required", "email"})

	// Test 3: Add JSON unmarshal
	validationBuilder.AddJSONUnmarshal()

	// Test 4: Add error handling
	validationBuilder.AddErrorHandling()

	// Verify results
	stmts := builder.stmts
	if len(stmts) == 0 {
		t.Error("Expected validation statements to be added")
	}

	// Check that required imports are present
	imports := builder.imports
	requiredImports := []string{
		"errors",
		"github.com/go-playground/validator/v10",
		"encoding/json",
	}

	for _, requiredImport := range requiredImports {
		if !imports[requiredImport] {
			t.Errorf("Expected import '%s' to be present", requiredImport)
		}
	}
}

func TestValidationBuilder_ChainedWorkflow(t *testing.T) {
	// Test chained validation building workflow
	config := BuilderConfig{
		PackageName:  "api",
		ImportPrefix: "github.com/example/api",
		UsePointers:  true,
	}

	builder := NewBuilder(config)

	validationConfig := ValidationConfig{
		PackageName:   "api",
		UsePointers:   true,
		ValidatorName: "validator",
	}

	validationBuilder := NewValidationBuilder(builder, validationConfig)

	// Test chained workflow
	validationBuilder.
		WithPackageName("custom").
		WithUsePointers(false).
		WithValidatorName("customValidator").
		AddRequiredFieldsValidation([]string{"id", "name"}).
		AddFieldValidation("email", []string{"required", "email"}).
		AddJSONUnmarshal().
		AddErrorHandling()

	// Verify that all operations were performed
	stmts := builder.stmts
	if len(stmts) == 0 {
		t.Error("Expected validation statements from chained workflow")
	}

	// Verify configuration changes
	if validationBuilder.config.PackageName != "custom" {
		t.Errorf("Expected package name 'custom', got '%s'", validationBuilder.config.PackageName)
	}

	if validationBuilder.config.UsePointers {
		t.Error("Expected UsePointers to be false")
	}

	if validationBuilder.config.ValidatorName != "customValidator" {
		t.Errorf("Expected validator name 'customValidator', got '%s'", validationBuilder.config.ValidatorName)
	}
}

func TestValidationBuilder_ObjectValidationIntegration(t *testing.T) {
	// Test object validation integration
	config := BuilderConfig{
		PackageName:  "api",
		ImportPrefix: "github.com/example/api",
		UsePointers:  true,
	}

	builder := NewBuilder(config)

	validationConfig := ValidationConfig{
		PackageName:   "api",
		UsePointers:   true,
		ValidatorName: "validator",
	}

	validationBuilder := NewValidationBuilder(builder, validationConfig)

	// Create comprehensive object schema
	schema := &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type: &openapi3.Types{"object"},
			Properties: openapi3.Schemas{
				"id": &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: &openapi3.Types{"integer"},
						Min:  func() *float64 { v := 1.0; return &v }(),
					},
				},
				"name": &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type:      &openapi3.Types{"string"},
						MinLength: 1,
						MaxLength: func() *uint64 { v := uint64(100); return &v }(),
					},
				},
				"email": &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type:   &openapi3.Types{"string"},
						Format: "email",
					},
				},
				"age": &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: &openapi3.Types{"integer"},
						Min:  func() *float64 { v := 0.0; return &v }(),
						Max:  func() *float64 { v := 120.0; return &v }(),
					},
				},
				"tags": &openapi3.SchemaRef{
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
				},
			},
			Required: []string{"id", "name", "email"},
		},
	}

	err := validationBuilder.BuildObjectValidation("User", schema)
	if err != nil {
		t.Fatalf("BuildObjectValidation failed: %v", err)
	}

	// Verify that declarations were added
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

func TestValidationBuilder_ArrayValidationIntegration(t *testing.T) {
	// Test array validation integration
	config := BuilderConfig{
		PackageName:  "api",
		ImportPrefix: "github.com/example/api",
		UsePointers:  true,
	}

	builder := NewBuilder(config)

	validationConfig := ValidationConfig{
		PackageName:   "api",
		UsePointers:   true,
		ValidatorName: "validator",
	}

	validationBuilder := NewValidationBuilder(builder, validationConfig)

	// Create comprehensive array schema
	schema := &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type: &openapi3.Types{"array"},
			Items: &openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: &openapi3.Types{"object"},
					Properties: openapi3.Schemas{
						"id": &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type: &openapi3.Types{"integer"},
							},
						},
						"name": &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type: &openapi3.Types{"string"},
							},
						},
					},
					Required: []string{"id", "name"},
				},
			},
			MinItems:    1,
			MaxItems:    func() *uint64 { v := uint64(100); return &v }(),
			UniqueItems: true,
		},
	}

	err := validationBuilder.BuildArrayValidation("UserList", schema)
	if err != nil {
		t.Fatalf("BuildArrayValidation failed: %v", err)
	}

	// Verify that declarations were added
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

func TestValidationBuilder_FieldValidationTypes(t *testing.T) {
	// Test different field validation types
	config := BuilderConfig{
		PackageName:  "api",
		ImportPrefix: "github.com/example/api",
		UsePointers:  true,
	}

	builder := NewBuilder(config)

	validationConfig := ValidationConfig{
		PackageName:   "api",
		UsePointers:   true,
		ValidatorName: "validator",
	}

	validationBuilder := NewValidationBuilder(builder, validationConfig)

	// Test string field with various validations
	stringSchema := &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type:      &openapi3.Types{"string"},
			MinLength: 5,
			MaxLength: func() *uint64 { v := uint64(50); return &v }(),
			Pattern:   "^[a-zA-Z0-9]+$",
			Format:    "email",
			Enum:      []interface{}{"active", "inactive", "pending"},
		},
	}

	rules := validationBuilder.BuildFieldValidation("status", stringSchema)
	expectedRules := []string{"min=5", "max=50", "regexp=^[a-zA-Z0-9]+$", "email", "oneof=active inactive pending"}

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

func TestValidationBuilder_ComplexObjectValidation(t *testing.T) {
	// Test complex object validation with nested structures
	config := BuilderConfig{
		PackageName:  "api",
		ImportPrefix: "github.com/example/api",
		UsePointers:  true,
	}

	builder := NewBuilder(config)

	validationConfig := ValidationConfig{
		PackageName:   "api",
		UsePointers:   true,
		ValidatorName: "validator",
	}

	validationBuilder := NewValidationBuilder(builder, validationConfig)

	// Create complex object schema
	schema := &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type: &openapi3.Types{"object"},
			Properties: openapi3.Schemas{
				"user": &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: &openapi3.Types{"object"},
						Properties: openapi3.Schemas{
							"id": &openapi3.SchemaRef{
								Value: &openapi3.Schema{
									Type: &openapi3.Types{"integer"},
									Min:  func() *float64 { v := 1.0; return &v }(),
								},
							},
							"profile": &openapi3.SchemaRef{
								Value: &openapi3.Schema{
									Type: &openapi3.Types{"object"},
									Properties: openapi3.Schemas{
										"firstName": &openapi3.SchemaRef{
											Value: &openapi3.Schema{
												Type:      &openapi3.Types{"string"},
												MinLength: 1,
												MaxLength: func() *uint64 { v := uint64(50); return &v }(),
											},
										},
										"lastName": &openapi3.SchemaRef{
											Value: &openapi3.Schema{
												Type:      &openapi3.Types{"string"},
												MinLength: 1,
												MaxLength: func() *uint64 { v := uint64(50); return &v }(),
											},
										},
									},
									Required: []string{"firstName", "lastName"},
								},
							},
						},
						Required: []string{"id", "profile"},
					},
				},
				"metadata": &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: &openapi3.Types{"object"},
						AdditionalProperties: openapi3.AdditionalProperties{
							Schema: &openapi3.SchemaRef{
								Value: &openapi3.Schema{
									Type: &openapi3.Types{"string"},
								},
							},
						},
					},
				},
			},
			Required: []string{"user"},
		},
	}

	err := validationBuilder.BuildObjectValidation("ComplexUser", schema)
	if err != nil {
		t.Fatalf("BuildObjectValidation failed: %v", err)
	}

	// Verify that declarations were added
	decls := builder.decls
	if len(decls) == 0 {
		t.Error("Expected validation declarations to be added")
	}
}

func TestValidationBuilder_EmptySchema(t *testing.T) {
	// Test handling of empty schema
	config := BuilderConfig{
		PackageName:  "api",
		ImportPrefix: "github.com/example/api",
		UsePointers:  true,
	}

	builder := NewBuilder(config)

	validationConfig := ValidationConfig{
		PackageName:   "api",
		UsePointers:   true,
		ValidatorName: "validator",
	}

	validationBuilder := NewValidationBuilder(builder, validationConfig)

	// Test with empty schema
	emptySchema := &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type:       &openapi3.Types{"object"},
			Properties: openapi3.Schemas{},
		},
	}

	err := validationBuilder.BuildObjectValidation("EmptyModel", emptySchema)
	if err != nil {
		t.Fatalf("BuildObjectValidation with empty schema failed: %v", err)
	}

	// Should still add basic validation components
	decls := builder.decls
	if len(decls) == 0 {
		t.Error("Expected at least one declaration for empty schema")
	}
}

func TestValidationBuilder_MultipleValidators(t *testing.T) {
	// Test multiple validators on the same field
	config := BuilderConfig{
		PackageName:  "api",
		ImportPrefix: "github.com/example/api",
		UsePointers:  true,
	}

	builder := NewBuilder(config)

	validationConfig := ValidationConfig{
		PackageName:   "api",
		UsePointers:   true,
		ValidatorName: "validator",
	}

	validationBuilder := NewValidationBuilder(builder, validationConfig)

	// Add multiple validators for the same field
	validationBuilder.
		AddFieldValidation("email", []string{"required", "email"}).
		AddFieldValidation("password", []string{"required", "min=8", "max=128"}).
		AddFieldValidation("age", []string{"required", "min=18", "max=120"})

	// Verify that statements were added
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

func TestValidationBuilder_ValidationRuleGeneration(t *testing.T) {
	// Test validation rule generation for different field types
	config := BuilderConfig{
		PackageName:  "api",
		ImportPrefix: "github.com/example/api",
		UsePointers:  true,
	}

	builder := NewBuilder(config)

	validationConfig := ValidationConfig{
		PackageName:   "api",
		UsePointers:   true,
		ValidatorName: "validator",
	}

	validationBuilder := NewValidationBuilder(builder, validationConfig)

	// Test cases for different field types
	testCases := []struct {
		name          string
		schema        *openapi3.SchemaRef
		expectedRules []string
	}{
		{
			name: "string with length constraints",
			schema: &openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:      &openapi3.Types{"string"},
					MinLength: 5,
					MaxLength: func() *uint64 { v := uint64(100); return &v }(),
				},
			},
			expectedRules: []string{"min=5", "max=100"},
		},
		{
			name: "integer with range constraints",
			schema: &openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: &openapi3.Types{"integer"},
					Min:  func() *float64 { v := 0.0; return &v }(),
					Max:  func() *float64 { v := 100.0; return &v }(),
				},
			},
			expectedRules: []string{"min=0", "max=100"},
		},
		{
			name: "array with item constraints",
			schema: &openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:     &openapi3.Types{"array"},
					MinItems: 1,
					MaxItems: func() *uint64 { v := uint64(10); return &v }(),
				},
			},
			expectedRules: []string{"min=1", "max=10"},
		},
		{
			name: "string with format",
			schema: &openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:   &openapi3.Types{"string"},
					Format: "email",
				},
			},
			expectedRules: []string{"email"},
		},
		{
			name: "string with enum",
			schema: &openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: &openapi3.Types{"string"},
					Enum: []interface{}{"red", "green", "blue"},
				},
			},
			expectedRules: []string{"oneof=red green blue"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rules := validationBuilder.BuildFieldValidation("testField", tc.schema)

			if len(rules) != len(tc.expectedRules) {
				t.Errorf("Expected %d rules, got %d", len(tc.expectedRules), len(rules))
			}

			for _, expectedRule := range tc.expectedRules {
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
		})
	}
}
