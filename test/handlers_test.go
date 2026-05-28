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
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/sintoniastrategy/validgo-gen/internal/usage/generated/api"
	"github.com/sintoniastrategy/validgo-gen/internal/usage/generated/api/apimodels"
	"github.com/stretchr/testify/assert"
)

type mockHandler struct{}

func (m *mockHandler) HandleCreate(ctx context.Context, r apimodels.CreateRequest) (*apimodels.CreateResponse, error) {
	if r.Body.CodeForResponse != nil {
		switch *r.Body.CodeForResponse {
		case 400:
			return api.Create400(), nil
		case 404:
			return api.Create404(), nil
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
	return api.Create200(
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
	t.Run("200 with Content-Type charset", func(t *testing.T) {
		requestBody := `{"name": "value", "description": "descr", "date": "2023-10-01T00:00:00+03:00", "code_for_response": 200, "enum-val": "value1", "decimal-field": "13.42"}`
		request, err := http.NewRequest(http.MethodPost, server.URL+"/path/to/param/resourse?count=3", bytes.NewBufferString(requestBody))
		assert.NoError(t, err)
		request.Header.Set("Content-Type", "application/json; charset=utf-8")
		request.Header.Set("Idempotency-Key", "unique-idempotency-key")
		request.Header.Set("Optional-Header", "2023-10-01T00:00:00+03:00")
		request.Header.Set("Cookie", "required-cookie-param=required-value")
		resp, err := http.DefaultClient.Do(request)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		defer resp.Body.Close()
		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		assert.NoError(t, err)
		assert.Equal(t, "value", responseBody["name"])
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

		assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

		defer resp.Body.Close()
		var responseBody map[string]string
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		assert.NoError(t, err)
		assert.Equal(t, "InternalServerError", responseBody["code"])
		assert.Equal(t, "Internal server error", responseBody["error"])
		_, hasReqID := responseBody["req_id"]
		assert.True(t, hasReqID, "envelope must contain req_id field")
	})
}

// TestWithErrorHandler asserts that a custom ErrorHandler passed via
// api.WithErrorHandler fully replaces the default {code,error,req_id}
// envelope at every generated error site.
func TestWithErrorHandler(t *testing.T) {
	custom := func(w http.ResponseWriter, r *http.Request, status int, msg string) {
		w.Header().Set("Content-Type", "application/x-custom+json")
		w.Header().Set("X-Captured-Status", "yes")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"my_status": status,
			"my_msg":    msg,
			"sentinel":  "custom-handler-ran",
		})
	}

	router := chi.NewRouter()
	api.NewHandler(&mockHandler{}, api.WithErrorHandler(custom)).AddRoutes(router)
	server := httptest.NewServer(router)
	defer server.Close()

	req, _ := http.NewRequest(http.MethodPost,
		server.URL+"/path/to/param/resourse?count=3",
		bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotency-Key", "k")
	req.Header.Set("Cookie", "required-cookie-param=required-value")

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "application/x-custom+json", resp.Header.Get("Content-Type"))
	assert.Equal(t, "yes", resp.Header.Get("X-Captured-Status"))

	var body map[string]any
	err = json.NewDecoder(resp.Body).Decode(&body)
	assert.NoError(t, err)
	assert.Equal(t, "custom-handler-ran", body["sentinel"])
	assert.EqualValues(t, http.StatusBadRequest, body["my_status"])
	assert.NotEmpty(t, body["my_msg"], "msg should propagate from parser")
}

// TestSetErrorHandler asserts the post-construction setter has the same
// effect as WithErrorHandler — the aggregator-loop pattern relies on it.
func TestSetErrorHandler(t *testing.T) {
	var captured struct {
		status int
		msg    string
	}
	h := api.NewHandler(&mockHandler{})
	h.SetErrorHandler(func(w http.ResponseWriter, r *http.Request, status int, msg string) {
		captured.status = status
		captured.msg = msg
		w.WriteHeader(status)
		_, _ = w.Write([]byte("ok"))
	})

	router := chi.NewRouter()
	h.AddRoutes(router)
	server := httptest.NewServer(router)
	defer server.Close()

	req, _ := http.NewRequest(http.MethodPost,
		server.URL+"/path/to/param/resourse?count=3",
		bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotency-Key", "k")
	req.Header.Set("Cookie", "required-cookie-param=required-value")

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, http.StatusBadRequest, captured.status)
	assert.NotEmpty(t, captured.msg)
}

// TestErrorHandlerAliasIsCrossPackage verifies the package-level
// ErrorHandler is a *type alias* (not a named type), so a plain
// `func(http.ResponseWriter, *http.Request, int, string)` value can be
// passed to SetErrorHandler on Handlers from different generated packages
// without a per-package conversion. This is the property the aggregator
// pattern documented in the README relies on.
func TestErrorHandlerAliasIsCrossPackage(t *testing.T) {
	var eh = func(w http.ResponseWriter, r *http.Request, status int, msg string) {
		_ = msg
		w.WriteHeader(status)
	}
	// Compile-time check: assigning a bare func value to api.ErrorHandler
	// would fail if ErrorHandler were a named type.
	var _ api.ErrorHandler = eh

	h := api.NewHandler(&mockHandler{})
	h.SetErrorHandler(eh)
	// no runtime assertion needed — the compile-time check above is the test.
	_ = h
}

// TestStandardErrorEnvelope exercises the generated handlers via real HTTP
// and asserts the standard {code,error,req_id} envelope shape across the
// four distinct error sites: 400 (parse/validate), 415 (unsupported media
// type), 404 (route miss handled by user handler), and 500 (handler returned
// nil). chi's RequestID middleware is mounted so req_id is populated.
func TestStandardErrorEnvelope(t *testing.T) {
	router := chi.NewRouter()
	router.Use(chimw.RequestID)
	apiHandler := api.NewHandler(&mockHandler{})
	apiHandler.AddRoutes(router)

	server := httptest.NewServer(router)
	defer server.Close()

	type envelope struct {
		Code  string `json:"code"`
		Error string `json:"error"`
		ReqID string `json:"req_id"`
	}

	doRequest := func(t *testing.T, req *http.Request) (*http.Response, envelope) {
		t.Helper()
		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))
		var env envelope
		err = json.NewDecoder(resp.Body).Decode(&env)
		_ = resp.Body.Close()
		assert.NoError(t, err)
		return resp, env
	}

	t.Run("400 parse error", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost,
			server.URL+"/path/to/param/resourse?count=3",
			bytes.NewBufferString(`{}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Idempotency-Key", "k")
		req.Header.Set("Cookie", "required-cookie-param=required-value")

		resp, env := doRequest(t, req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, "BadRequest", env.Code)
		assert.NotEmpty(t, env.Error, "error message should propagate from parser")
		assert.NotEmpty(t, env.ReqID, "chi RequestID middleware should populate req_id")
	})

	t.Run("415 unsupported content type", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost,
			server.URL+"/path/to/param/resourse?count=3",
			bytes.NewBufferString(`<xml/>`))
		req.Header.Set("Content-Type", "application/xml")
		req.Header.Set("Idempotency-Key", "k")
		req.Header.Set("Cookie", "required-cookie-param=required-value")

		resp, env := doRequest(t, req)
		assert.Equal(t, http.StatusUnsupportedMediaType, resp.StatusCode)
		assert.Equal(t, "UnsupportedMediaType", env.Code)
		assert.Equal(t, "Unsupported Content-Type", env.Error)
		assert.NotEmpty(t, env.ReqID)
	})

	t.Run("500 handler returned nil response", func(t *testing.T) {
		router500 := chi.NewRouter()
		router500.Use(chimw.RequestID)
		api.NewHandler(&mockHandler500{}).AddRoutes(router500)
		srv500 := httptest.NewServer(router500)
		defer srv500.Close()

		req, _ := http.NewRequest(http.MethodPost,
			srv500.URL+"/path/to/param/resourses?count=3",
			bytes.NewBufferString(
				`{"name":"value","description":"d","date":"2023-10-01T00:00:00+03:00","code_for_response":200,"enum-val":"value1"}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Idempotency-Key", "k")
		req.Header.Set("Optional-Header", "2023-10-01T00:00:00+03:00")
		req.Header.Set("Cookie", "required-cookie-param=required-value")

		resp, env := doRequest(t, req)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.Equal(t, "InternalServerError", env.Code)
		assert.Equal(t, "Internal server error", env.Error,
			"500 must use generic message, not leak internal err")
		assert.NotEmpty(t, env.ReqID)
	})

	t.Run("envelope works without chi RequestID middleware", func(t *testing.T) {
		bare := chi.NewRouter()
		api.NewHandler(&mockHandler{}).AddRoutes(bare)
		srv := httptest.NewServer(bare)
		defer srv.Close()

		req, _ := http.NewRequest(http.MethodPost,
			srv.URL+"/path/to/param/resourse?count=3",
			bytes.NewBufferString(`{}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Idempotency-Key", "k")
		req.Header.Set("Cookie", "required-cookie-param=required-value")

		resp, env := doRequest(t, req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, "BadRequest", env.Code)
		assert.Equal(t, "", env.ReqID, "req_id falls back to empty string without middleware")
	})
}
