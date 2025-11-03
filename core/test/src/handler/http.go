package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
	"strings"
)

type Files map[string]MyFile

type MyFile struct {
	Name    string
	Content []byte
}

type httpActions struct {
	Error           error
	Status          int
	ResBody         map[string]any
	ResHeaders      map[string][]string
	url             string
	method          string
	expectedStatus  int
	httpRes         *http.Response
	reqHeaders      http.Header // Store headers before sending the request
	rawResponseBody []byte      // Store the raw response body
	log_level       string      // 'silent' or 'verbose'
}

func NewHttpClient() *httpActions {
	log_level := os.Getenv("HTTP_TEST_LOG_LEVEL")
	if log_level == "" {
		log_level = "silent" // Default log level
	}
	return &httpActions{
		ResBody:    make(map[string]any),
		ResHeaders: make(map[string][]string),
		log_level:  log_level,
	}
}

func (h *httpActions) URL(url string) *httpActions {
	AppPort := os.Getenv("APP_PORT")
	BaseUrl := fmt.Sprintf("http://localhost:%s", AppPort)
	if strings.HasPrefix(url, "http://localhost:") || strings.HasPrefix(url, "https://localhost:") {
		BaseUrl = ""
	}

	// Auto-prepend /api for API endpoints (except public pages)
	if !strings.HasPrefix(url, "/api") &&
		!strings.HasPrefix(url, "http") &&
		url != "/" &&
		url != "/verify-email" &&
		!strings.HasPrefix(url, "/translations") {
		url = "/api" + url
	}

	FullUrl := fmt.Sprintf(BaseUrl + url)
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

// Sends a request to the specified URL with the specified body
// and returns the response. If an error occurs, it sets the error message.
// The response body is stored in ResBody, and the status code is stored in Status.
// It will defer res.Body.Close() to ensure the response body is closed.
func (h *httpActions) Send(body any) *httpActions {
	if h.Error != nil {
		return h
	}

	var bodyBytes []byte
	var err error

	switch b := body.(type) {
	case Files:
		bodyBytes, err = h.parseFiles(b) // Isso já cuida do header
		if err != nil {
			return h.set_error(fmt.Sprintf("failed to parse files: %v", err))
		}
	case []byte:
		bodyBytes = b
	default:
		var err error
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
			decoder := json.NewDecoder(bytes.NewReader(bodyBytes))
			decoder.UseNumber()
			if err := decoder.Decode(&h.ResBody); err != nil {
				return h.set_error(fmt.Sprintf("failed to unmarshal response body: %v", err))
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

	// Se for []byte, retorna o corpo bruto da resposta
	if b, ok := to.(*[]byte); ok {
		*b = h.rawResponseBody
		return h
	}

	if h.ResBody == nil {
		return h
	}

	resBodyBytes, err := json.Marshal(h.ResBody)
	if err != nil {
		return h.set_error(fmt.Sprintf("failed to marshal response body: %v", err))
	}

	if err := json.Unmarshal(resBodyBytes, to); err != nil {
		return h.set_error(fmt.Sprintf("failed to unmarshal response body: %v", err))
	}

	return h
}

func (h *httpActions) parseFiles(files Files) ([]byte, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("no files provided for multipart request")
	}
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	for field, f := range files {
		part, err := writer.CreateFormFile(field, f.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to create form file for field %s: %v", field, err)
		}
		if _, err := part.Write(f.Content); err != nil {
			return nil, fmt.Errorf("failed to write content for field %s: %v", field, err)
		}
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %v", err)
	}

	h.Header("Content-Type", writer.FormDataContentType())
	return b.Bytes(), nil
}

func (h *httpActions) set_error(e string) *httpActions {
	if h.Error == nil {
		h.Error = errors.New(e)
	}
	return h
}
