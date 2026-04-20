package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type StorageType int

const (
	JsonFile StorageType = iota
	SQLite
)

type ItemManager struct {
	Items       []Item
	StorageType StorageType
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

func (im *ItemManager) loadItemsJson() error {
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

func (im *ItemManager) saveItemsJson() error {
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
	switch im.StorageType {
	case JsonFile:
		return im.addItemJson(name, count)
	default:
		return fmt.Errorf("invalid storage type")
	}
}

func (im *ItemManager) addItemJson(name string, count int) error {
	if err := im.loadItemsJson(); err != nil {
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

	if err := im.saveItemsJson(); err != nil {
		return fmt.Errorf("failed to save items: %v", err)
	}

	return nil
}

func (im *ItemManager) GetItems() ([]Item, error) {
	switch im.StorageType {
	case JsonFile:
		if err := im.loadItemsJson(); err != nil {
			return nil, fmt.Errorf("failed to load items: %v", err)
		}

		return im.Items, nil
	default:
		return nil, fmt.Errorf("invalid storage type")
	}
}

func (im *ItemManager) DeleteItem(name string) error {
	switch im.StorageType {
	case JsonFile:
		return im.deleteItemJson(name)
	default:
		return fmt.Errorf("invalid storage type")
	}
}

func (im *ItemManager) deleteItemJson(name string) error {
	if err := im.loadItemsJson(); err != nil {
		return fmt.Errorf("failed to load items: %v", err)
	}

	for i, item := range im.Items {
		if item.Name == name {
			if i < len(im.Items)-1 {
				im.Items = append(im.Items[:i], im.Items[i+1:]...)
			} else {
				im.Items = im.Items[:i]
			}

			if err := im.saveItemsJson(); err != nil {
				return fmt.Errorf("failed to save items: %v", err)
			}

			return nil
		}
	}

	return fmt.Errorf("item not found")
}
