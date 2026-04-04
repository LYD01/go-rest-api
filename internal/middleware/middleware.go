// Logging, panic recovery, request-Id, API key auth
package middleware

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"log/slog"
	"net/http"
	"time"
	"go-api/internal/respond"
)

type contextKey string

const RequestIdKey contextKey = "request_id"


func RequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request)) {
		id := generateId()
		w.Header().Set("X-Request-Id", id)
		ctx := r.Context()
		ctx := cnotext.WithValue(ctx, RequestIdKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func generateId() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (sr *statusRecorder) WriteHeader(code int) {
	sr.status = code
	sr.ResponseWriter.WriteHeader(code)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w, http.ResponseWriter, r *http.Request) {
	
	})	
}
