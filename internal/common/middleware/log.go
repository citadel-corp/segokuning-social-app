package middleware

import (
	"bytes"
	"log/slog"
	"net/http"
	"time"
)

type LogResponseWriter struct {
	http.ResponseWriter
	statusCode int
	buf        bytes.Buffer
}

func NewLogResponseWriter(w http.ResponseWriter) *LogResponseWriter {
	return &LogResponseWriter{ResponseWriter: w}
}

func (w *LogResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *LogResponseWriter) Write(body []byte) (int, error) {
	w.buf.Write(body)
	return w.ResponseWriter.Write(body)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		logRespWriter := NewLogResponseWriter(w)
		next.ServeHTTP(logRespWriter, r)
		slog.Debug("request information",
			slog.Duration("duration", time.Since(startTime)),
			slog.Int("status", logRespWriter.statusCode),
			slog.String("uri", r.RequestURI),
		)
	})
}
