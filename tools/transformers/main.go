package main

import (
	"encoding/json"
	"fmt"
	"os"

	"go.uber.org/zap"
)

type Logger interface {
	Debug(args ...interface{})
	Error(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
}

type ZapLogger struct {
	logger *zap.SugaredLogger
}

func NewLogger() *ZapLogger {
	new_logger, err := zap.NewProduction()
	if err != nil {
		fmt.Println("Couldn't initialize zap logging")
	}
	sugar := new_logger.Sugar()
	sugar.Debug("Logger initialized")

	return &ZapLogger{logger: sugar}
}

type Patterns struct {
	NamePattern        string `json:"name_pattern"`
	PricePattern       string `json:"price_pattern"`
	DescriptionPattern string `json:"description_pattern"`
}

type Product struct {
	Product_name        string
	Product_price       float64
	Product_description string
	logger              Logger
}

func NewProduct(logger Logger) *Product {
	p := &Product{
		Product_name:        "name",
		Product_price:       000000,
		Product_description: "description",
		logger:              logger,
	}

	return p
}

func (p *Product) FilterJsonToProductByPatterns(json_data map[string]interface{},
	patterns Patterns) *Product {

	for key, value := range json_data {
		switch key {
		case patterns.DescriptionPattern:
			if v, ok := value.(string); !ok {
				p.logger.Warn("Invalid description type. Must be strong. Maybe missing quotes?")
			} else {
				p.Product_description = v
			}
		case patterns.NamePattern:
			if v, ok := value.(string); !ok {
				p.logger.Warn("Invalid name type. Must be string. Maybe missing quotes?")
			} else {
				p.Product_name = v
			}
		case patterns.PricePattern:
			if v, ok := value.(float64); !ok {
				p.logger.Warn("Invalid price type. Must be float. Maybe ambitious quotes?")
			} else {
				p.Product_price = v
			}
		default:
			p.logger.Debug("Key not matching any of the patterns")

		}
	}

	return p
}

type ProductList struct {
	products []Product
	logger   Logger
}

func NewProductList(logger Logger) *ProductList {
	logger.Debug("Init new product list")
	var p []Product

	pl := &ProductList{
		products: p,
		logger:   logger,
	}

	return pl
}

func (pl *ProductList) UpdateProductList(json_data []map[string]interface{},
	search_patterns Patterns) *ProductList {

	for _, product := range json_data {
		p := NewProduct(pl.logger)
		product := p.FilterJsonToProductByPatterns(product, search_patterns)

		pl.products = append(pl.products, *product)

	}

	return pl
}

func main() {

	logger := NewLogger()
	log := logger.logger
	log.Info("logger started")

	// Load patterns
	var patterns_struct Patterns
	var patterns_file string = "patterns.json"

	patterns, err := os.ReadFile(patterns_file)
	if err != nil {
		log.Error("Error reading pattern file")
	}

	if err := json.Unmarshal(patterns, &patterns_struct); err != nil {
		log.Error("Error unmarshal pattern file")
	}

	// Load data
	var file_name string = "data.json"
	jsonData, err := os.ReadFile(file_name)
	if err != nil {
		log.Error("Couldn't read file: %s", file_name)
	}
	var data []map[string]interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		fmt.Printf("Unmarshal failed %s", err)
	}

	// Create convention product list
	pl := NewProductList(log)
	product_list := pl.UpdateProductList(data, patterns_struct)

	fmt.Println(product_list)
}
