package main

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v4"
)

type Balance struct {
	Date             string            `yaml:"date"`
	CurrencyBalances []CurrencyBalance `yaml:"balances"`
}

type CurrencyBalance struct {
	Currency string `yaml:"currency"`
	Total    int    `yaml:"total"`
}

func main() {
	file, err := os.ReadFile("balance.yaml")
	if err != nil {
		fmt.Printf("error while reading file: %v\n", err)
		os.Exit(1)
	}

	var balance Balance
	if err = yaml.Unmarshal(file, &balance); err != nil {
		fmt.Printf("error while unmarshalling yaml: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%v\n", balance)
}
