package data

import (
	"io"
	"time"
)

type Product struct {
    ID int `json:"id"` 
    Name string `json:"name"`
    Description string `json:"description"`
    Price float32 `json:"price"`
    SKU string `json:"sku"`
    CreatedOn string `json:"-"`
    UpdatedOn string `json:"-"`
    DeletedOn string `json:"-"`
}

type Products []*Product


func (p *Products)ToJson(w io.Writer) error  {
    return ToJson(p, w)
}

func (p *Product)FromJson(r io.Reader) error  {
    return FromJson(p, r)
}

func GetProducts() Products {
    return productList
}


func AddProduct(p *Product) int {
    p.ID = getNextID() 
    productList = append(productList, p)
    return p.ID
}


func getNextID() int {
    lastProduct := productList[len(productList) -1]
    return lastProduct.ID + 1
}


var productList = Products{
    	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc323",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	&Product{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fjd34",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}



