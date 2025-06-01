package handlerT

import "testing"

type testErrorHandler struct {
	t *testing.T
}

func NewTestErrorHandler(t *testing.T) *testErrorHandler {
	return &testErrorHandler{t: t}
}

func (h *testErrorHandler) Test(e error) {
	if e != nil {
		h.t.Errorf("Test failed with error: %v", e)
	}
}