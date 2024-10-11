package handlers

import (
	"fmt"
	"net/http"
	"testing"
)

type Tester struct {
	Entity    string
	BaseURL   string
	PostBody  map[string]interface{}
	PatchBody map[string]interface{}
	EntityID  int
}

func (test *Tester) POST(t *testing.T) {
	HTTP := HttpClient{}
	h := HTTP.SetTest(t)
	url := fmt.Sprintf("%s/%s", test.BaseURL, test.Entity)
	h.
		URL(url).
		ExpectStatus(http.StatusCreated).
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
	url := fmt.Sprintf("%s/%s/%d", test.BaseURL, test.Entity, test.EntityID)
	HTTP := HttpClient{}
	HTTP.
		SetTest(t).
		URL(url).
		Method(http.MethodPatch).
		ExpectStatus(http.StatusOK).
		Send(test.PatchBody)
}

func (test *Tester) DELETE(t *testing.T) {
	idMsg := validateId(test.EntityID)
	if idMsg != "" {
		t.Fatalf(idMsg)
	}
	url := fmt.Sprintf("%s/%s/%d", test.BaseURL, test.Entity, test.EntityID)
	HTTP := HttpClient{}
	HTTP.
		SetTest(t).
		URL(url).
		Method(http.MethodDelete).
		ExpectStatus(http.StatusNoContent).
		Send(nil)
}

func validateId(id int) string {
	if id != 0 {
		return ""
	}
	return "EntityID is empty or invalid. Either: 1. Create test method hasn't been called 2. Create test method failed 3. EntityID is not set manually."
}
