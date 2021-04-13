package internalhttp

import (
	"log"
	"net/http"
	"time"

	"github.com/FedoseevAlex/OtusGolangHomeWork/hw12_13_14_15_calendar/internal/app"
)

type ResponseWriterWrapper struct {
	http.ResponseWriter
	status int
}

func (rw *ResponseWriterWrapper) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func loggingMiddleware(next http.Handler, logger app.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrappedResponseWriter := &ResponseWriterWrapper{ResponseWriter: w, status: http.StatusOK}
		log.Println("wrappedResponseWriter ", wrappedResponseWriter)

		begin := time.Now()
		next.ServeHTTP(wrappedResponseWriter, r)
		duration := time.Since(begin)

		params := map[string]interface{}{
			"ip":          r.RemoteAddr,
			"timestamp":   begin.Format(time.RFC822Z),
			"method":      r.Method,
			"path":        r.URL.Path,
			"HTTP ver.":   r.Proto,
			"status code": wrappedResponseWriter.status,
			"latency":     duration.String(),
			"user agent":  r.UserAgent(),
		}
		logger.Trace("Request", params)
	})
}
