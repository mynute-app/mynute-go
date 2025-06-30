package modelT

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	mJSON "agenda-kaki-go/core/config/db/model/json"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/lib"
	FileBytes "agenda-kaki-go/core/lib/file_bytes"
	handler "agenda-kaki-go/core/test/handlers"
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Company struct {
	Created   *model.CompanyMerged
	Owner     *Employee
	Employees []*Employee
	Branches  []*Branch
	Services  []*Service
}

// Sets the company with 4 employees (1 owner),
// 2 branches, and 2 services.
// This method is used in tests to create a company with a predefined structure.
func (c *Company) Set() error {
	if err := c.Create(200); err != nil {
		return err
	}
	cOwnerToken := c.Owner.X_Auth_Token

	for range 2 {
		service := &Service{}
		service.Company = c
		if err := service.Create(200, cOwnerToken, nil); err != nil {
			return err
		}
		if err := service.GetById(200, cOwnerToken, nil); err != nil {
			return err
		}
		c.Services = append(c.Services, service)
	}

	for range 1 {
		branch := &Branch{}
		branch.Company = c
		if err := branch.Create(200, c.Owner.X_Auth_Token, nil); err != nil {
			return err
		}
		if err := branch.GetById(200, c.Owner.X_Auth_Token, nil); err != nil {
			return err
		}
		for _, service := range c.Services {
			if err := branch.AddService(200, service, c.Owner.X_Auth_Token, nil); err != nil {
				return fmt.Errorf("failed to assign service %s to branch %s: %v", service.Created.Name, branch.Created.Name, err)
			}
		}
		branchWorkSchedule := GetExampleBranchWorkSchedule(branch.Created.ID, c.Services)
		if err := branch.CreateWorkSchedule(200, branchWorkSchedule, c.Owner.X_Auth_Token, nil); err != nil {
			return fmt.Errorf("failed to create work schedule for branch %s: %v", branch.Created.Name, err)
		}
		c.Branches = append(c.Branches, branch)
	}

	for range 3 {
		employee := &Employee{}
		employee.Company = c
		if err := employee.Create(200, &cOwnerToken, nil); err != nil {
			return err
		}
		if err := employee.VerifyEmail(200, nil); err != nil {
			return err
		}
		if err := employee.Login(200, nil); err != nil {
			return err
		}
		if err := employee.GetById(200, nil, nil); err != nil {
			return err
		}
		for _, service := range c.Services {
			if err := employee.AddService(200, service, &c.Owner.X_Auth_Token, nil); err != nil {
				return fmt.Errorf("failed to assign employee %s to service %s: %v", employee.Created.Email, service.Created.Name, err)
			}
		}
		for _, branch := range c.Branches {
			if err := employee.AddBranch(200, branch, &c.Owner.X_Auth_Token, nil); err != nil {
				return fmt.Errorf("failed to assign employee %s to branch %s: %v", employee.Created.Email, branch.Created.Name, err)
			}
		}
		employeeID := employee.Created.ID
		branchID := c.Branches[0].Created.ID
		employeeWorkSchedule := GetExampleEmployeeWorkSchedule(employeeID, branchID, c.Services)
		if err := employee.CreateWorkSchedule(200, employeeWorkSchedule, nil, nil); err != nil {
			return fmt.Errorf("failed to create work schedule for employee %s: %v", employee.Created.Email, err)
		}
		c.Employees = append(c.Employees, employee)
	}

	if err := c.GetById(200, c.Owner.X_Auth_Token, nil); err != nil {
		return err
	}

	if err := c.UploadImages(200, map[string][]byte{
		"logo": FileBytes.PNG_FILE_1,
	}, c.Owner.X_Auth_Token, nil); err != nil {
		return err
	}

	return nil
}

// --- Randomized Company Setup Method ---

// CreateCompanyRandomly replaces the static Set method.
// It creates the Company, owner, and then generates the specified number
// of employees, branches, and services, linking them randomly.
func (c *Company) CreateCompanyRandomly(numEmployees, numBranches, numServices int) error {
	if os.Getenv("APP_ENV") == "prod" { // Make sure it never runs in production
		panic("Cannot run tests in production environment. Set APP_ENV to 'test' or 'dev'")
	}

	log.Printf("Setting up randomized Company with %d employees, %d branches, %d services", numEmployees, numBranches, numServices)

	if c.Created != nil && c.Created.ID != uuid.Nil {
		return fmt.Errorf("this function should not be called if the company has already been created")
	}

	//  --- Creating company ---

	if err := c.Create(200); err != nil {
		return err
	}

	// --- Generating Employees, Branches, and Services ---

	if err := c.GenerateEmployees(numEmployees); err != nil {
		return err
	}

	if err := c.GenerateBranches(numBranches); err != nil {
		return err
	}

	if err := c.GenerateServices(numServices); err != nil {
		return err
	}

	if err := c.GetById(200, c.Owner.X_Auth_Token, nil); err != nil {
		return err
	}

	// --- Random Relationship Assignments ---

	if err := c.RandomlyAssignServicesToEmployees(); err != nil {
		return err
	}

	if err := c.RandomlyAssignServicesToBranches(); err != nil {
		return err
	}

	if err := c.RandomlyAssignEmployeesToBranches(); err != nil {
		return err
	}

	if err := c.RandomlyAssignWorkScheduleToEmployees(); err != nil {
		return err
	}

	log.Println("Randomized Company setup completed")
	return nil
}

