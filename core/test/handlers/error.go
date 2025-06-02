package handlerT

import "testing"

type testErrorHandler struct {
	t *testing.T
}

func NewTestErrorHandler(t *testing.T) *testErrorHandler {
	return &testErrorHandler{t: t}
}

func (h *testErrorHandler) Test(e error, it string) {
	if e != nil {
		h.t.Fatalf("%s failed with error: %v", it, e)
	}
}