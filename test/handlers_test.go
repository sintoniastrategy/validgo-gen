package test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/sintoniastrategy/validgo-gen/internal/usage/generated/api"
	"github.com/sintoniastrategy/validgo-gen/internal/usage/generated/api/apimodels"
	"github.com/stretchr/testify/assert"
)

type mockHandler struct{}

func (m *mockHandler) HandleCreate(ctx context.Context, r apimodels.CreateRequest) (*apimodels.CreateResponse, error) {
	if r.Body.CodeForResponse != nil {
		switch *r.Body.CodeForResponse {
		case 400:
			return api.Create400Response(), nil
		case 404:
			return api.Create404Response(), nil
		}
	}
	var date *time.Time
	if r.Body.Date != nil {
		date = new(time.Time)
		*date = r.Body.Date.UTC()
	}
	var date2 *time.Time
	if r.Headers.OptionalHeader != nil {
		date2 = new(time.Time)
		*date2 = r.Headers.OptionalHeader.UTC()
	}
	return api.Create200Response(
		apimodels.NewResourseResponse{
			Count:        r.Query.Count,
			Description:  r.Body.Description,
			Name:         r.Body.Name,
			Param:        r.Path.Param,
			Date:         date,
			Date2:        date2,
			EnumVal:      r.Body.EnumVal,
			DecimalField: r.Body.DecimalField,
		},
		apimodels.CreateResponse200Headers{
			IdempotencyKey: &r.Headers.IdempotencyKey,
		},
	), nil
}