// AddRandomizedEntitiesToCompany adds randomized employees, branches, and services to an existing company.
// It assumes the company has already been created and has a valid owner.
// This function is used to extend an existing company with additional randomized entities.
func (c *Company) AddRandomizedEntitiesToCompany(numEmployees, numBranches, numServices int) error {
	if os.Getenv("APP_ENV") == "prod" { // Make sure it never runs in production
		panic("Cannot run tests in production environment. Set APP_ENV to 'test' or 'dev'")
	}

	if c.Created == nil || c.Created.ID == uuid.Nil {
		return fmt.Errorf("this function should not be called if the company has not been created yet")
	}

	if c.Owner == nil || c.Owner.X_Auth_Token == "" {
		return fmt.Errorf("company owner is not set or has no Auth Token")
	}

	if err := c.GetById(200, c.Owner.X_Auth_Token, nil); err != nil {
		return fmt.Errorf("it looks like the company doesn't exist: %v", err)
	}

	// --- Generating Employees, Branches, and Services ---

	if err := c.GenerateEmployees(numEmployees); err != nil {
		return err
	}

	if err := c.GenerateBranches(numBranches); err != nil {
		return err
	}

	if err := c.GenerateServices(numServices); err != nil {
		return err
	}

	if err := c.GetById(200, c.Owner.X_Auth_Token, nil); err != nil {
		return err
	}

	// --- Random Relationship Assignments ---

	if err := c.RandomlyAssignServicesToEmployees(); err != nil {
		return err
	}

	if err := c.RandomlyAssignServicesToBranches(); err != nil {
		return err
	}

	if err := c.RandomlyAssignEmployeesToBranches(); err != nil {
		return err
	}

	if err := c.RandomlyAssignWorkScheduleToEmployees(); err != nil {
		return err
	}

	log.Println("Randomized Company setup completed")
	return nil
}

// --- Generation Functions ---

// GenerateEmployees creates n *additional* employees (owner already exists).
func (c *Company) GenerateEmployees(n int) error {
	log.Printf("Generating %d employees for Company ID %s", n, c.Created.ID)

	if n <= 0 {
		return fmt.Errorf("employee generation failed: n must be greater than 0")
	} else if c.Created.ID == uuid.Nil {
		return fmt.Errorf("employee generation failed: Company ID is nil")
	}

	initialEmployeeCount := len(c.Employees)
	createdCount := 0

	for i := range n {
		employee := &Employee{Company: c}
		if err := employee.Create(200, &c.Owner.X_Auth_Token, nil); err != nil {
			return err
		}

		if employee.Created.ID == uuid.Nil {
			return fmt.Errorf("failed to create employee %d/%d or retrieve ID", i+1, n)
		}

		if err := employee.VerifyEmail(200, nil); err != nil {
			return err
		}
		if err := employee.Login(200, nil); err != nil {
			return err
		}

		if employee.X_Auth_Token == "" {
			return fmt.Errorf("failed to login employee %d/%d or retrieve Auth Token", i+1, n)
		}
		createdCount++
		c.Employees = append(c.Employees, employee)
	}
	if createdCount != n {
		return fmt.Errorf("tried to create %d employees, but only %d succeeded", n, createdCount)
	}
	if len(c.Employees) != initialEmployeeCount+createdCount {
		return fmt.Errorf("company employee slice length (%d) does not match expected count (%d)", len(c.Employees), initialEmployeeCount+createdCount)
	}
	return nil
}

// GenerateBranches creates n branches for the Company.
func (c *Company) GenerateBranches(n int) error {
	log.Printf("Generating %d branches for Company ID %s", n, c.Created.ID)

	if n <= 0 {
		return fmt.Errorf("branch generation failed: n must be greater than 0")
	} else if c.Created.ID == uuid.Nil {
		return fmt.Errorf("branch generation failed: Company ID is nil")
	}

	initialBranchCount := len(c.Branches)
	createdCount := 0

	for i := range n {
		branch := &Branch{Company: c}
		if err := branch.Create(200, c.Owner.X_Auth_Token, nil); err != nil {
			return fmt.Errorf("failed to create branch %d/%d: %v", i+1, n, err)
		}

		if branch.Created.ID == uuid.Nil {
			return fmt.Errorf("failed to create branch %d/%d or retrieve ID", i+1, n)
		}
		createdCount++
		c.Branches = append(c.Branches, branch)
	}
	if createdCount != n {
		return fmt.Errorf("tried to create %d branches, but only %d succeeded", n, createdCount)
	}
	if len(c.Branches) != initialBranchCount+createdCount {
		return fmt.Errorf("company branch slice length (%d) does not match expected count (%d)", len(c.Branches), initialBranchCount+createdCount)
	}
	return nil
}

