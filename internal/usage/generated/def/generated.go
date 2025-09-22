package def

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type ExternalRef string
type ExternalRef2 struct {
	Subfield1 string `json:"subfield1,omitempty"`
}
type ExternalObject struct {
	Field1 string `json:"field1,omitempty"`
	Field2 string `json:"field2,omitempty"`
}
type Handler struct {
	validator *validator.Validate
}

func NewHandler(validator *validator.Validate) *Handler {
	return &Handler{validator: validator}
}
func AddRoutes(h *Handler, r *chi.Mux) {
}
