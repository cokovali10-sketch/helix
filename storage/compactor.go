// Package storage provides compaction
package storage

import (
	"fmt"
	"sync"
	"time"
)

// CompactionJob represents a compaction job
type CompactionJob struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"`
	StartedAt time.Time `json:"started_at"`
	FinishedAt time.Time `json:"finished_at"`
	InputFiles  int     `json:"input_files"`
	OutputFiles int     `json:"output_files"`
}

// Compactor manages compaction
type Compactor struct {
	jobs map[string]*CompactionJob
	mu   sync.Mutex
}

// NewCompactor creates compactor
func NewCompactor() *Compactor {
	return &Compactor{jobs: make(map[string]*CompactionJob)}
}

// StartJob starts a job
func (c *Compactor) StartJob(inputFiles int) *CompactionJob {
	c.mu.Lock()
	defer c.mu.Unlock()
	job := &CompactionJob{
		ID:         fmt.Sprintf("compact-%d", time.Now().UnixNano()),
		Status:     "running",
		StartedAt:  time.Now(),
		InputFiles: inputFiles,
	}
	c.jobs[job.ID] = job
	return job
}

// FinishJob finishes job
func (c *Compactor) FinishJob(id string, outputFiles int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if job, ok := c.jobs[id]; ok {
		job.Status = "completed"
		job.FinishedAt = time.Now()
		job.OutputFiles = outputFiles
	}
}

// FailJob marks failed
func (c *Compactor) FailJob(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if job, ok := c.jobs[id]; ok {
		job.Status = "failed"
		job.FinishedAt = time.Now()
	}
}

// GetJob returns job
func (c *Compactor) GetJob(id string) (*CompactionJob, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	job, ok := c.jobs[id]
	return job, ok
}

// ListJobs lists jobs
func (c *Compactor) ListJobs() []*CompactionJob {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := []*CompactionJob{}
	for _, j := range c.jobs {
		out = append(out, j)
	}
	return out
}

// Count returns count
func (c *Compactor) Count() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.jobs)
}
