package json

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	st "alkalax/item-manager/internal/storage"
)

type JsonFileStore struct {
	storageFile string
	items       []st.Item
}

func NewStore(storageFile string) st.StorageManager {
	return &JsonFileStore{
		storageFile: storageFile,
		items:       []st.Item{},
	}
}

func (sm *JsonFileStore) loadItems() error {
	_, err := os.Stat(sm.storageFile)
	if errors.Is(err, os.ErrNotExist) {
		sm.items = []st.Item{}
		if err = os.WriteFile(sm.storageFile, []byte("[]"), 0644); err != nil {
			return err
		}
		return nil
	} else if err != nil {
		return err
	}

	data, err := os.ReadFile(sm.storageFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &sm.items)
	if err != nil {
		return err
	}

	return nil
}

func (sm *JsonFileStore) saveItems() error {
	data, err := json.MarshalIndent(sm.items, "", "  ")
	if err != nil {
		return err
	}

	if err = os.WriteFile(sm.storageFile, data, 0644); err != nil {
		return err
	}

	return nil
}

func (sm *JsonFileStore) AddItem(name string, description string, count int) error {
	if err := sm.loadItems(); err != nil {
		return fmt.Errorf("failed to load items: %v", err)
	}

	for _, item := range sm.items {
		if item.Name == name {
			return fmt.Errorf("item '%s' already exists", name)
		}
	}

	sm.items = append(sm.items, st.Item{
		Name:        name,
		Description: description,
		Count:       count,
	})

	if err := sm.saveItems(); err != nil {
		return fmt.Errorf("failed to save items: %v", err)
	}

	return nil
}

func (sm *JsonFileStore) GetItems() ([]st.Item, error) {
	if err := sm.loadItems(); err != nil {
		return nil, fmt.Errorf("failed to load items: %v", err)
	}

	return sm.items, nil
}

func (sm *JsonFileStore) DeleteItem(name string) error {
	if err := sm.loadItems(); err != nil {
		return fmt.Errorf("failed to load items: %v", err)
	}

	for i, item := range sm.items {
		if item.Name == name {
			if i < len(sm.items)-1 {
				sm.items = append(sm.items[:i], sm.items[i+1:]...)
			} else {
				sm.items = sm.items[:i]
			}

			if err := sm.saveItems(); err != nil {
				return fmt.Errorf("failed to save items: %v", err)
			}

			return nil
		}
	}

	return fmt.Errorf("item not found")
}
