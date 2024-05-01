package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in:          "Invalid Input for function: Not a struct",
			expectedErr: ErrInvalidInput,
		},
		{
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "John",
				Age:    25,
				Email:  "test@example.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: ValidationErrors{},
		},
		{
			in: User{
				ID:     "12345678901234564567890123456",
				Name:   "John",
				Age:    25,
				Email:  "test@example.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: ValidationErrors{{Field: "ID", Err: ErrInvalidStrLen}},
		},
		{
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "John",
				Age:    25,
				Email:  "test@example.com",
				Role:   "not_valid role",
				Phones: []string{"12345678901"},
			},
			expectedErr: ValidationErrors{{Field: "Role", Err: ErrInvalidStrNotListed}},
		},
		{
			in: User{
				ID:     "123456789012345678901234567",
				Name:   "John",
				Age:    54,
				Email:  "testexample.com",
				Role:   "not_valid role",
				Phones: []string{"12345678901", "12345"},
			},
			expectedErr: ValidationErrors{
				{Field: "ID", Err: ErrInvalidStrLen},
				{Field: "Age", Err: ErrInvalidIntMax},
				{Field: "Email", Err: ErrInvalidStrValue},
				{Field: "Role", Err: ErrInvalidStrNotListed},
				{Field: "Phones", Err: ErrInvalidStrLen},
			},
		},
		{
			in: App{
				Version: "1.0.0",
			},
			expectedErr: ValidationErrors{},
		},
		{
			in: App{
				Version: "1.0.02",
			},
			expectedErr: ValidationErrors{{Field: "Version", Err: ErrInvalidStrLen}},
		},
		{
			in:          Response{Code: 700},
			expectedErr: ValidationErrors{},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			require.Equal(t, tt.expectedErr, err)
		})
	}
}