// GenerateServices creates n services for the Company.
func (c *Company) GenerateServices(n int) error {
	log.Printf("Generating %d services for Company ID %s", n, c.Created.ID)

	if n <= 0 {
		return fmt.Errorf("service generation failed: n must be greater than 0, Company ID must not be nil, and Auth Token must not be empty")
	} else if c.Created.ID == uuid.Nil {
		return fmt.Errorf("service generation failed: Company ID is nil")
	}

	initialServiceCount := len(c.Services)
	createdCount := 0

	for i := range n {
		service := &Service{Company: c}
		if err := service.Create(200, c.Owner.X_Auth_Token, nil); err != nil {
			return fmt.Errorf("failed to create service %d/%d: %v", i+1, n, err)
		}

		if service.Created.ID == uuid.Nil {
			return fmt.Errorf("failed to create service %d/%d or retrieve ID", i+1, n)
		}
		createdCount++
		c.Services = append(c.Services, service)
	}
	if createdCount != n {
		return fmt.Errorf("tried to create %d services, but only %d succeeded", n, createdCount)
	}
	if len(c.Services) != initialServiceCount+createdCount {
		return fmt.Errorf("Company service slice length (%d) does not match expected count (%d)", len(c.Services), initialServiceCount+createdCount)
	}
	return nil
}

// --- Random Assignment Functions ---

// RandomlyAssignEmployeesToBranches assigns each employee to 1 to N random branches.
func (c *Company) RandomlyAssignEmployeesToBranches() error {
	log.Println("Randomly assigning employees to branches...")
	if len(c.Employees) == 0 {
		log.Println("Warning: No employees to assign branches to in RandomlyAssignEmployeesToBranches.")
		return nil
	}
	if len(c.Branches) == 0 {
		log.Println("Warning: No branches available in the company to assign to employees in RandomlyAssignEmployeesToBranches.")
		return nil // Não é um erro se não houver filiais, mas os funcionários não serão atribuídos.
	}

	// Filtrar por filiais válidas (ID não nulo) antecipadamente
	validBranches := []*Branch{}
	for _, b := range c.Branches {
		if b.Created.ID != uuid.Nil {
			validBranches = append(validBranches, b)
		}
	}

	if len(validBranches) == 0 {
		// Se, após a filtragem, não houver filiais válidas, isso é um problema se esperávamos atribuir.
		return fmt.Errorf("random assignment error: no valid branches found in the company (total branches: %d) to assign to employees", len(c.Branches))
	}

	for i, employee := range c.Employees {
		if employee.Created.ID == uuid.Nil {
			log.Printf("Skipping branch assignment for employee (index %d, email: %s): Employee ID is nil\n", i, employee.Created.Email)
			continue // Pular este funcionário se o ID for nulo
		}

		// Determinar quantas filiais únicas atribuir a este funcionário.
		// Atribuir pelo menos 1, até todas as filiais válidas disponíveis, limitado a um máximo razoável (ex: 10).
		maxAssignmentsForThisEmployee := len(validBranches)
		// Pode-se adicionar um teto, ex: if maxAssignmentsForThisEmployee > 5 { maxAssignmentsForThisEmployee = 5 }

		numBranchesToAssign := 1
		if maxAssignmentsForThisEmployee > 1 {
			numBranchesToAssign = rand.Intn(maxAssignmentsForThisEmployee) + 1 // Gera número de 1 a maxAssignmentsForThisEmployee
		}
		// Garante que não tentaremos atribuir mais filiais do que as disponíveis.
		if numBranchesToAssign > len(validBranches) {
			numBranchesToAssign = len(validBranches)
		}

		// Criar uma lista de índices para validBranches [0, 1, ..., len(validBranches)-1]
		shuffledIndices := make([]int, len(validBranches))
		for j := range shuffledIndices {
			shuffledIndices[j] = j
		}
		rand.Shuffle(len(shuffledIndices), func(k, l int) {
			shuffledIndices[k], shuffledIndices[l] = shuffledIndices[l], shuffledIndices[k]
		})

		assignedCount := 0
		for k := 0; k < numBranchesToAssign && k < len(shuffledIndices); k++ {
			actualBranchIndexInValidBranches := shuffledIndices[k]
			branchToAssign := validBranches[actualBranchIndexInValidBranches]

			// Usar token do proprietário para privilégio ao atribuir funcionários a filiais
			if err := employee.AddBranch(200, branchToAssign, &c.Owner.X_Auth_Token, nil); err != nil {
				return fmt.Errorf("failed to assign employee %d (%s, ID: %s) to branch (Name: %s, ID: %s): %v",
					i, employee.Created.Email, employee.Created.ID,
					branchToAssign.Created.Name, branchToAssign.Created.ID, err)
			}
			assignedCount++
		}
		if assignedCount == 0 && numBranchesToAssign > 0 {
			return fmt.Errorf("employee %s was intended to be assigned %d branches, but 0 were assigned. Valid branches: %d. This might be due to previous errors or lack of valid branches", employee.Created.Email, numBranchesToAssign, len(validBranches))
		}
	}
	return nil
}

