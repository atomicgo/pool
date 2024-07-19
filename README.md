<!--



┌───────────────────────────────────────────────────────────────────┐
│                                                                   │
│                          IMPORTANT NOTE                           │
│                                                                   │
│               This file is automatically generated                │
│           All manual modifications will be overwritten            │
│                                                                   │
└───────────────────────────────────────────────────────────────────┘



-->

<h1 align="center">AtomicGo | pool</h1>

<p align="center">
<img src="https://img.shields.io/endpoint?url=https%3A%2F%2Fatomicgo.dev%2Fapi%2Fshields%2Fpool&style=flat-square" alt="Downloads">

<a href="https://github.com/atomicgo/pool/releases">
<img src="https://img.shields.io/github/v/release/atomicgo/pool?style=flat-square" alt="Latest Release">
</a>

<a href="https://codecov.io/gh/atomicgo/pool" target="_blank">
<img src="https://img.shields.io/github/actions/workflow/status/atomicgo/pool/go.yml?style=flat-square" alt="Tests">
</a>

<a href="https://codecov.io/gh/atomicgo/pool" target="_blank">
<img src="https://img.shields.io/codecov/c/gh/atomicgo/pool?color=magenta&logo=codecov&style=flat-square" alt="Coverage">
</a>

<a href="https://codecov.io/gh/atomicgo/pool">
<!-- unittestcount:start --><img src="https://img.shields.io/badge/Unit_Tests-7-magenta?style=flat-square" alt="Unit test count"><!-- unittestcount:end -->
</a>

<a href="https://opensource.org/licenses/MIT" target="_blank">
<img src="https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square" alt="License: MIT">
</a>
  
<a href="https://goreportcard.com/report/github.com/atomicgo/pool" target="_blank">
<img src="https://goreportcard.com/badge/github.com/atomicgo/pool?style=flat-square" alt="Go report">
</a>   

</p>

---

<p align="center">
<strong><a href="https://pkg.go.dev/atomicgo.dev/pool#section-documentation" target="_blank">Documentation</a></strong>
|
<strong><a href="https://github.com/atomicgo/atomicgo/blob/main/CONTRIBUTING.md" target="_blank">Contributing</a></strong>
|
<strong><a href="https://github.com/atomicgo/atomicgo/blob/main/CODE_OF_CONDUCT.md" target="_blank">Code of Conduct</a></strong>
</p>

---

<p align="center">
  <img src="https://raw.githubusercontent.com/atomicgo/atomicgo/main/assets/header.png" alt="AtomicGo">
</p>

<p align="center">
<table>
<tbody>
</tbody>
</table>
</p>
<h3  align="center"><pre>go get atomicgo.dev/pool</pre></h3>
<p align="center">
<table>
<tbody>
</tbody>
</table>
</p>

<!-- gomarkdoc:embed:start -->

<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# pool

```go
import "atomicgo.dev/pool"
```

Package pool provides a generic and concurrent worker pool implementation. It allows you to enqueue tasks and process them using a fixed number of workers.





```go
package main

import (
	"atomicgo.dev/pool"
	"context"
	"log"
	"time"
)

func main() {
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
```







```go
package main

import (
	"atomicgo.dev/pool"
	"context"
	"fmt"
	"log"
)

func main() {
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
```







```go
package main

import (
	"atomicgo.dev/pool"
	"context"
	"log"
	"time"
)

func main() {
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
```



## Index

