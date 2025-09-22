# Phase 1 Decomposition: AST Building Abstraction Layer

## Overview
Phase 1 focuses on creating a high-level abstraction layer for AST building to replace the current low-level, repetitive AST construction patterns throughout the codebase.

## Current State Analysis

### Identified AST Building Patterns

#### 1. **Parameter Parsing Functions** (Most Complex)
- `AddParseQueryParamsMethod()` - 100+ lines
- `AddParseHeadersMethod()` - 100+ lines  
- `AddParseCookiesMethod()` - 100+ lines
- `AddParsePathParamsMethod()` - Similar pattern
- `AddParseRequestBodyMethod()` - 100+ lines

**Common Pattern:**
```go
bodyList := []ast.Stmt{
    // 1. Declare struct variable
    &ast.DeclStmt{...},
}
for _, param := range params {
    // 2. Extract parameter value
    bodyList = append(bodyList, &ast.AssignStmt{...})
    // 3. Validate required parameters
    if param.Value.Required {
        bodyList = append(bodyList, &ast.IfStmt{...})
    }
    // 4. Assign to struct field
    bodyList = append(bodyList, g.AssignStringField(...)...)
}
// 5. Validate struct
bodyList = append(bodyList, &ast.AssignStmt{...})
// 6. Return result
bodyList = append(bodyList, Ret2(...))
```

#### 2. **Schema Building Functions**
- `AddSchema()` - Struct field building
- `AddTypeAlias()` - Type alias creation
- `AddSliceAlias()` - Slice type creation
- `AddParamsModel()` - Parameter model creation

#### 3. **Handler Building Functions**
- `InitHandlerStruct()` - Handler struct creation
- `InitHandlerConstructor()` - Constructor creation
- `AddHandlersInterface()` - Interface creation
- `AddWriteResponseMethodHandlers()` - Response writing

#### 4. **Validation Functions**
- `AddObjectValidate()` - Complex validation logic (120+ lines)
- `AddArrayValidate()` - Array validation
- `GetValidateFuncStmt()` - Validation function references

## Phase 1 Implementation Plan

### Task 1.1: Create Core AST Builder Package
**Duration:** 2-3 days  
**Priority:** Critical

#### 1.1.1 Create Package Structure
```
internal/generator/astbuilder/
├── builder.go          # Main Builder struct and core methods
├── expressions.go      # Expression building utilities
├── statements.go       # Statement building utilities
├── types.go           # Type building utilities
├── functions.go       # Function building utilities
└── patterns.go        # Common AST patterns
```

#### 1.1.2 Core Builder Interface
```go
// internal/generator/astbuilder/builder.go
type Builder struct {
    imports map[string]bool
    stmts   []ast.Stmt
    decls   []ast.Decl
}

type BuilderConfig struct {
    PackageName    string
    ImportPrefix   string
    UsePointers    bool
}

func NewBuilder(config BuilderConfig) *Builder
func (b *Builder) AddImport(path string) *Builder
func (b *Builder) AddStatement(stmt ast.Stmt) *Builder
func (b *Builder) AddDeclaration(decl ast.Decl) *Builder
func (b *Builder) Build() ([]ast.Stmt, []ast.Decl, []string)
```

### Task 1.2: Parameter Parsing Abstraction
**Duration:** 3-4 days  
**Priority:** Critical

#### 1.2.1 Create Parameter Parser
```go
// internal/generator/astbuilder/parameter_parser.go
type ParameterParser struct {
    builder *Builder
    config  ParameterConfig
}

type ParameterConfig struct {
    BaseName        string
    PackageName     string
    UsePointers     bool
    ParameterType   string // "Query", "Header", "Cookie", "Path"
}

func (p *ParameterParser) ParseParameters(params openapi3.Parameters) error
func (p *ParameterParser) ParseQueryParams(params openapi3.Parameters) error
func (p *ParameterParser) ParseHeaders(params openapi3.Parameters) error
func (p *ParameterParser) ParseCookies(params openapi3.Parameters) error
func (p *ParameterParser) ParsePathParams(params openapi3.Parameters) error
```

#### 1.2.2 High-Level Parameter Methods
```go
func (p *ParameterParser) DeclareStruct() *ParameterParser
func (p *ParameterParser) ExtractParameter(param *openapi3.Parameter) *ParameterParser
func (p *ParameterParser) ValidateRequired(param *openapi3.Parameter) *ParameterParser
func (p *ParameterParser) AssignToField(param *openapi3.Parameter) *ParameterParser
func (p *ParameterParser) ValidateStruct() *ParameterParser
func (p *ParameterParser) ReturnResult() *ParameterParser
```

