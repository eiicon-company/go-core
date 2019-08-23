package graceful

import (
	"sync"
)

var mu sync.Mutex
var preOnce, postOnce sync.Once
var (
	prehooks  = make([]func(), 0)
	posthooks = make([]func(), 0)
)

// PreHook registers a function to be called before any of this package's normal
// shutdown actions. All listeners will be called in the order they were added,
// from a single goroutine.
func PreHook(f func()) {
	mu.Lock()
	defer mu.Unlock()

	prehooks = append(prehooks, f)
}

// PostHook registers a function to be called after all of this package's normal
// shutdown actions. All listeners will be called in the order they were added,
// from a single goroutine, and are guaranteed to be called after all listening
// connections have been closed, but before Wait() returns.
//
// If you've Hijacked any connections that must be gracefully shut down in some
// other way (since this library disowns all hijacked connections), it's
// reasonable to use a PostHook to signal and wait for them.
func PostHook(f func()) {
	mu.Lock()
	defer mu.Unlock()

	posthooks = append(posthooks, f)
}

// Shutdown shouts down after closing processes.
func Shutdown() {
	preOnce.Do(func() {
		mu.Lock()
		defer mu.Unlock()
		for _, f := range prehooks {
			f()
		}
	})

	postOnce.Do(func() {
		mu.Lock()
		defer mu.Unlock()
		for _, f := range posthooks {
			f()
		}
	})
}
