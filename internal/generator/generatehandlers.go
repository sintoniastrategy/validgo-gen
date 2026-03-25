package generator

import (
	"mime"
	"sort"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-faster/errors"
)

const applicationJSONCT = "application/json"

// parseMediaType extracts the base media type from a content type string,
// stripping parameters (e.g. "charset=utf-8") and handling edge cases like
// comma-separated values. For example:
//   - "application/json; charset=utf-8" returns "application/json"
//   - "application/json; charset=utf-8, application/json" returns "application/json"
func parseMediaType(rawContentType string) string {
	mediaType, _, _ := mime.ParseMediaType(rawContentType)
	if mediaType == "" {
		return rawContentType
	}
	return mediaType
}

func (g *Generator) AddInterface(baseName string) {
	interfaceName := baseName + "Handler"
	methodName := "Handle" + baseName
	requestName := baseName + "Request"
	responseName := baseName + "Response"
	g.AddHandlersInterface(interfaceName, methodName, requestName, responseName)
}

func (g *Generator) AddDependencyToHandler(baseName string) {
	g.AddDependencyToHandlers(baseName)
}

func (g *Generator) AddRoute(baseName string, method string, pathName string) {
	g.AddRouteToRouter(baseName, method, pathName)
}

func (g *Generator) AddContentTypeToHandler(baseName string, rawContentType string) {
	if g.GetHandler(baseName) == nil {
		g.CreateHandler(baseName)
	}
	g.AddContentTypeHandler(baseName, rawContentType)
}

func (g *Generator) AddHandleOperationMethod(baseName string) {
	g.AddHandleOperationMethodHandlers(baseName)
}

func (g *Generator) AddResponseCodeModels(baseName string, code string, response *openapi3.ResponseRef) error {
	const op = "generator.AddResponseCodeModels"
	if len(response.Value.Content) > 1 {
		return errors.New("multiple response content types are not supported")
	}
	model := SchemaStruct{
		Name:   baseName + "Response" + code,
		Fields: []SchemaField{},
	}
	for _, content := range response.Value.Content {
		if content.Schema != nil {
			if content.Schema.Ref == "" {
				err := g.ProcessSchema(baseName+"Response"+code+"Body", content.Schema)
				if err != nil {
					return errors.Wrap(err, op)
				}
			}
			typeName := baseName + "Response" + code + "Body"
			if content.Schema.Ref != "" {
				var importPath string
				typeName, importPath = g.ParseRefTypeName(content.Schema.Ref)
				if importPath != "" {
					g.AddSchemasImport(importPath)
				}
			}
			model.Fields = append(model.Fields, SchemaField{
				Name:        "Body",
				Type:        typeName,
				TagJSON:     []string{},
				TagValidate: []string{},
				Required:    true,
			})
		}
	}
	if len(response.Value.Headers) > 0 {
		err := g.AddHeadersModel(baseName+"Response"+code, response.Value.Headers)
		if err != nil {
			return errors.Wrap(err, op)
		}
		model.Fields = append(model.Fields, SchemaField{
			Name:     "Headers",
			Type:     baseName + "Response" + code + "Headers",
			Required: true,
		})
	}
	g.AddSchema(model)
	err := g.AddCreateResponseModel(baseName, code, response)
	if err != nil {
		return errors.Wrapf(err, op)
	}

	return nil
}

func (g *Generator) AddResponseModel(baseName string, responseCodes []string) {
	model := SchemaStruct{
		Name: baseName + "Response",
		Fields: []SchemaField{
			{
				Name:     "StatusCode",
				Type:     "int",
				Required: true,
			},
		},
	}
	for _, code := range responseCodes {
		field := SchemaField{
			Name: "Response" + code,
			Type: baseName + "Response" + code,
		}
		model.Fields = append(model.Fields, field)
	}
	g.AddSchema(model)
}

