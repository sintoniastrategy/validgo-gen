package astbuilder

import (
	"context"
	"fmt"
	"go/format"
	"go/token"
	"io"
	"strings"

	"github.com/jolfzverb/codegen/internal/generator"
	"github.com/jolfzverb/codegen/internal/generator/options"
)

// MigrationWrapper provides a bridge between the old Generator and new AST builder
type MigrationWrapper struct {
	generator *generator.Generator
	builder   *Builder
	config    MigrationConfig
}

// MigrationConfig holds configuration for the migration process
type MigrationConfig struct {
	PackageName    string
	ImportPrefix   string
	UseNewBuilders bool // Flag to enable/disable new builder usage
	MigrationMode  MigrationMode
	ValidateOutput bool // Flag to validate generated output
}

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

// NewMigrationWrapper creates a new migration wrapper
func NewMigrationWrapper(gen *generator.Generator, config MigrationConfig) *MigrationWrapper {
	builderConfig := BuilderConfig{
		PackageName:  config.PackageName,
		ImportPrefix: config.ImportPrefix,
	}

	return &MigrationWrapper{
		generator: gen,
		builder:   NewBuilder(builderConfig),
		config:    config,
	}
}

// MigrateParameterParsing migrates parameter parsing to use the new ParameterParser
func (m *MigrationWrapper) MigrateParameterParsing() error {
	if m.config.MigrationMode == MigrationModeLegacy {
		return nil // Skip migration in legacy mode
	}

	// Note: In a real implementation, you'd need to add a getter method on Generator
	// For now, we'll return an error indicating the spec needs to be set
	return fmt.Errorf("OpenAPI spec access not yet implemented - needs getter method on Generator")
}

// MigrateSchemaBuilding migrates schema building to use the new SchemaBuilder
func (m *MigrationWrapper) MigrateSchemaBuilding() error {
	if m.config.MigrationMode == MigrationModeLegacy {
		return nil // Skip migration in legacy mode
	}

	// Note: In a real implementation, you'd need to add a getter method on Generator
	// For now, we'll return an error indicating the spec needs to be set
	return fmt.Errorf("OpenAPI spec access not yet implemented - needs getter method on Generator")
}

// MigrateHandlerBuilding migrates handler building to use the new HandlerBuilder
func (m *MigrationWrapper) MigrateHandlerBuilding() error {
	if m.config.MigrationMode == MigrationModeLegacy {
		return nil // Skip migration in legacy mode
	}

	// Note: In a real implementation, you'd need to add a getter method on Generator
	// For now, we'll return an error indicating the spec needs to be set
	return fmt.Errorf("OpenAPI spec access not yet implemented - needs getter method on Generator")
}

// MigrateValidationBuilding migrates validation building to use the new ValidationBuilder
func (m *MigrationWrapper) MigrateValidationBuilding() error {
	if m.config.MigrationMode == MigrationModeLegacy {
		return nil // Skip migration in legacy mode
	}

	// Note: In a real implementation, you'd need to add a getter method on Generator
	// For now, we'll return an error indicating the spec needs to be set
	return fmt.Errorf("OpenAPI spec access not yet implemented - needs getter method on Generator")
}

// MigrateAll performs a complete migration of all components
func (m *MigrationWrapper) MigrateAll() error {
	// Migrate in dependency order
	if err := m.MigrateParameterParsing(); err != nil {
		return fmt.Errorf("parameter parsing migration failed: %w", err)
	}

	if err := m.MigrateSchemaBuilding(); err != nil {
		return fmt.Errorf("schema building migration failed: %w", err)
	}

	if err := m.MigrateHandlerBuilding(); err != nil {
		return fmt.Errorf("handler building migration failed: %w", err)
	}

	if err := m.MigrateValidationBuilding(); err != nil {
		return fmt.Errorf("validation building migration failed: %w", err)
	}

	return nil
}

// GenerateWithMigration generates code using the migration wrapper
func (m *MigrationWrapper) GenerateWithMigration(ctx context.Context) error {
	switch m.config.MigrationMode {
	case MigrationModeLegacy:
		// Use only the old generator
		return m.generator.Generate(ctx)

	case MigrationModeHybrid:
		// Use both old and new methods for comparison
		// First generate with old method
		oldErr := m.generator.Generate(ctx)

		// Then generate with new method
		newErr := m.MigrateAll()

		// Return old error if it exists, otherwise return new error
		if oldErr != nil {
			return fmt.Errorf("legacy generation failed: %w", oldErr)
		}
		if newErr != nil {
			return fmt.Errorf("new generation failed: %w", newErr)
		}

		return nil

	case MigrationModeNew:
		// Use only the new AST builder methods
		return m.MigrateAll()

	default:
		return fmt.Errorf("unknown migration mode: %d", m.config.MigrationMode)
	}
}

