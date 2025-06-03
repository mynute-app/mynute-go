package modelT

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	mJSON "agenda-kaki-go/core/config/db/model/json"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/lib"
	"agenda-kaki-go/core/lib/FileBytes"
	handler "agenda-kaki-go/core/test/handlers"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/google/uuid"
)

type Company struct {
	Created    model.CompanyMerged
	Owner      *Employee
	Employees  []*Employee
	Branches   []*Branch
	Services   []*Service
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
		c.Employees = append(c.Employees, employee)
	}

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
		c.Branches = append(c.Branches, branch)
	}

	if err := c.GetById(200, c.Owner.X_Auth_Token, nil); err != nil {
		return err
	}
	for _, employee := range c.Employees {
		for _, branch := range c.Branches {
			if err := employee.AddBranch(200, branch, &c.Owner.X_Auth_Token, nil); err != nil {
				return fmt.Errorf("failed to assign employee %s to branch %s: %v", employee.Created.Email, branch.Created.Name, err)
			}
		}
		for _, service := range c.Services {
			if err := employee.AddService(200, service, &c.Owner.X_Auth_Token, nil); err != nil {
				return fmt.Errorf("failed to assign employee %s to service %s: %v", employee.Created.Email, service.Created.Name, err)
			}
		}
		if err := employee.Update(200, map[string]any{"work_schedule": []mJSON.WorkSchedule{
			{
				Monday: []mJSON.WorkRange{
					{Start: "08:00", End: "12:00", BranchID: c.Branches[0].Created.ID},
					{Start: "13:00", End: "17:00", BranchID: c.Branches[0].Created.ID},
				},
				Tuesday: []mJSON.WorkRange{
					{Start: "09:00", End: "12:00", BranchID: c.Branches[0].Created.ID},
					{Start: "13:00", End: "18:00", BranchID: c.Branches[0].Created.ID},
				},
				Wednesday: []mJSON.WorkRange{
					{Start: "08:00", End: "12:00", BranchID: c.Branches[0].Created.ID},
					{Start: "13:00", End: "17:00", BranchID: c.Branches[0].Created.ID},
				},
				Thursday: []mJSON.WorkRange{
					{Start: "08:00", End: "12:00", BranchID: c.Branches[0].Created.ID},
					{Start: "13:00", End: "17:00", BranchID: c.Branches[0].Created.ID},
				},
				Friday: []mJSON.WorkRange{
					{Start: "08:00", End: "12:00", BranchID: c.Branches[0].Created.ID},
					{Start: "13:00", End: "17:00", BranchID: c.Branches[0].Created.ID},
				},
				Saturday: []mJSON.WorkRange{
					{Start: "08:00", End: "12:00", BranchID: c.Branches[0].Created.ID},
					{Start: "13:00", End: "17:00", BranchID: c.Branches[0].Created.ID},
				},
				Sunday: []mJSON.WorkRange{},
			},
		}}, nil, nil); err != nil {
			return err
		}
	}

	if err := c.UploadImages(200, map[string][]byte{
		"logo": FileBytes.PNG_FILE_1,
	}, c.Owner.X_Auth_Token, nil); err != nil {
		return err
	}

	return nil
}

// --- Randomized Company Setup Method ---

// SetupRandomized replaces the static Set method.
// It creates the Company, owner, and then generates the specified number
// of employees, branches, and services, linking them randomly.
func (c *Company) SetupRandomized(numEmployees, numBranches, numServices int) error {
	if os.Getenv("APP_ENV") == "prod" { // Make sure it never runs in production
		panic("Cannot run tests in production environment. Set APP_ENV to 'test' or 'dev'")
	}

	fmt.Printf("Setting up randomized Company with %d employees, %d branches, %d services", numEmployees, numBranches, numServices)

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

	fmt.Println("Randomized Company setup completed")
	return nil
}

// --- Generation Functions ---

