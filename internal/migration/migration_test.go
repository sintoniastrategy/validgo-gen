package migration

import (
	"context"
	"go/ast"
	"strings"
	"testing"

	"github.com/jolfzverb/codegen/internal/generator/options"
)

func TestNewMigrationWrapper(t *testing.T) {
	opts := &options.Options{
		PackagePrefix: "github.com/test",
		YAMLFiles:     []string{"test.yaml"},
	}

	wrapper := NewMigrationWrapper(opts, MigrationModeHybrid)
	if wrapper == nil {
		t.Fatal("NewMigrationWrapper returned nil")
	}

	if wrapper.legacyGenerator == nil {
		t.Fatal("Legacy generator is nil")
	}

	if wrapper.migratedGenerator == nil {
		t.Fatal("Migrated generator is nil")
	}

	if wrapper.mode != MigrationModeHybrid {
		t.Errorf("Expected mode %d, got %d", MigrationModeHybrid, wrapper.mode)
	}
}

func TestMigrationWrapper_Generate(t *testing.T) {
	opts := &options.Options{
		PackagePrefix: "github.com/test",
		YAMLFiles:     []string{"test.yaml"},
	}

	// Test legacy mode
	wrapper := NewMigrationWrapper(opts, MigrationModeLegacy)
	err := wrapper.Generate(context.Background())
	// This will likely fail due to missing YAML files, but we're testing the flow
	if err != nil && !strings.Contains(err.Error(), "no such file or directory") {
		t.Errorf("Unexpected error in legacy mode: %v", err)
	}

	// Test new mode
	wrapper = NewMigrationWrapper(opts, MigrationModeNew)
	err = wrapper.Generate(context.Background())
	// This will likely fail due to missing YAML files, but we're testing the flow
	if err != nil && !strings.Contains(err.Error(), "no such file or directory") {
		t.Errorf("Unexpected error in new mode: %v", err)
	}

	// Test hybrid mode
	wrapper = NewMigrationWrapper(opts, MigrationModeHybrid)
	err = wrapper.Generate(context.Background())
	// This will likely fail due to missing YAML files, but we're testing the flow
	if err != nil && !strings.Contains(err.Error(), "no such file or directory") {
		t.Errorf("Unexpected error in hybrid mode: %v", err)
	}
}

func TestMigrationWrapper_SetMigrationMode(t *testing.T) {
	opts := &options.Options{
		PackagePrefix: "github.com/test",
		YAMLFiles:     []string{"test.yaml"},
	}

	wrapper := NewMigrationWrapper(opts, MigrationModeLegacy)

	// Test initial mode
	if wrapper.mode != MigrationModeLegacy {
		t.Errorf("Expected initial mode %d, got %d", MigrationModeLegacy, wrapper.mode)
	}

	// Change mode
	wrapper.SetMigrationMode(MigrationModeNew)
	if wrapper.mode != MigrationModeNew {
		t.Errorf("Expected mode %d after change, got %d", MigrationModeNew, wrapper.mode)
	}

	// Test config update
	if !wrapper.config.UseNewBuilders {
		t.Error("Expected UseNewBuilders to be true in new mode")
	}
}

func TestMigrationWrapper_GetMigrationStatus(t *testing.T) {
	opts := &options.Options{
		PackagePrefix: "github.com/test",
		YAMLFiles:     []string{"test.yaml"},
	}

	wrapper := NewMigrationWrapper(opts, MigrationModeHybrid)
	status := wrapper.GetMigrationStatus()

	if status["mode"] != MigrationModeHybrid {
		t.Errorf("Expected mode %d, got %v", MigrationModeHybrid, status["mode"])
	}

	if status["use_new_builders"] != true {
		t.Error("Expected use_new_builders to be true")
	}

	if status["compare_outputs"] != true {
		t.Error("Expected compare_outputs to be true in hybrid mode")
	}
}

func TestMigrationWrapper_ValidateMigration(t *testing.T) {
	opts := &options.Options{
		PackagePrefix: "github.com/test",
		YAMLFiles:     []string{"test.yaml"},
	}

	wrapper := NewMigrationWrapper(opts, MigrationModeNew)

	// Test validation with empty builder (should fail)
	err := wrapper.ValidateMigration()
	if err == nil {
		t.Error("Expected validation to fail with empty builder")
	}

	// Add some content to the builder
	astBuilder := wrapper.GetASTBuilder()
	if astBuilder == nil {
		t.Fatal("AST builder is nil")
	}

	astBuilder.AddImport("fmt")
	astBuilder.AddDeclaration(&ast.FuncDecl{
		Name: ast.NewIdent("Test"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{},
		},
		Body: &ast.BlockStmt{},
	})

	// Test validation with content (should pass)
	err = wrapper.ValidateMigration()
	if err != nil {
		t.Errorf("Validation should pass with content: %v", err)
	}
}

func TestMigrationWrapper_CompareOutputs(t *testing.T) {
	opts := &options.Options{
		PackagePrefix: "github.com/test",
		YAMLFiles:     []string{"test.yaml"},
	}

	wrapper := NewMigrationWrapper(opts, MigrationModeHybrid)

	// Test comparison
	comparison, err := wrapper.CompareOutputs(context.Background())
	if err != nil {
		t.Errorf("CompareOutputs failed: %v", err)
	}

	if comparison == nil {
		t.Fatal("CompareOutputs returned nil comparison")
	}

	// Test with non-hybrid mode (should fail)
	wrapper.SetMigrationMode(MigrationModeNew)
	_, err = wrapper.CompareOutputs(context.Background())
	if err == nil {
		t.Error("Expected error when comparing in non-hybrid mode")
	}
}

