package validation

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"log"
	"os"
	"strconv"
	"strings"
)

type ErrorMsg struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func GetField(fe validator.FieldError) string {
	var builder strings.Builder

	for i, char := range fe.Field() {
		if i > 0 && char >= 'A' && char <= 'Z' {
			builder.WriteRune('_')
		}
		builder.WriteRune(char)
	}

	return strings.ToLower(builder.String())
}

func GetErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "The " + fe.Field() + " field must not be left blank."
	case "eqfield":
		return "Password and confirm password doesn't match"
	case "min":
		return "Should be greater than " + fe.Param()
	case "max":
		return "Should be less than " + fe.Param()
	case "secureDomain":
		return "The domain of email should be " + os.Getenv("EMAIL_DOMAIN")
	case "email":
		return "The email field must be an email"
	}

	return fe.Tag()
}

func AddValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("secureDomain", SecureDomain)

		if err != nil {
			log.Fatal("Error add custom validator")
		}
	}
}

var SecureDomain validator.Func = func(fl validator.FieldLevel) bool {
	if email, ok := fl.Field().Interface().(string); ok {
		domain := strings.Split(email, "@")[1]
		if secureEmailDomain, _ := strconv.ParseBool(os.Getenv("SECURE_EMAIL_DOMAIN")); secureEmailDomain == true && domain != os.Getenv("EMAIL_DOMAIN") {
			return false
		}
	}
	return true
}
