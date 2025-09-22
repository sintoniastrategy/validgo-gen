package migration

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/jolfzverb/codegen/internal/generator"
	"github.com/jolfzverb/codegen/internal/generator/astbuilder"
	"github.com/jolfzverb/codegen/internal/generator/options"
)

// MigrationMode defines the migration strategy
type MigrationMode int

const (
	// MigrationModeLegacy uses only the old generator methods
	MigrationModeLegacy MigrationMode = iota
	// MigrationModeHybrid uses both old and new methods for comparison
	MigrationModeHybrid
	// MigrationModeNew uses only the new AST builder methods
	MigrationModeNew
)

// MigrationWrapper provides a bridge between old and new generation methods
type MigrationWrapper struct {
	legacyGenerator   *generator.Generator
	migratedGenerator *MigratedGenerator
	mode              MigrationMode
	config            MigrationConfig
}

// MigrationConfig holds configuration for the migration
type MigrationConfig struct {
	PackageName    string
	ImportPrefix   string
	UseNewBuilders bool
	ValidateOutput bool
	CompareOutputs bool
}

// NewMigrationWrapper creates a new migration wrapper
func NewMigrationWrapper(opts *options.Options, mode MigrationMode) *MigrationWrapper {
	legacyGen := generator.NewGenerator(opts)
	migratedGen := NewMigratedGenerator(opts)

	config := MigrationConfig{
		PackageName:    "generated", // Default package name
		ImportPrefix:   legacyGen.ImportPrefix,
		UseNewBuilders: mode != MigrationModeLegacy,
		ValidateOutput: true,
		CompareOutputs: mode == MigrationModeHybrid,
	}

	return &MigrationWrapper{
		legacyGenerator:   legacyGen,
		migratedGenerator: migratedGen,
		mode:              mode,
		config:            config,
	}
}

// Generate generates code using the specified migration mode
func (mw *MigrationWrapper) Generate(ctx context.Context) error {
	switch mw.mode {
	case MigrationModeLegacy:
		return mw.generateLegacy(ctx)
	case MigrationModeHybrid:
		return mw.generateHybrid(ctx)
	case MigrationModeNew:
		return mw.generateNew(ctx)
	default:
		return fmt.Errorf("unknown migration mode: %d", mw.mode)
	}
}

// generateLegacy generates code using the legacy generator
func (mw *MigrationWrapper) generateLegacy(ctx context.Context) error {
	slog.Info("Generating code using legacy generator")
	return mw.legacyGenerator.Generate(ctx)
}

// generateNew generates code using the migrated generator
func (mw *MigrationWrapper) generateNew(ctx context.Context) error {
	slog.Info("Generating code using migrated generator")
	return mw.migratedGenerator.GenerateWithMigration(ctx)
}

// generateHybrid generates code using both methods for comparison
func (mw *MigrationWrapper) generateHybrid(ctx context.Context) error {
	slog.Info("Generating code using hybrid mode (legacy + migrated)")

	// Generate with legacy method
	legacyErr := mw.legacyGenerator.Generate(ctx)
	if legacyErr != nil {
		slog.Error("Legacy generation failed", "error", legacyErr)
		return fmt.Errorf("legacy generation failed: %w", legacyErr)
	}

	// Generate with migrated method
	migratedErr := mw.migratedGenerator.GenerateWithMigration(ctx)
	if migratedErr != nil {
		slog.Error("Migrated generation failed", "error", migratedErr)
		return fmt.Errorf("migrated generation failed: %w", migratedErr)
	}

	slog.Info("Both generation methods completed successfully")
	return nil
}

// WriteToOutput writes the generated code to output writers
func (mw *MigrationWrapper) WriteToOutput(modelsOutput io.Writer, handlersOutput io.Writer) error {
	switch mw.mode {
	case MigrationModeLegacy:
		return mw.legacyGenerator.WriteToOutput(modelsOutput, handlersOutput)
	case MigrationModeNew:
		return mw.migratedGenerator.WriteToOutputWithMigration(modelsOutput, handlersOutput)
	case MigrationModeHybrid:
		// In hybrid mode, we write using the migrated generator
		// but we could also write both for comparison
		return mw.migratedGenerator.WriteToOutputWithMigration(modelsOutput, handlersOutput)
	default:
		return fmt.Errorf("unknown migration mode: %d", mw.mode)
	}
}

// GetLegacyGenerator returns the legacy generator
func (mw *MigrationWrapper) GetLegacyGenerator() *generator.Generator {
	return mw.legacyGenerator
}

// GetMigratedGenerator returns the migrated generator
func (mw *MigrationWrapper) GetMigratedGenerator() *MigratedGenerator {
	return mw.migratedGenerator
}

// GetASTBuilder returns the AST builder from the migrated generator
func (mw *MigrationWrapper) GetASTBuilder() *astbuilder.Builder {
	return mw.migratedGenerator.GetASTBuilder()
}

