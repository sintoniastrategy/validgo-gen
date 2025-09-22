package main

import (
	"go/ast"
	"go/format"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/jolfzverb/codegen/internal/generator/astbuilder"
	"github.com/jolfzverb/codegen/internal/generator/options"
)

func main() {
	opts, err := options.GetOptions()
	if err != nil {
		log.Fatal("Failed to get options:", err)
	}

	log.Println("🚀 Generating code using new AST builder abstractions...")

	// Process each YAML file
	for _, yamlFile := range opts.YAMLFiles {
		log.Printf("📄 Processing file: %s", yamlFile)

		err := processYAMLFile(yamlFile, opts)
		if err != nil {
			log.Fatalf("Failed to process %s: %v", yamlFile, err)
		}

		log.Printf("✅ Successfully processed: %s", yamlFile)
	}

	log.Println("🎉 Code generation completed successfully!")
}

func processYAMLFile(yamlFile string, opts *options.Options) error {
	// Load OpenAPI spec
	loader := openapi3.NewLoader()
	spec, err := loader.LoadFromFile(yamlFile)
	if err != nil {
		return err
	}

	// Generate package name from file
	packageName := generatePackageName(yamlFile)

	// Create AST builder
	config := astbuilder.BuilderConfig{
		PackageName:  packageName,
		ImportPrefix: opts.PackagePrefix,
		UsePointers:  opts.RequiredFieldsArePointers,
	}
	builder := astbuilder.NewBuilder(config)

	// Add standard imports
	builder.AddImport("net/http")
	builder.AddImport("github.com/go-chi/chi/v5")
	builder.AddImport("github.com/go-playground/validator/v10")

	// Generate schemas
	err = generateSchemas(spec, builder)
	if err != nil {
		return err
	}

	// Generate handlers
	err = generateHandlers(spec, builder)
	if err != nil {
		return err
	}

	// Build and format the AST
	file := builder.BuildFile()
	if file == nil {
		return err
	}

	// Format the code
	formattedCode, err := formatASTFile(file)
	if err != nil {
		return err
	}

	// Create output directory
	outputDir := filepath.Join(opts.DirPrefix, "generated", packageName)
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		return err
	}

	// Write generated code
	outputFile := filepath.Join(outputDir, "generated.go")
	err = os.WriteFile(outputFile, formattedCode, 0644)
	if err != nil {
		return err
	}

	log.Printf("📝 Generated: %s", outputFile)
	return nil
}

func generateSchemas(spec *openapi3.T, builder *astbuilder.Builder) error {
	if spec.Components == nil || spec.Components.Schemas == nil {
		return nil
	}

	// Create schema builder
	schemaConfig := astbuilder.SchemaConfig{
		PackageName:  builder.GetConfig().PackageName,
		UsePointers:  builder.GetConfig().UsePointers,
		ImportPrefix: builder.GetConfig().ImportPrefix,
	}
	schemaBuilder := astbuilder.NewSchemaBuilder(builder, schemaConfig)

	// Process each schema
	for modelName, schemaRef := range spec.Components.Schemas {
		if schemaRef == nil || schemaRef.Value == nil {
			continue
		}

		err := schemaBuilder.BuildFromOpenAPISchema(modelName, schemaRef)
		if err != nil {
			return err
		}
	}

	return nil
}

func generateHandlers(spec *openapi3.T, builder *astbuilder.Builder) error {
	if spec.Paths == nil {
		return nil
	}

	// Create handler builder
	handlerConfig := astbuilder.HandlerConfig{
		PackageName: builder.GetConfig().PackageName,
		UsePointers: builder.GetConfig().UsePointers,
	}
	handlerBuilder := astbuilder.NewHandlerBuilder(builder, handlerConfig)

	// Build handlers from OpenAPI spec
	err := handlerBuilder.BuildFromOpenAPI(spec)
	if err != nil {
		return err
	}

	return nil
}

func generatePackageName(yamlFile string) string {
	// Extract filename without extension
	filename := filepath.Base(yamlFile)
	name := strings.TrimSuffix(filename, ".yaml")
	name = strings.TrimSuffix(name, ".yml")

	// Convert to valid Go package name
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, "-", "")
	name = strings.ReplaceAll(name, "_", "")

	if name == "" {
		name = "generated"
	}

	return name
}

func formatASTFile(file *ast.File) ([]byte, error) {
	var buf strings.Builder
	fset := token.NewFileSet()

	err := format.Node(&buf, fset, file)
	if err != nil {
		return nil, err
	}

	return []byte(buf.String()), nil
}
