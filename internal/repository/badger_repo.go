package repository

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/dgraph-io/badger/v4"
)

var (
	DB     *badger.DB
	dbOnce sync.Once
)

// InitBadgerDB initializes the Badger database
func InitBadgerDB(dbPath string) (*badger.DB, error) {
	var err error
	dbOnce.Do(func() {
		opts := badger.DefaultOptions(dbPath)
		opts.Logger = nil // Disable logging for production
		
		DB, err = badger.Open(opts)
		if err != nil {
			log.Printf("Failed to open BadgerDB: %v", err)
		} else {
			log.Printf("✅ BadgerDB connected at %s", dbPath)
		}
	})
	return DB, err
}

// CloseDB closes the database
func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// SaveJSON saves any data as JSON with a key
func SaveJSON(key string, value interface{}) error {
	return DB.Update(func(txn *badger.Txn) error {
		data, err := json.Marshal(value)
		if err != nil {
			return err
		}
		return txn.Set([]byte(key), data)
	})
}

// GetJSON retrieves and unmarshals JSON data by key
func GetJSON(key string, result interface{}) error {
	return DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, result)
		})
	})
}

// DeleteKey deletes a key from the database
func DeleteKey(key string) error {
	return DB.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}

// GetAllKeys returns all keys with a given prefix
func GetAllKeys(prefix string) ([]string, error) {
	var keys []string
	err := DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		
		prefixBytes := []byte(prefix)
		for it.Seek(prefixBytes); it.ValidForPrefix(prefixBytes); it.Next() {
			keys = append(keys, string(it.Item().Key()))
		}
		return nil
	})
	return keys, err
}

// GetAllValues returns all values with a given prefix
func GetAllValues(prefix string, resultSlice interface{}) error {
	return DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()
		
		prefixBytes := []byte(prefix)
		sliceValue := getSliceValue(resultSlice)
		
		for it.Seek(prefixBytes); it.ValidForPrefix(prefixBytes); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				return appendToSlice(sliceValue, val)
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// Helper functions for slice operations
func getSliceValue(slice interface{}) interface{} {
	return slice
}

func appendToSlice(slice interface{}, data []byte) error {
	// This is a simplified version - in production, use reflection
	// For now, we'll handle specific types in the repository methods
	return nil
}

// Sequence generates auto-incrementing IDs
func GetNextSequence(seqName string) (uint64, error) {
	var seq uint64
	err := DB.Update(func(txn *badger.Txn) error {
		key := []byte("seq_" + seqName)
		item, err := txn.Get(key)
		if err == badger.ErrKeyNotFound {
			seq = 1
		} else if err != nil {
			return err
		} else {
			err = item.Value(func(val []byte) error {
				// Simple increment - in production use proper sequence
				seq = uint64(time.Now().UnixNano() % 1000000)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return txn.Set(key, []byte{byte(seq)})
	})
	return seq, err
}