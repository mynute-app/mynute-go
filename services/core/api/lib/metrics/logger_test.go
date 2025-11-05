package myLogger

import (
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestLokiLogStructured_SuccessFirstTry(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://localhost:3100/loki/api/v1/push",
		httpmock.NewStringResponder(204, ""))

	logger := &Loki{}

	labels := map[string]string{
		"app":   "test-app",
		"level": "info",
		"type":  "test",
	}

	body := map[string]any{
		"message": "structured log test",
		"path":    "/test",
		"method":  "GET",
	}

	err := logger.LogV13(labels, body)
	assert.NoError(t, err)
}

func TestLokiLogStructured_SuccessAfterRetries(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	attempts := 0
	httpmock.RegisterResponder("POST", "http://localhost:3100/loki/api/v1/push",
		func(req *http.Request) (*http.Response, error) {
			attempts++
			if attempts < 3 {
				return httpmock.NewStringResponse(500, "Internal Server Error"), nil
			}
			return httpmock.NewStringResponse(204, ""), nil
		},
	)

	logger := &Loki{}

	labels := map[string]string{
		"app":   "retry-app",
		"level": "warning",
		"type":  "test",
	}

	body := map[string]any{
		"message": "retry structured log",
		"step":    "step-1",
	}

	start := time.Now()
	err := logger.LogV13(labels, body)
	duration := time.Since(start)

	assert.NoError(t, err)
	assert.Equal(t, 3, attempts)
	assert.GreaterOrEqual(t, duration, 2*time.Second)
}

func TestLokiLogStructured_FailAfterMaxRetries(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://localhost:3100/loki/api/v1/push",
		httpmock.NewStringResponder(500, "Internal Server Error"))

	logger := &Loki{}

	labels := map[string]string{
		"app":   "fail-app",
		"level": "error",
		"type":  "test",
	}

	body := map[string]any{
		"message": "this should fail",
	}

	err := logger.LogV13(labels, body)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed after 3 attempts")
}

