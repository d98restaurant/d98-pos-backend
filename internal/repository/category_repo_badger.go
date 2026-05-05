package repository

import (
	"encoding/json"
	"fmt"
	"time"

	"pos-backend/internal/models"

	"github.com/dgraph-io/badger/v4"
)

type CategoryRepository struct{}

func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{}
}

func (r *CategoryRepository) Create(category *models.Category) error {
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()
	seq, _ := GetNextSequence("category_id")
	category.ID = fmt.Sprintf("%d", seq)
	return SaveJSON("category:"+category.ID, category)
}

func (r *CategoryRepository) FindAll() ([]models.Category, error) {
	var categories []models.Category
	err := DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()
		prefix := []byte("category:")
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				var cat models.Category
				if err := json.Unmarshal(val, &cat); err != nil {
					return err
				}
				categories = append(categories, cat)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return categories, err
}

func (r *CategoryRepository) FindByID(id string) (*models.Category, error) {
	var cat models.Category
	err := GetJSON("category:"+id, &cat)
	if err == badger.ErrKeyNotFound {
		return nil, nil
	}
	return &cat, err
}

func (r *CategoryRepository) Update(id string, updates map[string]interface{}) error {
	var cat models.Category
	err := GetJSON("category:"+id, &cat)
	if err != nil {
		return err
	}
	if name, ok := updates["name"]; ok {
		cat.Name = name.(string)
	}
	if sortOrder, ok := updates["sortOrder"]; ok {
		cat.SortOrder = sortOrder.(int)
	}
	cat.UpdatedAt = time.Now()
	return SaveJSON("category:"+id, cat)
}

func (r *CategoryRepository) Delete(id string) error {
	return DeleteKey("category:" + id)
}
