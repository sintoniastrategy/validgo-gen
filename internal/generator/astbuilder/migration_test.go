package astbuilder

import (
	"context"
	"go/ast"
	"strings"
	"testing"

	"github.com/jolfzverb/codegen/internal/generator"
	"github.com/jolfzverb/codegen/internal/generator/options"
)

func TestNewMigrationWrapper(t *testing.T) {
	// Create migration config
	config := MigrationConfig{
		PackageName:    "test",
		ImportPrefix:   "github.com/test",
		UseNewBuilders: true,
		MigrationMode:  MigrationModeNew,
		ValidateOutput: true,
	}

	// Create migration wrapper
	wrapper := NewMigrationWrapper(nil, config)
	if wrapper == nil {
		t.Fatal("NewMigrationWrapper returned nil")
	}

	// Test configuration
	if wrapper.config.PackageName != "test" {
		t.Errorf("Expected package name 'test', got '%s'", wrapper.config.PackageName)
	}

	if wrapper.config.MigrationMode != MigrationModeNew {
		t.Errorf("Expected migration mode %d, got %d", MigrationModeNew, wrapper.config.MigrationMode)
	}

	// Test builder creation
	if wrapper.builder == nil {
		t.Fatal("Builder is nil")
	}

	if wrapper.builder.GetConfig().PackageName != "test" {
		t.Errorf("Expected builder package name 'test', got '%s'", wrapper.builder.GetConfig().PackageName)
	}
}

func TestNewMigrationWrapperFromOptions(t *testing.T) {
	// Create test options
	opts := &options.Options{
		PackagePrefix: "github.com/testpackage",
		YAMLFiles:     []string{"test.yaml"},
	}

	// Create migration wrapper from options
	wrapper := NewMigrationWrapperFromOptions(opts, MigrationModeHybrid)
	if wrapper == nil {
		t.Fatal("NewMigrationWrapperFromOptions returned nil")
	}

	// Test configuration
	if wrapper.config.PackageName != "generated" {
		t.Errorf("Expected package name 'generated', got '%s'", wrapper.config.PackageName)
	}

	if wrapper.config.MigrationMode != MigrationModeHybrid {
		t.Errorf("Expected migration mode %d, got %d", MigrationModeHybrid, wrapper.config.MigrationMode)
	}

	// Test generator creation
	if wrapper.generator == nil {
		t.Fatal("Generator is nil")
	}
}

func TestMigrationWrapper_GetBuilder(t *testing.T) {
	config := MigrationConfig{
		PackageName:   "test",
		MigrationMode: MigrationModeNew,
	}

	wrapper := NewMigrationWrapper(nil, config)
	builder := wrapper.GetBuilder()

	if builder == nil {
		t.Fatal("GetBuilder returned nil")
	}

	if builder.GetConfig().PackageName != "test" {
		t.Errorf("Expected package name 'test', got '%s'", builder.GetConfig().PackageName)
	}
}

func TestMigrationWrapper_GetGenerator(t *testing.T) {
	opts := &options.Options{
		PackagePrefix: "github.com/test",
	}
	gen := generator.NewGenerator(opts)

	config := MigrationConfig{
		PackageName:   "test",
		MigrationMode: MigrationModeNew,
	}

	wrapper := NewMigrationWrapper(gen, config)
	retrievedGen := wrapper.GetGenerator()

	if retrievedGen != gen {
		t.Error("GetGenerator returned different generator instance")
	}
}

func TestMigrationWrapper_GetConfig(t *testing.T) {
	config := MigrationConfig{
		PackageName:    "test",
		ImportPrefix:   "github.com/test",
		UseNewBuilders: true,
		MigrationMode:  MigrationModeNew,
		ValidateOutput: true,
	}

	wrapper := NewMigrationWrapper(nil, config)
	retrievedConfig := wrapper.GetConfig()

	if retrievedConfig.PackageName != config.PackageName {
		t.Errorf("Expected package name '%s', got '%s'", config.PackageName, retrievedConfig.PackageName)
	}

	if retrievedConfig.MigrationMode != config.MigrationMode {
		t.Errorf("Expected migration mode %d, got %d", config.MigrationMode, retrievedConfig.MigrationMode)
	}
}

