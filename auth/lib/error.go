package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"mynute-go/auth/config/namespace"

	"github.com/gofiber/fiber/v2"
)

// ErrorStruct defines the error structure returned to clients.// ErrorStruct defines the error structure returned to clients.
type ErrorStruct struct {
	DescriptionEn string         `json:"description_en"`
	DescriptionBr string         `json:"description_br"`
	HTTPStatus    int            `json:"http_status"`
	InnerError    map[int]string `json:"inner_error"` // optional, shown only in dev
}

// WithError attaches an internal error to the ErrorStruct
func (e ErrorStruct) WithError(err error) ErrorStruct {
	if e.InnerError == nil {
		e.InnerError = make(map[int]string)
	}

	newErr := ErrorStruct{
		InnerError:    make(map[int]string),
		DescriptionEn: e.DescriptionEn,
		DescriptionBr: e.DescriptionBr,
		HTTPStatus:    e.HTTPStatus,
	}

	// Copia os erros existentes
	i := 1
	for _, msg := range e.InnerError {
		newErr.InnerError[i] = msg
		i++
	}

	// Adiciona os novos erros
	if errStr, ok := err.(ErrorStruct); ok {
		newErr.DescriptionBr = errStr.DescriptionBr
		newErr.DescriptionEn = errStr.DescriptionEn
		newErr.HTTPStatus = errStr.HTTPStatus

		for _, msg := range errStr.InnerError {
			newErr.InnerError[i] = msg
			i++
		}
	} else if err != nil {
		newErr.InnerError[i] = err.Error()
	}

	return newErr
}

// Optional: If you want a printable version
func (e ErrorStruct) Error() string {
	e_byte, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf("ErrorStruct: %s", err.Error())
	}
	return string(e_byte)
}

// Optional: If you want a byte array version
func (e ErrorStruct) Byte() []byte {
	e_byte, err := json.Marshal(e)
	if err != nil {
		return []byte(fmt.Sprintf("ErrorStruct: %s", err.Error()))
	}
	return e_byte
}

// ToJSON converts the ErrorStruct to a JSON string.
func (e ErrorStruct) ToJSON() string {
	jsonData, err := json.Marshal(e)
	if err != nil {
		return `{"error": "failed to convert to JSON"}`
	}
	return string(jsonData)
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
	Auth               AuthErrors
	Appointment        AppointmentErrors
	AppointmentArchive AppointmentArchiveErrors
	Branch             BranchErrors
	Client             ClientErrors
	Company            CompanyErrors
	Employee           EmployeeErrors
	General            GeneralErrors
	Role               RoleErrors
	Validation         ValidationErrors
}

type AppointmentErrors struct {
	StartTimeInThePast           ErrorStruct
	EndTimeBeforeStart           ErrorStruct
	NotFound                     ErrorStruct
	MissingRequiredIDs           ErrorStruct // New: For missing FKs
	AssociationLoadFailed        ErrorStruct // New: General DB error loading related data
	InvalidServiceDuration       ErrorStruct // New: Service duration <= 0
	InvalidWorkScheduleFormat    ErrorStruct // New: Error parsing work schedule time
	UpdateFailed                 ErrorStruct // New: Generic update DB error
	CreateFailed                 ErrorStruct // New: Generic create DB error
	HistoryLoggingFailed         ErrorStruct // New: Failure during history log save
	HistoryManualUpdateForbidden ErrorStruct // New: Manual update of history log not allowed
	CancelledAppointmentUpdate   ErrorStruct // New: Attempt to modify a cancelled appointment
}

type AppointmentArchiveErrors struct {
	IdNotSet        ErrorStruct // New: For nil ID in appointment archive
	UpdateForbidden ErrorStruct
	DeleteForbidden ErrorStruct
}

// Grouped errors per domain
type AuthErrors struct {
	InvalidLogin         ErrorStruct
	NoToken              ErrorStruct
	InvalidToken         ErrorStruct
	Unauthorized         ErrorStruct
	EmailCodeInvalid     ErrorStruct
	CompanyHeaderMissing ErrorStruct // New: For missing X-Company-ID header
	CompanyHeaderInvalid ErrorStruct // New: For invalid X-Company-ID header
}

type BranchErrors struct {
	NotFound                  ErrorStruct
	ServiceDoesNotBelong      ErrorStruct
	MaxConcurrentAppointments ErrorStruct
	MaxCapacityReached        ErrorStruct
	MaxServiceCapacityReached ErrorStruct
}

