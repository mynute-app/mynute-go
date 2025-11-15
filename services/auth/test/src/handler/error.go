package handler

import (
	"fmt"
	"testing"
)

type TestErrorHandler struct {
	T            *testing.T
	description  string
	test_counter int
}

func NewTestErrorHandler(t *testing.T) *TestErrorHandler {
	return &TestErrorHandler{
		T:            t,
		test_counter: 0,
	}
}

func (te *TestErrorHandler) Describe(description string) *TestErrorHandler {
	te.description = description
	te.test_counter++
	return te
}

func (te *TestErrorHandler) Test(err error) {
	if err != nil {
		te.T.Fatalf("\n\n[Test %d] ❌ %s\nError: %v\n\n", te.test_counter, te.description, err)
	} else {
		fmt.Printf("[Test %d] ✓ %s\n", te.test_counter, te.description)
	}
}
