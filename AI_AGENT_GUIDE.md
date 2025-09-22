# AI Agent Guide for jolfzverb-codegen

This guide helps AI agents understand and efficiently work with the jolfzverb-codegen OpenAPI code generator.

## Project Overview

**jolfzverb-codegen** is an OpenAPI 3.0 code generator that creates type-safe Go code from OpenAPI specifications. It generates HTTP handlers, data models, and validation code for building robust REST APIs in Go.

## Core Architecture

### 1. Main Components

```
cmd/generate.go                    # CLI entry point
internal/generator/                # Core generation logic
├── api.go                        # Main Generator struct and orchestration
├── generator.go                  # Utility functions and ref parsing
├── generatehandlers.go           # Handler generation logic
├── generateschemas.go            # Schema generation logic
├── handlers.go                   # Handler AST building utilities
├── handlers2.go                  # Handler method implementations
├── schemas.go                    # Schema AST building utilities
├── nameutils.go                  # Go identifier formatting
├── utils.go                      # AST helper functions
└── options/options.go            # Command-line options
```

### 2. Data Flow

```
OpenAPI YAML → Generator → AST → Go Code
     ↓              ↓        ↓       ↓
  Parse &      Process    Build   Write to
  Validate     Schemas    AST     Files
```

### 3. Key Data Structures

#### Generator (main orchestrator)
```go
type Generator struct {
    Opts *options.Options           // Command-line options
    SchemasFile *SchemasFile        // Schema AST builder
    HandlersFile *HandlersFile      // Handler AST builder
    yaml *openapi3.T               // Parsed OpenAPI spec
    PackageName string              // Generated package name
    ImportPrefix string             // Import path prefix
    ModelsImportPath string         // Models import path
    CurrentYAMLFile string          // Currently processing file
    YAMLFilesToProcess []string     // Files to process
    YAMLFilesProcessed map[string]bool // Processed files
}
```

#### Schema Generation
```go
type SchemaStruct struct {
    Name   string
    Fields []SchemaField
}

type SchemaField struct {
    Name        string
    Type        string
    TagJSON     []string
    TagValidate []string
    Required    bool
}
```

## Code Generation Process

### 1. Entry Point (`cmd/generate.go`)
- Parses command-line options
- Creates Generator instance
- Calls `Generate()` method

### 2. Main Generation Loop (`api.go:Generate()`)
```go
func (g *Generator) Generate(ctx context.Context) error {
    for len(g.YAMLFilesToProcess) > 0 {
        g.CurrentYAMLFile = g.YAMLFilesToProcess[0]
        // Process file: PrepareFiles() → GenerateFiles() → WriteOutFiles()
    }
}
```

### 3. File Processing Pipeline
1. **PrepareFiles()**: Create directories, parse YAML, initialize AST builders
2. **GenerateFiles()**: Process paths and schemas, build AST
3. **WriteOutFiles()**: Write generated Go code to files

### 4. Schema Processing (`generateschemas.go`)
- Processes OpenAPI components/schemas
- Creates Go structs with validation tags
- Handles external references (`$ref`)

### 5. Handler Processing (`generatehandlers.go`)
- Processes OpenAPI paths and operations
- Creates HTTP handlers with parameter parsing
- Generates Chi router integration

## Key Functions for AI Agents

### Schema Generation
- `ProcessSchema()` - Main schema processing entry point
- `AddSchema()` - Adds schema to AST
- `ProcessSchemas()` - Processes all schemas in OpenAPI spec

### Handler Generation
- `ProcessPaths()` - Processes all OpenAPI paths
- `ProcessOperation()` - Processes individual operations
- `AddParseParamsMethods()` - Generates parameter parsing methods
- `AddWriteResponseMethod()` - Generates response writing methods

### AST Building Utilities (`utils.go`)
- `I()` - Creates identifier
- `Str()` - Creates string literal
- `Field()` - Creates struct field
- `Func()` - Creates function declaration
- `Sel()` - Creates selector expression

### Name Utilities (`nameutils.go`)
- `FormatGoLikeIdentifier()` - Converts OpenAPI names to Go identifiers
- `GoIdentLowercase()` - Converts to lowercase Go identifier

## Generated Code Patterns

### 1. Request/Response Models
```go
type OperationNameRequest struct {
    Path   OperationNamePathParams
    Query  OperationNameQueryParams
    Headers OperationNameHeaders
    Body   OperationNameRequestBody
}
```

