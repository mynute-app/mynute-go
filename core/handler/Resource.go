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

var ResourceRegistry = map[string]*EndPoint{}

func makeRegistryKey(path, method string) string {
	method = strings.ToUpper(method)
	return method + " " + path
}

type Route struct {
	DB *gorm.DB
}

func (r *Route) Build(rPub fiber.Router, rPrv fiber.Router, mdwPub []fiber.Handler, mdwPrv []fiber.Handler) error {
	var Resources []*model.EndPoint
	if err := r.DB.Find(&Resources).Error; err != nil {
		return err
	}
	for _, EndPoint := range Resources {
		dbRouteHandler := r.GetHandler(EndPoint.Path, EndPoint.Method)
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

func (r *Route) GetHandler(path, method string) fiber.Handler {
	key := makeRegistryKey(path, method)
	return ResourceRegistry[key].Handler
}

type EndPoint struct {
	Path        string
	Method      string
	Handler     fiber.Handler
	Description string
	Access      string
}

func (r *Route) BulkRegisterAndSave(resources []*EndPoint) {
	for _, endpoint := range resources {
		r.Register(endpoint).Save()
	}
}

func (r *Route) Register(endpoint *EndPoint) *ResourceToRegister {
	endpoint.Method = strings.ToUpper(endpoint.Method)
	if endpoint.Access != "public" && endpoint.Access != "private" {
		panic("EndPoint Route access must be either public or private")
	}
	key := makeRegistryKey(endpoint.Path, endpoint.Method)
	ResourceRegistry[key] = endpoint
	return &ResourceToRegister{
		Path:        endpoint.Path,
		Method:      endpoint.Method,
		Handler:     endpoint.Handler,
		Description: endpoint.Description,
		Access:      endpoint.Access,
		DB:          r.DB,
	}
}

type ResourceToRegister struct {
	Path        string
	Method      string
	Handler     fiber.Handler
	Description string
	Access      string
	DB          *gorm.DB
	RoleAccess  []string
}

func (rr *ResourceToRegister) SetRoleAccess(roles []string) *ResourceToRegister {
	rr.RoleAccess = roles
	return rr
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
	// Example: "agenda-kaki-go/core/controller.branch.(*BranchController).CreateBranch"
	parts := strings.Split(fullName, "/")
	short := parts[len(parts)-1]
	return short
}
