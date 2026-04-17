package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/tula-hack/voice-redaction/internal/service"
)

const maxUploadSize = 500 << 20 // 500 MB

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

	ext, ok := allowedAudioTypes[header.Header.Get("Content-Type")]
	if !ok {
		buf := make([]byte, 512)
		n, _ := file.Read(buf)
		ext, ok = allowedAudioTypes[http.DetectContentType(buf[:n])]
		if !ok {
			writeError(w, http.StatusUnsupportedMediaType, "unsupported audio format")
			return
		}
		if _, err := file.Seek(0, io.SeekStart); err != nil {
			writeError(w, http.StatusInternalServerError, "seek error")
			return
		}
	}

	jobID := uuid.NewString()
	destPath := filepath.Join(h.uploadDir, fmt.Sprintf("%s%s", jobID, ext))

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

	j := h.store.Create(jobID, destPath)
	h.queue.Enqueue(j)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"job_id": jobID})
}

func (h *Handler) JobStatus(w http.ResponseWriter, r *http.Request) {
	j, err := h.store.Get(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusNotFound, "job not found")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(j)
}


func writeError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
