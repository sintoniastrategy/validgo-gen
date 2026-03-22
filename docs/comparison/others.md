# go-swagger & openapi-generator

## go-swagger — The Validation Champion (Stuck on OpenAPI 2.0)

**GitHub**: github.com/go-swagger/go-swagger | **Stars**: ~9,956 | **License**: Apache-2.0

### What it does right

go-swagger has the **best validation** of any Go generator. Full JSON Schema Draft 4 compliance:

```go
// go-swagger generated validation — comprehensive
func (m *UpdatePetRequest) Validate(formats strfmt.Registry) error {
    if err := m.validateName(formats); err != nil {
        res = append(res, err)
    }
    if err := m.validateTags(formats); err != nil {
        res = append(res, err)
    }
    // ... every constraint checked
}

func (m *UpdatePetRequest) validateName(formats strfmt.Registry) error {
    if err := validate.RequiredString("name", "body", m.Name); err != nil {
        return err
    }
    if err := validate.MinLength("name", "body", m.Name, 1); err != nil {
        return err
    }
    if err := validate.MaxLength("name", "body", m.Name, 100); err != nil {
        return err
    }
    return nil
}
```

It also generates typed responders per status code, full security scheme handling, and uses denco's fast ternary search tree router.

### The dealbreaker

**go-swagger supports OpenAPI 2.0 only. No 3.x. No roadmap for it.**

If your specs are OpenAPI 3.0 (and they should be — 3.0 has been the standard since 2017), go-swagger is not an option. End of story.

This is tragic because its validation approach is exactly right. validgo-gen's two-layer validation was directly inspired by what go-swagger does well.

## openapi-generator — The Polyglot Giant

**GitHub**: github.com/OpenAPITools/openapi-generator | **Stars**: ~25,969 | **License**: Apache-2.0

### The Go reality

openapi-generator supports 50+ languages. Go is not one of its strengths.

```go
// openapi-generator's Go server handler signature
type PetsAPIServicer interface {
    GetPet(context.Context, string, string) (ImplResponse, error)
    UpdatePet(context.Context, string, UpdatePetRequest) (ImplResponse, error)
}

// ImplResponse — the type safety killer
type ImplResponse struct {
    Code int
    Body interface{}  // <-- any type, no compile-time checking
}

// Implementation
func (s *PetsService) UpdatePet(ctx context.Context, petId string, req UpdatePetRequest) (ImplResponse, error) {
    pet, err := s.store.UpdatePet(ctx, petId, req)
    if err != nil {
        // Nothing prevents you from returning a Pet when you should return an Error
        // The compiler cannot help you here
        return ImplResponse{Code: 400, Body: pet}, nil  // BUG: wrong type, compiles fine
    }
    return ImplResponse{Code: 200, Body: pet}, nil
}
```

### Why not for Go

- **Requires JVM** — you need Java installed to generate Go code
- **`interface{}` response body** — zero compile-time type safety
- **No security generation** — none of the 3 Go server generators support any auth scheme
- **No `allOf`/`anyOf`/`oneOf`** — no schema composition in Go generators
- **Non-idiomatic output** — Java-influenced patterns, `*string` for enums, flat package structure
- **Go is second-class** — recent release had 1 Go fix vs dozens for Java/Kotlin/TypeScript
- **5,600+ open issues** — Go generator issues compete with 50 other languages for attention
