package myLogger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type LokiEntry struct {
	Streams []LokiStream `json:"streams"`
}

type LokiStream struct {
	Stream map[string]string `json:"stream"`
	Values [][2]string        `json:"values"`
}

type Loki struct {}

// Deprecated: Use LogV13 instead.
// Log sends a log message to Loki with the specified labels.
func (l *Loki) LogV12(message string, labels map[string]string) error {
	const (
		maxRetries = 3
		retryDelay = 1 * time.Second
		lokiURL    = "http://localhost:3100/loki/api/v1/push"
	)

	entry := LokiEntry{
		Streams: []LokiStream{
			{
				Stream: labels,
				Values: [][2]string{
					{fmt.Sprintf("%d", time.Now().UnixNano()), message},
				},
			},
		},
	}

	payload, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal Loki payload: %w", err)
	}

	for attempt := 1; attempt <= maxRetries; attempt++ {
		if err := sendToLoki(lokiURL, payload); err == nil {
			return nil
		} else if attempt < maxRetries {
			time.Sleep(retryDelay)
		} else {
			return fmt.Errorf("failed after %d attempts: %w", maxRetries, err)
		}
	}

	return nil // nunca chega aqui, mas mantém assinatura válida
}

// LogV13 sends a structure log message to Loki based on the Schema V13 format.
// It includes a timestamp and retries sending the log up to 3 times with a 1 second delay between attempts.
// It uses the Loki HTTP API to push logs.
// The log body is expected to be a map with string keys and any values.
// The baseLabels parameter is used to set the stream labels for the log entry.
// It returns an error if the log could not be sent after all retries.
// Example usage:
//   logger := myLogger.Loki{}
//   err := logger.LogV13(map[string]string{"app": "myapp", "level": "info"}, map[string]any{"message": "This is a log message"})
func (l *Loki) LogV13(baseLabels map[string]string, body map[string]any) error {
	const (
		lokiURL    = "http://localhost:3100/loki/api/v1/push"
		maxRetries = 3
		retryDelay = 1 * time.Second
	)

	// Timestamp
	timestamp := time.Now().UnixNano()
	body["timestamp"] = time.Now().Format(time.RFC3339Nano)

	// Encode log body as JSON
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal log body: %w", err)
	}

	// Construct log entry
	entry := LokiEntry{
		Streams: []LokiStream{
			{
				Stream: baseLabels,
				Values: [][2]string{
					{fmt.Sprintf("%d", timestamp), string(jsonBody)},
				},
			},
		},
	}

	payload, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal Loki payload: %w", err)
	}

	// Retry logic
	for attempt := 1; attempt <= maxRetries; attempt++ {
		if err := sendToLoki(lokiURL, payload); err == nil {
			return nil
		} else if attempt < maxRetries {
			time.Sleep(retryDelay)
		} else {
			return fmt.Errorf("failed after %d attempts: %w", maxRetries, err)
		}
	}

	return nil
}

func sendToLoki(url string, data []byte) error {
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("post request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("loki responded with %d: %s", resp.StatusCode, string(body))
	}

	return nil
}