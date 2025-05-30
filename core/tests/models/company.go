package models_test

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	mJSON "agenda-kaki-go/core/config/db/model/json"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/lib"
	"agenda-kaki-go/core/lib/FileBytes"
	handler "agenda-kaki-go/core/tests/handlers"
	"bytes"
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/google/uuid"
)

type Company struct {
	Created    model.CompanyMerged
	Owner      *Employee
	Employees  []*Employee
	Branches   []*Branch
	Services   []*Service
	Auth_token string
}


func (c *Company) Set(t *testing.T) {
	c.Create(t, 200)
	c.Owner.Company = c
	c.Owner.VerifyEmail(t, 200)
	c.Owner.Login(t, 200)
	c.Auth_token = c.Owner.Auth_token
	c.Owner.GetById(t, 200)
	employee := &Employee{}
	employee.Company = c
	employee.Create(t, 200)
	employee.VerifyEmail(t, 200)
	employee.Login(t, 200)
	c.Employees = append(c.Employees, employee)
	branch := &Branch{}
	branch.Auth_token = c.Owner.Auth_token
	branch.Company = c
	branch.Create(t, 200)
	c.Branches = append(c.Branches, branch)
	service := &Service{}
	service.Auth_token = c.Owner.Auth_token
	service.Company = c
	service.Create(t, 200)
	c.Services = append(c.Services, service)
	c.GetById(t, 200)
	c.Employees[0].AddService(t, 200, c.Services[0])
	c.Employees[0].AddBranch(t, 200, c.Branches[0], &c.Owner.Auth_token)
	c.Branches[0].AddService(t, 200, c.Services[0], nil)
	c.Employees[0].Update(t, 200, map[string]any{"work_schedule": []mJSON.WorkSchedule{
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
	}})
	c.UploadImages(t, 200, map[string][]byte{
		"logo": FileBytes.PNG_FILE_1,
	})
}

// --- Randomized Company Setup Method ---

// SetupRandomized replaces the static Set method.
// It creates the Company, owner, and then generates the specified number
// of employees, branches, and services, linking them randomly.
func (c *Company) SetupRandomized(t *testing.T, numEmployees, numBranches, numServices int) {
	t.Logf("Setting up randomized Company with %d employees, %d branches, %d services", numEmployees, numBranches, numServices)

	// 1. Create Company and Owner (uses Company.Create from your e2e_test)
	c.Create(t, 200) // This populates c.Created and c.Owner according to your Create method.
	if c.Created.ID == uuid.Nil || c.Owner == nil || c.Owner.Created.ID == uuid.Nil {
		t.Fatal("Failed to create Company or owner structure")
	}
	t.Logf("Company Created: ID %s, Owner Email: %s", c.Created.ID, c.Owner.Created.Email)

	// Link owner back to Company and get token
	c.Owner.Company = c // Ensure back-reference for owner helper
	c.Owner.VerifyEmail(t, 200)
	c.Owner.Login(t, 200)             // Populates c.Owner.Auth_token
	c.Auth_token = c.Owner.Auth_token // Store owner's token in Company helper
	c.Owner.GetById(t, 200)           // Refresh owner data
	// Check if Create already added the owner to c.Created.Employees and maybe sync c.Employees
	foundOwner := false
	for _, empHelper := range c.Employees { // Assuming c.Employees is also populated by Create/Set logic initially if owner added
		if empHelper.Created.ID == c.Owner.Created.ID {
			foundOwner = true
			break
		}
	}
	if !foundOwner {
		c.Employees = append(c.Employees, c.Owner) // Add owner to the employee list if not already present
	}
	t.Log("Company owner configured and logged in.")

	// --- Entity Generation ---
	// These functions now use the Create methods from your e2e_test helpers

	c.GenerateEmployees(t, numEmployees) // Generate *additional* employees
	t.Logf("Generated %d additional employees. Total employees (incl. owner): %d", numEmployees, len(c.Employees))

	c.GenerateBranches(t, numBranches)
	t.Logf("Generated %d branches.", len(c.Branches))

	c.GenerateServices(t, numServices)
	t.Logf("Generated %d services.", len(c.Services))

	// Refresh Company data potentially, if needed after additions
	c.GetById(t, 200)

	// --- Random Relationship Assignments ---
	// These functions now use AddBranch/AddService from your e2e_test helpers

	if len(c.Employees) > 0 && len(c.Branches) > 0 {
		c.RandomlyAssignEmployeesToBranches(t)
		t.Log("Randomly assigned employees to branches.")
	} else {
		t.Log("Skipping employee-to-branch assignment (not enough employees or branches).")
	}

	if len(c.Employees) > 0 && len(c.Services) > 0 {
		c.RandomlyAssignServicesToEmployees(t)
		t.Log("Randomly assigned services to employees.")
	} else {
		t.Log("Skipping service-to-employee assignment (not enough employees or services).")
	}

	if len(c.Branches) > 0 && len(c.Services) > 0 {
		c.RandomlyAssignServicesToBranches(t)
		t.Log("Randomly assigned services to branches.")
	} else {
		t.Log("Skipping service-to-branch assignment (not enough branches or services).")
	}

	// --- Random Work Schedule Assignment ---
	// Uses Employee.Update from your e2e_test helper

	if len(c.Employees) > 0 && len(c.Branches) > 0 {
		c.RandomlyAssignWorkSchedules(t)
		t.Log("Randomly assigned work schedules to employees.")
	} else {
		t.Log("Skipping work schedule assignment (not enough employees or branches).")
	}

	t.Log("Randomized Company setup completed.")
}

