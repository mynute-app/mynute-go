package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

type HttpClient struct{}

type httpActions struct {
	Error          string
	Status         int
	ResBody        map[string]interface{}
	url            string
	method         string
	expectedStatus int
	test           *testing.T
	httpRes        *http.Response
}

func (c *HttpClient) SetTest(t *testing.T) *httpActions {
	h := &httpActions{}
	h.test = t
	h.expectedStatus = 0
	return h
}

func (h *httpActions) Method(method string) *httpActions {
	h.method = method
	return h
}

func (h *httpActions) URL(url string) *httpActions {
	h.url = url
	return h
}

func (h *httpActions) ExpectStatus(status int) *httpActions {
	h.expectedStatus = status
	return h
}

func (h *httpActions) Clear() {
	h.Error = ""
	h.Status = 0
	h.ResBody = nil
}

// Send sends a request to the specified URL with the specified body
// and returns the response. If an error occurs, it sets the error message.
// The response body is stored in ResBody, and the status code is stored in Status.
// It will defer res.Body.Close() to ensure the response body is closed.
func (h *httpActions) Send(body map[string]string) *httpActions {
	h.Error = ""
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		h.Error = fmt.Sprintf("failed to marshal payload: %v", err)
		h.test.Fatalf(h.Error)
		return nil
	}
	bodyBuffer := bytes.NewBuffer(bodyBytes)
	req, err := http.NewRequest(h.method, h.url, bodyBuffer)
	if err != nil {
		h.Error = fmt.Sprintf("failed to create request: %v", err)
		h.test.Fatalf(h.Error)
		return nil
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		h.Error = fmt.Sprintf("failed to send request: %v", err)
		h.test.Fatalf(h.Error)
		return nil
	}
	defer res.Body.Close()
	var response map[string]interface{}
	if res.ContentLength != 0 {
		if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
			h.Error = fmt.Sprintf("failed to decode response: %v", err)
			h.test.Fatalf(h.Error)
			return nil
		}
	}
	if h.expectedStatus != 0 && res.StatusCode != h.expectedStatus {
		h.Error = fmt.Sprintf("expected status code %d, got %d", h.expectedStatus, res.StatusCode)
		h.test.Fatalf(h.Error)
		return nil
	}
	h.ResBody = response
	h.Status = res.StatusCode
	h.httpRes = res
	return h
}
