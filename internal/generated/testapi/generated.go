package testapi

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type User struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `validate:"email" json:"email"`
}
type Handler struct {
	validator *validator.Validate
}

func NewHandler(validator *validator.Validate) *Handler {
	return &Handler{validator: validator}
}
func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
func AddRoutes(h *Handler, r *chi.Mux) {
	r.Get("/users", http.HandlerFunc(GetUsers))
}
