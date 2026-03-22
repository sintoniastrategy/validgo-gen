# Reference

## OpenAPI Constraint → Validator Tag Mapping

`utils.go` → `GetSchemaValidators()` translates OpenAPI schema properties to `go-playground/validator/v10` tags:

| OpenAPI Property | Validator Tag | Notes |
|---|---|---|
| `required: true` | `required` | On the field itself |
| `required: false` / missing | `omitempty` | Prefixed |
| `minLength: N` | `min=N` | Strings |
| `maxLength: N` | `max=N` | Strings |
| `minimum: N` | `min=N` | Numbers/integers |
| `maximum: N` | `max=N` | Numbers/integers |
| `enum: [a, b, c]` | `oneof=a b c` | Space-separated |
| `format: email` | `email` | |
| `format: ip` | `ip` | |
| `format: ipv4` | `ipv4` | |
| `format: ipv6` | `ipv6` | |
| `minItems: N` | `min=N` | Arrays |
| `maxItems: N` | `max=N` | Arrays |
| `uniqueItems: true` | `unique` | Arrays |
| Array items have validators | `dive,...` | Recursive into array items |

### Unsupported (logged as warnings)

| OpenAPI Property | Status |
|---|---|
| `pattern` | Warning, skipped |
| `exclusiveMinimum` | Warning, skipped |
| `exclusiveMaximum` | Warning, skipped |
| `multipleOf` | Warning, skipped |

## Configuration Flags

| Flag | Default | Description |
|---|---|---|
| `-d <dir>` | `internal` | Directory prefix for generated file output |
| `-p <prefix>` | `internal` | Go package import path prefix |
| `-pointers` | `false` | Generate required fields as pointers too (default: only optional fields are pointers) |
| `-allow-delete-with-body` | `false` | Allow DELETE operations to have a request body (normally errors) |
| `-allow-remote-addr-param` | `false` | Allow a fake `Remote-Addr` header parameter that maps to `r.RemoteAddr` |

### Positional arguments

All remaining arguments after flags are treated as YAML file paths to process.

### Examples

```bash
# Basic usage
go run ./cmd/generate.go api.yaml

# Custom output directory and package prefix
go run ./cmd/generate.go -d ./internal/api -p github.com/myorg/myservice/internal/api api.yaml

# All required fields as pointers (useful for PATCH semantics)
go run ./cmd/generate.go -pointers api.yaml

# Allow DELETE with body + remote addr parameter
go run ./cmd/generate.go -allow-delete-with-body -allow-remote-addr-param api.yaml

# Multiple YAML files (cross-referenced)
go run ./cmd/generate.go -d ./generated -p github.com/myorg/project/generated api.yaml definitions.yaml
```

## Testing Strategy

The project has three testing layers:

### Unit tests (`internal/generator/generator_test.go`)

~1200 lines testing individual generation features. Each test:
1. Constructs a minimal `openapi3.T` document programmatically
2. Runs the generator
3. Compares output strings against expected Go code

**Test coverage by feature:**

