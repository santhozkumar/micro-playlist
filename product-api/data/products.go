package data

import (
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

type Product struct {
	ID          int     `json:"id"`
    Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
    Price       float32 `json:"price" validate:"gt=0"`
    SKU         string  `json:"sku" validate:"required,sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

type Products []*Product

func (p *Products) ToJson(w io.Writer) error {
	return ToJson(p, w)
}

func (p *Product) FromJson(r io.Reader) error {
	return FromJson(p, r)
}

func (p *Product) Validate() error {
	v := validator.New(validator.WithRequiredStructEnabled())
	v.RegisterValidation("sku", skuValidator)
	err := v.Struct(p)
	if err != nil {
		return err
	}
	return nil 
}

func skuValidator(fl validator.FieldLevel) bool {
    reg := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
    matches := reg.FindAllString(fl.Field().String(), -1)
    if len(matches) != 1{
        return false
    }
	return true
}

func GetProducts() Products {
	time.Sleep(time.Millisecond * 800)
	return productList
}

func AddProduct(p *Product) int {
	p.ID = getNextID()
	productList = append(productList, p)
	return p.ID
}

func UpdateProduct(id int, p *Product) error {
	_, pos, err := findProduct(id)
	if err != nil {
		return err
	}
	p.ID = id
	productList[pos] = p

	return nil
}

var ErrProductNotFound = fmt.Errorf("Product Not Foound")

func findProduct(id int) (*Product, int, error) {
	for i, product := range productList {
		if product.ID == id {
			return product, i, nil
		}
	}

	return nil, -1, ErrProductNotFound
}

func getNextID() int {
	lastProduct := productList[len(productList)-1]
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
