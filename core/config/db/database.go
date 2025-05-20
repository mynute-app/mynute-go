package database

import (
	"agenda-kaki-go/core/lib"
	"fmt"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

type Database struct {
	Gorm  *gorm.DB
	Error error
}

type Test struct {
	*Database
	name string
}

// Connects to the database
func Connect() *Database {
	// Get environment variables
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	port := os.Getenv("POSTGRES_PORT")
	app_env := os.Getenv("APP_ENV")
	sslmode := "disable" // You can modify this based on your setup
	timeZone := "UTC"    // Default timezone
	LogLevel := logger.Warn

	if app_env == "test" {
		dbName = os.Getenv("POSTGRES_DB_TEST")
		LogLevel = logger.Info
	} else if app_env == "dev" {
		LogLevel = logger.Info
	} else if app_env != "production" && app_env != "dev" {
		log.Fatalf("Invalid APP_ENV: %s", app_env)
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

	return &Database{Gorm: db}
}

// Migrate the database schema
func (db *Database) Migrate(models any) *Database {
	if db.Error != nil {
		return db
	}
	if models == nil {
		db.Error = fmt.Errorf("models cannot be nil at Migrate function")
		return db
	}
	// Make sure the models is a slice of pointers to structs
	if reflect.TypeOf(models).Kind() != reflect.Slice {
		db.Error = fmt.Errorf("models must be a slice of pointers to structs at Migrate function")
		return db
	}
	for i := range reflect.ValueOf(models).Len() {
		newModel := reflect.ValueOf(models).Index(i).Interface()
		if newModel == nil {
			db.Error = fmt.Errorf("model at index %d is nil at Migrate function", i)
			return db
		}
		// Check if the model is a pointer to a struct
		if reflect.TypeOf(newModel).Kind() != reflect.Ptr {
			db.Error = fmt.Errorf("model at index %d is not a pointer to a struct at Migrate function", i)
			return db
		}

		// Check if the model is a struct
		if reflect.TypeOf(newModel).Elem().Kind() != reflect.Struct {
			db.Error = fmt.Errorf("model at index %d is not a struct at Migrate function", i)
			return db
		}

		if err := db.Gorm.AutoMigrate(newModel); err != nil {
			db.Error = fmt.Errorf("failed to migrate model at index %d: %v", i, err)
			return db
		}
	}

	log.Println("Migration finished!")
	return db
}

func (db *Database) Seed(name string, models any, query string, keys []string) *Database {
	if db.Error != nil {
		return db
	}

	if models == nil {
		db.Error = fmt.Errorf("models cannot be nil. seeding name: %s", name)
		return db
	}
	modelsVal := reflect.ValueOf(models) // Use modelsVal consistently
	modelsTyp := modelsVal.Type()        // Use modelsTyp consistently

	if modelsTyp.Kind() != reflect.Slice {
		db.Error = fmt.Errorf("models must be a slice. seeding name: %s. Got: %s", name, modelsTyp.Kind())
		return db
	}

	modelsLen := modelsVal.Len()

	if modelsLen == 0 {
		log.Printf("models slice is empty, nothing to seed for: %s", name)
		return db
	}

	// Check the type of the slice elements *once*
	elemType := modelsTyp.Elem()
	if elemType.Kind() != reflect.Ptr {
		db.Error = fmt.Errorf("models slice elements must be pointers. seeding name: %s. Got element kind: %s", name, elemType.Kind())
		return db
	}
	if elemType.Elem().Kind() != reflect.Struct {
		db.Error = fmt.Errorf("models slice elements must be pointers to structs. seeding name: %s. Got pointer to: %s", name, elemType.Elem().Kind())
		return db
	}

	tx := db.Gorm

	// Iterate over the slice of models
	for i := range modelsLen { // Correct loop condition
		newModelVal := modelsVal.Index(i)
		if newModelVal.IsNil() {
			db.Error = fmt.Errorf("model at index %d is a nil pointer. seeding name: %s", i, name)
			return db
		}
		newModel := newModelVal.Interface()
		underlyingStructType := newModelVal.Elem().Type()         // Get the struct type that newModel points to
		oldModel := reflect.New(underlyingStructType).Interface() // Create a new pointer to an instance of that struct type

		args := make([]any, 0, len(keys)) // Pre-allocate capacity
		for _, key := range keys {
			field := newModelVal.Elem().FieldByName(key) // Operate on newModelVal.Elem() which is the struct
			if !field.IsValid() {
				db.Error = fmt.Errorf("field '%s' does not exist in model %s at index %d. seeding name: %s", key, underlyingStructType.Name(), i, name)
				return db
			}
			args = append(args, field.Interface())
		}

		if err := tx.Where(query, args...).First(oldModel).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if errCreate := tx.Create(newModel).Error; errCreate != nil {
					db.Error = fmt.Errorf("failed to create model %s at index %d: %v. seeding name: %s", underlyingStructType.Name(), i, errCreate, name)
					return db
				}
			} else {
				db.Error = fmt.Errorf("failed to check if model %s at index %d exists: %v. seeding name: %s", underlyingStructType.Name(), i, err, name)
				return db
			}
		} else {
			// Model exists, update it
			if errUpdate := tx.Model(oldModel).Updates(newModel).Error; errUpdate != nil {
				db.Error = fmt.Errorf("failed to update model %s at index %d: %v. seeding name: %s", underlyingStructType.Name(), i, errUpdate, name)
				return db
			}
		}
	}

	return db
}

