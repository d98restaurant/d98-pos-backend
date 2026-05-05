package repository

import (
	"encoding/json"
	"fmt"
	"time"

	"pos-backend/internal/models"

	"github.com/dgraph-io/badger/v4"
)

type MenuRepository struct{}

func NewMenuRepository() *MenuRepository {
	return &MenuRepository{}
}

func (r *MenuRepository) Create(item *models.MenuItem) error {
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()
	seq, _ := GetNextSequence("menu_id")
	item.ID = fmt.Sprintf("%d", seq)
	return SaveJSON("menu:"+item.ID, item)
}

func (r *MenuRepository) FindAll() ([]models.MenuItem, error) {
	var items []models.MenuItem
	err := DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()
		prefix := []byte("menu:")
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				var menuItem models.MenuItem
				if err := json.Unmarshal(val, &menuItem); err != nil {
					return err
				}
				items = append(items, menuItem)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return items, err
}

func (r *MenuRepository) Update(id string, updates map[string]interface{}) error {
	var item models.MenuItem
	err := GetJSON("menu:"+id, &item)
	if err != nil {
		return err
	}
	if name, ok := updates["name"]; ok {
		item.Name = name.(string)
	}
	if price, ok := updates["price"]; ok {
		item.Price = price.(float64)
	}
	if category, ok := updates["category"]; ok {
		item.Category = category.(string)
	}
	if available, ok := updates["available"]; ok {
		item.Available = available.(bool)
	}
	item.UpdatedAt = time.Now()
	return SaveJSON("menu:"+id, item)
}

func (r *MenuRepository) Delete(id string) error {
	return DeleteKey("menu:" + id)
}
