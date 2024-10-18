package handlers

import (
	"fmt"
	"net/http"
	"testing"
)

type Tester struct {
	Entity         string
	BaseURL        string
	RelatedPath    string
	PostBody       map[string]interface{}
	PatchBody      map[string]interface{}
	EntityID       int
	expectedStatus int
}

func (test *Tester) ExpectStatus(status int) *Tester {
	test.expectedStatus = status
	return test
}

func (test *Tester) POST(t *testing.T) {
	HTTP := HttpClient{}
	h := HTTP.SetTest(t)
	url := fmt.Sprintf("%s/%s", test.BaseURL, test.RelatedPath)
	h.
		URL(url).
		ExpectStatus(test.expectedStatus).
		Method(http.MethodPost).
		Send(test.PostBody)
	idFloat, ok := h.ResBody["id"].(float64)
	if !ok {
		t.Fatalf("failed to assert EntityID as float64")
	}
	test.EntityID = int(idFloat)
}

func (test *Tester) PATCH(t *testing.T) {
	idMsg := validateId(test.EntityID)
	if idMsg != "" {
		t.Fatalf(idMsg)
	}
	url := fmt.Sprintf("%s/%s/%d", test.BaseURL, test.RelatedPath, test.EntityID)
	HTTP := HttpClient{}
	HTTP.
		SetTest(t).
		URL(url).
		Method(http.MethodPatch).
		ExpectStatus(test.expectedStatus).
		Send(test.PatchBody)
}

func (test *Tester) DELETE(t *testing.T) {
	idMsg := validateId(test.EntityID)
	if idMsg != "" {
		t.Fatalf(idMsg)
	}
	url := fmt.Sprintf("%s/%s/%d", test.BaseURL, test.RelatedPath, test.EntityID)
	HTTP := HttpClient{}
	HTTP.
		SetTest(t).
		URL(url).
		Method(http.MethodDelete).
		ExpectStatus(test.expectedStatus).
		Send(nil)
}

func (test *Tester) ForceDELETE(t *testing.T) {
	idMsg := validateId(test.EntityID)
	if idMsg != "" {
		t.Fatalf(idMsg)
	}
	url := fmt.Sprintf("%s/%s/%d/force", test.BaseURL, test.RelatedPath, test.EntityID)
	HTTP := HttpClient{}
	HTTP.
		SetTest(t).
		URL(url).
		Method(http.MethodDelete).
		ExpectStatus(test.expectedStatus).
		Send(nil)
}

func validateId(id int) string {
	if id != 0 {
		return ""
	}
	return "EntityID is empty or invalid. Either: 1. Create test method hasn't been called 2. Create test method failed 3. EntityID is not set manually."
}
