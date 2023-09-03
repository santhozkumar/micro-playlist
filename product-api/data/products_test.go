package data

import "testing"


func TestValidate(t *testing.T){

    product := Product{Name:"tea", SKU: "abc-def-ghi", Price: 10.0}
    if err:=  product.Validate(); err != nil{
        t.Fatal(err)
}
}


