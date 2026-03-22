# The Problem None of Them Solve

Let's send a malicious but common request to an `updatePet` endpoint:

```json
{
  "name": null,
  "price": "not-a-decimal",
  "tags": ["", "", "valid"],
  "address": {
    "zip": "not-a-zip"
  }
}
```

## What each generator does with this request

**ogen**: Catches `name: null` (bitmask tracking). Catches duplicate tags (`uniqueItems`). Does not support `decimal` format. Uses its own router — can't use chi middleware. Wraps optionals in `OptString` instead of pointers.

**oapi-codegen (without middleware)**: `json.Unmarshal` succeeds silently. `name` becomes `""` (zero value). `price` is just a string, so `"not-a-decimal"` is accepted. Empty strings in `tags` are accepted. `zip` pattern is ignored. **Every constraint in your spec is violated. No errors.**

**oapi-codegen (with middleware)**: The middleware catches these at the HTTP layer by re-validating against the spec. But after the middleware, the Go structs have no constraints — any code path that doesn't go through the middleware (tests, internal service calls, queue consumers) gets zero validation.

**go-swagger**: Would catch everything — but doesn't support OpenAPI 3.0.

**openapi-generator**: Basic constraint checks only. `interface{}` responses. Requires JVM.

## The three-way impossible choice

For Go teams using OpenAPI 3.0 with chi:

1. **ogen**: Best validation, wrong router, non-idiomatic types
2. **oapi-codegen**: Right router, idiomatic types, no validation
3. **Neither**: Write everything by hand

This is the gap validgo-gen fills.

## validgo-gen — Validation-First Code Generation

validgo-gen generates **two files per spec**: models (with `go-playground/validator` tags) and handlers (with chi routes and two-layer validation).

### Generated models

```go
// validgo-gen models — idiomatic Go + validation tags
type UpdatePetRequestBody struct {
    Name    string                        `json:"name" validate:"required,min=1,max=100"`
    Price   *decimal.Decimal              `json:"price,omitempty" validate:"omitempty"`
    Tags    []string                      `json:"tags,omitempty" validate:"omitempty,max=10,unique,dive,min=1"`
    Address *UpdatePetRequestBodyAddress  `json:"address,omitempty" validate:"omitempty"`
}

type UpdatePetRequestBodyAddress struct {
    City string  `json:"city"`
    Zip  *string `json:"zip,omitempty" validate:"omitempty"`
}

type UpdatePetPathParams struct {
    PetID string `json:"petId"`
}

// Composite request — everything in one place
type UpdatePetRequest struct {
    Path UpdatePetPathParams
    Body UpdatePetRequestBody
}

// Response with typed status-code variants
type UpdatePetResponse struct {
    StatusCode  int
    Response200 *UpdatePetResponse200
    Response400 *UpdatePetResponse400
    Response404 *UpdatePetResponse404
}

type UpdatePetResponse200 struct {
    Body Pet
}

type UpdatePetResponse400 struct {
    Body Error
}
```

Key differences from other generators:
- `validate` tags on every constrained field
- `decimal.Decimal` for `format: decimal`
- Pointer-based optionals (idiomatic Go)
- Nested inline objects get their own struct
- Composite request aggregates path/query/headers/cookies/body

### Generated handlers — per-operation interfaces

```go
// ONE interface per operation — not a monolithic mega-interface
type GetPetHandler interface {
    HandleGetPet(ctx context.Context, req *petstoremodels.GetPetRequest) (*petstoremodels.GetPetResponse, error)
}

type UpdatePetHandler interface {
    HandleUpdatePet(ctx context.Context, req *petstoremodels.UpdatePetRequest) (*petstoremodels.UpdatePetResponse, error)
}

// Handler struct accepts each interface independently
type Handler struct {
    validator        *validator.Validate
    getPetHandler    GetPetHandler
    updatePetHandler UpdatePetHandler
}

func NewHandler(
    getPetHandler    GetPetHandler,
    updatePetHandler UpdatePetHandler,
) *Handler {
    return &Handler{
        validator:        validator.New(validator.WithRequiredStructEnabled()),
        getPetHandler:    getPetHandler,
        updatePetHandler: updatePetHandler,
    }
}

// Routes register on chi.Router
func (h *Handler) AddRoutes(router chi.Router) {
    router.Get("/pets/{petId}", h.handleGetPet)
    router.Put("/pets/{petId}", h.handleUpdatePet)
}
```

