package generator_test

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/sintoniastrategy/validgo-gen/internal/generator"
	"github.com/stretchr/testify/assert"
)

func TestGetSchemaValidators_EnumWithEmptyString(t *testing.T) {
	schema := &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type: &openapi3.Types{openapi3.TypeString},
			Enum: []any{"", "WinXP", "Win78"},
		},
	}

	tags := generator.GetSchemaValidators(schema)
	assert.Equal(t, []string{"oneof='' WinXP Win78"}, tags)
}

func TestGetSchemaValidators_EnumWithSpaceInValue(t *testing.T) {
	schema := &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type: &openapi3.Types{openapi3.TypeString},
			Enum: []any{"hello world", "foo"},
		},
	}

	tags := generator.GetSchemaValidators(schema)
	assert.Equal(t, []string{"oneof='hello world' foo"}, tags)
}
