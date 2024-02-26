package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationFildError struct {
	Field string
	Err   error
}

type ValidationFildErrors []ValidationFildError

func (v ValidationFildErrors) Error() string {
	return fmt.Sprintf("%#v", v)
}

type ValidationSliceError struct {
	N   int
	Err error
}

type ValidationSliceErrors []ValidationSliceError

func (v ValidationSliceErrors) Error() string {
	return fmt.Sprintf("%#v", v)
}

type ErrInvalidValue struct {
	ErrMsg string
}

func (e ErrInvalidValue) Error() string {
	return e.ErrMsg
}

var (
	// Программные ошибки.
	ErrNotStruct        = errors.New("received value is not a struct")
	ErrInvalidTeg       = errors.New("invalid tag")
	ErrNotSupportedType = errors.New("not supported type")

	// Ошибки валидации.
	ErrLeng        = ErrInvalidValue{"length does not match"}
	ErrNotMatchReg = ErrInvalidValue{"value is not match regular expression"}
	ErrNotInclude  = ErrInvalidValue{"value is not included in the set"}
	ErrLessMin     = ErrInvalidValue{"value is less than the set minimum"}
	ErrGreatMax    = ErrInvalidValue{"value is greater than the set maximum"}
)

func Validate(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	return validateStruct(rv)
}

func validateStruct(v reflect.Value) error {
	var validErr ValidationFildErrors
	for i := 0; i < v.NumField(); i++ {
		fieldVal := v.Field(i)
		fieldInfo := v.Type().Field(i)

		if !fieldInfo.IsExported() {
			continue
		}

		validators, exist := fieldInfo.Tag.Lookup("validate")
		if !exist {
			continue
		}

		if err := validate(fieldVal, validators); err != nil {
			var (
				invalidVal    ErrInvalidValue
				invalidSlice  ValidationSliceErrors
				invalidStruct ValidationFildErrors
			)
			switch {
			case errors.As(err, &invalidVal), errors.As(err, &invalidSlice), errors.As(err, &invalidStruct):
				validErr = append(validErr, ValidationFildError{
					Field: fieldInfo.Name,
					Err:   err,
				})
			default:
				return err
			}
		}
	}
	if len(validErr) == 0 {
		return nil
	}
	return validErr
}

// validate вызывает разные валидаторы в зависимости от базового типа
// валидируемого значения.
func validate(v reflect.Value, validators string) error {
	var err error
	switch v.Kind() { //nolint:exhaustive
	case reflect.Struct:
		err = validateNestedStruct(v, validators)
	case reflect.Slice:
		err = validateSlice(v, validators)
	case reflect.String:
		err = validateString(v.String(), validators)
	case reflect.Int:
		err = validateInt(v.Int(), validators)
	default:
		err = ErrNotSupportedType
	}
	return err
}