// RandomlyAssignServicesToEmployees assigns each employee 1 to N random services.
func (c *Company) RandomlyAssignServicesToEmployees() error {
	log.Println("Randomly assigning services to employees...")
	if len(c.Employees) == 0 {
		log.Println("Warning: No employees to assign services to in RandomlyAssignServicesToEmployees.")
		return nil
	}
	if len(c.Services) == 0 {
		log.Println("Warning: No services available in the company to assign to employees in RandomlyAssignServicesToEmployees.")
		return nil
	}

	validServices := []*Service{}
	for _, s := range c.Services {
		if s.Created.ID != uuid.Nil {
			validServices = append(validServices, s)
		}
	}

	if len(validServices) == 0 {
		return fmt.Errorf("random assignment error: no valid services found in the company (total services: %d) to assign to employees", len(c.Services))
	}

	for i, employee := range c.Employees {
		if employee.Created.ID == uuid.Nil || employee.X_Auth_Token == "" {
			log.Printf("Skipping service assignment for employee (index %d, email: %s): Employee ID nil or not logged in\n", i, employee.Created.Email)
			continue
		}

		maxAssignmentsForThisEmployee := len(validServices)
		// Pode-se adicionar um teto aqui também, ex: if maxAssignmentsForThisEmployee > 10 { maxAssignmentsForThisEmployee = 10 }

		numServicesToAssign := 1
		if maxAssignmentsForThisEmployee > 1 {
			numServicesToAssign = rand.Intn(maxAssignmentsForThisEmployee) + 1
		}
		if numServicesToAssign > len(validServices) {
			numServicesToAssign = len(validServices)
		}

		shuffledIndices := make([]int, len(validServices))
		for j := range shuffledIndices {
			shuffledIndices[j] = j
		}
		rand.Shuffle(len(shuffledIndices), func(k, l int) {
			shuffledIndices[k], shuffledIndices[l] = shuffledIndices[l], shuffledIndices[k]
		})

		assignedCount := 0
		for k := 0; k < numServicesToAssign && k < len(shuffledIndices); k++ {
			actualServiceIndexInValidServices := shuffledIndices[k]
			serviceToAssign := validServices[actualServiceIndexInValidServices]

			// Usar Employee.AddService, o token do funcionário é usado (passando nil para X_Auth_Token em AddService)
			if err := employee.AddService(200, serviceToAssign, nil, nil); err != nil {
				return fmt.Errorf("failed to assign service (Name: %s, ID: %s) to employee %d (%s, ID: %s): %v",
					serviceToAssign.Created.Name, serviceToAssign.Created.ID,
					i, employee.Created.Email, employee.Created.ID, err)
			}
			assignedCount++
		}
		if assignedCount == 0 && numServicesToAssign > 0 {
			log.Printf("Warning: Employee %s was intended to be assigned %d services, but 0 were assigned. Valid services: %d.\n", employee.Created.Email, numServicesToAssign, len(validServices))
		}
	}
	return nil
}

// RandomlyAssignServicesToBranches assigns each branch 1 to N random services.
func (c *Company) RandomlyAssignServicesToBranches() error {
	if len(c.Branches) == 0 {
		log.Println("Warning: No branches to assign services to in RandomlyAssignServicesToBranches.")
		return nil
	}
	if len(c.Services) == 0 {
		log.Println("Warning: No services available in the company to assign to branches in RandomlyAssignServicesToBranches.")
		return nil
	}

	validBranches := []*Branch{}
	for _, b := range c.Branches {
		if b.Created.ID != uuid.Nil {
			validBranches = append(validBranches, b)
		}
	}
	if len(validBranches) == 0 {
		return fmt.Errorf("random assignment error: no valid branches found (total: %d)", len(c.Branches))
	}

	validServices := []*Service{}
	for _, s := range c.Services {
		if s.Created.ID != uuid.Nil {
			validServices = append(validServices, s)
		}
	}
	if len(validServices) == 0 {
		return fmt.Errorf("random assignment error: no valid services found (total: %d)", len(c.Services))
	}

	for i, branch := range validBranches { // Iterar sobre validBranches
		// Não precisa checar branch.Created.ID != uuid.Nil aqui, pois já foi filtrado.

		maxAssignmentsForThisBranch := len(validServices)
		// Pode-se adicionar um teto aqui, ex: if maxAssignmentsForThisBranch > 15 { maxAssignmentsForThisBranch = 15 }

		numServicesToAssign := 1
		if maxAssignmentsForThisBranch > 1 {
			numServicesToAssign = rand.Intn(maxAssignmentsForThisBranch) + 1
		}
		if numServicesToAssign > len(validServices) {
			numServicesToAssign = len(validServices)
		}

		shuffledIndices := make([]int, len(validServices))
		for j := range shuffledIndices {
			shuffledIndices[j] = j
		}
		rand.Shuffle(len(shuffledIndices), func(k, l int) {
			shuffledIndices[k], shuffledIndices[l] = shuffledIndices[l], shuffledIndices[k]
		})

		assignedCount := 0
		for k := 0; k < numServicesToAssign && k < len(shuffledIndices); k++ {
			actualServiceIndexInValidServices := shuffledIndices[k]
			serviceToAssign := validServices[actualServiceIndexInValidServices]

			if err := branch.AddService(200, serviceToAssign, c.Owner.X_Auth_Token, nil); err != nil {
				return fmt.Errorf("failed to assign service (Name: %s, ID: %s) to branch %d (%s, ID: %s): %v",
					serviceToAssign.Created.Name, serviceToAssign.Created.ID,
					i, branch.Created.Name, branch.Created.ID, err)
			}
			assignedCount++
		}
		if assignedCount == 0 && numServicesToAssign > 0 {
			return fmt.Errorf("branch %s was intended to be assigned %d services, but 0 were assigned. Valid services: %d", branch.Created.Name, numServicesToAssign, len(validServices))
		}
	}
	return nil
}

