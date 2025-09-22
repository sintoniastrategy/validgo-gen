# OpenAPI Code Generator for Go

A powerful OpenAPI 3.0 code generator that creates type-safe Go code from OpenAPI specifications. This tool generates HTTP handlers, data models, and validation code to help you build robust REST APIs in Go.

## Features

- **Type-Safe Models**: Generates Go structs with proper JSON tags and validation rules
- **HTTP Handlers**: Creates Chi router-compatible handlers with automatic parameter parsing
- **Request/Response Processing**: Handles path parameters, query parameters, headers, cookies, and request bodies
- **Validation**: Integrates with go-playground/validator for comprehensive input validation
- **External References**: Supports `$ref` to external YAML files for reusable schemas
- **Multiple Content Types**: Currently supports `application/json` content type
- **Flexible Configuration**: Command-line options for customization

## Generated Code Structure

For each OpenAPI specification, the generator creates:

- **Models Package** (`*models`): Go structs representing request/response schemas
- **Handlers Package**: HTTP handlers with routing and parameter parsing
- **Type Safety**: All generated code is strongly typed with proper validation

## Installation

```bash
go install github.com/jolfzverb/codegen/cmd/generate@latest
```

## Usage

### Command Line

```bash
# Basic usage
go run cmd/generate.go -d ./generated -p github.com/yourorg/api api.yaml

# With additional options
go run cmd/generate.go \
  -d ./generated \
  -p github.com/yourorg/api \
  --pointers \
  --allow-delete-with-body \
  api.yaml external-definitions.yml
```

### Command Line Options

- `-d, --dir-prefix`: Directory prefix for generated files (default: "internal")
- `-p, --package-prefix`: Package prefix for imports (default: "internal")
- `--pointers`: Generate required fields as pointers
- `--allow-delete-with-body`: Allow DELETE operations with request body
- `--allow-remote-addr-param`: Allow RemoteAddr fake parameter

### Go Generate Integration

Add to your Go files:

```go
//go:generate go run github.com/jolfzverb/codegen/cmd/generate.go -d ./generated -p github.com/yourorg/api api.yaml
```

Then run:

```bash
go generate ./...
```

## Example

Given an OpenAPI specification like:

```yaml
openapi: 3.0.0
info:
  title: Sample API
  version: 1.0.0

paths:
  /users/{id}:
    get:
      operationId: get_user
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
      responses:
        '200':
          description: User found
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/User'

components:
  schemas:
    User:
      type: object
      properties:
        id:
          type: string
        name:
          type: string
      required:
        - id
        - name
```

The generator creates:

**Models** (`usermodels/models.go`):
```go
type GetUserPathParams struct {
	ID string `json:"id" validate:"required"`
}
type GetUserRequest struct {
	Path GetUserPathParams
}
type GetUserResponse200 struct {
	Body User
}
type GetUserResponse struct {
	StatusCode  int
	Response200 *GetUserResponse200
}
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
```

**Handlers** (`handlers.go`):
```go
type GetUserHandler interface {
	HandleGetUser(ctx context.Context, r usersmodels.GetUserRequest) (*usersmodels.GetUserResponse, error)
}
type Handler struct {
	validator *validator.Validate
	getUser   GetUserHandler
}

func NewHandler(getUser GetUserHandler) *Handler {
	return &Handler{validator: validator.New(validator.WithRequiredStructEnabled()), getUser: getUser}
}
func (h *Handler) AddRoutes(router chi.Router) {
	router.Get("/users/{id}", h.handleGetUser)
}
```

## Supported OpenAPI Features

### ✅ Currently Supported
- OpenAPI 3.0 specification format
- Path parameters, query parameters, headers, cookies
- Request/response bodies with JSON content type
- Schema validation (min/max length, required fields, enums)
- External schema references (`$ref`)
- Multiple HTTP methods (GET, POST, PUT, PATCH, DELETE)
- Response headers
- Nested objects and arrays
- Decimal number support via shopspring/decimal

### 🚧 Planned Features
- Integer/number parameters
- Additional properties in schemas
- Components: parameters, request bodies, responses, headers
- Enhanced external reference support
- Additional content types beyond JSON

## Project Structure

```
├── cmd/
│   └── generate.go          # Main CLI entry point
├── internal/
│   ├── generator/           # Core generation logic
│   │   ├── generatehandlers.go
│   │   ├── generateschemas.go
│   │   └── options/
│   └── usage/               # Example usage and generated code
├── test/                    # Test files and examples
└── Makefile                 # Build and test commands
```

## Development

### Prerequisites
- Go 1.24 or later
- golangci-lint (for linting)

### Building and Testing

```bash
# Run tests
make test

# Run linter
make lint

# Generate example code
make generate

# Run all checks
make check
```

### Dependencies

- [kin-openapi](https://github.com/getkin/kin-openapi) - OpenAPI 3.0 parsing
- [go-chi/chi](https://github.com/go-chi/chi) - HTTP router
- [go-playground/validator](https://github.com/go-playground/validator) - Validation
- [shopspring/decimal](https://github.com/shopspring/decimal) - Decimal number support

## License

This project is licensed under the terms specified in the LICENSE file.

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests to improve the generator's functionality and OpenAPI specification support.