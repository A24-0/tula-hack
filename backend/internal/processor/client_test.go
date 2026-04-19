package processor

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProcess_OK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/process" || r.Method != http.MethodPost {
			t.Fatalf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
            "transcript":"привет 9001234567",
            "redacted_transcript":"привет [PHONE]",
            "redacted_audio_path":"/uploads/x_redacted.wav",
            "pii_events":[{"type":"phone","original":"9001234567","start_sec":0.5,"end_sec":1.2}],
            "words":[{"word":"привет","start":0.0,"end":0.3}]
        }`))
	}))
	defer srv.Close()

	c := NewClient(srv.URL)
	res, err := c.Process(context.Background(), "job-1", "/uploads/x.wav")
	if err != nil {
		t.Fatalf("Process: %v", err)
	}
	if len(res.Words) != 1 || res.Words[0].Word != "привет" {
		t.Fatalf("words not parsed: %#v", res.Words)
	}
	if len(res.PIIEvents) != 1 || res.PIIEvents[0].Type != "phone" {
		t.Fatalf("pii not parsed: %#v", res.PIIEvents)
	}
}

func TestProcess_5xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	defer srv.Close()
	if _, err := NewClient(srv.URL).Process(context.Background(), "j", "/x"); err == nil {
		t.Fatal("expected error on 5xx, got nil")
	}
}
