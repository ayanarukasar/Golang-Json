package main

import (
	"encoding/json"
	"fmt"

	"io/ioutil"
	"log"
	"net/http"
)

type Product struct {
	ProductId      int    `json:"productId"`
	Manufacturer   string `json:"manufacturer"`
	PricePerUnit   string `json:"pricePerUnit"`
	QuantityOnHand int    `json:"quantityOnHand"`
	ProductName    string `json:"productName"`
}

var productList []Product

func init() {

	// Let's first read the `Last.json` file
	productJSON, err := ioutil.ReadFile("./Last.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}
	err = json.Unmarshal([]byte(productJSON), &productList)
	if err != nil {
		log.Fatal(err)
	}
}

func getNextID() int {
	highestID := -1
	for _, product := range productList {
		if highestID < product.ProductId {
			highestID = product.ProductId
		}
	}
	return highestID + 1
}
func productsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		productsJson, err := json.Marshal(productList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application-json")
		w.Write(productsJson)
	case http.MethodPost:
		//add new product to the list
		var newProduct Product
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(bodyBytes, &newProduct)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if newProduct.ProductId != 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		newProduct.ProductId = getNextID()
		productList = append(productList, newProduct)

		//changes
		for _, x := range productList {
			fmt.Printf("%s \n", x.ProductName)
		}

		w.WriteHeader(http.StatusCreated)
		return

	}
}

func main() {

	http.HandleFunc("/products", productsHandler)
	http.ListenAndServe(":8080", nil)

}
