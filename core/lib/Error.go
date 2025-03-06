package lib

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

// ToJSON converts the ErrorStruct to a JSON string
func (e ErrorStruct) ToJSON() string {
	jsonData, err := json.Marshal(e)
	if err != nil {
		return `{"error": "failed to convert to JSON"}`
	}
	return string(jsonData)
}

func (e ErrorStruct) SendToClient(c *fiber.Ctx) error {
	res := &SendResponse{Ctx: c}
	return res.Send(e.HTTPStatus, e.ToJSON())
}

// ErrorStruct defines the error structure
type ErrorStruct struct {
	DescriptionEn string `json:"description_en"`
	DescriptionBr string `json:"description_br"`
	ID            int    `json:"id"`
	HTTPStatus    int    `json:"http_status"`
}

// ErrorTypes holds different error types
type ErrorTypes struct {
	InterfaceDataNotFound ErrorStruct
	InvalidLogin          ErrorStruct
	UserNotVerified       ErrorStruct
}

// MyErrors defines a list of predefined errors
var MyErrors = ErrorTypes{
	InterfaceDataNotFound: ErrorStruct{
		DescriptionEn: "Interface data not found",
		DescriptionBr: "Dados da interface não encontrados",
		ID:            1,
		HTTPStatus:    500,
	},
	InvalidLogin: ErrorStruct{
		DescriptionEn: "Invalid login",
		DescriptionBr: "Login inválido",
		ID:            2,
		HTTPStatus:    401,
	},
	UserNotVerified: ErrorStruct{
		DescriptionEn: "User not verified",
		DescriptionBr: "Usuário não verificado",
		ID:            3,
		HTTPStatus:    401,
	},
}
