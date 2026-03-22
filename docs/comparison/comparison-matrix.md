# Comparison Matrix

## Side-by-Side: The Same API in Every Generator

### What you write for `updatePet`

**ogen:**
```go
// You implement:
func (h *MyHandler) UpdatePet(ctx context.Context, req *api.UpdatePetRequest, params api.UpdatePetParams) (api.UpdatePetRes, error) {
    // req.Name is string (validated)
    // req.Price is OptString (not decimal)
    // req.Tags is []string (validated for uniqueItems)
    // params.PetId is string (from ogen router, not chi)
    pet, err := h.store.Update(ctx, params.PetId, fromOgenReq(req))
    if err != nil {
        return &api.UpdatePet404JSONResponse{Code: 404, Message: "not found"}, nil
    }
    return &toOgenPet(pet), nil  // convert *time.Time → OptDateTime, etc.
}
// Total: ~15 lines + conversion functions
// Chi middleware: incompatible
// Decimal: plain string
// Optional types: OptString wrappers
```

**oapi-codegen (strict mode):**
```go
// You implement:
func (s *PetServer) UpdatePet(ctx context.Context, request api.UpdatePetRequestObject) (api.UpdatePetResponseObject, error) {
    // request.Body.Name is string (NO validation)
    // request.Body.Price is *string (not decimal)
    // request.Body.Tags is *[]string (NO uniqueItems check)
    // request.PetId is string (from chi)
    if request.Body.Name == "" {  // manual check — spec says minLength: 1
        return api.UpdatePet400JSONResponse(api.Error{Code: 400, Message: "name required"}), nil
    }
    pet, err := s.store.Update(ctx, request.PetId, *request.Body)
    if err != nil {
        return api.UpdatePet404JSONResponse(api.Error{Code: 404, Message: "not found"}), nil
    }
    return api.UpdatePet200JSONResponse(api.Pet{...}), nil
}
// Total: ~20 lines + manual validation for every constraint
// Chi middleware: compatible
// Decimal: plain string
// Validation: manual or middleware
```

**openapi-generator:**
```go
// You implement:
func (s *PetsService) UpdatePet(ctx context.Context, petId string, req api.UpdatePetRequest) (api.ImplResponse, error) {
    // req.Name is string (basic validation only)
    // req.Price is *string (not decimal)
    // petId is a plain string arg (no struct)
    pet, err := s.store.Update(ctx, petId, req)
    if err != nil {
        return api.ImplResponse{Code: 400, Body: api.Error{...}}, nil
        // Body is interface{} — could put ANYTHING here, compiles fine
    }
    return api.ImplResponse{Code: 200, Body: pet}, nil
}
// Total: ~10 lines but zero type safety on responses
// JVM required
// Type safety: interface{} body
```

**validgo-gen:**
```go
// You implement:
func (h *petUpdateHandler) HandleUpdatePet(ctx context.Context, req *models.UpdatePetRequest) (*models.UpdatePetResponse, error) {
    // req.Body.Name — validated: non-empty, 1-100 chars
    // req.Body.Price — *decimal.Decimal
    // req.Body.Tags — validated: unique, max 10, each non-empty
    // req.Path.PetID — from chi
    pet, err := h.store.Update(ctx, req.Path.PetID, req.Body)
    if err != nil {
        if errors.Is(err, ErrNotFound) {
            return petstore.UpdatePet404Response(models.Error{Code: 404, Message: "not found"}), nil
        }
        return petstore.UpdatePet400Response(models.Error{Code: 400, Message: err.Error()}), nil
    }
    return petstore.UpdatePet200Response(models.Pet{...}), nil
}
// Total: ~15 lines, zero manual validation
// Chi middleware: fully compatible
// Decimal: native decimal.Decimal
// Validation: two-layer, automatic
// Type safety: per-status-code response types
// DI: per-operation interface
```

## When to Use What

### Use ogen when:
- Performance is your #1 priority (high-throughput API gateways, proxies)
- You don't have existing chi middleware to preserve
- You need both client and server from one spec
- You want built-in OpenTelemetry
- You're OK with `OptString` wrapper types throughout your codebase
- You need sophisticated `oneOf`/discriminator handling

### Use oapi-codegen when:
- You want maximum framework flexibility (8 routers)
- You're OK with runtime-only validation via middleware
- Your API is simple enough that manual validation is manageable
- You need template customization for unusual output patterns
- You want the largest community and most battle-tested solution

### Use openapi-generator when:
- You need to generate clients/servers in 50+ languages from one spec
- Go is a secondary language in a polyglot organization
- You already have JVM infrastructure
- Type safety in Go code is not a priority

