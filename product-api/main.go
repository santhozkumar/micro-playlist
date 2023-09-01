package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/santhozkumar/micro-playlist/product-api/handlers"
)

func main() {
	l := log.New(os.Stdout, "product-api", log.LstdFlags)

	hh := handlers.NewHello(l)
	// gh := handlers.NewGoodbye(l)
	ph := handlers.NewProducts(l)

	r := chi.NewRouter()
	r.Handle("/", hh)
	r.Mount("/product", ph.Routes())

	s := &http.Server{
		Addr:         ":9090",
		Handler:      r,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {

		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	ch := make(chan os.Signal)

	signal.Notify(ch, os.Kill)
	signal.Notify(ch, os.Interrupt)

	sig := <-ch
	l.Println("Received terminate graceful shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)

}
