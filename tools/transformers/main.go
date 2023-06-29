package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"

	"go.uber.org/zap"
)

type SearchPatterns struct {
	patterns map[string]string
}

func (sp *SearchPatterns) newSearchPatterns(patterns map[string]string) *SearchPatterns {
	p := SearchPatterns{patterns: patterns}
	return &p
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
	patterns SearchPatterns) *Product {

	for key, value := range json_data {
		for pattern_key, pattern_value := range patterns.patterns {

			if key == pattern_value {

				switch typeof := reflect.TypeOf(value).Kind(); typeof {
				case reflect.Float64:
					value_typed := value.(float64)
					p.Product_price = value_typed
				case reflect.String:
					value_typed := value.(string)

					switch pattern_key {
					case "Product_name":
						p.Product_name = value_typed
					case "Product_description":
						p.Product_description = value_typed
					default:
						fmt.Println("Unknown pattern")
					}

				default:
					fmt.Println("Unknown data type in JSON file")
				}
			}
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
	search_patterns SearchPatterns) *ProductList {

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

	var sp SearchPatterns

	my_patterns := make(map[string]string)
	my_patterns["Product_name"] = "Name"
	my_patterns["Product_price"] = "PSPrice"
	my_patterns["Product_description"] = "Notes"

	patterns := sp.newSearchPatterns(my_patterns)

	var file_name string = "data.json"
	jsonData, err := ioutil.ReadFile(file_name)
	if err != nil {
		sugar.Error("Couldn't read file: %s", file_name)
	}
	var data []map[string]interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		fmt.Printf("Unmarshal failed %s", err)
	}

	pl := NewProductList()
	product_list := pl.UpdateProductList(data, *patterns)

	fmt.Println(product_list)
}