### Task 1.3: Schema Building Abstraction
**Duration:** 2-3 days  
**Priority:** High

#### 1.3.1 Create Schema Builder
```go
// internal/generator/astbuilder/schema_builder.go
type SchemaBuilder struct {
    builder *Builder
    config  SchemaConfig
}

type SchemaConfig struct {
    PackageName string
    UsePointers bool
}

func (s *SchemaBuilder) BuildStruct(model SchemaStruct) error
func (s *SchemaBuilder) BuildTypeAlias(name, typeName string) error
func (s *SchemaBuilder) BuildSliceAlias(name, elementType string) error
func (s *SchemaBuilder) BuildField(field SchemaField) *ast.Field
```

#### 1.3.2 High-Level Schema Methods
```go
func (s *SchemaBuilder) AddStruct(name string, fields []SchemaField) *SchemaBuilder
func (s *SchemaBuilder) AddField(name, typeName string, tags map[string]string) *SchemaBuilder
func (s *SchemaBuilder) AddTypeAlias(name, typeName string) *SchemaBuilder
func (s *SchemaBuilder) AddSliceType(name, elementType string) *SchemaBuilder
```

### Task 1.4: Handler Building Abstraction
**Duration:** 2-3 days  
**Priority:** High

#### 1.4.1 Create Handler Builder
```go
// internal/generator/astbuilder/handler_builder.go
type HandlerBuilder struct {
    builder *Builder
    config  HandlerConfig
}

type HandlerConfig struct {
    PackageName string
    UsePointers bool
}

func (h *HandlerBuilder) BuildHandlerStruct() error
func (h *HandlerBuilder) BuildConstructor() error
func (h *HandlerBuilder) BuildInterface(name string, methods []MethodSpec) error
func (h *HandlerBuilder) BuildRoutesFunction() error
```

#### 1.4.2 High-Level Handler Methods
```go
func (h *HandlerBuilder) AddHandlerField(name, typeName string) *HandlerBuilder
func (h *HandlerBuilder) AddInterfaceMethod(name string, params, returns []FieldSpec) *HandlerBuilder
func (h *HandlerBuilder) AddRoute(method, path string, handlerName string) *HandlerBuilder
func (h *HandlerBuilder) AddResponseWriter(baseName string, codes []string) *HandlerBuilder
```

### Task 1.5: Validation Building Abstraction
**Duration:** 2-3 days  
**Priority:** Medium-High

#### 1.5.1 Create Validation Builder
```go
// internal/generator/astbuilder/validation_builder.go
type ValidationBuilder struct {
    builder *Builder
    config  ValidationConfig
}

func (v *ValidationBuilder) BuildObjectValidation(modelName string, schema *openapi3.SchemaRef) error
func (v *ValidationBuilder) BuildArrayValidation(modelName string, schema *openapi3.SchemaRef) error
func (v *ValidationBuilder) BuildFieldValidation(fieldName string, schema *openapi3.SchemaRef) []string
```

#### 1.5.2 High-Level Validation Methods
```go
func (v *ValidationBuilder) AddRequiredFieldsValidation(fields []string) *ValidationBuilder
func (v *ValidationBuilder) AddFieldValidation(fieldName string, validators []string) *ValidationBuilder
func (v *ValidationBuilder) AddJSONUnmarshal() *ValidationBuilder
func (v *ValidationBuilder) AddErrorHandling() *ValidationBuilder
```

### Task 1.6: Expression and Statement Utilities
**Duration:** 1-2 days  
**Priority:** Medium

#### 1.6.1 Expression Builder
```go
// internal/generator/astbuilder/expressions.go
func (b *Builder) Call(receiver, method string, args ...ast.Expr) ast.Expr
func (b *Builder) Select(receiver, field string) ast.Expr
func (b *Builder) Ident(name string) ast.Expr
func (b *Builder) String(value string) ast.Expr
func (b *Builder) Int(value int) ast.Expr
func (b *Builder) Bool(value bool) ast.Expr
func (b *Builder) Nil() ast.Expr
func (b *Builder) AddressOf(expr ast.Expr) ast.Expr
func (b *Builder) Deref(expr ast.Expr) ast.Expr
```

