package config

import (
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestWithLock_Basic(t *testing.T) {
	// Create temp HOME directory
	tempHome := t.TempDir()
	oldHome := os.Getenv("HOME")
	_ = os.Setenv("HOME", tempHome)
	t.Cleanup(func() {
		_ = os.Setenv("HOME", oldHome)
	})

	executed := false
	err := WithLock(func() error {
		executed = true
		return nil
	})

	require.NoError(t, err)
	require.True(t, executed)

	// Verify lock file is cleaned up
	lockPath := filepath.Join(tempHome, lockFileName)
	_, err = os.Stat(lockPath)
	require.True(t, os.IsNotExist(err), "lock file should be removed after completion")
}

func TestWithLock_PropagatesError(t *testing.T) {
	tempHome := t.TempDir()
	oldHome := os.Getenv("HOME")
	_ = os.Setenv("HOME", tempHome)
	t.Cleanup(func() {
		_ = os.Setenv("HOME", oldHome)
	})

	expectedErr := os.ErrInvalid
	err := WithLock(func() error {
		return expectedErr
	})

	require.ErrorIs(t, err, expectedErr)
}

func TestWithLock_Concurrent(t *testing.T) {
	tempHome := t.TempDir()
	oldHome := os.Getenv("HOME")
	_ = os.Setenv("HOME", tempHome)
	t.Cleanup(func() {
		_ = os.Setenv("HOME", oldHome)
	})

	// Test that concurrent access is serialized
	var counter int64
	var maxConcurrent int64
	var currentConcurrent int64
	var wg sync.WaitGroup

	numGoroutines := 5
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := WithLock(func() error {
				// Increment current concurrent count
				current := atomic.AddInt64(&currentConcurrent, 1)

				// Track max concurrent
				for {
					max := atomic.LoadInt64(&maxConcurrent)
					if current <= max || atomic.CompareAndSwapInt64(&maxConcurrent, max, current) {
						break
					}
				}

				// Simulate some work
				time.Sleep(10 * time.Millisecond)
				atomic.AddInt64(&counter, 1)

				// Decrement current concurrent count
				atomic.AddInt64(&currentConcurrent, -1)
				return nil
			})
			require.NoError(t, err)
		}()
	}

	wg.Wait()

	require.Equal(t, int64(numGoroutines), counter, "all goroutines should complete")
	require.Equal(t, int64(1), maxConcurrent, "only one goroutine should be in critical section at a time")
}

func TestWithLock_StaleLockRemoval(t *testing.T) {
	tempHome := t.TempDir()
	oldHome := os.Getenv("HOME")
	_ = os.Setenv("HOME", tempHome)
	t.Cleanup(func() {
		_ = os.Setenv("HOME", oldHome)
	})

	// Create a stale lock file
	lockPath := filepath.Join(tempHome, lockFileName)
	err := os.WriteFile(lockPath, []byte("stale"), 0600)
	require.NoError(t, err)

	// Set the modification time to be older than staleLockTimeout
	oldTime := time.Now().Add(-staleLockTimeout - time.Second)
	err = os.Chtimes(lockPath, oldTime, oldTime)
	require.NoError(t, err)

	// WithLock should be able to acquire the lock after removing the stale one
	executed := false
	err = WithLock(func() error {
		executed = true
		return nil
	})

	require.NoError(t, err)
	require.True(t, executed)
}
