package repository

import (
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/badger/v4"
)

type CartRepository struct{}

func NewCartRepository() *CartRepository {
	return &CartRepository{}
}

type Cart struct {
	UserID              string                 `json:"userId"`
	Items               []map[string]interface{} `json:"items"`
	SpecialInstructions map[string]string      `json:"specialInstructions"`
}

func (r *CartRepository) FindByUserID(userID string) (*Cart, error) {
	var cart Cart
	err := GetJSON("cart:"+userID, &cart)
	if err == badger.ErrKeyNotFound {
		return nil, nil
	}
	return &cart, err
}

func (r *CartRepository) Save(userID string, cartData map[string]interface{}) error {
	cart := Cart{
		UserID:              userID,
		Items:               []map[string]interface{}{},
		SpecialInstructions: map[string]string{},
	}
	if items, ok := cartData["items"]; ok {
		cart.Items = items.([]map[string]interface{})
	}
	if instructions, ok := cartData["specialInstructions"]; ok {
		cart.SpecialInstructions = instructions.(map[string]string)
	}
	return SaveJSON("cart:"+userID, cart)
}

func (r *CartRepository) AddItem(userID string, item map[string]interface{}) (*Cart, error) {
	cart, err := r.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	if cart == nil {
		cart = &Cart{UserID: userID, Items: []map[string]interface{}{}, SpecialInstructions: map[string]string{}}
	}
	itemID := item["id"].(string)
	found := false
	for i, existingItem := range cart.Items {
		if existingItem["id"] == itemID {
			qty, _ := existingItem["quantity"].(float64)
			cart.Items[i]["quantity"] = qty + 1
			found = true
			break
		}
	}
	if !found {
		if _, ok := item["quantity"]; !ok {
			item["quantity"] = 1
		}
		cart.Items = append(cart.Items, item)
	}
	return cart, SaveJSON("cart:"+userID, cart)
}

func (r *CartRepository) UpdateItemQuantity(userID, itemID string, quantity int) (*Cart, error) {
	cart, err := r.FindByUserID(userID)
	if err != nil || cart == nil {
		return cart, err
	}
	for i, item := range cart.Items {
		if item["id"] == itemID {
			if quantity <= 0 {
				cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)
			} else {
				cart.Items[i]["quantity"] = quantity
			}
			break
		}
	}
	return cart, SaveJSON("cart:"+userID, cart)
}

func (r *CartRepository) RemoveItem(userID, itemID string) (*Cart, error) {
	return r.UpdateItemQuantity(userID, itemID, 0)
}

func (r *CartRepository) Clear(userID string) error {
	return DeleteKey("cart:" + userID)
}
