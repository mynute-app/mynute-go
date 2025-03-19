package model

import (
	"agenda-kaki-go/core/lib"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Employee struct {
	gorm.Model
	Name             string        `gorm:"type:varchar(100);not null" json:"name"`
	Surname          string        `gorm:"type:varchar(100)" json:"surname"`
	Role             string        `gorm:"type:varchar(50);default:user;not null" json:"role"`
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
	CompanyID        uint          `gorm:"not null;index" json:"company_id"`
	Company          *Company      `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE;" json:"company"`
	Branches         []*Branch     `gorm:"many2many:employee_branches;" json:"branches"`
	Services         []*Service    `gorm:"many2many:employee_services;" json:"services"`
}

type WorkSchedule struct {
	Monday    []TimeRange `json:"monday"`
	Tuesday   []TimeRange `json:"tuesday"`
	Wednesday []TimeRange `json:"wednesday"`
	Thursday  []TimeRange `json:"thursday"`
	Friday    []TimeRange `json:"friday"`
	Saturday  []TimeRange `json:"saturday"`
	Sunday    []TimeRange `json:"sunday"`
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

// func (e *Employee) CheckAvailability(service Service, requestedTime time.Time) bool {
// 	serviceEnd := requestedTime.Add(time.Duration(service.Duration) * time.Minute)

// 	for _, slot := range e.AvailableSlots {
// 		if requestedTime.After(slot.Start) || requestedTime.Equal(slot.Start) {
// 			if serviceEnd.Before(slot.End) || serviceEnd.Equal(slot.End) {
// 				return true
// 			}
// 		}
// 	}
// 	return false
// }
