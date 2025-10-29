package model

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	DTO "mynute-go/core/src/config/api/dto"
	"mynute-go/core/src/config/db/model"
	mJSON "mynute-go/core/src/config/db/model/json"
	"mynute-go/core/src/config/namespace"
	"mynute-go/core/src/lib"
	FileBytes "mynute-go/core/src/lib/file_bytes"
	"mynute-go/test/src/handler"
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

	for range 3 {
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

	servicesID := make([]DTO.ServiceBase, len(c.Services))
	for i, service := range c.Services {
		servicesID[i] = DTO.ServiceBase{ID: service.Created.ID}
	}

	for range 3 {
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
		branchWorkSchedule := GetExampleBranchWorkSchedule(branch.Created.ID, servicesID)
		if err := branch.CreateWorkSchedule(200, branchWorkSchedule, c.Owner.X_Auth_Token, nil); err != nil {
			return fmt.Errorf("failed to create work schedule for branch %s: %v", branch.Created.Name, err)
		}
		if err := branch.UploadImages(200, map[string][]byte{
			"profile": FileBytes.PNG_FILE_1,
		}, c.Owner.X_Auth_Token, nil); err != nil {
			return err
		}
		c.Branches = append(c.Branches, branch)
	}

	branchID := c.Branches[0].Created.ID

	for range 3 {
		employee := &Employee{}
		employee.Company = c
		if err := employee.Create(200, &cOwnerToken, nil); err != nil {
			return err
		}
		// Use email code login for first login to verify the employee
		if err := employee.LoginWithEmailCode(200, nil); err != nil {
			return err
		}
		if err := employee.LoginWithPassword(200, nil); err != nil {
			return err
		}
		if err := employee.GetById(200, nil, nil); err != nil {
			return err
		}
		if err := employee.UploadImages(200, map[string][]byte{
			"profile": FileBytes.PNG_FILE_1,
		}, nil, nil); err != nil {
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
		employeeWorkSchedule := GetExampleEmployeeWorkSchedule(employeeID, branchID, servicesID)
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

	if err := c.RandomlyAssignWorkScheduleToBranches(); err != nil {
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

	if err := c.RandomlyAssignWorkScheduleToBranches(); err != nil {
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
		if lib.GenerateRandomIntFromRange(0, 1) == 1 {
			if err := employee.LoginWithPassword(200, nil); err != nil {
				return err
			}
		} else {
			if err := employee.LoginWithEmailCode(200, nil); err != nil {
				return err
			}
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

		EmployeeWorkSchedule, err := GenerateRandomEmployeeWorkRanges(employee.Branches, employee)
		if err != nil {
			return fmt.Errorf("failed to generate work schedule for employee %d (%s): %v", i, employee.Created.Email, err)
		}

		scheduleCreationAttempts := 1

		for {
			if len(EmployeeWorkSchedule.WorkRanges) > 0 {
				break
			} else if scheduleCreationAttempts > 50 {
				return fmt.Errorf("failed to generate a valid work schedule for employee %d (%s) after %d attempts", i, employee.Created.Email, scheduleCreationAttempts)
			}
			EmployeeWorkSchedule, err = GenerateRandomEmployeeWorkRanges(validEmployeeBranches, employee)
			if err != nil {
				return fmt.Errorf("failed to generate work schedule for employee %d (%s): %v", i, employee.Created.Email, err)
			}
			scheduleCreationAttempts++
		}

		// Call Employee.Update using owner's token (c.X_Auth_Token is implicitly used in helper via employee.Company.X_Auth_Token)
		if err := employee.CreateWorkSchedule(200, EmployeeWorkSchedule, nil, nil); err != nil {
			return fmt.Errorf("failed to create work schedule for employee %d (%s) via API: %v", i, employee.Created.Email, err)
		}
	}
	return nil
}

func (c *Company) RandomlyAssignWorkScheduleToBranches() error {
	log.Println("Randomly assigning work schedules to branches...")
	if len(c.Branches) == 0 {
		return fmt.Errorf("no branches to assign work schedules")
	}
	for i, branch := range c.Branches {
		if branch.Created.ID == uuid.Nil {
			return fmt.Errorf("skipping schedule assignment for branch %d (%s): ID nil", i, branch.Created.Name)
		}
		branchWorkSchedule, err := GenerateRandomBranchWorkRanges(branch)
		if err != nil {
			return fmt.Errorf("failed to generate work schedule for branch %d (%s): %v", i, branch.Created.Name, err)
		}
		if len(branchWorkSchedule.WorkRanges) == 0 {
			return fmt.Errorf("failed to generate a valid work schedule for branch %d (%s)", i, branch.Created.Name)
		}
		if err := branch.CreateWorkSchedule(200, branchWorkSchedule, c.Owner.X_Auth_Token, nil); err != nil {
			return fmt.Errorf("failed to create work schedule for branch %d (%s) via API: %v", i, branch.Created.Name, err)
		}
	}
	return nil
}

func GenerateRandomBranchWorkRanges(branch *Branch) (DTO.CreateBranchWorkSchedule, error) {
	var BranchWorkSchedule DTO.CreateBranchWorkSchedule

	// Range from 0 to 6 days (Sunday to Saturday)
	for day := range 7 {
		if err := generateWorkRangeForDay(&BranchWorkSchedule.WorkRanges, []*Branch{branch}, nil, model.Weekday(day), 0.8); err != nil {
			return DTO.CreateBranchWorkSchedule{}, fmt.Errorf("failed to generate work range for branch %d (%s) on day %d: %v", branch.Created.ID, branch.Created.Name, day, err)
		}
	}
	return BranchWorkSchedule, nil
}

// GenerateRandomEmployeeWorkRanges creates a randomized work schedule for an employee.
func GenerateRandomEmployeeWorkRanges(validBranches []*Branch, employee *Employee) (DTO.CreateEmployeeWorkSchedule, error) {
	var EmployeeWorkSchedule DTO.CreateEmployeeWorkSchedule

	// Range from 0 to 6 days (Sunday to Saturday)
	for day := range 7 {
		if err := generateWorkRangeForDay(&EmployeeWorkSchedule.WorkRanges, validBranches, employee, model.Weekday(day), 0.9); err != nil {
			return DTO.CreateEmployeeWorkSchedule{}, fmt.Errorf("failed to generate work range for employee %d (%s) on day %d: %v", employee.Created.ID, employee.Created.Name, day, err)
		}
	}
	return EmployeeWorkSchedule, nil
}

// generateWorkRangeForDay creates work ranges for a specific day.
//   - It randomly decides whether the employee works that day based on workProbability.
//   - If they do work, it generates 1 to 2 work ranges with random start and end times,
//     ensuring no overlaps and that they fit within business hours.
//   - It also ensures that the services offered in the work range are common between the employee and the branch.
//   - Returns a slice of CreateEmployeeWorkRange DTOs.
//   - If the employee is nil, it generates ranges without employee-specific services (for branches).
func generateWorkRangeForDay(ranges any, validBranches []*Branch, employee *Employee, weekday model.Weekday, workProbability float32) error {
	if employee == nil {
		log.Printf("Generating work range for branch on weekday %d\n", weekday)
	} else {
		log.Printf("Generating work range for employee %d (%s) on weekday %d\n", employee.Created.ID, employee.Created.Email, weekday)
	}

	if rand.Float32() > workProbability || len(validBranches) == 0 {
		return nil // Does not work this day or no valid branches
	}

	branch := validBranches[rand.Intn(len(validBranches))]

	var branchWorkRangesForDay []model.BranchWorkRange
	if employee != nil {
		// If it's an employee, we must generate a work range inside the branch's schedule
		for _, bws := range branch.Created.WorkSchedule {
			bwsWeekday := int(bws.Weekday)
			if bwsWeekday == int(weekday) {
				branchWorkRangesForDay = append(branchWorkRangesForDay, bws)
			}
		}
		if len(branchWorkRangesForDay) == 0 {
			return nil // Branch doesn't work on this day, so employee can't either.
		}
	}

	// Employees can have 1-2 ranges per day, branches can have 1-3
	numRanges := 1 + rand.Intn(2) // 1 or 2 for employees
	if employee == nil {
		numRanges = 1 + rand.Intn(3) // 1, 2, or 3 for branches
	}

	getMinutes := func() int {
		fifteenMinOpts := []int{0, 15, 30, 45}
		return fifteenMinOpts[rand.Intn(4)]
	}

	type AcceptableHours struct {
		min int
		max int
	}

	acceptableHours := AcceptableHours{
		min: 6,
		max: 23,
	}

	// For branches, we want sequential non-overlapping ranges, so lastEndHour persists across iterations
	branchLastEndHour := 0
	branchLastEndMinute := 0
	// For employees, also track sequential non-overlapping ranges
	var employeeLastEndHour int
	var employeeLastEndMinute int

	for range numRanges {
		var startHour int
		var startMinute int
		var endHour int
		var endMinute int

		if employee != nil {
			// Employee work ranges must be sequential and non-overlapping
			// Pick a branch range that can fit the employee's next sequential range
			var branchRange *model.BranchWorkRange

			for i := range branchWorkRangesForDay {
				br := &branchWorkRangesForDay[i]
				loc, _ := time.LoadLocation(br.TimeZone)
				brEnd := br.EndTime.In(loc)

				// Find a branch range that ends after the employee's last range ended
				if employeeLastEndHour == 0 || brEnd.Hour() > employeeLastEndHour ||
					(brEnd.Hour() == employeeLastEndHour && brEnd.Minute() > employeeLastEndMinute) {
					branchRange = br
					break
				}
			}

			if branchRange == nil {
				// No suitable branch range found
				continue
			}

			loc, err := time.LoadLocation(branchRange.TimeZone)
			if err != nil {
				return fmt.Errorf("failed to load timezone for branch range: %v", err)
			}

			branchStartTime := branchRange.StartTime.In(loc)
			branchEndTime := branchRange.EndTime.In(loc)

			// Generate start and end times within this branch's range, after the last employee range
			minStartHour := branchStartTime.Hour()
			maxEndHour := branchEndTime.Hour()

			// If this is the employee's second (or later) range, ensure it starts after the previous one
			if employeeLastEndHour > 0 {
				// Calculate minimum start based on last end time
				if employeeLastEndMinute > 0 {
					// If last range ended with minutes (e.g., 18:45), next can start in the next hour
					minStartHour = max(minStartHour, employeeLastEndHour+1)
				} else {
					// If last range ended on the hour (e.g., 18:00), next can start same hour or later
					minStartHour = max(minStartHour, employeeLastEndHour)
				}
			}

			if maxEndHour-minStartHour < 1 {
				continue // Not enough room in this branch range, skip
			}

			// Generate start time
			startHour = minStartHour + rand.Intn(max(1, maxEndHour-minStartHour))
			startMinute = getMinutes()

			// Ensure the start is within the branch range
			if startHour < minStartHour || (startHour == minStartHour && startMinute < branchStartTime.Minute()) {
				startHour = minStartHour
				startMinute = branchStartTime.Minute()
			}

			// Ensure this range starts after the previous employee range ended
			if employeeLastEndHour > 0 {
				if startHour < employeeLastEndHour || (startHour == employeeLastEndHour && startMinute <= employeeLastEndMinute) {
					// Adjust start time to be after the last range ended
					startHour = employeeLastEndHour
					startMinute = employeeLastEndMinute
					// Give a small gap (e.g., start at the next minute mark)
					if startMinute < 45 {
						startMinute += 15 // 15 minute gap
					} else {
						startHour++
						startMinute = 0
					}

					// After adjustment, check if there's still enough time in this branch range
					if startHour >= branchEndTime.Hour() {
						continue // No room left in this branch range
					}
				}
			}
			// Generate end time, ensuring it's after start time and within branch hours

			// Need at least 1 hour for a valid range
			if startHour >= branchEndTime.Hour() {
				continue // Cannot create a valid range
			}

			durationHours := 1 + rand.Intn(max(1, branchEndTime.Hour()-startHour))

			endHour = min(startHour+durationHours, branchEndTime.Hour())
			if endHour == startHour {
				endHour++
			}
			endMinute = getMinutes()

			// Ensure end time doesn't exceed branch end time
			if endHour > branchEndTime.Hour() || (endHour == branchEndTime.Hour() && endMinute > branchEndTime.Minute()) {
				endHour = branchEndTime.Hour()
				endMinute = branchEndTime.Minute()
			}

			if endHour >= 24 && endMinute > 0 {
				continue // Skip invalid end hours
			}

			if startHour >= endHour {
				return fmt.Errorf("for some reason the branch start hour generated (%d) is bigger than the end hour generated (%d). This should not happen", startHour, endHour)
			}

			// Update the last end time for the next employee range selection
			employeeLastEndHour = endHour
			employeeLastEndMinute = endMinute
		} else {
			// Original logic for branch work schedule generation (no employee)

			if branchLastEndHour == 0 {
				branchLastEndHour = acceptableHours.min
				branchLastEndMinute = 0
			}

			startHour = branchLastEndHour + rand.Intn(2)
			startMinute = getMinutes()

			if startHour == branchLastEndHour {
				startHour = startHour + 1
				startMinute = branchLastEndMinute
			}

			if startHour >= acceptableHours.max-1 {
				continue // Cannot create a valid range if start hour is the same as max hour
			}

			duration := 1 + rand.Intn(max(1, acceptableHours.max-startHour-1))

			endHour = min(startHour+duration, acceptableHours.max)
			if endHour == startHour {
				endHour++
			}
			endMinute = getMinutes()

			if endHour >= 24 && endMinute > 0 {
				continue // Skip invalid end hours
			}

			if startHour >= endHour {
				return fmt.Errorf("for some reason the branch start hour generated (%d) is bigger than the end hour generated (%d). This should not happen", startHour, endHour)
			}

			branchLastEndHour = endHour
			branchLastEndMinute = endMinute
		}

		log.Printf("Found startHour: %d, startMinute: %d, endHour: %d, endMinute: %d, weekday: %d, branchID: %s", startHour, startMinute, endHour, endMinute, weekday, branch.Created.ID)

		startTime := time.Date(2020, 1, 1, startHour, startMinute, 0, 0, time.UTC).UTC()
		endTime := time.Date(2020, 1, 1, endHour, endMinute, 0, 0, time.UTC).UTC()
		finalStartTimeStr := startTime.Format("15:04")
		finalEndTimeStr := endTime.Format("15:04")

		if startTime.After(endTime) || startTime.Equal(endTime) {
			continue // Skip invalid ranges
		}

		var commonServices []DTO.ServiceBase
		if employee != nil {
			if len(employee.Created.Services) == 0 {
				return fmt.Errorf("employee has no services assigned")
			}

			employeeServices := map[uuid.UUID]bool{}
			for _, s := range employee.Created.Services {
				employeeServices[s.ID] = true
			}

			for _, s := range branch.Services {
				if _, ok := employeeServices[s.Created.ID]; ok {
					commonServices = append(commonServices, DTO.ServiceBase{ID: s.Created.ID})
				}
			}
		} else {
			for _, s := range branch.Services {
				commonServices = append(commonServices, DTO.ServiceBase{ID: s.Created.ID})
			}
		}

		if len(commonServices) == 0 {
			continue
		}

		// Log generated start and end times for debugging

		switch rgs := ranges.(type) {
		case *[]DTO.CreateBranchWorkRange:
			log.Printf("Generated work range for branch (%s) on weekday %d: %s - %s\n", branch.Created.ID, weekday, finalStartTimeStr, finalEndTimeStr)
			*rgs = append(*rgs, DTO.CreateBranchWorkRange{
				Weekday:                 uint8(weekday),
				StartTime:               finalStartTimeStr,
				EndTime:                 finalEndTimeStr,
				TimeZone:                branch.Created.TimeZone,
				BranchID:                branch.Created.ID,
				BranchWorkRangeServices: DTO.BranchWorkRangeServices{Services: commonServices},
			})
		case *[]DTO.CreateEmployeeWorkRange:
			log.Printf("Generated work range for employee (%s) on weekday %d: %s - %s\n", employee.Created.ID, weekday, finalStartTimeStr, finalEndTimeStr)
			*rgs = append(*rgs, DTO.CreateEmployeeWorkRange{
				Weekday:                   uint8(weekday),
				StartTime:                 finalStartTimeStr,
				EndTime:                   finalEndTimeStr,
				TimeZone:                  branch.Created.TimeZone,
				EmployeeID:                employee.Created.ID,
				BranchID:                  branch.Created.ID,
				EmployeeWorkRangeServices: DTO.EmployeeWorkRangeServices{Services: commonServices},
			})
		}
	}
	return nil
}

// Creates the company with randomized data
// verifies the owner,
// logs in the owner,
// and sets the owner as the company owner.
func (c *Company) Create(status int) error {
	if c == nil {
		return fmt.Errorf("company receiver is nil")
	}
	ownerPswd := lib.GenerateValidPassword()
	ownerEmail := lib.GenerateRandomEmail("owner")
	if err := handler.NewHttpClient().
		Method("POST").
		URL("/company").
		ExpectedStatus(status).
		Send(DTO.CreateCompany{
			LegalName:      lib.GenerateRandomName("Company Legal Name"),
			TradeName:      lib.GenerateRandomName("Company Trade Name"),
			TaxID:          lib.GenerateRandomStrNumber(14),
			OwnerName:      lib.GenerateRandomName("Owner Name"),
			OwnerSurname:   lib.GenerateRandomName("Owner Surname"),
			OwnerEmail:     ownerEmail,
			OwnerPhone:     lib.GenerateRandomPhoneNumber(),
			OwnerPassword:  ownerPswd,
			OwnerTimeZone:  "America/Sao_Paulo",
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
	c.Owner.Created.Email = ownerEmail
	// Use email code login for first login to verify the owner
	if err := c.Owner.LoginWithEmailCode(200, nil); err != nil {
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
	cID, err := Get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	if err := handler.NewHttpClient().
		Method("GET").
		URL(fmt.Sprintf("/company/name/%s", c.Created.LegalName)).
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
		URL(fmt.Sprintf("/company/subdomain/%s", c.Created.Subdomains[0].Name)).
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
	cID, err := Get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	if err := handler.NewHttpClient().
		Method("GET").
		URL(fmt.Sprintf("/company/%s", c.Created.ID.String())).
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
	cID, err := Get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	if err := handler.NewHttpClient().
		Method("PATCH").
		URL(fmt.Sprintf("/company/%s", c.Created.ID.String())).
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
	cID, err := Get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	if err := handler.NewHttpClient().
		Method("DELETE").
		URL(fmt.Sprintf("/company/%s", c.Created.ID.String())).
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
	cID, err := Get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}

	if err := handler.NewHttpClient().
		Method("PATCH").
		URL(fmt.Sprintf("/company/%s/design/images", c.Created.ID.String())).
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

func (c *Company) DeleteImages(status int, image_types []string, x_auth_token string, x_company_id *string) error {
	if len(image_types) == 0 {
		return fmt.Errorf("no images provided to delete")
	}

	createdCompanyID := c.Created.ID.String()
	cID, err := Get_x_company_id(x_company_id, &createdCompanyID)
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

	base_url := fmt.Sprintf("/company/%s/design/images", c.Created.ID.String())
	for _, image_type := range image_types {
		image_url := base_url + "/" + image_type
		http.URL(image_url)
		http.Send(nil)
		http.ParseResponse(&c.Created.Design.Images)
		if http.Error != nil {
			return fmt.Errorf("failed to delete image %s: %w", image_type, http.Error)
		}
		url := c.Created.Design.Images.GetImageURL(image_type)
		if url != "" {
			return fmt.Errorf("image %s was not deleted successfully, expected empty URL but got %s", image_type, url)
		}
	}
	return nil
}

func (c *Company) ChangeColors(status int, colors mJSON.Colors, x_auth_token string, x_company_id *string) error {
	var companyIDStr = c.Created.ID.String()
	cID, err := Get_x_company_id(x_company_id, &companyIDStr)
	if err != nil {
		return err
	}
	if err := handler.NewHttpClient().
		Method("PUT").
		URL(fmt.Sprintf("/company/%s/design/colors", c.Created.ID.String())).
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
	if imageURL == "" {
		return fmt.Errorf("image URL cannot be empty")
	}
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
			return fmt.Errorf("received empty response for image (%s)", imageURL)
		} else if len(response) != len(*compareImgBytes) {
			return fmt.Errorf("image size mismatch for %s: expected %d bytes, got %d bytes", imageURL, len(*compareImgBytes), len(response))
		} else if !bytes.Equal(response, *compareImgBytes) {
			return fmt.Errorf("image content mismatch for %s", imageURL)
		}
	}
	return nil
}

func (c *Company) GetRandomService() (*Service, error) {
	if len(c.Services) == 0 {
		return nil, fmt.Errorf("no services available in the company to select a random service")
	}

	service := c.Services[lib.GenerateRandomIntFromRange(0, len(c.Services)-1)]
	if service.Created.ID == uuid.Nil {
		return nil, fmt.Errorf("selected service has nil ID")
	}
	return service, nil
}