func (g *Generator) AddWriteResponseMethod(baseName string, operation *openapi3.Operation) error {
	const op = "generator.AddWriteResponseMethod"
	var err error
	codes := make([]string, 0, len(operation.Responses.Map()))
	keys := make([]string, 0, len(operation.Responses.Map()))
	for code := range operation.Responses.Map() {
		keys = append(keys, code)
	}
	sort.Strings(keys)
	for _, code := range keys {
		response := operation.Responses.Value(code)
		err = g.AddResponseCodeModels(baseName, code, response)
		if err != nil {
			return errors.Wrap(err, op)
		}
		err = g.AddWriteResponseCode(baseName, code, response)
		if err != nil {
			return errors.Wrap(err, op)
		}
		if len(response.Value.Headers) > 0 {
			err = g.AddWriteHeadersForResponseCode(baseName, code, response)
			if err != nil {
				return errors.Wrap(err, op)
			}
		}
		codes = append(codes, code)

	}
	err = g.AddWriteResponseMethodHandlers(baseName, codes, operation)
	if err != nil {
		return errors.Wrap(err, op)
	}
	g.AddResponseModel(baseName, keys)

	return nil
}

func (g *Generator) GetOperationParamsByType(operation *openapi3.Operation, paramIn string) openapi3.Parameters {
	var result openapi3.Parameters
	for _, p := range operation.Parameters {
		if p.Value.In == paramIn {
			result = append(result, p)
		}
	}

	return result
}

func (g *Generator) AddParseParamsMethods(baseName string, contentType string, operation *openapi3.Operation) error {
	const op = "generator.AddParseParamsMethods"
	var err error

	pathParams := g.GetOperationParamsByType(operation, openapi3.ParameterInPath)
	if len(pathParams) > 0 {
		err = g.AddParamsModel(baseName, "PathParams", pathParams)
		if err != nil {
			return errors.Wrap(err, op)
		}
		err = g.AddParsePathParamsMethod(baseName, pathParams)
		if err != nil {
			return errors.Wrap(err, op)
		}
	}
	queryParams := g.GetOperationParamsByType(operation, openapi3.ParameterInQuery)
	if len(queryParams) > 0 {
		err = g.AddParamsModel(baseName, "QueryParams", queryParams)
		if err != nil {
			return errors.Wrap(err, op)
		}
		err = g.AddParseQueryParamsMethod(baseName, queryParams)
		if err != nil {
			return errors.Wrap(err, op)
		}
	}
	headerParams := g.GetOperationParamsByType(operation, openapi3.ParameterInHeader)
	if len(headerParams) > 0 {
		err = g.AddParamsModel(baseName, "Headers", headerParams)
		if err != nil {
			return errors.Wrap(err, op)
		}
		err = g.AddParseHeadersMethod(baseName, headerParams)
		if err != nil {
			return errors.Wrap(err, op)
		}
	}
	cookieParams := g.GetOperationParamsByType(operation, openapi3.ParameterInCookie)
	if len(cookieParams) > 0 {
		err = g.AddParamsModel(baseName, "Cookies", cookieParams)
		if err != nil {
			return errors.Wrap(err, op)
		}
		err = g.AddParseCookiesMethod(baseName, cookieParams)
		if err != nil {
			return errors.Wrap(err, op)
		}
	}
	if operation.RequestBody != nil && operation.RequestBody.Value != nil {
		content, ok := operation.RequestBody.Value.Content[contentType]
		if ok && content.Schema != nil {
			if content.Schema.Ref == "" {
				err = g.ProcessSchema(baseName+"RequestBody", content.Schema)
				if err != nil {
					return errors.Wrap(err, op)
				}
			}
			err = g.AddParseRequestBodyMethod(baseName, contentType, operation.RequestBody)
			if err != nil {
				return errors.Wrap(err, op)
			}
		}
	}
	g.AddParseRequestMethod(baseName, contentType,
		pathParams, queryParams, headerParams, cookieParams, operation.RequestBody,
	)
	g.GenerateRequestModel(baseName, contentType,
		pathParams, queryParams, headerParams, cookieParams, operation.RequestBody,
	)

	return nil
}

