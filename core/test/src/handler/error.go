package handler

import "testing"

type testErrorHandler struct {
	t           *testing.T
	description string
}

func NewTestErrorHandler(t *testing.T) *testErrorHandler {
	return &testErrorHandler{t: t}
}

func (h *testErrorHandler) Describe(it string) *testErrorHandler {
	h.description = it
	return h
}

func (h *testErrorHandler) Test(e error) *testErrorHandler {
	if e != nil {
		h.t.Fatalf("`%s` failed with error: %v", h.description, e)
	}
	return h
}