// --- Work Schedule Assignment ---

// RandomlyAssignWorkSchedules assigns a generated work schedule to each employee.
func (c *Company) RandomlyAssignWorkScheduleToEmployees() error {
	log.Println("Randomly assigning work schedules to employees...")
	if len(c.Employees) == 0 {
		return fmt.Errorf("no employees to assign work schedules")
	} else if len(c.Branches) == 0 {
		return fmt.Errorf("no branches to reference for work schedule assignment")
	} else if len(c.Services) == 0 {
		return fmt.Errorf("no services to reference for work schedule assignment")
	}

	for i, employee := range c.Employees {
		if employee.Created.ID == uuid.Nil { // No need to check Auth_token here, Update uses Company token.
			return fmt.Errorf("skipping schedule assignment for employee %d (%s): ID nil", i, employee.Created.Email)
		}

		if len(employee.Created.Branches) == 0 {
			log.Printf("Skipping work schedule for employee %d (%s): Not assigned to any branches.\n", i, employee.Created.Email)
			continue // Pula para o próximo funcionário
		}
		if len(employee.Created.Services) == 0 {
			log.Printf("Skipping work schedule for employee %d (%s): Does not offer any services.\n", i, employee.Created.Email)
			continue // Pula para o próximo funcionário
		}

		validEmployeeBranches := []*Branch{}
		for _, mb := range employee.Created.Branches {
			for _, mbT := range employee.Branches {
				if mbT.Created.ID == mb.ID {
					validEmployeeBranches = append(validEmployeeBranches, mbT)
				}
			}
		}

		// Make sure Employee.Branches matches Employee.Created.Branches
		if len(validEmployeeBranches) != len(employee.Created.Branches) {
			return fmt.Errorf("employee %d (%s) has %d branches, but only %d valid branches found", i, employee.Created.Email, len(employee.Created.Branches), len(validEmployeeBranches))
		} else if len(validEmployeeBranches) != len(employee.Branches) {
			return fmt.Errorf("employee %d (%s) has %d branches, but only %d valid branches found", i, employee.Created.Email, len(employee.Branches), len(validEmployeeBranches))
		}

		EmployeeWorkSchedule := GenerateRandomWorkRanges(employee.Branches, employee)

		scheduleCreationAttempts := 1

		for {
			if len(EmployeeWorkSchedule.WorkRanges) > 0 {
				break
			} else if scheduleCreationAttempts > 5 {
				return fmt.Errorf("failed to generate a valid work schedule for employee %d (%s) after %d attempts", i, employee.Created.Email, scheduleCreationAttempts)
			}
			EmployeeWorkSchedule = GenerateRandomWorkRanges(validEmployeeBranches, employee)
			scheduleCreationAttempts++
		}

		// Call Employee.Update using owner's token (c.X_Auth_Token is implicitly used in helper via employee.Company.X_Auth_Token)
		if err := employee.CreateWorkSchedule(200, EmployeeWorkSchedule, nil, nil); err != nil {
			return fmt.Errorf("failed to create work schedule for employee %d (%s) via API: %v", i, employee.Created.Email, err)
		}
	}
	return nil
}

func GenerateRandomWorkRanges(validBranches []*Branch, employee *Employee) DTO.CreateEmployeeWorkSchedule {
	var EmployeeWorkSchedule DTO.CreateEmployeeWorkSchedule

	// Range from 0 to 6 days (Sunday to Saturday)
	for day := range 7 {
		wr := generateWorkRangesForDay(validBranches, employee, model.Weekday(day), 0.9)
		EmployeeWorkSchedule.WorkRanges = append(EmployeeWorkSchedule.WorkRanges, wr...)
	}

	return EmployeeWorkSchedule
}

