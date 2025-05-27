package myLogger

import (
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestLokiLog_SuccessFirstTry(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://localhost:3100/loki/api/v1/push",
		httpmock.NewStringResponder(204, ""))

	logger := &Loki{}
	err := logger.Log("test message", map[string]string{"level": "info"})
	assert.NoError(t, err)
}

func TestLokiLog_SuccessAfterRetries(t *testing.T) {
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
	start := time.Now()
	err := logger.Log("retry message", map[string]string{"level": "warning"})
	duration := time.Since(start)

	assert.NoError(t, err)
	assert.Equal(t, 3, attempts)
	assert.GreaterOrEqual(t, duration, 2*time.Second)
}

func TestLokiLog_FailAfterMaxRetries(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://localhost:3100/loki/api/v1/push",
		httpmock.NewStringResponder(500, "Internal Server Error"))

	logger := &Loki{}
	err := logger.Log("fail message", map[string]string{"level": "error"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed after 3 attempts")
}
