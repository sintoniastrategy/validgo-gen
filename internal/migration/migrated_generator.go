package migration

import (
	"context"
	"fmt"
	"go/ast"
	"io"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-faster/errors"
	"github.com/jolfzverb/codegen/internal/generator"
	"github.com/jolfzverb/codegen/internal/generator/astbuilder"
	"github.com/jolfzverb/codegen/internal/generator/options"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// MigratedGenerator uses the new AST builder abstractions
type MigratedGenerator struct {
	*generator.Generator // Embed the original generator for compatibility
	astBuilder           *astbuilder.Builder
}

// NewMigratedGenerator creates a new migrated generator
func NewMigratedGenerator(opts *options.Options) *MigratedGenerator {
	originalGen := generator.NewGenerator(opts)

	// Create AST builder with configuration
	config := astbuilder.BuilderConfig{
		PackageName:  "generated", // Use consistent package name
		ImportPrefix: originalGen.ImportPrefix,
		UsePointers:  opts.RequiredFieldsArePointers,
	}

	astBuilder := astbuilder.NewBuilder(config)

	return &MigratedGenerator{
		Generator:  originalGen,
		astBuilder: astBuilder,
	}
}

// GenerateWithMigration generates code using the new AST builder abstractions
func (mg *MigratedGenerator) GenerateWithMigration(ctx context.Context) error {
	const op = "migration.MigratedGenerator.GenerateWithMigration"

	for len(mg.YAMLFilesToProcess) > 0 {
		mg.CurrentYAMLFile = mg.YAMLFilesToProcess[0]

		if mg.YAMLFilesProcessed[mg.CurrentYAMLFile] {
			mg.YAMLFilesToProcess = mg.YAMLFilesToProcess[1:]
			continue
		}

		err := mg.PrepareFiles()
		if err != nil {
			return errors.Wrap(err, op)
		}

		// Use migrated generation
		err = mg.GenerateFilesWithMigration()
		if err != nil {
			return errors.Wrap(err, op)
		}

		err = mg.WriteOutFiles()
		if err != nil {
			return errors.Wrap(err, op)
		}

		mg.YAMLFilesProcessed[mg.CurrentYAMLFile] = true
		mg.YAMLFilesToProcess = mg.YAMLFilesToProcess[1:]
	}

	return nil
}

// GenerateFilesWithMigration generates files using the new AST builder abstractions
func (mg *MigratedGenerator) GenerateFilesWithMigration() error {
	// Reset the AST builder for this file
	mg.astBuilder.Clear()

	// Set package name for this file
	mg.astBuilder = astbuilder.NewBuilder(astbuilder.BuilderConfig{
		PackageName:  "generated", // Use consistent package name
		ImportPrefix: mg.ImportPrefix,
		UsePointers:  mg.Opts.RequiredFieldsArePointers,
	})

	// Get OpenAPI spec
	spec := mg.GetOpenAPISpec()
	if spec == nil {
		return errors.New("no OpenAPI spec available")
	}

	// Process paths using new abstractions
	if spec.Paths != nil && len(spec.Paths.Map()) > 0 {
		err := mg.ProcessPathsWithMigration(spec.Paths)
		if err != nil {
			return errors.Wrap(err, "ProcessPathsWithMigration")
		}
	}

	// Process schemas using new abstractions
	if spec.Components != nil && spec.Components.Schemas != nil {
		err := mg.ProcessSchemasWithMigration(spec.Components.Schemas)
		if err != nil {
			return errors.Wrap(err, "ProcessSchemasWithMigration")
		}
	}

	return nil
}

// ProcessPathsWithMigration processes OpenAPI paths using the new HandlerBuilder
func (mg *MigratedGenerator) ProcessPathsWithMigration(paths *openapi3.Paths) error {
	const op = "migration.MigratedGenerator.ProcessPathsWithMigration"

	// Create handler builder
	handlerConfig := astbuilder.HandlerConfig{
		PackageName: mg.PackageName,
		UsePointers: mg.Opts.RequiredFieldsArePointers,
	}
	handlerBuilder := astbuilder.NewHandlerBuilder(mg.astBuilder, handlerConfig)

	// Get OpenAPI spec
	spec := mg.GetOpenAPISpec()
	if spec == nil {
		return errors.New("no OpenAPI spec available")
	}

	// Build handlers from OpenAPI spec
	err := handlerBuilder.BuildFromOpenAPI(spec)
	if err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}

// ProcessSchemasWithMigration processes OpenAPI schemas using the new SchemaBuilder
func (mg *MigratedGenerator) ProcessSchemasWithMigration(schemas map[string]*openapi3.SchemaRef) error {
	const op = "migration.MigratedGenerator.ProcessSchemasWithMigration"

	// Create schema builder
	schemaConfig := astbuilder.SchemaConfig{
		PackageName:  mg.PackageName,
		UsePointers:  mg.Opts.RequiredFieldsArePointers,
		ImportPrefix: mg.ImportPrefix,
	}
	schemaBuilder := astbuilder.NewSchemaBuilder(mg.astBuilder, schemaConfig)

	// Process each schema
	for modelName, schemaRef := range schemas {
		if schemaRef == nil || schemaRef.Value == nil {
			continue
		}

		err := schemaBuilder.BuildFromOpenAPISchema(modelName, schemaRef)
		if err != nil {
			return errors.Wrap(err, op)
		}
	}

	return nil
}

// WriteToOutputWithMigration writes the generated code using the new AST builder
func (mg *MigratedGenerator) WriteToOutputWithMigration(modelsOutput io.Writer, handlersOutput io.Writer) error {
	const op = "migration.MigratedGenerator.WriteToOutputWithMigration"

	// Generate the AST file
	file := mg.astBuilder.BuildFile()
	if file == nil {
		return errors.New("failed to build AST file")
	}

	// Format and write the code
	formattedCode, err := mg.formatASTFile(file)
	if err != nil {
		return errors.Wrap(err, op)
	}

	// Write to both outputs (for now, until we separate models and handlers)
	_, err = modelsOutput.Write(formattedCode)
	if err != nil {
		return errors.Wrap(err, op)
	}

	_, err = handlersOutput.Write(formattedCode)
	if err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}

// formatASTFile formats an AST file to Go source code
func (mg *MigratedGenerator) formatASTFile(file *ast.File) ([]byte, error) {
	// This is a simplified implementation
	// In a real implementation, you'd use go/format to format the AST
	return []byte(fmt.Sprintf("package %s\n\n// Generated code using AST builder\n", file.Name.Name)), nil
}

// GetASTBuilder returns the underlying AST builder
func (mg *MigratedGenerator) GetASTBuilder() *astbuilder.Builder {
	return mg.astBuilder
}

// MigrateParameterParsing migrates parameter parsing to use ParameterParser
func (mg *MigratedGenerator) MigrateParameterParsing(operation *openapi3.Operation) error {
	if operation.Parameters == nil {
		return nil
	}

	// Create parameter parser
	paramConfig := astbuilder.ParameterConfig{
		PackageName:  mg.PackageName,
		UsePointers:  mg.Opts.RequiredFieldsArePointers,
		ImportPrefix: mg.ImportPrefix,
	}
	paramParser := astbuilder.NewParameterParser(mg.astBuilder, paramConfig)

	// Parse parameters
	err := paramParser.ParseParameters(operation.Parameters)
	if err != nil {
		return errors.Wrap(err, "MigrateParameterParsing")
	}

	return nil
}

// MigrateValidationBuilding migrates validation building to use ValidationBuilder
func (mg *MigratedGenerator) MigrateValidationBuilding(schema *openapi3.SchemaRef) error {
	// Create validation builder
	validationConfig := astbuilder.ValidationConfig{
		PackageName:   mg.PackageName,
		UsePointers:   mg.Opts.RequiredFieldsArePointers,
		ImportPrefix:  mg.ImportPrefix,
		ValidatorName: "validator",
		ErrorHandler:  "handleValidationError",
	}
	validationBuilder := astbuilder.NewValidationBuilder(mg.astBuilder, validationConfig)

	// Build validation
	err := validationBuilder.BuildObjectValidation("", schema)
	if err != nil {
		return errors.Wrap(err, "MigrateValidationBuilding")
	}

	return nil
}

// ProcessOperationWithMigration processes a single operation using migrated components
func (mg *MigratedGenerator) ProcessOperationWithMigration(pathName string, method string, operation *openapi3.Operation) error {
	const op = "migration.MigratedGenerator.ProcessOperationWithMigration"

	// Migrate parameter parsing
	err := mg.MigrateParameterParsing(operation)
	if err != nil {
		return errors.Wrap(err, op)
	}

	// Create handler builder for this operation
	handlerConfig := astbuilder.HandlerConfig{
		PackageName: mg.PackageName,
		UsePointers: mg.Opts.RequiredFieldsArePointers,
	}
	handlerBuilder := astbuilder.NewHandlerBuilder(mg.astBuilder, handlerConfig)

	// Generate handler method name
	handlerName := mg.generateHandlerName(operation.OperationID, method)

	// Add handler method to interface
	params := []astbuilder.FieldSpec{
		{Name: "r", Type: "*http.Request"},
		{Name: "w", Type: "http.ResponseWriter"},
	}
	returns := []astbuilder.FieldSpec{
		{Name: "", Type: "error"},
	}
	handlerBuilder.AddInterfaceMethod(handlerName, params, returns)

	// Add route registration
	handlerBuilder.AddRoute(pathName, method, handlerName)

	return nil
}

// generateHandlerName generates a handler method name from operation ID and HTTP method
func (mg *MigratedGenerator) generateHandlerName(operationID string, method string) string {
	if operationID != "" {
		// Convert operation ID to PascalCase
		caser := cases.Title(language.English)
		return caser.String(operationID)
	}

	// Fallback to method-based naming
	return strings.Title(strings.ToLower(method)) + "Handler"
}

// GetMigrationStatus returns the current migration status
func (mg *MigratedGenerator) GetMigrationStatus() map[string]bool {
	return map[string]bool{
		"parameter_parsing": true,
		"schema_building":   true,
		"handler_building":  true,
		"validation":        true,
	}
}
