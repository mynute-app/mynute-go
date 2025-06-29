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
	Name                 string              `gorm:"type:varchar(100)" validate:"required,min=3,max=100" json:"name"`
	Surname              string              `gorm:"type:varchar(100)" validate:"max=100" json:"surname"`
	Email                string              `gorm:"type:varchar(100);uniqueIndex" validate:"required,email" json:"email"`
	Phone                string              `gorm:"type:varchar(20);uniqueIndex" validate:"required,e164" json:"phone"`
	Tags                 []string            `gorm:"type:json" json:"tags"`
	Password             string              `gorm:"type:varchar(255)" validate:"required,myPasswordValidation" json:"password"`
	ChangePassword       bool                `gorm:"default:false" json:"change_password"`
	VerificationCode     string              `gorm:"type:varchar(100)" json:"verification_code"`
	Verified             bool                `gorm:"default:false" json:"verified"`
	SlotTimeDiff         uint                `gorm:"default:30" json:"slot_time_diff"`
	EmployeeWorkSchedule []EmployeeWorkRange `gorm:"foreignKey:EmployeeID;constraint:OnDelete:CASCADE;" json:"work_schedule"`
	Appointments         []Appointment       `gorm:"foreignKey:EmployeeID;constraint:OnDelete:CASCADE;" json:"appointments"`
	CompanyID            uuid.UUID           `gorm:"not null;index" json:"company_id"`
	Branches             []*Branch           `gorm:"many2many:employee_branches;constraint:OnDelete:CASCADE;" json:"branches"`
	Services             []*Service          `gorm:"many2many:employee_services;constraint:OnDelete:CASCADE;" json:"services"`
	Roles                []*Role             `gorm:"many2many:employee_roles;constraint:OnDelete:CASCADE;" json:"roles"`
}

func (Employee) TableName() string  { return "employees" }
func (Employee) SchemaType() string { return "company" }

func (e *Employee) BeforeCreate(tx *gorm.DB) error {
	if err := lib.MyCustomStructValidator(e); err != nil {
		return err
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
	// if !e.EmployeeWorkSchedule.IsEmpty() && tx.Statement.Changed("work_schedule") {
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

func (e *Employee) GetWorkRangeForDay(day time.Weekday) []EmployeeWorkRange {
	var WorkRanges []EmployeeWorkRange
	if len(e.EmployeeWorkSchedule) == 0 {
		return WorkRanges
	}
	for _, wr := range e.EmployeeWorkSchedule {
		if wr.Weekday == day {
			WorkRanges = append(WorkRanges, wr)
		}
	}
	return WorkRanges
}

// ValidateEmployeeWorkRangeTime checks if the employee work range overlaps with existing work ranges for the employee.
func (e *Employee) ValidateEmployeeWorkRangeTime(tx *gorm.DB, ewr *EmployeeWorkRange) error {
	var emp_work_schedule []EmployeeWorkRange
	if err := tx.Find(&emp_work_schedule, "employee_id = ? AND weekday = ?", e.ID, ewr.Weekday).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return lib.Error.General.InternalError.WithError(err)
		}
	}

	// Check for overlapping work ranges
	for _, existing := range emp_work_schedule {
		overlaps, err := existing.Overlaps(ewr)
		if err != nil {
			return err
		}
		if overlaps {
			startTime, err := lib.Utc2LocalTime(ewr.TimeZone, ewr.StartTime)
			if err != nil {
				return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid start time %s: %w", ewr.StartTime, err))
			}
			endTime, err := lib.Utc2LocalTime(ewr.TimeZone, ewr.EndTime)
			if err != nil {
				return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid end time %s: %w", ewr.EndTime, err))
			}
			existingStartTime, err := lib.Utc2LocalTime(existing.TimeZone, existing.StartTime)
			if err != nil {
				return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid existing start time %s: %w", existing.StartTime, err))
			}
			existingEndTime, err := lib.Utc2LocalTime(existing.TimeZone, existing.EndTime)
			if err != nil {
				return lib.Error.General.BadRequest.WithError(fmt.Errorf("invalid existing end time %s: %w", existing.EndTime, err))
			}
			return lib.Error.General.BadRequest.WithError(fmt.Errorf("employee work range to create (%s ~ %s) overlaps with existing range %s (%s ~ %s)", startTime.Format("15:04"), endTime.Format("15:04"), existing.ID, existingStartTime.Format("15:04"), existingEndTime.Format("15:04")))
		}
	}

	return nil
}