// SetMigrationMode changes the migration mode
func (mw *MigrationWrapper) SetMigrationMode(mode MigrationMode) {
	mw.mode = mode
	mw.config.UseNewBuilders = mode != MigrationModeLegacy
	mw.config.CompareOutputs = mode == MigrationModeHybrid
}

// GetMigrationStatus returns the current migration status
func (mw *MigrationWrapper) GetMigrationStatus() map[string]interface{} {
	status := map[string]interface{}{
		"mode":             mw.mode,
		"use_new_builders": mw.config.UseNewBuilders,
		"validate_output":  mw.config.ValidateOutput,
		"compare_outputs":  mw.config.CompareOutputs,
	}

	// Add migrated generator status
	if mw.migratedGenerator != nil {
		status["migrated_components"] = mw.migratedGenerator.GetMigrationStatus()
	}

	return status
}

// ValidateMigration validates that the migration is working correctly
func (mw *MigrationWrapper) ValidateMigration() error {
	if !mw.config.ValidateOutput {
		return nil
	}

	// Basic validation - check that the migrated generator has content
	astBuilder := mw.GetASTBuilder()
	if astBuilder == nil {
		return fmt.Errorf("AST builder is nil")
	}

	// Check that the builder has content
	if astBuilder.DeclarationCount() == 0 && astBuilder.StatementCount() == 0 {
		return fmt.Errorf("migration produced no output")
	}

	// Check for common issues
	file := astBuilder.BuildFile()
	if file == nil {
		return fmt.Errorf("migration produced invalid AST file")
	}

	// Validate package name
	if file.Name == nil || file.Name.Name != mw.config.PackageName {
		return fmt.Errorf("migration produced incorrect package name")
	}

	return nil
}

// CompareOutputs compares the output of legacy and migrated generation
func (mw *MigrationWrapper) CompareOutputs(ctx context.Context) (*MigrationComparison, error) {
	if mw.mode != MigrationModeHybrid {
		return nil, fmt.Errorf("comparison only available in hybrid mode")
	}

	comparison := &MigrationComparison{
		LegacyOutput: make(map[string]string),
		NewOutput:    make(map[string]string),
		Differences:  make([]string, 0),
	}

	// Generate with legacy method
	legacyErr := mw.legacyGenerator.Generate(ctx)
	if legacyErr != nil {
		comparison.LegacyError = legacyErr.Error()
	}

	// Generate with migrated method
	migratedErr := mw.migratedGenerator.GenerateWithMigration(ctx)
	if migratedErr != nil {
		comparison.NewError = migratedErr.Error()
	}

	// Compare outputs if both succeeded
	if legacyErr == nil && migratedErr == nil {
		comparison.Success = true
		comparison.Differences = append(comparison.Differences, "Output comparison not yet implemented")
	}

	return comparison, nil
}

// MigrationComparison holds the results of comparing legacy and migrated generation
type MigrationComparison struct {
	Success      bool              `json:"success"`
	LegacyError  string            `json:"legacy_error,omitempty"`
	NewError     string            `json:"new_error,omitempty"`
	LegacyOutput map[string]string `json:"legacy_output,omitempty"`
	NewOutput    map[string]string `json:"new_output,omitempty"`
	Differences  []string          `json:"differences,omitempty"`
}

// CreateMigrationPlan creates a step-by-step migration plan
func CreateMigrationPlan() *MigrationPlan {
	return &MigrationPlan{
		Steps: []MigrationStep{
			{
				Phase:    "Phase 1: Core Infrastructure",
				Duration: "Week 1",
				Tasks: []string{
					"✅ Implement core AST builder",
					"✅ Implement parameter parsing abstraction",
					"✅ Create migration wrapper",
					"✅ Set up hybrid mode for comparison",
				},
			},
			{
				Phase:    "Phase 2: Schema Migration",
				Duration: "Week 2",
				Tasks: []string{
					"✅ Implement schema building abstraction",
					"🔄 Migrate schema generation functions",
					"🔄 Test schema output compatibility",
					"🔄 Update validation logic",
				},
			},
			{
				Phase:    "Phase 3: Handler Migration",
				Duration: "Week 3",
				Tasks: []string{
					"✅ Implement handler building abstraction",
					"🔄 Migrate handler generation functions",
					"🔄 Test handler output compatibility",
					"🔄 Update routing logic",
				},
			},
			{
				Phase:    "Phase 4: Validation Migration",
				Duration: "Week 4",
				Tasks: []string{
					"✅ Implement validation building abstraction",
					"🔄 Migrate validation generation functions",
					"🔄 Test validation output compatibility",
					"🔄 Complete migration and cleanup",
				},
			},
		},
	}
}

// MigrationPlan defines the step-by-step migration plan
type MigrationPlan struct {
	Steps []MigrationStep `json:"steps"`
}

// MigrationStep represents a single step in the migration plan
type MigrationStep struct {
	Phase    string   `json:"phase"`
	Duration string   `json:"duration"`
	Tasks    []string `json:"tasks"`
}
