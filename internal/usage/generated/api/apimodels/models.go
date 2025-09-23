package apimodels

import "time"

type RequestBody struct {
	Name                string                `json:"name" validate:"required"`
	Description         string                `json:"description" validate:"omitempty,min=1,max=10"`
	Date                *time.Time            `json:"date,omitempty" validate:"omitempty"`
	CodeForResponse     *int                  `json:"code_for_response,omitempty" validate:"omitempty,min=100,max=999"`
	EnumVal             string                `json:"enum-val" validate:"omitempty,oneof=value1 value2 value3"`
	DecimalField        string                `json:"decimal-field" validate:"omitempty"`
	FieldToValidateDive *ComplexObjectForDive `json:"field_to_validate_dive,omitempty" validate:"omitempty"`
}
type RequestHeaders struct {
	IdempotencyKey string     `json:"Idempotency-Key" validate:"required,min=1,max=100"`
	OptionalHeader *time.Time `json:"Optional-Header,omitempty" validate:"omitempty"`
}
type RequestQuery struct {
	Count string `json:"count" validate:"required"`
}
type RequestPath struct {
	Param string `json:"param" validate:"required"`
}
type RequestCookies struct {
	CookieParam         *string `json:"cookie-param,omitempty" validate:"omitempty,min=10,max=15"`
	RequiredCookieParam string  `json:"required-cookie-param" validate:"required,min=10,max=15"`
}
type Response200Data struct {
	Data    NewResourseResponse
	Headers CreateResponse200Headers
}
type Response400Data struct {
	Error string
}
type Response404Data struct {
	Error string
}
type CreateResponse200Headers struct {
	IdempotencyKey *string
}
type NewResourseResponse struct {
	Name         string     `json:"name"`
	Param        string     `json:"param"`
	Count        string     `json:"count"`
	Date         *time.Time `json:"date,omitempty"`
	Date2        *time.Time `json:"date2,omitempty"`
	DecimalField string     `json:"decimal-field"`
	Description  string     `json:"description"`
	EnumVal      string     `json:"enum-val"`
}
type ComplexObjectForDive struct {
	ObjectFieldRequired  string
	ArrayObjectsOptional []string
	ArrayObjectsRequired []string
	ArrayStringsOptional []string
	ArrayStringsRequired []string
	ArraysOfArrays       []string
	ObjectFieldOptional  string
}
type CreateRequest struct {
	Body    RequestBody
	Headers RequestHeaders
	Query   RequestQuery
	Path    RequestPath
	Cookies RequestCookies
}
type CreateResponse struct {
	StatusCode  int
	Response200 *Response200Data
	Response400 *Response400Data
	Response404 *Response404Data
}

func Create200Response(data NewResourseResponse, headers CreateResponse200Headers) *CreateResponse {
	return &CreateResponse{StatusCode: 200, Response200: &Response200Data{Data: data, Headers: headers}}
}
func Create400Response() *CreateResponse {
	return &CreateResponse{StatusCode: 400, Response400: &Response400Data{Error: "Bad Request"}}
}
func Create404Response() *CreateResponse {
	return &CreateResponse{StatusCode: 404, Response404: &Response404Data{Error: "Not Found"}}
}
