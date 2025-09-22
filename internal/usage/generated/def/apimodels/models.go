package apimodels

import "time"

type RequestBody struct {
	Name                string
	Description         string
	Date                *time.Time
	CodeForResponse     *int
	EnumVal             string
	DecimalField        string
	FieldToValidateDive *ComplexObjectForDive
}
type RequestHeaders struct {
	IdempotencyKey string
	OptionalHeader *time.Time
}
type RequestQuery struct {
	Count string
}
type RequestPath struct {
	Param string
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
	Name         string
	Param        string
	Count        string
	Date         *time.Time
	Date2        *time.Time
	DecimalField string
	Description  string
	EnumVal      string
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
