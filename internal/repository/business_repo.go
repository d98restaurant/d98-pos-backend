package repository

import (
	"encoding/json"
	"time"

	"github.com/dgraph-io/badger/v4"
)

type BusinessRepository struct{}

func NewBusinessRepository() *BusinessRepository {
	return &BusinessRepository{}
}

func (r *BusinessRepository) Find() (map[string]interface{}, error) {
	var business map[string]interface{}
	err := DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("business:settings"))
		if err == badger.ErrKeyNotFound {
			return nil
		}
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &business)
		})
	})
	return business, err
}

func (r *BusinessRepository) Update(business map[string]interface{}) error {
	business["updatedAt"] = time.Now()
	return DB.Update(func(txn *badger.Txn) error {
		data, err := json.Marshal(business)
		if err != nil {
			return err
		}
		return txn.Set([]byte("business:settings"), data)
	})
}
