package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/santhozkumar/micro-playlist/product-api/data"
)


type Products struct {
    l *log.Logger
}


func NewProducts(l *log.Logger) *Products {
    return &Products{l}
}


func (p *Products) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        p.getProduct(w,r) 
        return
    }
    if r.Method == http.MethodPost {
        p.addProduct(w,r)
        return 
    }
    if r.Method == http.MethodPut {
        //expect the id in the URI
        // path := r.URL.Path
        http.Error(w, "client side error", http.StatusBadRequest)
        return 
    }
}


func (p *Products) addProduct(w http.ResponseWriter, r *http.Request) {
    p.l.Println("Handle POST Product ")
    product := &data.Product{}
    err := product.FromJson(r.Body)
    if err != nil  {
        http.Error(w, "Error while reading the json", http.StatusInternalServerError)
        return
    }
    p.l.Printf("%v ", product)
    productID := data.AddProduct(product)
    fmt.Fprintf(w, "Id: %d", productID)
}

func (p *Products) getProduct(w http.ResponseWriter, r *http.Request) {
    p.l.Println("Handle GET Product ")
    lp := data.GetProducts()

    err := lp.ToJson(w)
    if err != nil {
        http.Error(w, "Error while reading the json", http.StatusInternalServerError)
        return
    }
}

