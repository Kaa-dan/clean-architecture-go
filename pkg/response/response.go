package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/kaa-dan/clean-architecture-go/pkg/errors"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func Success(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Data:    data,
	})
}

func Error(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, Response{
		Success: false,
		Error:   message,
	})
}

func HandleError(c *gin.Context, err error) {
	statusCode := errors.GetHTTPStatusCode(err)
	Error(c, statusCode, err.Error())
}

func ValidationError(c *gin.Context, err error) {
	var validationErrors []string

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, validationErr := range validationErrs {
			validationErrors = append(validationErrors, getValidationErrorMessage(validationErr))
		}
	} else {
		validationErrors = append(validationErrors, err.Error())
	}

	c.JSON(http.StatusBadRequest, Response{
		Success: false,
		Error:   "Validation failed",
		Data:    validationErrors,
	})
}

func getValidationErrorMessage(err validator.FieldError) string {
	field := err.Field()
	tag := err.Tag()

	switch tag {
	case "required":
		return field + " is required"
	case "email":
		return field + " must be a valid email address"
	case "min":
		return field + " must be at least " + err.Param() + " characters long"
	case "max":
		return field + " must be at most " + err.Param() + " characters long"
	case "alphanum":
		return field + " must contain only alphanumeric characters"
	default:
		return field + " is invalid"
	}
}
