package middleware

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/handler"
	"agenda-kaki-go/core/lib"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func DenyUnauthorized(c *fiber.Ctx) error {
	// --- Get Database Session ---
	tx, err := lib.Session(c)
	if err != nil {
		return lib.Error.General.InternalError.WithError(err)
	}
	method := c.Method()
	routePath := c.Route().Path // Use routePath for consistency

	var EndPoint model.EndPoint
	if err := tx.Where("method = ? AND path = ?", method, routePath).Preload("Resource").First(&EndPoint).Error; err != nil || EndPoint.ID == uuid.Nil {
		return lib.Error.Auth.Unauthorized.WithError(fmt.Errorf("endpoint not found: %s %s", method, routePath))
	}

	// --- Get Subject ---
	auth_claims := c.Locals(namespace.RequestKey.Auth_Claims)
	claim, ok := auth_claims.(*DTO.Claims)
	if !ok || claim.ID == uuid.Nil || claim.ID.Variant() != uuid.RFC4122 || !claim.Verified {
		return lib.Error.Auth.InvalidToken
	}

	var policies []*model.PolicyRule
	PoliciesWhereClause := "end_point_id = ?"
	if err := tx.Where(PoliciesWhereClause, EndPoint.ID).Find(&policies).Error; err != nil {
		// Note: gorm.ErrRecordNotFound is handled by the len(policies) == 0 check later
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("tx Error fetching policies: %v", err)
			return lib.Error.General.AuthError.WithError(err)
		}
	}

	if len(policies) == 0 {
		companyIDStr := "NULL"
		if claim.CompanyID != uuid.Nil {
			companyIDStr = claim.CompanyID.String()
		}
		return lib.Error.Auth.Unauthorized.WithError(fmt.Errorf("no policies found for endpoint: %s %s and company: %s", method, routePath, companyIDStr))
	}

	var user any
	if claim.CompanyID == uuid.Nil {
		user = &model.ClientFull{}
	} else {
		user = &model.Employee{}
	}
	subject_data := make(map[string]any)
	if err := tx.Model(user).Preload(clause.Associations).Where("id = ?", claim.ID).Take(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return lib.Error.Auth.Unauthorized
		}
		log.Printf("tx Error fetching subject %s: %v", claim.ID, err)
		return lib.Error.General.AuthError.WithError(err)
	}
	jsonDataSub, errSub := json.Marshal(user)
	if errSub != nil {
		log.Printf("Error marshaling subject: %v", errSub)
		return lib.Error.General.AuthError.WithError(fmt.Errorf("marshal subject error: %w", errSub))
	}
	if errSub = json.Unmarshal(jsonDataSub, &subject_data); errSub != nil {
		log.Printf("Error unmarshaling subject: %v", errSub)
		return lib.Error.General.AuthError.WithError(fmt.Errorf("unmarshal subject error: %w", errSub))
	}

	// --- Parse Body ONCE ---
	body_data := make(map[string]any)
	body_bytes := c.Request().Body()
	if len(body_bytes) > 0 {
		if err := json.Unmarshal(body_bytes, &body_data); err != nil {
			log.Printf("Error unmarshaling request body: %v", err)
			return lib.Error.General.AuthError.WithError(fmt.Errorf("unmarshal body error: %w", err))
		}
	} else {
		log.Println("No request body found")
	}

	// --- Find Resource Reference Key/Value ---
	var RequestVal string
	var ResourceReference model.ResourceReference
	resourceRefFound := false
