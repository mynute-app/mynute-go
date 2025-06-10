package controller

import (
	"agenda-kaki-go/core/db"
	"agenda-kaki-go/core/model"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"time"
)

func GetScheduleOptions(c *fiber.Ctx) error {
	get := c.Query("get")

	companyID, err := GetCompanyIDFromContext(c)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid company context"})
	}

	var (
		branchID, employeeID, serviceID *uuid.UUID
		weekday                         *time.Weekday
		timeStr                         *string
	)

	if val := c.Query("branch_id"); val != "" {
		if id, err := uuid.Parse(val); err == nil {
			branchID = &id
		}
	}
	if val := c.Query("employee_id"); val != "" {
		if id, err := uuid.Parse(val); err == nil {
			employeeID = &id
		}
	}
	if val := c.Query("service_id"); val != "" {
		if id, err := uuid.Parse(val); err == nil {
			serviceID = &id
		}
	}
	if val := c.Query("start_time"); val != "" {
		if parsed, err := time.Parse(time.RFC3339, val); err == nil {
			d := parsed.Weekday()
			s := parsed.Format("15:04")
			weekday = &d
			timeStr = &s
		} else {
			return c.Status(400).JSON(fiber.Map{"error": "invalid start_time, must be RFC3339"})
		}
	}

	switch get {
	case "services":
		result, err := filterServices(companyID, branchID, employeeID, weekday, timeStr)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(result)

	case "employees":
		result, err := filterEmployees(companyID, branchID, serviceID, weekday, timeStr)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(result)

	case "branches":
		result, err := filterBranches(companyID, employeeID, serviceID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(result)

	default:
		return c.Status(400).JSON(fiber.Map{"error": "invalid 'get' parameter"})
	}
}

func filterServices(
	companyID uuid.UUID,
	branchID, employeeID *uuid.UUID,
	weekday *time.Weekday,
	timeStr *string,
) ([]model.Service, error) {
	conn := db.Conn
	serviceMap := make(map[uuid.UUID]model.Service)

	if branchID == nil && employeeID == nil && weekday == nil && timeStr == nil {
		var all []model.Service
		if err := conn.Where("company_id = ?", companyID).Find(&all).Error; err != nil {
			return nil, err
		}
		return all, nil
	}

	if branchID != nil && employeeID == nil && weekday == nil {
		var services []model.Service
		err := conn.Joins("JOIN branch_services bs ON bs.service_id = services.id").
			Where("bs.branch_id = ?", *branchID).Find(&services).Error
		if err != nil {
			return nil, err
		}
		return services, nil
	}

	if branchID != nil && employeeID != nil && weekday == nil {
		var emp model.Employee
		if err := conn.Preload("Services").Preload("Branches").
			First(&emp, "id = ?", *employeeID).Error; err != nil {
			return nil, err
		}
		for _, s := range emp.Services {
			for _, b := range emp.Branches {
				if b.ID == *branchID {
					serviceMap[s.ID] = *s
				}
			}
		}
	}

	if branchID != nil && weekday != nil {
		var emps []model.Employee
		err := conn.Joins("JOIN employee_branches eb ON eb.employee_id = employees.id").
			Where("eb.branch_id = ?", *branchID).
			Preload("Services").
			Find(&emps).Error
		if err != nil {
			return nil, err
		}
		for _, emp := range emps {
			ranges := emp.WorkSchedule.GetRangesForDay(*weekday)
			for _, r := range ranges {
				if r.BranchID != *branchID {
					continue
				}
				if timeStr != nil && !(r.Start <= *timeStr && *timeStr < r.End) {
					continue
				}
				for _, sID := range r.Services {
					for _, svc := range emp.Services {
						if svc.ID == sID {
							serviceMap[svc.ID] = *svc
						}
					}
				}
			}
		}
	}

	if employeeID != nil {
		var emp model.Employee
		if err := conn.Preload("Services").
			First(&emp, "id = ?", *employeeID).Error; err != nil {
			return nil, err
		}
		ranges := emp.WorkSchedule.GetAllRanges()
		if weekday != nil {
			ranges = emp.WorkSchedule.GetRangesForDay(*weekday)
		}
		for _, r := range ranges {
			if timeStr != nil && !(r.Start <= *timeStr && *timeStr < r.End) {
				continue
			}
			for _, sID := range r.Services {
				for _, svc := range emp.Services {
					if svc.ID == sID {
						serviceMap[svc.ID] = *svc
					}
				}
			}
		}
	}

	if weekday != nil && timeStr != nil && branchID == nil && employeeID == nil {
		var emps []model.Employee
		err := conn.Preload("Services").Find(&emps).Error
		if err != nil {
			return nil, err
		}
		for _, emp := range emps {
			ranges := emp.WorkSchedule.GetRangesForDay(*weekday)
			for _, r := range ranges {
				if r.Start <= *timeStr && *timeStr < r.End {
					for _, sID := range r.Services {
						for _, svc := range emp.Services {
							if svc.ID == sID {
								serviceMap[svc.ID] = *svc
							}
						}
					}
				}
			}
		}
	}

	services := make([]model.Service, 0, len(serviceMap))
	for _, s := range serviceMap {
		services = append(services, s)
	}
	return services, nil
}

func filterEmployees(
	companyID uuid.UUID,
	branchID, serviceID *uuid.UUID,
	weekday *time.Weekday,
	timeStr *string,
) ([]model.Employee, error) {
	conn := db.Conn
	employeeMap := make(map[uuid.UUID]model.Employee)

	query := conn.Model(&model.Employee{}).
		Joins("JOIN employee_services es ON es.employee_id = employees.id").
		Where("employees.company_id = ?", companyID).
		Preload("Services")

	if serviceID != nil {
		query = query.Where("es.service_id = ?", *serviceID)
	}
	if branchID != nil {
		query = query.Joins("JOIN employee_branches eb ON eb.employee_id = employees.id").
			Where("eb.branch_id = ?", *branchID)
	}

	var emps []model.Employee
	if err := query.Find(&emps).Error; err != nil {
		return nil, err
	}

	for _, emp := range emps {
		if weekday != nil && timeStr != nil {
			ranges := emp.WorkSchedule.GetRangesForDay(*weekday)
			for _, r := range ranges {
				if r.Start <= *timeStr && *timeStr < r.End {
					employeeMap[emp.ID] = emp
					break
				}
			}
		} else {
			employeeMap[emp.ID] = emp
		}
	}

	result := make([]model.Employee, 0, len(employeeMap))
	for _, emp := range employeeMap {
		result = append(result, emp)
	}
	return result, nil
}

func filterBranches(
	companyID uuid.UUID,
	employeeID, serviceID *uuid.UUID,
) ([]model.Branch, error) {
	conn := db.Conn
	query := conn.Model(&model.Branch{}).
		Where("company_id = ?", companyID).
		Preload("Services")

	if employeeID != nil {
		query = query.Joins("JOIN employee_branches eb ON eb.branch_id = branches.id").
			Where("eb.employee_id = ?", *employeeID)
	}
	if serviceID != nil {
		query = query.Joins("JOIN branch_services bs ON bs.branch_id = branches.id").
			Where("bs.service_id = ?", *serviceID)
	}

	var branches []model.Branch
	if err := query.Find(&branches).Error; err != nil {
		return nil, err
	}
	return branches, nil
}
