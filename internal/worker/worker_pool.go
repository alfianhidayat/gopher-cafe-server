package worker

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/ajaibid/coin-common-golang/logger"
)

type WorkerPool struct {
	name       string
	jobs       chan JobInput
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
	numWorkers uint8
}

func NewWorkerPool(name string, workers uint8) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	wp := &WorkerPool{
		name:       name,
		jobs:       make(chan JobInput),
		ctx:        ctx,
		cancel:     cancel,
		numWorkers: workers,
	}

	return wp
}

func (wp *WorkerPool) start() {
	for i := range wp.numWorkers {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

func (wp *WorkerPool) stop() {
	wp.cancel()    // stop signal
	close(wp.jobs) // stop accepting jobs
	wp.wg.Wait()   // wait for workers to finish
}

func (wp *WorkerPool) worker(id uint8) {
	defer wp.wg.Done()

	for {
		select {
		case <-wp.ctx.Done():
			logger.Debugf("[%s] worker %d stopped", wp.name, id)
			return

		case job, ok := <-wp.jobs:
			if !ok {
				logger.Debugf("[%s] worker %d stopped, got channel closed", wp.name, id)
				return
			}
			logger.Debugf("[%s] worker %d doing job: %v start", wp.name, id, job.Job)
			time.Sleep(job.Job.Timer)
			job.Output <- JobOutput{
				job.Job,
				nil,
			}
			logger.Debugf("[%s] worker %d doing job: %v finish", wp.name, id, job.Job)
		}
	}
}

func (wp *WorkerPool) Submit(job Job) error {
	ji := JobInput{
		Job:    job,
		Output: make(chan JobOutput),
	}

	select {
	case wp.jobs <- ji:
	case <-wp.ctx.Done():
		return errors.New("pool closed")
	}

	// wait for response
	select {
	case res := <-ji.Output:
		return res.Err
	case <-wp.ctx.Done():
		return errors.New("shutdown")
	}
}
