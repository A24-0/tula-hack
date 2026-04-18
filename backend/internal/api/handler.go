package api

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/tula-hack/voice-redaction/internal/service"
)

const maxUploadSize = 500 << 20

var allowedAudioTypes = map[string]string{
	"audio/mpeg":  ".mp3",
	"audio/wav":   ".wav",
	"audio/x-wav": ".wav",
	"audio/ogg":   ".ogg",
	"audio/mp4":   ".m4a",
	"audio/flac":  ".flac",
}

type Handler struct {
	store     *service.Store
	queue     *service.Queue
	uploadDir string
}

func NewHandler(store *service.Store, queue *service.Queue, uploadDir string) *Handler {
	return &Handler{store: store, queue: queue, uploadDir: uploadDir}
}

func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "failed to parse form: "+err.Error())
		return
	}

	file, header, err := r.FormFile("audio")
	if err != nil {
		writeError(w, http.StatusBadRequest, "field 'audio' is required")
		return
	}
	defer file.Close()

	ct := header.Header.Get("Content-Type")
	ext, ok := allowedAudioTypes[ct]
	if !ok {
		buf := make([]byte, 512)
		n, _ := file.Read(buf)
		sniffed := http.DetectContentType(buf[:n])
		ext, ok = allowedAudioTypes[sniffed]
		if !ok {
			writeError(w, http.StatusUnsupportedMediaType, "unsupported audio format")
			return
		}
		file.Seek(0, io.SeekStart)
	}

	jobID := uuid.NewString()
	destPath := filepath.Join(h.uploadDir, jobID+ext)

	dst, err := os.Create(destPath)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not save file")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		os.Remove(destPath)
		writeError(w, http.StatusInternalServerError, "could not write file")
		return
	}

	job := h.store.Create(jobID, destPath)
	h.queue.Enqueue(job)

	w.WriteHeader(http.StatusAccepted)
	writeJSON(w, map[string]string{"job_id": jobID})
}

func (h *Handler) JobStatus(w http.ResponseWriter, r *http.Request) {
	job, err := h.store.Get(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusNotFound, "job not found")
		return
	}
	writeJSON(w, job)
}

func (h *Handler) Transcript(w http.ResponseWriter, r *http.Request) {
	job, err := h.store.Get(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusNotFound, "job not found")
		return
	}
	if job.Status != service.StatusDone || job.Result == nil {
		writeError(w, http.StatusConflict, "not ready yet")
		return
	}
	writeJSON(w, job.Result)
}

func (h *Handler) RedactedAudio(w http.ResponseWriter, r *http.Request) {
	job, err := h.store.Get(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusNotFound, "job not found")
		return
	}
	if job.Result == nil || job.Result.RedactedAudioPath == "" {
		writeError(w, http.StatusNotFound, "redacted audio not available")
		return
	}
	http.ServeFile(w, r, job.Result.RedactedAudioPath)
}

func (h *Handler) Logs(w http.ResponseWriter, r *http.Request) {
	job, err := h.store.Get(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusNotFound, "job not found")
		return
	}
	events := make([]service.PIIEvent, 0)
	if job.Result != nil {
		events = job.Result.PIIEvents
	}
	writeJSON(w, events)
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
