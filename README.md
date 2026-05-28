# validgo-gen

**OpenAPI 3.0 → Go code generator with two-layer validation**

validgo-gen reads OpenAPI 3.0 YAML specs and generates a complete, validated Go HTTP layer — models with `go-playground/validator` tags and chi-based handlers with per-operation interfaces.

## The problem

Go's `encoding/json` silently turns missing fields, explicit nulls, and empty strings into the same zero value. Every OpenAPI generator for Go either skips validation entirely (oapi-codegen), couples it to a non-standard runtime (ogen), or doesn't support OpenAPI 3.0 (go-swagger).

## The solution

validgo-gen generates **two-layer validation**:

1. **Layer 1 (pre-deserialization)** — checks raw JSON for missing required fields, explicit nulls, and nested structure before `json.Unmarshal` ever runs
2. **Layer 2 (post-deserialization)** — `go-playground/validator` struct tags enforce constraints like `min`, `max`, `oneof`, `email`, `unique`

## Key features

| Feature | Detail |
|---|---|
| **Chi-native routing** | Generates `chi.Router` integration — works with your existing middleware stack |
| **Two-layer validation** | Pre-deserialization JSON checks + struct tag validation |
| **Per-operation interfaces** | One Go interface per operation — clean dependency injection, no monolithic handler |
| **Idiomatic Go types** | `*string` for optionals, `decimal.Decimal` for decimals, `time.Time` for dates |
| **go-playground/validator** | Standard validation library — same tags you already use |
| **Go AST generation** | Code built as `go/ast` nodes, formatted via `go/format` — always valid, always `gofmt` |

## Quick start

```bash
go build -o validgo-gen ./cmd/generate.go

./validgo-gen \
  -spec api/petstore.yaml \
  -out internal/petstore \
  -name petstore
```

This generates two files:

| File | Package | Contains |
|---|---|---|
| `models.go` | `petstoremodels` | Request/response structs with `json` + `validate` tags |
| `handlers.go` | `petstore` | Handler interfaces, chi routes, request parsing, JSON validation, response writing |

## What you write vs what's generated

```go
// This is ALL you write — everything else is generated
func (h *petHandler) HandleUpdatePet(
    ctx context.Context,
    req *models.UpdatePetRequest,
) (*models.UpdatePetResponse, error) {
    // req.Body.Name — validated: non-empty, 1-100 chars
    // req.Body.Price — *decimal.Decimal (or nil)
    // req.Body.Tags — validated: unique, max 10, each non-empty
    // req.Path.PetID — extracted from chi URL params

    pet, err := h.store.Update(ctx, req.Path.PetID, req.Body)
    if err != nil {
        return petstore.UpdatePet404Response(models.Error{
            Code: 404, Message: "not found",
        }), nil
    }
    return petstore.UpdatePet200Response(pet), nil
}
```

## Error envelope

Every internal failure emitted by a generated handler — request parse errors,
unsupported `Content-Type`, the handler returning `nil`, and JSON encode
failures inside the response writers — is routed through a single
`ErrorHandler` callback. The default reproduces the body the generator
emitted before the hook existed:

```
HTTP/1.1 400 Bad Request
Content-Type: text/plain; charset=utf-8

{"error":"field email is required"}
```

(Body shape is `{"error":"<msg>"}`; status comes from the call site. 500
sites pass `"InternalServerError"` and never leak the original
`err.Error()`; 415 sites pass `"Unsupported Content-Type"`.)

### Customizing the envelope

Each generated package exposes:

```go
type ErrorHandler = func(w http.ResponseWriter, r *http.Request, status int, msg string)

func WithErrorHandler(eh ErrorHandler) Option
func (h *Handler) SetErrorHandler(eh ErrorHandler)
var DefaultErrorHandler ErrorHandler
```

`ErrorHandler` is a **type alias** (not a named type), so a single bare
`func(http.ResponseWriter, *http.Request, int, string)` value assigns to
every package's setter without a conversion — which is what makes the
aggregator-loop pattern below possible.

**Constructor option** — idiomatic for one-off handlers:

```go
api.NewHandler(impl, api.WithErrorHandler(func(w http.ResponseWriter, r *http.Request, status int, msg string) {
    utils.WriteErr(w, r, utils.ErrorResponse{
        Code:    canonicalCode(status),
        Message: msg,
        ReqID:   chimw.GetReqID(r.Context()),
        TraceID: tracing.FromCtx(r.Context()),
    })
}))
```

**Wrapping the default** — extra behavior, same body:

```go
api.NewHandler(impl, api.WithErrorHandler(func(w http.ResponseWriter, r *http.Request, status int, msg string) {
    if status >= 500 {
        slog.ErrorContext(r.Context(), "api 5xx", "status", status, "msg", msg, "path", r.URL.Path)
    }
    api.DefaultErrorHandler(w, r, status, msg)
}))
```

## Documentation

- **[Design & Usage](docs/design/)** — Full architecture reference: code generation pipeline, AST helpers, two-layer validation, OpenAPI→validator tag mapping, handler interfaces, and test strategy.

- **[Why validgo-gen](docs/comparison/)** — In-depth comparison of ogen, oapi-codegen, go-swagger, and openapi-generator with code examples — and why none of them solve the validation problem.

## TODO

- **AST builder refactor** (`feat/prepare-for-vibecoding`) — Replace direct `go/ast` node construction in `handlers.go` and `schemas.go` with a structured builder API (`internal/generator/astbuilder/`). The current generator builds AST nodes inline, which works but makes the code hard for AI coding assistants to modify safely. The builder provides a higher-level API (struct builder, handler builder, validation builder, etc.) so that future feature work — especially AI-assisted — can manipulate generation logic without understanding raw `go/ast` internals.

## License

MIT (Usishchev Yury, 2025)
