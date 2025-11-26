# Work Schedule Generation Reliability Issue - Detailed Context

## Problem Summary
The test infrastructure for `CreateCompanyRandomly()` has ~60% reliability when generating work schedules with 100% probability (`workProbability = 1.0`). Despite setting the probability to 100% to ensure every weekday has a work schedule, the complex constraint logic in `generateWorkRangeForDay()` causes intermittent failures where some weekdays end up with **zero work ranges**.

## Impact
- E2E tests fail intermittently with errors like: `"no work schedule was found that could contain the appointment from (2025-12-01T12:00:00-03:00) to (2025-12-01T13:00:00-03:00) on day (Monday)"`
- This makes the test suite unreliable and blocks CI/CD pipelines
- Developers cannot trust test results

## Root Cause Analysis

### Current Flow
1. `CreateCompanyRandomly()` creates a randomized company structure
2. `RandomlyAssignWorkScheduleToBranches()` - creates branch work schedules with `workProbability=1.0`
3. Branches are refreshed via `GetById()` to load work schedules into memory
4. `RandomlyAssignWorkScheduleToEmployees()` - creates employee work schedules with `workProbability=1.0`
5. Employees are refreshed via `GetById()` to load work schedules into memory

### The Core Problem: `generateWorkRangeForDay()` Function

**Location:** `test/src/model/company.go` lines ~745-1070

This function is called for each weekday (0-6) and attempts to generate work ranges. Even with `workProbability=1.0`, it can fail to create ANY ranges for a given day due to complex constraints:

#### Failure Points (lines that can cause early exit with `continue`):

**For Employees:**
1. **Line ~835**: `if branchRange == nil { continue }` - No suitable branch range found
2. **Line ~878**: `if maxEndHour-minStartHour < 1 { continue }` - Not enough time in branch range
3. **Line ~909**: `if startHour >= branchEndTime.Hour() { continue }` - No room after adjustment
4. **Line ~919**: `if startHour >= branchEndTime.Hour() { continue }` - Cannot create valid range
5. **Line ~1020**: `if len(commonServices) == 0 { continue }` - No common services between employee and branch

**For Branches:**
1. **Line ~960**: `if startHour >= acceptableHours.max-1 { continue }` - Start hour too late
2. **Line ~1020**: `if len(commonServices) == 0 { continue }` - No services assigned to branch

### Why This Is a Problem

The function tries to create `numRanges` work ranges per day (currently hardcoded to 1 when `workProbability=1.0`). However:

1. It loops `for range numRanges` (1 iteration when probability is 100%)
2. Inside the loop, multiple `continue` statements can skip range creation
3. If the single iteration hits a `continue`, **zero ranges are created** for that day
4. The function returns successfully without error (no validation that at least 1 range was created)
5. Tests later fail when trying to create appointments on that weekday

### Attempted Fixes (What We've Done So Far)

