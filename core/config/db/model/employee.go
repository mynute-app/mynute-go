package model

import (
	"agenda-kaki-go/core/lib"
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	WorkSchedule     []WorkRange   `gorm:"foreignKey:EmployeeID;constraint:OnDelete:CASCADE;" json:"work_schedule"`
	Appointments     []Appointment `gorm:"foreignKey:EmployeeID;constraint:OnDelete:CASCADE;" json:"appointments"`
	CompanyID        uuid.UUID     `gorm:"not null;index" json:"company_id"`
	Branches         []*Branch     `gorm:"many2many:employee_branches;constraint:OnDelete:CASCADE;" json:"branches"`
	Services         []*Service    `gorm:"many2many:employee_services;constraint:OnDelete:CASCADE;" json:"services"`
	Roles            []*Role       `gorm:"many2many:employee_roles;constraint:OnDelete:CASCADE;" json:"roles"`
}

func (Employee) TableName() string  { return "employees" }
func (Employee) SchemaType() string { return "company" }

func (e *Employee) BeforeCreate(tx *gorm.DB) error {
	if err := lib.ValidatorV10.Struct(e); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			BadReq := lib.Error.General.BadRequest
			for _, fieldErr := range validationErrors {
				// You can customize the message
				BadReq.WithError(
					fmt.Errorf("field '%s' failed on the '%s' rule", fieldErr.Field(), fieldErr.Tag()),
				)
			}
			return BadReq
		} else {
			return lib.Error.General.InternalError.WithError(err)
		}
	}
	if err := e.HashPassword(); err != nil {
		return err
	}
	return nil
}

