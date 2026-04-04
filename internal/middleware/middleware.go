// Logging, panic recovery, request-Id, API key auth
package middleware

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/LYD01/my-go-app/internal/respond"
)

type contextKey string

const RequestIdKey contextKey = "request_id"

func RequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := generateId()
		w.Header().Set("X-Request-Id", id)
		ctx := context.WithValue(r.Context(), RequestIdKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
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
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// record start time
		// get a status with the statusRecorder
		// serve the http message
		// log info
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rec, r)

		slog.Info("request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rec.status,
			"duration", time.Since(start).String(),
		)
	})
}

// write a recoverer func that returns an http.Handler
// in that returned handler func body 
// write a deffer annonymous func that 
// checks for nill error after calling recovery()
// slog an error with pannic recovered, set "panic" to v, set stack to string(debug.Stack())
// then respond with an error saying internal server error, curry that
// then call next.ServeHTTP
func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if v := recover(); v != nil {
				slog.Error("panic recovered",
					"panic", v,
					"stack", string(debug.Stack()),
				)
				respond.Error(w, http.StatusInternalServerError, "internal server error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// write a function that returns a function called require api key
// input should be a string 
// returned func type should be a http handler 
// body of this fun should return a func that reutrns the result of calling http.HandlerFunc 
// and passing in an annonymous func 
// inside that get the value of the X-API-Key WriteHeader
// and check if subtle.ConstantTimeCompare([]byte(provided), []byte(apiKey)) does not equal 1
// respond with invalid or missing key
// the call next and server the request
func RequireApiKey(key string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			provided := r.Header.Get("X-API-Key")
			if subtle.ConstantTimeCompare([]byte(provided), []byte(key)) != 1 {
				respond.Error(w, http.StatusUnauthorized, "invalid or missing API key")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
