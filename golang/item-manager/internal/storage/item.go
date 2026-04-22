package storage

type StorageType int

const (
	JsonFile StorageType = iota
	SQLite
)

type StorageManager interface {
	AddItem(name string, description string, count int) error
	DeleteItem(name string) error
	GetItems() ([]Item, error)
}

type Item struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Count       int    `json:"count"`
}
