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

func (im *ItemManager) loadItems() error {
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

func (im *ItemManager) saveItems() error {
	data, err := json.MarshalIndent(im.Items, "", "  ")
	if err != nil {
		return err
	}

	if err = os.WriteFile(im.storageFile, data, 0644); err != nil {
		return err
	}

	return nil
}

func (im *ItemManager) AddItem(name string, count int) error {
	if err := im.loadItems(); err != nil {
		return fmt.Errorf("failed to load items: %v", err)
	}

	for _, item := range im.Items {
		if item.Name == name {
			return fmt.Errorf("item '%s' already exists", name)
		}
	}

	im.Items = append(im.Items, Item{
		Name:  name,
		Count: count,
	})

	if err := im.saveItems(); err != nil {
		return fmt.Errorf("failed to save items: %v", err)
	}

	return nil
}

func (im *ItemManager) GetItems() ([]Item, error) {
	if err := im.loadItems(); err != nil {
		return nil, fmt.Errorf("failed to load items: %v", err)
	}

	return im.Items, nil
}
