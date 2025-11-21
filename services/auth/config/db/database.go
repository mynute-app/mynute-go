package database

import (
	"fmt"
	"log"
	"mynute-go/services/auth/config/db/model"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	Gorm  *gorm.DB // Auth database connection
	Error error
}

type Test struct {
	*Database
	name string
}

// Connects to the main business database
func Connect() *Database {
	// Get environment variables
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	port := os.Getenv("POSTGRES_PORT")

	app_env := os.Getenv("APP_ENV")
	db_log_level := os.Getenv("POSTGRES_LOG_LEVEL")
	LogLevel := logger.Warn

	dbName := ""
	switch app_env {
	case "test":
		dbName = os.Getenv("POSTGRES_DB_TEST")
		if dbName == "" {
			dbName = "testdb"
		}
		LogLevel = logger.Info
	case "dev":
		dbName = os.Getenv("POSTGRES_DB_DEV")
		if dbName == "" {
			dbName = "devdb"
		}
		LogLevel = logger.Warn
	case "prod":
		dbName = os.Getenv("POSTGRES_DB_PROD")
		if dbName == "" {
			dbName = "maindb"
		}
	default:
		panic("APP_ENV must be one of 'dev', 'test', or 'prod'")
	}

	sslmode := "disable" // You can modify this based on your setup
	timeZone := "UTC"    // Default time_zone

	switch db_log_level {
	case "info":
		LogLevel = logger.Info
	case "error":
		LogLevel = logger.Error
	case "silent":
		LogLevel = logger.Silent
	case "warn":
		LogLevel = logger.Warn
	}

	log.Printf("Running in %s environment. Database: %s\n", app_env, dbName)

	// Build the DSN (Data Source Name)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		host, user, password, dbName, port, sslmode, timeZone)

	customGormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  LogLevel,
			Colorful:                  true,
			IgnoreRecordNotFoundError: true,
		},
	)

	gormConfig := &gorm.Config{
		Logger: customGormLogger,
	}

	// Connect to the database
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}

	// Set the connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database connection pool: ", err)
	}

	sqlDB.SetMaxIdleConns(20)                  // Max number of idle connections in the pool
	sqlDB.SetMaxOpenConns(100)                 // Max number of open connections to the database
	sqlDB.SetConnMaxLifetime(15 * time.Minute) // Max lifetime of a connection in the pool
	sqlDB.SetConnMaxIdleTime(2 * time.Second)  // Max idle time for a connection in the pool

	// NOTE: Core service does NOT connect to auth database
	// All auth operations should go through the auth service API at http://localhost:4001

	dbWrapper := &Database{
		Gorm:  db,
		Error: nil,
	}

	if app_env == "test" {
		dbWrapper.Test().Clear()
	}

	return dbWrapper
}

// Migrate runs database migrations for auth models
func (d *Database) Migrate(models []interface{}) error {
	return d.Gorm.AutoMigrate(models...)
}

// WithDB allows using a specific database connection
func (d *Database) WithDB(db *gorm.DB) *Database {
	return &Database{
		Gorm:  db,
		Error: d.Error,
	}
}

// Disconnect closes the database connection
func (d *Database) Disconnect() {
	sqlDB, err := d.Gorm.DB()
	if err != nil {
		log.Println("Failed to get database connection for closing:", err)
		return
	}
	if err := sqlDB.Close(); err != nil {
		log.Println("Failed to close database connection:", err)
	}
}

// Test returns a Test instance for testing utilities
func (d *Database) Test() *Test {
	return &Test{
		Database: d,
		name:     "auth-test",
	}
}

// Clear clears all data from auth tables (for testing)
func (t *Test) Clear() {
	log.Println("Clearing auth test database...")

	// Delete all records from auth tables
	t.Gorm.Exec("DELETE FROM admin_policies")
	t.Gorm.Exec("DELETE FROM admin_roles")
	t.Gorm.Exec("DELETE FROM admin_users")
	t.Gorm.Exec("DELETE FROM tenant_policies")
	t.Gorm.Exec("DELETE FROM tenant_general_policies")
	t.Gorm.Exec("DELETE FROM tenant_roles")
	t.Gorm.Exec("DELETE FROM tenant_general_roles")
	t.Gorm.Exec("DELETE FROM tenant_users")
	t.Gorm.Exec("DELETE FROM client_policies")
	t.Gorm.Exec("DELETE FROM client_roles")
	t.Gorm.Exec("DELETE FROM client_users")
	t.Gorm.Exec("DELETE FROM properties")
	t.Gorm.Exec("DELETE FROM resources")
	t.Gorm.Exec("DELETE FROM endpoints")

	log.Println("Auth test database cleared")
}