func TestMigrationWrapper_SetMigrationMode(t *testing.T) {
	config := MigrationConfig{
		PackageName:   "test",
		MigrationMode: MigrationModeLegacy,
	}

	wrapper := NewMigrationWrapper(nil, config)

	// Test initial mode
	if wrapper.config.MigrationMode != MigrationModeLegacy {
		t.Errorf("Expected initial mode %d, got %d", MigrationModeLegacy, wrapper.config.MigrationMode)
	}

	// Change mode
	wrapper.SetMigrationMode(MigrationModeNew)

	if wrapper.config.MigrationMode != MigrationModeNew {
		t.Errorf("Expected mode %d after change, got %d", MigrationModeNew, wrapper.config.MigrationMode)
	}
}

func TestMigrationWrapper_MigrateParameterParsing(t *testing.T) {
	config := MigrationConfig{
		PackageName:   "test",
		MigrationMode: MigrationModeNew,
	}

	wrapper := NewMigrationWrapper(nil, config)

	// Test migration in legacy mode (should skip)
	wrapper.SetMigrationMode(MigrationModeLegacy)
	err := wrapper.MigrateParameterParsing()
	if err != nil {
		t.Errorf("Migration in legacy mode should not fail: %v", err)
	}

	// Test migration in new mode (will fail without proper spec setup)
	wrapper.SetMigrationMode(MigrationModeNew)
	err = wrapper.MigrateParameterParsing()
	if err == nil {
		t.Error("Expected error in new mode without proper spec setup")
	}
}

func TestMigrationWrapper_MigrateSchemaBuilding(t *testing.T) {
	config := MigrationConfig{
		PackageName:   "test",
		MigrationMode: MigrationModeNew,
	}

	wrapper := NewMigrationWrapper(nil, config)

	// Test migration in legacy mode (should skip)
	wrapper.SetMigrationMode(MigrationModeLegacy)
	err := wrapper.MigrateSchemaBuilding()
	if err != nil {
		t.Errorf("Migration in legacy mode should not fail: %v", err)
	}

	// Test migration in new mode (will fail without proper spec setup)
	wrapper.SetMigrationMode(MigrationModeNew)
	err = wrapper.MigrateSchemaBuilding()
	if err == nil {
		t.Error("Expected error in new mode without proper spec setup")
	}
}

func TestMigrationWrapper_MigrateHandlerBuilding(t *testing.T) {
	config := MigrationConfig{
		PackageName:   "test",
		MigrationMode: MigrationModeNew,
	}

	wrapper := NewMigrationWrapper(nil, config)

	// Test migration in legacy mode (should skip)
	wrapper.SetMigrationMode(MigrationModeLegacy)
	err := wrapper.MigrateHandlerBuilding()
	if err != nil {
		t.Errorf("Migration in legacy mode should not fail: %v", err)
	}

	// Test migration in new mode (will fail without proper spec setup)
	wrapper.SetMigrationMode(MigrationModeNew)
	err = wrapper.MigrateHandlerBuilding()
	if err == nil {
		t.Error("Expected error in new mode without proper spec setup")
	}
}

func TestMigrationWrapper_MigrateValidationBuilding(t *testing.T) {
	config := MigrationConfig{
		PackageName:   "test",
		MigrationMode: MigrationModeNew,
	}

	wrapper := NewMigrationWrapper(nil, config)

	// Test migration in legacy mode (should skip)
	wrapper.SetMigrationMode(MigrationModeLegacy)
	err := wrapper.MigrateValidationBuilding()
	if err != nil {
		t.Errorf("Migration in legacy mode should not fail: %v", err)
	}

	// Test migration in new mode (will fail without proper spec setup)
	wrapper.SetMigrationMode(MigrationModeNew)
	err = wrapper.MigrateValidationBuilding()
	if err == nil {
		t.Error("Expected error in new mode without proper spec setup")
	}
}

