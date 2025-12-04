package middleware

import (
	"fmt"
	"log"
	"mynute-go/core/src/config/db/model"
	"mynute-go/core/src/handler"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

var EndpointControllers = make(map[string]fiber.Handler)

type Endpoint struct {
	DB *handler.Gorm
}

// BuildWithRetry attempts to build endpoints with retry logic
// Waits for migrations to complete instead of crashing
func (ep *Endpoint) BuildWithRetry(r fiber.Router, maxRetries int, retryDelay int) error {
	db := ep.DB.DB

	for attempt := 1; attempt <= maxRetries; attempt++ {
		// Check if endpoints table exists
		var tableExists bool
		err := db.Raw("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'endpoints')").Scan(&tableExists).Error
		if err != nil {
			log.Printf("⚠️  Attempt %d/%d: Failed to check endpoints table: %v", attempt, maxRetries, err)
			if attempt < maxRetries {
				time.Sleep(time.Duration(retryDelay) * time.Second)
				continue
			}
			return fmt.Errorf("failed to check if endpoints table exists after %d attempts: %w", maxRetries, err)
		}

		if !tableExists {
			if attempt == 1 {
				log.Println("⚠️  Endpoints table does not exist yet. Waiting for migrations...")
				log.Println("   If migrations haven't been run, execute: docker exec <container> ./migrate-tool up")
			}
			log.Printf("   Attempt %d/%d: Waiting %d seconds for migrations to complete...", attempt, maxRetries, retryDelay)
			if attempt < maxRetries {
				time.Sleep(time.Duration(retryDelay) * time.Second)
				continue
			}
			return fmt.Errorf("endpoints table not found after %d attempts - please run migrations manually", maxRetries)
		}

		// Table exists, try to load endpoints
		var edps []*model.EndPoint
		if err := db.Find(&edps).Error; err != nil {
			log.Printf("⚠️  Attempt %d/%d: Failed to load endpoints: %v", attempt, maxRetries, err)
			if attempt < maxRetries {
				time.Sleep(time.Duration(retryDelay) * time.Second)
				continue
			}
			return fmt.Errorf("failed to load endpoints after %d attempts: %w", maxRetries, err)
		}

		log.Printf("✅ Successfully loaded %d endpoints", len(edps))

		// Build routes with loaded endpoints
		return ep.buildRoutes(r, edps)
	}

	return fmt.Errorf("unexpected error in BuildWithRetry")
}

// Build attempts to build endpoints once (kept for backward compatibility)
func (ep *Endpoint) Build(r fiber.Router) error {
	db := ep.DB.DB

	// Check if endpoints table exists before querying
	var tableExists bool
	err := db.Raw("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'endpoints')").Scan(&tableExists).Error
	if err != nil {
		return fmt.Errorf("failed to check if endpoints table exists: %w", err)
	}

	if !tableExists {
		log.Println("⚠️  Endpoints table does not exist yet. Please run migrations first.")
		log.Println("   Run: docker compose -f docker-compose.prod.yml run --rm migrate")
		return fmt.Errorf("endpoints table not found - run migrations first")
	}

	var edps []*model.EndPoint
	if err := db.Find(&edps).Error; err != nil {
		return err
	}

	return ep.buildRoutes(r, edps)
}

// buildRoutes is the internal method that actually builds the routes
func (ep *Endpoint) buildRoutes(r fiber.Router, edps []*model.EndPoint) error {
	db := ep.DB.DB

	r.Use(WhoAreYou)

	for _, e := range edps {
		handlers := []fiber.Handler{}

		// Sessão
		if e.NeedsCompanyId {
			handlers = append(handlers, SaveCompanySession(db))
		} else {
			handlers = append(handlers, SavePublicSession(db))
		}

		// Autorização
		if e.DenyUnauthorized {
			handlers = append(handlers, DenyUnauthorized)
		}

		// Schema
		if e.NeedsCompanyId {
			handlers = append(handlers, ChangeToCompanySchema)
		} else {
			handlers = append(handlers, ChangeToPublicSchema)
		}

		controller, err := ep.GetControllerFnc(e.ControllerName)
		if err != nil {
			panic(err)
		}

		// Handler final
		handlers = append(handlers, controller)

		method := strings.ToUpper(e.Method)
		r.Add(method, e.Path, handlers...)
	}

	log.Println("Routes build finished!")
	return nil
}

func (ep *Endpoint) GetControllerFnc(ctrlName string) (fiber.Handler, error) {
	if ctrlName == "" {
		return nil, fmt.Errorf("controller name is empty")
	}
	if controller, ok := EndpointControllers[ctrlName]; !ok {
		return nil, fmt.Errorf("controller '%s' not found", ctrlName)
	} else {
		return controller, nil
	}
}

func (ep *Endpoint) BulkRegisterHandler(handlers []fiber.Handler) {
	for _, h := range handlers {
		ep.RegisterHandler(h)
	}
}

func (ep *Endpoint) RegisterHandler(handler fiber.Handler) {
	handlerName := getEndpointControllerName(handler)
	if handlerName == "" {
		panic("Couldn't get handler name")
	}
	EndpointControllers[handlerName] = handler
}

func getEndpointControllerName(fn fiber.Handler) string {
	fullName := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	// Example: "mynute-go/core/src/controller.(*appointment_controller).CreateAppointment-fm"
	parts := strings.Split(fullName, ".")
	if len(parts) == 0 {
		return fullName
	}
	last := parts[len(parts)-1]
	// Remove suffix like "-fm" if present
	last = strings.Split(last, "-")[0]
	return last
}
