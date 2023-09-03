// Package classification at product api
//
// Documentation for Product API
// 
// Schemes: http
// Host: localhost
// BasePath: /product
// version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//swagger:meta



package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/santhozkumar/micro-playlist/product-api/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) Routes() chi.Router {

	r := chi.NewRouter()
	r.Get("/", p.GetProduct)
	r.With(ProductValidationMiddleWare).Post("/", p.AddProduct)
	r.With(ProductValidationMiddleWare).Put("/{id:[0-9]+}", p.UpdateProduct)

	r.Route("/{id:[0-9]+}", func(r chi.Router) {
		// r.Use(ProductCtx)
		// r.Get("/", h http.HandlerFunc)

	})
    r.Options("/", http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
        w.Header().Set("Allow", fmt.Sprintf("%s, %s, %s, %s", http.MethodGet, http.MethodPut, http.MethodPost, http.MethodOptions ))
        w.WriteHeader(http.StatusNoContent)
    }))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	return r
}

func (p *Products) AddProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Product ")

	prod := r.Context().Value(ProductKey{}).(data.Product)
	productID := data.AddProduct(&prod)
	fmt.Fprintf(w, "Id: %d", productID)
}

func (p *Products) GetProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Product ")
	ctx, cancel := context.WithTimeout(r.Context(), time.Millisecond*1000)
	defer cancel()

	productsChannel := make(chan data.Products)
	r = r.WithContext(ctx)

	go func() {
		productsChannel <- data.GetProducts()
	}()

	var lp data.Products
	for {
		select {
		case <-r.Context().Done():
			http.Error(w, "took to long for the product", http.StatusInternalServerError)
			return
		case products := <-productsChannel:
			lp = products
			err := lp.ToJson(w)
			if err != nil {
				http.Error(w, "Error while reading the json", http.StatusInternalServerError)
				return
			}
			return
		}
	}

}

func (p *Products) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT Product ")
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid Integer", http.StatusBadRequest)
		return
	}

	prod := r.Context().Value(ProductKey{}).(data.Product)

	err = data.UpdateProduct(id, &prod)
	if err == data.ErrProductNotFound {
		http.Error(w, "Product Not Found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "Product Not Found", http.StatusInternalServerError)
		return

	}

}

type ProductKey struct{}

func ProductValidationMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prod := data.Product{}
		err := prod.FromJson(r.Body)
		if err != nil {
			http.Error(w, "Error while reading the json", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), ProductKey{}, prod)
		fmt.Println(prod)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