func (e *Employee) BeforeUpdate(tx *gorm.DB) error {
	if tx.Statement.Changed("CompanyID") {
		return lib.Error.General.UpdatedError.WithError(errors.New("the CompanyID cannot be changed after creation"))
	}
	if e.Password != "" {
		db_e := &Employee{}
		tx.First(db_e, e.ID)
		if e.Password != db_e.Password && !e.MatchPassword(db_e.Password) {
			if err := lib.ValidatorV10.Var(e.Password, "myPasswordValidation"); err != nil {
				if _, ok := err.(validator.ValidationErrors); ok {
					return lib.Error.General.BadRequest.WithError(fmt.Errorf("password invalid"))
				} else {
					return lib.Error.General.InternalError.WithError(err)
				}
			}
			if err := e.HashPassword(); err != nil {
				return err
			}
		}
	}
	// if !e.WorkSchedule.IsEmpty() && tx.Statement.Changed("work_schedule") {
	// 	if err := e.ValiateWorkSchedule(tx); err != nil {
	// 		return err
	// 	}
	// }
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

// func (e *Employee) ValiateWorkSchedule(tx *gorm.DB) error {
// 	// Avoid running this function if the employee is not created yet
// 	if e.ID == uuid.Nil {
// 		return nil
// 	}
// 	var associatedBranchIDs []uuid.UUID
// 	if err := tx.Table("employee_branches").
// 		Where("employee_id = ?", e.ID).
// 		Pluck("branch_id", &associatedBranchIDs).Error; err != nil {
// 		return fmt.Errorf("failed to fetch employee branches for validation: %w", err)
// 	}
// 	branchSet := make(map[uuid.UUID]bool)
// 	for _, id := range associatedBranchIDs {
// 		branchSet[id] = true
// 	}

// 	var employeeServiceIDs []uuid.UUID
// 	if err := tx.Table("employee_services").
// 		Where("employee_id = ?", e.ID).
// 		Pluck("service_id", &employeeServiceIDs).Error; err != nil {
// 		return fmt.Errorf("failed to fetch employee services for validation: %w", err)
// 	}

// 	employeeServiceSet := make(map[uuid.UUID]bool)
// 	for _, id := range employeeServiceIDs {
// 		employeeServiceSet[id] = true
// 	}

// 	allRanges := e.WorkSchedule.GetAllRanges()                   // Use the existing helper on WorkSchedule
// 	branchServiceCache := make(map[uuid.UUID]map[uuid.UUID]bool) // Cache branch services

// 	for _, wr := range allRanges {
// 		// --- Basic Range Checks ---
// 		if wr.BranchID == uuid.Nil {
// 			return fmt.Errorf("work range has nil BranchID")
// 		} else if wr.Start == "" || wr.End == "" {
// 			return fmt.Errorf("work range (Branch: %s) has empty start/end time", wr.BranchID)
// 		} else if wr.Start == wr.End {
// 			return fmt.Errorf("work range (Branch: %s) has start time equal to end time", wr.BranchID)
// 		} else if wr.Start > wr.End {
// 			return fmt.Errorf("work range (Branch: %s) has start time after end time", wr.BranchID)
// 		} else if _, err := time.Parse("15:04", wr.Start); err != nil {
// 			return fmt.Errorf("work range (Branch: %s) has invalid start time format: %s", wr.BranchID, err.Error())
// 		} else if _, err := time.Parse("15:04", wr.End); err != nil {
// 			return fmt.Errorf("work range (Branch: %s) has invalid end time format: %s", wr.BranchID, err.Error())
// 		}

// 		// --- Check 1: Branch Association ---
// 		if _, exists := branchSet[wr.BranchID]; !exists {
// 			return fmt.Errorf("work range BranchID %s is not associated with the employee %s", wr.BranchID, e.ID)
// 		}

// 		// --- Check 2: Service Validity (if specified in WorkRange) ---
// 		if len(wr.Services) > 0 {
// 			branchServices, ok := branchServiceCache[wr.BranchID]

// 			// If not in cache, fetch services for the branch
// 			if !ok {
// 				var branchServiceIDs []uuid.UUID
// 				if err := tx.Table("branch_services").
// 					Where("branch_id = ?", wr.BranchID).
// 					Pluck("service_id", &branchServiceIDs).Error; err != nil {
// 					return fmt.Errorf("failed to fetch services for branch %s: %w", wr.BranchID, err)
// 				}
// 				branchServices = make(map[uuid.UUID]bool)
// 				for _, id := range branchServiceIDs {
// 					branchServices[id] = true
// 				}
// 				branchServiceCache[wr.BranchID] = branchServices
// 			}

// 			// Check each specified service in the work range
// 			for _, serviceID := range wr.Services {
// 				if serviceID == uuid.Nil {
// 					return fmt.Errorf("work range (Branch: %s) contains a nil ServiceID", wr.BranchID)
// 				}
// 				// Must be offered by Employee
// 				if _, exists := employeeServiceSet[serviceID]; !exists {
// 					return fmt.Errorf("work range (Branch: %s) lists ServiceID %s which the employee %s does not provide", wr.BranchID, serviceID, e.ID)
// 				}
// 				// Must be offered by Branch
// 				if _, exists := branchServices[serviceID]; !exists {
// 					return fmt.Errorf("work range (Branch: %s) lists ServiceID %s which the branch does not offer", wr.BranchID, serviceID)
// 				}
// 			}
// 		}
// 	}

// 	return nil // All validations passed
// }

func (e *Employee) GetWorkRangeForDay(day time.Weekday) []*WorkRange {
	var WorkRanges []*WorkRange
	if len(e.WorkSchedule) == 0 {
		return WorkRanges
	}
	for _, wr := range e.WorkSchedule {
		if wr.Weekday == day {
			WorkRanges = append(WorkRanges, &wr)
		}
	}
	return WorkRanges
}

func (e *Employee) HasOverlappingWorkRange(tx *gorm.DB, wr *WorkRange) error {
	var emp_work_schedule []WorkRange
	if err := tx.Find(&emp_work_schedule, "employee_id = ? AND weekday = ?", e.ID, wr.Weekday).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return lib.Error.General.InternalError.WithError(err)
		}
	}

	// Check for overlapping work ranges
	for _, existing := range emp_work_schedule {
		if existing.Overlaps(wr) {
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("work range to create (%s ~ %s) overlaps with existing range %s (%s ~ %s)", wr.StartTime, wr.EndTime, existing.ID, existing.StartTime, existing.EndTime))
		}
	}

	return nil
}

func (e *Employee) AddWorkRange(tx *gorm.DB, wr *WorkRange) error {
	if wr.EmployeeID != e.ID {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("work range employee ID does not match employee ID"))
	}

	if err := tx.Create(wr).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	if err := tx.Preload(clause.Associations).First(&e, "id = ?", e.ID.String()).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.UpdatedError.WithError(fmt.Errorf("employee not found"))
		}
		return lib.Error.General.InternalError.WithError(err)
	}
	return nil
}

func (e *Employee) RemoveWorkRange(tx *gorm.DB, wr *WorkRange) error {
	if wr.EmployeeID != e.ID {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("work range employee ID does not match employee ID"))
	}

	if err := tx.Delete(wr).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	if err := tx.Preload(clause.Associations).First(&e, "id = ?", e.ID.String()).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.UpdatedError.WithError(fmt.Errorf("employee not found"))
		}
		return lib.Error.General.InternalError.WithError(err)
	}
	return nil
}

