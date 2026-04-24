package tools

import (
	"sync"
	"testing"
)

func TestFileMutationQueueSerializes(t *testing.T) {
	ClearFileMutationQueues()

	counter := 0
	var mu sync.Mutex

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = WithFileMutationQueue("test.txt", func() (interface{}, error) {
				mu.Lock()
				counter++
				mu.Unlock()
				return nil, nil
			})
		}()
	}
	wg.Wait()

	if counter != 10 {
		t.Errorf("expected 10 mutations, got %d", counter)
	}
}

func TestFileMutationQueueParallelPaths(t *testing.T) {
	ClearFileMutationQueues()

	var wg sync.WaitGroup
	completed := make(chan string, 4)

	// Different paths should run in parallel
	for _, path := range []string{"a.txt", "b.txt", "c.txt"} {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			_, _ = WithFileMutationQueue(p, func() (interface{}, error) {
				completed <- p
				return nil, nil
			})
		}(path)
	}

	wg.Wait()
	close(completed)

	count := 0
	for range completed {
		count++
	}
	if count != 3 {
		t.Errorf("expected 3 completions, got %d", count)
	}
}
