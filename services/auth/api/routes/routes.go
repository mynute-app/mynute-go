package routes

import (
	"mynute-go/services/auth/api/controller"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// SetupAuthRoutes configures all authentication and authorization routes
func SetupAuthRoutes(app *fiber.App, authDB *gorm.DB) {
	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"service": "auth",
		})
	})

	// Authentication routes
	auth := app.Group("/auth")
	{
		// Client authentication
		auth.Post("/client/login", controller.LoginClientByPassword)
		auth.Post("/client/login-with-code", controller.LoginClientByEmailCode)
		auth.Post("/client/send-login-code/email/:email", controller.SendClientLoginValidationCodeByEmail)

		// Employee authentication
		auth.Post("/employee/login", controller.LoginEmployeeByPassword)
		auth.Post("/employee/login-with-code", controller.LoginEmployeeByEmailCode)
		auth.Post("/employee/send-login-code/email/:email", controller.SendEmployeeLoginValidationCodeByEmail)

		// Admin authentication
		auth.Post("/admin/login", controller.AdminLoginByPassword)

		// Token validation endpoints (for business service to call)
		auth.Post("/validate", controller.ValidateToken)
		auth.Post("/validate-admin", controller.ValidateAdminToken)
	}

	// User management routes
	users := app.Group("/users")
	{
		// Client management
		users.Post("/client", controller.CreateClient)
		users.Get("/client/email/:email", controller.GetClientByEmail)
		users.Get("/client/:id", controller.GetClientById)
		users.Patch("/client/:id", controller.UpdateClientById)
		users.Delete("/client/:id", controller.DeleteClientById)

		// Employee management
		users.Post("/employee", controller.CreateEmployee)
		users.Get("/employee/email/:email", controller.GetEmployeeByEmail)
		users.Get("/employee/:id", controller.GetEmployeeById)
		users.Patch("/employee/:id", controller.UpdateEmployeeById)
		users.Delete("/employee/:id", controller.DeleteEmployeeById)

		// Admin management
		users.Get("/admin/are_there_any_superadmin", controller.AreThereAnyAdmin)
		users.Post("/admin/first_superadmin", controller.CreateFirstAdmin)
		users.Get("/admin", controller.ListAdmins)
		users.Post("/admin", controller.CreateAdmin)
		users.Get("/admin/:id", controller.GetAdminById)
		users.Patch("/admin/:id", controller.UpdateAdminById)
		users.Delete("/admin/:id", controller.DeleteAdminById)
	}

	// Policy management routes (for admin use)
	policies := app.Group("/policies")
	{
		policies.Get("/", controller.ListPolicies)
		policies.Post("/", controller.CreatePolicy)
		policies.Get("/:id", controller.GetPolicyById)
		policies.Patch("/:id", controller.UpdatePolicyById)
		policies.Delete("/:id", controller.DeletePolicyById)
	}

	// Endpoint management routes (for admin use)
	endpoints := app.Group("/endpoints")
	{
		endpoints.Get("/", controller.ListEndpoints)
		endpoints.Post("/", controller.CreateEndpoint)
		endpoints.Get("/:id", controller.GetEndpointById)
		endpoints.Patch("/:id", controller.UpdateEndpointById)
		endpoints.Delete("/:id", controller.DeleteEndpointById)
	}

	// Authorization routes (runtime access control checks)
	authorize := app.Group("/authorize")
	{
		// Check access by HTTP method and path
		authorize.Post("/by-method-and-path", controller.CheckAccess)

		// Evaluate a single policy (admin only, for testing)
		authorize.Post("/test-policy/:id", controller.EvaluatePolicy)
	}

	// Role management routes
	roles := app.Group("/roles")
	_ = roles // TODO: Remove this when endpoints are implemented
	{
		// TODO: Implement role endpoints
		// roles.Get("/", controllers.ListRoles)
		// roles.Post("/", controllers.CreateRole)
		// roles.Get("/:id", controllers.GetRole)
		// roles.Put("/:id", controllers.UpdateRole)
		// roles.Delete("/:id", controllers.DeleteRole)
	}
}
