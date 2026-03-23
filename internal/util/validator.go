package util

import (
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var mobileValidator validator.Func = func(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	pattern := `^1[3-9]\d{9}$`
	regex := regexp.MustCompile(pattern)

	return regex.MatchString(mobile)
}

func RegisterCustomValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("mobile", mobileValidator)
		if err != nil {
			panic(err)
		}
	}
}
