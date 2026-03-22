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

## Documentation

- **[Design & Usage](docs/design/)** — Full architecture reference: code generation pipeline, AST helpers, two-layer validation, OpenAPI→validator tag mapping, handler interfaces, and test strategy.

- **[Why validgo-gen](docs/comparison/)** — In-depth comparison of ogen, oapi-codegen, go-swagger, and openapi-generator with code examples — and why none of them solve the validation problem.

## License

MIT (Usishchev Yury, 2025)