func generateWorkRangesForDay(validBranches []*Branch, employee *Employee, weekday model.Weekday, workProbability float32) []DTO.CreateEmployeeWorkRange {
	var ranges []DTO.CreateEmployeeWorkRange

	if rand.Float32() > workProbability || len(validBranches) == 0 {
		return ranges
	}

	numRanges := 1 + rand.Intn(2)
	lastEndTimeStr := "00:00"

	for r := 0; r < numRanges; r++ {
		branch := validBranches[rand.Intn(len(validBranches))]
		startHourLower := 7

		if r > 0 {
			var hourPart, minutePart int
			fmt.Sscanf(lastEndTimeStr, "%02d:%02d", &hourPart, &minutePart)
			startHourLower = hourPart
			if minutePart > 0 {
				startHourLower++
			}
			startHourLower++
			if startHourLower < 13 && hourPart >= 12 {
				startHourLower = 13
			}
		}
		if startHourLower > 19 {
			continue
		}

		startHour := startHourLower + rand.Intn(2)
		startTimeStr := fmt.Sprintf("%02d:%02d", startHour, 0)

		if r > 0 && startTimeStr <= lastEndTimeStr {
			startHour++
			if startHour > 20 {
				continue
			}
			startTimeStr = fmt.Sprintf("%02d:%02d", startHour, 0)
		}

		durationHours := 2 + rand.Intn(4)
		endHour := min(startHour+durationHours, 22)
		endTimeStr := fmt.Sprintf("%02d:%02d", endHour, 0)

		if endTimeStr <= startTimeStr {
			endHour++
			if endHour > 22 {
				endTimeStr = "23:30"
			} else {
				endTimeStr = fmt.Sprintf("%02d:%02d", endHour, 0)
			}
		}

		loc, err := branch.Created.GetTimeZone()
		if err != nil {
			panic(err)
		}
		startTime := time.Date(1, 1, 1, startHour, 0, 0, 0, loc).UTC()
		endTime := time.Date(1, 1, 1, endHour, 0, 0, 0, loc).UTC()
		finalStartTimeStr := startTime.Format("15:04")
		finalEndTimeStr := endTime.Format("15:04")

		employeeServices := map[uuid.UUID]bool{}
		for _, s := range employee.Created.Services {
			employeeServices[s.ID] = true
		}
		var commonServices []DTO.ServiceID
		for _, s := range branch.Services {
			if _, ok := employeeServices[s.Created.ID]; ok {
				commonServices = append(commonServices, DTO.ServiceID{ID: s.Created.ID})
			}
		}
		if len(commonServices) == 0 {
			continue
		}

		ranges = append(ranges, DTO.CreateEmployeeWorkRange{
			Weekday:    uint8(weekday),
			StartTime:  finalStartTimeStr,
			EndTime:    finalEndTimeStr,
			TimeZone:   branch.Created.TimeZone,
			EmployeeID: employee.Created.ID,
			BranchID:   branch.Created.ID,
			Services:   commonServices,
		})
		lastEndTimeStr = endTimeStr
	}

	return ranges
}

// @deprecated
// GenerateRandomModelWorkSchedule creates a *mJSON.EmployeeWorkSchedule* struct
// func GenerateRandomModelWorkSchedule(validBranches []*Branch, employee *Employee) mJSON.EmployeeWorkSchedule {
// 	schedule := mJSON.EmployeeWorkSchedule{}

// 	randomTimeStringHHMM := func(minHour, maxHour int) string {
// 		hour := min(max(minHour+rand.Intn(maxHour-minHour+1), 6), 22) // Ensure hour is between 6 and 22
// 		minute := rand.Intn(4) * 15
// 		return fmt.Sprintf("%02d:%02d", hour, minute)
// 	}

// 	schedule.Monday = generateRangesForDayModel(validBranches, employee, randomTimeStringHHMM, 0.9)
// 	schedule.Tuesday = generateRangesForDayModel(validBranches, employee, randomTimeStringHHMM, 0.9)
// 	schedule.Wednesday = generateRangesForDayModel(validBranches, employee, randomTimeStringHHMM, 0.9)
// 	schedule.Thursday = generateRangesForDayModel(validBranches, employee, randomTimeStringHHMM, 0.9)
// 	schedule.Friday = generateRangesForDayModel(validBranches, employee, randomTimeStringHHMM, 0.9)
// 	schedule.Saturday = generateRangesForDayModel(validBranches, employee, randomTimeStringHHMM, 0.4)
// 	schedule.Sunday = generateRangesForDayModel(validBranches, employee, randomTimeStringHHMM, 0.1)

// 	return schedule
// }

// @deprecated
// Helper for GenerateRandomModelWorkSchedule, returns []mJSON.WorkRange
// func generateRangesForDayModel(validBranches []*Branch, employee *Employee, randomTime func(int, int) string, workProbability float32) []mJSON.WorkRange {
// 	// Use global rand.Float32() - auto-seeded
// 	if rand.Float32() > workProbability || len(validBranches) == 0 {
// 		return []mJSON.WorkRange{}
// 	}

// 	ranges := []mJSON.WorkRange{}
// 	numRanges := 1 + rand.Intn(2) // Use global rand.Intn()

// 	lastEndTimeStr := "00:00"

// 	for r := range numRanges {
// 		// Use global rand.Intn()
// 		targetBranchHelper := validBranches[rand.Intn(len(validBranches))]

// 		startHourLower := 7
// 		if r > 0 {
// 			hourPart := 0
// 			minutePart := 0
// 			fmt.Sscanf(lastEndTimeStr, "%02d:%02d", &hourPart, &minutePart) // Ensure Sscanf matches HH:MM
// 			startHourLower = hourPart
// 			if minutePart > 0 {
// 				startHourLower++
// 			}
// 			startHourLower++ // Buffer hour
// 			if startHourLower < 13 && hourPart >= 12 {
// 				startHourLower = 13
// 			}
// 		}
// 		if startHourLower > 19 {
// 			continue
// 		}

// 		// Use global rand.Intn()
// 		startHour := startHourLower + rand.Intn(2)
// 		startTime := randomTime(startHour, startHour)

// 		if r > 0 && startTime <= lastEndTimeStr {
// 			startHour++
// 			if startHour > 20 {
// 				continue
// 			}
// 			startTime = randomTime(startHour, startHour)
// 		}

// 		// Use global rand.Intn()
// 		durationHours := 2 + rand.Intn(4)
// 		endHour := min(startHour+durationHours, 22)
// 		endTime := randomTime(endHour, endHour)

// 		if endTime <= startTime {
// 			if endHour < 22 {
// 				endHour++
// 				endTime = randomTime(endHour, endHour)
// 			} else {
// 				endTime = "23:00"
// 			}
// 			if endTime <= startTime {
// 				endTime = "23:30"
// 			}
// 		}

