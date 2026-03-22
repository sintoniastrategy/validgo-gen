# Architecture

## Project layout

```
cmd/generate.go                     CLI entry point
internal/generator/
  options/options.go                CLI flag parsing
  api.go                            Generator struct, main loop, file I/O
  generator.go                      Ref parsing, import resolution, YAML tracking
  schemas.go                        Schema/model processing → Go AST
  handlers.go                       Handler AST construction (interfaces, structs, routing, responses)
  handlers2.go                      Query/header/cookie/body parsing, JSON validation
  generatehandlers.go               Orchestration: paths → operations → handler pipeline
  generateschemas.go                Orchestration: components/schemas iteration
  bb.go                             AST builder helpers (I, Str, Star, Sel, Amp, Func, Field, ...)
  nameutils.go                      Go identifier formatting, common initialisms
  utils.go                          OpenAPI constraint → validator tag conversion
```

## Core types

```go
// Generator — the main orchestrator
type Generator struct {
    Options             *options.Options
    YAMLFilesToProcess  []string           // discovery queue
    YAMLFilesProcessed  map[string]bool    // already processed
    SchemasFile         *SchemasFile       // accumulates models AST
    HandlersFile        *HandlersFile      // accumulates handlers AST
    SchemaStructs       []*SchemaStruct    // processed struct metadata
}

// SchemasFile — models.go builder
type SchemasFile struct {
    File    *ast.File             // the Go AST file node
    Imports map[string]string     // path → alias
    Schemas []*SchemaStruct       // struct definitions
}

// HandlersFile — handlers.go builder
type HandlersFile struct {
    File            *ast.File
    Imports         map[string]string
    HandlerStruct   *ast.TypeSpec     // Handler struct accumulator
    ConstructorFunc *ast.FuncDecl     // NewHandler() accumulator
    RoutesFunc      *ast.FuncDecl     // AddRoutes() accumulator
    // ... switch statements for content-type/response dispatching
}

// SchemaStruct — metadata for a generated struct
type SchemaStruct struct {
    Name   string
    Fields []*SchemaField
}

// SchemaField — metadata for a struct field
type SchemaField struct {
    Name       string
    Type       string
    Required   bool
    Nullable   bool
    IsObject   bool
    IsArray    bool
    ArrayItem  *SchemaField
    // ... validators, format info
}
```

## Data flow

```
┌─────────────────┐
│   YAML files    │
└────────┬────────┘
         │
    kin-openapi parse + validate
         │
         ▼
┌─────────────────┐
│  OpenAPI 3.0    │
│  Document       │
│  (in-memory)    │
└────────┬────────┘
         │
    ┌────┴────┐
    │         │
    ▼         ▼
ProcessPaths  ProcessSchemas
    │              │
    │         ┌────┴────┐
    │         │         │
    │    ProcessSchema  AddTypeAlias/
    │    (objects)      AddSliceAlias
    │         │
    │    AddSchema()
    │    (→ AST struct decl)
    │         │
    ▼         ▼
┌──────────────────────────┐
│  SchemasFile (AST)       │  → models.go
│  HandlersFile (AST)      │  → handlers.go
└────────────┬─────────────┘
             │
        go/format.Node()
             │
             ▼
        .go source files
```

## Generation Pipeline — Step by Step

The `Generate(ctx)` method on `Generator` drives the full pipeline:

### Phase 1: File discovery loop

```
for each YAML file in queue (including discovered $ref targets):
    if already processed → skip
    mark as processed
    Phase 2–4 for this file
```

### Phase 2: `PrepareFiles()`

1. Initialize fresh `SchemasFile` and `HandlersFile` with package declarations, base imports
2. Initialize handler struct, constructor, route function, imports AST nodes

### Phase 3: `GenerateFiles()`

1. **`ProcessPaths(doc.Paths)`** — for each path + method:
   - Extract operationId → Go identifier
   - `AddInterface()` — creates handler interface (`type XxxHandler interface { HandleXxx(...) }`)
   - `AddDependencyToHandler()` — adds field to Handler struct + param to constructor
   - `AddRoute()` — adds `router.Method(path, h.handleXxx)` to AddRoutes
   - For each parameter type (path/query/header/cookie) → `AddParseParamsMethods()`
   - For request body → `ProcessApplicationJSONOperation()`:
     - Resolves body schema (inline or `$ref`)
     - Generates request model, parse method, JSON validation
   - `AddResponseCodeModels()` + `AddWriteResponseMethod()` — response structs + writers
   - `FinalizeHandlerSwitches()` — closes content-type and response-code switch statements

2. **`ProcessSchemas(doc.Components.Schemas)`** — for each component schema:
   - `ProcessSchema()` routes to `ProcessObjectSchema()`, `ProcessTypeAlias()`, or `ProcessArraySchema()`
   - Each produces AST struct declarations + JSON validation functions

### Phase 4: `WriteOutFiles()`

1. `WriteSchemasToOutput()` — renders `SchemasFile.File` via `go/format.Node()` → `models.go`
2. `WriteHandlersToOutput()` — renders `HandlersFile.File` via `go/format.Node()` → `handlers.go`
