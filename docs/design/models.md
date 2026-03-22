# Generated Models

## Request parameter structs

For an operation `Create` with path params, query params, headers, cookies:

```go
type CreatePathParams struct {
    Param  string `json:"param"`
    Suffix string `json:"suffix" validate:"oneof=json xml yaml"`
}

type CreateQueryParams struct {
    Search *string `json:"search,omitempty" validate:"omitempty,min=3,max=100"`
}

type CreateHeaders struct {
    XRequestID string `json:"X-Request-ID"`
}

type CreateCookies struct {
    SessionID string `json:"session_id"`
    Theme     *string `json:"theme,omitempty" validate:"omitempty"`
}
```

**Rules:**
- Required fields → value types, no `omitempty`
- Optional fields → pointer types, `omitempty` in json tag, `omitempty` prefix in validate tag
- When `-pointers` flag is set, required fields are also pointers

## Request body structs

```go
type CreateRequestBody struct {
    Name        string                      `json:"name" validate:"required,min=1,max=255"`
    Tags        []string                    `json:"tags" validate:"required,min=1,max=10,unique,dive,min=1,max=50"`
    Metadata    *CreateRequestBodyMetadata  `json:"metadata,omitempty" validate:"omitempty"`
    Amount      decimal.Decimal             `json:"amount"`
    CreatedAt   time.Time                   `json:"created_at"`
}

type CreateRequestBodyMetadata struct {
    Key   string `json:"key"`
    Value string `json:"value"`
}
```

**Inline objects** (anonymous objects defined directly in the schema) get their own struct with a name derived from the parent field: `<Operation>RequestBody<FieldName>`.

## Composite request/response structs

```go
type CreateRequest struct {
    Path    CreatePathParams
    Query   CreateQueryParams
    Headers CreateHeaders
    Cookies CreateCookies
    Body    CreateRequestBody
}

type CreateResponse struct {
    StatusCode  int
    Response200 *CreateResponse200
    Response400 *CreateResponse400
    Response404 *CreateResponse404
}

type CreateResponse200 struct {
    Body SomeResponseModel
}
```

## Type aliases

Simple schemas (non-object, non-array):
```go
type ExternalRef string
type StatusCode int
```

Array schemas:
```go
type ItemList []Item
type StringList []string
```

## Type mapping

| OpenAPI type + format | Go type |
|---|---|
| `string` | `string` |
| `string` + `date-time` | `time.Time` |
| `string` + `decimal` | `decimal.Decimal` |
| `string` + `ip` | `string` (validate: `ip`) |
| `string` + `ipv4` | `string` (validate: `ipv4`) |
| `string` + `ipv6` | `string` (validate: `ipv6`) |
| `string` + `email` | `string` (validate: `email`) |
| `integer` | `int` |
| `integer` + `int8/16/32/64` | `int8/16/32/64` |
| `integer` + `uint8/16/32/64` | `uint8/16/32/64` |
| `number` | `float64` |
| `number` + `float` | `float32` |
| `number` + `double` | `float64` |
| `boolean` | `bool` |
| `object` | Generated struct |
| `array` | `[]<ItemType>` |