// 		employeeServices := []uuid.UUID{}
// 		for _, svc := range employee.Created.Services {
// 			if svc.ID != uuid.Nil {
// 				employeeServices = append(employeeServices, svc.ID)
// 			}
// 		}

// 		branchServices := []uuid.UUID{}
// 		for _, svc := range targetBranchHelper.Services {
// 			if svc.Created.ID != uuid.Nil {
// 				branchServices = append(branchServices, svc.Created.ID)
// 			}
// 		}

// 		commonServices := intersectUUIDs(employeeServices, branchServices)

// 		if len(commonServices) == 0 {
// 			continue // pula esse range se não houver serviços em comum
// 		}

// 		ranges = append(ranges, mJSON.WorkRange{
// 			Start:    startTime,
// 			End:      endTime,
// 			BranchID: targetBranchHelper.Created.ID,
// 			Services: commonServices,
// 		})

// 		lastEndTimeStr = endTime
// 	}
// 	return ranges
// }

// Creates the company with randomized data
// verifies the owner,
// logs in the owner,
// and sets the owner as the company owner.
func (c *Company) Create(status int) error {
	if c == nil {
		return fmt.Errorf("company receiver is nil")
	}
	ownerPswd := "1SecurePswd!"
	if err := handler.NewHttpClient().
		Method("POST").
		URL("/Company").
		ExpectedStatus(status).
		Send(DTO.CreateCompany{
			LegalName:      lib.GenerateRandomName("Company Legal Name"),
			TradeName:      lib.GenerateRandomName("Company Trade Name"),
			TaxID:          lib.GenerateRandomStrNumber(14),
			OwnerName:      lib.GenerateRandomName("Owner Name"),
			OwnerSurname:   lib.GenerateRandomName("Owner Surname"),
			OwnerEmail:     lib.GenerateRandomEmail("owner"),
			OwnerPhone:     lib.GenerateRandomPhoneNumber(),
			OwnerPassword:  ownerPswd,
			StartSubdomain: strings.ToLower(lib.GenerateRandomString(12)),
		}).
		ParseResponse(&c.Created).
		Error; err != nil {
		return fmt.Errorf("failed to create company: %w", err)
	}
	owner := c.Created.Employees[0]
	owner.Password = ownerPswd
	c.Owner = &Employee{
		Company: c,
		Created: &owner,
	}
	if err := c.Owner.VerifyEmail(200, nil); err != nil {
		return fmt.Errorf("failed to verify owner email: %w", err)
	}
	if err := c.Owner.Login(200, nil); err != nil {
		return fmt.Errorf("failed to login owner: %w", err)
	}
	c.Employees = append(c.Employees, c.Owner)
	if err := c.Owner.GetById(200, nil, nil); err != nil {
		return err
	}
	if c.Created.ID == uuid.Nil {
		return fmt.Errorf("failed to create Company: ID is nil")
	} else if c.Owner == nil {
		return fmt.Errorf("failed to create Company: Owner is nil")
	} else if c.Owner.Created.ID == uuid.Nil {
		return fmt.Errorf("failed to create Company: Owner ID is nil")
	} else if c.Owner.X_Auth_Token == "" {
		return fmt.Errorf("failed to create Company: Owner Auth Token is empty")
	}
	return nil
}

func (c *Company) GetByName(status int, x_auth_token string, x_company_id *string) error {
	companyIDStr := c.Created.ID.String()
	cID, err := get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	if err := handler.NewHttpClient().
		Method("GET").
		URL(fmt.Sprintf("/Company/name/%s", c.Created.LegalName)).
		ExpectedStatus(status).
		Header(namespace.HeadersKey.Auth, x_auth_token).
		Header(namespace.HeadersKey.Company, cID).
		Send(nil).
		ParseResponse(&c.Created).
		Error; err != nil {
		return fmt.Errorf("failed to get company by name: %w", err)
	}

	return nil
}

func (c *Company) GetBySubdomain(status int) error {
	if err := handler.NewHttpClient().
		Method("GET").
		URL(fmt.Sprintf("/Company/subdomain/%s", c.Created.Subdomains[0].Name)).
		ExpectedStatus(status).
		Send(nil).
		ParseResponse(&c.Created).
		Error; err != nil {
		return fmt.Errorf("failed to get company by subdomain: %w", err)
	}

	return nil
}

func (c *Company) GetById(status int, x_auth_token string, x_company_id *string) error {
	companyIDStr := c.Created.ID.String()
	cID, err := get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	if err := handler.NewHttpClient().
		Method("GET").
		URL(fmt.Sprintf("/Company/%s", c.Created.ID.String())).
		ExpectedStatus(status).
		Header(namespace.HeadersKey.Auth, x_auth_token).
		Header(namespace.HeadersKey.Company, cID).
		Send(nil).
		ParseResponse(&c.Created).
		Error; err != nil {
		return fmt.Errorf("failed to get company by id: %w", err)
	}

	return nil
}

