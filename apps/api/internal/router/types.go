package router

import "net/http"

type Handlers struct {
	Items         ItemsHandler
}

// AuthMiddleware holds the auth middleware configuration
type AuthMiddleware struct {
	JWTSecret   string
	SupabaseURL string
}

type ItemsHandler interface {
	CreateItem(w http.ResponseWriter, r *http.Request)
	GetItem(w http.ResponseWriter, r *http.Request)
	ListItems(w http.ResponseWriter, r *http.Request)
	UpdateItem(w http.ResponseWriter, r *http.Request)
	DeleteItem(w http.ResponseWriter, r *http.Request)
	UpdateItemPrice(w http.ResponseWriter, r *http.Request)
	GetPriceHistory(w http.ResponseWriter, r *http.Request)
	GetPriceStats(w http.ResponseWriter, r *http.Request)
	CheckPriceDrop(w http.ResponseWriter, r *http.Request)
}
