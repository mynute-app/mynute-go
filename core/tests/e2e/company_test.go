package e2e_test

import (
	"agenda-kaki-go/core"
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	mJSON "agenda-kaki-go/core/config/db/model/json"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/lib"
	handler "agenda-kaki-go/core/tests/handlers"
	"fmt"
	"math/rand"
	"testing"

	"github.com/google/uuid"
)

type Company struct {
	created    model.CompanyMerged
	owner      *Employee
	employees  []*Employee
	branches   []*Branch
	services   []*Service
	auth_token string
}

func Test_Company(t *testing.T) {
	server := core.NewServer().Run("test")
	defer server.Shutdown()
	company := &Company{}
	company.Create(t, 200)
	company.owner.VerifyEmail(t, 200)
	company.owner.Login(t, 200)
	company.auth_token = company.owner.auth_token
	company.Update(t, 200, map[string]any{"name": "Updated Company Name"})
	company.GetById(t, 200)
	company.GetByName(t, 200)
	company.Delete(t, 200)
}

func (c *Company) Set(t *testing.T) {
	c.Create(t, 200)
	c.owner.company = c
	c.owner.VerifyEmail(t, 200)
	c.owner.Login(t, 200)
	c.auth_token = c.owner.auth_token
	c.owner.GetById(t, 200)
	employee := &Employee{}
	employee.company = c
	employee.Create(t, 200)
	employee.VerifyEmail(t, 200)
	employee.Login(t, 200)
	c.employees = append(c.employees, employee)
	branch := &Branch{}
	branch.auth_token = c.owner.auth_token
	branch.company = c
	branch.Create(t, 200)
	c.branches = append(c.branches, branch)
	service := &Service{}
	service.auth_token = c.owner.auth_token
	service.company = c
	service.Create(t, 200)
	c.services = append(c.services, service)
	c.GetById(t, 200)
	c.employees[0].AddService(t, 200, c.services[0])
	c.employees[0].AddBranch(t, 200, c.branches[0], &c.owner.auth_token)
	c.branches[0].AddService(t, 200, c.services[0], nil)
	c.employees[0].Update(t, 200, map[string]any{"work_schedule": []mJSON.WorkSchedule{
		{
			Monday: []mJSON.WorkRange{
				{Start: "08:00", End: "12:00", BranchID: c.branches[0].created.ID},
				{Start: "13:00", End: "17:00", BranchID: c.branches[0].created.ID},
			},
			Tuesday: []mJSON.WorkRange{
				{Start: "09:00", End: "12:00", BranchID: c.branches[0].created.ID},
				{Start: "13:00", End: "18:00", BranchID: c.branches[0].created.ID},
			},
			Wednesday: []mJSON.WorkRange{
				{Start: "08:00", End: "12:00", BranchID: c.branches[0].created.ID},
				{Start: "13:00", End: "17:00", BranchID: c.branches[0].created.ID},
			},
			Thursday: []mJSON.WorkRange{
				{Start: "08:00", End: "12:00", BranchID: c.branches[0].created.ID},
				{Start: "13:00", End: "17:00", BranchID: c.branches[0].created.ID},
			},
			Friday: []mJSON.WorkRange{
				{Start: "08:00", End: "12:00", BranchID: c.branches[0].created.ID},
				{Start: "13:00", End: "17:00", BranchID: c.branches[0].created.ID},
			},
			Saturday: []mJSON.WorkRange{
				{Start: "08:00", End: "12:00", BranchID: c.branches[0].created.ID},
				{Start: "13:00", End: "17:00", BranchID: c.branches[0].created.ID},
			},
			Sunday: []mJSON.WorkRange{},
		},
	}})

}

// --- Randomized Company Setup Method ---

