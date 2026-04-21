package main

import (
	"alkalax/item-manager/internal/storage"
	"flag"
	"fmt"
	"os"
)

// const storageFile = "items.json"
const storageFile = "items.db"

func main() {
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	createItemName := createCmd.String("name", "", "Name of the item to create")
	itemDescription := createCmd.String("description", "", "Description of the item")
	itemCount := createCmd.Int("count", 1, "Number of items to create")

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)

	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	deleteItemName := deleteCmd.String("name", "", "Name of the item to delete")

	if len(os.Args) < 2 {
		fmt.Println("expected 'create' or other subcommands")
		os.Exit(1)
	}

	itemManager := storage.NewItemManager(storage.SQLite, storageFile)

	switch os.Args[1] {
	case "create":
		err := createCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Printf("error while parsing: %v\n", err)
			os.Exit(1)
		}

		if *createItemName == "" {
			fmt.Println("error: name is required")
			createCmd.PrintDefaults()
			os.Exit(1)
		}

		fmt.Println("Creating item...")
		if err = itemManager.AddItem(*createItemName, *itemDescription, *itemCount); err != nil {
			fmt.Printf("failed to add item: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Item '%s' saved.\n", *createItemName)
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

		fmt.Printf("%-15s %-10s %-20s\n", "NAME", "COUNT", "DESCRIPTION")
		fmtString := "%-15s %-10d %-20s\n"
		for _, item := range items {
			fmt.Printf(fmtString, item.Name, item.Count, item.Description)
		}
	case "delete":
		if err := deleteCmd.Parse(os.Args[2:]); err != nil {
			fmt.Printf("error while parsing: %v\n", err)
			os.Exit(1)
		}

		if *deleteItemName == "" {
			fmt.Println("error: name is required")
			deleteCmd.PrintDefaults()
			os.Exit(1)
		}

		fmt.Println("Deleting item...")
		if err := itemManager.DeleteItem(*deleteItemName); err != nil {
			fmt.Printf("failed to delete item: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Item '%s' deleted.\n", *deleteItemName)
	default:
		fmt.Printf("unknown subcommand: %s\n", os.Args[1])
		os.Exit(1)
	}
}
