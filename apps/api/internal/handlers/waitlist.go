package handlers

import "net/http"

// WaitlistHandler implements the router.WaitlistHandler interface
type WaitlistHandler struct{}

// NewWaitlistHandler creates a new WaitlistHandler instance
func NewWaitlistHandler() *WaitlistHandler {
	return &WaitlistHandler{}
}

// Join handles POST /api/waitlist/join requests
func (h *WaitlistHandler) Join(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement Join handler
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Join handler not implemented"))
}