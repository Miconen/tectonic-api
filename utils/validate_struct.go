package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func ValidateStruct(s interface{}) error {
	structType := reflect.TypeOf(s)
	if structType.Kind() != reflect.Struct {
		return errors.New("Input should be a struct")
	}

	structVal := reflect.ValueOf(s)
	fieldNum := structVal.NumField()

	var e strings.Builder

	for i := 0; i < fieldNum; i++ {
		field := structVal.Field(i)
		fieldName := structType.Field(i).Name

		if !field.IsValid() || field.IsZero() {
			if e.Len() > 0 {
				e.WriteString("; ")
			}
			e.WriteString(fmt.Sprintf("%s is not set", fieldName))
		}
	}

	if e.Len() > 0 {
		return errors.New(e.String())
	}

	return nil
}
