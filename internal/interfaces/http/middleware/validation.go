package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func ValidateRequest(obj interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.ShouldBindJSON(obj); err != nil {
			var errors []ValidationError

			if verr, ok := err.(validator.ValidationErrors); ok {
				for _, f := range verr {
					err := ValidationError{
						Field:   f.Field(),
						Message: getErrorMsg(f),
					}
					errors = append(errors, err)
				}
			} else {
				errors = append(errors, ValidationError{
					Field:   "request",
					Message: err.Error(),
				})
			}

			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "validation_failed",
				"details": errors,
			})
			c.Abort()
			return
		}
		c.Set("validated_data", obj)
		c.Next()
	}
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Value is too short"
	case "max":
		return "Value is too long"
	case "alphanum":
		return "Only alphanumeric characters are allowed"
	default:
		return "Invalid value"
	}
}