type ClientErrors struct {
	NotVerified       ErrorStruct
	EmailExists       ErrorStruct
	InvalidClientName ErrorStruct
	InvalidEmail      ErrorStruct
	NotFound          ErrorStruct
	CompanyLimit      ErrorStruct
	CompanyNotFound   ErrorStruct
	ScheduleConflict  ErrorStruct
	CompanyIdNotFound ErrorStruct
}

type CompanyErrors struct {
	NotFound              ErrorStruct
	CnpjAlreadyExists     ErrorStruct
	NotSame               ErrorStruct
	BranchDoesNotBelong   ErrorStruct
	EmployeeDoesNotBelong ErrorStruct
	ServiceDoesNotBelong  ErrorStruct
	CompanyMismatch       ErrorStruct
	IdUpdateForbidden     ErrorStruct
	CouldNotCreateOwner   ErrorStruct
	SchemaIsEmpty         ErrorStruct
}

type EmployeeErrors struct {
	NotFound                 ErrorStruct
	ServiceDoesNotBelong     ErrorStruct
	BranchDoesNotBelong      ErrorStruct
	NoWorkScheduleForDay     ErrorStruct
	NotAvailableOnDate       ErrorStruct
	ScheduleConflict         ErrorStruct
	LacksService             ErrorStruct // New (More specific than ServiceDoesNotBelong)
	NotAvailableWorkSchedule ErrorStruct // New (More specific than NotAvailableOnDate)
}

type GeneralErrors struct {
	InterfaceDataNotFound ErrorStruct
	BadRequest            ErrorStruct
	RecordNotFound        ErrorStruct
	CreatedError          ErrorStruct
	UpdatedError          ErrorStruct
	DeletedError          ErrorStruct
	ResourceNotFoundError ErrorStruct
	InternalError         ErrorStruct
	AuthError             ErrorStruct
	SessionNotFound       ErrorStruct // New: For missing session in DB
	DatabaseError         ErrorStruct // New: General DB error
	TooManyRequests       ErrorStruct // New: For rate limiting
}

type RoleErrors struct {
	NameReserved ErrorStruct
	NilCompanyID ErrorStruct
}

type ValidationErrors struct {
	Failed ErrorStruct // New: General validation failure bucket
}