### 2. Handler Interface
```go
type OperationNameHandler interface {
    HandleOperationName(ctx context.Context, r models.OperationNameRequest) (*models.OperationNameResponse, error)
}
```

### 3. Handler Implementation
```go
type Handler struct {
    validator *validator.Validate
    operationName OperationNameHandler
}

func (h *Handler) AddRoutes(router chi.Router) {
    router.Method("/path", h.handleOperationName)
}
```

## Common Modification Patterns

### 1. Adding New OpenAPI Feature Support

1. **Identify processing point**: Usually in `generatehandlers.go` or `generateschemas.go`
2. **Add parsing logic**: Extract data from OpenAPI spec
3. **Add AST generation**: Use utilities in `utils.go` to build Go AST
4. **Add validation**: Include validation tags if needed
5. **Update tests**: Add test cases in `test/` directory

### 2. Modifying Generated Code Structure

1. **Update AST builders**: Modify `handlers.go` or `schemas.go`
2. **Update generation logic**: Modify `generatehandlers.go` or `generateschemas.go`
3. **Update utilities**: Modify helper functions in `utils.go`
4. **Update tests**: Ensure tests reflect new structure

### 3. Adding New Parameter Types

1. **Add parsing method**: Create new `AddParse*ParamsMethod()` function
2. **Add type conversion**: Handle OpenAPI type to Go type mapping
3. **Add validation**: Include appropriate validation tags
4. **Update operation processing**: Call new parsing method in `ProcessOperation()`

## Testing Strategy

### 1. Test Structure
- `test/generator_test.go` - Main generator tests
- `test/handlers_test.go` - Handler-specific tests
- `test/testdata/` - Test OpenAPI specs and expected output

### 2. Test Patterns
- **Golden file tests**: Compare generated code with expected output
- **Unit tests**: Test individual functions
- **Integration tests**: Test full generation pipeline

### 3. Running Tests
```bash
make test          # Run all tests
go test ./...      # Run tests for specific package
```

## Common Pitfalls and Solutions

### 1. AST Building
- **Issue**: Incorrect AST structure causes compilation errors
- **Solution**: Use utilities in `utils.go` and follow existing patterns

### 2. External References
- **Issue**: `$ref` to external files not processed
- **Solution**: Ensure files are added to `YAMLFilesToProcess` in `ParseRefTypeName()`

### 3. Validation Tags
- **Issue**: Missing or incorrect validation tags
- **Solution**: Check `SchemaField.TagValidate` population in schema processing

### 4. Import Management
- **Issue**: Missing or incorrect imports in generated code
- **Solution**: Use `AddSchemasImport()` or `AddHandlersImport()` methods

## File Modification Guidelines

### 1. When to Modify Core Files
- **`api.go`**: Main orchestration changes, new processing steps
- **`generator.go`**: Utility functions, ref parsing, naming
- **`generatehandlers.go`**: Handler generation logic
- **`generateschemas.go`**: Schema generation logic
- **`handlers.go`**: Handler AST building utilities
- **`schemas.go`**: Schema AST building utilities
- **`utils.go`**: AST helper functions

### 2. When to Add New Files
- **New major features**: Create new generation modules
- **New utilities**: Add helper functions
- **New tests**: Add test files for new functionality

### 3. Code Style Guidelines
- Follow existing error handling patterns with `errors.Wrap()`
- Use descriptive function and variable names
- Add comments for complex logic
- Follow Go naming conventions

## Debugging Tips

### 1. AST Inspection
- Use `go/ast` package to inspect generated AST
- Print AST nodes for debugging
- Compare with working examples

### 2. Generated Code Issues
- Check generated code in `internal/usage/generated/`
- Verify imports and package names
- Check validation tag correctness

### 3. OpenAPI Processing
- Validate OpenAPI spec before processing
- Check external reference resolution
- Verify parameter and schema extraction

## Dependencies

- **kin-openapi**: OpenAPI 3.0 parsing and validation
- **go-chi/chi**: HTTP router for generated handlers
- **go-playground/validator**: Validation for generated models
- **shopspring/decimal**: Decimal number support
- **golang.org/x/text**: Text processing utilities

## Command Line Options

- `-d, --dir-prefix`: Output directory prefix
- `-p, --package-prefix`: Package import prefix
- `--pointers`: Generate required fields as pointers
- `--allow-delete-with-body`: Allow DELETE with request body
- `--allow-remote-addr-param`: Allow RemoteAddr parameter

This guide should help AI agents understand the codebase structure and make informed modifications to the OpenAPI code generator.
