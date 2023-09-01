package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type Hello struct {
	l *log.Logger
}

func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

func (h *Hello) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	d, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "unable to read the data", http.StatusBadRequest)
		return
	}
	h.l.Printf("Hello Processed: %s", d)
	fmt.Fprintf(w, "Hello %s", d)

}