### Two-layer validation (the key differentiator)

**Layer 1: Raw JSON validation — before `json.Unmarshal`**

```go
// Generated function — operates on json.RawMessage, not structs
func ValidateUpdatePetRequestBodyJSON(data json.RawMessage) error {
    mapData := make(map[string]json.RawMessage)
    if err := json.Unmarshal(data, &mapData); err != nil {
        return errors.Wrap(err, "UpdatePetRequestBody")
    }

    // Required field "name" — must exist
    if _, ok := mapData["name"]; !ok {
        return errors.New("UpdatePetRequestBody: field 'name' is required")
    }

    // Required field "name" — must not be null
    if containsNull(mapData["name"]) {
        return errors.New("UpdatePetRequestBody: field 'name' must not be null")
    }

    // Nested object "address" — validate recursively if present
    if addressData, ok := mapData["address"]; ok && !containsNull(addressData) {
        if err := ValidateUpdatePetRequestBodyAddressJSON(addressData); err != nil {
            return errors.Wrap(err, "UpdatePetRequestBody.address")
        }
    }

    return nil
}

func ValidateUpdatePetRequestBodyAddressJSON(data json.RawMessage) error {
    mapData := make(map[string]json.RawMessage)
    if err := json.Unmarshal(data, &mapData); err != nil {
        return errors.Wrap(err, "UpdatePetRequestBodyAddress")
    }
    if _, ok := mapData["city"]; !ok {
        return errors.New("UpdatePetRequestBodyAddress: field 'city' is required")
    }
    if containsNull(mapData["city"]) {
        return errors.New("UpdatePetRequestBodyAddress: field 'city' must not be null")
    }
    return nil
}
```

**Why this matters — the Go JSON problem:**

```go
type Request struct {
    Name string `json:"name"`  // required in spec
}

// All three of these inputs produce Name == "" after json.Unmarshal:
json.Unmarshal([]byte(`{}`), &req)              // missing field
json.Unmarshal([]byte(`{"name": null}`), &req)  // explicit null
json.Unmarshal([]byte(`{"name": ""}`), &req)    // empty string

// Without pre-deserialization validation, you CANNOT tell these apart.
// But they have very different semantic meanings:
// - Missing: client forgot the field → 400 "field 'name' is required"
// - Null: client explicitly nulled it → 400 "field 'name' must not be null"
// - Empty: client sent empty string → 400 from validator "min=1" tag
```

Layer 1 runs on raw JSON, distinguishing all three cases. Layer 2 (struct tags) catches constraint violations after deserialization.

**Layer 2: Struct tag validation — after `json.Unmarshal`**

```go
// Generated parse method — the full pipeline
func (h *Handler) parseUpdatePetRequestBody(r *http.Request) (*petstoremodels.UpdatePetRequestBody, error) {
    // Read raw body
    data, err := io.ReadAll(r.Body)
    if err != nil {
        return nil, errors.Wrap(err, "read body")
    }

    // Layer 1: raw JSON structure validation
    if err := ValidateUpdatePetRequestBodyJSON(data); err != nil {
        return nil, err  // → 400: "field 'name' is required" / "must not be null"
    }

    // Deserialize
    var body petstoremodels.UpdatePetRequestBody
    if err := json.Unmarshal(data, &body); err != nil {
        return nil, errors.Wrap(err, "unmarshal body")
    }

    // Layer 2: struct constraint validation
    if err := h.validator.Struct(body); err != nil {
        return nil, err  // → 400: "min=1" / "max=100" / "unique" / "oneof=..."
    }

    return &body, nil
}
```

### The full generated handler lifecycle

