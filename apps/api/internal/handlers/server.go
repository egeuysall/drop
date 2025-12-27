package handlers

import (
	"net/http"

	"github.com/egeuysall/drop/internal/utils"
)

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	utils.CheckStatus(w, "Drop API", http.StatusOK)
}

func HandleHealth(w http.ResponseWriter, r *http.Request) {
	utils.CheckStatus(w, "Healthy", http.StatusOK)
}
