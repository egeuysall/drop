package utils

import (
	"encoding/json"
	"log"
	"net/http"

	generated "github.com/egeuysall/drop/internal/supabase/generated"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

var Queries *generated.Queries

func Init(q *generated.Queries) {
	Queries = q
}

func SendJson(w http.ResponseWriter, message interface{}, statusCode int) {
	w.WriteHeader(statusCode)

	response := map[string]interface{}{"data": message}
	err := json.NewEncoder(w).Encode(response)

	if err != nil {
		SendError(w, "Failed to encode JSON response", http.StatusInternalServerError)
	}
}

func CheckStatus(w http.ResponseWriter, message interface{}, statusCode int) {
	w.WriteHeader(statusCode)

	response := map[string]interface{}{"status": message}
	err := json.NewEncoder(w).Encode(response)

	if err != nil {
		SendError(w, "Failed to encode JSON response", http.StatusInternalServerError)
	}
}

func SendError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)

	errorResponse := map[string]string{"error": message}
	err := json.NewEncoder(w).Encode(errorResponse)

	if err != nil {
		log.Printf("SendError encoding failed: %v", err)
	}
}

func ParseUUID(str string) (pgtype.UUID, error) {
	var id pgtype.UUID
	err := id.Scan(str)
	return id, err
}

func UUIDToString(u pgtype.UUID) string {
	if !u.Valid {
		return ""
	}
	return uuid.UUID(u.Bytes).String()
}
