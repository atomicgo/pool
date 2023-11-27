package pool_test

import (
	"atomicgo.dev/pool"
	"context"
	"fmt"
	"log"
	"time"
)

func Example_demo() {
	// Create a new pool with 3 workers
	p := pool.New[int](pool.Config{
		MaxWorkers: 3,
	})

	// Set the task handler to process integers, simulating a 1-second task
	p.SetHandler(func(ctx context.Context, i int) error {
		log.Printf("Processing %d", i)
		time.Sleep(time.Second)
		return nil
	})

	// Start the pool
	p.Start()

	// Add 10 tasks to the pool
	for i := 0; i < 10; i++ {
		p.Add(i)
	}

	// Close the pool and wait for all tasks to complete
	p.Close()
}

func Example_errorHandling() {
	// Create a new pool
	p := pool.New[int](pool.Config{
		MaxWorkers: 2,
	})

	// Set the task handler, which returns an error for even numbers
	p.SetHandler(func(ctx context.Context, i int) error {
		if i%2 == 0 {
			return fmt.Errorf("error processing %d", i)
		}

		log.Printf("Successfully processed %d", i)

		return nil
	})

	// Set a custom error handler that logs errors
	p.SetErrorHandler(func(err error, pool *pool.Pool[int]) {
		log.Printf("Encountered an error: %v", err)
		// Optional: you can call pool.Kill() here if you want to stop the whole pool on the first error
	})

	// Start the pool
	p.Start()

	// Add 5 tasks to the pool
	for i := 0; i < 5; i++ {
		p.Add(i)
	}

	// Close the pool and wait for all tasks to complete
	p.Close()
}

func Example_timeout() {
	// Create a new pool with a short timeout of 2 seconds
	p := pool.New[int](pool.Config{
		MaxWorkers: 3,
		Timeout:    time.Second * 2,
	})

	// Set the task handler with a task that may exceed the timeout
	p.SetHandler(func(ctx context.Context, i int) error {
		log.Printf("Processing %d", i)
		time.Sleep(time.Second * 3) // This sleep is longer than the timeout
		return nil
	})

	// Start the pool
	p.Start()

	// Add tasks to the pool
	for i := 0; i < 5; i++ {
		p.Add(i)
	}

	// Close the pool and wait for all tasks to complete
	p.Close()
}
