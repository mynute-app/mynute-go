package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func SendLogToLoki(message string, labels map[string]string) error {
	now := time.Now().UnixNano()
	entry := LokiEntry{
		Streams: []LokiStream{
			{
				Stream: labels,
				Values: [][2]string{
					{fmt.Sprintf("%d", now), message},
				},
			},
		},
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	resp, err := http.Post("http://localhost:3100/loki/api/v1/push", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("loki returned status %d", resp.StatusCode)
	}
	return nil
}