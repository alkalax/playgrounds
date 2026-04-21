package storage

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
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
	Name        string `json:"name"`
	Description string `json:"description"`
	Count       int    `json:"count"`
}

func NewItemManager(storageType StorageType, storageFile string) *ItemManager {
	return &ItemManager{
		Items:       []Item{},
		StorageType: storageType,
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

func (im *ItemManager) loadItemsSQLite() error {
	db, err := sql.Open("sqlite", im.storageFile)
	if err != nil {
		return fmt.Errorf("failed to open database file: %v", err)
	}
	defer db.Close()

	if _, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS Items (
		Name TEXT PRIMARY KEY,
		Description TEXT,
		Count INTEGER
	)
	`); err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}

	rows, err := db.Query("SELECT Name, Description, Count FROM Items")
	if err != nil {
		return fmt.Errorf("failed to get items from table: %v", err)
	}

	for rows.Next() {
		var item Item
		if err = rows.Scan(&item.Name, &item.Description, &item.Count); err != nil {
			return fmt.Errorf("failed to scan row: %v", err)
		}
		im.Items = append(im.Items, item)
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

func (im *ItemManager) AddItem(name string, description string, count int) error {
	switch im.StorageType {
	case JsonFile:
		return im.addItemJson(name, description, count)
	case SQLite:
		return im.addItemSQLite(name, description, count)
	default:
		return fmt.Errorf("invalid storage type")
	}
}

func (im *ItemManager) addItemJson(name string, description string, count int) error {
	if err := im.loadItemsJson(); err != nil {
		return fmt.Errorf("failed to load items: %v", err)
	}

	for _, item := range im.Items {
		if item.Name == name {
			return fmt.Errorf("item '%s' already exists", name)
		}
	}

	im.Items = append(im.Items, Item{
		Name:        name,
		Description: description,
		Count:       count,
	})

	if err := im.saveItemsJson(); err != nil {
		return fmt.Errorf("failed to save items: %v", err)
	}

	return nil
}

func (im *ItemManager) addItemSQLite(name string, description string, count int) error {
	if err := im.loadItemsSQLite(); err != nil {
		return err
	}

	for _, item := range im.Items {
		if item.Name == name {
			return fmt.Errorf("item '%s' already exists", name)
		}
	}

	db, err := sql.Open("sqlite", im.storageFile)
	if err != nil {
		return fmt.Errorf("failed to open database file: %v", err)
	}
	defer db.Close()

	if _, err := db.Exec("INSERT INTO Items (Name, Description, Count) VALUES (?, ?, ?)", name, description, count); err != nil {
		return fmt.Errorf("failed to create item: %v", err)
	}

	im.Items = append(im.Items, Item{
		Name:        name,
		Description: description,
		Count:       count,
	})

	return nil
}

func (im *ItemManager) GetItems() ([]Item, error) {
	switch im.StorageType {
	case JsonFile:
		if err := im.loadItemsJson(); err != nil {
			return nil, fmt.Errorf("failed to load items: %v", err)
		}

	case SQLite:
		if err := im.loadItemsSQLite(); err != nil {
			return nil, fmt.Errorf("failed to load items: %v", err)
		}
	default:
		return nil, fmt.Errorf("invalid storage type")
	}

	return im.Items, nil
}

func (im *ItemManager) DeleteItem(name string) error {
	switch im.StorageType {
	case JsonFile:
		return im.deleteItemJson(name)
	case SQLite:
		return im.deleteItemSQLite(name)
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

func (im *ItemManager) deleteItemSQLite(name string) error {
	if err := im.loadItemsSQLite(); err != nil {
		return err
	}

	db, err := sql.Open("sqlite", im.storageFile)
	if err != nil {
		return fmt.Errorf("failed to open database file: %v", err)
	}
	defer db.Close()

	if _, err := db.Exec("DELETE FROM Items WHERE Name = ?", name); err != nil {
		return fmt.Errorf("failed to delete item: %v", err)
	}

	for i, item := range im.Items {
		if item.Name == name {
			if i < len(im.Items)-1 {
				im.Items = append(im.Items[:i], im.Items[i+1:]...)
			} else {
				im.Items = im.Items[:i]
			}

			return nil
		}
	}

	return fmt.Errorf("item not found")
}
