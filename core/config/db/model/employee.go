package model

import (
	"agenda-kaki-go/core/lib"
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Employee struct {
	BaseModel
	Name             string        `gorm:"type:varchar(100);not null" json:"name"`
	Surname          string        `gorm:"type:varchar(100)" json:"surname"`
	Email            string        `gorm:"type:varchar(100);not null;uniqueIndex" json:"email" validate:"required,email"`
	Phone            string        `gorm:"type:varchar(20);not null;uniqueIndex" json:"phone" validate:"required,e164"`
	Tags             []string      `gorm:"type:json" json:"tags"`
	Password         string        `gorm:"type:varchar(255);not null" json:"password" validate:"required,myPasswordValidation"`
	ChangePassword   bool          `gorm:"default:false;not null" json:"change_password"`
	VerificationCode string        `gorm:"type:varchar(100)" json:"verification_code"`
	Verified         bool          `gorm:"default:false;not null" json:"verified"`
	SlotTimeDiff     uint          `gorm:"default:30;not null" json:"slot_time_diff"`
	WorkSchedule     WorkSchedule  `gorm:"type:jsonb" json:"work_schedule"`
	Appointments     []Appointment `gorm:"foreignKey:EmployeeID;constraint:OnDelete:CASCADE;" json:"appointments"`
	CompanyID        uuid.UUID     `gorm:"not null;index" json:"company_id"`
	Company          Company       `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE;" json:"company"`
	Branches         []*Branch     `gorm:"many2many:employee_branches;constraint:OnDelete:CASCADE;" json:"branches"`
	Services         []*Service    `gorm:"many2many:employee_services;constraint:OnDelete:CASCADE;" json:"services"`
	Roles            []*Role       `gorm:"many2many:employee_roles;constraint:OnDelete:CASCADE;" json:"roles"`
}

type WorkSchedule struct {
	Monday    []WorkRange `json:"monday"`
	Tuesday   []WorkRange `json:"tuesday"`
	Wednesday []WorkRange `json:"wednesday"`
	Thursday  []WorkRange `json:"thursday"`
	Friday    []WorkRange `json:"friday"`
	Saturday  []WorkRange `json:"saturday"`
	Sunday    []WorkRange `json:"sunday"`
}

type WorkRange struct {
	Start    string    `json:"start"` // Store as "15:30:00"
	End      string    `json:"end"`   // Store as "18:00:00"
	BranchID uuid.UUID `json:"branch_id"`
}

// Implement driver.Valuer
func (ws WorkSchedule) Value() (driver.Value, error) {
	return json.Marshal(ws)
}

// Implement sql.Scanner
func (ws *WorkSchedule) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan WorkSchedule: expected []byte")
	}

	return json.Unmarshal(bytes, ws)
}

func (e *Employee) BeforeCreate(tx *gorm.DB) error {
	if err := lib.ValidatorV10.Struct(e); err != nil {
		return err
	}
	if err := e.HashPassword(); err != nil {
		return err
	}
	return nil
}

func (e *Employee) BeforeUpdate(tx *gorm.DB) error {
	if e.Password != "" {
		db_e := &Employee{}
		tx.First(db_e, e.ID)
		if e.Password != db_e.Password && !e.MatchPassword(db_e.Password) {
			if err := lib.ValidatorV10.Var(e.Password, "myPasswordValidation"); err != nil {
				return err
			}
			if err := e.HashPassword(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *Employee) MatchPassword(hashedPass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(e.Password))
	return err == nil
}

// Method to set hashed password:
func (e *Employee) HashPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(e.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	e.Password = string(hash)
	return nil
}
