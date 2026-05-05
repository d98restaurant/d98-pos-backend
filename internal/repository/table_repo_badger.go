package repository

import (
	"encoding/json"
	"fmt"
	"time"

	"pos-backend/internal/models"

	"github.com/dgraph-io/badger/v4"
)

type TableRepository struct{}

func NewTableRepository() *TableRepository {
	return &TableRepository{}
}

func (r *TableRepository) Create(table *models.Table) error {
	table.CreatedAt = time.Now()
	table.UpdatedAt = time.Now()
	table.Status = models.TableStatusAvailable
	return SaveJSON(fmt.Sprintf("table:%d", table.TableNumber), table)
}

func (r *TableRepository) FindAll() ([]models.Table, error) {
	var tables []models.Table
	err := DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()
		prefix := []byte("table:")
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				var table models.Table
				if err := json.Unmarshal(val, &table); err != nil {
					return err
				}
				tables = append(tables, table)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return tables, err
}

func (r *TableRepository) FindByNumber(tableNumber int) (*models.Table, error) {
	var table models.Table
	err := GetJSON(fmt.Sprintf("table:%d", tableNumber), &table)
	if err == badger.ErrKeyNotFound {
		return nil, nil
	}
	return &table, err
}

func (r *TableRepository) UpdateByNumber(tableNumber int, updates map[string]interface{}) error {
	table, err := r.FindByNumber(tableNumber)
	if err != nil {
		return err
	}
	if table == nil {
		return nil
	}
	if capacity, ok := updates["capacity"]; ok {
		table.Capacity = capacity.(int)
	}
	if status, ok := updates["status"]; ok {
		table.Status = status.(models.TableStatus)
	}
	table.UpdatedAt = time.Now()
	return SaveJSON(fmt.Sprintf("table:%d", tableNumber), table)
}

func (r *TableRepository) DeleteByNumber(tableNumber int) error {
	return DeleteKey(fmt.Sprintf("table:%d", tableNumber))
}

func (r *TableRepository) IncrementRunningCount(tableNumber int, amount float64) error {
	table, err := r.FindByNumber(tableNumber)
	if err != nil || table == nil {
		return err
	}
	table.RunningOrderCount++
	table.TotalRunningAmount += amount
	return SaveJSON(fmt.Sprintf("table:%d", tableNumber), table)
}

func (r *TableRepository) DecrementRunningCount(tableNumber int, amount float64) error {
	table, err := r.FindByNumber(tableNumber)
	if err != nil || table == nil {
		return err
	}
	table.RunningOrderCount--
	table.TotalRunningAmount -= amount
	return SaveJSON(fmt.Sprintf("table:%d", tableNumber), table)
}

func (r *TableRepository) ResetRunningOrders(tableNumber int) error {
	table, err := r.FindByNumber(tableNumber)
	if err != nil || table == nil {
		return err
	}
	table.RunningOrderCount = 0
	table.TotalRunningAmount = 0
	return SaveJSON(fmt.Sprintf("table:%d", tableNumber), table)
}
