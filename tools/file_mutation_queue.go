// Package tools implements a serialized file mutation queue.
// Ported from superai/pi-cli/packages/coding-agent/src/core/tools/file-mutation-queue.ts
//
// When multiple edit operations target the same file concurrently (e.g., from
// parallel tool calls), the mutation queue ensures they execute serially rather
// than racing. Each file gets its own mutex, so independent files can be
// edited in parallel.
//
// This is critical for correctness: without serialization, two edits to the
// same file could interleave reads and writes, causing one edit to be silently
// lost or applied to stale content.
package tools

import (
	"sync"
)

// fileMutexes provides per-file mutexes for serializing edits.
var fileMutexes sync.Map

// WithFileMutationQueue serializes access to a file path.
// It uses a per-path mutex so that edits to different files can proceed in parallel,
// but edits to the same file are strictly sequential.
func WithFileMutationQueue(filePath string, fn func() (interface{}, error)) (interface{}, error) {
	// Get or create a mutex for this specific file path
	mu, _ := fileMutexes.LoadOrStore(filePath, &sync.Mutex{})
	mutex := mu.(*sync.Mutex)

	mutex.Lock()
	defer mutex.Unlock()

	return fn()
}

// ClearFileMutationQueues removes all queued mutexes (for testing).
func ClearFileMutationQueues() {
	fileMutexes.Range(func(key, _ interface{}) bool {
		fileMutexes.Delete(key)
		return true
	})
}
