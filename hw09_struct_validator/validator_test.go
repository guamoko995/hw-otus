package hw09structvalidator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:10"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		Meta   App      `validate:"nested"`
	}

	App struct {
		Version string `validate:"len:5"`
		Name    string
	}

	WhithFailTag struct {
		V string `validate:"len::5"`
	}

	WhithUnsupportedType struct {
		B byte `validate:"nested"`
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

var testData = []struct {
	in          interface{}
	expectedErr error
}{
	{ // case 0
		in: User{
			ID:     "1234567890",
			Age:    18,
			Email:  "example@example.ex",
			Role:   "admin",
			Phones: []string{"+1234567890", "+2345678901"},
			Meta: App{
				Name:    "Приложение",
				Version: "v.1.0",
			},
		},
		expectedErr: nil,
	},
	{ // case 1
		in: User{
			ID:     "123", // invalid len (validate:"len:10")
			Age:    50,
			Email:  "exampl.e@example.ex", // not match reg (validate:"regexp:^\\w+@\\w+\\.\\w+$")
			Role:   "stuff",
			Phones: []string{"+123", "+2345678901"}, // [0] invalid len (validate:"len:11")
			Meta: App{ // validate:"nested"
				Name:    "Приложение",
				Version: "v.1.12", // invalid len (validate:"len:5")
			},
		},
		expectedErr: ValidationFildErrors{
			ValidationFildError{
				Field: "ID",
				Err:   ErrLeng,
			},
			ValidationFildError{
				Field: "Email",
				Err:   ErrNotMatchReg,
			},
			ValidationFildError{
				Field: "Phones",
				Err: ValidationSliceErrors{
					ValidationSliceError{
						N:   0,
						Err: ErrLeng,
					},
				},
			},
			ValidationFildError{
				Field: "Meta",
				Err: ValidationFildErrors{
					ValidationFildError{
						Field: "Version",
						Err:   ErrLeng,
					},
				},
			},
		},
	},
	{ // case 2
		in: User{
			ID:     "1234567890",
			Age:    17, // less min (validate:"min:18|max:50")
			Email:  "example@example.ex",
			Role:   "user",                                  // not include  (validate:"in:admin,stuff")
			Phones: []string{"+1234567890", "+23456789012"}, // [1] invalid len (validate:"len:11")
			Meta: App{
				Name:    "Приложение",
				Version: "v.1.0",
			},
		},
		expectedErr: ValidationFildErrors{
			ValidationFildError{
				Field: "Age",
				Err:   ErrLessMin,
			},
			ValidationFildError{
				Field: "Role",
				Err:   ErrNotInclude,
			},
			ValidationFildError{
				Field: "Phones",
				Err: ValidationSliceErrors{
					ValidationSliceError{
						N:   1,
						Err: ErrLeng,
					},
				},
			},
		},
	},
	{ // case 3
		in: User{
			ID:     "1234567890",
			Age:    51, // great max (validate:"min:18|max:50")
			Email:  "example@example.ex",
			Role:   "admin",
			Phones: []string{"+123", "+23456789012"}, // [0,1] invalid len (validate:"len:11")
			Meta: App{
				Name:    "Приложение",
				Version: "v.1.0",
			},
		},
		expectedErr: ValidationFildErrors{
			ValidationFildError{
				Field: "Age",
				Err:   ErrGreatMax,
			},
			ValidationFildError{
				Field: "Phones",
				Err: ValidationSliceErrors{
					ValidationSliceError{
						N:   0,
						Err: ErrLeng,
					},
					ValidationSliceError{
						N:   1,
						Err: ErrLeng,
					},
				},
			},
		},
	},
	{ // case 4
		in: Response{
			Code: 200,
			Body: "",
		},
		expectedErr: nil,
	},
	{ // case 5
		in: Response{
			Code: 404,
			Body: "404 err",
		},
		expectedErr: nil,
	},
	{ // case 6
		in: Response{
			Code: 500,
			Body: "internal server error",
		},
		expectedErr: nil,
	},
	{ // case 7
		in: Response{
			Code: 201, // not include (validate:"in:200,404,500")
			Body: "created",
		},
		expectedErr: ValidationFildErrors{
			ValidationFildError{
				Field: "Code",
				Err:   ErrNotInclude,
			},
		},
	},
	{ // case 8
		in: Token{
			Header:    []byte("123"),
			Payload:   []byte("123"),
			Signature: []byte("123"),
		},
		expectedErr: nil,
	},
	{ // case 9
		in: WhithFailTag{
			V: "12345", // invalid tag (validate:"len::5")
		},
		expectedErr: ErrInvalidTeg,
	},
	{ // case 10
		in:          "not struct",
		expectedErr: ErrNotStruct,
	},
	{ // case 11
		in: WhithUnsupportedType{
			B: 1,
		},
		expectedErr: ErrNotSupportedType,
	},
}

func TestValidate(t *testing.T) {
	for i, tt := range testData {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			require.Equal(t, tt.expectedErr, err)
		})
	}
}
