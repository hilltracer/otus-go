package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
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
			in: User{
				ID:     "01234567-89ab-cdef-0123-456789abcdef",
				Name:   "denis",
				Age:    38,
				Email:  "hilltracer@qweqwe.ru",
				Role:   "admin",
				Phones: []string{"12345678901", "00000000000"},
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:     "short-id",
				Age:    10,
				Email:  "not-an-email",
				Role:   "unknown",
				Phones: []string{"123"},
			},
			expectedErr: ErrValidation,
		},
		{
			in: App{
				Version: "1.2.3",
			},
			expectedErr: nil,
		},
		{
			in: App{
				Version: "v1",
			},
			expectedErr: ErrValidation,
		},
		{
			in: Response{
				Code: 200,
				Body: "ok",
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 418,
			},
			expectedErr: ErrValidation,
		},
		{
			in:          123, // not a struct
			expectedErr: ErrUnsupportedType,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			if tt.expectedErr == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.expectedErr != nil && !errors.Is(err, tt.expectedErr) {
				t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}
