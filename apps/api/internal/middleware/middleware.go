package middleware

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/egeuysall/drop/internal/utils"
	"github.com/go-chi/cors"
	"github.com/golang-jwt/jwt/v5"
)

// Context

type contextKey string

const userIDKey = contextKey("userID")

// JWKS Cache

type jwksCache struct {
	keys      map[string]*ecdsa.PublicKey
	expiresAt time.Time
	mu        sync.RWMutex
}

var jwks = &jwksCache{
	keys: make(map[string]*ecdsa.PublicKey),
}

// JWKS Fetching

func fetchJWKS(issuer string) (map[string]*ecdsa.PublicKey, error) {
	issuer = strings.TrimRight(issuer, "/")
	jwksURL := issuer + "/auth/v1/.well-known/jwks.json"

	resp, err := http.Get(jwksURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("JWKS endpoint returned %d", resp.StatusCode)
	}

	var payload struct {
		Keys []struct {
			KID string `json:"kid"`
			KTY string `json:"kty"`
			Alg string `json:"alg"`
			Crv string `json:"crv"`
			X   string `json:"x"`
			Y   string `json:"y"`
		} `json:"keys"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	keys := make(map[string]*ecdsa.PublicKey)

	for _, k := range payload.Keys {
		if k.KTY != "EC" || k.Alg != "ES256" || k.Crv != "P-256" {
			continue
		}

		xBytes, err := base64.RawURLEncoding.DecodeString(k.X)
		if err != nil {
			log.Printf("Failed to decode X coordinate for kid %s: %v", k.KID, err)
			continue
		}
		yBytes, err := base64.RawURLEncoding.DecodeString(k.Y)
		if err != nil {
			log.Printf("Failed to decode Y coordinate for kid %s: %v", k.KID, err)
			continue
		}

		keys[k.KID] = &ecdsa.PublicKey{
			Curve: elliptic.P256(),
			X:     new(big.Int).SetBytes(xBytes),
			Y:     new(big.Int).SetBytes(yBytes),
		}
	}

	if len(keys) == 0 {
		return nil, fmt.Errorf("no valid ES256 keys found in JWKS")
	}

	return keys, nil
}

func getPublicKey(issuer, kid string) (*ecdsa.PublicKey, error) {
	jwks.mu.RLock()
	if time.Now().Before(jwks.expiresAt) {
		if key, ok := jwks.keys[kid]; ok {
			jwks.mu.RUnlock()
			return key, nil
		}
	}
	jwks.mu.RUnlock()

	keys, err := fetchJWKS(issuer)
	if err != nil {
		return nil, err
	}

	jwks.mu.Lock()
	jwks.keys = keys
	jwks.expiresAt = time.Now().Add(10 * time.Minute)
	jwks.mu.Unlock()

	key, ok := keys[kid]
	if !ok {
		return nil, fmt.Errorf("kid %s not found in JWKS", kid)
	}

	return key, nil
}

// Auth Middleware

func RequireAuth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			issuer := os.Getenv("SUPABASE_ISSUER")
			if issuer == "" {
				log.Println("SUPABASE_ISSUER not set")
				utils.SendError(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			audience := os.Getenv("SUPABASE_AUDIENCE")
			if audience == "" {
				audience = "authenticated"
			}

			auth := r.Header.Get("Authorization")
			if auth == "" {
				utils.SendError(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(auth, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				utils.SendError(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			tokenStr := parts[1]

			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
				if token.Method != jwt.SigningMethodES256 {
					return nil, fmt.Errorf("unexpected signing method: %s", token.Method.Alg())
				}

				kid, ok := token.Header["kid"].(string)
				if !ok {
					return nil, fmt.Errorf("missing kid in token header")
				}

				return getPublicKey(issuer, kid)
			})

			if err != nil || !token.Valid {
				log.Printf("JWT validation error: %v", err)
				utils.SendError(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				utils.SendError(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if iss, ok := claims["iss"].(string); !ok || !strings.HasPrefix(iss, issuer) {
				utils.SendError(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			validAud := false
			switch aud := claims["aud"].(type) {
			case string:
				validAud = aud == audience
			case []any:
				for _, v := range aud {
					if s, ok := v.(string); ok && s == audience {
						validAud = true
						break
					}
				}
			}

			if !validAud {
				utils.SendError(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			sub, ok := claims["sub"].(string)
			if !ok || sub == "" {
				utils.SendError(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, sub)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Helpers

func UserIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(userIDKey).(string)
	return id, ok
}

func Cors() func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://www.drop.egeuysal.com", "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           3600,
	})
}

func SetContentType() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	}
}
