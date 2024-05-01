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
)

func (v ValidationErrors) Error() string {
	errString := ""

	for _, err := range v {
		errString += fmt.Sprintf("%s: %s\n", err.Field, err.Err)
	}

	return errString
}

func Validate(v interface{}) error {
	iv := reflect.ValueOf(v)
	validationErrors := ValidationErrors{}

	if iv.Kind() != reflect.Struct {
		return ErrInvalidInput
	}

	t := iv.Type()

	for i := 0; i < iv.NumField(); i++ {
		field := t.Field(i)
		value := iv.Field(i)

		valueType := value.Type()
		valueTypeKind := value.Type().Kind()

		validateTag := field.Tag.Get("validate")

		if validateTag == "" {
			continue
		}

		validationRestrictions := strings.Split(validateTag, "|")

		for _, restriction := range validationRestrictions {
			validationKey := strings.Split(restriction, ":")[0]
			validationValue := strings.Split(restriction, ":")[1]

			if valueTypeKind == reflect.String {
				err := validateString(value.String(), validationKey, validationValue)

				if err != nil {
					validationErrors = append(validationErrors, ValidationError{
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
						validationErrors = append(validationErrors, ValidationError{
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
					validationErrors = append(validationErrors, ValidationError{
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
						validationErrors = append(validationErrors, ValidationError{
							Field: field.Name,
							Err:   err,
						})
						break
					}
				}
			}
		}
	}

	// Place your code here.
	return validationErrors
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
			return err
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
			return err
		}
		if value < minValue {
			return ErrInvalidIntMin
		}
	case "max":
		maxValue, err := strconv.Atoi(validationValue)
		if err != nil {
			return err
		}
		if value > maxValue {
			return ErrInvalidIntMax
		}
	case "range":
		minValue, err := strconv.Atoi(strings.Split(validationValue, ",")[0])
		if err != nil {
			return err
		}
		maxValue, err := strconv.Atoi(strings.Split(validationValue, ",")[1])
		if err != nil {
			return err
		}
		if value < minValue || value > maxValue {
			return ErrInvalidIntRange
		}
	}

	return nil
}