1. ✅ Changed `workProbability` from 0.8/0.9 to 1.0 (100%)
2. ✅ Added early return check: `if workProbability < 1.0 && rand.Float32() > workProbability { return nil }`
3. ✅ Set `numRanges = 1` when `workProbability >= 1.0` (simplified from random 1-3 ranges)
4. ✅ Added branch/employee data refresh after work schedule creation
5. ✅ Added simplified logic for branches: when `workProbability >= 1.0`, use fixed hours (6:00-23:00)
6. ✅ Added simplified logic for employees: when `workProbability >= 1.0`, use entire branch range
7. ⚠️ Added validation at end of function to check if at least 1 range was created (but this only logs error, doesn't retry)

**Current State:** Even with all these fixes, we still see ~40% failure rate because:
- The simplified logic for employees still relies on `branchRange` selection which can fail
- The `continue` statements in the middle of complex time calculations can still skip range creation
- No retry mechanism if range creation fails

## File Structure

### Key Files
- `test/src/model/company.go` - Contains all work schedule generation logic
  - `CreateCompanyRandomly()` - Main entry point (line ~160)
  - `RandomlyAssignWorkScheduleToBranches()` - Line ~680
  - `RandomlyAssignWorkScheduleToEmployees()` - Line ~650
  - `GenerateRandomBranchWorkRanges()` - Line ~712
  - `GenerateRandomEmployeeWorkRanges()` - Line ~725
  - `generateWorkRangeForDay()` - Line ~745 (THE PROBLEMATIC FUNCTION)

- `test/e2e/service_availability_client_filter_test.go` - Tests that fail intermittently

### Data Flow
```
CreateCompanyRandomly()
  ├─> Create branches
  ├─> Create employees  
  ├─> Assign services to branches
  ├─> RandomlyAssignServicesToBranches()
  ├─> RandomlyAssignWorkScheduleToBranches()
  │   ├─> For each branch: GenerateRandomBranchWorkRanges()
  │   │   └─> For each weekday (0-6): generateWorkRangeForDay(workProbability=1.0)
  │   └─> API call to create branch work schedule
  ├─> Refresh all branches (GetById to load WorkSchedule data)
  ├─> RandomlyAssignEmployeesToBranches()
  ├─> RandomlyAssignWorkScheduleToEmployees()
  │   ├─> For each employee: GenerateRandomEmployeeWorkRanges()
  │   │   └─> For each weekday (0-6): generateWorkRangeForDay(workProbability=1.0)
  │   └─> API call to create employee work schedule
  └─> Refresh all employees (GetById to load WorkSchedule data)
```

## Business Rules & Constraints

### Work Schedule Rules
1. **Branches** can work 0-7 days per week, with 1-3 ranges per day
2. **Employees** must work within branch hours (cannot work when branch is closed)
3. **Employees** can work 0-7 days per week, with 1-2 ranges per day
4. Work ranges must be **non-overlapping** and **sequential** (later range starts after earlier one ends)
5. Work ranges must have **at least 1 hour duration**
6. Employee work ranges must only offer **services that both the employee AND branch support**
7. Acceptable hours: 6:00 AM to 11:00 PM (23:00)
8. Time slots are in 15-minute increments (0, 15, 30, 45 minutes)

### Database Schema
- `public.branches` table with `work_schedule` JSONB field
- `company_schema.branch_work_ranges` table with:
  - `weekday` (0=Sunday, 1=Monday, ..., 6=Saturday)
  - `start_time` and `end_time` (stored as time-only values like "08:00")
  - `branch_id` FK
  - Many-to-many relationship with services via `branch_work_range_services`

- `company_schema.employee_work_ranges` table with:
  - Same structure as branch_work_ranges
  - `employee_id` FK
  - `branch_id` FK (which branch the employee is working at)
  - Many-to-many relationship with services via `employee_work_range_services`

## Reproducibility

### To Reproduce the Issue
```powershell
# Run tests multiple times - you'll see intermittent failures
for ($i=1; $i -le 10; $i++) { 
  Write-Host "`n=== RUN $i ==="; 
  go test ./test/e2e -run Test_ServiceAvailability_ClientFilter -v 2>&1 | 
    Select-String -Pattern "--- PASS:|--- FAIL:|--- SKIP:" 
}
```

Expected: All 10 runs should pass
Actual: ~6/10 runs pass, 4/10 fail with work schedule errors

### Common Error Messages
- `"no work schedule was found that could contain the appointment from (...) on day (Monday/Thursday/etc)"`
- `"failed to create appointment: expected status code: 200 | received status code: 400"`
- Tests SKIP when they can't find available slots (because work schedule is missing for that day)

### When It Fails
- Randomly on any weekday (Monday through Sunday)
- More common on certain days (Monday, Thursday observed more frequently)
- No consistent pattern - true randomness based on constraint collision

## Requirements for a Solution

### Must Have
1. **100% reliability** - When `workProbability=1.0`, EVERY weekday (0-6) MUST have at least 1 work range
2. **Maintain business rules** - Don't break the constraints (employees within branch hours, common services, etc.)
3. **No breaking changes** - Should work with existing API endpoints and database schema
4. **Performance** - Should not significantly slow down test execution

### Nice to Have
1. Configurable complexity - Allow tests to request "simple" schedules (all 7 days, single range) vs "complex" schedules (multiple ranges, gaps)
2. Deterministic mode - Option to use seed for reproducible schedules in debugging
3. Better error messages - If generation fails, explain which constraint was violated

## Potential Solutions to Explore

### Option 1: Guaranteed Simple Path (Recommended)
When `workProbability >= 1.0`, bypass all complex logic and use a guaranteed simple approach:
- **Branches**: Create one range per day: 8:00-20:00 (covers most appointment times)
- **Employees**: Copy the branch's range exactly for each day
- **Services**: Include ALL services from branch/employee intersection
- Only use complex random logic when probability < 1.0

**Pros:** Guaranteed to work, simple to implement
**Cons:** Loses randomization for 100% probability cases (but that's acceptable for deterministic tests)

### Option 2: Retry with Backoff
Keep current logic but add retry mechanism:
- If `generateWorkRangeForDay()` creates zero ranges, retry with relaxed constraints
- Retry 1: Remove sequential/non-overlapping requirement
- Retry 2: Use longer time windows
- Retry 3: Fall back to simple guaranteed range
- Fail hard if all retries exhausted

**Pros:** Maintains randomization, graceful degradation
**Cons:** More complex, slower on failures

### Option 3: Constraint Solver Approach
Pre-analyze constraints and generate valid schedules upfront:
- Check which services are common between employee and branch
- Calculate available time windows that satisfy all constraints
- Generate ranges from known-valid options
- Fail fast if no valid schedule is possible (report which constraint is impossible)

**Pros:** Most robust, best error reporting
**Cons:** Complex to implement, might be overkill

### Option 4: Two-Phase Generation
1. **Phase 1 (Guaranteed):** Create minimal valid schedule for each day (one simple range)
2. **Phase 2 (Enhancement):** If `workProbability < 1.0`, randomly remove days or add complexity

**Pros:** Separates "must work" from "randomization", easier to debug
**Cons:** Slightly more complex flow

## Code Samples

### Current Problematic Code Pattern
```go
// Current logic in generateWorkRangeForDay (simplified)
for range numRanges {
    // Complex logic with many continue statements
    if someCondition {
        continue // <-- Can skip range creation
    }
    if anotherCondition {
        continue // <-- Can skip range creation
    }
    // ... more conditions ...
    
    // Only reaches here if all conditions pass
    *rgs = append(*rgs, workRange)
}
// No validation that at least 1 range was created!
return nil
```

### Example of Desired Pattern
```go
// Proposed improved logic
func generateWorkRangeForDay(...) error {
    // Special handling for 100% probability
    if workProbability >= 1.0 {
        return generateGuaranteedWorkRange(ranges, validBranches, employee, weekday)
    }
    
    // Existing complex random logic for < 100% probability
    // ...
}