// GenerateEmployees creates n *additional* employees (owner already exists).
func (c *Company) GenerateEmployees(n int) error {
	fmt.Printf("Generating %d employees for Company ID %s", n, c.Created.ID)

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

		createdCount++
		fmt.Printf("Generated employee %d/%d: ID %s, Email %s", i+1, n, employee.Created.ID, employee.Created.Email)

		if err := employee.VerifyEmail(200, nil); err != nil {
			return err
		}
		if err := employee.Login(200, nil); err != nil {
			return err
		}

		if employee.X_Auth_Token == "" {
			return fmt.Errorf("failed to login employee %d/%d or retrieve Auth Token", i+1, n)
		}

		c.Employees = append(c.Employees, employee)
	}
	if createdCount != n {
		return fmt.Errorf("tried to create %d employees, but only %d succeeded", n, createdCount)
	}
	if len(c.Employees) != initialEmployeeCount+createdCount {
		return fmt.Errorf("company employee slice length (%d) does not match expected count (%d)", len(c.Employees), initialEmployeeCount+createdCount)
	}
	fmt.Printf("Generated %d employees. Initial count was %d", createdCount, initialEmployeeCount)
	fmt.Printf("Total employees after generation: %d", len(c.Employees))
	return nil
}

// GenerateBranches creates n branches for the Company.
func (c *Company) GenerateBranches(n int) error {
	fmt.Printf("Generating %d branches for Company ID %s", n, c.Created.ID)

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
		fmt.Printf("Generated branch %d/%d: ID %s, Name %s", i+1, n, branch.Created.ID, branch.Created.Name)
		c.Branches = append(c.Branches, branch)
	}
	if createdCount != n {
		return fmt.Errorf("tried to create %d branches, but only %d succeeded", n, createdCount)
	}
	if len(c.Branches) != initialBranchCount+createdCount {
		return fmt.Errorf("company branch slice length (%d) does not match expected count (%d)", len(c.Branches), initialBranchCount+createdCount)
	}
	fmt.Printf("Generated %d branches. Initial count was %d", createdCount, initialBranchCount)
	fmt.Printf("Total branches after generation: %d", len(c.Branches))
	return nil
}

// GenerateServices creates n services for the Company.
func (c *Company) GenerateServices(n int) error {
	fmt.Printf("Generating %d services for Company ID %s", n, c.Created.ID)

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
		fmt.Printf("Generated service %d/%d: ID %s, Name %s", i+1, n, service.Created.ID, service.Created.Name)
		c.Services = append(c.Services, service)
	}
	if createdCount != n {
		return fmt.Errorf("tried to create %d services, but only %d succeeded", n, createdCount)
	}
	if len(c.Services) != initialServiceCount+createdCount {
		return fmt.Errorf("Company service slice length (%d) does not match expected count (%d)", len(c.Services), initialServiceCount+createdCount)
	}
	fmt.Printf("Generated %d services. Initial count was %d", createdCount, initialServiceCount)
	fmt.Printf("Total services after generation: %d", len(c.Services))
	return nil
}

// --- Random Assignment Functions ---