// SetupRandomized replaces the static Set method.
// It creates the company, owner, and then generates the specified number
// of employees, branches, and services, linking them randomly.
func (c *Company) SetupRandomized(t *testing.T, numEmployees, numBranches, numServices int) {
	t.Logf("Setting up randomized company with %d employees, %d branches, %d services", numEmployees, numBranches, numServices)

	// 1. Create Company and Owner (uses Company.Create from your e2e_test)
	c.Create(t, 200) // This populates c.created and c.owner according to your Create method.
	if c.created.ID == uuid.Nil || c.owner == nil || c.owner.created.ID == uuid.Nil {
		t.Fatal("Failed to create company or owner structure")
	}
	t.Logf("Company created: ID %s, Owner Email: %s", c.created.ID, c.owner.created.Email)

	// Link owner back to company and get token
	c.owner.company = c // Ensure back-reference for owner helper
	c.owner.VerifyEmail(t, 200)
	c.owner.Login(t, 200)             // Populates c.owner.auth_token
	c.auth_token = c.owner.auth_token // Store owner's token in company helper
	c.owner.GetById(t, 200)           // Refresh owner data
	// Check if Create already added the owner to c.created.Employees and maybe sync c.employees
	foundOwner := false
	for _, empHelper := range c.employees { // Assuming c.employees is also populated by Create/Set logic initially if owner added
		if empHelper.created.ID == c.owner.created.ID {
			foundOwner = true
			break
		}
	}
	if !foundOwner {
		c.employees = append(c.employees, c.owner) // Add owner to the employee list if not already present
	}
	t.Log("Company owner configured and logged in.")

	// --- Entity Generation ---
	// These functions now use the Create methods from your e2e_test helpers

	c.GenerateEmployees(t, numEmployees) // Generate *additional* employees
	t.Logf("Generated %d additional employees. Total employees (incl. owner): %d", numEmployees, len(c.employees))

	c.GenerateBranches(t, numBranches)
	t.Logf("Generated %d branches.", len(c.branches))

	c.GenerateServices(t, numServices)
	t.Logf("Generated %d services.", len(c.services))

	// Refresh company data potentially, if needed after additions
	c.GetById(t, 200)

	// --- Random Relationship Assignments ---
	// These functions now use AddBranch/AddService from your e2e_test helpers

	if len(c.employees) > 0 && len(c.branches) > 0 {
		c.RandomlyAssignEmployeesToBranches(t)
		t.Log("Randomly assigned employees to branches.")
	} else {
		t.Log("Skipping employee-to-branch assignment (not enough employees or branches).")
	}

	if len(c.employees) > 0 && len(c.services) > 0 {
		c.RandomlyAssignServicesToEmployees(t)
		t.Log("Randomly assigned services to employees.")
	} else {
		t.Log("Skipping service-to-employee assignment (not enough employees or services).")
	}

	if len(c.branches) > 0 && len(c.services) > 0 {
		c.RandomlyAssignServicesToBranches(t)
		t.Log("Randomly assigned services to branches.")
	} else {
		t.Log("Skipping service-to-branch assignment (not enough branches or services).")
	}

	// --- Random Work Schedule Assignment ---
	// Uses Employee.Update from your e2e_test helper

	if len(c.employees) > 0 && len(c.branches) > 0 {
		c.RandomlyAssignWorkSchedules(t)
		t.Log("Randomly assigned work schedules to employees.")
	} else {
		t.Log("Skipping work schedule assignment (not enough employees or branches).")
	}

	t.Log("Randomized company setup completed.")
}

// --- Generation Functions ---

