package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
	"github.com/santhozkumar/micro-playlist/product-api/files"
	"github.com/santhozkumar/micro-playlist/product-api/handlers"
)

func main() {
	// l := log.New(os.Stdout, "product-api", log.LstdFlags)
	l := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	store, err := files.NewLocalStorage(files.UploadPath)
	if err != nil {
		return
	}
	hh := handlers.NewHello(l)
	// gh := handlers.NewGoodbye(l)
	ph := handlers.NewProducts(l)
	fh := handlers.NewFiles(store, l)

	r := chi.NewRouter()
	cors_instance := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut},
		AllowCredentials: false,
		MaxAge:           int(time.Hour * 24),
	})
	r.Use(cors_instance.Handler)
	r.Handle("/", hh)
	r.Mount("/product", ph.Routes())

	r.Mount("/file", fh.Routes())



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
			l.Error(err.Error())
		}
	}()

	ch := make(chan os.Signal)

	signal.Notify(ch, os.Kill)
	signal.Notify(ch, os.Interrupt)

	sig := <-ch
	l.Info("Received terminate graceful shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)

}
