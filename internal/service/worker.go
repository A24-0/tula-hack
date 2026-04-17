package service

import (
	"context"
	"log/slog"

	"github.com/tula-hack/voice-redaction/internal/processor"
)

type Queue struct {
	tasks chan *Job
	store *Store
	proc  *processor.Client
}

func NewQueue(bufSize, workers int, store *Store, proc *processor.Client) *Queue {
	q := &Queue{
		tasks: make(chan *Job, bufSize),
		store: store,
		proc:  proc,
	}
	for range workers {
		go q.run()
	}
	return q
}

func (q *Queue) Enqueue(j *Job) {
	q.tasks <- j
}

func (q *Queue) run() {
	for j := range q.tasks {
		resp, err := q.proc.Process(context.Background(), j.ID, j.FilePath)
		if err != nil {
			slog.Error("processor failed", "job_id", j.ID, "err", err)
			q.store.SetError(j.ID, err)
			continue
		}

		events := make([]PIIEvent, len(resp.PIIEvents))
		for i, e := range resp.PIIEvents {
			events[i] = PIIEvent{Type: e.Type, Original: e.Original, StartSec: e.StartSec, EndSec: e.EndSec}
		}
		q.store.SetDone(j.ID, &Result{
			Transcript:         resp.Transcript,
			RedactedTranscript: resp.RedactedTranscript,
			RedactedAudioPath:  resp.RedactedAudioPath,
			PIIEvents:          events,
		})
		slog.Info("job done", "job_id", j.ID, "pii_count", len(events))
	}
}
