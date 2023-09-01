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
	r.Post("/", p.AddProduct)
	r.Put("/{id:[0-9]+}", p.UpdateProduct)

	return r
}

func (p *Products) AddProduct(w http.ResponseWriter, r *http.Request) {
	// context.WithValue(parent context.Context, key any, val any)
	p.l.Println("Handle POST Product ")
	product := &data.Product{}
	err := product.FromJson(r.Body)
	if err != nil {
		http.Error(w, "Error while reading the json", http.StatusInternalServerError)
		return
	}
	p.l.Printf("%v ", product)
	productID := data.AddProduct(product)
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
	prod := &data.Product{}
	err = prod.FromJson(r.Body)
	if err != nil {
		http.Error(w, "Error while reading the json", http.StatusInternalServerError)
		return
	}

	err = data.UpdateProduct(id, prod)
	if err == data.ErrProductNotFound {
		http.Error(w, "Product Not Found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Product Not Found", http.StatusInternalServerError)
		return

	}

}