// Close connection to the database
func (db *Database) Disconnect() {
	sqlDB, err := db.Gorm.DB()
	if err != nil {
		log.Fatal("Failed to close the database: ", err)
	}
	sqlDB.Close()
}

func (db *Database) Test() *Test {
	dbName := os.Getenv("POSTGRES_DB_NAME")
	app_env := os.Getenv("APP_ENV")
	dbName = fmt.Sprintf("%s-%s", dbName, app_env)
	return &Test{Database: db, name: dbName}
}

func (t *Test) Clear() {
	if os.Getenv("APP_ENV") != "test" {
		return
	}

	// Step 1: Drop all schemas except 'public'
	dropSchemasSQL := `
		DO $$ DECLARE
			schema_name text;
		BEGIN
			FOR schema_name IN
				SELECT nspname FROM pg_namespace
				WHERE nspname NOT IN ('pg_catalog', 'information_schema', 'public')
				  AND nspname NOT LIKE 'pg_toast%'
			LOOP
				EXECUTE format('DROP SCHEMA IF EXISTS %I CASCADE', schema_name);
			END LOOP;
		END $$;
	`

	// Step 2: Drop and recreate 'public' just in case
	resetPublicSQL := `
		DROP SCHEMA IF EXISTS public CASCADE;
		CREATE SCHEMA public;
	`

	// Execute both
	if err := t.Gorm.Exec(dropSchemasSQL).Error; err != nil {
		log.Fatalf("Failed to drop non-public schemas: %v", err)
	}
	if err := t.Gorm.Exec(resetPublicSQL).Error; err != nil {
		log.Fatalf("Failed to reset public schema: %v", err)
	}

	log.Printf("Erased all schemas on %s database.\n", t.name)
}

/*
Handle transaction rollback and commit.
It should be deferred after starting a transaction.
*/
func Defer(tx *gorm.DB) {
	if r := recover(); r != nil {
		_ = tx.Rollback()
		if err, ok := r.(error); ok {
			log.Printf("ContextTransaction rolled back due to panic: %v", err)
		} else {
			log.Println("ContextTransaction rolled back due to unknown panic.")
		}
		panic(r) // re-throw
	} else {
		// Commit and log error if commit fails
		if err := tx.Commit().Error; err != nil {
			_ = tx.Rollback()
			log.Printf("Commit failed, transaction rolled back: %v", err)
		}
	}
}

/*
 * Callback function to handle transaction rollback and commit.
 It should be deferred after starting a transaction.
 * @param tx *gorm.DB - The transaction session
*/
// @return func() - The function to be defered
func DeferCallback(tx *gorm.DB) func() {
	return func() {
		Defer(tx)
	}
}

func Transaction(db *gorm.DB) (*gorm.DB, func(), error) {
	tx := db.Begin()
	if tx.Error != nil {
		return nil, nil, lib.Error.General.DatabaseError.WithError(tx.Error)
	}
	return tx, DeferCallback(tx), nil
}

/*
 * Opens a transaction session for the current request.
 Recomended when you need to perform multiple database operations dependant of each other.
 * @return *gorm.DB - The transaction session
 * @return func() - The function to end the transaction
 * @return error - The error if any
*/
// @example
//	tx, end, err := database.ContextTransaction(c)
//	defer end()
//	if err != nil {
//		return err
//	}
//	// Then use the transaction session (tx) for your database operations
func ContextTransaction(c *fiber.Ctx) (*gorm.DB, func(), error) {
	session, err := lib.Session(c)
	if err != nil {
		return nil, nil, err
	}
	return Transaction(session)
}

// Locks the record for update using the given transaction and model.
// It uses the "UPDATE" locking strength to prevent other transactions
// from modifying the record until the current transaction is completed.
// It will also retrieve the record with the specified ID from the database.
func LockForUpdate(tx *gorm.DB, model any, key, val string) error {
	where := fmt.Sprintf("%s = ?", key)
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(where, val).First(model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return lib.Error.General.RecordNotFound.WithError(err)
		}
		return lib.Error.General.InternalError.WithError(err)
	}
	return nil
}
