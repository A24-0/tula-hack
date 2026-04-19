package service

import (
	"context"
	"log/slog"

	"github.com/tula-hack/voice-redaction/internal/processor"
)

type Processor interface {
	Process(ctx context.Context, jobID, filePath string) (*processor.Response, error)
}

type Queue struct {
	tasks chan *Job
	store *Store
	proc  Processor
}

func NewQueue(bufSize, workers int, store *Store, proc Processor) *Queue {
	q := &Queue{
		tasks: make(chan *Job, bufSize),
		store: store,
		proc:  proc,
	}
	for i := 0; i < workers; i++ {
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
			slog.Error("processor failed", "job", j.ID, "err", err)
			q.store.SetError(j.ID, err)
			continue
		}

		var events []PIIEvent
		for _, e := range resp.PIIEvents {
			events = append(events, PIIEvent{
				Type:     e.Type,
				Original: e.Original,
				StartSec: e.StartSec,
				EndSec:   e.EndSec,
			})
		}

		var words []Word
		for _, w := range resp.Words {
			words = append(words, Word{
				Word:  w.Word,
				Start: w.Start,
				End:   w.End,
			})
		}

		q.store.SetDone(j.ID, &Result{
			Transcript:         resp.Transcript,
			RedactedTranscript: resp.RedactedTranscript,
			RedactedAudioPath:  resp.RedactedAudioPath,
			PIIEvents:          events,
			Words:              words,
		})
		slog.Info("job done", "job", j.ID)
	}
}