// --- Generation Functions ---

// GenerateEmployees creates n *additional* employees (owner already exists).
func (c *Company) GenerateEmployees(t *testing.T, n int) {
	if n <= 0 || c.Created.ID == uuid.Nil {
		t.Log("Skipping employee generation: n <= 0 or Company ID is nil.")
		return
	}

	initialEmployeeCount := len(c.Employees)
	createdCount := 0

	for i := range n {
		employee := &Employee{Company: c}
		employee.Create(t, 200)

		if employee.Created.ID == uuid.Nil {
			t.Errorf("Failed to create employee %d/%d or retrieve ID.", i+1, n)
			continue
		}
		createdCount++
		t.Logf("Generated employee %d/%d: ID %s, Email %s", i+1, n, employee.Created.ID, employee.Created.Email)

		employee.VerifyEmail(t, 200)
		employee.Login(t, 200)
		if employee.Auth_token == "" {
			t.Logf("Warning: Employee %s failed to login after creation.", employee.Created.Email)
		}

		c.Employees = append(c.Employees, employee)
	}
	if createdCount != n {
		t.Logf("Warning: Tried to create %d employees, but only %d succeeded.", n, createdCount)
	}
	if len(c.Employees) != initialEmployeeCount+createdCount {
		t.Logf("Warning: Company employee slice length (%d) does not match expected count (%d).", len(c.Employees), initialEmployeeCount+createdCount)
	}
}

// GenerateBranches creates n branches for the Company.
func (c *Company) GenerateBranches(t *testing.T, n int) {
	if n <= 0 || c.Created.ID == uuid.Nil || c.Auth_token == "" {
		t.Log("Skipping branch generation: Prerequisite missing (n>0, Company ID, Owner Auth Token).")
		return
	}

	initialBranchCount := len(c.Branches)
	createdCount := 0

	for i := range n {
		branch := &Branch{Company: c, Auth_token: c.Auth_token}
		branch.Create(t, 200)

		if branch.Created.ID == uuid.Nil {
			t.Errorf("Failed to create branch %d/%d or retrieve ID.", i+1, n)
			continue
		}
		createdCount++
		t.Logf("Generated branch %d/%d: ID %s, Name %s", i+1, n, branch.Created.ID, branch.Created.Name)
		c.Branches = append(c.Branches, branch)
	}
	if createdCount != n {
		t.Logf("Warning: Tried to create %d branches, but only %d succeeded.", n, createdCount)
	}
	if len(c.Branches) != initialBranchCount+createdCount {
		t.Logf("Warning: Company branch slice length (%d) does not match expected count (%d).", len(c.Branches), initialBranchCount+createdCount)
	}
}

