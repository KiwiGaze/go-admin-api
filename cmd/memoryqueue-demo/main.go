package main

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"go-admin-api/sdk"
	"go-admin-api/storage"
)

func main() {
	runBasicAndUUIDCase()
	runPrefixStampingCase()
	runRetryBackoffCase()
	runRetryCeilingCase()
	runRunShutdownCase()
}

func runBasicAndUUIDCase() {
	stream := "greet"
	q := sdk.Runtime.GetMemoryQueue("basic-host")
	received := make(chan storage.Messager, 1)

	q.Register(stream, func(m storage.Messager) error {
		fmt.Println("got:", m.GetValues())
		received <- m
		return nil
	})

	msg, err := sdk.Runtime.GetStreamMessage("", stream, map[string]interface{}{"hello": "world"})
	if err != nil {
		printFail("basic round-trip", "unexpected GetStreamMessage error: %v", err)
		printFail("UUID assignment", "basic message setup failed before append")
		return
	}

	if err := q.Append(msg); err != nil {
		printFail("basic round-trip", "unexpected Append error: %v", err)
		printFail("UUID assignment", "append failed before delivery")
		return
	}

	select {
	case got := <-received:
		values := got.GetValues()
		if values["hello"] == "world" && got.GetPrefix() == "basic-host" {
			printPass("basic round-trip")
		} else {
			printFail("basic round-trip", "expected hello=world and prefix=basic-host, got values=%v prefix=%q", values, got.GetPrefix())
		}

		if got.GetID() != "" {
			printPass("UUID assignment")
		} else {
			printFail("UUID assignment", "expected non-empty message ID, got empty string")
		}
	case <-time.After(2 * time.Second):
		printFail("basic round-trip", "timed out waiting for consumer delivery")
		printFail("UUID assignment", "timed out before observing message ID")
	}
}

func runPrefixStampingCase() {
	stream := "greet-shared"
	queueA := sdk.Runtime.GetMemoryQueue("tenant-a")
	queueB := sdk.Runtime.GetMemoryQueue("tenant-b")

	var mu sync.Mutex
	prefixes := make([]string, 0, 2)
	delivered := make(chan struct{}, 2)

	queueA.Register(stream, func(m storage.Messager) error {
		fmt.Printf("shared consumer received prefix=%q values=%v\n", m.GetPrefix(), m.GetValues())
		mu.Lock()
		prefixes = append(prefixes, m.GetPrefix())
		mu.Unlock()
		delivered <- struct{}{}
		return nil
	})

	msgA, _ := sdk.Runtime.GetStreamMessage("", stream, map[string]interface{}{"from": "a"})
	msgB, _ := sdk.Runtime.GetStreamMessage("", stream, map[string]interface{}{"from": "b"})
	_ = queueA.Append(msgA)
	_ = queueB.Append(msgB)

	timeout := time.After(2 * time.Second)
	for i := 0; i < 2; i++ {
		select {
		case <-delivered:
		case <-timeout:
			mu.Lock()
			current := append([]string(nil), prefixes...)
			mu.Unlock()
			printFail("prefix stamping", "timed out waiting for both messages, got prefixes=%v", current)
			return
		}
	}

	mu.Lock()
	current := append([]string(nil), prefixes...)
	mu.Unlock()

	if contains(current, "tenant-a") && contains(current, "tenant-b") {
		printPass("prefix stamping")
		return
	}

	printFail("prefix stamping", "expected prefixes tenant-a and tenant-b, got %v", current)
}

func runRetryBackoffCase() {
	stream := "retry-success"
	q := sdk.Runtime.GetMemoryQueue("")
	var attempt int32
	timestamps := make([]time.Time, 0, 3)
	var mu sync.Mutex
	done := make(chan struct{}, 1)

	q.Register(stream, func(m storage.Messager) error {
		now := time.Now()
		currentAttempt := atomic.AddInt32(&attempt, 1)
		mu.Lock()
		timestamps = append(timestamps, now)
		mu.Unlock()
		fmt.Printf("retry attempt %d at %s\n", currentAttempt, now.Format(time.RFC3339Nano))
		if currentAttempt < 3 {
			return errors.New("retry me")
		}
		done <- struct{}{}
		return nil
	})

	msg, _ := sdk.Runtime.GetStreamMessage("", stream, map[string]interface{}{"case": "retry"})
	start := time.Now()
	_ = q.Append(msg)

	select {
	case <-done:
		elapsed := time.Since(start)
		mu.Lock()
		current := append([]time.Time(nil), timestamps...)
		mu.Unlock()
		if atomic.LoadInt32(&attempt) == 3 && elapsed >= 3*time.Second {
			printPass("retry with backoff")
			return
		}
		printFail("retry with backoff", "expected 3 attempts and >=3s elapsed, got attempts=%d elapsed=%s timestamps=%v", atomic.LoadInt32(&attempt), elapsed, current)
	case <-time.After(5 * time.Second):
		mu.Lock()
		current := append([]time.Time(nil), timestamps...)
		mu.Unlock()
		printFail("retry with backoff", "timed out waiting for third attempt, got attempts=%d timestamps=%v", atomic.LoadInt32(&attempt), current)
	}
}

func runRetryCeilingCase() {
	stream := "retry-drop"
	q := sdk.Runtime.GetMemoryQueue("")
	var attempts int32
	timestamps := make([]time.Time, 0, 4)
	var mu sync.Mutex

	q.Register(stream, func(m storage.Messager) error {
		now := time.Now()
		currentAttempt := atomic.AddInt32(&attempts, 1)
		mu.Lock()
		timestamps = append(timestamps, now)
		mu.Unlock()
		fmt.Printf("drop attempt %d at %s\n", currentAttempt, now.Format(time.RFC3339Nano))
		return errors.New("always fail")
	})

	msg, _ := sdk.Runtime.GetStreamMessage("", stream, map[string]interface{}{"case": "drop"})
	_ = q.Append(msg)

	time.Sleep(5500 * time.Millisecond)

	mu.Lock()
	current := append([]time.Time(nil), timestamps...)
	mu.Unlock()

	if atomic.LoadInt32(&attempts) == 3 {
		printPass("retry ceiling")
		return
	}

	printFail("retry ceiling", "expected 3 attempts total, got attempts=%d timestamps=%v", atomic.LoadInt32(&attempts), current)
}

func runRunShutdownCase() {
	q := sdk.Runtime.GetMemoryQueue("")
	done := make(chan struct{})

	go func() {
		q.Run()
		close(done)
	}()

	time.Sleep(200 * time.Millisecond)
	q.Shutdown()

	select {
	case <-done:
		fmt.Println("Run returned after Shutdown")
		printPass("Run/Shutdown parking")
	case <-time.After(2 * time.Second):
		printFail("Run/Shutdown parking", "Run did not return within 2s after Shutdown")
	}
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func printPass(name string) {
	fmt.Printf("PASS: %s\n", name)
}

func printFail(name string, format string, args ...interface{}) {
	fmt.Printf("FAIL: %s - %s\n", name, fmt.Sprintf(format, args...))
}
