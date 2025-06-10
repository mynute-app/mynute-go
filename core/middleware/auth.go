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
	"strings"

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

	change_schema := func(schema string) error {
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
		return nil
	}

	var user any
	var userIs string
	var schema string
	if claim.CompanyID == uuid.Nil {
		userIs = "client"
		user = &model.Client{ClientMeta: model.ClientMeta{BaseModel: model.BaseModel{ID: claim.ID}}}
		schema = "public"
	} else {
		userIs = "employee"
		user = &model.Employee{BaseModel: model.BaseModel{ID: claim.ID}}
		schema = "company"
	}

	subject_data := make(map[string]any)
	if err := change_schema(schema); err != nil {
		return err
	}
	if err := tx.Model(user).Preload(clause.Associations).First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return lib.Error.Auth.Unauthorized.WithError(fmt.Errorf("subject '%s' with ID '%s' not found at '%s' schema", userIs, claim.ID, schema))
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
	contentType := c.Get("Content-Type")

	if strings.HasPrefix(contentType, "application/json") {
		body_bytes := c.Request().Body()
		if len(body_bytes) > 0 {
			if err := json.Unmarshal(body_bytes, &body_data); err != nil {
				log.Printf("Error unmarshaling request body: %v", err)
				return lib.Error.General.AuthError.WithError(fmt.Errorf("unmarshal body error: %w", err))
			}
		} else {
			log.Println("No request body found")
		}
	} else if strings.HasPrefix(contentType, "multipart/form-data") {
		log.Println("Multipart form data not supported in this context, skipping body parsing")
		// Note: If you need to handle multipart/form-data, you would need to parse it differently.
		// Fiber provides methods for handling multipart forms, but this is not implemented here.
		// You might want to return an error or handle it as needed.
	} else if contentType == "" {
		log.Println("No Content-Type header found, assuming no body data")
		// No body data to parse, continue with an empty body_data map
	} else {
		log.Printf("Unsupported Content-Type: %s", contentType)
		return lib.Error.General.AuthError.WithError(fmt.Errorf("unsupported content type: %s", contentType))
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
	resource_data := make(map[string]any)
	var resourceFetchError error

	if resourceRefFound {
		var errUnescape error
		RequestVal, errUnescape = url.QueryUnescape(RequestVal) // Unescape value (e.g., from path/query)
		if errUnescape != nil {
			log.Printf("Error unescaping request val '%s': %v", RequestVal, errUnescape)
			return lib.Error.Auth.Unauthorized.WithError(fmt.Errorf("invalid resource identifier format"))
		}

		resource, schema, err := model.GetModelFromTableName(EndPoint.Resource.Table)
		if err != nil {
			return err
		}

		if err := change_schema(schema); err != nil {
			return err
		}

		// Fetch the resource from tx
		resourceFetchError = tx.Model(resource).Where(ResourceReference.DatabaseKey+" = ?", RequestVal).Preload(clause.Associations).Take(&resource_data).Error
		if resourceFetchError != nil {
			if resourceFetchError == gorm.ErrRecordNotFound {
				return lib.Error.Auth.Unauthorized.WithError(fmt.Errorf("resource not found for %s=%s in table %s", ResourceReference.DatabaseKey, RequestVal, EndPoint.Resource.Table))
			} else {
				log.Printf("tx Error fetching resource %s=%s in table %s: %v", ResourceReference.DatabaseKey, RequestVal, EndPoint.Resource.Table, resourceFetchError)
				return lib.Error.General.AuthError.WithError(resourceFetchError).WithError(fmt.Errorf("resource fetch error for %s=%s in table %s", ResourceReference.DatabaseKey, RequestVal, EndPoint.Resource.Table))
			}
		}
		// if resourceFetchError == nil {
		// 	// Convert fetched resource struct to map
		// 	jsonDataRes, errRes := json.Marshal(resource)
		// 	if errRes != nil {
		// 		log.Printf("Error marshaling resource: %v", errRes)
		// 		return lib.Error.General.AuthError.WithError(fmt.Errorf("marshal resource error: %w", errRes))
		// 	}
		// 	if errRes = json.Unmarshal(jsonDataRes, &resource_data); errRes != nil {
		// 		log.Printf("Error unmarshaling resource: %v", errRes)
		// 		return lib.Error.General.AuthError.WithError(fmt.Errorf("unmarshal resource error: %w", errRes))
		// 	}
		// } else {
		// 	if !errors.Is(resourceFetchError, gorm.ErrRecordNotFound) {
		// 		// Handle unexpected tx errors
		// 		log.Printf("tx Error fetching resource %s=%s in table %s: %v", ResourceReference.DatabaseKey, RequestVal, EndPoint.Resource.Table, resourceFetchError)
		// 		return lib.Error.General.AuthError.WithError(resourceFetchError)
		// 	} else {
		// 		// Handle case where resource is not found
		// 		log.Printf("Resource %s=%s not found in table %s", ResourceReference.DatabaseKey, RequestVal, EndPoint.Resource.Table)
		// 		return lib.Error.Auth.Unauthorized.WithError(fmt.Errorf("resource not found: %s=%s, in table %s", ResourceReference.DatabaseKey, RequestVal, EndPoint.Resource.Table))
		// 	}
		// }
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
		// log.Printf("DEBUG: Evaluating policy '%s' (ID: %s)", policy.Name, policy.ID) // Debug log
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
			detailedErr := fmt.Errorf("policy '%s' evaluation error: %w", policy.Name, decision.Error)
			return lib.Error.General.AuthError.WithError(detailedErr)
		} else if !decision.Allowed {
			deniedErr := fmt.Errorf("policy '%s' denied access", policy.Name)
			detailedReason := fmt.Errorf("%w. Reason: %s", deniedErr, decision.Reason) // Wrap reason
			return lib.Error.Auth.Unauthorized.WithError(detailedReason)
		}
	}

	// If loop finished and no policy explicitly denied (and no errors occurred), access is granted
	// log.Printf("INFO: Access granted for Endpoint %s %s (Subject: %s)", method, routePath, claim.ID)
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