forLoop: // Label is optional but can improve readability
	for _, ref := range EndPoint.Resource.References {
		switch ref.RequestRef {
		case "query":
			val := c.Query(ref.RequestKey)
			if val != "" {
				ResourceReference = ref
				RequestVal = val
				resourceRefFound = true
				break forLoop
			}
		case "body":
			// Use the pre-parsed body_data map
			if val, exists := body_data[ref.RequestKey]; exists && val != nil && fmt.Sprintf("%v", val) != "" {
				ResourceReference = ref
				RequestVal = fmt.Sprintf("%v", val) // Convert potentially varied types to string
				resourceRefFound = true
				break forLoop
			}
		case "header":
			// Fiber's Get respects Canonical-MIME-Header-Key format
			val := c.Get(ref.RequestKey)
			if val != "" {
				ResourceReference = ref
				RequestVal = val
				resourceRefFound = true
				break forLoop
			}
		case "path":
			val := c.Params(ref.RequestKey)
			if val != "" {
				ResourceReference = ref
				RequestVal = val
				resourceRefFound = true
				break forLoop
			}
		default:
			log.Printf("Error: Invalid RequestRef '%s' in Resource %s", ref.RequestRef, EndPoint.Resource.ID)
			return fmt.Errorf("invalid request reference type: %s", ref.RequestRef) // Consider returning AuthError
		}
	}

	// Check if a reference value was actually found, otherwise resource lookup is impossible
	// Allow CREATE style endpoints where maybe no resource is expected/needed for the lookup step.
	// Need a way to know if resource lookup is required vs optional for this endpoint.
	// Let's assume for now Resource lookup IS required if references are defined.
	if EndPoint.Resource != nil && len(EndPoint.Resource.References) > 0 && !resourceRefFound {
		log.Printf("Auth Error: Could not find required resource reference value for Endpoint %s %s (Resource ID %s)", method, routePath, EndPoint.Resource.ID)
		return lib.Error.Auth.Unauthorized.WithError(fmt.Errorf("required resource identifier missing in request"))
	}

	// --- Fetch Actual Resource (based on reference found) ---
	resource_data := make(map[string]any) // Initialize as empty map
	var resourceFetchError error

	if resourceRefFound { // Only attempt fetch if we have the reference value
		var errUnescape error
		RequestVal, errUnescape = url.QueryUnescape(RequestVal) // Unescape value (e.g., from path/query)
		if errUnescape != nil {
			log.Printf("Error unescaping request val '%s': %v", RequestVal, errUnescape)
			return lib.Error.Auth.Unauthorized.WithError(fmt.Errorf("invalid resource identifier format"))
		}

		var resource any
		var schema string
		switch EndPoint.Resource.Table { // Determine model struct based on table name
		case "appointments":
			resource = &model.Appointment{}
			schema = "company"
		case "branches":
			resource = &model.Branch{}
			schema = "company"
		case "clients":
			resource = &model.ClientFull{}
			schema = "public"
		case "companies":
			resource = &model.Company{}
			schema = "public"
		case "employees":
			resource = &model.Employee{}
			schema = "company"
		case "holidays":
			resource = &model.Holiday{}
			schema = "public"
		case "policy_rules":
			resource = &model.PolicyRule{}
			schema = "public"
		case "roles":
			resource = &model.Role{}
			schema = "public"
		case "sectors":
			resource = &model.Sector{}
			schema = "public"
		case "services":
			resource = &model.Service{}
			schema = "company"
		case "subdomains":
			resource = &model.Subdomain{}
			schema = "public"
		default:
			log.Printf("Error: Invalid resource table '%s'", EndPoint.Resource.Table)
			return lib.Error.General.AuthError.WithError(fmt.Errorf("invalid resource table: %s", EndPoint.Resource.Table))
		}
		if schema == "public" {
			if err := lib.ChangeToPublicSchemaByContext(c); err != nil {
				return lib.Error.General.InternalError.WithError(err)
			}
		} else if schema == "company" {
			if err := lib.ChangeToCompanySchemaByContext(c); err != nil {
				return lib.Error.General.InternalError.WithError(err)
			}
		} else {
			return lib.Error.General.AuthError.WithError(fmt.Errorf("invalid schema type: %s", schema))
		}

		// Fetch the resource from tx
		resourceFetchError = tx.Model(resource).Where(ResourceReference.DatabaseKey+" = ?", RequestVal).Preload(clause.Associations).Take(resource).Error

		if resourceFetchError == nil {
			// Convert fetched resource struct to map
			jsonDataRes, errRes := json.Marshal(resource)
			if errRes != nil {
				log.Printf("Error marshaling resource: %v", errRes)
				return lib.Error.General.AuthError.WithError(fmt.Errorf("marshal resource error: %w", errRes))
			}
			if errRes = json.Unmarshal(jsonDataRes, &resource_data); errRes != nil {
				log.Printf("Error unmarshaling resource: %v", errRes)
				return lib.Error.General.AuthError.WithError(fmt.Errorf("unmarshal resource error: %w", errRes))
			}
		} else if !errors.Is(resourceFetchError, gorm.ErrRecordNotFound) {
			// Handle unexpected tx errors
			log.Printf("tx Error fetching resource %s=%s in table %s: %v", ResourceReference.DatabaseKey, RequestVal, EndPoint.Resource.Table, resourceFetchError)
			return lib.Error.General.AuthError.WithError(resourceFetchError)
		}
	}

	// --- Collect Path, Query, Header data ---
	path_data := make(map[string]any)
	if route := c.Route(); route != nil {
		for _, paramName := range route.Params {
			paramValue := c.Params(paramName)
			if paramValue != "" {
				path_data[paramName] = paramValue
			}
		}
	}

	query_data := make(map[string]any)
	for key, value := range c.Queries() {
		query_data[key] = value
	}

	header_data := make(map[string]any)
	for key, value := range c.GetReqHeaders() {
		// Use canonical keys provided by Fiber
		header_data[key] = value
	}

	// --- Evaluate Policies ---

	PolicyEngine := handler.NewPolicyEngine(tx)

	for _, policy := range policies {
		log.Printf("DEBUG: Evaluating policy '%s' (ID: %s)", policy.Name, policy.ID) // Debug log
		decision := PolicyEngine.CanAccess(
			subject_data,  // The subject (user) data
			resource_data, // The data fetched (or empty map)
			path_data,     // Path parameters
			body_data,     // Parsed request body (or empty map)
			query_data,    // Query parameters
			header_data,   // Request headers
			policy,        // The policy rule being checked
		)

		if decision.Error != nil {
			log.Printf("ERROR: Policy '%s' evaluation failed: %v", policy.Name, decision.Error)
			// Provide more specific error details if possible
			detailedErr := fmt.Errorf("policy '%s' evaluation error: %w", policy.Name, decision.Error)
			return lib.Error.General.AuthError.WithError(detailedErr)
		} else if !decision.Allowed {
			// Denied access
			log.Printf("INFO: Access denied by policy '%s'. Reason: %s", policy.Name, decision.Reason)
			// Create structured error for potential upstream logging/handling
			deniedErr := fmt.Errorf("policy '%s' denied access", policy.Name)
			detailedReason := fmt.Errorf("%w. Reason: %s", deniedErr, decision.Reason) // Wrap reason
			return lib.Error.Auth.Unauthorized.WithError(detailedReason)
		}
		// If allowed by this policy, continue (implicitly, maybe next policy needs check - currently loop allows if ANY policy permits)
		// TODO: Consider if ALL policies must pass, or just one. Current logic allows if *any* policy evaluates to allowed.
		// If you need ALL policies to pass, you might need different logic (e.g., track if any deny, and only allow if none deny)
		// For now, assume any allowing policy is sufficient.
		log.Printf("DEBUG: Policy '%s' allowed access.", policy.Name)
	}

	// If loop finished and no policy explicitly denied (and no errors occurred), access is granted
	log.Printf("INFO: Access granted for Endpoint %s %s (Subject: %s)", method, routePath, claim.ID)
	return c.Next()
}

func WhoAreYou(c *fiber.Ctx) error {
	authorization := c.Get(namespace.HeadersKey.Auth)
	if authorization == "" {
		return c.Next()
	}
	user, err := handler.JWT(c).WhoAreYou()
	if err != nil {
		return err
	} else if user == nil {
		return c.Next()
	}
	c.Locals(namespace.RequestKey.Auth_Claims, user)
	return c.Next()
}
