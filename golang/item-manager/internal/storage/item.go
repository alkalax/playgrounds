package storage

import (
	"encoding/json"
	"fmt"
	"os"
)

type Item struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func SaveItem(name string, count int) error {
	item := Item{
		Name:  name,
		Count: count,
	}

	data, err := json.MarshalIndent(item, "", "  ")
	if err != nil {
		return fmt.Errorf("could not encode json: %w", err)
	}

	filename := fmt.Sprintf("%s.json", item.Name)

	return os.WriteFile(filename, data, 0644)
}
