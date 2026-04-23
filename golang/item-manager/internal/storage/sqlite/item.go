package storage

import (
	"database/sql"
	"fmt"

	st "alkalax/item-manager/internal/storage"
	_ "modernc.org/sqlite"
)

type SQLiteStore struct {
	storageFile string
	items       []st.Item
}

func NewStore(storageFile string) *SQLiteStore {
	return &SQLiteStore{
		storageFile: storageFile,
		items:       []st.Item{},
	}
}

func (sm *SQLiteStore) loadItems() error {
	db, err := sql.Open("sqlite", sm.storageFile)
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
		var item st.Item
		if err = rows.Scan(&item.Name, &item.Description, &item.Count); err != nil {
			return fmt.Errorf("failed to scan row: %v", err)
		}
		sm.items = append(sm.items, item)
	}

	return nil
}

func (sm *SQLiteStore) AddItem(name string, description string, count int) error {
	if err := sm.loadItems(); err != nil {
		return err
	}

	for _, item := range sm.items {
		if item.Name == name {
			return fmt.Errorf("item '%s' already exists", name)
		}
	}

	db, err := sql.Open("sqlite", sm.storageFile)
	if err != nil {
		return fmt.Errorf("failed to open database file: %v", err)
	}
	defer db.Close()

	if _, err := db.Exec("INSERT INTO Items (Name, Description, Count) VALUES (?, ?, ?)", name, description, count); err != nil {
		return fmt.Errorf("failed to create item: %v", err)
	}

	sm.items = append(sm.items, st.Item{
		Name:        name,
		Description: description,
		Count:       count,
	})

	return nil
}

func (sm *SQLiteStore) GetItems() ([]st.Item, error) {
	if err := sm.loadItems(); err != nil {
		return nil, fmt.Errorf("failed to load items: %v", err)
	}

	return sm.items, nil
}

func (sm *SQLiteStore) DeleteItem(name string) error {
	if err := sm.loadItems(); err != nil {
		return err
	}

	db, err := sql.Open("sqlite", sm.storageFile)
	if err != nil {
		return fmt.Errorf("failed to open database file: %v", err)
	}
	defer db.Close()

	if _, err := db.Exec("DELETE FROM Items WHERE Name = ?", name); err != nil {
		return fmt.Errorf("failed to delete item: %v", err)
	}

	for i, item := range sm.items {
		if item.Name == name {
			if i < len(sm.items)-1 {
				sm.items = append(sm.items[:i], sm.items[i+1:]...)
			} else {
				sm.items = sm.items[:i]
			}

			return nil
		}
	}

	return fmt.Errorf("item not found")
}
