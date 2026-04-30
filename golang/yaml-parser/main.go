package main

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v4"
)

type Money struct {
	Amount   int    `yaml:"amount"`
	Currency string `yaml:"currency"`
}

type Asset struct {
	Name  string `yaml:"name"`
	Value Money  `yaml:"value"`
}

type Transation struct {
	Description string `yaml:"description"`
	Value       Money  `yaml:"value"`
}

type AccountReport struct {
	Date                 string       `yaml:"date"`
	Assets               []Asset      `yaml:"assets"`
	Transations          []Transation `yaml:"transations"`
	ExchangeRateEURToRSD float64      `yaml:"eur_to_rsd"`
}

type FinanceData struct {
	Reports []AccountReport `yaml:"reports"`
}

func main() {
	file, err := os.ReadFile("balance.yaml")
	if err != nil {
		fmt.Printf("error while reading file: %v\n", err)
		os.Exit(1)
	}

	var balance FinanceData
	if err = yaml.Unmarshal(file, &balance); err != nil {
		fmt.Printf("error while unmarshalling yaml: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%v\n", balance)
}
