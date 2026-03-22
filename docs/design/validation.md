# Two-Layer Validation System

This is the most distinctive feature of validgo-gen.

## Layer 1: Raw JSON Validation (pre-deserialization)

Generated `Validate<Type>JSON(data json.RawMessage) error` functions check:

1. **Required fields exist** — unmarshals to `map[string]json.RawMessage`, checks key presence
2. **Non-nullable fields aren't null** — calls `containsNull()` to detect JSON `null` literals
3. **Nested objects validated recursively** — calls `Validate<NestedType>JSON()` on nested raw messages
4. **Array items validated** — iterates array elements, validates each item

```go
func ValidateCreateRequestBodyJSON(data json.RawMessage) error {
    mapData := make(map[string]json.RawMessage)
    if err := json.Unmarshal(data, &mapData); err != nil {
        return errors.Wrap(err, "CreateRequestBody")
    }

    // Check required field "name" exists and is not null
    if _, ok := mapData["name"]; !ok {
        return errors.New("CreateRequestBody: field 'name' is required")
    }
    if containsNull(mapData["name"]) {
        return errors.New("CreateRequestBody: field 'name' must not be null")
    }

    // Recursively validate nested object
    if metadataData, ok := mapData["metadata"]; ok && !containsNull(metadataData) {
        if err := ValidateCreateRequestBodyMetadataJSON(metadataData); err != nil {
            return errors.Wrap(err, "CreateRequestBody.metadata")
        }
    }

    return nil
}
```

**Why this matters:** Standard Go JSON unmarshaling silently accepts missing required fields (they get zero values) and null values for non-pointer types. This layer catches those issues *before* the data hits your structs.

## Layer 2: Struct Tag Validation (post-deserialization)

After JSON is unmarshaled into the typed struct, `go-playground/validator/v10` validates via struct tags:

```go
// In parseCreateRequestBody:
var body apimodels.CreateRequestBody
json.Unmarshal(data, &body)
err = h.validator.Struct(body)  // validates min, max, oneof, email, ip, etc.
```

## Validation flow in generated parse methods

```
HTTP Request
    │
    ▼
Read body → json.RawMessage
    │
    ▼
ValidateCreateRequestBodyJSON(raw)     ← Layer 1: required/null/structure
    │ error? → 400 Bad Request
    ▼
json.Unmarshal(raw, &body)
    │ error? → 400 Bad Request
    ▼
validator.Struct(body)                  ← Layer 2: constraints (min/max/oneof/etc)
    │ error? → 400 Bad Request
    ▼
Populate CreateRequest.Body
```