### Use go-swagger when:
- Your specs are OpenAPI 2.0 (and will stay 2.0)
- You need the most comprehensive validation possible
- You want bidirectional spec-to-code generation

### Use validgo-gen when:
- You use **chi** and have chi middleware you can't abandon
- You need **generated request validation** (not just middleware)
- You use **`go-playground/validator`** elsewhere in your codebase
- You want **per-operation handler interfaces** for clean dependency injection
- You need **`decimal.Decimal`**, **cookie parameters**, or **client IP extraction**
- You prefer **pointer-based optionals** (`*string`) over wrapper types (`OptString`)
- You want the generated code to be **a complete, working HTTP layer** — not a skeleton you fill in

## Why We Built validgo-gen

### The core problem: no existing tool fits our architecture

Our projects have specific constraints that no off-the-shelf generator satisfies simultaneously:

| Requirement | Why | Who satisfies it? |
|---|---|---|
| Chi router | 100% of our HTTP stack (middleware, auth, rate limiting) depends on chi | oapi-codegen, openapi-generator |
| Pre-deserialization JSON validation | We need to reject malformed JSON before it hits structs (required-but-null, missing required fields get zero values in Go) | **validgo-gen**, ogen, go-swagger |
| go-playground/validator tags | Our validation library wraps go-playground/validator; all manual models already use it | **validgo-gen** only |
| Pointer-based optionals | Entire codebase uses `*string` / `*int` for optional fields; ogen's `OptString` wrappers would infect every layer | **validgo-gen**, oapi-codegen |
| Per-operation handler interfaces | Clean DI — each domain handler implements only its operations, not a monolithic interface | **validgo-gen** only |
| `decimal.Decimal` support | Payment/billing domain uses shopspring/decimal extensively | **validgo-gen** only |
| Remote-Addr as parameter | API services need client IP from `r.RemoteAddr` (rewritten by RealIP middleware) | **validgo-gen** only |
| DELETE with request body | Several endpoints use DELETE with body | **validgo-gen**, ogen, oapi-codegen |
| Cookie parameters | Session auth and web login use cookies | **validgo-gen**, ogen, oapi-codegen |
| OpenAPI 3.0 | All our specs are 3.0 | Everything except go-swagger |

**No single existing generator satisfies more than 4 of these 10 requirements.** validgo-gen satisfies all 10.

### Ecosystem compatibility

Our projects have deep investments in:

| Library | Usage | Generator compatibility |
|---|---|---|
| `go-chi/chi/v5` | All routing, all middleware | validgo-gen, oapi-codegen |
| `go-playground/validator/v10` | Validation library, manual models | **validgo-gen only** |
| `shopspring/decimal` | Billing/payment domain | **validgo-gen only** |
| `go-faster/errors` | Error wrapping throughout | validgo-gen, ogen |
| Pointer-based optionals | Every storage model, every API model | validgo-gen, oapi-codegen |

Adopting ogen would require replacing chi, converting all pointer optionals to wrappers, removing go-playground/validator, and adding decimal conversion layers. Adopting oapi-codegen would require adding middleware-based validation, losing pre-deserialization null checks, switching to a monolithic handler interface, and losing decimal/Remote-Addr support.

### Production proof

Our API services run in production with validgo-gen generating:
- **30+ sub-packages** from OpenAPI specs
- **Per-operation handler interfaces** wired via DI
- **Two-layer validation** catching malformed requests before business logic
- **Chi route registration** that plugs into our middleware stack (auth, RealIP, rate limiting, CORS, metrics)
- **Decimal fields** in payment/billing endpoints
- **Cookie parameters** in session management
- **Remote-Addr extraction** for client IP logging

## Roadmap & Limitations

### What needs to be added

| Feature | Priority | Effort | Notes |
|---|---|---|---|
| Non-string params (int, float, bool) | **P0** | Medium | Most path/query params are IDs (int64) or flags (bool) |
| `pattern` validation | P1 | Low | Regex constraints on strings — add `regexp` to validation |
| `additionalProperties` | P1 | Medium | Map types for dynamic config objects |
| `exclusiveMin/Max`, `multipleOf` | P2 | Low | Numeric constraints — extend `GetSchemaValidators()` |
| Component-level `$ref` reuse | P2 | Medium | Shared parameters, responses, headers |
| `oneOf`/`anyOf` (basic) | P3 | High | Union types for polymorphic endpoints |
| Response header generation | P3 | Low | Already partially implemented |

### What we deliberately skip

