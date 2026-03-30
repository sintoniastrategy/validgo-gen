package generator

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

func GetSchemaValidators(schema *openapi3.SchemaRef) []string {
	var validateTags []string
	switch {
	case schema.Value.Type.Permits(openapi3.TypeString):
		if schema.Value.MinLength > 0 {
			validateTags = append(validateTags, "min="+strconv.FormatUint(schema.Value.MinLength, 10))
		}
		if schema.Value.MaxLength != nil {
			validateTags = append(validateTags, "max="+strconv.FormatUint(*schema.Value.MaxLength, 10))
		}
		if schema.Value.Pattern != "" {
			slog.Warn("pattern validator is not supported", slog.String("pattern", schema.Value.Pattern))
		}
		if len(schema.Value.Enum) > 0 {
			enumStrValues := make([]string, 0, len(schema.Value.Enum))
			for _, enumValue := range schema.Value.Enum {
				var enumStrValue string
				if strValue, ok := enumValue.(string); ok {
					enumStrValue = strValue
				} else {
					slog.Warn("enum value is not a string", slog.Any("value", enumValue))
					enumStrValue = fmt.Sprintf("%v", enumValue)
				}
				if enumStrValue == "" || strings.Contains(enumStrValue, " ") {
					enumStrValue = "'" + enumStrValue + "'"
				}
				enumStrValues = append(enumStrValues, enumStrValue)
			}
			joinedEnum := strings.Join(enumStrValues, " ")
			validateTags = append(validateTags, "oneof="+joinedEnum)
		}
		switch schema.Value.Format {
		case "ip":
			validateTags = append(validateTags, "ip")
		case "ipv4":
			validateTags = append(validateTags, "ipv4")
		case "ipv6":
			validateTags = append(validateTags, "ipv6")
		case "email":
			validateTags = append(validateTags, "email")
		}

	case schema.Value.Type.Permits(openapi3.TypeInteger):
		if schema.Value.Min != nil {
			validateTags = append(validateTags, "min="+fmt.Sprint(*schema.Value.Min))
		}
		if schema.Value.Max != nil {
			validateTags = append(validateTags, "max="+fmt.Sprint(*schema.Value.Max))
		}
		if schema.Value.MultipleOf != nil {
			slog.Warn("multipleOf validator is not supported")
		}
		if schema.Value.ExclusiveMax {
			slog.Warn("exclusiveMax validator is not supported")
		}
		if schema.Value.ExclusiveMin {
			slog.Warn("exclusiveMin validator is not supported")
		}
		if len(schema.Value.Enum) > 0 {
			enumStrValues := make([]string, 0, len(schema.Value.Enum))
			for _, enumValue := range schema.Value.Enum {
				enumStrValue := fmt.Sprintf("%v", enumValue)
				enumStrValues = append(enumStrValues, enumStrValue)
			}
			joinedEnum := strings.Join(enumStrValues, " ")
			validateTags = append(validateTags, "oneof="+joinedEnum)
		}

	case schema.Value.Type.Permits(openapi3.TypeNumber):
		if schema.Value.Min != nil {
			validateTags = append(validateTags, "min="+fmt.Sprint(*schema.Value.Min))
		}
		if schema.Value.Max != nil {
			validateTags = append(validateTags, "max="+fmt.Sprint(*schema.Value.Max))
		}
		if schema.Value.MultipleOf != nil {
			slog.Warn("multipleOf validator is not supported")
		}
		if schema.Value.ExclusiveMax {
			slog.Warn("exclusiveMax validator is not supported")
		}
		if schema.Value.ExclusiveMin {
			slog.Warn("exclusiveMin validator is not supported")
		}
		if len(schema.Value.Enum) > 0 {
			enumStrValues := make([]string, 0, len(schema.Value.Enum))
			for _, enumValue := range schema.Value.Enum {
				enumStrValue := fmt.Sprintf("%v", enumValue)
				enumStrValues = append(enumStrValues, enumStrValue)
			}
			joinedEnum := strings.Join(enumStrValues, " ")
			validateTags = append(validateTags, "oneof="+joinedEnum)
		}

	case schema.Value.Type.Permits(openapi3.TypeArray):
		if schema.Value.MinItems > 0 {
			validateTags = append(validateTags, "min="+strconv.FormatUint(schema.Value.MinItems, 10))
		}
		if schema.Value.MaxItems != nil {
			validateTags = append(validateTags, "max="+strconv.FormatUint(*schema.Value.MaxItems, 10))
		}
		if schema.Value.UniqueItems {
			validateTags = append(validateTags, "unique")
		}
		validateTags = append(validateTags, "dive")
		itemsValidators := GetSchemaValidators(schema.Value.Items)
		validateTags = append(validateTags, itemsValidators...)
	}

	return validateTags
}