// GenerateEmployees creates n *additional* employees (owner already exists).
func (c *Company) GenerateEmployees(t *testing.T, n int) {
	if n <= 0 || c.created.ID == uuid.Nil {
		t.Log("Skipping employee generation: n <= 0 or Company ID is nil.")
		return
	}
	initialEmployeeCount := len(c.employees)
	createdCount := 0
	for i := 0; i < n; i++ {
		employee := &Employee{company: c}
		employee.Create(t, 200)

		if employee.created.ID == uuid.Nil {
			t.Errorf("Failed to create employee %d/%d or retrieve ID.", i+1, n)
			continue
		}
		createdCount++
		t.Logf("Generated employee %d/%d: ID %s, Email %s", i+1, n, employee.created.ID, employee.created.Email)

		employee.VerifyEmail(t, 200)
		employee.Login(t, 200)
		if employee.auth_token == "" {
			t.Logf("Warning: Employee %s failed to login after creation.", employee.created.Email)
		}

		c.employees = append(c.employees, employee)
	}
	if createdCount != n {
		t.Logf("Warning: Tried to create %d employees, but only %d succeeded.", n, createdCount)
	}
	if len(c.employees) != initialEmployeeCount+createdCount {
		t.Logf("Warning: Company employee slice length (%d) does not match expected count (%d).", len(c.employees), initialEmployeeCount+createdCount)
	}
}

// GenerateBranches creates n branches for the company.
func (c *Company) GenerateBranches(t *testing.T, n int) {
	if n <= 0 || c.created.ID == uuid.Nil || c.auth_token == "" {
		t.Log("Skipping branch generation: Prerequisite missing (n>0, Company ID, Owner Auth Token).")
		return
	}
	initialBranchCount := len(c.branches)
	createdCount := 0
	for i := 0; i < n; i++ {
		branch := &Branch{company: c, auth_token: c.auth_token}
		branch.Create(t, 200)

		if branch.created.ID == uuid.Nil {
			t.Errorf("Failed to create branch %d/%d or retrieve ID.", i+1, n)
			continue
		}
		createdCount++
		t.Logf("Generated branch %d/%d: ID %s, Name %s", i+1, n, branch.created.ID, branch.created.Name)
		c.branches = append(c.branches, branch)
	}
	if createdCount != n {
		t.Logf("Warning: Tried to create %d branches, but only %d succeeded.", n, createdCount)
	}
	if len(c.branches) != initialBranchCount+createdCount {
		t.Logf("Warning: Company branch slice length (%d) does not match expected count (%d).", len(c.branches), initialBranchCount+createdCount)
	}
}

// GenerateServices creates n services for the company.
func (c *Company) GenerateServices(t *testing.T, n int) {
	if n <= 0 || c.created.ID == uuid.Nil || c.auth_token == "" {
		t.Log("Skipping service generation: Prerequisite missing (n>0, Company ID, Owner Auth Token).")
		return
	}
	initialServiceCount := len(c.services)
	createdCount := 0
	for i := 0; i < n; i++ {
		service := &Service{company: c, auth_token: c.auth_token}
		service.Create(t, 200)

		if service.created.ID == uuid.Nil {
			t.Errorf("Failed to create service %d/%d or retrieve ID.", i+1, n)
			continue
		}
		createdCount++
		t.Logf("Generated service %d/%d: ID %s, Name %s", i+1, n, service.created.ID, service.created.Name)
		c.services = append(c.services, service)
	}
	if createdCount != n {
		t.Logf("Warning: Tried to create %d services, but only %d succeeded.", n, createdCount)
	}
	if len(c.services) != initialServiceCount+createdCount {
		t.Logf("Warning: Company service slice length (%d) does not match expected count (%d).", len(c.services), initialServiceCount+createdCount)
	}
}

// --- Random Assignment Functions ---

