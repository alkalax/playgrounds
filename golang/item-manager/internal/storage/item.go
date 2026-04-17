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

func (im *ItemManager) AddItem(name string, count int) error {
	for _, item := range im.Items {
		if item.Name == name {
			return fmt.Errorf("item '%s' already exists", name)
		}
	}

	im.Items = append(im.Items, Item{
		Name:  name,
		Count: count,
	})

	return nil
}
