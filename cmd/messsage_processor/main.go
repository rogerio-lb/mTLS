package main

import (
	"fmt"
	"sync"
	"time"
)

// MessageBatcher handles batching messages over a time interval
type MessageBatcher struct {
	batchInterval time.Duration
	processor     func([]string) // Function to process the batch of messages
	messages      []string
	mutex         sync.Mutex
	timer         *time.Timer
}

// NewMessageBatcher creates a new message batcher
func NewMessageBatcher(interval time.Duration, processor func([]string)) *MessageBatcher {
	return &MessageBatcher{
		batchInterval: interval,
		processor:     processor,
		messages:      make([]string, 0),
	}
}

// AddMessage adds a message to the batch
func (mb *MessageBatcher) AddMessage(message string) {
	mb.mutex.Lock()
	defer mb.mutex.Unlock()

	// Add message to the current batch
	mb.messages = append(mb.messages, message)

	// If this is the first message in the batch, start the timer
	if len(mb.messages) == 1 {
		mb.timer = time.AfterFunc(mb.batchInterval, mb.processBatch)
	}

	if len(mb.messages) == 10 {
		mb.timer.Stop()
		go mb.processBatch()
	}
}

// processBatch processes the current batch of messages
func (mb *MessageBatcher) processBatch() {
	mb.mutex.Lock()
	defer mb.mutex.Unlock()

	if len(mb.messages) > 0 {
		// Create a copy of messages to send
		batch := make([]string, len(mb.messages))
		copy(batch, mb.messages)

		// Clear the messages slice for the next batch
		mb.messages = mb.messages[:0]

		// Process the batch (call this in a goroutine to avoid blocking)
		go mb.processor(batch)
	}
}

// Flush immediately processes any pending messages
func (mb *MessageBatcher) Flush() {
	mb.mutex.Lock()
	defer mb.mutex.Unlock()

	if mb.timer != nil {
		mb.timer.Stop()
	}

	if len(mb.messages) > 0 {
		batch := make([]string, len(mb.messages))
		copy(batch, mb.messages)
		mb.messages = mb.messages[:0]
		go mb.processor(batch)
	}
}

// Example processor function
func processMessageBatch(messages []string) {
	fmt.Printf("Processing batch of %d messages:\n", len(messages))
	for i, msg := range messages {
		fmt.Printf("  %d: %s\n", i+1, msg)
	}
	fmt.Println("Batch processing completed")
	fmt.Println("---")
}

func main() {
	// Create a message batcher with 5-second interval and max 10 messages
	batcher := NewMessageBatcher(5*time.Second, processMessageBatch)

	// Simulate receiving messages
	fmt.Println("Starting message batcher demo...")
	fmt.Println("Messages will be batched every 5 seconds OR when 10 messages are reached")
	fmt.Println()

	// Send some messages quickly (less than 10)
	fmt.Println("Sending 5 messages quickly...")
	for i := 1; i <= 5; i++ {
		batcher.AddMessage(fmt.Sprintf("Quick message %d", i))
		time.Sleep(100 * time.Millisecond)
	}

	// Wait for the time-based batch to be processed
	fmt.Println("Waiting for time-based processing...")
	time.Sleep(6 * time.Second)

	// Send exactly 10 messages to trigger size-based processing
	fmt.Println("Sending 10 messages to trigger size-based processing...")
	for i := 1; i <= 10; i++ {
		batcher.AddMessage(fmt.Sprintf("Size-trigger message %d", i))
		time.Sleep(50 * time.Millisecond)
	}

	// Send a few more after the batch was processed
	time.Sleep(1 * time.Second)
	fmt.Println("Sending 3 more messages...")
	for i := 1; i <= 3; i++ {
		batcher.AddMessage(fmt.Sprintf("Final message %d", i))
	}

	// Demonstrate flush
	time.Sleep(1 * time.Second)
	fmt.Println("Flushing remaining messages...")
	batcher.Flush()

	// Give time for processing to complete
	time.Sleep(1 * time.Second)
	fmt.Println("Demo completed")
}
