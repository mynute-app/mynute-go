package model

import (
	"agenda-kaki-go/core/lib"
	"fmt"
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
	&EmployeeWorkRange{},
	&BranchWorkRange{},
	&Appointment{},
	&AppointmentArchive{},
	&Branch{},
	&Employee{},
	&Service{},
	&Payment{},
}

var GeneralModels = []any{
	&Sector{},
	&Company{},
	&Holiday{},
	&Client{},
	&EndPoint{},
	&Role{},
	&PolicyRule{},
	&Resource{},
	&Property{},
	&Subdomain{},
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
