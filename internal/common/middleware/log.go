package middleware

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/citadel-corp/segokuning-social-app/internal/common/id"
	"github.com/citadel-corp/segokuning-social-app/internal/common/response"
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
		requestID := id.GenerateStringID(12)
		logRespWriter := NewLogResponseWriter(w)
		next.ServeHTTP(logRespWriter, r)

		slog.Debug("request information",
			slog.Duration("duration", time.Since(startTime)),
			slog.Int("status", logRespWriter.statusCode),
			slog.String("uri", r.RequestURI),
			slog.String("requestID", requestID),
			slog.String("method", r.Method),
		)
		var resp response.ResponseBody
		err := json.NewDecoder(&logRespWriter.buf).Decode(&resp)
		if logRespWriter.statusCode >= 500 && err == nil {
			slog.Error("internal server error on request",
				slog.String("error", resp.Error),
			)
		}
	})
}