func TestHandler(t *testing.T) {
	router := chi.NewRouter()
	handler := api.NewHandler(
		&mockHandler{},
	)
	handler.AddRoutes(router)

	// Create a test server using the chi router
	server := httptest.NewServer(router)
	defer server.Close()

	t.Run("200 Success", func(t *testing.T) {
		requestBody := `{"name": "value", "description": "descr", "date": "2023-10-01T00:00:00+03:00", "code_for_response": 200, "enum-val": "value1", "decimal-field": "13.42"}`
		request, err := http.NewRequest(http.MethodPost, server.URL+"/path/to/param/resourse?count=3", bytes.NewBufferString(requestBody))
		assert.NoError(t, err)
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Idempotency-Key", "unique-idempotency-key")
		request.Header.Set("Optional-Header", "2023-10-01T00:00:00+03:00")
		request.Header.Set("Cookie", "required-cookie-param=required-value")
		resp, err := http.DefaultClient.Do(request)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

		defer resp.Body.Close()
		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		assert.NoError(t, err)
		assert.Equal(t, "unique-idempotency-key", resp.Header.Get("Idempotency-Key"))
		assert.Equal(t, "3", responseBody["count"])
		assert.Equal(t, "descr", responseBody["description"])
		assert.Equal(t, "value", responseBody["name"])
		assert.Equal(t, "2023-09-30T21:00:00Z", responseBody["date"])
		assert.Equal(t, "2023-09-30T21:00:00Z", responseBody["date2"])
		assert.Equal(t, "value1", responseBody["enum-val"])
		assert.Equal(t, "13.42", responseBody["decimal-field"])
	})
	t.Run("404", func(t *testing.T) {
		requestBody := `{"name": "value", "description": "descr", "code_for_response": 404}`
		request, err := http.NewRequest(http.MethodPost, server.URL+"/path/to/param/resourse?count=3", bytes.NewBufferString(requestBody))
		assert.NoError(t, err)
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Idempotency-Key", "unique-idempotency-key")
		request.Header.Set("Cookie", "required-cookie-param=required-value")
		resp, err := http.DefaultClient.Do(request)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
	t.Run("400 No name", func(t *testing.T) {
		requestBody := `{}`
		request, err := http.NewRequest(http.MethodPost, server.URL+"/path/to/param/resourse?count=3", bytes.NewBufferString(requestBody))
		assert.NoError(t, err)
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Idempotency-Key", "unique-idempotency-key")
		request.Header.Set("Cookie", "required-cookie-param=required-value")
		resp, err := http.DefaultClient.Do(request)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		defer resp.Body.Close()
		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		assert.NoError(t, err)
	})
	t.Run("400 number enum", func(t *testing.T) {
		requestBody := `{"name": "value", "enum-int": 15}`
		request, err := http.NewRequest(http.MethodPost, server.URL+"/path/to/param/resourse?count=3", bytes.NewBufferString(requestBody))
		assert.NoError(t, err)
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Idempotency-Key", "unique-idempotency-key")
		request.Header.Set("Cookie", "required-cookie-param=required-value")
		resp, err := http.DefaultClient.Do(request)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
	t.Run("400 required cookie", func(t *testing.T) {
		requestBody := `{"name": "value"}`
		request, err := http.NewRequest(http.MethodPost, server.URL+"/path/to/param/resourse?count=3", bytes.NewBufferString(requestBody))
		assert.NoError(t, err)
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Idempotency-Key", "unique-idempotency-key")
		resp, err := http.DefaultClient.Do(request)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
	t.Run("400 cookie validation 1", func(t *testing.T) {
		requestBody := `{"name": "value"}`
		request, err := http.NewRequest(http.MethodPost, server.URL+"/path/to/param/resourse?count=3", bytes.NewBufferString(requestBody))
		assert.NoError(t, err)
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Idempotency-Key", "unique-idempotency-key")
		request.Header.Set("Cookie", "required-cookie-param=required-value; cookie-param=short")
		resp, err := http.DefaultClient.Do(request)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
	t.Run("400 cookie validation 2", func(t *testing.T) {
		requestBody := `{"name": "value"}`
		request, err := http.NewRequest(http.MethodPost, server.URL+"/path/to/param/resourse?count=3", bytes.NewBufferString(requestBody))
		assert.NoError(t, err)
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Idempotency-Key", "unique-idempotency-key")
		request.Header.Set("Cookie", "required-cookie-param=required-value-too-long")
		resp, err := http.DefaultClient.Do(request)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
	t.Run("400 invalid suffix", func(t *testing.T) {
		requestBody := `{"name": "value"}`
		request, err := http.NewRequest(http.MethodPost, server.URL+"/path/to/param/resourseee?count=3", bytes.NewBufferString(requestBody))
		assert.NoError(t, err)
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Idempotency-Key", "unique-idempotency-key")
		request.Header.Set("Cookie", "required-cookie-param=required-value")
		resp, err := http.DefaultClient.Do(request)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		defer resp.Body.Close()
		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		assert.NoError(t, err)
	})
	t.Run("404 no suffix", func(t *testing.T) {
		requestBody := `{"name": "value"}`
		request, err := http.NewRequest(http.MethodPost, server.URL+"/path/to/param/resours?count=3", bytes.NewBufferString(requestBody))
		assert.NoError(t, err)
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Idempotency-Key", "unique-idempotency-key")
		request.Header.Set("Cookie", "required-cookie-param=required-value")
		resp, err := http.DefaultClient.Do(request)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
	t.Run("200 on dive 1", func(t *testing.T) {
		requestBody := `{"name": "value", "description": "descr", "date": "2023-10-01T00:00:00+03:00", "code_for_response": 200, "enum-val": "value1", "decimal-field": "13.42",
		"field_to_validate_dive": {
		  "object_field_required": {
		  	"field1": "minimum:5"
		  },
		  "array_objects_required": [{"field1":"minumum:5"}],
		  "array_strings_required": ["minimum:5"]
		}}`
		request, err := http.NewRequest(http.MethodPost, server.URL+"/path/to/param/resourse?count=3", bytes.NewBufferString(requestBody))
		assert.NoError(t, err)
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Idempotency-Key", "unique-idempotency-key")
		request.Header.Set("Cookie", "required-cookie-param=required-value")
		resp, err := http.DefaultClient.Do(request)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		defer resp.Body.Close()
		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		assert.NoError(t, err)
	})
	t.Run("200 on dive with optional fields", func(t *testing.T) {
		requestBody := `{"name": "value", "description": "descr", "date": "2023-10-01T00:00:00+03:00", "code_for_response": 200, "enum-val": "value1", "decimal-field": "13.42",
		"field_to_validate_dive": {
			"object_field_required": {
				"field1": "minimum:5"
			},
			"object_field_optional": {
				"field1": "minimum:5"
			},
			"array_objects_required": [{"field1":"minumum:5"}],
			"array_objects_optional": [{"field1":"minumum:5"}],
			"array_strings_required": ["minimum:5"],
			"array_strings_optional": ["minimum:5"]
		}}`
		request, err := http.NewRequest(http.MethodPost, server.URL+"/path/to/param/resourse?count=3", bytes.NewBufferString(requestBody))
		assert.NoError(t, err)
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Idempotency-Key", "unique-idempotency-key")
		request.Header.Set("Cookie", "required-cookie-param=required-value")
		resp, err := http.DefaultClient.Do(request)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		defer resp.Body.Close()
		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		assert.NoError(t, err)
	})

	for _, tc := range []struct {
		name      string
		diveField string
	}{
		{
			name: "400 on dive required fields 1",
			diveField: `"field_to_validate_dive": {
				"object_field_required": {"field1": "min"},
				"array_objects_required": [{"field1":"minumum:5"}],
				"array_strings_required": ["minimum:5"]
			}`,
		},
		{
			name: "400 on dive required fields 2",
			diveField: `"field_to_validate_dive": {
				"object_field_required": {"field1": "minimum:5"},
				"array_objects_required": [{"field1":"min"}],
				"array_strings_required": ["minimum:5"]
			}`,
		},

		{
			name: "400 on dive required fields 3",
			diveField: `"field_to_validate_dive": {
				"object_field_required": {"field1": "minimum:5"},
				"array_objects_required": [{}],
				"array_strings_required": ["minimum:5"]
			}`,
		},
		{
			name: "400 on dive required fields 4",
			diveField: `"field_to_validate_dive": {
				"object_field_required": {"field1": "minimum:5"},
				"array_objects_required": [],
				"array_strings_required": ["minimum:5"]
			}`,
		},
		{
			name: "400 on dive required fields 5",
			diveField: `"field_to_validate_dive": {
				"object_field_required": {"field1": "minimum:5"},
				"array_objects_required": [{"field1":"minumum:5"}],
				"array_strings_required": ["min"]
			}`,
		},
		{
			name: "400 on dive required fields 6",
			diveField: `"field_to_validate_dive": {
				"object_field_required": {"field1": "minimum:5"},
				"array_objects_required": [{"field1":"minumum:5"}],
				"array_strings_required": []
			}`,
		},
		{
			name: "400 on dive optional fields 1",
			diveField: `"field_to_validate_dive": {
				"object_field_required": {"field1": "minimum:5"},
				"object_field_optional": {"field1": "min"},
				"array_objects_required": [{"field1":"minumum:5"}],
				"array_objects_optional": [{"field1":"minumum:5"}],
				"array_strings_required": ["minimum:5"],
				"array_strings_optional": ["minimum:5"]
			}`,
		},
		{
			name: "400 on dive optional fields 2",
			diveField: `"field_to_validate_dive": {
				"object_field_required": {"field1": "minimum:5"},
				"object_field_optional": {"field1": "minimum:5"},
				"array_objects_required": [{"field1":"minumum:5"}],
				"array_objects_optional": [{"field1":"min"}],
				"array_strings_required": ["minimum:5"],
				"array_strings_optional": ["minimum:5"]
			}`,
		},
		{
			name: "400 on dive optional fields 3",
			diveField: `"field_to_validate_dive": {
				"object_field_required": {"field1": "minimum:5"},
				"object_field_optional": {"field1": "minimum:5"},
				"array_objects_required": [{}],
				"array_objects_optional": [{"field1":"min"}],
				"array_strings_required": ["minimum:5"],
				"array_strings_optional": ["minimum:5"]
			}`,
		},
		{
			name: "400 on dive optional fields 4",
			diveField: `"field_to_validate_dive": {
				"object_field_required": {"field1": "minimum:5"},
				"object_field_optional": {"field1": "minimum:5"},
				"array_objects_required": [],
				"array_objects_optional": [{"field1":"min"}],
				"array_strings_required": ["minimum:5"],
				"array_strings_optional": ["minimum:5"]
			}`,
		},
		{
			name: "400 on dive optional fields 5",
			diveField: `"field_to_validate_dive": {
				"object_field_required": {"field1": "minimum:5"},
				"object_field_optional": {"field1": "minimum:5"},
				"array_objects_required": [{"field1":"minumum:5"}],
				"array_objects_optional": [{"field1":"minumum:5"}],
				"array_strings_required": ["minimum:5"],
				"array_strings_optional": ["min"]
			}`,
		},
		{
			name: "400 on dive optional fields 6",
			diveField: `"field_to_validate_dive": {
				"object_field_required": {"field1": "minimum:5"},
				"object_field_optional": {"field1": "minimum:5"},
				"array_objects_required": [{"field1":"minumum:5"}],
				"array_objects_optional": [{"field1":"minumum:5"}],
				"array_strings_required": ["minimum:5"],
				"array_strings_optional": []
			}`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			requestBody := `{"name": "value", "description": "descr", "date": "2023-10-01T00:00:00+03:00", "code_for_response": 200, "enum-val": "value1", "decimal-field": "13.42",` +
				tc.diveField + `}`
			request, err := http.NewRequest(http.MethodPost, server.URL+"/path/to/param/resourse?count=3", bytes.NewBufferString(requestBody))
			assert.NoError(t, err)
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Idempotency-Key", "unique-idempotency-key")
			request.Header.Set("Cookie", "required-cookie-param=required-value")
			resp, err := http.DefaultClient.Do(request)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

			defer resp.Body.Close()
			var responseBody map[string]any
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			assert.NoError(t, err)
		})
	}
}

type mockHandler500 struct{}

func (m *mockHandler500) HandleCreate(ctx context.Context, r apimodels.CreateRequest) (*apimodels.CreateResponse, error) {
	return &apimodels.CreateResponse{
		StatusCode:  http.StatusOK,
		Response200: nil,
		Response400: nil,
		Response404: nil,
	}, nil

}

func Test500(t *testing.T) {
	router := chi.NewRouter()
	handler := api.NewHandler(
		&mockHandler500{},
	)
	handler.AddRoutes(router)

	// Create a test server using the chi router
	server := httptest.NewServer(router)
	defer server.Close()

	t.Run("500 Internal Server Error", func(t *testing.T) {
		requestBody := `{"name": "value", "description": "descr", "date": "2023-10-01T00:00:00+03:00", "code_for_response": 200, "enum-val": "value1"}`
		request, err := http.NewRequest(http.MethodPost, server.URL+"/path/to/param/resourses?count=3", bytes.NewBufferString(requestBody))
		assert.NoError(t, err)
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Idempotency-Key", "unique-idempotency-key")
		request.Header.Set("Optional-Header", "2023-10-01T00:00:00+03:00")
		request.Header.Set("Cookie", "required-cookie-param=required-value")
		resp, err := http.DefaultClient.Do(request)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		defer resp.Body.Close()
		var responseBody map[string]string
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		assert.NoError(t, err)
		assert.Equal(t, "InternalServerError", responseBody["error"])
	})
}
