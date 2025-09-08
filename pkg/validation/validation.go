package validation

import (
	"encoding/json"

	"github.com/danielreinhard1129/fiber-clean-arch/pkg/exception"
	"github.com/go-playground/validator/v10"
)

func Validate(modelValidate interface{}) {
	validate := validator.New()
	err := validate.Struct(modelValidate)
	if err != nil {
		var messages []map[string]interface{}
		for _, err := range err.(validator.ValidationErrors) {
			messages = append(messages, map[string]interface{}{
				"field":   err.Field(),
				"message": err.Tag(),
			})
		}

		jsonMessage, errJson := json.Marshal(messages)
		exception.PanicLogging(errJson)

		panic(exception.ValidationError{
			Message: string(jsonMessage),
		})
	}
}