// RandomlyAssignEmployeesToBranches assigns each employee to 1 to N random branches.
func (c *Company) RandomlyAssignEmployeesToBranches() error {
	if len(c.Employees) == 0 || len(c.Branches) == 0 {
		return fmt.Errorf("no employees or branches to assign")
	}

	// Limit the number of branches assigned to each employee
	maxBranchesPerEmployee := min(len(c.Branches), 10)

	for i, employee := range c.Employees {
		if employee.Created.ID == uuid.Nil {
			return fmt.Errorf("skipping branch assignment for employee %d (email: %s): Employee ID is nil", i, employee.Created.Email)
		}

		numBranchesToAssign := 1
		if maxBranchesPerEmployee > 1 {
			numBranchesToAssign = rand.Intn(maxBranchesPerEmployee) + 1
		}
		assignedBranchIndices := make(map[int]bool)
		assignedCount := 0

		for k := 0; k < numBranchesToAssign && assignedCount < len(c.Branches); k++ { // Use assignedCount guard
			branchIndex := -1
			for attempts := 0; attempts < len(c.Branches)*2; attempts++ { // Limit attempts
				potentialIndex := rand.Intn(len(c.Branches))
				if !assignedBranchIndices[potentialIndex] && c.Branches[potentialIndex].Created.ID != uuid.Nil {
					branchIndex = potentialIndex
					break
				}
			}
			if branchIndex == -1 {
				return fmt.Errorf("could not find unique valid branch for employee %s after attempts", employee.Created.Email)
			}

			branch := c.Branches[branchIndex]
			assignedBranchIndices[branchIndex] = true
			assignedCount++ // Increment count of successfully assigned unique branches

			fmt.Printf("Assigning employee %d (%s, ID: %s) to branch %d (%s, ID: %s)",
				i, employee.Created.Email, employee.Created.ID,
				branchIndex, branch.Created.Name, branch.Created.ID)

			// Use owner token for privilege when assigning employees to branches
			if err := employee.AddBranch(200, branch, &c.Owner.X_Auth_Token, nil); err != nil {
				return fmt.Errorf("failed to assign employee %d (%s, ID: %s) to branch %d (%s, ID: %s): %v",
					i, employee.Created.Email, employee.Created.ID,
					branchIndex, branch.Created.Name, branch.Created.ID, err)
			}
		}
	}
	return nil
}

// RandomlyAssignServicesToEmployees assigns each employee 1 to N random services.
func (c *Company) RandomlyAssignServicesToEmployees() error {
	if len(c.Employees) == 0 || len(c.Services) == 0 {
		return fmt.Errorf("no employees or services to assign")
	}

	// Limit the number of services assigned to each employee
	maxServicesPerEmployee := min(len(c.Services), 10)

	for i, employee := range c.Employees {
		if employee.Created.ID == uuid.Nil || employee.X_Auth_Token == "" {
			return fmt.Errorf("skipping service assignment for employee %d (email: %s): Employee ID nil or not logged in", i, employee.Created.Email)
		}

		numServicesToAssign := 1
		if maxServicesPerEmployee > 1 {
			numServicesToAssign = rand.Intn(maxServicesPerEmployee) + 1
		}
		assignedServiceIndices := make(map[int]bool)
		assignedCount := 0

		for k := 0; k < numServicesToAssign && assignedCount < len(c.Services); k++ {
			serviceIndex := -1
			for range len(c.Services) * 2 {
				potentialIndex := rand.Intn(len(c.Services))
				if !assignedServiceIndices[potentialIndex] && c.Services[potentialIndex].Created.ID != uuid.Nil {
					serviceIndex = potentialIndex
					break
				}
			}
			if serviceIndex == -1 {
				return fmt.Errorf("could not find unique valid service for employee %s after attempts", employee.Created.Email)
			}

			service := c.Services[serviceIndex]
			assignedServiceIndices[serviceIndex] = true
			assignedCount++

			fmt.Printf("Assigning service %d (%s, ID: %s) to employee %d (%s, ID: %s)",
				serviceIndex, service.Created.Name, service.Created.ID,
				i, employee.Created.Email, employee.Created.ID)

			// Use Employee.AddService, assumes employee's token is used
			if err := employee.AddService(200, service, nil, nil); err != nil {
				return fmt.Errorf("failed to assign service %d (%s, ID: %s) to employee %d (%s, ID: %s): %v",
					serviceIndex, service.Created.Name, service.Created.ID,
					i, employee.Created.Email, employee.Created.ID, err)
			}
		}
	}
	return nil
}