func (e *Employee) RemoveWorkRange(tx *gorm.DB, wr *EmployeeWorkRange) error {
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

func (e *Employee) AddServicesToWorkRange(tx *gorm.DB, wr_id string, services []Service) error {
	if len(services) == 0 {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("no services provided to add to work range"))
	}

	var wr EmployeeWorkRange
	if err := tx.First(&wr, "id = ? AND employee_id = ?", wr_id, e.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.RecordNotFound.WithError(fmt.Errorf("work range not found for this employee ID (%s)", e.ID))
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	if wr.EmployeeID != e.ID {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("work range employee ID does not match employee ID"))
	}

	var branch Branch
	if err := tx.First(&branch, "id = ?", wr.BranchID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.RecordNotFound.WithError(fmt.Errorf("branch not found"))
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	for _, service := range services {
		if err := branch.HasService(tx, &service); err != nil {
			return err
		}
		if err := tx.Exec("INSERT INTO work_schedule_services (work_range_id, service_id) VALUES (?, ?)", wr.ID, service.ID).Error; err != nil {
			return lib.Error.General.InternalError.WithError(err)
		}
	}

	return nil
}

func (e *Employee) RemoveServiceFromWorkRange(tx *gorm.DB, wr_id string, service_id string) error {
	var wr EmployeeWorkRange
	if err := tx.First(&wr, "id = ? AND employee_id = ?", wr_id, e.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.RecordNotFound.WithError(fmt.Errorf("work range not found for this employee ID (%s)", e.ID))
		}
		return lib.Error.General.InternalError.WithError(err)
	}

	if err := tx.Exec("DELETE FROM work_schedule_services WHERE work_range_id = ? AND service_id = ?", wr.ID, service_id).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}

	return nil
}

func (e *Employee) HasServices(tx *gorm.DB, services []*Service) error {
	for _, service := range services {
		if err := e.HasService(tx, service); err != nil {
			return err
		}
	}
	return nil
}

func (e *Employee) HasService(tx *gorm.DB, service *Service) error {
	if service == nil {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("service passed is nil when validating employee (%s) services", e.ID))
	}
	var count int64
	if err := tx.Raw("SELECT COUNT(*) FROM employee_services WHERE employee_id = ? AND service_id = ?", e.ID, service.ID).Scan(&count).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	if count == 0 {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("employee %s does not have service %s", e.ID, service.ID))
	}
	return nil
}

func (e *Employee) AddService(tx *gorm.DB, service *Service) error {
	if service.CompanyID != e.CompanyID {
		return lib.Error.Company.NotSame
	}

	eID := e.ID.String()
	sID := service.ID.String()

	// Check if the employee already has the service
	var count int64
	if err := tx.Raw("SELECT COUNT(*) FROM employee_services WHERE employee_id = ? AND service_id = ?", eID, sID).Scan(&count).Error; err != nil {
		return lib.Error.General.InternalError.WithError(fmt.Errorf("failed to check if employee %s has service %s: %w", eID, sID, err))
	}
	if count > 0 {
		return lib.Error.General.UpdatedError.WithError(fmt.Errorf("employee %s already has service %s", eID, sID))
	}
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

	// Check if the employee has the service
	var count int64
	if err := tx.Raw("SELECT COUNT(*) FROM employee_services WHERE employee_id = ? AND service_id = ?", eID, sID).Scan(&count).Error; err != nil {
		return lib.Error.General.InternalError.WithError(fmt.Errorf("failed to check if employee %s has service %s: %w", eID, sID, err))
	}
	if count == 0 {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("employee %s does not have service %s", eID, sID))
	}
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

func (e *Employee) HasBranch(tx *gorm.DB, branchID uuid.UUID) error {
	var count int64
	// Check if the employee exists in the branch
	if err := tx.Raw("SELECT COUNT(*) FROM employee_branches WHERE employee_id = ? AND branch_id = ?", e.ID, branchID).Scan(&count).Error; err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	if count == 0 {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("employee %s not found in branch %s", e.ID, branchID))
	}
	return nil
}

func (e *Employee) AddBranch(tx *gorm.DB, branch *Branch) error {
	if branch.CompanyID != e.CompanyID {
		return lib.Error.Company.NotSame
	}

	eID := e.ID.String()
	bID := branch.ID.String()

	// Check if the employee already exists in the branch
	var count int64
	if err := tx.Raw("SELECT COUNT(*) FROM employee_branches WHERE employee_id = ? AND branch_id = ?", eID, bID).Scan(&count).Error; err != nil {
		return lib.Error.General.InternalError.WithError(fmt.Errorf("failed to check if employee %s is already in branch %s: %w", eID, bID, err))
	}
	if count > 0 {
		return lib.Error.General.UpdatedError.WithError(fmt.Errorf("employee %s already exists in branch %s", eID, bID))
	}
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

	// Check if the employee exists in the branch
	var count int64
	if err := tx.Raw("SELECT COUNT(*) FROM employee_branches WHERE employee_id = ? AND branch_id = ?", eID, bID).Scan(&count).Error; err != nil {
		return lib.Error.General.InternalError.WithError(fmt.Errorf("failed to check if employee %s is in branch %s: %w", eID, bID, err))
	}
	if count == 0 {
		return lib.Error.General.BadRequest.WithError(fmt.Errorf("employee %s is not in branch %s", eID, bID))
	}
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
