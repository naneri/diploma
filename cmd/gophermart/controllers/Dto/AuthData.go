package Dto

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type AuthData struct {
	Login    string `dipjson:"required"`
	Password string `dipjson:"required"`
}

func (a *AuthData) CheckForEmptyFields() error {
	fmt.Println(a)
	fields := reflect.ValueOf(a).Elem()
	for i := 0; i < fields.NumField(); i++ {

		tags := fields.Type().Field(i).Tag.Get("dipjson")
		if strings.Contains(tags, "required") && fields.Field(i).IsZero() {
			return errors.New("required field is missing")
		}

	}
	return nil
}
