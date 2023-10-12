package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Для входящего содержимого
		contentEncoding := r.Header.Get("Content-Encoding")
		if contentEncoding == "gzip" {
			reader, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			defer reader.Close()
			r.Body = reader
		}

		// Подготовка к сжатию исходящего содержимого
		acceptEncoding := r.Header.Get("Accept-Encoding")
		if strings.Contains(acceptEncoding, "gzip") {
			gzw := gzip.NewWriter(w)
			defer gzw.Close()

			w.Header().Set("Content-Encoding", "gzip")
			w = &wrapResponseWriter{Writer: gzw, ResponseWriter: w}
		}

		next.ServeHTTP(w, r)
	})
}

type wrapResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *wrapResponseWriter) Write(data []byte) (int, error) {
	return w.Writer.Write(data)
}
