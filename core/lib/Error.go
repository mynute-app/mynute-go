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
	EmailExists           ErrorStruct
	InvalidUserName       ErrorStruct
	InvalidEmail          ErrorStruct
	CompanyIDNotFound     ErrorStruct
	Unauthroized          ErrorStruct
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
	EmailExists: ErrorStruct{
		DescriptionEn: "Email already exists",
		DescriptionBr: "Email já cadastrado",
		ID:            4,
		HTTPStatus:    409,
	},
	InvalidUserName: ErrorStruct{
		DescriptionEn: "Invalid user name",
		DescriptionBr: "Nome de usuário inválido",
		ID:            5,
		HTTPStatus:    400,
	},
	InvalidEmail: ErrorStruct{
		DescriptionEn: "Invalid email",
		DescriptionBr: "Email inválido",
		ID:            6,
		HTTPStatus:    400,
	},
	CompanyIDNotFound: ErrorStruct{
		DescriptionEn: "Company ID not found at the request body",
		DescriptionBr: "ID da empresa não encontrado no corpo da requisição",
		ID:            7,
		HTTPStatus:    400,
	},
	Unauthroized: ErrorStruct{
		DescriptionEn: "You are not authorized to access this resource",
		DescriptionBr: "Você não está autorizado a acessar este recurso",
		ID:            8,
		HTTPStatus:    401,
	},
}
