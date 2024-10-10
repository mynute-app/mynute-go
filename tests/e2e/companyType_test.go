package e2e_test

import (
	"agenda-kaki-go/tests/lib"
	"fmt"
	"net/http"
	"testing"
)

const baseURL = "http://localhost:3000"

// Run the test in debug mode to avoid cache.

func TestCompanyTypeFlow(t *testing.T) {
	var companyTypeID string
	HTTP := lib.HttpClient{}

	t.Run("CreateCompanyType", func(t *testing.T) {
		h := HTTP.SetTest(t)
		body := map[string]string{
			"name": lib.GenerateRandomName("CreateCompanyType"),
		}
		url := fmt.Sprintf("%s/companyType", baseURL)
		h.
			URL(url).
			ExpectStatus(http.StatusCreated).
			Method(http.MethodPost).
			Send(body)
		companyTypeID = fmt.Sprintf("%v", h.ResBody["id"])
		t.Logf("Created companyType with ID: %s", companyTypeID)
	})

	t.Run("UpdateCompanyType", func(t *testing.T) {
		if companyTypeID == "" {
			t.Fatalf("companyTypeID is empty, previous test may have failed")
		}
		body := map[string]string{
			"name": lib.GenerateRandomName("UpdateCompanyType"),
		}
		url := fmt.Sprintf("%s/companyType/%s", baseURL, companyTypeID)
		HTTP.
			SetTest(t).
			URL(url).
			Method(http.MethodPatch).
			ExpectStatus(http.StatusOK).
			Send(body)
	})

	t.Run("DeleteCompanyType", func(t *testing.T) {
		if companyTypeID == "" {
			t.Fatalf("companyTypeID is empty, previous tests may have failed")
		}
		url := fmt.Sprintf("%s/companyType/%s", baseURL, companyTypeID)
		HTTP.
			SetTest(t).
			ExpectStatus(http.StatusNoContent).
			URL(url).
			Method(http.MethodDelete).
			Send(nil)
	})
}
