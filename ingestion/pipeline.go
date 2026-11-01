// Package ingestion provides pipeline
package ingestion

import (
	"sync"
)

// PipelineStage represents a stage
type PipelineStage func(map[string]interface{}) (map[string]interface{}, error)

// Pipeline processes records through stages
type Pipeline struct {
	stages []PipelineStage
	mu     sync.Mutex
}

// NewPipeline creates pipeline
func NewPipeline() *Pipeline {
	return &Pipeline{stages: []PipelineStage{}}
}

// AddStage adds stage
func (p *Pipeline) AddStage(stage PipelineStage) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.stages = append(p.stages, stage)
}

// Process processes record
func (p *Pipeline) Process(record map[string]interface{}) (map[string]interface{}, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	current := record
	for _, stage := range p.stages {
		next, err := stage(current)
		if err != nil {
			return nil, err
		}
		current = next
	}
	return current, nil
}

// StageCount returns count
func (p *Pipeline) StageCount() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.stages)
}

// Clear clears stages
func (p *Pipeline) Clear() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.stages = []PipelineStage{}
}
