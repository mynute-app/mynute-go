package lib

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// ErrorStruct defines the error structure returned to clients.// ErrorStruct defines the error structure returned to clients.
type ErrorStruct struct {
	DescriptionEn string   `json:"description_en"`
	DescriptionBr string   `json:"description_br"`
	HTTPStatus    int      `json:"http_status"`
	InnerError    []string `json:"inner_error"` // optional, shown only in dev
}

// WithError attaches an internal error to the ErrorStruct
func (e ErrorStruct) WithError(err error) ErrorStruct {
	if err != nil {
		e.InnerError = append(e.InnerError, err.Error())
	}
	return e
}

// Optional: If you want a printable version
func (e ErrorStruct) Error() string {
	if len(e.InnerError) > 0 {
		errText := ""
		for index, innerErr := range e.InnerError {
			errText += fmt.Sprintf("%d. Inner error : %s \n", index+1, innerErr)
		}
		return fmt.Sprintf("%s \n%s", e.DescriptionEn, errText)
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
	Auth        AuthErrors
	Appointment AppointmentErrors
	Branch      BranchErrors
	Client      ClientErrors
	Company     CompanyErrors
	Employee    EmployeeErrors
	General     GeneralErrors
	Role        RoleErrors
}

type AppointmentErrors struct {
	StartTimeInThePast ErrorStruct
	EndTimeBeforeStart ErrorStruct
}

// Grouped errors per domain
type AuthErrors struct {
	InvalidLogin     ErrorStruct
	NoToken          ErrorStruct
	InvalidToken     ErrorStruct
	Unauthorized     ErrorStruct
	EmailCodeInvalid ErrorStruct
}

type BranchErrors struct {
	ServiceDoesNotBelong                ErrorStruct
	MaxConcurrentAppointments           ErrorStruct
	MaxConcurrentAppointmentsForService ErrorStruct
	MaxConcurrentAppointmentsGeneral    ErrorStruct
}

type ClientErrors struct {
	NotVerified       ErrorStruct
	EmailExists       ErrorStruct
	InvalidClientName ErrorStruct
	InvalidEmail      ErrorStruct
	NotFoundById      ErrorStruct
	CompanyLimit      ErrorStruct
	CompanyIdNotFound ErrorStruct
	ScheduleConflict  ErrorStruct
}

type CompanyErrors struct {
	IdNotFound            ErrorStruct
	CnpjAlreadyExists     ErrorStruct
	NotSame               ErrorStruct
	BranchDoesNotBelong   ErrorStruct
	EmployeeDoesNotBelong ErrorStruct
	ServiceDoesNotBelong  ErrorStruct
}

type EmployeeErrors struct {
	ServiceDoesNotBelong ErrorStruct
	BranchDoesNotBelong  ErrorStruct
	NoWorkScheduleForDay ErrorStruct
	NotAvailableOnDate   ErrorStruct
	ScheduleConflict     ErrorStruct
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

type RoleErrors struct {
	NameReserved ErrorStruct
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
	Appointment: AppointmentErrors{
		StartTimeInThePast: NewError("Start time is in the past", "A data de início está no passado", fiber.StatusBadRequest),
		EndTimeBeforeStart: NewError("End time is before start time", "A data de término é anterior à data de início", fiber.StatusBadRequest),
	},
	Branch: BranchErrors{
		ServiceDoesNotBelong:                NewError("Service is not registered in the branch", "Serviço não está registrado na filial", fiber.StatusBadRequest),
		MaxConcurrentAppointmentsForService: NewError("Maximum concurrent appointments reached for the service in this branch", "Máximo de atendimentos simultâneos atingido para o serviço nesta filial", fiber.StatusBadRequest),
		MaxConcurrentAppointmentsGeneral:    NewError("Maximum concurrent appointments reached for the branch", "Máximo de atendimentos simultâneos atingido para a filial", fiber.StatusBadRequest),
	},
	Client: ClientErrors{
		NotVerified:       NewError("Client not verified", "Usuário não verificado", fiber.StatusUnauthorized),
		EmailExists:       NewError("Email already exists", "Email já cadastrado", fiber.StatusBadRequest),
		InvalidClientName: NewError("Invalid client name", "Nome de usuário inválido", fiber.StatusBadRequest),
		InvalidEmail:      NewError("Invalid email", "Email inválido", fiber.StatusBadRequest),
		NotFoundById:      NewError("Could not find client by ID", "Não foi possível encontrar o usuário pelo ID", fiber.StatusNotFound),
		CompanyLimit:      NewError("Client already has a company associated", "Usuário já possui uma empresa associada", fiber.StatusBadRequest),
		CompanyIdNotFound: NewError("Client company ID not found. This is an internal error", "ID da empresa do usuário não encontrado. Este é um erro interno", fiber.StatusInternalServerError),
		ScheduleConflict:  NewError("Client already has a schedule on this date and time", "Usuário já tem compromisso nesse dia e horário", fiber.StatusBadRequest),
	},
	Company: CompanyErrors{
		IdNotFound:            NewError("Company ID not found or malformed at the request body", "ID da empresa não encontrado ou malformado no corpo da requisição", fiber.StatusBadRequest),
		CnpjAlreadyExists:     NewError("Company CNPJ already exists", "Empresa já cadastrada", fiber.StatusBadRequest),
		NotSame:               NewError("The CompanyID of entities are not all equal.", "O CompanyID das entidades não são todos iguais.", fiber.StatusBadRequest),
		BranchDoesNotBelong:   NewError("Branch does not belong to the company", "Filial não pertence à empresa", fiber.StatusBadRequest),
		EmployeeDoesNotBelong: NewError("Employee does not belong to the company", "Funcionário não pertence à empresa", fiber.StatusBadRequest),
		ServiceDoesNotBelong:  NewError("Service does not belong to the company", "Serviço não pertence à empresa", fiber.StatusBadRequest),
	},
	Employee: EmployeeErrors{
		ServiceDoesNotBelong: NewError("Employee does not have the service registered", "Funcionário não possui o serviço cadastrado", fiber.StatusBadRequest),
		BranchDoesNotBelong:  NewError("Employee does not belong to the branch", "Funcionário não pertence à filial", fiber.StatusBadRequest),
		NoWorkScheduleForDay: NewError("Employee does not have a work schedule for the selected day", "Funcionário não possui um horário de trabalho para o dia selecionado", fiber.StatusBadRequest),
		NotAvailableOnDate:   NewError("Employee is not available on the selected date", "Funcionário não está disponível na data selecionada", fiber.StatusBadRequest),
		ScheduleConflict:     NewError("Employee already has a schedule on this date and time", "Funcionário já tem compromisso nesse dia e horário", fiber.StatusBadRequest),
	},
	General: GeneralErrors{
		InternalError:         NewError("Internal server error while processing the request", "Erro interno do servidor ao processar a requisição", fiber.StatusInternalServerError),
		InterfaceDataNotFound: NewError("Interface data not found", "Dados da interface não encontrados", fiber.StatusInternalServerError),
		RecordNotFound:        NewError("Could not find the specific record", "Não foi possivel encontrar o registro especifico", fiber.StatusNotFound),
		CreatedError:          NewError("Error creating record", "Erro ao criar registro", fiber.StatusBadRequest),
		UpdatedError:          NewError("Error updating record", "Erro ao atualizar registro", fiber.StatusBadRequest),
		DeletedError:          NewError("Error deleting record", "Erro ao deletar registro", fiber.StatusBadRequest),
		NotFoundError:         NewError("Resource not found", "Recurso não encontrado", fiber.StatusNotFound),
	},
	Role: RoleErrors{
		NameReserved: NewError("This role name is reserved for system usage", "Esse nome de cargo é reservado para uso do sistema", fiber.StatusBadRequest),
	},
}
