# Migration Guide: AST Builder Abstraction Layer

This guide explains how to migrate from the existing generator code to the new AST builder abstraction layer.

## Overview

The migration strategy provides a gradual transition from the existing generator methods to the new AST builder abstractions. This allows for:

- **Backward Compatibility**: Existing code continues to work during migration
- **Gradual Migration**: Components can be migrated one at a time
- **Comparison Testing**: Old and new methods can be compared side-by-side
- **Rollback Safety**: Easy rollback if issues are discovered

## Migration Modes

### 1. Legacy Mode (`MigrationModeLegacy`)
- Uses only the existing generator methods
- No changes to current behavior
- Safe for production use during migration

### 2. Hybrid Mode (`MigrationModeHybrid`)
- Uses both old and new methods
- Allows comparison of outputs
- Useful for testing and validation
- Generates output from both methods

### 3. New Mode (`MigrationModeNew`)
- Uses only the new AST builder methods
- Full migration to new abstractions
- Recommended for new projects

## Migration Components

### 1. Parameter Parsing Migration
**Old Method**: Direct AST manipulation in `handlers.go` and `handlers2.go`
**New Method**: `ParameterParser` abstraction

```go
// Old way
func (g *Generator) AddParseQueryParamsMethod(baseName string, params openapi3.Parameters) error {
    // Direct AST manipulation
}

// New way
paramParser := NewParameterParser(builder, config)
_, err := paramParser.ParseParameters(params)
```

### 2. Schema Building Migration
**Old Method**: Direct AST manipulation in `schemas.go`
**New Method**: `SchemaBuilder` abstraction

```go
// Old way
func (g *Generator) GenerateImportsSpecsSchemas(imp []string) ([]*ast.ImportSpec, []ast.Spec) {
    // Direct AST manipulation
}

// New way
schemaBuilder := NewSchemaBuilder(builder, config)
_, err := schemaBuilder.BuildFromOpenAPISchema(name, schema)
```

### 3. Handler Building Migration
**Old Method**: Direct AST manipulation in `handlers.go`
**New Method**: `HandlerBuilder` abstraction

```go
// Old way
func (g *Generator) GenerateHandlersFile() *ast.File {
    // Direct AST manipulation
}

// New way
handlerBuilder := NewHandlerBuilder(builder, config)
_, err := handlerBuilder.BuildFromOpenAPI(spec)
```

### 4. Validation Building Migration
**Old Method**: Manual validation code generation
**New Method**: `ValidationBuilder` abstraction

```go
// Old way
// Manual validation code generation scattered throughout

// New way
validationBuilder := NewValidationBuilder(builder, config)
_, err := validationBuilder.BuildObjectValidation(name, schema)
```

## Step-by-Step Migration Process

### Phase 1: Core Infrastructure (Week 1)
1. **Implement Core Builder**
   - ✅ Core AST builder with configuration
   - ✅ Import management
   - ✅ Statement and declaration management

2. **Implement Parameter Parsing**
   - ✅ ParameterParser abstraction
   - ✅ Fluent interface for parameter handling
   - ✅ Support for query, header, cookie, and path parameters

3. **Create Migration Wrapper**
   - ✅ MigrationWrapper struct
   - ✅ Migration modes (Legacy, Hybrid, New)
   - ✅ Comparison functionality

4. **Set Up Hybrid Mode**
   - ✅ Side-by-side comparison
   - ✅ Output validation
   - ✅ Error handling

### Phase 2: Schema Migration (Week 2)
1. **Implement Schema Building**
   - ✅ SchemaBuilder abstraction
   - ✅ OpenAPI schema to Go struct conversion
   - ✅ Type mapping and validation tags

2. **Migrate Schema Functions**
   - Migrate `GenerateImportsSpecsSchemas`
   - Migrate `GenerateRequestModel`
   - Update schema generation logic

3. **Test Schema Compatibility**
   - Compare old vs new schema output
   - Validate type mappings
   - Test edge cases

### Phase 3: Handler Migration (Week 3)
1. **Implement Handler Building**
   - ✅ HandlerBuilder abstraction
   - ✅ HTTP handler generation
   - ✅ Route registration

2. **Migrate Handler Functions**
   - Migrate `GenerateHandlersFile`
   - Migrate parameter parsing methods
   - Update handler generation logic

3. **Test Handler Compatibility**
   - Compare old vs new handler output
   - Validate HTTP handler structure
   - Test routing functionality

### Phase 4: Validation Migration (Week 4)
1. **Implement Validation Building**
   - ✅ ValidationBuilder abstraction
   - ✅ OpenAPI validation to Go validation conversion
   - ✅ Error handling patterns

2. **Migrate Validation Functions**
   - Migrate validation code generation
   - Update error handling
   - Consolidate validation logic

