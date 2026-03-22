# Generated Handlers

## Handler interface (one per operation)

```go
type CreateHandler interface {
    HandleCreate(ctx context.Context, req *apimodels.CreateRequest) (*apimodels.CreateResponse, error)
}
```

Each operation gets its own interface — implement only what you need.

## Handler struct & constructor

```go
type Handler struct {
    validator     *validator.Validate
    createHandler CreateHandler
    // ... one field per operation
}

func NewHandler(createHandler CreateHandler, /* ... */) *Handler {
    return &Handler{
        validator:     validator.New(validator.WithRequiredStructEnabled()),
        createHandler: createHandler,
        // ...
    }
}
```

## Route registration

```go
func (h *Handler) AddRoutes(router chi.Router) {
    router.Post("/path/to/{param}/resource{suffix}", h.handleCreate)
    // ... one line per operation
}
```

## Per-operation generated methods

For an operation `Create`, the generator produces these methods on `*Handler`:

| Method | Purpose |
|---|---|
| `parseCreatePathParams(r)` | Extract chi URL params, validate |
| `parseCreateQueryParams(r)` | Extract query string values |
| `parseCreateHeaders(r)` | Extract HTTP headers (with date-time parsing) |
| `parseCreateCookies(r)` | Extract cookies (required vs optional) |
| `parseCreateRequestBody(r)` | Decode JSON → raw validate → unmarshal → struct validate |
| `parseCreateRequest(r)` | Orchestrate all parse methods → `*CreateRequest` |
| `handleCreate(w, r)` | Content-type switch → delegates to `handleCreateRequest` |
| `handleCreateRequest(w, r)` | Parse → call handler → write response |
| `writeCreateResponse(w, resp)` | Status code switch → per-code writer |
| `writeCreate200Response(w, resp)` | JSON encode + set headers for 200 |
| `Create200Response(body)` | Convenience constructor: `&CreateResponse{StatusCode: 200, Response200: &CreateResponse200{Body: body}}` |

## Example: implementing a handler

```go
type myCreateHandler struct {
    db *sql.DB
}

func (h *myCreateHandler) HandleCreate(
    ctx context.Context,
    req *apimodels.CreateRequest,
) (*apimodels.CreateResponse, error) {
    // All parsing and validation already done.
    // req.Path, req.Query, req.Headers, req.Cookies, req.Body are populated.

    result, err := h.db.ExecContext(ctx, "INSERT INTO resources ...", req.Body.Name)
    if err != nil {
        return nil, err
    }

    return api.Create200Response(apimodels.NewResourceResponse{
        Name: req.Body.Name,
    }), nil
}
```

## Wiring it up with chi

```go
func main() {
    r := chi.NewRouter()

    handler := api.NewHandler(
        &myCreateHandler{db: db},
        // ... other operation handlers
    )
    handler.AddRoutes(r)

    http.ListenAndServe(":8080", r)
}
```
