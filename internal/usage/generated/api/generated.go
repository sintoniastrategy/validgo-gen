package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type NewResourseResponse struct {
	Date2        string `json:"date2,omitempty"`
	DecimalField string `json:"decimal-field,omitempty"`
	Description  string `json:"description,omitempty"`
	EnumVal      string `json:"enum-val,omitempty"`
	Name         string `json:"name"`
	Param        string `json:"param"`
	Count        string `json:"count"`
	Date         string `json:"date,omitempty"`
}
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
	validator *validator.Validate
}

func NewHandler(validator *validator.Validate) *Handler {
	return &Handler{validator: validator}
}
func Create(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
func AddRoutes(h *Handler, r *chi.Mux) {
	r.Post("/path/to/{param}/resours{suffix}", http.HandlerFunc(Create))
}
