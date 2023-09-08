package handlers

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type Hello struct {
	l *slog.Logger
}

func NewHello(l *slog.Logger) *Hello {
	return &Hello{l}
}

func (h *Hello) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	d, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "unable to read the data", http.StatusBadRequest)
		return
	}

	h.l.Log(r.Context(), slog.LevelInfo, fmt.Sprintf("Hello Processed: %s", d))
	fmt.Fprintf(w, "Hello %s", d)

}
