package complexapi

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type User struct {
	Age   int    `json:"age,omitempty" validate:"min=0.000000,max=120.000000"`
	Email string `json:"email" validate:"email"`
	Id    int    `json:"id"`
	Name  string `json:"name"`
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
func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
func GetUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
func AddRoutes(h *Handler, r *chi.Mux) {
	r.Get("/users", http.HandlerFunc(GetUsers))
	r.Post("/users", http.HandlerFunc(CreateUser))
	r.Delete("/users/{id}", http.HandlerFunc(DeleteUser))
	r.Get("/users/{id}", http.HandlerFunc(GetUser))
	r.Put("/users/{id}", http.HandlerFunc(UpdateUser))
}