| Feature | Why skip it |
|---|---|
| Custom JSON (jx) | `encoding/json` is fast enough; jx couples to go-faster ecosystem |
| Static radix router | Chi is sufficient and battle-tested; no need to own routing |
| Client generation | Use ogen or oapi-codegen for clients; validgo-gen is server-only |
| OpenTelemetry codegen | OTel via middleware is more flexible than generated code |
| `oneOf` discriminator inference | Massive complexity (ogen's `schema_gen_sum.go` is 47KB); use tagged unions manually where needed |

### Cost-benefit vs migration

| Approach | Cost | Risk | Fit |
|---|---|---|---|
| **Extend validgo-gen** | Low-Medium | Low (already in production) | Perfect |
| Migrate to ogen | High (rewrite middleware, change optional patterns, new router) | High | Poor |
| Migrate to oapi-codegen | Medium (add validation middleware, monolithic handler) | Medium | Partial |
| Migrate to openapi-generator | High (JVM dep, non-idiomatic code, no validation) | High | Poor |
| Migrate to go-swagger | Impossible (OpenAPI 2.0 only) | N/A | N/A |

## Feature Matrix

| Feature | **validgo-gen** | ogen | oapi-codegen | go-swagger | openapi-generator |
|---|---|---|---|---|---|
| **OpenAPI version** | 3.0 | 3.0 | 3.0 | 2.0 only | 2.0 + 3.x |
| **Router** | chi | own (static) | chi + 7 others | denco | mux/chi/gin/echo |
| **JSON** | encoding/json | jx (~1 GB/s) | encoding/json | encoding/json | encoding/json |
| | | | | | |
| **VALIDATION** | | | | | |
| Body validation generated | **Two-layer** | Code-generated | **None** | Code-generated | Basic |
| Pre-deser null check | **Yes** | Yes (bitmask) | No | Yes | No |
| Struct tag validators | **go-playground** | Own | None | Own (go-openapi) | Minimal |
| `minLength/maxLength` | **Yes** | Yes | Middleware only | Yes | Partial |
| `enum` | **Yes** (oneof=) | Yes | Middleware only | Yes | No |
| `uniqueItems` | **Yes** | Yes | Middleware only | Yes | No |
| `pattern` (regex) | No (planned) | Yes (regexp2) | Middleware only | Yes | Partial |
| `exclusiveMin/Max` | No (planned) | Yes | Middleware only | Yes | No |
| | | | | | |
| **TYPES** | | | | | |
| Optional = pointer | **Yes** | No (wrappers) | Yes | Yes | Yes |
| `decimal.Decimal` | **Yes** | No | No | No | No |
| `time.Time` | **Yes** | Yes | Yes | Yes | Yes |
| Inline objects → structs | **Yes** | Yes | Yes | Yes | Yes |
| | | | | | |
| **ARCHITECTURE** | | | | | |
| Handler pattern | **1 interface/op** | 1 monolith | 1 monolith | 1 func/op | 1 func/op |
| Chi middleware compat | **Full** | None | Full | None | Partial |
| External `$ref` | **Yes** | Yes | Yes | Yes | Yes |
| Cookie params | **Yes** | Yes | Yes | Yes | No |
| Client IP extraction | **Yes** | No | No | No | No |
| DELETE with body | **Yes** (flag) | Yes | Yes | No | Yes |
| Client generation | No | Yes | Yes | Yes | Yes |
| OTel integration | No | Yes (built-in) | No | No | No |
| Security generation | No | Yes | No | Yes | No (Go) |
| `oneOf`/`anyOf` | No | Yes | Awkward | allOf only | No (Go) |
| | | | | | |
| **ECOSYSTEM** | | | | | |
| Stars | ~0 (new) | ~2K | ~8K | ~10K | ~26K |
| Written in | Go | Go | Go | Go | Java |
| Requires JVM | No | No | No | No | **Yes** |
| Code gen approach | go/ast | text/template | text/template | text/template | Mustache |
| Output formatting | gofmt by construction | goimports | goimports | gofmt | Mustache |

## Summary

Every Go OpenAPI generator makes trade-offs. The four mainstream options each leave a significant gap:

- **ogen** sacrifices router compatibility and idiomatic optional types for maximum performance
- **oapi-codegen** sacrifices validation for framework flexibility and community size
- **go-swagger** sacrifices OpenAPI 3.0 support for the best validation
- **openapi-generator** sacrifices Go idiomaticity for polyglot reach

validgo-gen occupies the intersection: **chi router + generated two-layer validation + idiomatic Go types + per-operation DI interfaces**. It doesn't try to be the fastest (ogen wins), the most flexible (oapi-codegen wins), or the most feature-complete (go-swagger wins for 2.0). It solves a specific, common problem: building validated, well-structured Go HTTP APIs from OpenAPI 3.0 specs with chi.

If that's your problem, validgo-gen is the only tool that solves it without compromises.
