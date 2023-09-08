package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
)

type Goodbye struct {
	l *slog.Logger
}

func NewGoodbye(l *slog.Logger) *Goodbye {
	return &Goodbye{l}
}

func (g *Goodbye) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.l.Log(r.Context(), slog.LevelInfo, "goddbye processed")
	fmt.Fprintf(w, "Good bye")

}
