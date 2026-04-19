package service

import "time"

type Status string

const (
	StatusPending Status = "pending"
	StatusDone    Status = "done"
	StatusError   Status = "error"
)

type PIIEvent struct {
	Type     string  `json:"type"`
	Original string  `json:"original"`
	StartSec float64 `json:"start_sec"`
	EndSec   float64 `json:"end_sec"`
}

type Word struct {
	Word  string  `json:"word"`
	Start float64 `json:"start"`
	End   float64 `json:"end"`
}

type Result struct {
	Transcript         string     `json:"transcript"`
	RedactedTranscript string     `json:"redacted_transcript"`
	RedactedAudioPath  string     `json:"redacted_audio_path,omitempty"`
	PIIEvents          []PIIEvent `json:"pii_events"`
	Words              []Word     `json:"words,omitempty"`
}

type Job struct {
	ID        string    `json:"id"`
	Status    Status    `json:"status"`
	FilePath  string    `json:"-"`
	Result    *Result   `json:"result,omitempty"`
	ErrMsg    string    `json:"error,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