// GenerateServices creates n services for the Company.
func (c *Company) GenerateServices(t *testing.T, n int) {
	if n <= 0 || c.Created.ID == uuid.Nil || c.Auth_token == "" {
		t.Log("Skipping service generation: Prerequisite missing (n>0, Company ID, Owner Auth Token).")
		return
	}

	initialServiceCount := len(c.Services)
	createdCount := 0

	for i := range n {
		service := &Service{Company: c, Auth_token: c.Auth_token}
		service.Create(t, 200)

		if service.Created.ID == uuid.Nil {
			t.Errorf("Failed to create service %d/%d or retrieve ID.", i+1, n)
			continue
		}
		createdCount++
		t.Logf("Generated service %d/%d: ID %s, Name %s", i+1, n, service.Created.ID, service.Created.Name)
		c.Services = append(c.Services, service)
	}
	if createdCount != n {
		t.Logf("Warning: Tried to create %d services, but only %d succeeded.", n, createdCount)
	}
	if len(c.Services) != initialServiceCount+createdCount {
		t.Logf("Warning: Company service slice length (%d) does not match expected count (%d).", len(c.Services), initialServiceCount+createdCount)
	}
}

// --- Random Assignment Functions ---

// RandomlyAssignEmployeesToBranches assigns each employee to 1 to N random branches.
func (c *Company) RandomlyAssignEmployeesToBranches(t *testing.T) {
	if len(c.Employees) == 0 || len(c.Branches) == 0 {
		t.Log("No employees or branches to assign.")
		return
	}

	// Limit the number of branches assigned to each employee
	maxBranchesPerEmployee := min(len(c.Branches), 10)

	for i, employee := range c.Employees {
		if employee.Created.ID == uuid.Nil {
			t.Logf("Skipping branch assignment for employee %d (email: %s): Employee ID is nil.", i, employee.Created.Email)
			continue
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
				t.Logf("Could not find unique valid branch for employee %s after attempts.", employee.Created.Email)
				break
			}

			branch := c.Branches[branchIndex]
			assignedBranchIndices[branchIndex] = true
			assignedCount++ // Increment count of successfully assigned unique branches

			t.Logf("Assigning employee %d (%s, ID: %s) to branch %d (%s, ID: %s)",
				i, employee.Created.Email, employee.Created.ID,
				branchIndex, branch.Created.Name, branch.Created.ID)

			// Use owner token for privilege when assigning employees to branches
			employee.AddBranch(t, 200, branch, &c.Auth_token)
		}
	}
}

// RandomlyAssignServicesToEmployees assigns each employee 1 to N random services.
func (c *Company) RandomlyAssignServicesToEmployees(t *testing.T) {
	if len(c.Employees) == 0 || len(c.Services) == 0 {
		t.Log("No employees or services to assign.")
		return
	}

	// Limit the number of services assigned to each employee
	maxServicesPerEmployee := min(len(c.Services), 10)

	for i, employee := range c.Employees {
		if employee.Created.ID == uuid.Nil || employee.Auth_token == "" {
			t.Logf("Skipping service assignment for employee %d (email: %s): Employee ID nil or not logged in.", i, employee.Created.Email)
			continue
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
				t.Logf("Could not find unique valid service for employee %s after attempts.", employee.Created.Email)
				break
			}

			service := c.Services[serviceIndex]
			assignedServiceIndices[serviceIndex] = true
			assignedCount++

			t.Logf("Assigning service %d (%s, ID: %s) to employee %d (%s, ID: %s)",
				serviceIndex, service.Created.Name, service.Created.ID,
				i, employee.Created.Email, employee.Created.ID)

			// Use Employee.AddService, assumes employee's token is used
			employee.AddService(t, 200, service)
		}
	}
}

