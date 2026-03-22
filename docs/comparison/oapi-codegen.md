# oapi-codegen — The Community Favorite

**GitHub**: github.com/oapi-codegen/oapi-codegen | **Stars**: ~8,166 | **License**: Apache-2.0

## What it generates

```go
// oapi-codegen models — clean, idiomatic Go
type Pet struct {
    ID        string      `json:"id"`
    Name      string      `json:"name"`
    Price     string      `json:"price"`
    BirthDate *time.Time  `json:"birthDate,omitempty"`   // pointer = optional
    Tags      *[]string   `json:"tags,omitempty"`
    Address   *PetAddress `json:"address,omitempty"`
}

type PetAddress struct {
    City string  `json:"city"`
    Zip  *string `json:"zip,omitempty"`
}

type UpdatePetRequest struct {
    Name    string      `json:"name"`
    Price   *string     `json:"price,omitempty"`   // no decimal.Decimal support
    Tags    *[]string   `json:"tags,omitempty"`
    Address *UpdatePetRequestAddress `json:"address,omitempty"`
}

type Error struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}
```

## Non-strict mode: you handle everything

```go
// ServerInterface — raw http.ResponseWriter + *http.Request
type ServerInterface interface {
    GetPet(w http.ResponseWriter, r *http.Request, petId string, params GetPetParams)
    UpdatePet(w http.ResponseWriter, r *http.Request, petId string)
}

// Implementation — you decode JSON yourself, you encode responses yourself
type PetServer struct {
    db *sql.DB
}

func (s *PetServer) UpdatePet(w http.ResponseWriter, r *http.Request, petId string) {
    // YOU must decode the body
    var req UpdatePetRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        w.WriteHeader(400)
        json.NewEncoder(w).Encode(Error{Code: 400, Message: "invalid JSON"})
        return
    }

    // YOU must validate (oapi-codegen generates NO validators)
    // These checks don't exist in the generated code:
    // - Is name non-empty? (minLength: 1)
    // - Is name under 100 chars? (maxLength: 100)
    // - Are tags unique? (uniqueItems: true)
    // - Are tag items non-empty? (items.minLength: 1)
    // - Are there at most 10 tags? (maxItems: 10)
    // Answer: you must write all of this manually or use middleware

    // YOU must encode the response
    pet, err := s.db.UpdatePet(r.Context(), petId, req)
    if err != nil {
        w.WriteHeader(500)
        json.NewEncoder(w).Encode(Error{Code: 500, Message: err.Error()})
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(pet)
}
```

## Strict mode: typed but ceremonial

```go
// StrictServerInterface — typed request/response objects
type StrictServerInterface interface {
    GetPet(ctx context.Context, request GetPetRequestObject) (GetPetResponseObject, error)
    UpdatePet(ctx context.Context, request UpdatePetRequestObject) (UpdatePetResponseObject, error)
}

// Request wrapper
type UpdatePetRequestObject struct {
    PetId string `json:"petId"`
    Body  *UpdatePetRequest
}

// Multiple response types — one per status code
type UpdatePet200JSONResponse Pet
type UpdatePet400JSONResponse Error
type UpdatePet404JSONResponse Error

// Implementation
func (s *PetServer) UpdatePet(ctx context.Context, request api.UpdatePetRequestObject) (api.UpdatePetResponseObject, error) {
    pet, err := s.db.UpdatePet(ctx, request.PetId, *request.Body)
    if err != nil {
        if errors.Is(err, ErrNotFound) {
            return api.UpdatePet404JSONResponse(api.Error{Code: 404, Message: "not found"}), nil
        }
        return api.UpdatePet400JSONResponse(api.Error{Code: 400, Message: err.Error()}), nil
    }
    return api.UpdatePet200JSONResponse(api.Pet{
        Id:   pet.ID,
        Name: pet.Name,
        // ... map all fields
    }), nil
}

// Wiring — uses chi
func main() {
    petServer := &PetServer{db: db}
    strictHandler := api.NewStrictHandler(petServer, nil)
    r := chi.NewRouter()
    api.HandlerFromMux(strictHandler, r)
    http.ListenAndServe(":8080", r)
}
```

## What oapi-codegen does well

- **8 router frameworks**: chi, echo, gin, fiber, iris, gorilla, stdlib, echo5
- **Strict mode**: compile-time response type safety
- **Idiomatic Go**: pointers for optionals, `encoding/json` tags
- **Template overrides**: customize any generated template
- **OpenAPI Overlay**: modify specs non-invasively before generation
- **Embedded spec**: bundles the spec for runtime validation middleware
- **Most mature**: 7 years, 250+ contributors, battle-tested

## Where oapi-codegen hurts

### Problem 1: No request body validation. At all.

This is the critical gap. oapi-codegen generates **zero** validation code for request bodies. Your OpenAPI spec can say `minLength: 1, maxLength: 100, uniqueItems: true` and the generated code completely ignores it.

```go
// Your spec says: name is required, minLength 1, maxLength 100
// oapi-codegen generates this struct:
type UpdatePetRequest struct {
    Name string `json:"name"`    // No validate tag. Nothing.
}

// What happens at runtime:
// Client sends: {"name": ""}
// json.Unmarshal succeeds, Name == ""
// Your spec says minLength: 1 — violated. No error.

// Client sends: {}
// json.Unmarshal succeeds, Name == ""
// Your spec says required: true — violated. No error.

// Client sends: {"name": null}
// json.Unmarshal succeeds, Name == ""
// Your spec says required, non-nullable — violated. No error.
```

The recommended workaround is validation middleware:

```go
swagger, _ := api.GetSwagger()
r.Use(chimiddleware.OapiRequestValidator(swagger))
```

This re-parses the entire incoming request against the full OpenAPI spec using `kin-openapi`. For every request. It works, but:

- **Runtime overhead**: the full spec is walked for every request
- **No response validation**: only incoming requests are checked
- **Mismatch risk**: the middleware validates against the spec, but your Go structs have no constraints — if someone bypasses the middleware (tests, internal calls), they get zero validation
- **No struct-level safety net**: after deserialization, there's nothing stopping invalid data from flowing through your business logic

### Problem 2: Monolithic handler interface

Same as ogen — one `ServerInterface` for all operations:

```go
type ServerInterface interface {
    GetPet(w http.ResponseWriter, r *http.Request, petId string, params GetPetParams)
    UpdatePet(w http.ResponseWriter, r *http.Request, petId string)
    DeletePet(w http.ResponseWriter, r *http.Request, petId string)
    ListPets(w http.ResponseWriter, r *http.Request, params ListPetsParams)
    CreatePet(w http.ResponseWriter, r *http.Request)
    GetOwner(w http.ResponseWriter, r *http.Request, ownerId string)
    // ... every operation on one interface
}
```

You can embed `Unimplemented` for incremental implementation, but you cannot decompose this into independent domain handlers for dependency injection. All endpoints share one implementation struct.

### Problem 3: No `decimal.Decimal` support

`format: decimal` in your spec produces a plain `string`. You must convert manually:

```go
// oapi-codegen generates:
type Pet struct {
    Price string `json:"price"`  // just a string
}

// You must convert in every handler:
func (s *PetServer) GetPet(...) {
    price, err := decimal.NewFromString(dbPet.Price.String())
    // error handling...
    return api.GetPet200JSONResponse(api.Pet{
        Price: price.String(), // and convert back
    }), nil
}
```
