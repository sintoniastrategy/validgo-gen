# ogen — The Performance King

**GitHub**: github.com/ogen-go/ogen | **Stars**: ~2,028 | **License**: Apache-2.0

## What it generates

ogen produces ~20 files from a single spec. Here's what the models look like:

```go
// ogen's generated Pet type
type Pet struct {
    ID        string          `json:"id"`
    Name      string          `json:"name"`
    Price     string          `json:"price"`
    BirthDate OptDateTime     `json:"birthDate"`  // NOT *time.Time
    Tags      []string        `json:"tags"`
    Address   OptPetAddress   `json:"address"`     // NOT *PetAddress
}

// ogen's optional wrapper — used instead of pointers
type OptDateTime struct {
    Value time.Time
    Set   bool
}

type OptPetAddress struct {
    Value PetAddress
    Set   bool
}

// To check if a field was sent:
if pet.BirthDate.Set {
    fmt.Println(pet.BirthDate.Value)
}

// Compare with idiomatic Go:
if pet.BirthDate != nil {
    fmt.Println(*pet.BirthDate)
}
```

## How you implement handlers

```go
// ogen generates ONE interface for ALL operations
type Handler interface {
    GetPet(ctx context.Context, params GetPetParams) (GetPetRes, error)
    UpdatePet(ctx context.Context, req *UpdatePetRequest, params UpdatePetParams) (UpdatePetRes, error)
}

// You implement it on a single struct
type PetHandler struct {
    db *sql.DB
}

func (h *PetHandler) GetPet(ctx context.Context, params api.GetPetParams) (api.GetPetRes, error) {
    pet, err := h.db.GetPet(ctx, params.PetID)
    if err != nil {
        return &api.GetPetNotFound{Code: 404, Message: "not found"}, nil
    }
    return &api.Pet{ID: pet.ID, Name: pet.Name}, nil
}

func (h *PetHandler) UpdatePet(ctx context.Context, req *api.UpdatePetRequest, params api.UpdatePetParams) (api.UpdatePetRes, error) {
    // ...
}

// Wiring — ogen has its OWN router, NOT chi
func main() {
    handler := &PetHandler{db: db}
    srv, err := api.NewServer(handler, securityHandler)
    http.ListenAndServe(":8080", srv) // srv IS the http.Handler
}
```

## What ogen does well

- **Fastest router**: code-generated static radix tree, ~18 ns/op, 0 allocations
- **Fastest JSON**: uses `jx` (go-faster jsoniter fork), ~1 GB/s parsing
- **Code-generated validation**: bitmask required-field tracking, no reflection
- **Client + server**: generates matching typed client from the same spec
- **OpenTelemetry**: built-in tracing and metrics per operation
- **`oneOf`/discriminator**: most sophisticated union type inference of any Go generator
- **Streaming JSON**: `x-ogen-json-streaming` for large payloads

## Where ogen hurts

### Problem 1: It replaces your router

ogen generates its own static radix router. You cannot use chi, echo, or any other router. If your project has chi middleware:

```go
// Your existing middleware stack — ALL of this needs rethinking with ogen
r := chi.NewRouter()
r.Use(middleware.RealIP)
r.Use(middleware.RequestID)
r.Use(middleware.Recoverer)
r.Use(middleware.Timeout(30 * time.Second))
r.Use(corsMiddleware)
r.Use(rateLimiter.Middleware)

// Per-route middleware
r.Route("/admin", func(r chi.Router) {
    r.Use(adminAuth)
    r.Get("/users", listUsers)
})
```

With ogen, you'd need to convert all of this to either `net/http` middleware wrapping the entire server, or ogen's experimental middleware API:

```go
// ogen middleware — experimental, different signature, limited
func ogenMiddleware(
    req middleware.Request,
    next func(req middleware.Request) (middleware.Response, error),
) (middleware.Response, error) {
    // No access to chi.RouteContext
    // No chi.URLParam()
    // No r.Context().Value(chi.RouteCtxKey)
    // Different error handling model
    return next(req)
}
```

### Problem 2: `OptString` wrappers infect your entire codebase

Every layer that touches ogen types must use `OptString`, `OptInt`, `OptDateTime`, etc. This means your storage layer, your business logic, and your API layer all speak different languages:

```go
// Storage model (idiomatic Go)
type Pet struct {
    Name      string
    BirthDate *time.Time    // pointer = optional
    Tags      []string
}

// ogen API model
type APIPet struct {
    Name      string
    BirthDate OptDateTime   // wrapper = optional
    Tags      []string
}

// You need conversion everywhere
func toAPIPet(p Pet) APIPet {
    result := APIPet{Name: p.Name}
    if p.BirthDate != nil {
        result.BirthDate = OptDateTime{Value: *p.BirthDate, Set: true}
    }
    return result
}

func fromAPIPet(p APIPet) Pet {
    result := Pet{Name: p.Name}
    if p.BirthDate.Set {
        result.BirthDate = &p.BirthDate.Value
    }
    return result
}

// Repeat for EVERY optional field, EVERY model, EVERY endpoint.
// In a 50-endpoint API with 200 models, this is thousands of lines of conversion code.
```

### Problem 3: Monolithic handler interface doesn't scale

With 50+ endpoints, you implement one massive interface:

```go
type Handler interface {
    GetPet(...)
    UpdatePet(...)
    DeletePet(...)
    ListPets(...)
    CreatePet(...)
    GetOwner(...)
    UpdateOwner(...)
    // ... 45 more methods
}

// One struct must implement ALL of them
type MyHandler struct {
    petService   PetService
    ownerService OwnerService
    orderService OrderService
    // ... all dependencies for all domains
}
```

`x-ogen-operation-group` helps by splitting into sub-interfaces, but they're still composed into one `Handler` interface. You can't wire pet handlers and owner handlers as independent components.

### Problem 4: Requires latest Go version

ogen's `go.mod` specifies `go 1.25.0`, which can be a hard blocker for teams on LTS or enterprise Go versions.
