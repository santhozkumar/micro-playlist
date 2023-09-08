package handlers

import (
	"compress/gzip"
	"net/http"
	"strings"
)

func GzipMiddleWare(next http.Handler) http.Handler {
	// return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gzip_rw := NewGzipResponseWriter(w)
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
            next.ServeHTTP(w, r)
        } else {
	        gzip_rw.rw.Header().Set("Content-Encoding", "gzip")
			next.ServeHTTP(gzip_rw, r)
			defer gzip_rw.Flush()
		}
	})
}

type GzipResponseWriter struct {
	rw http.ResponseWriter
	gw *gzip.Writer
}

func NewGzipResponseWriter(rw http.ResponseWriter) *GzipResponseWriter {
	gw := gzip.NewWriter(rw)
	return &GzipResponseWriter{rw: rw, gw: gw}
}

func (grw *GzipResponseWriter) Header() http.Header {
	return grw.rw.Header()
}

func (grw *GzipResponseWriter) Write(p []byte) (int, error) {
	return grw.gw.Write(p)
}

func (grw *GzipResponseWriter) WriteHeader(statusCode int) {
	grw.rw.WriteHeader(statusCode)
}
func (grw *GzipResponseWriter) Flush() {
	grw.gw.Flush()
	grw.gw.Close()
}
