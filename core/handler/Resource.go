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

var ResourceRegistry = map[string]*Resource{}

func makeRegistryKey(path, method string) string {
	method = strings.ToUpper(method)
	return method + " " + path
}

type Route struct {
	DB *gorm.DB
}

func (r *Route) Build(rPub fiber.Router, rPrv fiber.Router, mdwPub []fiber.Handler, mdwPrv []fiber.Handler) error {
	var Resources []*model.Resource
	if err := r.DB.Find(&Resources).Error; err != nil {
		return err
	}
	for _, Resource := range Resources {
		dbRouteHandler := r.GetHandler(Resource.Path, Resource.Method)
		method := strings.ToUpper(Resource.Method)

		if Resource.IsPublic {
			handlers := append(mdwPub, dbRouteHandler)
			rPub.Add(method, Resource.Path, handlers...)
		} else {
			handlers := append(mdwPrv, dbRouteHandler)
			rPrv.Add(method, Resource.Path, handlers...)
		}
	}
	log.Println("Routes build finished!")
	return nil
}

func (r *Route) GetHandler(path, method string) fiber.Handler {
	key := makeRegistryKey(path, method)
	return ResourceRegistry[key].Handler
}

type Resource struct {
	Path        string
	Method      string
	Handler     fiber.Handler
	Description string
	Access      string
}

func (r *Route) BulkRegisterAndSave(resources []*Resource) {
	for _, resource := range resources {
		r.Register(resource).Save()
	}
}

func (r *Route) Register(resource *Resource) *ResourceToRegister {
	resource.Method = strings.ToUpper(resource.Method)
	if resource.Access != "public" && resource.Access != "private" {
		panic("Resource Route access must be either public or private")
	}
	key := makeRegistryKey(resource.Path, resource.Method)
	ResourceRegistry[key] = resource
	return &ResourceToRegister{
		Path:        resource.Path,
		Method:      resource.Method,
		Handler:     resource.Handler,
		Description: resource.Description,
		Access:      resource.Access,
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
		Model(&model.Resource{}).
		Where("method = ? AND path = ?", rr.Method, rr.Path).
		Count(&count)
	if count > 0 {
		panic(fmt.Sprintf("Resource %s %s already exists", rr.Method, rr.Path))
	}
	isPublic := rr.Access == "public"
	handlerName := getHandlerName(rr.Handler)
	resource := model.Resource{
		Handler:     handlerName,
		Description: rr.Description,
		Method:      rr.Method,
		Path:        rr.Path,
		IsPublic:    isPublic,
	}
	if err := rr.DB.Create(&resource); err.Error != nil {
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
