package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

var (
	ErrInvalidInput        = errors.New("validation failed: input is not a struct")
	ErrInvalidStrLen       = errors.New("invalid length")
	ErrInvalidStrNotListed = errors.New("value not listed in allowed values")
	ErrInvalidStrValue     = errors.New("invalid value")
	ErrInvalidIntMin       = errors.New("value is smaller than min")
	ErrInvalidIntMax       = errors.New("value is bigger than max")
	ErrInvalidIntRange     = errors.New("value is out of range")
	ErrUnknown             = errors.New("unknown error")
)

func (v ValidationErrors) Error() string {
	errString := strings.Builder{}

	for i, err := range v {
		errString.WriteString(fmt.Sprintf("%s: %s", err.Field, err.Err))
		if i != len(v)-1 {
			errString.WriteString("\n")
		}
	}

	return errString.String()
}

func Validate(v interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if rErr, ok := r.(error); ok {
				err = rErr
			} else {
				err = ErrUnknown
			}
		}
	}()

	iv := reflect.ValueOf(v)
	validationErrors := ValidationErrors{}

	if iv.Kind() != reflect.Struct {
		return ErrInvalidInput
	}

	t := iv.Type()
	for i := 0; i < iv.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		value := iv.Field(i)

		validateTag := field.Tag.Get("validate")

		if validateTag == "" {
			continue
		}

		validationRestrictions := strings.Split(validateTag, "|")

		for _, restriction := range validationRestrictions {
			validationErrors = validateByRestriction(field, value, restriction, validationErrors)
		}
	}

	return validationErrors
}

func validateByRestriction(
	field reflect.StructField,
	value reflect.Value,
	restriction string,
	ve ValidationErrors,
) ValidationErrors {
	validationKey := strings.Split(restriction, ":")[0]
	validationValue := strings.Split(restriction, ":")[1]

	valueType := value.Type()
	valueTypeKind := value.Type().Kind()

	if valueTypeKind == reflect.String {
		err := validateString(value.String(), validationKey, validationValue)
		if err != nil {
			ve = append(ve, ValidationError{
				Field: field.Name,
				Err:   err,
			})
		}
	}
	if valueTypeKind == reflect.Slice && valueType.Elem().Kind() == reflect.String {
		var err error
		for i := 0; i < value.Len(); i++ {
			err = validateString(value.Index(i).String(), validationKey, validationValue)
			if err != nil {
				ve = append(ve, ValidationError{
					Field: field.Name,
					Err:   err,
				})
				break
			}
		}
	}

	if valueTypeKind == reflect.Int {
		err := validateInt(int(value.Int()), validationKey, validationValue)
		if err != nil {
			ve = append(ve, ValidationError{
				Field: field.Name,
				Err:   err,
			})
		}
	}

	if valueTypeKind == reflect.Slice && valueType.Elem().Kind() == reflect.Int {
		var err error
		for i := 0; i < value.Len(); i++ {
			err = validateInt(int(value.Int()), validationKey, validationValue)
			if err != nil {
				ve = append(ve, ValidationError{
					Field: field.Name,
					Err:   err,
				})
				break
			}
		}
	}

	return ve
}

func validateString(value string, validationKey, validationValue string) error {
	switch validationKey {
	case "in":
		allowedValues := strings.Split(validationValue, ",")
		for _, v := range allowedValues {
			if value == v {
				return nil
			}
		}
		return ErrInvalidStrNotListed

	case "len":
		strLen, err := strconv.Atoi(validationValue)
		if err != nil {
			panic(err)
		}

		if strLen != len(value) {
			return ErrInvalidStrLen
		}
	case "regexp":
		if !regexp.MustCompile(validationValue).MatchString(value) {
			return ErrInvalidStrValue
		}
	}

	return nil
}

func validateInt(value int, validationKey, validationValue string) error {
	switch validationKey {
	case "min":
		minValue, err := strconv.Atoi(validationValue)
		if err != nil {
			panic(err)
		}
		if value < minValue {
			return ErrInvalidIntMin
		}
	case "max":
		maxValue, err := strconv.Atoi(validationValue)
		if err != nil {
			panic(err)
		}
		if value > maxValue {
			return ErrInvalidIntMax
		}
	case "range":
		minValue, err := strconv.Atoi(strings.Split(validationValue, ",")[0])
		if err != nil {
			panic(err)
		}
		maxValue, err := strconv.Atoi(strings.Split(validationValue, ",")[1])
		if err != nil {
			panic(err)
		}
		if value < minValue || value > maxValue {
			return ErrInvalidIntRange
		}
	}

	return nil
}
