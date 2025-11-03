package middleware

import (
	"fmt"
	"log"
	authModel "mynute-go/auth/model"
	"mynute-go/core/src/handler"
	"reflect"
	"runtime"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var EndpointControllers = make(map[string]fiber.Handler)

type Endpoint struct {
	DB *handler.Gorm
}

func (ep *Endpoint) Build(r fiber.Router) error {
	db := ep.DB.DB

	var edps []*authModel.EndPoint
	if err := db.Find(&edps).Error; err != nil {
		return err
	}

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
