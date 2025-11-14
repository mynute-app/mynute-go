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

		// Tenant User authentication
		auth.Post("/tenant/login", controller.LoginTenantByPassword)
		auth.Post("/tenant/login-with-code", controller.LoginTenantByEmailCode)
		auth.Post("/tenant/send-login-code/email/:email", controller.SendTenantLoginValidationCodeByEmail)

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

		// Tenant User management
		users.Post("/tenant", controller.CreateTenantUser)
		users.Get("/tenant/email/:email", controller.GetTenantUserByEmail)
		users.Get("/tenant/:id", controller.GetTenantUserById)
		users.Patch("/tenant/:id", controller.UpdateTenantUserById)
		users.Delete("/tenant/:id", controller.DeleteTenantUserById)

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
		// Tenant policies
		policies.Get("/tenant", controller.ListTenantPolicies)
		policies.Post("/tenant", controller.CreateTenantPolicy)
		policies.Get("/tenant/:id", controller.GetTenantPolicyById)
		policies.Patch("/tenant/:id", controller.UpdateTenantPolicyById)
		policies.Delete("/tenant/:id", controller.DeleteTenantPolicyById)

		// Client policies
		policies.Get("/client", controller.ListClientPolicies)
		policies.Post("/client", controller.CreateClientPolicy)
		policies.Get("/client/:id", controller.GetClientPolicyById)
		policies.Patch("/client/:id", controller.UpdateClientPolicyById)
		policies.Delete("/client/:id", controller.DeleteClientPolicyById)

		// Admin policies
		policies.Get("/admin", controller.ListAdminPolicies)
		policies.Post("/admin", controller.CreateAdminPolicy)
		policies.Get("/admin/:id", controller.GetAdminPolicyById)
		policies.Patch("/admin/:id", controller.UpdateAdminPolicyById)
		policies.Delete("/admin/:id", controller.DeleteAdminPolicyById)
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
		// Tenant authorization
		authorize.Post("/tenant", controller.AuthorizeTenant)

		// Client authorization
		authorize.Post("/client", controller.AuthorizeClient)

		// Admin authorization
		authorize.Post("/admin", controller.AuthorizeAdmin)
	}

	// Role management routes
	roles := app.Group("/roles")
	{
		// Tenant roles
		roles.Get("/tenant", controller.ListTenantRoles)
		roles.Post("/tenant", controller.CreateTenantRole)
		roles.Get("/tenant/:id", controller.GetTenantRoleById)
		roles.Patch("/tenant/:id", controller.UpdateTenantRoleById)
		roles.Delete("/tenant/:id", controller.DeleteTenantRoleById)

		// Client roles
		roles.Get("/client", controller.ListClientRoles)
		roles.Post("/client", controller.CreateClientRole)
		roles.Get("/client/:id", controller.GetClientRoleById)
		roles.Patch("/client/:id", controller.UpdateClientRoleById)
		roles.Delete("/client/:id", controller.DeleteClientRoleById)

		// Admin roles
		roles.Get("/admin", controller.ListAdminRoles)
		roles.Post("/admin", controller.CreateAdminRole)
		roles.Get("/admin/:id", controller.GetAdminRoleById)
		roles.Patch("/admin/:id", controller.UpdateAdminRoleById)
		roles.Delete("/admin/:id", controller.DeleteAdminRoleById)
	}
}
