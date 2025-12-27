package items

import (
	"encoding/json"
	"fmt"
	"net/http"

	appmid "github.com/egeuysall/drop/internal/middleware"
	"github.com/egeuysall/drop/internal/utils"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// CreateItem handles POST /items - Create a new tracked item
func (h *Handler) CreateItem(w http.ResponseWriter, r *http.Request) {
	var req CreateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, ok := appmid.UserIDFromContext(r.Context())
	if !ok {
		utils.SendError(w, "Unauthorized: user not authenticated", http.StatusUnauthorized)
		return
	}

	response, err := h.service.CreateItem(r.Context(), userID, req)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendJSON(w, response, http.StatusCreated)
}

// GetItem handles GET /items/{id} - Get a specific item
func (h *Handler) GetItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID, ok := appmid.UserIDFromContext(r.Context())
	if !ok {
		utils.SendError(w, "Unauthorized: user not authenticated", http.StatusUnauthorized)
		return
	}

	response, err := h.service.GetItemByID(r.Context(), id, userID)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusNotFound)
		return
	}

	utils.SendJSON(w, response, http.StatusOK)
}

// ListItems handles GET /items - List all user's items
func (h *Handler) ListItems(w http.ResponseWriter, r *http.Request) {
	userID, ok := appmid.UserIDFromContext(r.Context())
	if !ok {
		utils.SendError(w, "Unauthorized: user not authenticated", http.StatusUnauthorized)
		return
	}

	responses, err := h.service.ListItemsByUserID(r.Context(), userID)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendJSON(w, ListItemsResponse{Items: responses, Total: len(responses)}, http.StatusOK)
}

// UpdateItem handles PUT /items/{id} - Update an item
func (h *Handler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID, ok := appmid.UserIDFromContext(r.Context())
	if !ok {
		utils.SendError(w, "Unauthorized: user not authenticated", http.StatusUnauthorized)
		return
	}

	var req UpdateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response, err := h.service.UpdateItemByID(r.Context(), id, userID, req)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendJSON(w, response, http.StatusOK)
}

// DeleteItem handles DELETE /items/{id} - Delete an item
func (h *Handler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID, ok := appmid.UserIDFromContext(r.Context())
	if !ok {
		utils.SendError(w, "Unauthorized: user not authenticated", http.StatusUnauthorized)
		return
	}

	if err := h.service.DeleteItem(r.Context(), id, userID); err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdateItemPrice handles PATCH /items/{id}/price - Update item price
func (h *Handler) UpdateItemPrice(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID, ok := appmid.UserIDFromContext(r.Context())
	if !ok {
		utils.SendError(w, "Unauthorized: user not authenticated", http.StatusUnauthorized)
		return
	}

	var req UpdatePriceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response, err := h.service.UpdateItemPrice(r.Context(), id, userID, req.CurrentPrice, req.InStock)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendJSON(w, response, http.StatusOK)
}

// GetPriceHistory handles GET /items/{id}/history - Get price history
func (h *Handler) GetPriceHistory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	days := 30 // default

	if daysParam := r.URL.Query().Get("days"); daysParam != "" {
		fmt.Sscanf(daysParam, "%d", &days)
	}

	history, err := h.service.GetPriceHistory(r.Context(), id, days)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendJSON(w, history, http.StatusOK)
}

// GetPriceStats handles GET /items/{id}/stats - Get price statistics
func (h *Handler) GetPriceStats(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	stats, err := h.service.GetPriceStats(r.Context(), id)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendJSON(w, stats, http.StatusOK)
}

// CheckPriceDrop handles GET /items/{id}/check - Check for price drops
func (h *Handler) CheckPriceDrop(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID, ok := appmid.UserIDFromContext(r.Context())
	if !ok {
		utils.SendError(w, "Unauthorized: user not authenticated", http.StatusUnauthorized)
		return
	}

	check, err := h.service.CheckPriceDrop(r.Context(), id, userID)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendJSON(w, check, http.StatusOK)
}
