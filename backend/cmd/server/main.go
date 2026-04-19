package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/tula-hack/voice-redaction/internal/api"
	"github.com/tula-hack/voice-redaction/internal/processor"
	"github.com/tula-hack/voice-redaction/internal/service"
)

func main() {
	cfg := loadConfig()

	if err := os.MkdirAll(cfg.uploadDir, 0o755); err != nil {
		slog.Error("cannot create upload dir", "err", err)
		os.Exit(1)
	}

	proc := processor.NewClient(cfg.processorURL)
	store := service.NewStore()
	queue := service.NewQueue(cfg.queueBuffer, cfg.workers, store, proc)
	h := api.NewHandler(store, queue, cfg.uploadDir)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/upload", h.Upload)
	r.Get("/health", h.Health)
	r.Get("/jobs/{id}", h.JobStatus)
	r.Get("/transcript/{id}", h.Transcript)
	r.Get("/audio/{id}/redacted", h.RedactedAudio)
	r.Get("/audio/{id}/original", h.OriginalAudio)
	r.Get("/logs/{id}", h.Logs)

	slog.Info("starting server", "addr", cfg.addr)
	if err := http.ListenAndServe(cfg.addr, r); err != nil {
		slog.Error("server stopped", "err", err)
		os.Exit(1)
	}
}

type config struct {
	addr         string
	uploadDir    string
	processorURL string
	workers      int
	queueBuffer  int
}

func loadConfig() config {
	return config{
		addr:         envOr("ADDR", ":8080"),
		uploadDir:    envOr("UPLOAD_DIR", "./uploads"),
		processorURL: envOr("PROCESSOR_URL", "http://localhost:8000"),
		workers:      4,
		queueBuffer:  64,
	}
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
