package api

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/jolfzverb/codegen/internal/usage/generated/api/apimodels"
	"net/http"
)

type ComplexObjectForDive struct {
	Arrayobjectsrequired []string `json:"array_objects_required"`
	Arraystringsoptional []string `json:"array_strings_optional,omitempty"`
	Arraystringsrequired []string `json:"array_strings_required"`
	Arraysofarrays       []string `json:"arrays_of_arrays,omitempty"`
	Objectfieldoptional  string   `json:"object_field_optional,omitempty"`
	Objectfieldrequired  string   `json:"object_field_required"`
	Arrayobjectsoptional []string `json:"array_objects_optional,omitempty"`
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
func Create(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
func (h *Handler) AddRoutes(r *chi.Mux) {
	r.Post("/path/to/{param}/resours{suffix}", http.HandlerFunc(Create))
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
