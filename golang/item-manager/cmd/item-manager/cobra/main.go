package main

import (
	"alkalax/item-manager/cmd/item-manager/cmd"
	"alkalax/item-manager/internal/storage/json"
)

func main() {
	cmd.Execute(json.NewStore("items.json"))
}