// RandomlyAssignEmployeesToBranches assigns each employee to 1 to N random branches.
func (c *Company) RandomlyAssignEmployeesToBranches(t *testing.T) {
	if len(c.employees) == 0 || len(c.branches) == 0 {
		t.Log("No employees or branches to assign.")
		return
	}
	maxBranchesPerEmployee := 3
	if len(c.branches) < maxBranchesPerEmployee {
		maxBranchesPerEmployee = len(c.branches)
	}

	for i, employee := range c.employees {
		if employee.created.ID == uuid.Nil {
			t.Logf("Skipping branch assignment for employee %d (email: %s): Employee ID is nil.", i, employee.created.Email)
			continue
		}

		numBranchesToAssign := 1
		if maxBranchesPerEmployee > 1 {
			numBranchesToAssign = rand.Intn(maxBranchesPerEmployee) + 1
		}
		assignedBranchIndices := make(map[int]bool)
		assignedCount := 0

		for k := 0; k < numBranchesToAssign && assignedCount < len(c.branches); k++ { // Use assignedCount guard
			branchIndex := -1
			for attempts := 0; attempts < len(c.branches)*2; attempts++ { // Limit attempts
				potentialIndex := rand.Intn(len(c.branches))
				if !assignedBranchIndices[potentialIndex] && c.branches[potentialIndex].created.ID != uuid.Nil {
					branchIndex = potentialIndex
					break
				}
			}
			if branchIndex == -1 {
				t.Logf("Could not find unique valid branch for employee %s after attempts.", employee.created.Email)
				break
			}

			branch := c.branches[branchIndex]
			assignedBranchIndices[branchIndex] = true
			assignedCount++ // Increment count of successfully assigned unique branches

			t.Logf("Assigning employee %d (%s, ID: %s) to branch %d (%s, ID: %s)",
				i, employee.created.Email, employee.created.ID,
				branchIndex, branch.created.Name, branch.created.ID)

			// Use owner token for privilege when assigning employees to branches
			employee.AddBranch(t, 200, branch, &c.auth_token)
		}
	}
}

// RandomlyAssignServicesToEmployees assigns each employee 1 to N random services.
func (c *Company) RandomlyAssignServicesToEmployees(t *testing.T) {
	if len(c.employees) == 0 || len(c.services) == 0 {
		t.Log("No employees or services to assign.")
		return
	}
	maxServicesPerEmployee := 5
	if len(c.services) < maxServicesPerEmployee {
		maxServicesPerEmployee = len(c.services)
	}

	for i, employee := range c.employees {
		if employee.created.ID == uuid.Nil || employee.auth_token == "" {
			t.Logf("Skipping service assignment for employee %d (email: %s): Employee ID nil or not logged in.", i, employee.created.Email)
			continue
		}

		numServicesToAssign := 1
		if maxServicesPerEmployee > 1 {
			numServicesToAssign = rand.Intn(maxServicesPerEmployee) + 1
		}
		assignedServiceIndices := make(map[int]bool)
		assignedCount := 0

		for k := 0; k < numServicesToAssign && assignedCount < len(c.services); k++ {
			serviceIndex := -1
			for attempts := 0; attempts < len(c.services)*2; attempts++ {
				potentialIndex := rand.Intn(len(c.services))
				if !assignedServiceIndices[potentialIndex] && c.services[potentialIndex].created.ID != uuid.Nil {
					serviceIndex = potentialIndex
					break
				}
			}
			if serviceIndex == -1 {
				t.Logf("Could not find unique valid service for employee %s after attempts.", employee.created.Email)
				break
			}

			service := c.services[serviceIndex]
			assignedServiceIndices[serviceIndex] = true
			assignedCount++

			t.Logf("Assigning service %d (%s, ID: %s) to employee %d (%s, ID: %s)",
				serviceIndex, service.created.Name, service.created.ID,
				i, employee.created.Email, employee.created.ID)

			// Use Employee.AddService, assumes employee's token is used
			employee.AddService(t, 200, service)
		}
	}
}

