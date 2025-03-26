package handler

import (
	"agenda-kaki-go/core/config/db/model"
	"log"
	"reflect"
	"runtime"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var RouteRegistry = map[string]fiber.Handler{}

func makeRegistryKey(path, method string) string {
	return method + " " + path
}

type Route struct {
	DB *gorm.DB
}

func (r *Route) Build(rPub fiber.Router, rPrv fiber.Router) error {
	var dbRoutes []*model.Route
	if err := r.DB.Find(&dbRoutes).Error; err != nil {
		return err
	}
	for _, dbRoute := range dbRoutes {
		dbRouteHandler := r.GetHandler(dbRoute.Path, dbRoute.Method)
		if dbRoute.IsPublic {
			rPub.Add(dbRoute.Method, dbRoute.Path, dbRouteHandler)
		} else {
			rPrv.Add(dbRoute.Method, dbRoute.Path, dbRouteHandler)
		}
	}
	return nil
}

func (r *Route) GetHandler(path, method string) fiber.Handler {
	key := makeRegistryKey(path, method)
	return RouteRegistry[key]
}

func (r *Route) Register(path, method, access string, handler fiber.Handler, description string) *RouteToRegister {
	if access != "public" && access != "private" {
		panic("Route access must be either public or private")
	}
	key := makeRegistryKey(path, method)
	RouteRegistry[key] = handler
	return &RouteToRegister{
		Path:        path,
		Method:      method,
		Handler:     handler,
		Description: description,
		Access:      access,
		DB:          r.DB,
	}
}

type RouteToRegister struct {
	Path        string
	Method      string
	Handler     fiber.Handler
	Description string
	Access      string
	DB          *gorm.DB
}

func (rr *RouteToRegister) Save() {
	var count int64
	rr.DB.
		Model(&model.Route{}).
		Where("method = ? AND path = ?", rr.Method, rr.Path).
		Count(&count)
	if count == 0 {
		isPublic := rr.Access == "public"
		handlerName := getHandlerName(rr.Handler)
		route := model.Route{
			Handler:     handlerName,
			Description: rr.Description,
			Method:      rr.Method,
			Path:        rr.Path,
			IsPublic:    isPublic,
		}
		if err := rr.DB.Create(&route); err.Error != nil {
			panic(err.Error)
		}
		log.Printf("Route %s %s saved on database", rr.Method, rr.Path)
	}
}

func getHandlerName(fn fiber.Handler) string {
	fullName := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	// Example: "agenda-kaki-go/core/controller.branch.(*BranchController).CreateBranch"
	parts := strings.Split(fullName, "/")
	short := parts[len(parts)-1]
	return short
}
