package generator_test

import (
	"encoding/json"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestStrings(t *testing.T) {
	type StringModel string
	type RequiredString struct {
		Str *StringModel `json:"str" validate:"required,min=3"`
	}

	type OptionalString struct {
		Str *StringModel `json:"str,omitempty" validate:"omitempty,min=3"`
	}
	t.Run("TestEmptyJson", func(t *testing.T) {
		jsonData := `{}`

		var optionalResult OptionalString
		err := json.Unmarshal([]byte(jsonData), &optionalResult)
		assert.NoError(t, err)

		v := validator.New(validator.WithRequiredStructEnabled())
		err = v.Struct(&optionalResult)
		assert.NoError(t, err)

		var requiredResult RequiredString
		err = json.Unmarshal([]byte(jsonData), &requiredResult)
		assert.NoError(t, err)

		err = v.Struct(&requiredResult)
		assert.Error(t, err)
	})
	t.Run("TestValidJson", func(t *testing.T) {
		jsonData := `{"str":"test"}`

		var optionalResult OptionalString
		err := json.Unmarshal([]byte(jsonData), &optionalResult)
		assert.NoError(t, err)

		v := validator.New(validator.WithRequiredStructEnabled())
		err = v.Struct(&optionalResult)
		assert.NoError(t, err)

		var requiredResult RequiredString
		err = json.Unmarshal([]byte(jsonData), &requiredResult)
		assert.NoError(t, err)

		err = v.Struct(&requiredResult)
		assert.NoError(t, err)
	})
	t.Run("TestInvalidJson", func(t *testing.T) {
		jsonData := `{"str":"t"}`

		var optionalResult OptionalString
		err := json.Unmarshal([]byte(jsonData), &optionalResult)
		assert.NoError(t, err)

		v := validator.New(validator.WithRequiredStructEnabled())
		err = v.Struct(&optionalResult)
		assert.Error(t, err)

		var requiredResult RequiredString
		err = json.Unmarshal([]byte(jsonData), &requiredResult)
		assert.NoError(t, err)

		err = v.Struct(&requiredResult)
		assert.Error(t, err)
	})
}

func TestEnumWithEmptyString(t *testing.T) {
	type Model struct {
		Field *string `json:"field,omitempty" validate:"omitempty,oneof='' WinXP Win78 Win10"`
	}
	v := validator.New(validator.WithRequiredStructEnabled())

	ptr := func(s string) *string { return &s }
	for _, tc := range []struct {
		name  string
		field *string
		valid bool
	}{
		{"empty string", ptr(""), true},
		{"known value", ptr("WinXP"), true},
		{"unknown value", ptr("Linux"), false},
		{"nil omitempty", nil, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := v.Struct(&Model{Field: tc.field})
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestArrayValidatorsDive(t *testing.T) {
	type ObjectModelArrayField []string
	type ObjectModel struct {
		ArrayField *ObjectModelArrayField `json:"array_field,omitempty" validate:"omitempty,min=1,max=3,unique,dive,min=3,max=10"`
	}
	t.Run("TestEmptyJson", func(t *testing.T) {
		jsonData := `{}`

		var optionalResult ObjectModel
		err := json.Unmarshal([]byte(jsonData), &optionalResult)
		assert.NoError(t, err)

		v := validator.New(validator.WithRequiredStructEnabled())
		err = v.Struct(&optionalResult)
		assert.NoError(t, err)
	})
	t.Run("TestNullArray", func(t *testing.T) {
		jsonData := `{"array_field": null}`

		var optionalResult ObjectModel
		err := json.Unmarshal([]byte(jsonData), &optionalResult)
		assert.NoError(t, err)

		v := validator.New(validator.WithRequiredStructEnabled())
		err = v.Struct(&optionalResult)
		assert.NoError(t, err)
	})
	t.Run("TestEmptyArray", func(t *testing.T) {
		jsonData := `{"array_field": []}`

		var optionalResult ObjectModel
		err := json.Unmarshal([]byte(jsonData), &optionalResult)
		assert.NoError(t, err)

		v := validator.New(validator.WithRequiredStructEnabled())
		err = v.Struct(&optionalResult)
		assert.Error(t, err)
	})
	t.Run("TestValidArray", func(t *testing.T) {
		jsonData := `{"array_field": ["test"]}`

		var optionalResult ObjectModel
		err := json.Unmarshal([]byte(jsonData), &optionalResult)
		assert.NoError(t, err)

		v := validator.New(validator.WithRequiredStructEnabled())
		err = v.Struct(&optionalResult)
		assert.NoError(t, err)
	})
	t.Run("TestIntArray", func(t *testing.T) {
		jsonData := `{"array_field": [1]}`

		var optionalResult ObjectModel
		err := json.Unmarshal([]byte(jsonData), &optionalResult)
		assert.Error(t, err)
	})
	t.Run("TestTooShortStringArray", func(t *testing.T) {
		jsonData := `{"array_field": ["t"]}`

		var optionalResult ObjectModel
		err := json.Unmarshal([]byte(jsonData), &optionalResult)
		assert.NoError(t, err)

		v := validator.New(validator.WithRequiredStructEnabled())
		err = v.Struct(&optionalResult)
		assert.Error(t, err)
	})
	t.Run("TestTooLongStringArray", func(t *testing.T) {
		jsonData := `{"array_field": ["too_long_string_in_array"]}`

		var optionalResult ObjectModel
		err := json.Unmarshal([]byte(jsonData), &optionalResult)
		assert.NoError(t, err)

		v := validator.New(validator.WithRequiredStructEnabled())
		err = v.Struct(&optionalResult)
		assert.Error(t, err)
	})
	t.Run("TestNonUniqueArray", func(t *testing.T) {
		jsonData := `{"array_field": ["test", "test"]}`

		var optionalResult ObjectModel
		err := json.Unmarshal([]byte(jsonData), &optionalResult)
		assert.NoError(t, err)

		v := validator.New(validator.WithRequiredStructEnabled())
		err = v.Struct(&optionalResult)
		assert.Error(t, err)
	})
	t.Run("TestNonUniqueArray", func(t *testing.T) {
		jsonData := `{"array_field": ["test1", "test2", "test3", "test4"]}`

		var optionalResult ObjectModel
		err := json.Unmarshal([]byte(jsonData), &optionalResult)
		assert.NoError(t, err)

		v := validator.New(validator.WithRequiredStructEnabled())
		err = v.Struct(&optionalResult)
		assert.Error(t, err)
	})
}
