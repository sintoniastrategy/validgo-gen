# Internals

## External `$ref` Resolution

The generator supports cross-file `$ref` references:

```yaml
# api.yaml
components:
  schemas:
    CreateRequestBody:
      properties:
        external:
          $ref: 'def.yml#/components/schemas/ExternalRef'
```

### How it works

1. `refIsExternal(ref)` checks if the ref starts with a filename (not `#/`)
2. `parseFilenameFromRef(ref)` extracts the filename
3. The filename is added to `YAMLFilesToProcess` queue
4. On the next iteration of the main loop, `def.yml` gets its own `SchemasFile` and `HandlersFile`
5. Import paths are computed: `GetModelsImportForFile("def.yml")` → `"<prefix>/def/defmodels"`
6. The referencing file gets an import statement and uses the qualified type name

### Generated import in the referencing file

```go
// api/apimodels/models.go
import "github.com/myorg/project/generated/def/defmodels"

type CreateRequestBody struct {
    External defmodels.ExternalRef `json:"external"`
}
```

### Limitations

- Only schema-level `$ref` is supported (not component-level parameters, responses, headers)
- External refs must be relative file paths (not URLs)
- Circular refs are prevented by the `YAMLFilesProcessed` map

## AST Builder (`bb.go`)

The `bb.go` file provides terse helper functions for constructing `go/ast` nodes:

| Helper | Produces | Example |
|---|---|---|
| `I(name)` | `*ast.Ident` | `I("ctx")` → identifier `ctx` |
| `Str(val)` | `*ast.BasicLit` (string) | `Str("application/json")` → `"application/json"` |
| `Star(expr)` | `*ast.StarExpr` | `Star(I("Handler"))` → `*Handler` |
| `Sel(x, sel)` | `*ast.SelectorExpr` | `Sel(I("h"), "validator")` → `h.validator` |
| `Amp(expr)` | `&expr` (unary) | `Amp(I("resp"))` → `&resp` |
| `Ne(x, y)` | `x != y` | `Ne(I("err"), I("nil"))` → `err != nil` |
| `Eq(x, y)` | `x == y` | |
| `Ret(results...)` | `return` stmt | `Ret(I("nil"), I("err"))` → `return nil, err` |
| `Ret1(expr)` | `return expr` | |
| `Ret2(a, b)` | `return a, b` | |
| `Func(recv, name, params, results, body)` | `*ast.FuncDecl` | Full function declaration |
| `Field(names, type)` | `*ast.Field` | Struct or param field |
| `FieldA(names, type, tag)` | `*ast.Field` with tag | Struct field with `json`/`validate` tags |

### Example: how a parse method is built

```go
// This Go code in the generator:
Func(
    Field([]string{"h"}, Star(I("Handler"))),     // receiver: (h *Handler)
    "parseCreatePathParams",                        // method name
    []*ast.Field{Field([]string{"r"}, Star(Sel(I("http"), "Request")))},  // params
    []*ast.Field{Field(nil, Star(I("CreatePathParams"))), Field(nil, I("error"))},  // returns
    []ast.Stmt{ /* body statements */ },
)

// Produces this Go output:
func (h *Handler) parseCreatePathParams(r *http.Request) (*CreatePathParams, error) {
    // ...
}
```

## Name Formatting & Initialisms

`nameutils.go` handles OpenAPI identifiers → Go identifiers.

### `FormatGoLikeIdentifier(name string) string`

Converts any string to PascalCase with proper Go initialisms:

| Input | Output |
|---|---|
| `user_id` | `UserID` |
| `http_url` | `HTTPURL` |
| `x-request-id` | `XRequestID` |
| `created_at` | `CreatedAt` |
| `enum-val` | `EnumVal` |
| `api` | `API` |
| `json` | `JSON` |

### Common initialisms recognized

`ACL`, `API`, `ASCII`, `CPU`, `CSS`, `DNS`, `EOF`, `GUID`, `HTML`, `HTTP`, `HTTPS`, `ID`, `IP`, `JSON`, `QPS`, `RAM`, `RPC`, `SLA`, `SMTP`, `SQL`, `SSH`, `TCP`, `TLS`, `TTL`, `UDP`, `UI`, `GID`, `UID`, `UUID`, `URI`, `URL`, `UTF8`, `VM`, `XML`, `XMPP`, `XSRF`, `XSS`, `SIP`, `RTP`, `AMQP`, `DB`, `TS`

### `GoIdentLowercase(name string) string`

Same as above but first character lowercased — used for unexported identifiers:

| Input | Output |
|---|---|
| `UserID` | `userID` |
| `CreateHandler` | `createHandler` |
