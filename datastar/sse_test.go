package datastar

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSSEContextCancellation(t *testing.T) {
	// Create a test request with a cancellable context
	ctx, cancel := context.WithCancel(context.Background())
	req := httptest.NewRequest("GET", "/test", nil).WithContext(ctx)
	
	// Create a response recorder
	w := httptest.NewRecorder()
	
	// Create SSE handler
	sse := NewSSE(w, req)
	
	// Verify IsClosed returns false initially
	if sse.IsClosed() {
		t.Error("Expected IsClosed to return false for active connection")
	}
	
	// Cancel the context
	cancel()
	
	// Small delay to ensure context cancellation propagates
	time.Sleep(10 * time.Millisecond)
	
	// Verify IsClosed returns true after cancellation
	if !sse.IsClosed() {
		t.Error("Expected IsClosed to return true after context cancellation")
	}
	
	// Try to send an event - should return context error
	err := sse.Send(EventTypePatchElements, []string{"test"})
	if err == nil {
		t.Error("Expected error when sending to cancelled context")
	}
	
	// Verify the error message contains "context cancelled"
	if err != nil && err.Error() != "context cancelled: context canceled" {
		t.Errorf("Expected 'context cancelled' error, got: %v", err)
	}
}

func TestSSESendWithActiveConnection(t *testing.T) {
	// Create a test request
	req := httptest.NewRequest("GET", "/test", nil)
	
	// Create a response recorder
	w := httptest.NewRecorder()
	
	// Create SSE handler
	sse := NewSSE(w, req)
	
	// Send an event - should succeed
	err := sse.Send(EventTypePatchElements, []string{"test data"})
	if err != nil {
		t.Errorf("Expected no error when sending to active connection, got: %v", err)
	}
	
	// Verify headers were set correctly
	contentType := w.Header().Get("Content-Type")
	if contentType != "text/event-stream" {
		t.Errorf("Expected Content-Type to be 'text/event-stream', got: %s", contentType)
	}
}