#### 1.6.2 Statement Builder
```go
// internal/generator/astbuilder/statements.go
func (b *Builder) DeclareVar(name, typeName string, value ast.Expr) ast.Stmt
func (b *Builder) Assign(lhs, rhs ast.Expr) ast.Stmt
func (b *Builder) If(cond ast.Expr, body []ast.Stmt) ast.Stmt
func (b *Builder) IfElse(cond ast.Expr, ifBody, elseBody []ast.Stmt) ast.Stmt
func (b *Builder) Return(values ...ast.Expr) ast.Stmt
func (b *Builder) CallStmt(receiver, method string, args ...ast.Expr) ast.Stmt
```

### Task 1.7: Migration Strategy
**Duration:** 3-4 days  
**Priority:** Critical

#### 1.7.1 Create Migration Wrapper
```go
// internal/generator/astbuilder/migration.go
type MigrationWrapper struct {
    generator *Generator
    builder   *Builder
}

func (m *MigrationWrapper) MigrateParameterParsing() error
func (m *MigrationWrapper) MigrateSchemaBuilding() error
func (m *MigrationWrapper) MigrateHandlerBuilding() error
func (m *MigrationWrapper) MigrateValidationBuilding() error
```

#### 1.7.2 Gradual Migration Plan
1. **Week 1**: Implement core builder and parameter parsing
2. **Week 2**: Migrate parameter parsing functions
3. **Week 3**: Implement schema and handler building
4. **Week 4**: Migrate remaining functions and cleanup

### Task 1.8: Testing Strategy
**Duration:** 2-3 days  
**Priority:** High

#### 1.8.1 Unit Tests
```go
// internal/generator/astbuilder/builder_test.go
func TestBuilder_BasicOperations(t *testing.T)
func TestParameterParser_QueryParams(t *testing.T)
func TestSchemaBuilder_StructCreation(t *testing.T)
func TestHandlerBuilder_InterfaceCreation(t *testing.T)
```

#### 1.8.2 Integration Tests
```go
// internal/generator/astbuilder/integration_test.go
func TestParameterParsingMigration(t *testing.T)
func TestSchemaBuildingMigration(t *testing.T)
func TestHandlerBuildingMigration(t *testing.T)
```

#### 1.8.3 Golden File Tests
```go
// internal/generator/astbuilder/golden_test.go
func TestGeneratedCodeMatchesExpected(t *testing.T)
```

## Implementation Timeline

### Week 1: Core Infrastructure
- **Day 1-2**: Task 1.1 - Core AST Builder Package
- **Day 3-4**: Task 1.2 - Parameter Parsing Abstraction
- **Day 5**: Task 1.6 - Expression and Statement Utilities

### Week 2: Migration and Testing
- **Day 1-2**: Task 1.7 - Migration Strategy Implementation
- **Day 3-4**: Task 1.8 - Testing Strategy
- **Day 5**: Integration testing and bug fixes

## Success Metrics

### Code Complexity Reduction
- **Before**: 100+ line functions with complex nested AST
- **After**: 10-20 line functions using high-level abstractions
- **Target**: 70% reduction in function complexity

### AI Agent Benefits
- **Understandability**: AI can work with high-level operations instead of low-level AST
- **Consistency**: All AST building follows the same patterns
- **Maintainability**: Changes to AST patterns only require updating the builder
- **Extensibility**: New AST patterns can be easily added to the builder

### Code Quality Improvements
- **Reduced Duplication**: Common patterns are abstracted into reusable methods
- **Better Error Handling**: Centralized error handling in the builder
- **Improved Readability**: Code intent is clear from method names
- **Easier Testing**: Each builder component can be tested independently

## Risk Mitigation

### Backward Compatibility
- Keep existing functions during migration
- Use feature flags to switch between old and new implementations
- Gradual migration with rollback capability

### Performance Impact
- Builder methods should be lightweight
- Avoid unnecessary AST node creation
- Profile performance before and after migration

### Testing Coverage
- Maintain 100% test coverage during migration
- Use golden file tests to ensure output consistency
- Integration tests to verify end-to-end functionality

## Next Steps

1. **Start with Task 1.1**: Create the core builder package
2. **Implement Task 1.2**: Focus on parameter parsing as it's the most complex
3. **Create comprehensive tests**: Ensure the abstraction works correctly
4. **Begin migration**: Start with the most problematic functions
5. **Iterate and refine**: Based on testing and usage feedback

This decomposition provides a clear roadmap for implementing Phase 1 while maintaining code quality and ensuring AI agents can effectively work with the new abstraction layer.
