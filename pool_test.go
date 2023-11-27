package pool

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestPoolRun(t *testing.T) {
	configs := []Config{
		{
			MaxWorkers: 2,
		},
		{
			MaxWorkers: 2,
			Timeout:    time.Millisecond * 10,
		},
		{}, // Empty config
	}

	for _, config := range configs {
		t.Run(fmt.Sprintf("%+v", config), func(t *testing.T) {
			testPool(t, config, 0, true, false)
		})
	}
}

func TestPoolTimeout(t *testing.T) {
	configs := []Config{
		{
			MaxWorkers: 2,
			Timeout:    time.Millisecond * 10,
		},
	}

	for _, config := range configs {
		t.Run(fmt.Sprintf("%+v", config), func(t *testing.T) {
			testPool(t, config, time.Millisecond*100, false, true)
		})
	}
}

func TestErrorHandling(t *testing.T) {
	config := Config{
		MaxWorkers: 2,
	}

	testPool(t, config, 0, true, true)
}

func testPool(t *testing.T, config Config, writeSpeed time.Duration, shouldPass bool, expectsError bool) {
	t.Helper()

	// Create map with 10 booleans
	var m sync.Map
	for i := 0; i < 10; i++ {
		m.Store(i, false)
	}

	// Create a new pool
	p := New[int](config)

	// Set the task handler to process integers
	p.SetHandler(func(ctx context.Context, i int) error {
		// Simulate write speed
		if writeSpeed > 0 {
			time.Sleep(writeSpeed)
		}

		// Set the map value to true
		m.Store(i, true)

		if expectsError {
			return fmt.Errorf("error")
		}

		return nil
	})

	var hasErrors bool

	// Set error handler
	p.SetErrorHandler(func(err error, p *Pool[int]) {
		hasErrors = true
	})

	// Start the pool
	p.Start()

	// Add map keys to the pool
	for i := 0; i < 10; i++ {
		p.Add(i)
	}

	// Close the pool and wait for all tasks to complete
	p.Close()

	m.Range(func(key, value interface{}) bool {
		if !value.(bool) && shouldPass {
			t.Errorf("Expected map value to be true, got false")
		} else if value.(bool) && !shouldPass {
			t.Errorf("Expected map value to be false, got true")
		}

		return true
	})

	if expectsError && !hasErrors {
		t.Errorf("Expected errors to be true, got false")
	} else if !expectsError && hasErrors {
		t.Errorf("Expected errors to be false, got true")
	}
}