// Global error instances
var Error = ErrorCategory{
	Auth: AuthErrors{
		InvalidLogin:         NewError("Invalid login credentials", "Credenciais de login inválidas", fiber.StatusUnauthorized),
		NoToken:              NewError("Authorization token not provided", "Token de autorização não fornecido", fiber.StatusUnauthorized),
		InvalidToken:         NewError("Authorization token is invalid or expired", "Token de autorização inválido ou expirado", fiber.StatusUnauthorized),
		Unauthorized:         NewError("You are not authorized for this action", "Você não está autorizado para esta ação", fiber.StatusForbidden),
		EmailCodeInvalid:     NewError("Email verification code is invalid", "Código de verificação do email inválido", fiber.StatusBadRequest),
		CompanyHeaderMissing: NewError(fmt.Sprintf(`%s required at headers`, namespace.HeadersKey.Company), fmt.Sprintf(`%s requerido no cabeçalho`, namespace.HeadersKey.Company), fiber.StatusBadRequest),
		CompanyHeaderInvalid: NewError(fmt.Sprintf(`%s invalid at headers`, namespace.HeadersKey.Company), fmt.Sprintf(`%s inválido no cabeçalho`, namespace.HeadersKey.Company), fiber.StatusBadRequest),
	},
	Appointment: AppointmentErrors{
		StartTimeInThePast:           NewError("Appointment start time cannot be in the past", "A data de início do compromisso não pode ser no passado", fiber.StatusBadRequest),
		EndTimeBeforeStart:           NewError("Appointment end time must be after start time", "A data de término do compromisso deve ser posterior à data de início", fiber.StatusBadRequest),
		NotFound:                     NewError("Appointment not found", "Compromisso não encontrado", fiber.StatusNotFound),
		MissingRequiredIDs:           NewError("Appointment creation requires valid ServiceID, EmployeeID, ClientID, BranchID, and CompanyID", "Criação do compromisso requer ServiceID, EmployeeID, ClientID, BranchID e CompanyID válidos", fiber.StatusBadRequest),
		AssociationLoadFailed:        NewError("Failed to load associated data for appointment validation", "Falha ao carregar dados associados para validação do compromisso", fiber.StatusInternalServerError),
		InvalidServiceDuration:       NewError("Service duration must be positive", "A duração do serviço deve ser positiva", fiber.StatusBadRequest),
		InvalidWorkScheduleFormat:    NewError("Employee work schedule time format is invalid, expected HH:MM", "Formato de horário de trabalho do funcionário inválido, esperado HH:MM", fiber.StatusBadRequest),
		UpdateFailed:                 NewError("Failed to update appointment in database", "Falha ao atualizar compromisso no banco de dados", fiber.StatusInternalServerError),
		CreateFailed:                 NewError("Failed to create appointment in database", "Falha ao criar compromisso no banco de dados", fiber.StatusInternalServerError),
		HistoryLoggingFailed:         NewError("Failed to save appointment history log", "Falha ao salvar histórico do compromisso", fiber.StatusInternalServerError),
		CancelledAppointmentUpdate:   NewError("Cannot modify a cancelled appointment", "Não é possível modificar um compromisso cancelado", fiber.StatusForbidden),
		HistoryManualUpdateForbidden: NewError("Manual update of appointment log is not allowed", "Atualização manual do histórico não é permitida", fiber.StatusForbidden),
	},
	AppointmentArchive: AppointmentArchiveErrors{
		IdNotSet:        NewError("Appointment archive ID cannot be nil", "ID do arquivo de compromisso não pode ser nulo", fiber.StatusBadRequest),
		UpdateForbidden: NewError("Can not update archived appointments", "Não é possível atualizar compromissos arquivados", fiber.StatusForbidden),
		DeleteForbidden: NewError("Can not delete archived appointments", "Não é possível deletar compromissos arquivados", fiber.StatusForbidden),
	},
	Branch: BranchErrors{
		NotFound:                  NewError("Branch not found", "Filial não encontrada", fiber.StatusNotFound),
		ServiceDoesNotBelong:      NewError("The selected service is not offered by this branch", "O serviço selecionado não é oferecido por esta filial", fiber.StatusBadRequest),
		MaxCapacityReached:        NewError("Branch maximum concurrent appointment capacity reached", "Capacidade máxima de compromissos simultâneos da filial atingida", fiber.StatusConflict),                                 // 409 Conflict better?
		MaxServiceCapacityReached: NewError("Branch maximum concurrent capacity for this specific service reached", "Capacidade máxima de compromissos simultâneos da filial para este serviço atingida", fiber.StatusConflict), // 409 Conflict better?
	},
	Client: ClientErrors{
		NotFound:          NewError("Client not found", "Cliente não encontrado", fiber.StatusNotFound),
		ScheduleConflict:  NewError("Client already has a conflicting appointment", "Cliente já possui um compromisso conflitante", fiber.StatusConflict), // 409 Conflict
		NotVerified:       NewError("Client not verified", "Usuário não verificado", fiber.StatusUnauthorized),
		EmailExists:       NewError("Email already exists", "Email já cadastrado", fiber.StatusBadRequest),
		InvalidClientName: NewError("Invalid client name", "Nome de usuário inválido", fiber.StatusBadRequest),
		InvalidEmail:      NewError("Invalid email", "Email inválido", fiber.StatusBadRequest),
		CompanyLimit:      NewError("Client already has a company associated", "Usuário já possui uma empresa associada", fiber.StatusBadRequest),
		CompanyIdNotFound: NewError("Client company ID not found. This is an internal error", "ID da empresa do usuário não encontrado. Este é um erro interno", fiber.StatusInternalServerError),
	},
	Company: CompanyErrors{
		NotFound:              NewError("Company not found", "Empresa não encontrada", fiber.StatusNotFound),
		BranchDoesNotBelong:   NewError("Branch does not belong to the specified company", "Filial não pertence à empresa especificada", fiber.StatusBadRequest),
		EmployeeDoesNotBelong: NewError("Employee does not belong to the specified company", "Funcionário não pertence à empresa especificada", fiber.StatusBadRequest),
		ServiceDoesNotBelong:  NewError("Service does not belong to the specified company", "Serviço não pertence à empresa especificada", fiber.StatusBadRequest),
		CnpjAlreadyExists:     NewError("Company CNPJ already exists", "Empresa já cadastrada", fiber.StatusBadRequest),
		CompanyMismatch:       NewError("Company ID mismatch between related entities", "Incompatibilidade de ID da empresa entre entidades relacionadas", fiber.StatusBadRequest),
		IdUpdateForbidden:     NewError("Company ID cannot be updated", "ID da empresa não pode ser atualizado", fiber.StatusBadRequest),
		CouldNotCreateOwner:   NewError("Could not create owner account for the company", "Não foi possível criar a conta de proprietário para a empresa", fiber.StatusInternalServerError),
		SchemaIsEmpty:         NewError("Company schema is empty", "Esquema da empresa está vazio", fiber.StatusBadRequest),
	},
	Employee: EmployeeErrors{
		NotFound:                 NewError("Employee not found", "Funcionário não encontrado", fiber.StatusNotFound),
		ServiceDoesNotBelong:     NewError("Employee does not provide the selected service (legacy check)", "Funcionário não oferece o serviço selecionado (verificação legada)", fiber.StatusBadRequest), // Keep? Or rely on LacksService?
		BranchDoesNotBelong:      NewError("Employee is not assigned to the selected branch", "Funcionário não está alocado na filial selecionada", fiber.StatusBadRequest),
		NoWorkScheduleForDay:     NewError("Employee has no work schedule defined for the selected day", "Funcionário não tem horário de trabalho definido para o dia selecionado", fiber.StatusBadRequest),                                   // Kept for initial check
		NotAvailableOnDate:       NewError("Employee is not available according to work schedule at the requested time", "Funcionário não está disponível de acordo com o horário de trabalho no horário solicitado", fiber.StatusBadRequest), // More specific than NotAvailableWorkSchedule
		ScheduleConflict:         NewError("Employee already has a conflicting appointment", "Funcionário já possui um compromisso conflitante", fiber.StatusConflict),                                                                        // 409 Conflict
		LacksService:             NewError("Employee does not provide the specified service", "Funcionário não oferece o serviço especificado", fiber.StatusBadRequest),
		NotAvailableWorkSchedule: NewError("Employee is not scheduled to work at the requested time/branch", "Funcionário não está escalado para trabalhar no horário/filial solicitados", fiber.StatusBadRequest),
	},
	General: GeneralErrors{
		InternalError:         NewError("Internal server error while processing the request", "Erro interno do servidor ao processar a requisição", fiber.StatusInternalServerError),
		InterfaceDataNotFound: NewError("Interface data not found", "Dados da interface não encontrados", fiber.StatusInternalServerError),
		RecordNotFound:        NewError("Could not find the specific record", "Não foi possivel encontrar o registro especifico", fiber.StatusNotFound),
		CreatedError:          NewError("Error creating record", "Erro ao criar registro", fiber.StatusBadRequest),
		UpdatedError:          NewError("Error updating record", "Erro ao atualizar registro", fiber.StatusBadRequest),
		DeletedError:          NewError("Error deleting record", "Erro ao deletar registro", fiber.StatusBadRequest),
		ResourceNotFoundError: NewError("Resource not found", "Recurso não encontrado", fiber.StatusNotFound),
		BadRequest:            NewError("Bad request", "Requisição inválida", fiber.StatusBadRequest),
		AuthError:             NewError("Internal Server Error while authenticating", "Erro Interno de Servidor enquanto autenticando", fiber.StatusInternalServerError),
		SessionNotFound:       NewError("Database session not found in context", "Sessão de banco de dados não encontrada no contexto", fiber.StatusInternalServerError),
		DatabaseError:         NewError("An internal error occurred regarding the database", "Ocorreu um erro interno relacionado ao banco de dados", fiber.StatusInternalServerError),
		TooManyRequests:       NewError("Too many requests, please try again later", "Muitas requisições, por favor tente novamente mais tarde", fiber.StatusTooManyRequests),
	},
	Role: RoleErrors{
		NameReserved: NewError("This role name is reserved for system usage", "Esse nome de cargo é reservado para uso do sistema", fiber.StatusBadRequest),
		NilCompanyID: NewError("The role has a nil company ID", "O cargo tem um ID de empresa nulo", fiber.StatusBadRequest),
	},
	Validation: ValidationErrors{
		Failed: NewError("Input validation failed", "Falha na validação dos dados de entrada", fiber.StatusBadRequest),
	},
}

