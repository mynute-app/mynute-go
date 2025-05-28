package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"
)

type HttpClient struct{}

type httpActions struct {
	Error           string
	Status          int
	ResBody         map[string]any
	ResHeaders      map[string][]string
	url             string
	method          string
	expectedStatus  int
	test            *testing.T
	httpRes         *http.Response
	headers         http.Header // Store headers before sending the request
	rawResponseBody []byte      // Store the raw response body
}

func (c *HttpClient) SetTest(t *testing.T) *httpActions {
	h := &httpActions{}
	h.test = t
	h.expectedStatus = 0
	return h
}

type MyFile struct {
	Name    string
	Content []byte
}

type Files map[string]MyFile

func (h *httpActions) parseFiles(files Files) []byte {
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	for field, f := range files {
		part, err := writer.CreateFormFile(field, f.Name)
		if err != nil {
			h.test.Fatalf("failed to create form file for field %s: %v", field, err)
		}
		if _, err := part.Write(f.Content); err != nil {
			h.test.Fatalf("failed to write content for field %s: %v", field, err)
		}
	}

	if err := writer.Close(); err != nil {
		h.test.Fatalf("failed to close multipart writer: %v", err)
	}

	h.Header("Content-Type", writer.FormDataContentType())
	return b.Bytes()
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
	if strings.HasPrefix(url, "http://localhost:") || strings.HasPrefix(url, "https://localhost:") {
		BaseUrl = ""
	}
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

func (h *httpActions) ParseResponse(to any) {
	if reflect.TypeOf(to).Kind() != reflect.Ptr {
		h.test.Fatalf("expected a pointer to a struct")
	}

	// Se for []byte, retorna o corpo bruto da resposta
	if b, ok := to.(*[]byte); ok {
		*b = h.rawResponseBody
		return
	}

	if h.ResBody == nil {
		h.test.Fatalf("response body is nil")
	}

	resBodyBytes, err := json.Marshal(h.ResBody)
	if err != nil {
		h.test.Fatalf("failed to marshal response body: %v", err)
	}

	if err := json.Unmarshal(resBodyBytes, to); err != nil {
		h.test.Fatalf("failed to unmarshal response body: %v", err)
	}
}

// Send sends a request to the specified URL with the specified body
// and returns the response. If an error occurs, it sets the error message.
// The response body is stored in ResBody, and the status code is stored in Status.
// It will defer res.Body.Close() to ensure the response body is closed.
func (h *httpActions) Send(body any) *httpActions {
	var bodyBytes []byte

	switch b := body.(type) {
	case Files:
		bodyBytes = h.parseFiles(b) // Isso já cuida do header
	case []byte:
		bodyBytes = b
	default:
		var err error
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			h.test.Fatalf("failed to marshal payload: %v", err)
			return nil
		}
	}

	h.test.Logf(">>>>>>>>>> Sending %s request to %s", h.method, h.url)
	h.test.Logf("request body: %+v", body)
	h.test.Logf("request headers: %+v", h.headers)

	bodyBuffer := bytes.NewBuffer(bodyBytes)

	req, err := http.NewRequest(h.method, h.url, bodyBuffer)
	if err != nil {
		h.test.Fatalf("failed to create request: %v", err)
		return nil
	}

	if h.headers.Get("Content-Type") == "" {
		h.Header("Content-Type", "application/json")
	}

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
	defer res.Body.Close()

	h.ResHeaders = res.Header

	if res.ContentLength != 0 && res.Body != nil {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			h.test.Fatalf("failed to read response body: %v", err)
		}
		h.rawResponseBody = bodyBytes // <--- Aqui salva

		if !json.Valid(bodyBytes) {
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
		h.test.Fatalf("expected status code: %d | received status code: %d \n", h.expectedStatus, res.StatusCode)
	}

	h.Status = res.StatusCode
	h.httpRes = res
	return h
}
