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
	UserNotFoundById      ErrorStruct
	InvalidToken          ErrorStruct
	UserCompanyLimit      ErrorStruct
}

type AuthErrors struct {
	InvalidLogin ErrorStruct
	NoToken      ErrorStruct
	InvalidToken ErrorStruct
}

type UserErrors struct {
	NotVerified     ErrorStruct
	EmailExists     ErrorStruct
	InvalidUserName ErrorStruct
	InvalidEmail    ErrorStruct
	Unauthroized    ErrorStruct
	NotFoundById    ErrorStruct
	CompanyLimit    ErrorStruct
}

type CompanyErrors struct {
	IdNotFound        ErrorStruct
	CnpjAlreadyExists ErrorStruct
}

type GeneralErrors struct {
	InterfaceDataNotFound ErrorStruct
}

type ErrorCategory struct {
	Auth    AuthErrors
	User    UserErrors
	Company CompanyErrors
	General GeneralErrors
}

var Error = ErrorCategory{
	Auth: AuthErrors{
		InvalidLogin: ErrorStruct{
			DescriptionEn: "Invalid login",
			DescriptionBr: "Login inválido",
			HTTPStatus:    401,
		},
		NoToken: ErrorStruct{
			DescriptionEn: "No token provided",
			DescriptionBr: "Nenhum token fornecido",
			HTTPStatus:    401,
		},
		InvalidToken: ErrorStruct{
			DescriptionEn: "Invalid token",
			DescriptionBr: "Token inválido",
			HTTPStatus:    401,
		},
	},
	User: UserErrors{
		NotVerified: ErrorStruct{
			DescriptionEn: "User not verified",
			DescriptionBr: "Usuário não verificado",
			HTTPStatus:    401,
		},
		EmailExists: ErrorStruct{
			DescriptionEn: "Email already exists",
			DescriptionBr: "Email já cadastrado",
			HTTPStatus:    400,
		},
		InvalidUserName: ErrorStruct{
			DescriptionEn: "Invalid user name",
			DescriptionBr: "Nome de usuário inválido",
			HTTPStatus:    400,
		},
		InvalidEmail: ErrorStruct{
			DescriptionEn: "Invalid email",
			DescriptionBr: "Email inválido",
			HTTPStatus:    400,
		},
		Unauthroized: ErrorStruct{
			DescriptionEn: "You are not authorized to access this resource",
			DescriptionBr: "Você não está autorizado a acessar este recurso",
			HTTPStatus:    401,
		},
		NotFoundById: ErrorStruct{
			DescriptionEn: "Could not find user by ID",
			DescriptionBr: "Não foi possível encontrar o usuário pelo ID",
			HTTPStatus:    404,
		},
		CompanyLimit: ErrorStruct{
			DescriptionEn: "User already has a company associated",
			DescriptionBr: "Usuário já possui uma empresa associada",
			HTTPStatus:    400,
		},
	},
	Company: CompanyErrors{
		IdNotFound: ErrorStruct{
			DescriptionEn: "Company ID not found at the request body",
			DescriptionBr: "ID da empresa não encontrado no corpo da requisição",
			HTTPStatus:    400,
		},
		CnpjAlreadyExists: ErrorStruct{
			DescriptionEn: "Company CNPJ already exists",
			DescriptionBr: "Empresa já cadastrada",
			HTTPStatus:    400,
		},
	},
	General: GeneralErrors{
		InterfaceDataNotFound: ErrorStruct{
			DescriptionEn: "Interface data not found",
			DescriptionBr: "Dados da interface não encontrados",
			HTTPStatus:    500,
		},
	},
}