- [type Config](<#Config>)
- [type ErrorHandler](<#ErrorHandler>)
- [type Pool](<#Pool>)
  - [func New\[T any\]\(config Config\) \*Pool\[T\]](<#New>)
  - [func \(p \*Pool\[T\]\) Add\(item T\)](<#Pool[T].Add>)
  - [func \(p \*Pool\[T\]\) Close\(\)](<#Pool[T].Close>)
  - [func \(p \*Pool\[T\]\) Kill\(\)](<#Pool[T].Kill>)
  - [func \(p \*Pool\[T\]\) SetErrorHandler\(handler ErrorHandler\[T\]\) \*Pool\[T\]](<#Pool[T].SetErrorHandler>)
  - [func \(p \*Pool\[T\]\) SetHandler\(handler func\(context.Context, T\) error\) \*Pool\[T\]](<#Pool[T].SetHandler>)
  - [func \(p \*Pool\[T\]\) Start\(\)](<#Pool[T].Start>)


<a name="Config"></a>
## type [Config](<https://github.com/atomicgo/pool/blob/main/pool.go#L17-L20>)

Config struct defines the configuration parameters for the Pool. \- MaxWorkers: The maximum number of concurrent workers in the pool. \- Timeout: The maximum duration for processing a single task. If a task takes longer, it will be terminated.

If the Timeout is set to 0, tasks will not be terminated. This is the default behavior. If MaxWorkers is set to 0, it will be set to the number of logical CPUs on the machine.

```go
type Config struct {
    MaxWorkers int
    Timeout    time.Duration
}
```

<a name="ErrorHandler"></a>
## type [ErrorHandler](<https://github.com/atomicgo/pool/blob/main/pool.go#L24>)

ErrorHandler is a function type for handling errors. It provides a way for users to define custom error handling logic. The function receives an error and a pointer to the pool, allowing users to log errors, modify the pool, or perform other actions.

```go
type ErrorHandler[T any] func(error, *Pool[T])
```

<a name="Pool"></a>
## type [Pool](<https://github.com/atomicgo/pool/blob/main/pool.go#L28-L36>)

Pool struct represents a pool of workers processing tasks of type T. It encapsulates the task handling logic, error handling, and synchronization mechanisms.

```go
type Pool[T any] struct {
    // contains filtered or unexported fields
}
```

<a name="New"></a>
### func [New](<https://github.com/atomicgo/pool/blob/main/pool.go#L41>)

```go
func New[T any](config Config) *Pool[T]
```

New creates and returns a new pool with the specified configuration. It initializes the internal structures but does not start the worker goroutines. \- config: Configuration settings for the pool, including max workers and task timeout.

<a name="Pool[T].Add"></a>
### func \(\*Pool\[T\]\) [Add](<https://github.com/atomicgo/pool/blob/main/pool.go#L139>)

```go
func (p *Pool[T]) Add(item T)
```

Add enqueues a task into the pool. If the pool's worker goroutines are running, the task will be picked up for processing. If the pool is not running or has been closed, the behavior of Add is undefined and may result in a deadlock or panic. \- item: The task to be added to the pool for processing.

<a name="Pool[T].Close"></a>
### func \(\*Pool\[T\]\) [Close](<https://github.com/atomicgo/pool/blob/main/pool.go#L145>)

```go
func (p *Pool[T]) Close()
```

Close gracefully shuts down the pool. It stops accepting new tasks and waits for all ongoing tasks to complete. This method should be called to ensure a clean shutdown of the pool.

<a name="Pool[T].Kill"></a>
### func \(\*Pool\[T\]\) [Kill](<https://github.com/atomicgo/pool/blob/main/pool.go#L152>)

```go
func (p *Pool[T]) Kill()
```

Kill immediately stops all workers in the pool. It cancels the context, causing all worker goroutines to exit. Any ongoing tasks may be left unfinished. This method is useful for emergency shutdown scenarios.

<a name="Pool[T].SetErrorHandler"></a>
### func \(\*Pool\[T\]\) [SetErrorHandler](<https://github.com/atomicgo/pool/blob/main/pool.go#L70>)

```go
func (p *Pool[T]) SetErrorHandler(handler ErrorHandler[T]) *Pool[T]
```

SetErrorHandler sets a custom error handling function for the pool. It is optional. The error handler allows custom logic to be executed when a task processing results in an error. \- handler: The function to be called when a task encounters an error.

<a name="Pool[T].SetHandler"></a>
### func \(\*Pool\[T\]\) [SetHandler](<https://github.com/atomicgo/pool/blob/main/pool.go#L62>)

```go
func (p *Pool[T]) SetHandler(handler func(context.Context, T) error) *Pool[T]
```

SetHandler sets the task handling function for the pool. It must be called before starting the pool. The handler function takes a context \(for timeout control\) and a task of type T, and returns an error. \- handler: The function that will be called to process each task.

<a name="Pool[T].Start"></a>
### func \(\*Pool\[T\]\) [Start](<https://github.com/atomicgo/pool/blob/main/pool.go#L77>)

```go
func (p *Pool[T]) Start()
```

Start initiates the worker goroutines. It should be called after setting up the task handler and optionally the error handler. This method spins up MaxWorkers number of goroutines, each listening for tasks to process.

Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)


<!-- gomarkdoc:embed:end -->

---

> [AtomicGo.dev](https://atomicgo.dev) &nbsp;&middot;&nbsp;
> with ❤️ by [@MarvinJWendt](https://github.com/MarvinJWendt) |
> [MarvinJWendt.com](https://marvinjwendt.com)