func (g *Generator) ProcessApplicationJSONOperation(pathName string, method string, contentType string,
	operation *openapi3.Operation,
) error {
	const op = "generator.ProcessApplicationJsonOperation"
	if contentType == "" {
		contentType = applicationJSONCT
	}
	handlerBaseName := FormatGoLikeIdentifier(method) + FormatGoLikeIdentifier(pathName)
	if operation.OperationID != "" {
		handlerBaseName = FormatGoLikeIdentifier(operation.OperationID)
	}

	g.AddInterface(handlerBaseName)
	g.AddDependencyToHandler(handlerBaseName)
	g.AddRoute(handlerBaseName, method, pathName)
	err := g.AddParseParamsMethods(handlerBaseName, contentType, operation)
	if err != nil {
		return errors.Wrap(err, op)
	}
	err = g.AddWriteResponseMethod(handlerBaseName, operation)
	if err != nil {
		return errors.Wrap(err, op)
	}
	g.AddHandleOperationMethod(handlerBaseName)
	if operation.RequestBody != nil {
		g.AddContentTypeToHandler(handlerBaseName, contentType)
	} else {
		g.CreateDirectHandler(handlerBaseName)
	}

	return nil
}

func (g *Generator) ProcessOperation(pathName string, method string, operation *openapi3.Operation) error {
	const op = "generator.ProcessOperation"

	if operation.RequestBody != nil {
		contentKeys := make([]string, 0, len(operation.RequestBody.Value.Content))
		for contentType := range operation.RequestBody.Value.Content {
			contentKeys = append(contentKeys, contentType)
		}
		sort.Strings(contentKeys)
		for _, contentType := range contentKeys {
			switch parseMediaType(contentType) {
			case applicationJSONCT:
				err := g.ProcessApplicationJSONOperation(pathName, method, contentType, operation)
				if err != nil {
					return errors.Wrap(err, op)
				}
			default:
				return errors.New("unsupported content type")
			}
		}
	} else {
		err := g.ProcessApplicationJSONOperation(pathName, method, "", operation)
		if err != nil {
			return errors.Wrap(err, op)
		}
	}

	return nil
}

func (g *Generator) ProcessPaths(paths *openapi3.Paths) error {
	const op = "generator.ProcessPaths"
	g.AddHandlersImport(g.ModelsImportPath)
	g.AddHandlersImport("context")
	g.AddHandlersImport("net/http")
	for _, pathName := range paths.InMatchingOrder() {
		pathItem := paths.Value(pathName)
		if pathItem.Get != nil {
			if pathItem.Get.RequestBody != nil {
				return errors.New("GET method should not have request body")
			}
			err := g.ProcessOperation(pathName, "Get", pathItem.Get)
			if err != nil {
				return errors.Wrap(err, op)
			}
		}
		if pathItem.Post != nil {
			err := g.ProcessOperation(pathName, "Post", pathItem.Post)
			if err != nil {
				return errors.Wrap(err, op)
			}
		}
		if pathItem.Delete != nil {
			if !g.Opts.AllowDeleteWithBody && pathItem.Delete.RequestBody != nil {
				return errors.New("DELETE method should not have request body")
			}
			err := g.ProcessOperation(pathName, "Delete", pathItem.Delete)
			if err != nil {
				return errors.Wrap(err, op)
			}
		}
		if pathItem.Put != nil {
			err := g.ProcessOperation(pathName, "Put", pathItem.Put)
			if err != nil {
				return errors.Wrap(err, op)
			}
		}
		if pathItem.Patch != nil {
			err := g.ProcessOperation(pathName, "Patch", pathItem.Patch)
			if err != nil {
				return errors.Wrap(err, op)
			}
		}
	}

	return nil
}
