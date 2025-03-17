package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
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
	headers        http.Header // Store headers before sending the request
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

func (h *httpActions) Header(key, value string) *httpActions {
	if h.headers == nil {
		h.headers = http.Header{}
	}
	h.headers.Set(key, value)
	return h
}

func (h *httpActions) URL(url string) *httpActions {
	AppPort := os.Getenv("APP_PORT")
	BaseUrl := fmt.Sprintf("http://localhost:%s", AppPort)
	FullUrl := fmt.Sprintf(BaseUrl + url)
	h.url = FullUrl
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

func (h *httpActions) ParseResponse(to interface{}) {
	if reflect.TypeOf(to).Kind() != reflect.Ptr {
		h.test.Fatalf("expected a pointer to a struct")
	}
	if h.ResBody == nil {
		h.test.Fatalf("response body is nil")
	}

	// Marshal the map into JSON bytes
	resBodyBytes, err := json.Marshal(h.ResBody)
	if err != nil {
		h.test.Fatalf("failed to marshal response body: %v", err)
	}

	// Unmarshal JSON bytes directly into your struct
	if err := json.Unmarshal(resBodyBytes, to); err != nil {
		h.test.Fatalf("failed to unmarshal response body: %v", err)
	}
}

// Send sends a request to the specified URL with the specified body
// and returns the response. If an error occurs, it sets the error message.
// The response body is stored in ResBody, and the status code is stored in Status.
// It will defer res.Body.Close() to ensure the response body is closed.
func (h *httpActions) Send(body any) *httpActions {
	h.Error = ""
	h.test.Logf(">>>>>>>>>> Sending %s request to %s", h.method, h.url)
	h.test.Logf("request body: %+v", body)
	h.test.Logf("request headers: %+v", h.headers)
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		h.test.Fatalf("failed to marshal payload: %v", err)
		return nil
	}
	bodyBuffer := bytes.NewBuffer(bodyBytes)
	req, err := http.NewRequest(h.method, h.url, bodyBuffer)
	if err != nil {
		h.test.Fatalf("failed to create request: %v", err)
		return nil
	}
	req.Header.Set("Content-Type", "application/json")
	for key, values := range h.headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		h.test.Fatalf("failed to send request: %v", err)
		return nil
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	h.ResHeaders = res.Header
	if res.ContentLength != 0 && res.Body != nil {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			h.test.Fatalf("failed to read response body: %v", err)
		}

		if !json.Valid(bodyBytes) {
			h.ResBody = map[string]any{"data": string(bodyBytes)} // or handle accordingly
		} else if len(bodyBytes) > 0 && bodyBytes[0] == '"' && bodyBytes[len(bodyBytes)-1] == '"' {
			h.ResBody = map[string]any{"data": string(bodyBytes)}
		} else {
			decoder := json.NewDecoder(bytes.NewReader(bodyBytes))
			decoder.UseNumber()
			if err := decoder.Decode(&h.ResBody); err != nil {
				h.test.Fatalf("failed to unmarshal response body: %v", err)
			}
		}
	}
	h.test.Logf("response body: %+v", h.ResBody)
	h.test.Logf("response headers: %+v", h.ResHeaders)

	if h.expectedStatus != 0 && res.StatusCode != h.expectedStatus {
		h.test.Fatalf("expected status code: %d | received status code: %d", h.expectedStatus, res.StatusCode)
	}

	h.Status = res.StatusCode
	h.httpRes = res
	return h
}
