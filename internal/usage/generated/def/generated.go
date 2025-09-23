package def

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jolfzverb/codegen/internal/usage/generated/api/apimodels"
)

type ExternalObject struct {
	Field1 string `json:"field1,omitempty"`
	Field2 string `json:"field2,omitempty"`
}
type ExternalRef string
type ExternalRef2 struct {
	Subfield1 string `json:"subfield1,omitempty"`
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
func (h *Handler) AddRoutes(r *chi.Mux) {
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