```go
// Generated — you never write this code
func (h *Handler) handleUpdatePetRequest(w http.ResponseWriter, r *http.Request) {
    // Parse everything — path params, query params, headers, cookies, body
    req, err := h.parseUpdatePetRequest(r)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
        return
    }

    // Call YOUR handler — the only code you write
    resp, err := h.updatePetHandler.HandleUpdatePet(r.Context(), req)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
        return
    }

    // Write response — dispatches by status code
    h.writeUpdatePetResponse(w, resp)
}

func (h *Handler) writeUpdatePetResponse(w http.ResponseWriter, resp *petstoremodels.UpdatePetResponse) {
    switch resp.StatusCode {
    case 200:
        h.writeUpdatePet200Response(w, resp.Response200)
    case 400:
        h.writeUpdatePet400Response(w, resp.Response400)
    case 404:
        h.writeUpdatePet404Response(w, resp.Response404)
    }
}

// Convenience constructors
func UpdatePet200Response(body petstoremodels.Pet) *petstoremodels.UpdatePetResponse {
    return &petstoremodels.UpdatePetResponse{
        StatusCode:  200,
        Response200: &petstoremodels.UpdatePetResponse200{Body: body},
    }
}
```

### Your handler implementation

```go
// This is ALL you write. Everything else is generated.
type petUpdateHandler struct {
    petStore PetStore
}

func (h *petUpdateHandler) HandleUpdatePet(
    ctx context.Context,
    req *petstoremodels.UpdatePetRequest,
) (*petstoremodels.UpdatePetResponse, error) {
    // req.Path.PetID — already extracted from chi URL params
    // req.Body.Name — already validated: non-empty, 1-100 chars
    // req.Body.Price — already a *decimal.Decimal (or nil if not sent)
    // req.Body.Tags — already validated: unique, max 10, each non-empty
    // req.Body.Address.City — already validated: required if address is present

    pet, err := h.petStore.Update(ctx, req.Path.PetID, req.Body)
    if err != nil {
        if errors.Is(err, ErrNotFound) {
            return petstore.UpdatePet404Response(petstoremodels.Error{
                Code: 404, Message: "pet not found",
            }), nil
        }
        return petstore.UpdatePet400Response(petstoremodels.Error{
            Code: 400, Message: err.Error(),
        }), nil
    }

    return petstore.UpdatePet200Response(petstoremodels.Pet{
        ID:   pet.ID,
        Name: pet.Name,
        // ...
    }), nil
}
```

### Wiring with dependency injection

Because each operation is its own interface, you can wire handlers independently:

```go
func main() {
    r := chi.NewRouter()

    // Standard chi middleware — works perfectly
    r.Use(middleware.RealIP)
    r.Use(middleware.RequestID)
    r.Use(middleware.Recoverer)
    r.Use(middleware.Timeout(30 * time.Second))

    // Each handler is independent — different struct, different deps
    getPetHandler := &petGetHandler{
        petStore: petStore,
        cache:    cache,
    }
    updatePetHandler := &petUpdateHandler{
        petStore:  petStore,
        publisher: eventPublisher,
    }

    // Wire them into the generated handler
    handler := petstore.NewHandler(
        getPetHandler,
        updatePetHandler,
    )
    handler.AddRoutes(r)

    http.ListenAndServe(":8080", r)
}
```

Compare with ogen/oapi-codegen where **one struct** must hold **all** dependencies:

```go
// ogen / oapi-codegen pattern — everything on one struct
type MonolithicHandler struct {
    petStore     PetStore
    ownerStore   OwnerStore
    orderStore   OrderStore
    cache        Cache
    publisher    EventPublisher
    emailSender  EmailSender
    // Every dependency for every endpoint lives here
}

// With DI frameworks like wire, dig, fx — this struct becomes the bottleneck.
// Adding a new endpoint means adding deps to the monolith.
// Testing one endpoint means constructing the entire monolith.
```

With validgo-gen's per-operation interfaces:

```go
// Testing is trivial — mock only what this handler needs
func TestUpdatePet(t *testing.T) {
    mockStore := &MockPetStore{...}
    handler := petstore.NewHandler(
        nil,  // getPet — not needed for this test
        &petUpdateHandler{petStore: mockStore},
    )
    // Only the updatePet handler is real; everything else is nil.
    // The test is focused and fast.
}
```
