package service

import (
	"fmt"
	"sync"
	"time"
)

type Store struct {
	mu   sync.RWMutex
	jobs map[string]*Job
}

func NewStore() *Store {
	return &Store{jobs: make(map[string]*Job)}
}

func (s *Store) Create(id, filePath string) *Job {
	now := time.Now()
	j := &Job{
		ID:        id,
		Status:    StatusPending,
		FilePath:  filePath,
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.mu.Lock()
	s.jobs[id] = j
	s.mu.Unlock()
	return j
}

func (s *Store) Get(id string) (*Job, error) {
	s.mu.RLock()
	j, ok := s.jobs[id]
	s.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("job %s not found", id)
	}
	return j, nil
}

func (s *Store) SetDone(id string, result *Result) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if j, ok := s.jobs[id]; ok {
		j.Status = StatusDone
		j.Result = result
		j.UpdatedAt = time.Now()
	}
}

func (s *Store) SetError(id string, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if j, ok := s.jobs[id]; ok {
		j.Status = StatusError
		j.ErrMsg = err.Error()
		j.UpdatedAt = time.Now()
	}
}
