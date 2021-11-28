package models

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator"
)

// Model is the base definition for every app model.
type Model struct {
	ID uint `gorm:"primarykey"`
}

// FieldError is an error used to indicate there is field model validation error.
type FieldError struct {
	Field   string
	Type    string
	Param   string
	Message string
}

// Error returns an humanized error message.
func (fe *FieldError) Error() string {
	if len(fe.Message) > 0 {
		return fe.Message
	}

	if len(fe.Field+fe.Type) == 0 {
		return "Unknown error"
	}

	return fmt.Sprintf("Error:Field validation for '%s' failed on the '%s' tag", fe.Field, fe.Type)
}

// ValidationErrors is an error used to indicate there is model validation error.
type ValidationErrors []FieldError

// Error concatenates all the contained field errors messages.
func (v ValidationErrors) Error() string {
	raw := ""
	for _, e := range v {
		if len(raw) > 0 {
			raw += "\n"
		}

		raw += e.Error()
	}

	return raw
}

type recordValidator interface {
	Validate() (bool, ValidationErrors)
}

// Validate run the model tags validations and the Validate() method if the record satisfy a recordValidator.
func Validate(record interface{}) (bool, ValidationErrors) {
	tagValidator := validator.New()
	tagValidator.SetTagName("validate")

	var fieldErrors []FieldError

	// add validations error evaluating the `validate` tag in models
	if err := tagValidator.Struct(record); err != nil {
		var valErr validator.ValidationErrors
		if ok := errors.As(err, &valErr); ok {
			for _, v := range valErr {
				fieldErrors = append(fieldErrors, FieldError{
					Field: v.Field(),
					Type:  v.Tag(),
					Param: v.Param(),
				})
			}
		} else {
			fieldErrors = append(fieldErrors, FieldError{})
		}
	}

	// add validations errors evaluating the `Validate` method if it supported
	if validator, ok := record.(recordValidator); ok {
		_, err := validator.Validate()
		fieldErrors = append(fieldErrors, err...)
	}

	return len(fieldErrors) == 0, fieldErrors
}
