# Choosing a Go OpenAPI Code Generator: Why validgo-gen Exists

> A practitioner's comparison of Go OpenAPI 3.0 code generators, with code examples showing where each one breaks down — and why we built validgo-gen.

## The Landscape

There are four serious OpenAPI code generators for Go server development:

| Generator | Stars | Age | OpenAPI | Written in |
|---|---|---|---|---|
| **ogen** | ~2K | 5 yrs | 3.0 | Go |
| **oapi-codegen** | ~8K | 7 yrs | 3.0 | Go |
| **go-swagger** | ~10K | 9 yrs | 2.0 only | Go |
| **openapi-generator** | ~26K | 8 yrs | 2.0 + 3.x | Java |

Each makes fundamentally different trade-offs. None of them solve the same problem that validgo-gen targets: **correct request validation with idiomatic Go output and chi router integration**.

This document uses a single OpenAPI spec as a running example to show exactly how each generator behaves — with real code, not marketing claims.

## In this section

- [The Problem](the-problem.md) — what happens when you send a malicious request to each generator
- [ogen](ogen.md) — the performance king: what it does well and where it hurts
- [oapi-codegen](oapi-codegen.md) — the community favorite: clean models, zero validation
- [go-swagger & openapi-generator](others.md) — one stuck on 2.0, one stuck on Java
- [Comparison Matrix](comparison-matrix.md) — side-by-side feature table, when to use what, roadmap

## The Running Example

```yaml
# petstore.yaml
openapi: 3.0.0
info:
  title: Pet Store
  version: 1.0.0
paths:
  /pets/{petId}:
    get:
      operationId: getPet
      parameters:
        - name: petId
          in: path
          required: true
          schema:
            type: string
        - name: include
          in: query
          required: false
          schema:
            type: string
            enum: [owner, vaccinations, both]
      responses:
        '200':
          description: OK
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Pet'
        '404':
          description: Not found
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
    put:
      operationId: updatePet
      parameters:
        - name: petId
          in: path
          required: true
          schema:
            type: string
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/UpdatePetRequest'
      responses:
        '200':
          description: Updated
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Pet'
        '400':
          description: Validation error
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
        '404':
          description: Not found
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
components:
  schemas:
    Pet:
      type: object
      required: [id, name, price]
      properties:
        id:
          type: string
        name:
          type: string
          minLength: 1
          maxLength: 100
        price:
          type: string
          format: decimal
        birthDate:
          type: string
          format: date-time
        tags:
          type: array
          items:
            type: string
            minLength: 1
          minItems: 0
          maxItems: 10
          uniqueItems: true
        address:
          type: object
          required: [city]
          properties:
            city:
              type: string
            zip:
              type: string
              pattern: '^\d{5}$'
    UpdatePetRequest:
      type: object
      required: [name]
      properties:
        name:
          type: string
          minLength: 1
          maxLength: 100
        price:
          type: string
          format: decimal
        tags:
          type: array
          items:
            type: string
            minLength: 1
          maxItems: 10
          uniqueItems: true
        address:
          type: object
          required: [city]
          properties:
            city:
              type: string
            zip:
              type: string
    Error:
      type: object
      required: [code, message]
      properties:
        code:
          type: integer
        message:
          type: string
```