// WriteToOutput writes the generated code to output writers
func (m *MigrationWrapper) WriteToOutput(modelsOutput io.Writer, handlersOutput io.Writer) error {
	if m.config.MigrationMode == MigrationModeLegacy {
		// Use the old generator's WriteToOutput method
		return m.generator.WriteToOutput(modelsOutput, handlersOutput)
	}

	// Use the new AST builder to generate output
	file := m.builder.BuildFile()

	// Format the AST
	fset := token.NewFileSet()
	var buf strings.Builder
	err := format.Node(&buf, fset, file)
	if err != nil {
		return fmt.Errorf("failed to format generated code: %w", err)
	}
	formatted := buf.String()

	// Write to both outputs (for now, until we separate models and handlers)
	_, err = modelsOutput.Write([]byte(formatted))
	if err != nil {
		return fmt.Errorf("failed to write models output: %w", err)
	}

	_, err = handlersOutput.Write([]byte(formatted))
	if err != nil {
		return fmt.Errorf("failed to write handlers output: %w", err)
	}

	return nil
}

// CompareOutputs compares the output of old and new generation methods
func (m *MigrationWrapper) CompareOutputs(ctx context.Context) (*MigrationComparison, error) {
	if m.config.MigrationMode != MigrationModeHybrid {
		return nil, fmt.Errorf("comparison only available in hybrid mode")
	}

	comparison := &MigrationComparison{
		LegacyOutput: make(map[string]string),
		NewOutput:    make(map[string]string),
		Differences:  make([]string, 0),
	}

	// Generate with legacy method
	legacyErr := m.generator.Generate(ctx)
	if legacyErr != nil {
		comparison.LegacyError = legacyErr.Error()
	}

	// Generate with new method
	newErr := m.MigrateAll()
	if newErr != nil {
		comparison.NewError = newErr.Error()
	}

	// Compare outputs if both succeeded
	if legacyErr == nil && newErr == nil {
		// This is a simplified comparison - in practice, you'd want to
		// capture the actual output and compare it
		comparison.Success = true
		comparison.Differences = append(comparison.Differences, "Output comparison not yet implemented")
	}

	return comparison, nil
}

// MigrationComparison holds the results of comparing old and new generation methods
type MigrationComparison struct {
	Success      bool              `json:"success"`
	LegacyError  string            `json:"legacy_error,omitempty"`
	NewError     string            `json:"new_error,omitempty"`
	LegacyOutput map[string]string `json:"legacy_output,omitempty"`
	NewOutput    map[string]string `json:"new_output,omitempty"`
	Differences  []string          `json:"differences,omitempty"`
}

// GetBuilder returns the underlying AST builder
func (m *MigrationWrapper) GetBuilder() *Builder {
	return m.builder
}

// GetGenerator returns the underlying generator
func (m *MigrationWrapper) GetGenerator() *generator.Generator {
	return m.generator
}

// GetConfig returns the migration configuration
func (m *MigrationWrapper) GetConfig() MigrationConfig {
	return m.config
}

// SetMigrationMode changes the migration mode
func (m *MigrationWrapper) SetMigrationMode(mode MigrationMode) {
	m.config.MigrationMode = mode
}

// ValidateMigration validates that the migration was successful
func (m *MigrationWrapper) ValidateMigration() error {
	if !m.config.ValidateOutput {
		return nil
	}

	// Basic validation - check that the builder has content
	if m.builder.DeclarationCount() == 0 && m.builder.StatementCount() == 0 {
		return fmt.Errorf("migration produced no output")
	}

	// Check for common issues
	file := m.builder.BuildFile()
	if file == nil {
		return fmt.Errorf("migration produced invalid AST file")
	}

	// Validate package name
	if file.Name == nil || file.Name.Name != m.config.PackageName {
		return fmt.Errorf("migration produced incorrect package name")
	}

	return nil
}

// CreateMigrationPlan creates a step-by-step migration plan
func CreateMigrationPlan() *MigrationPlan {
	return &MigrationPlan{
		Steps: []MigrationStep{
			{
				Phase:    "Phase 1: Core Infrastructure",
				Duration: "Week 1",
				Tasks: []string{
					"Implement core AST builder",
					"Implement parameter parsing abstraction",
					"Create migration wrapper",
					"Set up hybrid mode for comparison",
				},
			},
			{
				Phase:    "Phase 2: Schema Migration",
				Duration: "Week 2",
				Tasks: []string{
					"Implement schema building abstraction",
					"Migrate schema generation functions",
					"Test schema output compatibility",
					"Update validation logic",
				},
			},
			{
				Phase:    "Phase 3: Handler Migration",
				Duration: "Week 3",
				Tasks: []string{
					"Implement handler building abstraction",
					"Migrate handler generation functions",
					"Test handler output compatibility",
					"Update routing logic",
				},
			},
			{
				Phase:    "Phase 4: Validation Migration",
				Duration: "Week 4",
				Tasks: []string{
					"Implement validation building abstraction",
					"Migrate validation generation functions",
					"Test validation output compatibility",
					"Complete migration and cleanup",
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

// Helper function to create a migration wrapper from options
func NewMigrationWrapperFromOptions(opts *options.Options, migrationMode MigrationMode) *MigrationWrapper {
	// Create a generator instance
	gen := generator.NewGenerator(opts)

	// Create migration config
	config := MigrationConfig{
		PackageName:    "generated", // Default package name
		ImportPrefix:   opts.PackagePrefix,
		UseNewBuilders: migrationMode != MigrationModeLegacy,
		MigrationMode:  migrationMode,
		ValidateOutput: true,
	}

	return NewMigrationWrapper(gen, config)
}
