package api

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jolfzverb/codegen/internal/usage/generated/api/apimodels"
	"net/http"
	"time"
)

type ComplexObjectForDive struct {
	Arraysofarrays       []string `json:"arrays_of_arrays,omitempty"`
	Objectfieldoptional  string   `json:"object_field_optional,omitempty"`
	Objectfieldrequired  string   `json:"object_field_required"`
	Arrayobjectsoptional []string `json:"array_objects_optional,omitempty"`
	Arrayobjectsrequired []string `json:"array_objects_required"`
	Arraystringsoptional []string `json:"array_strings_optional,omitempty"`
	Arraystringsrequired []string `json:"array_strings_required"`
}
type Handler struct {
	validator *validator.Validate
	create    CreateHandler
}
type CreateHandler interface {
	HandleCreate(ctx context.Context, r apimodels.CreateRequest) (*apimodels.CreateResponse, error)
}

func NewHandler(create CreateHandler) *Handler {
	return &Handler{validator: validator.New(validator.WithRequiredStructEnabled()), create: create}
}
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var param string = chi.URLParam(r, "param")
	var count string = r.URL.Query().Get("count")
	var idempotencyKey string = r.Header.Get("Idempotency-Key")
	var optionalHeader *time.Time
	if r.Header.Get("Optional-Header") != "" {
		var optionalHeaderStr string = r.Header.Get("Optional-Header")
		var parsedTime time.Time
		var err error
		parsedTime, err = time.Parse("2006-01-02T15:04:05Z07:00", optionalHeaderStr)
		if err == nil {
			optionalHeader = &parsedTime
		}
	}
	var requiredCookieParam *http.Cookie
	var cookieParam *string
	var err error
	requiredCookieParam, err = r.Cookie("required-cookie-param")
	if err != nil {
		http.Error(w, "{\"error\":\"required-cookie-param cookie is required\"}", http.StatusBadRequest)
		return
	}
	var requiredCookieParamValue string = requiredCookieParam.Value
	var cookie *http.Cookie
	cookie, err = r.Cookie("cookie-param")
	if err != nil && err != http.ErrNoCookie {
		http.Error(w, "{\"error\":\"Invalid cookie\"}", http.StatusBadRequest)
		return
	}
	if err == nil {
		cookieParam = &cookie.Value
	}
	var body apimodels.RequestBody
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "{\"error\":\"Invalid JSON\"}", http.StatusBadRequest)
		return
	}
	if body.Name == "" {
		http.Error(w, "{\"error\":\"name is required\"}", http.StatusBadRequest)
		return
	}
	var req apimodels.CreateRequest = apimodels.CreateRequest{Body: body, Headers: apimodels.RequestHeaders{IdempotencyKey: idempotencyKey, OptionalHeader: optionalHeader}, Query: apimodels.RequestQuery{Count: count}, Path: apimodels.RequestPath{Param: param}, Cookies: apimodels.RequestCookies{RequiredCookieParam: requiredCookieParamValue, CookieParam: cookieParam}}
	var response *apimodels.CreateResponse
	response, err = h.create.HandleCreate(r.Context(), req)
	if err != nil {
		http.Error(w, "{\"error\":\"Internal Server Error\"}", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if response.Response200 != nil {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(response.Response200.Data)
	}
	if response.Response400 != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(response.Response400)
	}
	if response.Response404 != nil {
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(response.Response404)
	}
}
func (h *Handler) AddRoutes(r *chi.Mux) {
	r.Post("/path/to/{param}/resourse", h.Create)
}
func parseTime(timeStr string) *time.Time {
	if timeStr == "" {
		return nil
	}
	var t time.Time
	var err error
	t, err = time.Parse("2006-01-02T15:04:05Z07:00", timeStr)
	if err != nil {
		return nil
	}
	return &t
}
func Create200Response(data apimodels.NewResourseResponse, headers apimodels.CreateResponse200Headers) *apimodels.CreateResponse {
	return &apimodels.CreateResponse{StatusCode: 200, Response200: &apimodels.Response200Data{Data: data, Headers: headers}}
}
func Create400Response() *apimodels.CreateResponse {
	return &apimodels.CreateResponse{StatusCode: 400, Response400: &apimodels.Response400Data{Error: "Bad Request"}}
}
func Create404Response() *apimodels.CreateResponse {
	return &apimodels.CreateResponse{StatusCode: 404, Response404: &apimodels.Response404Data{Error: "Not Found"}}
}
