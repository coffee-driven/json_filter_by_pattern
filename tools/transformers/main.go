package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"go.uber.org/zap"
)

type Patterns struct {
	NamePattern        string `json:"name_pattern"`
	PricePattern       string `json:"price_pattern"`
	DescriptionPattern string `json:"description_pattern"`
}

type Product struct {
	Product_name        string
	Product_price       float64
	Product_description string
}

func NewProduct() *Product {
	p := &Product{
		Product_name:        "name",
		Product_price:       000000,
		Product_description: "description",
	}

	return p
}

func (p *Product) FilterJsonToProductByPatterns(json_data map[string]interface{},
	patterns Patterns) *Product {

	for key, value := range json_data {
		switch key {
		case patterns.NamePattern:
			p.Product_name = value.(string)
		case patterns.DescriptionPattern:
			p.Product_description = value.(string)
		case patterns.PricePattern:
			p.Product_price = value.(float64)
		default:
			fmt.Printf("Key not matching any of the patterns")

		}
	}

	return p
}

type ProductList struct {
	products []Product
}

func NewProductList() *ProductList {
	var p []Product

	pl := &ProductList{
		products: p,
	}

	return pl
}

func (pl *ProductList) UpdateProductList(json_data []map[string]interface{},
	search_patterns Patterns) *ProductList {

	for _, product := range json_data {
		p := NewProduct()
		product := p.FilterJsonToProductByPatterns(product, search_patterns)

		pl.products = append(pl.products, *product)

	}

	return pl
}

func main() {

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	sugar := logger.Sugar()

	sugar.Info("Hello from zap logger")

	// Load patterns
	var patterns_struct Patterns
	var patterns_file string = "patterns.json"

	patterns, err := os.ReadFile(patterns_file)
	if err != nil {
		sugar.Error("Error reading pattern file")
	}

	if err := json.Unmarshal(patterns, &patterns_struct); err != nil {
		sugar.Error("Error unmarshaling pattern file")
	}

	/*
		my_patterns := make(map[string]string)
		my_patterns["Product_name"] = "Name"
		my_patterns["Product_price"] = "PSPrice"
		my_patterns["Product_description"] = "Notes"

		patterns := sp.newSearchPatterns(my_patterns)
	*/

	// Load data
	var file_name string = "data.json"
	jsonData, err := os.ReadFile(file_name)
	if err != nil {
		sugar.Error("Couldn't read file: %s", file_name)
	}
	var data []map[string]interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		fmt.Printf("Unmarshal failed %s", err)
	}

	// Create convention product list
	pl := NewProductList()
	product_list := pl.UpdateProductList(data, patterns_struct)

	fmt.Println(product_list)
}
