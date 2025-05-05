package model

import (
	mJSON "agenda-kaki-go/core/config/db/model/json"
	"agenda-kaki-go/core/lib"
	"fmt"
	"time"

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
	WorkSchedule     mJSON.WorkSchedule  `gorm:"type:jsonb" json:"work_schedule"`
	Appointments     []Appointment `gorm:"foreignKey:EmployeeID;constraint:OnDelete:CASCADE;" json:"appointments"`
	CompanyID        uuid.UUID     `gorm:"not null;index" json:"company_id"`
	Branches         []*Branch     `gorm:"many2many:employee_branches;constraint:OnDelete:CASCADE;" json:"branches"`
	Services         []*Service    `gorm:"many2many:employee_services;constraint:OnDelete:CASCADE;" json:"services"`
	Roles            []*Role       `gorm:"many2many:employee_roles;constraint:OnDelete:CASCADE;" json:"roles"`
}


func (e *Employee) BeforeCreate(tx *gorm.DB) error {
	e.WorkSchedule = mJSON.WorkSchedule{} // Do not ever let someone create an employee with WorkSchedule set.
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
	if !e.WorkSchedule.IsEmpty() && tx.Statement.Changed("work_schedule") {
		if err := e.ValiateWorkSchedule(tx); err != nil {
			return err
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

func (e *Employee) ValiateWorkSchedule(tx *gorm.DB) error {
	// Avoid running this function if the employee is not created yet
	if e.ID == uuid.Nil {
		return nil
	}
	var associatedBranchIDs []uuid.UUID
	if err := tx.Table("employee_branches").
		Where("employee_id = ?", e.ID).
		Pluck("branch_id", &associatedBranchIDs).Error; err != nil {
		return fmt.Errorf("failed to fetch employee branches for validation: %w", err)
	}
	branchSet := make(map[uuid.UUID]bool)
	for _, id := range associatedBranchIDs {
		branchSet[id] = true
	}

	var employeeServiceIDs []uuid.UUID
	if err := tx.Table("employee_services").
		Where("employee_id = ?", e.ID).
		Pluck("service_id", &employeeServiceIDs).Error; err != nil {
		return fmt.Errorf("failed to fetch employee services for validation: %w", err)
	}

	employeeServiceSet := make(map[uuid.UUID]bool)
	for _, id := range employeeServiceIDs {
		employeeServiceSet[id] = true
	}

	allRanges := e.WorkSchedule.GetAllRanges()                   // Use the existing helper on WorkSchedule
	branchServiceCache := make(map[uuid.UUID]map[uuid.UUID]bool) // Cache branch services

	for _, wr := range allRanges {
		// --- Basic Range Checks ---
		if wr.BranchID == uuid.Nil {
			return fmt.Errorf("work range has nil BranchID")
		} else if wr.Start == "" || wr.End == "" {
			return fmt.Errorf("work range (Branch: %s) has empty start/end time", wr.BranchID)
		} else if wr.Start == wr.End {
			return fmt.Errorf("work range (Branch: %s) has start time equal to end time", wr.BranchID)
		} else if wr.Start > wr.End {
			return fmt.Errorf("work range (Branch: %s) has start time after end time", wr.BranchID)
		} else if _, err := time.Parse("15:04", wr.Start); err != nil {
			return fmt.Errorf("work range (Branch: %s) has invalid start time format: %s", wr.BranchID, err.Error())
		} else if _, err := time.Parse("15:04", wr.End); err != nil {
			return fmt.Errorf("work range (Branch: %s) has invalid end time format: %s", wr.BranchID, err.Error())
		}

		// --- Check 1: Branch Association ---
		if _, exists := branchSet[wr.BranchID]; !exists {
			return fmt.Errorf("work range BranchID %s is not associated with the employee %s", wr.BranchID, e.ID)
		}

		// --- Check 2: Service Validity (if specified in WorkRange) ---
		if len(wr.Services) > 0 {
			branchServices, ok := branchServiceCache[wr.BranchID]

			// If not in cache, fetch services for the branch
			if !ok {
				var branchServiceIDs []uuid.UUID
				if err := tx.Table("branch_services").
					Where("branch_id = ?", wr.BranchID).
					Pluck("service_id", &branchServiceIDs).Error; err != nil {
					return fmt.Errorf("failed to fetch services for branch %s: %w", wr.BranchID, err)
				}
				branchServices = make(map[uuid.UUID]bool)
				for _, id := range branchServiceIDs {
					branchServices[id] = true
				}
				branchServiceCache[wr.BranchID] = branchServices
			}

			// Check each specified service in the work range
			for _, serviceID := range wr.Services {
				if serviceID == uuid.Nil {
					return fmt.Errorf("work range (Branch: %s) contains a nil ServiceID", wr.BranchID)
				}
				// Must be offered by Employee
				if _, exists := employeeServiceSet[serviceID]; !exists {
					return fmt.Errorf("work range (Branch: %s) lists ServiceID %s which the employee %s does not provide", wr.BranchID, serviceID, e.ID)
				}
				// Must be offered by Branch
				if _, exists := branchServices[serviceID]; !exists {
					return fmt.Errorf("work range (Branch: %s) lists ServiceID %s which the branch does not offer", wr.BranchID, serviceID)
				}
			}
		}
	}

	return nil // All validations passed
}