// RandomlyAssignServicesToBranches assigns each branch 1 to N random services.
func (c *Company) RandomlyAssignServicesToBranches(t *testing.T) {
	if len(c.Branches) == 0 || len(c.Services) == 0 {
		t.Log("No branches or services to assign.")
		return
	}
	// Limit the number of services assigned to each branch
	maxServicesPerBranch := min(len(c.Services), 20)

	for i, branch := range c.Branches {
		if branch.Created.ID == uuid.Nil {
			t.Logf("Skipping service assignment for branch %d (%s): Branch ID is nil.", i, branch.Created.Name)
			continue
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
				t.Logf("Could not find unique valid service for branch %s after attempts.", branch.Created.Name)
				break
			}

			service := c.Services[serviceIndex]
			assignedServiceIndices[serviceIndex] = true
			assignedCount++

			t.Logf("Assigning service %d (%s, ID: %s) to branch %d (%s, ID: %s)",
				serviceIndex, service.Created.Name, service.Created.ID,
				i, branch.Created.Name, branch.Created.ID)

			// Use Branch.AddService method. Use owner's token implicitly (branch.Auth_token or nil param).
			branch.AddService(t, 200, service, nil)
		}
	}
}

// --- Work Schedule Assignment ---

// RandomlyAssignWorkSchedules assigns a generated work schedule to each employee.
func (c *Company) RandomlyAssignWorkSchedules(t *testing.T) {
	if len(c.Employees) == 0 || len(c.Branches) == 0 {
		t.Log("No employees or branches for work schedule assignment.")
		return
	}

	validBranches := []*Branch{}
	for _, b := range c.Branches {
		if b.Created.ID != uuid.Nil {
			validBranches = append(validBranches, b)
		}
	}

	if len(validBranches) == 0 {
		t.Log("Skipping work schedule assignment: No valid branches found.")
		return
	}

	for i, employee := range c.Employees {
		if employee.Created.ID == uuid.Nil { // No need to check Auth_token here, Update uses Company token.
			t.Logf("Skipping schedule assignment for employee %d (%s): ID nil.", i, employee.Created.Email)
			continue
		}

		fmt.Println("👀 Employee services:", employee.Created.Services)

		scheduleModel := GenerateRandomModelWorkSchedule(validBranches, employee)
		t.Logf("Generated work schedule for employee %d (%s), referencing %d valid branch(es).", i, employee.Created.Email, len(validBranches))

		// Payload format for Employee.Update
		payload := map[string]any{
			"work_schedule": []mJSON.WorkSchedule{scheduleModel},
		}

		// Call Employee.Update using owner's token (c.Auth_token is implicitly used in helper via employee.Company.Auth_token)
		employee.Update(t, 200, payload)
		t.Logf("Attempted to update work schedule for employee %d (%s) via API.", i, employee.Created.Email)

		// Optional: Refresh employee data locally if needed, though Update should handle API state
		// employee.GetById(t, 200)
	}
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

func (c *Company) Create(t *testing.T, status int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("POST")
	http.URL("/Company")
	http.ExpectStatus(status)
	ownerPswd := "1SecurePswd!"
	http.Send(DTO.CreateCompany{
		LegalName:      lib.GenerateRandomName("Company Legal Name"),
		TradeName:      lib.GenerateRandomName("Company Trade Name"),
		TaxID:          lib.GenerateRandomStrNumber(14),
		OwnerName:      lib.GenerateRandomName("Owner Name"),
		OwnerSurname:   lib.GenerateRandomName("Owner Surname"),
		OwnerEmail:     lib.GenerateRandomEmail("owner"),
		OwnerPhone:     lib.GenerateRandomPhoneNumber(),
		OwnerPassword:  ownerPswd,
		StartSubdomain: strings.ToLower(lib.GenerateRandomString(12)),
	})
	http.ParseResponse(&c.Created)
	owner := c.Created.Employees[0]
	owner.Password = ownerPswd
	c.Owner = &Employee{
		Company: c,
		Created: owner,
	}
}

func (c *Company) GetByName(t *testing.T, status int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("GET")
	http.URL(fmt.Sprintf("/Company/name/%s", c.Created.LegalName))
	http.ExpectStatus(status)
	http.Send(nil)
	http.ParseResponse(&c.Created)
}

func (c *Company) GetBySubdomain(t *testing.T, status int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("GET")
	http.URL(fmt.Sprintf("/Company/subdomain/%s", c.Created.Subdomains[0].Name))
	http.ExpectStatus(status)
	http.Send(nil)
	http.ParseResponse(&c.Created)
}

func (c *Company) GetById(t *testing.T, status int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("GET")
	http.URL(fmt.Sprintf("/Company/%s", c.Created.ID.String()))
	http.ExpectStatus(status)
	http.Send(nil)
	http.ParseResponse(&c.Created)
}

func (c *Company) Update(t *testing.T, status int, changes map[string]any) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("PATCH")
	http.URL(fmt.Sprintf("/Company/%s", c.Created.ID.String()))
	http.ExpectStatus(status)
	http.Header(namespace.HeadersKey.Auth, c.Auth_token)
	http.Header(namespace.HeadersKey.Company, c.Created.ID.String())
	http.Send(changes)
}

