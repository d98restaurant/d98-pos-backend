package repository

import (
	"encoding/json"
	"time"

	"github.com/dgraph-io/badger/v4"
)

type SettingsRepository struct{}

func NewSettingsRepository() *SettingsRepository {
	return &SettingsRepository{}
}

func (r *SettingsRepository) Find() (map[string]interface{}, error) {
	var settings map[string]interface{}
	err := DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("settings:general"))
		if err == badger.ErrKeyNotFound {
			return nil
		}
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &settings)
		})
	})
	return settings, err
}

func (r *SettingsRepository) Update(settings map[string]interface{}) error {
	settings["updatedAt"] = time.Now()
	return DB.Update(func(txn *badger.Txn) error {
		data, err := json.Marshal(settings)
		if err != nil {
			return err
		}
		return txn.Set([]byte("settings:general"), data)
	})
}
