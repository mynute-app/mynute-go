package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
)

type httpActions struct {
	Error           error
	Status          int
	ResBody         map[string]any
	ResHeaders      map[string][]string
	url             string
	method          string
	expectedStatus  int
	httpRes         *http.Response
	reqHeaders      http.Header
	rawResponseBody []byte
	log_level       string
}

func NewHttpClient() *httpActions {
	log_level := os.Getenv("HTTP_TEST_LOG_LEVEL")
	if log_level == "" {
		log_level = "silent"
	}
	return &httpActions{
		ResBody:    make(map[string]any),
		ResHeaders: make(map[string][]string),
		log_level:  log_level,
	}
}

func (h *httpActions) URL(url string) *httpActions {
	AppPort := os.Getenv("AUTH_SERVICE_PORT")
	if AppPort == "" {
		AppPort = "4001" // Default auth service port
	}
	BaseUrl := fmt.Sprintf("http://localhost:%s", AppPort)
	if strings.HasPrefix(url, "http://localhost:") || strings.HasPrefix(url, "https://localhost:") {
		BaseUrl = ""
	}

	// Auth service doesn't use /api prefix - routes are mounted directly
	// No need to auto-prepend /api

	FullUrl := BaseUrl + url
	h.url = FullUrl
	return h
}

func (h *httpActions) Method(method string) *httpActions {
	h.method = method
	return h
}

func (h *httpActions) ExpectedStatus(status int) *httpActions {
	h.expectedStatus = status
	return h
}

func (h *httpActions) Header(key string, value string) *httpActions {
	if h.reqHeaders == nil {
		h.reqHeaders = make(http.Header)
	}
	h.reqHeaders.Add(key, value)
	return h
}

func (h *httpActions) Send(body any) *httpActions {
	if h.Error != nil {
		return h
	}

	var bodyBytes []byte
	var err error

	switch b := body.(type) {
	case []byte:
		bodyBytes = b
	default:
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return h.set_error(fmt.Sprintf("failed to marshal body payload: %v", err))
		}
	}

	h.ResBody = nil
	h.ResHeaders = nil

	if h.log_level == "verbose" {
		fmt.Printf("\n--------------- NEW REQUEST ---------------\n")
		fmt.Printf("\n>>>>>>>>>> Sending %s request to %s <<<<<<<<<<\n", h.method, h.url)
		fmt.Printf("\n>>>>>>>>>> Request body: %+v\n", body)
		fmt.Printf("\n>>>>>>>>>> Request headers: %+v\n\n", h.reqHeaders)
	}

	bodyBuffer := bytes.NewBuffer(bodyBytes)

	req, err := http.NewRequest(h.method, h.url, bodyBuffer)
	if err != nil {
		return h.set_error(fmt.Sprintf("failed to create request: %v", err))
	}

	if h.reqHeaders.Get("Content-Type") == "" {
		h.Header("Content-Type", "application/json")
	}

	for key, values := range h.reqHeaders {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return h.set_error(fmt.Sprintf("failed to send request: %v", err))
	}
	defer res.Body.Close()

	h.ResHeaders = res.Header

	if res.ContentLength != 0 && res.Body != nil {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return h.set_error(fmt.Sprintf("failed to read response body: %v", err))
		}
		h.rawResponseBody = bodyBytes

		if !json.Valid(bodyBytes) {
			h.ResBody = map[string]any{"data": string(bodyBytes)}
		} else {
			// Try to decode as map first
			decoder := json.NewDecoder(bytes.NewReader(bodyBytes))
			decoder.UseNumber()
			if err := decoder.Decode(&h.ResBody); err != nil {
				// If it fails (e.g., array response), we'll handle it in ParseResponse
				// Just store the raw body and continue
				h.ResBody = nil
			}
		}
	}

	if h.log_level == "verbose" {
		fmt.Printf("\n>>>>>>>>>> Response body: %+v\n", h.ResBody)
		fmt.Printf("\n>>>>>>>>>> Response headers: %+v\n", h.ResHeaders)
		fmt.Printf("\n---------------- x ----------------\n")
	}

	if h.expectedStatus != 0 && res.StatusCode != h.expectedStatus {
		return h.set_error(fmt.Sprintf("expected status code: %d | received status code: %d \n response: %+v", h.expectedStatus, res.StatusCode, h.ResBody))
	}

	h.Status = res.StatusCode
	h.httpRes = res
	return h
}

func (h *httpActions) ParseResponse(to any) *httpActions {
	if h.Error != nil {
		return h
	}

	if h.Status == 0 {
		return h.set_error("no response received, please call Send() first")
	}

	if reflect.TypeOf(to).Kind() != reflect.Ptr {
		return h.set_error(fmt.Sprintf("the parse destination must be a pointer, but instead got %T", to))
	}

	// If it's []byte, return raw response body
	if b, ok := to.(*[]byte); ok {
		*b = h.rawResponseBody
		return h
	}

	// If it's map[string]interface{}, return ResBody directly
	if m, ok := to.(*map[string]interface{}); ok {
		*m = h.ResBody
		return h
	}

	// Otherwise, marshal ResBody and unmarshal to destination
	jsonBytes, err := json.Marshal(h.ResBody)
	if err != nil {
		return h.set_error(fmt.Sprintf("failed to marshal response body: %v", err))
	}

	if err := json.Unmarshal(jsonBytes, to); err != nil {
		return h.set_error(fmt.Sprintf("failed to unmarshal response body: %v", err))
	}

	return h
}

func (h *httpActions) set_error(err string) *httpActions {
	h.Error = fmt.Errorf("%s", err)
	return h
}
