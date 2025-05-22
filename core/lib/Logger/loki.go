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

func (l *Loki) Log(message string, labels map[string]string) error {
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

func sendToLoki(url string, data []byte) error {
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("post request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 || resp.StatusCode < 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("loki responded with %d: %s", resp.StatusCode, string(body))
	}

	return nil
}