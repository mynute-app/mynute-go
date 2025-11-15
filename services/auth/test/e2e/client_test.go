package e2e_test

import (
	"fmt"
	"mynute-go/services/core"
	"mynute-go/services/core/api/lib"
	"mynute-go/services/core/test/src/handler"
	"mynute-go/services/core/test/src/model"
	"testing"
)

func Test_Client(t *testing.T) {
	server := core.NewServer().Run("parallel")
	defer server.Shutdown()

	tt := handler.NewTestErrorHandler(t)

	// Create first client
	client1 := &model.Client{}
	tt.Describe("Create first client").Test(client1.Create(200))

	// Test client verification
	tt.Describe("Send verification email to client").Test(client1.SendVerificationEmail(200))

	tt.Describe("Verify client email with code").Test(func() error {
		code, err := client1.GetVerificationCodeFromEmail()
		if err != nil {
			return fmt.Errorf("failed to get verification code: %w", err)
		}
		return client1.VerifyEmailByCode(200, code)
	}())

	// Test client login with password
	tt.Describe("Login client with password").Test(client1.LoginByPassword(200, client1.Password))

	// Test invalid password login
	tt.Describe("Login client with wrong password fails").Test(client1.LoginByPassword(401, "wrong-password"))

	// Test login with email code
	tt.Describe("Send login code to client email").Test(client1.SendLoginValidationCodeByEmail(200))

	tt.Describe("Login client with email code").Test(func() error {
		code, err := client1.GetValidationCodeFromEmail()
		if err != nil {
			return fmt.Errorf("failed to get login code: %w", err)
		}
		return client1.LoginByEmailCode(200, code)
	}())

	// Create second client for authorization tests
	client2 := &model.Client{}
	tt.Describe("Create second client").Test(client2.Create(200))

	tt.Describe("Verify second client email").Test(func() error {
		if err := client2.SendVerificationEmail(200); err != nil {
			return err
		}
		code, err := client2.GetVerificationCodeFromEmail()
		if err != nil {
			return err
		}
		return client2.VerifyEmailByCode(200, code)
	}())

	tt.Describe("Login second client").Test(client2.LoginByPassword(200, client2.Password))

	// Test client authorization - can access own profile
	tt.Describe("Client can access own profile").Test(func() error {
		authReq := map[string]interface{}{
			"method": "GET",
			"path":   "/client/:id",
			"subject": map[string]interface{}{
				"user_id": client1.Created.UserID.String(),
				"role":    "client",
			},
			"resource": map[string]interface{}{
				"client_id": client1.Created.UserID.String(),
			},
			"path_params": map[string]interface{}{
				"id": client1.Created.UserID.String(),
			},
		}

		response, err := makeAuthorizationRequest("client", authReq, "")
		if err != nil {
			return err
		}

		allowed, ok := response["allowed"].(bool)
		if !ok || !allowed {
			return fmt.Errorf("expected authorization to be allowed, got response: %v", response)
		}

		return nil
	}())

	// Test client authorization - cannot access another client's profile
	tt.Describe("Client cannot access other client profile").Test(func() error {
		authReq := map[string]interface{}{
			"method": "GET",
			"path":   "/client/:id",
			"subject": map[string]interface{}{
				"user_id": client1.Created.UserID.String(),
				"role":    "client",
			},
			"resource": map[string]interface{}{
				"client_id": client2.Created.UserID.String(),
			},
			"path_params": map[string]interface{}{
				"id": client2.Created.UserID.String(),
			},
		}

		response, err := makeAuthorizationRequest("client", authReq, "")
		if err != nil {
			return err
		}

		allowed, ok := response["allowed"].(bool)
		if !ok {
			return fmt.Errorf("unexpected response format: %v", response)
		}

		// Should be denied (clients can only access their own data)
		if allowed {
			return fmt.Errorf("expected authorization to be denied for accessing other client")
		}

		return nil
	}())

	// Test client authorization - can create own appointments
	tt.Describe("Client can create own appointments").Test(func() error {
		authReq := map[string]interface{}{
			"method": "POST",
			"path":   "/appointment",
			"subject": map[string]interface{}{
				"user_id": client1.Created.UserID.String(),
				"role":    "client",
			},
			"body": map[string]interface{}{
				"client_id": client1.Created.UserID.String(),
			},
		}

		response, err := makeAuthorizationRequest("client", authReq, "")
		if err != nil {
			return err
		}

		allowed, ok := response["allowed"].(bool)
		if !ok || !allowed {
			return fmt.Errorf("expected authorization to be allowed for creating own appointment, got: %v", response)
		}

		return nil
	}())

	// Test client authorization - cannot create appointments for another client
	tt.Describe("Client cannot create appointments for other clients").Test(func() error {
		authReq := map[string]interface{}{
			"method": "POST",
			"path":   "/appointment",
			"subject": map[string]interface{}{
				"user_id": client1.Created.UserID.String(),
				"role":    "client",
			},
			"body": map[string]interface{}{
				"client_id": client2.Created.UserID.String(), // Different client!
			},
		}

		response, err := makeAuthorizationRequest("client", authReq, "")
		if err != nil {
			return err
		}

		allowed, ok := response["allowed"].(bool)
		if !ok {
			return fmt.Errorf("unexpected response format: %v", response)
		}

		if allowed {
			return fmt.Errorf("expected authorization to be denied for creating appointment for other client")
		}

		return nil
	}())

	// Test client profile update
	tt.Describe("Update client profile").Test(func() error {
		newName := lib.GenerateRandomName("Updated Client")
		changes := map[string]any{
			"name": newName,
		}
		if err := client1.Update(200, changes); err != nil {
			return err
		}
		if client1.Created.Name != newName {
			return fmt.Errorf("expected name to be updated to %s, got %s", newName, client1.Created.Name)
		}
		return nil
	}())

	// Test password reset
	tt.Describe("Reset client password by email").Test(func() error {
		if err := client1.SendPasswordResetEmail(200); err != nil {
			return err
		}
		newPassword, err := client1.GetNewPasswordFromEmail()
		if err != nil {
			return err
		}
		client1.Password = newPassword
		return client1.LoginByPassword(200, newPassword)
	}())

	// Test client authorization - list available services (public endpoint)
	tt.Describe("Client can list available services").Test(func() error {
		authReq := map[string]interface{}{
			"method": "GET",
			"path":   "/service",
			"subject": map[string]interface{}{
				"user_id": client1.Created.UserID.String(),
				"role":    "client",
			},
		}

		response, err := makeAuthorizationRequest("client", authReq, "")
		if err != nil {
			return err
		}

		allowed, ok := response["allowed"].(bool)
		if !ok || !allowed {
			return fmt.Errorf("expected authorization to be allowed for listing services, got: %v", response)
		}

		return nil
	}())

	// Test client can view their own appointments
	tt.Describe("Client can view own appointments").Test(func() error {
		authReq := map[string]interface{}{
			"method": "GET",
			"path":   "/appointment",
			"subject": map[string]interface{}{
				"user_id": client1.Created.UserID.String(),
				"role":    "client",
			},
			"query": map[string]interface{}{
				"client_id": client1.Created.UserID.String(),
			},
		}

		response, err := makeAuthorizationRequest("client", authReq, "")
		if err != nil {
			return err
		}

		allowed, ok := response["allowed"].(bool)
		if !ok || !allowed {
			return fmt.Errorf("expected authorization to be allowed for viewing own appointments, got: %v", response)
		}

		return nil
	}())

	// Test client cannot view another client's appointments
	tt.Describe("Client cannot view other client appointments").Test(func() error {
		authReq := map[string]interface{}{
			"method": "GET",
			"path":   "/appointment",
			"subject": map[string]interface{}{
				"user_id": client1.Created.UserID.String(),
				"role":    "client",
			},
			"query": map[string]interface{}{
				"client_id": client2.Created.UserID.String(), // Different client!
			},
		}

		response, err := makeAuthorizationRequest("client", authReq, "")
		if err != nil {
			return err
		}

		allowed, ok := response["allowed"].(bool)
		if !ok {
			return fmt.Errorf("unexpected response format: %v", response)
		}

		if allowed {
			return fmt.Errorf("expected authorization to be denied for viewing other client's appointments")
		}

		return nil
	}())

	// Test client deletion
	tt.Describe("Delete client").Test(client2.Delete(200))

	// Test deleted client cannot login
	tt.Describe("Deleted client cannot login").Test(client2.LoginByPassword(401, client2.Password))

	// Test retrieving client by email
	tt.Describe("Get client by email").Test(client1.GetByEmail(200))

	// Test client profile image upload
	tt.Describe("Upload client profile image").Test(func() error {
		return client1.UploadImages(200, map[string][]byte{
			"profile": []byte("fake-image-data"),
		}, nil)
	}())
}
