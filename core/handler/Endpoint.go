package handler

import (
	"agenda-kaki-go/core/config/db/model"
	"fmt"
	"log"
	"reflect"
	"runtime"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var EndpointHandlers = make(map[string]fiber.Handler)

func makeRegistryKey(path, method string) string {
	method = strings.ToUpper(method)
	return method + " " + path
}

type Endpoint struct {
	DB *gorm.DB
}

func (ep *Endpoint) Build(rPub fiber.Router, rPrv fiber.Router, mdwPub []fiber.Handler, mdwPrv []fiber.Handler) error {
	var EndPoints []*model.EndPoint
	db := ep.DB
	if err := db.Find(&EndPoints).Error; err != nil {
		return err
	}
	for _, EndPoint := range EndPoints {
		dbRouteHandler := getHandler(EndPoint.Path, EndPoint.Method)
		method := strings.ToUpper(EndPoint.Method)

		if EndPoint.IsPublic {
			handlers := append(mdwPub, dbRouteHandler)
			rPub.Add(method, EndPoint.Path, handlers...)
		} else {
			handlers := append(mdwPrv, dbRouteHandler)
			rPrv.Add(method, EndPoint.Path, handlers...)
		}
	}
	log.Println("Routes build finished!")
	return nil
}

func getHandler(path, method string) fiber.Handler {
	key := makeRegistryKey(path, method)
	return EndpointHandlers[key]
}

func (ep *Endpoint) BulkRegisterAndSave(handlers []fiber.Handler) {
	for _, h := range handlers {
		ep.RegisterHandler(h)
	}
}

func (ep *Endpoint) RegisterHandler(handler fiber.Handler) fiber.Handler {
	return handler
}

// func (r *Route) Register(endpoint *EndPoint) *ResourceToRegister {
// 	endpoint.Method = strings.ToUpper(endpoint.Method)
// 	if endpoint.Access != "public" && endpoint.Access != "private" {
// 		panic("EndPoint Route access must be either public or private")
// 	}
// 	key := makeRegistryKey(endpoint.Path, endpoint.Method)
// 	EndpointRegistry[key] = endpoint
// 	return &ResourceToRegister{
// 		Path:        endpoint.Path,
// 		Method:      endpoint.Method,
// 		Handler:     endpoint.Handler,
// 		Description: endpoint.Description,
// 		Access:      endpoint.Access,
// 		DB:          r.DB,
// 	}
// }

type ResourceToRegister struct {
	Path        string
	Method      string
	Handler     fiber.Handler
	Description string
	Access      string
	DB          *gorm.DB
	RoleAccess  []string
}

func (rr *ResourceToRegister) Save() {
	rr.Method = strings.ToUpper(rr.Method)
	var count int64
	rr.DB.
		Model(&model.EndPoint{}).
		Where("method = ? AND path = ?", rr.Method, rr.Path).
		Count(&count)
	if count > 0 {
		panic(fmt.Sprintf("EndPoint %s %s already exists", rr.Method, rr.Path))
	}
	isPublic := rr.Access == "public"
	handlerName := getHandlerName(rr.Handler)
	endpoint := model.EndPoint{
		Handler:     handlerName,
		Description: rr.Description,
		Method:      rr.Method,
		Path:        rr.Path,
		IsPublic:    isPublic,
	}
	if err := rr.DB.Create(&endpoint); err.Error != nil {
		panic(err.Error)
	}
}

func getHandlerName(fn fiber.Handler) string {
	fullName := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	// Example: "agenda-kaki-go/core/controller.branch.(*BranchController).CreateBranch-fm"
	parts := strings.Split(fullName, "/")
	short := parts[len(parts)-1]
	return short
}