// RandomlyAssignServicesToBranches assigns each branch 1 to N random services.
func (c *Company) RandomlyAssignServicesToBranches() error {
	if len(c.Branches) == 0 || len(c.Services) == 0 {
		return fmt.Errorf("no branches or services to assign")
	}
	// Limit the number of services assigned to each branch
	maxServicesPerBranch := min(len(c.Services), 20)

	for i, branch := range c.Branches {
		if branch.Created.ID == uuid.Nil {
			return fmt.Errorf("skipping service assignment for branch %d (%s): Branch ID is nil", i, branch.Created.Name)
		}

		numServicesToAssign := 1
		if maxServicesPerBranch > 1 {
			numServicesToAssign = rand.Intn(maxServicesPerBranch) + 1
		}
		assignedServiceIndices := make(map[int]bool)
		assignedCount := 0

		for k := 0; k < numServicesToAssign && assignedCount < len(c.Services); k++ {
			serviceIndex := -1
			for attempts := 0; attempts < len(c.Services)*2; attempts++ {
				potentialIndex := rand.Intn(len(c.Services))
				if !assignedServiceIndices[potentialIndex] && c.Services[potentialIndex].Created.ID != uuid.Nil {
					serviceIndex = potentialIndex
					break
				}
			}
			if serviceIndex == -1 {
				return fmt.Errorf("could not find unique valid service for branch %s after attempts", branch.Created.Name)
			}

			service := c.Services[serviceIndex]
			assignedServiceIndices[serviceIndex] = true
			assignedCount++

			fmt.Printf("Assigning service %d (%s, ID: %s) to branch %d (%s, ID: %s)",
				serviceIndex, service.Created.Name, service.Created.ID,
				i, branch.Created.Name, branch.Created.ID)

			// Use Branch.AddService method. Use owner's token implicitly (branch.X_Auth_Token or nil param).
			if err := branch.AddService(200, service, c.Owner.X_Auth_Token, nil); err != nil {
				return fmt.Errorf("failed to assign service %d (%s, ID: %s) to branch %d (%s, ID: %s): %v",
					serviceIndex, service.Created.Name, service.Created.ID,
					i, branch.Created.Name, branch.Created.ID, err)
			}
		}
	}
	return nil
}

// --- Work Schedule Assignment ---

// RandomlyAssignWorkSchedules assigns a generated work schedule to each employee.
func (c *Company) RandomlyAssignWorkScheduleToEmployees() error {
	if len(c.Employees) == 0 {
		return fmt.Errorf("no employees to assign work schedules")
	} else if len(c.Branches) == 0 {
		return fmt.Errorf("no branches to reference for work schedule assignment")
	} else if len(c.Services) == 0 {
		return fmt.Errorf("no services to reference for work schedule assignment")
	}

	validBranches := []*Branch{}
	for _, b := range c.Branches {
		if b.Created.ID != uuid.Nil {
			validBranches = append(validBranches, b)
		}
	}

	if len(validBranches) == 0 {
		return fmt.Errorf("no valid branches found for work schedule assignment")
	}

	for i, employee := range c.Employees {
		if employee.Created.ID == uuid.Nil { // No need to check Auth_token here, Update uses Company token.
			return fmt.Errorf("skipping schedule assignment for employee %d (%s): ID nil", i, employee.Created.Email)
		}

		fmt.Println("👀 Employee services:", employee.Created.Services)

		scheduleModel := GenerateRandomModelWorkSchedule(validBranches, employee)
		fmt.Printf("Generated work schedule for employee %d (%s), referencing %d valid branch(es)", i, employee.Created.Email, len(validBranches))

		// Payload format for Employee.Update
		payload := map[string]any{
			"work_schedule": []mJSON.WorkSchedule{scheduleModel},
		}

		// Call Employee.Update using owner's token (c.X_Auth_Token is implicitly used in helper via employee.Company.X_Auth_Token)
		if err := employee.Update(200, payload, nil, nil); err != nil {
			return fmt.Errorf("failed to update work schedule for employee %d (%s) via API: %v",
				i, employee.Created.Email, err)
		}
		fmt.Printf("Successfully updated work schedule for employee %d (%s)", i, employee.Created.Email)
	}
	return nil
}

