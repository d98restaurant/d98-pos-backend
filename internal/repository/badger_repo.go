package repository

import (
	"encoding/json"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/dgraph-io/badger/v4"
)

var (
	DB     *badger.DB
	dbOnce sync.Once
)

func InitBadgerDB(dbPath string) (*badger.DB, error) {
	var err error
	dbOnce.Do(func() {
		opts := badger.DefaultOptions(dbPath)
		opts.Logger = nil
		
		DB, err = badger.Open(opts)
		if err != nil {
			log.Printf("Failed to open BadgerDB: %v", err)
		} else {
			log.Printf("✅ BadgerDB connected at %s", dbPath)
		}
	})
	return DB, err
}

func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

func SaveJSON(key string, value interface{}) error {
	return DB.Update(func(txn *badger.Txn) error {
		data, err := json.Marshal(value)
		if err != nil {
			return err
		}
		return txn.Set([]byte(key), data)
	})
}

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

func DeleteKey(key string) error {
	return DB.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}

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

func GetNextSequence(seqName string) (uint64, error) {
	var seq uint64
	err := DB.Update(func(txn *badger.Txn) error {
		key := []byte("seq_" + seqName)
		item, err := txn.Get(key)
		if err == badger.ErrKeyNotFound {
			// Use timestamp-based ID to ensure uniqueness
			seq = uint64(time.Now().UnixNano())
		} else if err != nil {
			return err
		} else {
			err = item.Value(func(val []byte) error {
				seq, err = strconv.ParseUint(string(val), 10, 64)
				if err != nil {
					seq = uint64(time.Now().UnixNano())
				}
				seq++
				return nil
			})
			if err != nil {
				return err
			}
		}
		return txn.Set(key, []byte(strconv.FormatUint(seq, 10)))
	})
	return seq, err
}
