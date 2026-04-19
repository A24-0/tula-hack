package api

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/tula-hack/voice-redaction/internal/service"
)

// newTestQueue creates a Queue with 0 workers so Enqueue blocks on the channel
// but never calls the processor.
func newTestQueue(store *service.Store) *service.Queue {
	return service.NewQueue(16, 0, store, nil)
}

func TestHealth(t *testing.T) {
	h := &Handler{}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	h.Health(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", rec.Code)
	}
	if got := rec.Body.String(); got != `{"status":"ok"}` {
		t.Fatalf("unexpected body: %q", got)
	}
}

// minimalWAV returns a minimal 44-byte WAV header that DetectContentType
// will recognise as "audio/wav".
func minimalWAV() []byte {
	// RIFF....WAVEfmt  + 16 bytes of PCM format chunk
	b := make([]byte, 44)
	copy(b[0:], []byte("RIFF"))
	// chunk size (little-endian uint32) — value doesn't matter for sniffing
	b[4] = 0x24
	copy(b[8:], []byte("WAVE"))
	copy(b[12:], []byte("fmt "))
	return b
}

func buildMultipartRequest(t *testing.T, fieldName, filename, ct string, body []byte) *http.Request {
	t.Helper()
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)

	h := textproto.MIMEHeader{}
	h.Set("Content-Disposition", `form-data; name="`+fieldName+`"; filename="`+filename+`"`)
	h.Set("Content-Type", ct)

	part, err := mw.CreatePart(h)
	if err != nil {
		t.Fatalf("create part: %v", err)
	}
	part.Write(body)
	mw.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload", &buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	return req
}

func TestUpload_HappyPath(t *testing.T) {
	tmpDir := t.TempDir()
	store := service.NewStore()
	queue := newTestQueue(store)
	h := NewHandler(store, queue, tmpDir)

	req := buildMultipartRequest(t, "audio", "x.wav", "audio/wav", minimalWAV())
	rec := httptest.NewRecorder()
	h.Upload(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("want 202, got %d; body: %s", rec.Code, rec.Body.String())
	}

	var resp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	jobID, ok := resp["job_id"]
	if !ok || jobID == "" {
		t.Fatalf("missing job_id in response: %v", resp)
	}

	// Job must exist in store with pending status.
	job, err := store.Get(jobID)
	if err != nil {
		t.Fatalf("store.Get(%q): %v", jobID, err)
	}
	if job.Status != service.StatusPending {
		t.Fatalf("want status pending, got %q", job.Status)
	}

	// File must exist on disk.
	expected := filepath.Join(tmpDir, jobID+".wav")
	if _, err := os.Stat(expected); os.IsNotExist(err) {
		t.Fatalf("expected file %q to exist", expected)
	}
}

func TestUpload_BadContentType(t *testing.T) {
	tmpDir := t.TempDir()
	store := service.NewStore()
	queue := newTestQueue(store)
	h := NewHandler(store, queue, tmpDir)

	// Use Content-Type that is not audio and bytes that don't look like any audio format.
	body := []byte("this is definitely not audio data at all!!!!!")
	req := buildMultipartRequest(t, "audio", "x.bin", "application/octet-stream", body)
	rec := httptest.NewRecorder()
	h.Upload(rec, req)

	if rec.Code < 400 {
		t.Fatalf("want 4xx, got %d", rec.Code)
	}
}

// withChiID injects a chi URL param "id" into the request context.
func withChiID(r *http.Request, id string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", id)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}

func TestOriginalAudio(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.wav")
	content := []byte("fake-audio-content")
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	store := service.NewStore()
	store.Create("job-orig", filePath)

	queue := newTestQueue(store)
	h := NewHandler(store, queue, tmpDir)

	req := httptest.NewRequest(http.MethodGet, "/audio/job-orig/original", nil)
	req = withChiID(req, "job-orig")
	rec := httptest.NewRecorder()
	h.OriginalAudio(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("want 200, got %d; body: %s", rec.Code, rec.Body.String())
	}
	if got := rec.Body.Bytes(); !bytes.Equal(got, content) {
		t.Fatalf("body mismatch: got %q, want %q", got, content)
	}
}

func TestRedactedAudio(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test_redacted.wav")
	content := []byte("redacted-audio-content")
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	store := service.NewStore()
	store.Create("job-red", filepath.Join(tmpDir, "test.wav"))
	store.SetDone("job-red", &service.Result{
		RedactedAudioPath: filePath,
	})

	queue := newTestQueue(store)
	h := NewHandler(store, queue, tmpDir)

	req := httptest.NewRequest(http.MethodGet, "/audio/job-red/redacted", nil)
	req = withChiID(req, "job-red")
	rec := httptest.NewRecorder()
	h.RedactedAudio(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("want 200, got %d; body: %s", rec.Code, rec.Body.String())
	}
	if got := rec.Body.Bytes(); !bytes.Equal(got, content) {
		t.Fatalf("body mismatch: got %q, want %q", got, content)
	}
}

func TestJobStatus_NotFound(t *testing.T) {
	store := service.NewStore()
	queue := newTestQueue(store)
	h := NewHandler(store, queue, t.TempDir())

	req := httptest.NewRequest(http.MethodGet, "/jobs/nonexistent", nil)
	req = withChiID(req, "nonexistent")
	rec := httptest.NewRecorder()
	h.JobStatus(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("want 404, got %d", rec.Code)
	}
	body := strings.TrimSpace(rec.Body.String())
	if body != `{"error":"job not found"}` {
		t.Fatalf("unexpected body: %q", body)
	}
}

func TestLogs_NoResult(t *testing.T) {
	store := service.NewStore()
	store.Create("job-nores", "/some/path.wav")
	queue := newTestQueue(store)
	h := NewHandler(store, queue, t.TempDir())

	req := httptest.NewRequest(http.MethodGet, "/logs/job-nores", nil)
	req = withChiID(req, "job-nores")
	rec := httptest.NewRecorder()
	h.Logs(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", rec.Code)
	}
	body := strings.TrimSpace(rec.Body.String())
	if body != `[]` {
		t.Fatalf("unexpected body: %q", body)
	}
}
