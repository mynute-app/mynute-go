package middleware

import (
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/handler"
	"log"
	"reflect"
	"runtime"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var EndpointHandlers = make(map[string]fiber.Handler)

type Endpoint struct {
	DB *handler.Gorm
}

func (ep *Endpoint) Build(r fiber.Router) error {
	var EndPoints []*model.EndPoint
	db := ep.DB.DB
	if err := db.Find(&EndPoints).Error; err != nil {
		return err
	}
	r.Use(WhoAreYou)
	for _, EndPoint := range EndPoints {
		dbRouteHandler := getHandler(EndPoint.Handler)
		method := strings.ToUpper(EndPoint.Method)
		funcs := []fiber.Handler{}
		if EndPoint.NeedsCompanyId {
			funcs = append(funcs, SaveCompanySession(db))
		} else {
			funcs = append(funcs, SavePublicSession(db))
		}
		if EndPoint.DenyUnauthorized {
			funcs = append(funcs, DenyUnauthorized)
		}
		funcs = append(funcs, dbRouteHandler)
		r.Add(method, EndPoint.Path, funcs...)
	}
	log.Println("Routes build finished!")
	return nil
}

func getHandler(handlerName string) fiber.Handler {
	return EndpointHandlers[handlerName]
}

func (ep *Endpoint) BulkRegisterHandler(handlers []fiber.Handler) {
	for _, h := range handlers {
		ep.RegisterHandler(h)
	}
}

func (ep *Endpoint) RegisterHandler(handler fiber.Handler) {
	handlerName := getHandlerName(handler)
	if handlerName == "" {
		panic("Couldn't get handler name")
	}
	EndpointHandlers[handlerName] = handler
}

func getHandlerName(fn fiber.Handler) string {
	fullName := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	// Example: "agenda-kaki-go/core/controller.(*appointment_controller).CreateAppointment-fm"
	parts := strings.Split(fullName, ".")
	if len(parts) == 0 {
		return fullName
	}
	last := parts[len(parts)-1]
	// Remove suffix like "-fm" if present
	last = strings.Split(last, "-")[0]
	return last
}
