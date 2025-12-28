package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	generated "github.com/egeuysall/drop/internal/supabase/generated"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

var Queries *generated.Queries

func Init(q *generated.Queries) {
	Queries = q
}

func SendJSON(w http.ResponseWriter, message any, statusCode int) {
	w.WriteHeader(statusCode)

	response := map[string]any{"data": message}
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

func Float64ToNumeric(f float64) pgtype.Numeric {
    scaledValue := int64(f * 100)

    return pgtype.Numeric{
        Int: big.NewInt(scaledValue),
        Exp: -2, // Two decimal places
        Valid: true,
    }
}

func Float64PtrToNumeric(f *float64) pgtype.Numeric {
    // If f is nil, return invalid Numeric
    if f == nil {
        return pgtype.Numeric{Valid: false}
    }

    // Otherwise convert it
    return Float64ToNumeric(*f)
}

func BoolPtrToPgBool(b *bool) pgtype.Bool {
    // If b is nil, return invalid Bool
    if b == nil {
        return pgtype.Bool{Valid: false}
    }

    // Otherwise convert it
    return pgtype.Bool{
        Bool: *b,
        Valid: true,
    }
}

func BoolToPgBool(b bool) pgtype.Bool {
    return pgtype.Bool{
        Bool:  b,
        Valid: true,
    }
}

func NumericToFloat64(n pgtype.Numeric) float64 {
    ptr := NumericToFloat64Ptr(n)
    if ptr == nil {
        return 0.0
    }

    return *ptr
}

func NumericToFloat64Ptr(n pgtype.Numeric) *float64 {
    if !n.Valid {
        return nil
    }

    if n.Int == nil {
        return nil
    }

    // Convert big.Int to float64 and apply the exponent
    floatValue := float64(n.Int.Int64())
    if n.Exp < 0 {
        // For negative exponent (like -2), we need to divide
        divisor := float64(1)
        for i := 0; i < int(-n.Exp); i++ {
            divisor *= 10
        }
        floatValue /= divisor
    } else if n.Exp > 0 {
        // For positive exponent, we need to multiply
        multiplier := float64(1)
        for i := 0; i < int(n.Exp); i++ {
            multiplier *= 10
        }
        floatValue *= multiplier
    }

    return &floatValue
}

func PgBoolToBoolPtr(b pgtype.Bool) *bool {
    // If not valid, return nil
    if !b.Valid {
        return nil
    }

    // Return pointer to bool
    return &b.Bool
}

// ValidateURL performs basic URL format validation
func ValidateURL(urlStr string) error {
    // Simple format validation
    parsedURL, err := url.ParseRequestURI(urlStr)
    if err != nil {
        return fmt.Errorf("invalid URL format: %w", err)
    }

    // Basic scheme check
    if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
        return fmt.Errorf("URL must use http:// or https:// protocol")
    }

    // Basic host check
    if parsedURL.Host == "" {
        return fmt.Errorf("URL must include a host")
    }

    return nil
}

func GetEnvironment() string {
	env := os.Getenv("ENV")
	if env == "" {
		return "development"
	}
	return env
}

func ParseDuration(envVar, defaultValue string) time.Duration {
    durationStr := os.Getenv(envVar)

    if durationStr == "" {
        durationStr = defaultValue
    }

    duration, err := time.ParseDuration(durationStr)

    if err != nil {
        log.Fatalf("Error parsing %s: %v", envVar, err)
    }

    return duration
}

func GetEnvInt(envVar string, defaultValue int) int {
    valStr := os.Getenv(envVar)
    if valStr == "" {
        return defaultValue
    }

    val, err := strconv.Atoi(valStr)
    if err != nil {
        log.Printf("Warning: Invalid %s value '%s', using default %d", envVar, valStr, defaultValue)
        return defaultValue
    }

    if val <= 0 {
        log.Printf("Warning: %s must be positive, using default %d", envVar, defaultValue)
        return defaultValue
    }

    return val
}