| Test | What it covers |
|---|---|
| `TestGenerateValidInput` | End-to-end generation from a YAML file |
| `TestGeneratorFeatures/string` | String field → `string` type, json tag |
| `TestGeneratorFeatures/int` | Integer field → `int` type |
| `TestGeneratorFeatures/number` | Number field → `float64` type |
| `TestGeneratorFeatures/bool` | Boolean field → `bool` type |
| `TestGeneratorFeatures/simpleObject` | Nested object → separate struct |
| `TestGeneratorFeatures/emptyObject` | Empty object (no properties) |
| `TestGeneratorFeatures/validators` | All OpenAPI constraints → validate tags |
| `TestGeneratorFeatures/ref` | `$ref` to another component schema |
| `TestGeneratorFeatures/objectWithNestedObject` | Deeply nested inline objects |
| `TestGeneratorFeatures/stringArray` | `[]string` |
| `TestGeneratorFeatures/intArray` | `[]int` |
| `TestGeneratorFeatures/numberArray` | `[]float64` |
| `TestGeneratorFeatures/boolArray` | `[]bool` |
| `TestGeneratorFeatures/refArray` | `[]ReferencedType` |
| `TestGeneratorFeatures/objectArray` | `[]InlineStruct` |
| `TestGeneratorFeatures/nestedArray` | `[][]string` |
| `TestGeneratorFeatures/nestedArrayOfObjects` | `[][]InlineStruct` |
| `TestGeneratorFeatures/nestedNestedArray` | `[][][]string` |
| `TestGeneratePaths` | Full path → handler generation |
| `TestGenerateFeatures` | Body `$ref` → handler |
| `TestGenerateFeatures2` | OperationID formatting |
| `TestGenerateCookies` | Required + optional cookie params |
| `TestGenerateExternal` | External `$ref` across files |

### Validator tests (`internal/generator/validator_test.go`)

~150 lines testing that generated validation tags work correctly at runtime:
- `TestStrings` — required vs optional string validation
- `TestArrayValidatorsDive` — array min/max/unique + dive into item validators

### Integration tests (`test/`)

**Golden file tests** (`test/generator_test.go`):
- Runs generator against 4 YAML specs (`api.yaml`, `api2.yaml`, `api3.yaml`, `def.yaml`)
- Compares all 8 generated files against goldie snapshots in `test/testdata/generated/`
- Uses `goldie.WithFixtureDir()` for snapshot management
- Run `go test ./test/ -update` to update snapshots

**HTTP handler tests** (`test/handlers_test.go`):
- Spins up `httptest.Server` with generated handlers + chi router
- Tests actual HTTP requests and responses:
  - 200 success with valid request
  - 404 for missing resources
  - 400 for validation failures (missing name, invalid enum, bad cookie, invalid path param suffix)
  - 400 for dive validation failures on nested array items
  - 500 when handler returns nil response

## Supported & Unsupported OpenAPI Features

### Fully supported

| Feature | Notes |
|---|---|
| `openapi: 3.0.x` | Parsed via kin-openapi |
| `paths` with `get`, `post`, `put`, `patch`, `delete` | DELETE needs `-allow-delete-with-body` for bodies |
| `operationId` | Used as Go identifier base |
| Path parameters (`in: path`) | String type only |
| Query parameters (`in: query`) | String type only |
| Header parameters (`in: header`) | String type, with `date-time` parsing to `time.Time` |
| Cookie parameters (`in: cookie`) | Required vs optional |
| `application/json` request/response bodies | |
| `$ref` to `#/components/schemas/*` | Local and external file refs |
| `type: string/integer/number/boolean/object/array` | |
| `format: date-time` | → `time.Time` |
| `format: decimal` | → `shopspring/decimal.Decimal` |
| `format: int8/16/32/64, uint8/16/32/64` | Precise integer types |
| `format: float/double` | → `float32`/`float64` |
| `format: email/ip/ipv4/ipv6` | Validator tags |
| `required` fields | Value types (or pointers with `-pointers`) |
| `nullable` fields | Pointer types, null check in JSON validation |
| `minLength/maxLength` | Validator tags |
| `minimum/maximum` | Validator tags |
| `enum` | → `oneof=` validator tag |
| `minItems/maxItems` | Array validator tags |
| `uniqueItems` | → `unique` validator tag |
| Inline (anonymous) object schemas | Named by parent context |
| Nested arrays (`array of array of ...`) | Recursive processing |
| Multiple response codes (200, 400, 404, etc.) | Per-code response structs + writers |
| Response headers | Generated writer methods set headers |

### Not supported (TODO or limitation)

