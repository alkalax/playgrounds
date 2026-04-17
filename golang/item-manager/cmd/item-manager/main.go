package main

import (
	"alkalax/item-manager/internal/storage"
	"flag"
	"fmt"
	"os"
)

const storageFile = "items.json"

func main() {
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	itemName := createCmd.String("name", "", "Name of the item to create")
	itemCount := createCmd.Int("count", 1, "Number of items to create")

	if len(os.Args) < 2 {
		fmt.Println("expected 'create' or other subcommands")
		os.Exit(1)
	}

	itemManager := storage.NewItemManager(storageFile)
	if err := itemManager.LoadItems(); err != nil {
		fmt.Printf("failed to load items: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(len(itemManager.Items))
	for _, item := range itemManager.Items {
		fmt.Println(item)
	}

	switch os.Args[1] {
	case "create":
		err := createCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Printf("error while parsing: %v\n", err)
			os.Exit(1)
		}

		if *itemName == "" {
			fmt.Println("error: name is required")
			createCmd.PrintDefaults()
			os.Exit(1)
		}

		fmt.Println("Creating item...")
		err = storage.SaveItem(*itemName, *itemCount)
		if err != nil {
			fmt.Printf("failed to save item: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Item '%s' saved.\n", *itemName)
	default:
		fmt.Printf("unknown subcommand: %s\n", os.Args[1])
		os.Exit(1)
	}
}