func generateGuaranteedWorkRange(...) error {
    // Simple, guaranteed-to-work logic
    // No continues, no complex constraints
    // Always creates exactly 1 range
    // ...
    return nil
}
```

## Additional Context

### Why This Matters
These tests are critical for validating the client appointment filtering feature, which ensures clients aren't double-booked across multiple companies. The filtering logic itself is working perfectly (we fixed a bug where `ClientAppointment.EndTime` was missing). Now we need reliable test infrastructure to continuously validate this feature.

### What Makes This Hard
The work schedule generation was designed to create realistic, varied test data by randomizing:
- Which days employees work
- How many ranges per day
- Start/end times of each range
- Which services are offered during which ranges

This complexity is valuable for testing edge cases, but it creates a huge constraint satisfaction problem that fails ~40% of the time.

### Test Requirements
For the client filtering tests to work, we need:
- At least 1 employee working on the target appointment day
- That employee must work during hours that cover the appointment time
- The employee must offer the service being booked
- All of this must be deterministic when `workProbability=1.0`

## Success Criteria

A successful solution will:
1. ✅ Pass 10/10 consecutive test runs of `Test_ServiceAvailability_ClientFilter`
2. ✅ Generate work schedules for all 7 weekdays when `workProbability=1.0`
3. ✅ Complete in reasonable time (< 10 seconds per test run)
4. ✅ Maintain backward compatibility with existing tests
5. ✅ Have clear error messages when constraints cannot be satisfied

## Next Steps

Please analyze this problem and propose a concrete solution that achieves 100% reliability for work schedule generation when `workProbability=1.0`. Focus on:
1. Identifying the exact failure modes in `generateWorkRangeForDay()`
2. Designing a solution (pick from options above or propose new one)
3. Implementing the solution with minimal changes to existing code
4. Validating with 20+ consecutive successful test runs

## Environment
- Go 1.23+
- PostgreSQL with multi-tenant schema (public + per-company schemas)
- GORM ORM
- Testing framework: standard Go testing
- Current timestamp: November 26, 2025
