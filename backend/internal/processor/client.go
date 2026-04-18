package processor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Minute,
		},
	}
}

type processRequest struct {
	JobID    string `json:"job_id"`
	FilePath string `json:"file_path"`
}

type PIIEvent struct {
	Type     string  `json:"type"`
	Original string  `json:"original"`
	StartSec float64 `json:"start_sec"`
	EndSec   float64 `json:"end_sec"`
}

type Response struct {
	Transcript         string     `json:"transcript"`
	RedactedTranscript string     `json:"redacted_transcript"`
	RedactedAudioPath  string     `json:"redacted_audio_path"`
	PIIEvents          []PIIEvent `json:"pii_events"`
}

func (c *Client) Process(ctx context.Context, jobID, filePath string) (*Response, error) {
	body, err := json.Marshal(processRequest{JobID: jobID, FilePath: filePath})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/process", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("processor unavailable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("processor returned %d", resp.StatusCode)
	}

	var result Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &result, nil
}
