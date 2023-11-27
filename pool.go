package pool

import (
	"context"
	"errors"
	"sync"
	"time"
)

// Config struct defines the configuration parameters for the Pool.
// - MaxWorkers: The maximum number of concurrent workers in the pool.
// - Timeout: The maximum duration for processing a single task. If a task takes longer, it will be terminated.
type Config struct {
	MaxWorkers int
	Timeout    time.Duration
}

// ErrorHandler is a function type for handling errors. It provides a way for users to define custom error handling logic.
// The function receives an error and a pointer to the pool, allowing users to log errors, modify the pool, or perform other actions.
type ErrorHandler[T any] func(error, *Pool[T])

// Pool struct represents a pool of workers processing tasks of type T.
// It encapsulates the task handling logic, error handling, and synchronization mechanisms.
type Pool[T any] struct {
	config       Config
	handler      func(context.Context, T) error
	errorHandler ErrorHandler[T]
	queue        chan T
	cancel       context.CancelFunc
	wg           sync.WaitGroup
	ctx          context.Context
}

// New creates and returns a new pool with the specified configuration.
// It initializes the internal structures but does not start the worker goroutines.
// - config: Configuration settings for the pool, including max workers and task timeout.
func New[T any](config Config) *Pool[T] {
	ctx, cancel := context.WithCancel(context.Background())
	return &Pool[T]{
		config: config,
		queue:  make(chan T),
		cancel: cancel,
		ctx:    ctx,
	}
}

// SetHandler sets the task handling function for the pool. It must be called before starting the pool.
// The handler function takes a context (for timeout control) and a task of type T, and returns an error.
// - handler: The function that will be called to process each task.
func (p *Pool[T]) SetHandler(handler func(context.Context, T) error) *Pool[T] {
	p.handler = handler
	return p
}

// SetErrorHandler sets a custom error handling function for the pool. It is optional.
// The error handler allows custom logic to be executed when a task processing results in an error.
// - handler: The function to be called when a task encounters an error.
func (p *Pool[T]) SetErrorHandler(handler ErrorHandler[T]) *Pool[T] {
	p.errorHandler = handler
	return p
}

// Start initiates the worker goroutines. It should be called after setting up the task handler and optionally the error handler.
// This method spins up MaxWorkers number of goroutines, each listening for tasks to process.
func (p *Pool[T]) Start() {
	for i := 0; i < p.config.MaxWorkers; i++ {
		p.wg.Add(1)
		go p.worker()
	}
}

func (p *Pool[T]) worker() {
	defer p.wg.Done()
	for {
		select {
		case item, ok := <-p.queue:
			if !ok {
				return
			}
			p.processTask(item)
		case <-p.ctx.Done():
			return
		}
	}
}

func (p *Pool[T]) processTask(item T) {
	var workerCtx context.Context
	var cancel context.CancelFunc
	if p.config.Timeout > 0 {
		workerCtx, cancel = context.WithTimeout(p.ctx, p.config.Timeout)
	} else {
		workerCtx, cancel = context.WithCancel(p.ctx)
	}
	defer cancel()

	done := make(chan error)
	go func() {
		done <- p.handler(workerCtx, item)
	}()

	select {
	case <-workerCtx.Done():
		if !errors.Is(workerCtx.Err(), context.Canceled) {
			p.handleError(workerCtx.Err())
		}
	case err := <-done:
		p.handleError(err)
	}
}

func (p *Pool[T]) handleError(err error) {
	if err != nil && p.errorHandler != nil {
		p.errorHandler(err, p)
	}
}

// Add enqueues a task into the pool. If the pool's worker goroutines are running, the task will be picked up for processing.
// If the pool is not running or has been closed, the behavior of Add is undefined and may result in a deadlock or panic.
// - item: The task to be added to the pool for processing.
func (p *Pool[T]) Add(item T) {
	p.queue <- item
}

// Close gracefully shuts down the pool. It stops accepting new tasks and waits for all ongoing tasks to complete.
// This method should be called to ensure a clean shutdown of the pool.
func (p *Pool[T]) Close() {
	close(p.queue)
	p.wg.Wait() // Wait for all workers to finish processing current tasks
}

// Kill immediately stops all workers in the pool. It cancels the context, causing all worker goroutines to exit.
// Any ongoing tasks may be left unfinished. This method is useful for emergency shutdown scenarios.
func (p *Pool[T]) Kill() {
	p.cancel()  // Cancel the context, signaling all workers to stop
	p.wg.Wait() // Wait for all workers to acknowledge cancellation and exit
}
