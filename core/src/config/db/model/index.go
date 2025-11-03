package model

import (
	"fmt"
	authModel "mynute-go/auth/config/db/model"
	"mynute-go/core/src/lib"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid();<-:create" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (m *BaseModel) BeforeSave(tx *gorm.DB) (err error) {
	if m.ID != uuid.Nil && m.ID.Variant() != uuid.RFC4122 {
		errMsg := fmt.Errorf("BeforeSave: Invalid UUID variant for ID %s in %T", m.ID.String(), m)
		return lib.Error.General.UpdatedError.WithError(errMsg)
	}
	return nil
}

var TenantModels = []any{
	&Appointment{},
	&AppointmentArchive{},
	&BranchServiceDensity{},
	&BranchWorkRange{},
	&Branch{},
	&EmployeeServiceDensity{},
	&EmployeeWorkRange{},
	&Employee{},
	&Service{},
	&Payment{},
}

// MainDBModels are business models that live in the main database
var MainDBModels = []any{
	&Sector{},
	&Company{},
	&Holiday{},
	&Client{},   // Moved to auth DB (users)
	&Employee{}, // Moved to auth DB (users)
	&Role{},     // Company-specific roles stay in main DB
	&Subdomain{},
	&ClientAppointment{},
}

// AuthDBModels are authentication/authorization models that live in the auth database
var AuthDBModels = []any{
	&authModel.EndPoint{},
	&authModel.PolicyRule{},
	&authModel.Resource{},
	&authModel.Property{},
	&Client{},    // User authentication
	&Employee{},  // User authentication
	&Admin{},     // Admin users
	&RoleAdmin{}, // Admin roles
	&Role{},      // System roles (Owner, GM, BM, Employee)
}

// GeneralModels combines all models (for backwards compatibility and utilities)
var GeneralModels = []any{
	&Sector{},
	&Company{},
	&Holiday{},
	&Client{},
	&authModel.EndPoint{},
	&Role{},
	&authModel.PolicyRule{},
	&authModel.Resource{},
	&authModel.Property{},
	&Subdomain{},
	&ClientAppointment{},
	&Admin{},
	&RoleAdmin{},
}

func GetModelFromTableName(tableName string) (any, string, error) {
	type has_TableName_fnc interface {
		TableName() string
		SchemaType() string
	}
	for _, rawModel := range TenantModels {
		if model, ok := rawModel.(has_TableName_fnc); ok {
			if model.TableName() == tableName {
				return model, model.SchemaType(), nil
			}
		} else {
			return nil, "", lib.Error.General.InternalError.WithError(fmt.Errorf("model %T does not implement TableName() or SchemaType()", rawModel))
		}
	}
	for _, rawModel := range GeneralModels {
		if model, ok := rawModel.(has_TableName_fnc); ok {
			if model.TableName() == tableName {
				return model, model.SchemaType(), nil
			}
		} else {
			return nil, "", lib.Error.General.InternalError.WithError(fmt.Errorf("model %T does not implement TableName() or SchemaType()", rawModel))
		}
	}
	return nil, "", lib.Error.General.InternalError.WithError(fmt.Errorf("model not found for table name: %s", tableName))
}