// RandomlyAssignServicesToBranches assigns each branch 1 to N random services.
func (c *Company) RandomlyAssignServicesToBranches(t *testing.T) {
	if len(c.branches) == 0 || len(c.services) == 0 {
		t.Log("No branches or services to assign.")
		return
	}
	maxServicesPerBranch := 10
	if len(c.services) < maxServicesPerBranch {
		maxServicesPerBranch = len(c.services)
	}

	for i, branch := range c.branches {
		if branch.created.ID == uuid.Nil {
			t.Logf("Skipping service assignment for branch %d (%s): Branch ID is nil.", i, branch.created.Name)
			continue
		}

		numServicesToAssign := 1
		if maxServicesPerBranch > 1 {
			numServicesToAssign = rand.Intn(maxServicesPerBranch) + 1
		}
		assignedServiceIndices := make(map[int]bool)
		assignedCount := 0

		for k := 0; k < numServicesToAssign && assignedCount < len(c.services); k++ {
			serviceIndex := -1
			for attempts := 0; attempts < len(c.services)*2; attempts++ {
				potentialIndex := rand.Intn(len(c.services))
				if !assignedServiceIndices[potentialIndex] && c.services[potentialIndex].created.ID != uuid.Nil {
					serviceIndex = potentialIndex
					break
				}
			}
			if serviceIndex == -1 {
				t.Logf("Could not find unique valid service for branch %s after attempts.", branch.created.Name)
				break
			}

			service := c.services[serviceIndex]
			assignedServiceIndices[serviceIndex] = true
			assignedCount++

			t.Logf("Assigning service %d (%s, ID: %s) to branch %d (%s, ID: %s)",
				serviceIndex, service.created.Name, service.created.ID,
				i, branch.created.Name, branch.created.ID)

			// Use Branch.AddService method. Use owner's token implicitly (branch.auth_token or nil param).
			branch.AddService(t, 200, service, nil)
		}
	}
}

// --- Work Schedule Assignment ---

// RandomlyAssignWorkSchedules assigns a generated work schedule to each employee.
func (c *Company) RandomlyAssignWorkSchedules(t *testing.T) {
	if len(c.employees) == 0 || len(c.branches) == 0 {
		t.Log("No employees or branches for work schedule assignment.")
		return
	}

	validBranches := []*Branch{}
	for _, b := range c.branches {
		if b.created.ID != uuid.Nil {
			validBranches = append(validBranches, b)
		}
	}
	if len(validBranches) == 0 {
		t.Log("Skipping work schedule assignment: No valid branches found.")
		return
	}

	for i, employee := range c.employees {
		if employee.created.ID == uuid.Nil { // No need to check auth_token here, Update uses Company token.
			t.Logf("Skipping schedule assignment for employee %d (%s): ID nil.", i, employee.created.Email)
			continue
		}

		scheduleModel := GenerateRandomModelWorkSchedule(validBranches)
		t.Logf("Generated work schedule for employee %d (%s), referencing %d valid branch(es).", i, employee.created.Email, len(validBranches))

		// Payload format for Employee.Update
		payload := map[string]any{
			"work_schedule": []mJSON.WorkSchedule{scheduleModel},
		}

		// Call Employee.Update using owner's token (c.auth_token is implicitly used in helper via employee.company.auth_token)
		employee.Update(t, 200, payload)
		t.Logf("Attempted to update work schedule for employee %d (%s) via API.", i, employee.created.Email)

		// Optional: Refresh employee data locally if needed, though Update should handle API state
		// employee.GetById(t, 200)
	}
}

// GenerateRandomModelWorkSchedule creates a *mJSON.WorkSchedule* struct
func GenerateRandomModelWorkSchedule(validBranches []*Branch) mJSON.WorkSchedule {
	schedule := mJSON.WorkSchedule{}

	randomTimeStringHHMM := func(minHour, maxHour int) string {
		hour := minHour + rand.Intn(maxHour-minHour+1)
		if hour < 6 {
			hour = 6
		}
		if hour > 21 {
			hour = 21
		}
		minute := rand.Intn(4) * 15
		return fmt.Sprintf("%02d:%02d", hour, minute)
	}

	schedule.Monday = generateRangesForDayModel(validBranches, randomTimeStringHHMM, 0.9)
	schedule.Tuesday = generateRangesForDayModel(validBranches, randomTimeStringHHMM, 0.9)
	schedule.Wednesday = generateRangesForDayModel(validBranches, randomTimeStringHHMM, 0.9)
	schedule.Thursday = generateRangesForDayModel(validBranches, randomTimeStringHHMM, 0.9)
	schedule.Friday = generateRangesForDayModel(validBranches, randomTimeStringHHMM, 0.9)
	schedule.Saturday = generateRangesForDayModel(validBranches, randomTimeStringHHMM, 0.4)
	schedule.Sunday = generateRangesForDayModel(validBranches, randomTimeStringHHMM, 0.1)

	return schedule
}

