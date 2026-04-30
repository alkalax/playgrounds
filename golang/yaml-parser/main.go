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

func (report *AccountReport) GetTotalAssetsCurrency(currency string) int {
	total := 0
	for _, asset := range report.Assets {
		if asset.Value.Currency == currency {
			total += asset.Value.Amount
		}
	}

	return total
}

func ConvertEURToRSD(amount int, exchangeRate float64) int {
	return int(float64(amount) * exchangeRate)
}

func ConvertRSDToEUR(amount int, exchangeRate float64) int {
	return int(float64(amount) / exchangeRate)
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

	assetsRSD := balance.Reports[0].GetTotalAssetsCurrency("RSD")
	fmt.Printf("Assets RSD: %d RSD\n", assetsRSD)

	assetsEUR := balance.Reports[0].GetTotalAssetsCurrency("EUR")
	fmt.Printf("Assets EUR: %d EUR\n", assetsEUR)

	totalRSD := assetsRSD + ConvertEURToRSD(assetsEUR, balance.Reports[0].ExchangeRateEURToRSD)
	totalEUR := ConvertRSDToEUR(totalRSD, balance.Reports[0].ExchangeRateEURToRSD)
	fmt.Printf("Total: %d RSD (%d EUR)\n", totalRSD, totalEUR)

}
