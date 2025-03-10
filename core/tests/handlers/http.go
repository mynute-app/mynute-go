package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

type HttpClient struct{}

type httpActions struct {
	Error          string
	Status         int
	ResBody        map[string]any
	ResHeaders     map[string][]string
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
func (h *httpActions) Send(body any) *httpActions {
	h.Error = ""
	h.test.Logf("sending %s request to %s", h.method, h.url)
	h.test.Logf("body: %+v", body)
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		h.Error = fmt.Sprintf("failed to marshal payload: %v", err)
		h.test.Fatalf(h.Error)
		return nil
	}
	bodyBuffer := bytes.NewBuffer(bodyBytes)
	req, err := http.NewRequest(h.method, h.url, bodyBuffer)
	req.Header.Set("Content-Type", "application/json")
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
	if res.Body != nil {
		defer res.Body.Close()
	}
	if h.expectedStatus != 0 && res.StatusCode != h.expectedStatus {
		h.Error = fmt.Sprintf("expected status code: %d | received status code: %d", h.expectedStatus, res.StatusCode)
	}
	h.ResHeaders = res.Header
	if res.ContentLength != 0 && res.Body != nil {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			h.Error = fmt.Sprintf("failed to read response body: %v", err)
			h.test.Fatalf(h.Error)
			return nil
		}
		var response map[string]any
		if err := json.Unmarshal(bodyBytes, &response); err != nil {
			h.ResBody = map[string]any{
				"string": string(bodyBytes),
			}
		} else {
			h.ResBody = response
		}
	}
	h.test.Logf("response: %+v", h.ResBody)
	if h.Error != "" {
		h.test.Fatalf(h.Error)
	}
	h.Status = res.StatusCode
	h.httpRes = res
	return h
}
