package handlers

import (
	"fmt"
	"net/http"
	"testing"
)

type Tester struct {
	Entity    string
	BaseURL   string
	PostBody  map[string]string
	PatchBody map[string]string
	EntityID  string
}

func (test *Tester) POST(t *testing.T) {
	t.Logf("POST: /%s", test.Entity)
	HTTP := HttpClient{}
	h := HTTP.SetTest(t)
	url := fmt.Sprintf("%s/%s", test.BaseURL, test.Entity)
	h.
		URL(url).
		ExpectStatus(http.StatusCreated).
		Method(http.MethodPost).
		Send(test.PostBody)
	test.EntityID = fmt.Sprintf("%v", h.ResBody["id"])
}

func (test *Tester) PATCH(t *testing.T) {
	t.Logf("PATCH: /%s/%s", test.Entity, test.EntityID)
	idMsg := validateId(test.EntityID)
	if idMsg != "" {
		t.Fatalf(idMsg)
	}
	url := fmt.Sprintf("%s/%s/%s", test.BaseURL, test.Entity, test.EntityID)
	HTTP := HttpClient{}
	HTTP.
		SetTest(t).
		URL(url).
		Method(http.MethodPatch).
		ExpectStatus(http.StatusOK).
		Send(test.PatchBody)
}

func (test *Tester) DELETE(t *testing.T) {
	t.Logf("DELETE: /%s/%s", test.Entity, test.EntityID)
	idMsg := validateId(test.EntityID)
	if idMsg != "" {
		t.Fatalf(idMsg)
	}
	url := fmt.Sprintf("%s/%s/%s", test.BaseURL, test.Entity, test.EntityID)
	HTTP := HttpClient{}
	HTTP.
		SetTest(t).
		URL(url).
		Method(http.MethodDelete).
		ExpectStatus(http.StatusNoContent).
		Send(nil)
}

func validateId(id string) string {
	if id != "" {
		return ""
	}
	return "EntityID is empty. Either: 1. Create test method hasn't been called 2. Create test method failed 3. EntityID is not set manually."
}