3. **Complete Migration**
   - Switch to new mode by default
   - Remove old code
   - Update documentation

## Usage Examples

### Basic Migration Setup

```go
// Create migration wrapper
config := MigrationConfig{
    PackageName:    "myapi",
    ImportPrefix:   "github.com/myorg/myapi",
    UseNewBuilders: true,
    MigrationMode:  MigrationModeHybrid,
    ValidateOutput: true,
}

wrapper := NewMigrationWrapper(generator, config)

// Generate with migration
err := wrapper.GenerateWithMigration(ctx)
if err != nil {
    log.Fatal(err)
}
```

### Gradual Migration

```go
// Start with legacy mode
wrapper.SetMigrationMode(MigrationModeLegacy)
err := wrapper.GenerateWithMigration(ctx)

// Switch to hybrid mode for comparison
wrapper.SetMigrationMode(MigrationModeHybrid)
err = wrapper.GenerateWithMigration(ctx)

// Compare outputs
comparison, err := wrapper.CompareOutputs(ctx)
if err != nil {
    log.Fatal(err)
}

// Switch to new mode when ready
wrapper.SetMigrationMode(MigrationModeNew)
err = wrapper.GenerateWithMigration(ctx)
```

### Component-Specific Migration

```go
// Migrate only parameter parsing
err := wrapper.MigrateParameterParsing()

// Migrate only schema building
err := wrapper.MigrateSchemaBuilding()

// Migrate only handler building
err := wrapper.MigrateHandlerBuilding()

// Migrate only validation building
err := wrapper.MigrateValidationBuilding()
```

## Testing Strategy

### Unit Tests
- Test each migration component individually
- Verify AST structure correctness
- Test error handling and edge cases

### Integration Tests
- Test complete migration workflows
- Compare old vs new outputs
- Test different migration modes

### Golden File Tests
- Compare generated code with expected output
- Ensure backward compatibility
- Validate code formatting and structure

## Rollback Strategy

If issues are discovered during migration:

1. **Immediate Rollback**
   ```go
   wrapper.SetMigrationMode(MigrationModeLegacy)
   ```

2. **Component Rollback**
   - Disable specific migration components
   - Use hybrid mode to identify problematic components

3. **Full Rollback**
   - Revert to previous version
   - Use legacy mode exclusively

## Best Practices

### 1. Start with Hybrid Mode
- Always begin migration in hybrid mode
- Compare outputs before switching to new mode
- Validate all generated code

### 2. Migrate Incrementally
- Migrate one component at a time
- Test thoroughly after each migration
- Keep old code until migration is complete

### 3. Validate Outputs
- Enable output validation during migration
- Compare generated code with expected output
- Test with real OpenAPI specifications

### 4. Monitor Performance
- Compare generation speed between old and new methods
- Optimize if necessary
- Document performance characteristics

## Troubleshooting

### Common Issues

1. **AST Structure Mismatch**
   - Check that builder methods are called correctly
   - Verify AST node types and relationships
   - Use AST inspection tools

2. **Import Management Issues**
   - Ensure imports are added correctly
   - Check for duplicate imports
   - Verify import paths

3. **Type Mapping Problems**
   - Verify OpenAPI to Go type mapping
   - Check for missing type conversions
   - Test with different schema types

4. **Validation Errors**
   - Check validation tag generation
   - Verify validation rule correctness
   - Test with invalid inputs

### Debugging Tips

1. **Enable Debug Logging**
   ```go
   config := MigrationConfig{
       ValidateOutput: true,
       // Add debug flags
   }
   ```

2. **Inspect Generated AST**
   ```go
   file := wrapper.GetBuilder().BuildFile()
   ast.Print(token.NewFileSet(), file)
   ```

3. **Compare Outputs**
   ```go
   comparison, err := wrapper.CompareOutputs(ctx)
   if err != nil {
       log.Printf("Comparison error: %v", err)
   }
   ```

## Migration Checklist

### Pre-Migration
- [ ] Backup existing code
- [ ] Set up test environment
- [ ] Prepare test OpenAPI specifications
- [ ] Review migration plan

### During Migration
- [ ] Start with hybrid mode
- [ ] Migrate one component at a time
- [ ] Test after each migration
- [ ] Compare outputs
- [ ] Document issues and solutions

### Post-Migration
- [ ] Switch to new mode
- [ ] Remove old code
- [ ] Update documentation
- [ ] Performance testing
- [ ] User acceptance testing

## Support

For questions or issues during migration:

1. Check this guide for common solutions
2. Review test cases for examples
3. Examine the migration wrapper code
4. Create issues with detailed information

## Conclusion

The migration strategy provides a safe, gradual path from the existing generator code to the new AST builder abstractions. By following this guide and using the provided tools, you can migrate your code generation system with confidence and minimal risk.