// InitialSeed runs initial seeding for auth database
func (d *Database) InitialSeed() {
	// Seed test endpoints and policies for authorization testing

	// Enable endpoint creation
	model.AllowEndpointCreation = true
	defer func() {
		model.AllowEndpointCreation = false
	}()

	// Create test endpoints for admin authorization
	adminUserEndpoint := model.EndPoint{
		Method:      "GET",
		Path:        "/admin/users",
		Description: "Admin user management endpoint",
	}
	d.Gorm.FirstOrCreate(&adminUserEndpoint, model.EndPoint{Method: adminUserEndpoint.Method, Path: adminUserEndpoint.Path})

	createAdminEndpoint := model.EndPoint{
		Method:      "POST",
		Path:        "/admin/users",
		Description: "Create admin user endpoint",
	}
	d.Gorm.FirstOrCreate(&createAdminEndpoint, model.EndPoint{Method: createAdminEndpoint.Method, Path: createAdminEndpoint.Path})

	companyEndpoint := model.EndPoint{
		Method:      "GET",
		Path:        "/company",
		Description: "Company management endpoint",
	}
	d.Gorm.FirstOrCreate(&companyEndpoint, model.EndPoint{Method: companyEndpoint.Method, Path: companyEndpoint.Path})

	employeeEndpoint := model.EndPoint{
		Method:      "GET",
		Path:        "/employee/:id",
		Description: "Employee profile endpoint",
	}
	d.Gorm.FirstOrCreate(&employeeEndpoint, model.EndPoint{Method: employeeEndpoint.Method, Path: employeeEndpoint.Path})

	serviceEndpoint := model.EndPoint{
		Method:      "GET",
		Path:        "/service",
		Description: "Service listing endpoint",
	}
	d.Gorm.FirstOrCreate(&serviceEndpoint, model.EndPoint{Method: serviceEndpoint.Method, Path: serviceEndpoint.Path})

	appointmentEndpoint := model.EndPoint{
		Method:      "GET",
		Path:        "/appointment",
		Description: "Appointment listing endpoint",
	}
	d.Gorm.FirstOrCreate(&appointmentEndpoint, model.EndPoint{Method: appointmentEndpoint.Method, Path: appointmentEndpoint.Path})

	createAppointmentEndpoint := model.EndPoint{
		Method:      "POST",
		Path:        "/appointment",
		Description: "Create appointment endpoint",
	}
	d.Gorm.FirstOrCreate(&createAppointmentEndpoint, model.EndPoint{Method: createAppointmentEndpoint.Method, Path: createAppointmentEndpoint.Path})

	clientProfileEndpoint := model.EndPoint{
		Method:      "GET",
		Path:        "/client/:id",
		Description: "Client profile endpoint",
	}
	d.Gorm.FirstOrCreate(&clientProfileEndpoint, model.EndPoint{Method: clientProfileEndpoint.Method, Path: clientProfileEndpoint.Path})

	// Create admin roles
	superadminRole := model.AdminRole{
		Name:        "superadmin",
		Description: "Super administrator with full access",
	}
	d.Gorm.FirstOrCreate(&superadminRole, model.AdminRole{Name: superadminRole.Name})

	supportRole := model.AdminRole{
		Name:        "support",
		Description: "Support administrator with limited access",
	}
	d.Gorm.FirstOrCreate(&supportRole, model.AdminRole{Name: supportRole.Name})

	// Create admin policies
	// Superadmin can access all admin endpoints
	d.Gorm.FirstOrCreate(&model.AdminPolicy{
		Name:        "superadmin-view-users",
		Description: "Superadmins can view admin users",
		Effect:      "Allow",
		EndPointID:  adminUserEndpoint.ID,
		Conditions:  []byte(`{"logic_type":"AND","children":[{"leaf":{"attribute":"subject.role","operator":"Equals","value":"superadmin"}}]}`),
	}, model.AdminPolicy{Name: "superadmin-view-users"})

	d.Gorm.FirstOrCreate(&model.AdminPolicy{
		Name:        "superadmin-create-users",
		Description: "Superadmins can create admin users",
		Effect:      "Allow",
		EndPointID:  createAdminEndpoint.ID,
		Conditions:  []byte(`{"logic_type":"AND","children":[{"leaf":{"attribute":"subject.role","operator":"Equals","value":"superadmin"}}]}`),
	}, model.AdminPolicy{Name: "superadmin-create-users"})

	d.Gorm.FirstOrCreate(&model.AdminPolicy{
		Name:        "admin-view-companies",
		Description: "All admins can view companies",
		Effect:      "Allow",
		EndPointID:  companyEndpoint.ID,
		Conditions:  []byte(`{"logic_type":"OR","children":[{"leaf":{"attribute":"subject.role","operator":"Equals","value":"superadmin"}},{"leaf":{"attribute":"subject.role","operator":"Equals","value":"support"}}]}`),
	}, model.AdminPolicy{Name: "admin-view-companies"})

	// Support admin cannot create other admins
	d.Gorm.FirstOrCreate(&model.AdminPolicy{
		Name:        "support-deny-create-users",
		Description: "Support admins cannot create admin users",
		Effect:      "Deny",
		EndPointID:  createAdminEndpoint.ID,
		Conditions:  []byte(`{"logic_type":"AND","children":[{"leaf":{"attribute":"subject.role","operator":"Equals","value":"support"}}]}`),
	}, model.AdminPolicy{Name: "support-deny-create-users"})

	// Create tenant policies for testing (general policies that apply to all tenants)
	d.Gorm.FirstOrCreate(&model.TenantGeneralPolicy{
		Name:        "employee-view-own-profile",
		Description: "Employees can view their own profile",
		Effect:      "Allow",
		EndPointID:  employeeEndpoint.ID,
		Conditions:  []byte(`{"logic_type":"AND","children":[{"leaf":{"attribute":"subject.user_id","operator":"Equals","resource_attribute":"resource.employee_id"}}]}`),
	}, model.TenantGeneralPolicy{Name: "employee-view-own-profile"})

	d.Gorm.FirstOrCreate(&model.TenantGeneralPolicy{
		Name:        "employee-list-services",
		Description: "Employees can list company services",
		Effect:      "Allow",
		EndPointID:  serviceEndpoint.ID,
		Conditions:  []byte(`{"logic_type":"AND","children":[{"leaf":{"attribute":"subject.company_id","operator":"Equals","resource_attribute":"query.company_id"}}]}`),
	}, model.TenantGeneralPolicy{Name: "employee-list-services"})

	// Create client policies for testing
	d.Gorm.FirstOrCreate(&model.ClientPolicy{
		Name:        "client-view-own-appointments",
		Description: "Clients can view their own appointments",
		Effect:      "Allow",
		EndPointID:  appointmentEndpoint.ID,
		Conditions:  []byte(`{"logic_type":"AND","children":[{"leaf":{"attribute":"subject.user_id","operator":"Equals","resource_attribute":"query.client_id"}}]}`),
	}, model.ClientPolicy{Name: "client-view-own-appointments"})

	d.Gorm.FirstOrCreate(&model.ClientPolicy{
		Name:        "client-create-own-appointments",
		Description: "Clients can create their own appointments",
		Effect:      "Allow",
		EndPointID:  createAppointmentEndpoint.ID,
		Conditions:  []byte(`{"logic_type":"AND","children":[{"leaf":{"attribute":"subject.user_id","operator":"Equals","resource_attribute":"body.client_id"}}]}`),
	}, model.ClientPolicy{Name: "client-create-own-appointments"})

	d.Gorm.FirstOrCreate(&model.ClientPolicy{
		Name:        "client-view-own-profile",
		Description: "Clients can view their own profile",
		Effect:      "Allow",
		EndPointID:  clientProfileEndpoint.ID,
		Conditions:  []byte(`{"logic_type":"AND","children":[{"leaf":{"attribute":"subject.user_id","operator":"Equals","resource_attribute":"resource.client_id"}}]}`),
	}, model.ClientPolicy{Name: "client-view-own-profile"})

	d.Gorm.FirstOrCreate(&model.ClientPolicy{
		Name:        "client-list-services",
		Description: "Clients can list available services",
		Effect:      "Allow",
		EndPointID:  serviceEndpoint.ID,
		Conditions:  []byte(`{"logic_type":"AND","children":[{"leaf":{"attribute":"subject.role","operator":"Equals","value":"client"}}]}`),
	}, model.ClientPolicy{Name: "client-list-services"})

	// Create test superadmin user for authentication in tests
	// Password: "test123456" (hashed)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("test123456"), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Warning: Failed to hash test admin password: %v", err)
	} else {
		testSuperAdmin := model.AdminUser{
			Email:    "test-superadmin@mynute.local",
			Password: string(hashedPassword),
			Verified: true,
		}
		d.Gorm.FirstOrCreate(&testSuperAdmin, model.AdminUser{Email: testSuperAdmin.Email})
	}

	log.Println("Auth database seeding completed")
}