| Feature | Status |
|---|---|
| Non-string path/query/header/cookie params | TODO |
| `additionalProperties` | TODO |
| Component-level `parameters` | TODO |
| Component-level `requestBodies` | TODO |
| Component-level `responses` | TODO |
| Component-level `headers` | TODO |
| External `$ref` at component level (non-schema) | TODO |
| `pattern` (regex) | Logged warning, skipped |
| `exclusiveMinimum/exclusiveMaximum` | Logged warning, skipped |
| `multipleOf` | Logged warning, skipped |
| Non-JSON content types (multipart, form, XML, etc.) | Errors during generation |
| Multiple content types per response code | Errors during generation |
| `oneOf/anyOf/allOf` composition | Not handled |
| `discriminator` | Not handled |
| Security schemes | Not handled |
| Server definitions | Not handled |
| OpenAPI 3.1 | Uses kin-openapi 3.0 types |
| Callbacks, webhooks, links | Not handled |

## Dependencies

### Build-time (generator itself)

| Dependency | Purpose |
|---|---|
| `github.com/getkin/kin-openapi` | OpenAPI 3.0 YAML parsing and validation |
| `github.com/go-faster/errors` | Error wrapping in generator code |
| `golang.org/x/text` | Unicode-aware text processing for name formatting |

### Runtime (generated code depends on)

| Dependency | Purpose |
|---|---|
| `github.com/go-chi/chi/v5` | HTTP router — `AddRoutes(chi.Router)` |
| `github.com/go-playground/validator/v10` | Struct validation via tags |
| `github.com/go-faster/errors` | Error wrapping in generated handlers |
| `github.com/shopspring/decimal` | Decimal type (when `format: decimal` is used) |

### Test-only

| Dependency | Purpose |
|---|---|
| `github.com/sebdah/goldie/v2` | Golden file snapshot testing |
| `github.com/stretchr/testify` | Test assertions |

## Extending the Generator

### Adding a new type format

1. In `schemas.go` → `GetStringType()` or `GetIntegerType()`: add the format mapping
2. In `schemas.go` → `GetFieldTypeFromSchema()`: handle any special import needs
3. In `utils.go` → `GetSchemaValidators()`: add validator tag if applicable
4. In `handlers2.go`: if the type needs special parsing (like `time.Time` from strings), add parsing logic

### Adding a new parameter location

1. In `generatehandlers.go` → `GetOperationParamsByType()`: filter params by the new location
2. In `schemas.go` → `AddParamsModel()`: generate the model struct
3. In `handlers2.go`: create an `AddParse<Location>Method()` following the pattern of existing ones
4. In `handlers2.go` → `AddParseRequestMethod()`: call the new parse method

### Adding a new content type

1. In `generatehandlers.go` → `ProcessOperation()`: handle the new content type alongside `application/json`
2. Create a `Process<ContentType>Operation()` method
3. Add request body parsing in `handlers2.go`
4. Update `AddContentTypeToHandler()` for content-type switching

### Adding non-string parameter types (the main TODO)

This is the most impactful extension. The approach:

1. In `schemas.go` → `AddParamsModel()`: instead of always using `string`, call `GetFieldTypeFromSchema()` for param schemas
2. In `handlers.go` → `AddParsePathParamsMethod()`: instead of `chi.URLParam(r, name)` (returns string), add `strconv.Atoi()` / `strconv.ParseFloat()` etc. based on the param type
3. In `handlers2.go` → `AddParseQueryParamsMethod()` and `AddParseHeadersMethod()`: similarly add type conversion from the string value
4. Add proper error handling for parse failures (400 Bad Request)

### Key patterns to follow when extending

- Use `bb.go` helpers for all AST construction
- Add imports via `AddSchemasImport()` / `AddHandlersImport()` — they deduplicate
- Follow the `const op = "..."` error wrapping pattern
- Add unit tests in `generator_test.go` constructing minimal OpenAPI documents
- Update golden files with `go test ./test/ -update`
- Run `make check` to verify lint + tests
