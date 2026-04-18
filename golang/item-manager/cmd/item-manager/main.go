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

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("expected 'create' or other subcommands")
		os.Exit(1)
	}

	itemManager := storage.NewItemManager(storageFile)

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
		if err = itemManager.AddItem(*itemName, *itemCount); err != nil {
			fmt.Printf("failed to add item: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Item '%s' saved.\n", *itemName)
	case "list":
		if err := listCmd.Parse(os.Args[2:]); err != nil {
			fmt.Printf("error while parsing: %v\n", err)
			os.Exit(1)
		}

		items, err := itemManager.GetItems()
		if err != nil {
			fmt.Printf("failed to get items: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("%-15s %-10s\n", "NAME", "COUNT")
		fmtString := "%-15s %-10d\n"
		for _, item := range items {
			fmt.Printf(fmtString, item.Name, item.Count)
		}
	default:
		fmt.Printf("unknown subcommand: %s\n", os.Args[1])
		os.Exit(1)
	}
}