func (c *Company) Update(status int, changes map[string]any, x_auth_token string, x_company_id *string) error {
	var companyIDStr = c.Created.ID.String()
	cID, err := get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	if err := handler.NewHttpClient().
		Method("PATCH").
		URL(fmt.Sprintf("/Company/%s", c.Created.ID.String())).
		ExpectedStatus(status).
		Header(namespace.HeadersKey.Auth, x_auth_token).
		Header(namespace.HeadersKey.Company, cID).
		Send(changes).
		ParseResponse(&c.Created).Error; err != nil {
		return fmt.Errorf("failed to update company: %w", err)
	}
	if status > 200 && status < 300 {
		if err := ValidateUpdateChanges("Company", c.Created, changes); err != nil {
			return err
		}
	}

	return nil
}

func (c *Company) Delete(status int, x_auth_token string, x_company_id *string) error {
	var companyIDStr = c.Created.ID.String()
	cID, err := get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	if err := handler.NewHttpClient().
		Method("DELETE").
		URL(fmt.Sprintf("/Company/%s", c.Created.ID.String())).
		ExpectedStatus(status).
		Header(namespace.HeadersKey.Auth, x_auth_token).
		Header(namespace.HeadersKey.Company, cID).
		Send(nil).
		Error; err != nil {
		return fmt.Errorf("failed to delete company: %w", err)
	}

	return nil
}

func (c *Company) UploadImages(status int, files map[string][]byte, x_auth_token string, x_company_id *string) error {
	var fileMap = make(handler.Files)
	for field, content := range files {
		fileMap[field] = handler.MyFile{
			Name:    field + "_" + lib.GenerateRandomString(6) + ".png",
			Content: content,
		}
	}

	companyIDStr := c.Created.ID.String()
	cID, err := get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}

	if err := handler.NewHttpClient().
		Method("PATCH").
		URL(fmt.Sprintf("/Company/%s/design/images", c.Created.ID.String())).
		ExpectedStatus(status).
		Header(namespace.HeadersKey.Auth, x_auth_token).
		Header(namespace.HeadersKey.Company, cID).
		Send(fileMap).
		ParseResponse(&c.Created.Design.Images).
		Error; err != nil {
		return fmt.Errorf("failed to upload company images: %w", err)
	}

	return nil
}

func (c *Company) DeleteImages(status int, images []string, x_auth_token string, x_company_id *string) error {
	if len(images) == 0 {
		return fmt.Errorf("no images provided to delete")
	}

	createdCompanyID := c.Created.ID.String()
	cID, err := get_x_company_id(x_company_id, &createdCompanyID)
	if err != nil {
		return fmt.Errorf("failed to get company ID for deletion: %w", err)
	}

	http := handler.NewHttpClient()

	if err := http.
		Method("DELETE").
		ExpectedStatus(status).
		Header(namespace.HeadersKey.Auth, x_auth_token).
		Header(namespace.HeadersKey.Company, cID).
		Error; err != nil {
		return fmt.Errorf("failed to prepare delete images request: %w", err)
	}

	base_url := fmt.Sprintf("/Company/%s/design/images", c.Created.ID.String())
	for _, field := range images {
		image_url := base_url + "/" + field
		http.URL(image_url)
		http.Send(nil)
		http.ParseResponse(&c.Created.Design)
		if http.Error != nil {
			return fmt.Errorf("failed to delete image %s: %w", field, http.Error)
		}
		url := c.Created.Design.Images.GetImageURL(field)
		if url != "" {
			return fmt.Errorf("image %s was not deleted successfully, expected empty URL but got %s", field, url)
		}
	}
	return nil
}

func (c *Company) ChangeColors(status int, colors mJSON.Colors, x_auth_token string, x_company_id *string) error {
	var companyIDStr = c.Created.ID.String()
	cID, err := get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	if err := handler.NewHttpClient().
		Method("PUT").
		URL(fmt.Sprintf("/Company/%s/design/colors", c.Created.ID.String())).
		ExpectedStatus(status).
		Header(namespace.HeadersKey.Auth, x_auth_token).
		Header(namespace.HeadersKey.Company, cID).
		Send(colors).
		ParseResponse(&c.Created.Design.Colors).
		Error; err != nil {
		return fmt.Errorf("failed to change colors: %w", err)
	}

	return nil
}

func (c *Company) GetImage(status int, imageURL string, compareImgBytes *[]byte) error {
	http := handler.NewHttpClient()
	http.Method("GET")
	http.URL(imageURL)
	http.ExpectedStatus(status)
	http.Send(nil)
	// Compare the response bytes with the expected image bytes
	if compareImgBytes != nil {
		var response []byte
		http.ParseResponse(&response)
		if len(response) == 0 {
			return fmt.Errorf("received empty response for image %s", imageURL)
		} else if len(response) != len(*compareImgBytes) {
			return fmt.Errorf("image size mismatch for %s: expected %d bytes, got %d bytes", imageURL, len(*compareImgBytes), len(response))
		} else if !bytes.Equal(response, *compareImgBytes) {
			return fmt.Errorf("image content mismatch for %s", imageURL)
		}
	}
	return nil
}

func intersectUUIDs(a, b []uuid.UUID) []uuid.UUID {
	set := make(map[uuid.UUID]bool)
	for _, id := range a {
		set[id] = true
	}

	var intersection []uuid.UUID
	for _, id := range b {
		if set[id] {
			intersection = append(intersection, id)
		}
	}
	return intersection
}
