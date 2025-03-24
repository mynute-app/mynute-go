package lib

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// ErrorStruct defines the error structure returned to clients.// ErrorStruct defines the error structure returned to clients.
type ErrorStruct struct {
	DescriptionEn string `json:"description_en"`
	DescriptionBr string `json:"description_br"`
	HTTPStatus    int    `json:"http_status"`
	InnerError    string `json:"inner_error"` // optional, shown only in dev
}

// WithError attaches an internal error to the ErrorStruct
func (e ErrorStruct) WithError(err error) ErrorStruct {
	e.InnerError = err.Error()
	return e
}

// Optional: If you want a printable version
func (e ErrorStruct) Error() string {
	if e.InnerError != "" {
		return fmt.Sprintf("%s: %s", e.DescriptionEn, e.InnerError)
	}
	return e.DescriptionEn
}

// ToJSON converts the ErrorStruct to a JSON string.
func (e ErrorStruct) ToJSON() string {
	jsonData, err := json.Marshal(e)
	if err != nil {
		return `{"error": "failed to convert to JSON"}`
	}
	return string(jsonData)
}

// SendToClient sends the error response to the client.
func (e ErrorStruct) SendToClient(c *fiber.Ctx) error {
	return c.Status(e.HTTPStatus).JSON(e)
}

// Unwrap allows error comparison via `errors.Is(...)`
func (e ErrorStruct) Unwrap() error {
	return errors.New(e.DescriptionEn)
}

// Helper to create an ErrorStruct.
func NewError(en, br string, status int) ErrorStruct {
	return ErrorStruct{
		DescriptionEn: en,
		DescriptionBr: br,
		HTTPStatus:    status,
	}
}

// ErrorCategory groups all domain-specific error types.
type ErrorCategory struct {
	Auth    AuthErrors
	Client  ClientErrors
	Company CompanyErrors
	General GeneralErrors
}

// Grouped errors per domain
type AuthErrors struct {
	InvalidLogin     ErrorStruct
	NoToken          ErrorStruct
	InvalidToken     ErrorStruct
	Unauthorized     ErrorStruct
	EmailCodeInvalid ErrorStruct
}

type ClientErrors struct {
	NotVerified       ErrorStruct
	EmailExists       ErrorStruct
	InvalidClientName ErrorStruct
	InvalidEmail      ErrorStruct
	NotFoundById      ErrorStruct
	CompanyLimit      ErrorStruct
	CompanyIdNotFound ErrorStruct
}

type CompanyErrors struct {
	IdNotFound        ErrorStruct
	CnpjAlreadyExists ErrorStruct
	NotSame           ErrorStruct
}

type GeneralErrors struct {
	InterfaceDataNotFound ErrorStruct
	RecordNotFound        ErrorStruct
	CreatedError          ErrorStruct
	UpdatedError          ErrorStruct
	DeletedError          ErrorStruct
	NotFoundError         ErrorStruct
	InternalError         ErrorStruct
}

// Global error instances
var Error = ErrorCategory{
	Auth: AuthErrors{
		InvalidLogin:     NewError("Invalid login", "Login inválido", fiber.StatusUnauthorized),
		NoToken:          NewError("No token provided", "Nenhum token fornecido", fiber.StatusUnauthorized),
		InvalidToken:     NewError("Invalid token", "Token inválido", fiber.StatusUnauthorized),
		Unauthorized:     NewError("You are not authorized to access this resource", "Você não está autorizado a acessar este recurso", fiber.StatusUnauthorized),
		EmailCodeInvalid: NewError("Email's verification code is invalid", "Código de verificação do email inválido", fiber.StatusBadRequest),
	},
	Client: ClientErrors{
		NotVerified:       NewError("Client not verified", "Usuário não verificado", fiber.StatusUnauthorized),
		EmailExists:       NewError("Email already exists", "Email já cadastrado", fiber.StatusBadRequest),
		InvalidClientName: NewError("Invalid client name", "Nome de usuário inválido", fiber.StatusBadRequest),
		InvalidEmail:      NewError("Invalid email", "Email inválido", fiber.StatusBadRequest),
		NotFoundById:      NewError("Could not find client by ID", "Não foi possível encontrar o usuário pelo ID", fiber.StatusNotFound),
		CompanyLimit:      NewError("Client already has a company associated", "Usuário já possui uma empresa associada", fiber.StatusBadRequest),
		CompanyIdNotFound: NewError("Client company ID not found. This is an internal error", "ID da empresa do usuário não encontrado. Este é um erro interno", fiber.StatusInternalServerError),
	},
	Company: CompanyErrors{
		IdNotFound:        NewError("Company ID not found or malformed at the request body", "ID da empresa não encontrado ou malformado no corpo da requisição", fiber.StatusBadRequest),
		CnpjAlreadyExists: NewError("Company CNPJ already exists", "Empresa já cadastrada", fiber.StatusBadRequest),
		NotSame:           NewError("The CompanyID of entities are not all equal.", "O CompanyID das entidades não são todos iguais.", fiber.StatusBadRequest),
	},
	General: GeneralErrors{
		InternalError:         NewError("Internal server error while processing the request", "Erro interno do servidor ao processar a requisição", fiber.StatusInternalServerError),
		InterfaceDataNotFound: NewError("Interface data not found", "Dados da interface não encontrados", fiber.StatusInternalServerError),
		RecordNotFound:        NewError("Record not found", "Registro não encontrado", fiber.StatusNotFound),
		CreatedError:          NewError("Error creating record", "Erro ao criar registro", fiber.StatusInternalServerError),
		UpdatedError:          NewError("Error updating record", "Erro ao atualizar registro", fiber.StatusInternalServerError),
		DeletedError:          NewError("Error deleting record", "Erro ao deletar registro", fiber.StatusInternalServerError),
		NotFoundError:         NewError("Resource not found", "Recurso não encontrado", fiber.StatusNotFound),
	},
}