// GenerateRandomModelWorkSchedule creates a *mJSON.WorkSchedule* struct
func GenerateRandomModelWorkSchedule(validBranches []*Branch, employee *Employee) mJSON.WorkSchedule {
	schedule := mJSON.WorkSchedule{}

	randomTimeStringHHMM := func(minHour, maxHour int) string {
		hour := min(max(minHour+rand.Intn(maxHour-minHour+1), 6), 22) // Ensure hour is between 6 and 22
		minute := rand.Intn(4) * 15
		return fmt.Sprintf("%02d:%02d", hour, minute)
	}

	schedule.Monday = generateRangesForDayModel(validBranches, employee, randomTimeStringHHMM, 0.9)
	schedule.Tuesday = generateRangesForDayModel(validBranches, employee, randomTimeStringHHMM, 0.9)
	schedule.Wednesday = generateRangesForDayModel(validBranches, employee, randomTimeStringHHMM, 0.9)
	schedule.Thursday = generateRangesForDayModel(validBranches, employee, randomTimeStringHHMM, 0.9)
	schedule.Friday = generateRangesForDayModel(validBranches, employee, randomTimeStringHHMM, 0.9)
	schedule.Saturday = generateRangesForDayModel(validBranches, employee, randomTimeStringHHMM, 0.4)
	schedule.Sunday = generateRangesForDayModel(validBranches, employee, randomTimeStringHHMM, 0.1)

	return schedule
}

// Helper for GenerateRandomModelWorkSchedule, returns []mJSON.WorkRange
func generateRangesForDayModel(validBranches []*Branch, employee *Employee, randomTime func(int, int) string, workProbability float32) []mJSON.WorkRange {
	// Use global rand.Float32() - auto-seeded
	if rand.Float32() > workProbability || len(validBranches) == 0 {
		return []mJSON.WorkRange{}
	}

	ranges := []mJSON.WorkRange{}
	numRanges := 1 + rand.Intn(2) // Use global rand.Intn()

	lastEndTimeStr := "00:00"

	for r := range numRanges {
		// Use global rand.Intn()
		targetBranchHelper := validBranches[rand.Intn(len(validBranches))]

		fmt.Println("🏢 Branch services:", targetBranchHelper.Services)

		startHourLower := 7
		if r > 0 {
			hourPart := 0
			minutePart := 0
			fmt.Sscanf(lastEndTimeStr, "%02d:%02d", &hourPart, &minutePart) // Ensure Sscanf matches HH:MM
			startHourLower = hourPart
			if minutePart > 0 {
				startHourLower++
			}
			startHourLower++ // Buffer hour
			if startHourLower < 13 && hourPart >= 12 {
				startHourLower = 13
			}
		}
		if startHourLower > 19 {
			continue
		}

		// Use global rand.Intn()
		startHour := startHourLower + rand.Intn(2)
		startTime := randomTime(startHour, startHour)

		if r > 0 && startTime <= lastEndTimeStr {
			startHour++
			if startHour > 20 {
				continue
			}
			startTime = randomTime(startHour, startHour)
		}

		// Use global rand.Intn()
		durationHours := 2 + rand.Intn(4)
		endHour := min(startHour+durationHours, 22)
		endTime := randomTime(endHour, endHour)

		if endTime <= startTime {
			if endHour < 22 {
				endHour++
				endTime = randomTime(endHour, endHour)
			} else {
				endTime = "23:00"
			}
			if endTime <= startTime {
				endTime = "23:30"
			}
		}

		employeeServices := []uuid.UUID{}
		for _, svc := range employee.Created.Services {
			if svc.ID != uuid.Nil {
				employeeServices = append(employeeServices, svc.ID)
			}
		}

		branchServices := []uuid.UUID{}
		for _, svc := range targetBranchHelper.Services {
			if svc.Created.ID != uuid.Nil {
				branchServices = append(branchServices, svc.Created.ID)
			}
		}

		commonServices := intersectUUIDs(employeeServices, branchServices)

		fmt.Println("🔀 Intersected services:", commonServices)

		if len(commonServices) == 0 {
			continue // pula esse range se não houver serviços em comum
		}

		ranges = append(ranges, mJSON.WorkRange{
			Start:    startTime,
			End:      endTime,
			BranchID: targetBranchHelper.Created.ID,
			Services: commonServices,
		})

		lastEndTimeStr = endTime
	}
	return ranges
}

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
		Created: owner,
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
		Error; err != nil {
		return fmt.Errorf("failed to update company: %w", err)
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