func TestNewMigratedGenerator(t *testing.T) {
	opts := &options.Options{
		PackagePrefix: "github.com/test",
		YAMLFiles:     []string{"test.yaml"},
	}

	migratedGen := NewMigratedGenerator(opts)
	if migratedGen == nil {
		t.Fatal("NewMigratedGenerator returned nil")
	}

	if migratedGen.Generator == nil {
		t.Fatal("Embedded Generator is nil")
	}

	if migratedGen.astBuilder == nil {
		t.Fatal("AST builder is nil")
	}
}

func TestMigratedGenerator_GetASTBuilder(t *testing.T) {
	opts := &options.Options{
		PackagePrefix: "github.com/test",
		YAMLFiles:     []string{"test.yaml"},
	}

	migratedGen := NewMigratedGenerator(opts)
	astBuilder := migratedGen.GetASTBuilder()

	if astBuilder == nil {
		t.Fatal("GetASTBuilder returned nil")
	}

	if astBuilder.GetConfig().PackageName != "generated" {
		t.Errorf("Expected package name 'generated', got '%s'",
			astBuilder.GetConfig().PackageName)
	}
}

func TestMigratedGenerator_GetMigrationStatus(t *testing.T) {
	opts := &options.Options{
		PackagePrefix: "github.com/test",
		YAMLFiles:     []string{"test.yaml"},
	}

	migratedGen := NewMigratedGenerator(opts)
	status := migratedGen.GetMigrationStatus()

	expectedComponents := []string{"parameter_parsing", "schema_building", "handler_building", "validation"}
	for _, component := range expectedComponents {
		if !status[component] {
			t.Errorf("Expected %s to be migrated", component)
		}
	}
}

func TestCreateMigrationPlan(t *testing.T) {
	plan := CreateMigrationPlan()

	if plan == nil {
		t.Fatal("CreateMigrationPlan returned nil")
	}

	if len(plan.Steps) != 4 {
		t.Errorf("Expected 4 migration steps, got %d", len(plan.Steps))
	}

	// Check first step
	firstStep := plan.Steps[0]
	if firstStep.Phase != "Phase 1: Core Infrastructure" {
		t.Errorf("Expected first phase 'Phase 1: Core Infrastructure', got '%s'", firstStep.Phase)
	}

	if firstStep.Duration != "Week 1" {
		t.Errorf("Expected first duration 'Week 1', got '%s'", firstStep.Duration)
	}

	if len(firstStep.Tasks) == 0 {
		t.Error("Expected first step to have tasks")
	}
}

func TestMigrationModes(t *testing.T) {
	// Test migration mode constants
	if MigrationModeLegacy != 0 {
		t.Errorf("Expected MigrationModeLegacy to be 0, got %d", MigrationModeLegacy)
	}

	if MigrationModeHybrid != 1 {
		t.Errorf("Expected MigrationModeHybrid to be 1, got %d", MigrationModeHybrid)
	}

	if MigrationModeNew != 2 {
		t.Errorf("Expected MigrationModeNew to be 2, got %d", MigrationModeNew)
	}
}

func TestMigrationComparison(t *testing.T) {
	comparison := &MigrationComparison{
		Success:      true,
		LegacyError:  "",
		NewError:     "",
		LegacyOutput: map[string]string{"file1": "content1"},
		NewOutput:    map[string]string{"file1": "content1"},
		Differences:  []string{"No differences"},
	}

	if !comparison.Success {
		t.Error("Expected comparison to be successful")
	}

	if len(comparison.Differences) != 1 {
		t.Errorf("Expected 1 difference, got %d", len(comparison.Differences))
	}
}

func TestMigrationWrapper_Integration(t *testing.T) {
	// Test complete integration workflow
	opts := &options.Options{
		PackagePrefix: "github.com/test",
		YAMLFiles:     []string{"test.yaml"},
	}

	wrapper := NewMigrationWrapper(opts, MigrationModeNew)

	// Test configuration
	if wrapper.config.PackageName == "" {
		t.Errorf("Expected package name to be set, got empty string. Config: %+v", wrapper.config)
	}

	// Test generator access
	if wrapper.GetLegacyGenerator() == nil {
		t.Fatal("Legacy generator should not be nil")
	}

	if wrapper.GetMigratedGenerator() == nil {
		t.Fatal("Migrated generator should not be nil")
	}

	// Test AST builder access
	astBuilder := wrapper.GetASTBuilder()
	if astBuilder == nil {
		t.Fatal("AST builder should not be nil")
	}

	// Test mode switching
	wrapper.SetMigrationMode(MigrationModeLegacy)
	if wrapper.mode != MigrationModeLegacy {
		t.Error("Migration mode should have changed to legacy")
	}

	// Test validation with empty builder
	err := wrapper.ValidateMigration()
	if err == nil {
		t.Error("Validation should fail with empty builder")
	}

	// Add some content and test again
	astBuilder.AddImport("fmt")
	astBuilder.AddDeclaration(&ast.FuncDecl{
		Name: ast.NewIdent("Test"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{},
		},
		Body: &ast.BlockStmt{},
	})

	err = wrapper.ValidateMigration()
	if err != nil {
		t.Errorf("Validation should pass with content: %v", err)
	}
}