func (e *Employee) HasService(tx *gorm.DB, serviceID uuid.UUID) bool {
	var count int64
	if err := tx.Raw("SELECT COUNT(*) FROM employee_services WHERE employee_id = ? AND service_id = ?", e.ID, serviceID).Scan(&count).Error; err != nil {
		return false // Error occurred, assume service does not exist
	}
	return count > 0
}

func (e *Employee) AddService(tx *gorm.DB, service *Service) error {
	if service.CompanyID != e.CompanyID {
		return lib.Error.Company.NotSame
	}

	eID := e.ID.String()
	sID := service.ID.String()

	if err := tx.Exec("INSERT INTO employee_services (employee_id, service_id) VALUES (?, ?)", eID, sID).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	if err := tx.Preload(clause.Associations).First(&e, "id = ?", eID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.UpdatedError.WithError(fmt.Errorf("employee not found"))
		}
		return lib.Error.General.InternalError.WithError(err)
	}
	return nil
}

func (e *Employee) RemoveService(tx *gorm.DB, service *Service) error {
	if service.CompanyID != e.CompanyID {
		return lib.Error.Company.NotSame
	}

	eID := e.ID.String()
	sID := service.ID.String()

	if err := tx.Exec("DELETE FROM employee_services WHERE employee_id = ? AND service_id = ?", eID, sID).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	if err := tx.Preload(clause.Associations).First(&e, "id = ?", eID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.UpdatedError.WithError(fmt.Errorf("employee not found"))
		}
		return lib.Error.General.InternalError.WithError(err)
	}
	return nil
}

func (e *Employee) HasBranch(tx *gorm.DB, branchID uuid.UUID) bool {
	var count int64
	// Check if the employee exists in the branch
	if err := tx.Raw("SELECT COUNT(*) FROM employee_branches WHERE employee_id = ? AND branch_id = ?", e.ID, branchID).Scan(&count).Error; err != nil {
		return false // Error occurred, assume employee does not exist in the branch
	}
	return count > 0
}

func (e *Employee) AddBranch(tx *gorm.DB, branch *Branch) error {
	if branch.CompanyID != e.CompanyID {
		return lib.Error.Company.NotSame
	}

	eID := e.ID.String()
	bID := branch.ID.String()

	if err := tx.Exec("INSERT INTO employee_branches (employee_id, branch_id) VALUES (?, ?)", eID, bID).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	if err := tx.Preload(clause.Associations).First(&e, "id = ?", eID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.UpdatedError.WithError(fmt.Errorf("employee not found"))
		}
		return lib.Error.General.InternalError.WithError(err)
	}
	return nil
}

func (e *Employee) RemoveBranch(tx *gorm.DB, branch *Branch) error {
	if branch.CompanyID != e.CompanyID {
		return lib.Error.Company.NotSame
	}

	eID := e.ID.String()
	bID := branch.ID.String()

	if err := tx.Exec("DELETE FROM employee_branches WHERE employee_id = ? AND branch_id = ?", eID, bID).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	if err := tx.Preload(clause.Associations).First(&e, "id = ?", eID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.UpdatedError.WithError(fmt.Errorf("employee not found"))
		}
		return lib.Error.General.InternalError.WithError(err)
	}
	return nil
}

func (e *Employee) AddRole(tx *gorm.DB, role *Role) error {
	if role.CompanyID != nil && *role.CompanyID != e.CompanyID {
		return lib.Error.Company.NotSame
	}

	eID := e.ID.String()
	rID := role.ID.String()

	if err := tx.Exec("INSERT INTO employee_roles (employee_id, role_id) VALUES (?, ?)", eID, rID).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	if err := tx.Preload(clause.Associations).First(&e, "id = ?", eID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.UpdatedError.WithError(fmt.Errorf("employee not found"))
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

func (e *Employee) RemoveRole(tx *gorm.DB, role *Role) error {
	if role.CompanyID != nil && *role.CompanyID != e.CompanyID {
		return lib.Error.Company.NotSame
	}

	eID := e.ID.String()
	rID := role.ID.String()

	if err := tx.Exec("DELETE FROM employee_roles WHERE employee_id = ? AND role_id = ?", eID, rID).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	if err := tx.Preload(clause.Associations).First(&e, "id = ?", eID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.UpdatedError.WithError(fmt.Errorf("employee not found"))
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}