func (c *Company) Delete(t *testing.T, status int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("DELETE")
	http.URL(fmt.Sprintf("/Company/%s", c.Created.ID.String()))
	http.ExpectStatus(status)
	http.Header(namespace.HeadersKey.Auth, c.Auth_token)
	http.Header(namespace.HeadersKey.Company, c.Created.ID.String())
	http.Send(nil)
}

func (c *Company) UploadImages(t *testing.T, status int, files map[string][]byte) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("PATCH")
	http.URL(fmt.Sprintf("/Company/%s/design/images", c.Created.ID.String()))
	http.ExpectStatus(status)
	http.Header(namespace.HeadersKey.Auth, c.Auth_token)
	http.Header(namespace.HeadersKey.Company, c.Created.ID.String())
	var fileMap = make(handler.Files)
	for field, content := range files {
		fileMap[field] = handler.MyFile{
			Name:    field + "_" + lib.GenerateRandomString(6) + ".png",
			Content: content,
		}
	}
	http.Send(fileMap)
	http.ParseResponse(&c.Created.Design.Images)
}

func (c *Company) DeleteImages(t *testing.T, status int, images []string) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("DELETE")
	base_url := fmt.Sprintf("/Company/%s/design/images", c.Created.ID.String())
	http.ExpectStatus(status)
	http.Header(namespace.HeadersKey.Auth, c.Auth_token)
	http.Header(namespace.HeadersKey.Company, c.Created.ID.String())
	for _, field := range images {
		image_url := base_url + "/" + field
		http.URL(image_url)
		http.Send(nil)
		http.ParseResponse(&c.Created.Design)
		if c.Created.Design.Images.GetImageURL(field) != "" {
			t.Errorf("Image %s was not deleted successfully. Expected empty URL, got %s", field, c.Created.Design.Images.GetImageURL(field))
		} else {
			t.Logf("Image %s deleted successfully.", field)
		}
	}
}

func (c *Company) ChangeColors(t *testing.T, status int, colors mJSON.Colors) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("PUT")
	http.URL(fmt.Sprintf("/Company/%s/design/colors", c.Created.ID.String()))
	http.ExpectStatus(status)
	http.Header(namespace.HeadersKey.Auth, c.Auth_token)
	http.Header(namespace.HeadersKey.Company, c.Created.ID.String())
	http.Send(colors)
	http.ParseResponse(&c.Created.Design.Colors)
	if c.Created.Design.Colors != colors {
		t.Errorf("Colors were not updated correctly. Expected %v, got %v", colors, c.Created.Design.Colors)
	} else {
		t.Logf("Colors updated successfully to %v", colors)
	}
}

func (c *Company) GetImage(t *testing.T, status int, imageURL string, compareImgBytes *[]byte) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("GET")
	http.URL(imageURL)
	http.ExpectStatus(status)
	http.Header(namespace.HeadersKey.Auth, c.Auth_token)
	http.Header(namespace.HeadersKey.Company, c.Created.ID.String())
	http.Send(nil)
	// Compare the response bytes with the expected image bytes
	if compareImgBytes != nil {
		var response []byte
		http.ParseResponse(&response)
		if len(response) == 0 {
			t.Error("Received empty response for image")
		} else if len(response) != len(*compareImgBytes) {
			t.Errorf("Image size mismatch: expected %d bytes, got %d bytes", len(*compareImgBytes), len(response))
		} else if !bytes.Equal(response, *compareImgBytes) {
			t.Logf("Image content mismatch for %s", imageURL)
		}
	} else {
		t.Logf("Image fetched successfully from %s", imageURL)
	}
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