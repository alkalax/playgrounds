package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type ItemManager struct {
	Items       []Item
	storageFile string
}

type Item struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func NewItemManager(storageFile string) *ItemManager {
	return &ItemManager{
		Items:       []Item{},
		storageFile: storageFile,
	}
}

func (im *ItemManager) LoadItems() error {
	_, err := os.Stat(im.storageFile)
	if errors.Is(err, os.ErrNotExist) {
		im.Items = []Item{}
		if err = os.WriteFile(im.storageFile, []byte("[]"), 0644); err != nil {
			return err
		}
		return nil
	} else if err != nil {
		return err
	}

	data, err := os.ReadFile(im.storageFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &im.Items)
	if err != nil {
		return err
	}

	return nil
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