func TestMigrationWrapper_MigrateAll(t *testing.T) {
	config := MigrationConfig{
		PackageName:   "test",
		MigrationMode: MigrationModeNew,
	}

	wrapper := NewMigrationWrapper(nil, config)

	// Test migration in legacy mode (should skip)
	wrapper.SetMigrationMode(MigrationModeLegacy)
	err := wrapper.MigrateAll()
	if err != nil {
		t.Errorf("Migration in legacy mode should not fail: %v", err)
	}

	// Test migration in new mode (will fail without proper spec setup)
	wrapper.SetMigrationMode(MigrationModeNew)
	err = wrapper.MigrateAll()
	if err == nil {
		t.Error("Expected error in new mode without proper spec setup")
	}
}

func TestMigrationWrapper_ValidateMigration(t *testing.T) {
	config := MigrationConfig{
		PackageName:    "test",
		MigrationMode:  MigrationModeNew,
		ValidateOutput: true,
	}

	wrapper := NewMigrationWrapper(nil, config)

	// Test validation with empty builder (should fail)
	err := wrapper.ValidateMigration()
	if err == nil {
		t.Error("Expected validation to fail with empty builder")
	}

	// Add some content to the builder (but not enough to pass validation)
	wrapper.builder.AddImport("fmt")

	// Test validation with content (should still fail due to no declarations)
	err = wrapper.ValidateMigration()
	if err == nil {
		t.Error("Expected validation to fail with no declarations")
	}

	// Create a proper file structure
	wrapper.builder.Clear()
	wrapper.builder.AddImport("fmt")
	wrapper.builder.AddDeclaration(&ast.FuncDecl{
		Name: ast.NewIdent("Test"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{},
		},
		Body: &ast.BlockStmt{},
	})

	// Test validation with proper content (should pass now)
	err = wrapper.ValidateMigration()
	if err != nil {
		t.Errorf("Validation should pass with proper content: %v", err)
	}
}

func TestMigrationWrapper_WriteToOutput(t *testing.T) {
	config := MigrationConfig{
		PackageName:   "test",
		MigrationMode: MigrationModeNew, // Skip legacy mode to avoid panic
	}

	wrapper := NewMigrationWrapper(nil, config)

	// Test with new mode
	wrapper.builder.AddImport("fmt")
	wrapper.builder.AddDeclaration(&ast.FuncDecl{
		Name: ast.NewIdent("Test"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{},
		},
		Body: &ast.BlockStmt{},
	})

	var modelsOutput, handlersOutput strings.Builder
	err := wrapper.WriteToOutput(&modelsOutput, &handlersOutput)
	if err != nil {
		t.Errorf("WriteToOutput failed in new mode: %v", err)
	}

	// Check that output was written
	if modelsOutput.Len() == 0 {
		t.Error("Expected models output to be written")
	}

	if handlersOutput.Len() == 0 {
		t.Error("Expected handlers output to be written")
	}
}

func TestMigrationWrapper_CompareOutputs(t *testing.T) {
	config := MigrationConfig{
		PackageName:   "test",
		MigrationMode: MigrationModeHybrid,
	}

	opts := &options.Options{
		PackagePrefix: "github.com/test",
	}
	gen := generator.NewGenerator(opts)

	wrapper := NewMigrationWrapper(gen, config)

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
	config := MigrationConfig{
		PackageName:    "testpackage",
		ImportPrefix:   "github.com/testpackage",
		UseNewBuilders: true,
		MigrationMode:  MigrationModeNew,
		ValidateOutput: true,
	}

	wrapper := NewMigrationWrapper(nil, config)

	// Test configuration
	if wrapper.GetConfig().PackageName != "testpackage" {
		t.Errorf("Expected package name 'testpackage', got '%s'", wrapper.GetConfig().PackageName)
	}

	// Test builder access
	builder := wrapper.GetBuilder()
	if builder == nil {
		t.Fatal("Builder should not be nil")
	}

	// Test mode switching
	wrapper.SetMigrationMode(MigrationModeLegacy)
	if wrapper.GetConfig().MigrationMode != MigrationModeLegacy {
		t.Error("Migration mode should have changed to legacy")
	}

	// Test validation with empty builder
	err := wrapper.ValidateMigration()
	if err == nil {
		t.Error("Validation should fail with empty builder")
	}

	// Add some content and test again
	builder.AddImport("fmt")
	builder.AddDeclaration(&ast.FuncDecl{
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
