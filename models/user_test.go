package models

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestUserValidation(t *testing.T) {
	v := validator.New()

	// https://blog.mamansoft.net/2018/10/15/golang-struct-validation/
	if err := v.Struct(User{Email: "", PasswordDigest: "test"}); err == nil {
		t.Error(err)
	}
	if err := v.Struct(User{Email: "test", PasswordDigest: ""}); err == nil {
		t.Error(err)
	}
	if err := v.Struct(User{Email: "test", PasswordDigest: "test"}); err == nil {
		t.Error(err)
	}
	if err := v.Struct(User{Email: "test@example.com", PasswordDigest: "test"}); err != nil {
		t.Error(err)
	}
}
