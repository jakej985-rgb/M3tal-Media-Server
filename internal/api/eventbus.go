package api

import (
	"sync"
	"time"
)

// Event represents a system-wide event payload.
type Event struct {
	Type      string `json:"type"`
	Payload   any    `json:"payload,omitempty"`
	Timestamp int64  `json:"timestamp"`
}

// EventBus implements a thread-safe publish-subscribe pattern.
type EventBus struct {
	mu          sync.RWMutex
	subscribers map[chan Event]bool
}

// NewEventBus creates a new EventBus.
func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[chan Event]bool),
	}
}

// Subscribe registers a new channel to receive published events.
func (eb *EventBus) Subscribe() chan Event {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	ch := make(chan Event, 100)
	eb.subscribers[ch] = true
	return ch
}

// Unsubscribe removes a channel from the subscribers list.
func (eb *EventBus) Unsubscribe(ch chan Event) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	if _, ok := eb.subscribers[ch]; ok {
		delete(eb.subscribers, ch)
		close(ch)
	}
}

// Publish broadcasts an event to all active subscribers.
func (eb *EventBus) Publish(eventType string, payload any) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	ev := Event{
		Type:      eventType,
		Payload:   payload,
		Timestamp: time.Now().Unix(),
	}

	for ch := range eb.subscribers {
		select {
		case ch <- ev:
		default:
			// Buffer full, drop event to prevent blocking publisher
		}
	}
}

// GlobalEventBus is the shared bus for the entire application.
var GlobalEventBus = NewEventBus()