// validateSlice производит валидацию  каждого значения в слайсе в
// соответствии с переданной строкой валидаторов, разделенных "|".
func validateSlice(s reflect.Value, validators string) error {
	var errs ValidationSliceErrors
	for i := 0; i < s.Len(); i++ {
		v := s.Index(i)
		err := validate(v, validators)
		if err != nil {
			var invalidVal ErrInvalidValue
			if errors.As(err, &invalidVal) {
				errs = append(errs, ValidationSliceError{
					N:   i,
					Err: err,
				})
			} else {
				return err
			}
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errs
}

func validateNestedStruct(v reflect.Value, validators string) error {
	switch validators {
	case "nested":
		return validateStruct(v)
	default:
		return ErrInvalidTeg
	}
}

// validateString производит валидацию строкового значения в соответствии
// с переданной строкой валидаторов, разделенных "|".
//
// Поддерживаемые валидаторы:
// 		"len:%d" 	где %d - целое число. Строка валидна, если ее длина
// 					ровно %d символов.
//
//		"regexp:%s" где %s - регулярное выражение. Строка считается
//				 	валидной, если соответствует %s.
//
//		"in:%s" 	где %s - перечесление множества строк, разделенных ",".
//					Cтрока валидна, если соответствует любой строке из
//					множества %s.
func validateString(v string, validators string) error {
	validatorsList := strings.Split(validators, "|")
	for _, validator := range validatorsList {
		nameArg := strings.SplitN(validator, ":", 2)
		if len(nameArg) != 2 {
			return ErrInvalidTeg
		}
		name, arg := nameArg[0], nameArg[1]
		switch name {
		case "len":
			if err := validatorStringLen(v, arg); err != nil {
				return err
			}
		case "regexp":
			if err := validatorStringReg(v, arg); err != nil {
				return err
			}
		case "in":
			if err := validatorStringIn(v, arg); err != nil {
				return err
			}
		}
	}
	return nil
}

func validatorStringLen(val string, arg string) error {
	vLen, err := strconv.ParseInt(arg, 10, 64)
	if err != nil {
		return ErrInvalidTeg
	}
	if int(vLen) != len(val) {
		return ErrLeng
	}
	return nil
}

func validatorStringReg(val string, arg string) error {
	re, err := regexp.Compile(strings.ReplaceAll(arg, `\\`, `\`))
	if err != nil {
		return ErrInvalidTeg
	}
	if !re.MatchString(val) {
		return ErrNotMatchReg
	}
	return nil
}

func validatorStringIn(val string, arg string) error {
	setStr := strings.Split(arg, ",")
	if !func() bool { // if setStr not contain v
		for _, str := range setStr {
			if str == val {
				return true
			}
		}
		return false
	}() {
		return ErrNotInclude
	}
	return nil
}

// validateInt производит валидацию целочисленного значения в соответствии
// с переданной строкой валидаторов, разделенных "|".
//
// Поддерживаемые валидаторы:
// 		"min:%d"	где %d - целое число. Значение валидно, если не меньше
// 					%d.
//
//		"max:%d"	где %d - целое число. Значение валидно, если не больше
// 					%d.
//
//		"in:%s" 	где %s - строка-перечесление множества целых чисел,
//					разделенных ",". Значение валидно, если соответствует
//					любому значению из множества %s.
func validateInt(v int64, validators string) error {
	validatorsList := strings.Split(validators, "|")
	for _, validator := range validatorsList {
		nameArg := strings.SplitN(validator, ":", 2)
		if len(nameArg) != 2 {
			return ErrInvalidTeg
		}
		name, arg := nameArg[0], nameArg[1]
		switch name {
		case "min":
			if err := validatorIntMin(v, arg); err != nil {
				return err
			}
		case "max":
			if err := validatorIntMax(v, arg); err != nil {
				return err
			}
		case "in":
			if err := validatorIntIn(v, arg); err != nil {
				return err
			}
		}
	}
	return nil
}

func validatorIntMin(val int64, arg string) error {
	min, err := strconv.ParseInt(arg, 10, 64)
	if err != nil {
		return ErrInvalidTeg
	}
	if val < min {
		return ErrLessMin
	}
	return nil
}

func validatorIntMax(val int64, arg string) error {
	max, err := strconv.ParseInt(arg, 10, 64)
	if err != nil {
		return ErrInvalidTeg
	}
	if val > max {
		return ErrGreatMax
	}
	return nil
}

func validatorIntIn(val int64, arg string) error {
	setStr := strings.Split(arg, ",")
	contain, err := func() (bool, error) {
		for _, str := range setStr {
			d, err := strconv.ParseInt(str, 10, 64)
			if err != nil {
				return false, ErrInvalidTeg
			}
			if d == val {
				return true, nil
			}
		}
		return false, nil
	}()
	if err != nil {
		return ErrInvalidTeg
	}
	if !contain {
		return ErrNotInclude
	}
	return nil
}
