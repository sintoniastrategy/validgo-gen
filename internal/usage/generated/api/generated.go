package api

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/jolfzverb/codegen/internal/usage/generated/api/apimodels"
	"net/http"
	"time"
)

type ComplexObjectForDive struct {
	Arraystringsoptional []string `json:"array_strings_optional,omitempty"`
	Arraystringsrequired []string `json:"array_strings_required"`
	Arraysofarrays       []string `json:"arrays_of_arrays,omitempty"`
	Objectfieldoptional  string   `json:"object_field_optional,omitempty"`
	Objectfieldrequired  string   `json:"object_field_required"`
	Arrayobjectsoptional []string `json:"array_objects_optional,omitempty"`
	Arrayobjectsrequired []string `json:"array_objects_required"`
}
type Handler struct {
	handler HandlerInterface
}
type HandlerInterface interface {
	HandleCreate(ctx context.Context, r apimodels.CreateRequest) (*apimodels.CreateResponse, error)
}

func NewHandler(handler HandlerInterface) *Handler {
	return &Handler{handler: handler}
}
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var param string = chi.URLParam(r, "param")
	var count string = r.URL.Query().Get("count")
	var idempotencyKey string = r.Header.Get("Idempotency-Key")
	var body apimodels.RequestBody
	var err error
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)()
		return
	}
	var req apimodels.CreateRequest = apimodels.CreateRequest{Body: body, Headers: apimodels.RequestHeaders{IdempotencyKey: idempotencyKey, OptionalHeader: nil}, Query: apimodels.RequestQuery{Count: count}, Path: apimodels.RequestPath{Param: param}}
	var response *apimodels.CreateResponse
	response, err = h.handler.HandleCreate(r.Context(), req)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)()
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
	r.Post("/path/to/{param}/resours{suffix}", http.HandlerFunc(h.Create))
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