// Helper for GenerateRandomModelWorkSchedule, returns []mJSON.WorkRange
func generateRangesForDayModel(validBranches []*Branch, randomTime func(int, int) string, workProbability float32) []mJSON.WorkRange {
	// Use global rand.Float32() - auto-seeded
	if rand.Float32() > workProbability || len(validBranches) == 0 {
		return []mJSON.WorkRange{}
	}

	ranges := []mJSON.WorkRange{}
	numRanges := 1 + rand.Intn(2) // Use global rand.Intn()

	lastEndTimeStr := "00:00"

	for r := 0; r < numRanges; r++ {
		// Use global rand.Intn()
		targetBranchHelper := validBranches[rand.Intn(len(validBranches))]

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
		endHour := startHour + durationHours
		if endHour > 22 {
			endHour = 22
		}
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

		ranges = append(ranges, mJSON.WorkRange{
			Start:    startTime,
			End:      endTime,
			BranchID: targetBranchHelper.created.ID, // Get UUID from the created model within the helper
		})
		lastEndTimeStr = endTime
	}
	return ranges
}

func (c *Company) Create(t *testing.T, status int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("POST")
	http.URL("/company")
	http.ExpectStatus(status)
	http.Header(namespace.HeadersKey.Auth, c.auth_token)
	ownerEmail := lib.GenerateRandomEmail("owner")
	ownerPswd := "1SecurePswd!"
	http.Send(DTO.CreateCompany{
		LegalName:     lib.GenerateRandomName("Company Legal Name"),
		TradeName:     lib.GenerateRandomName("Company Trade Name"),
		TaxID:         lib.GenerateRandomStrNumber(14),
		OwnerName:     lib.GenerateRandomName("Owner Name"),
		OwnerSurname:  lib.GenerateRandomName("Owner Surname"),
		OwnerEmail:    ownerEmail,
		OwnerPhone:    lib.GenerateRandomPhoneNumber(),
		OwnerPassword: ownerPswd,
	})
	http.ParseResponse(&c.created)
	owner := c.created.Employees[0]
	owner.Password = ownerPswd
	c.owner = &Employee{
		company: c,
		created: owner,
	}
}

func (c *Company) GetByName(t *testing.T, status int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("GET")
	http.URL(fmt.Sprintf("/company/name/%s", c.created.LegalName))
	http.ExpectStatus(status)
	http.Send(nil)
	http.ParseResponse(&c.created)
}

func (c *Company) GetById(t *testing.T, status int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("GET")
	http.URL(fmt.Sprintf("/company/%s", c.created.ID.String()))
	http.ExpectStatus(status)
	http.Send(nil)
	http.ParseResponse(&c.created)
}

func (c *Company) Update(t *testing.T, status int, changes map[string]any) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("PATCH")
	http.URL(fmt.Sprintf("/company/%s", c.created.ID.String()))
	http.ExpectStatus(status)
	http.Header(namespace.HeadersKey.Auth, c.auth_token)
	http.Send(changes)
}

func (c *Company) Delete(t *testing.T, status int) {
	http := (&handler.HttpClient{}).SetTest(t)
	http.Method("DELETE")
	http.URL(fmt.Sprintf("/company/%s", c.created.ID.String()))
	http.ExpectStatus(status)
	http.Header(namespace.HeadersKey.Auth, c.auth_token)
	http.Send(nil)
}
