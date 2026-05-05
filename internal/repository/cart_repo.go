package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CartRepository struct {
	collection *mongo.Collection
}

func NewCartRepository(client *mongo.Client, dbName string) *CartRepository {
	return &CartRepository{
		collection: client.Database(dbName).Collection("carts"),
	}
}

func (r *CartRepository) FindByUserID(ctx context.Context, userID string) (map[string]interface{}, error) {
	var cart map[string]interface{}
	err := r.collection.FindOne(ctx, bson.M{"userId": userID}).Decode(&cart)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return cart, err
}

func (r *CartRepository) Save(ctx context.Context, userID string, cartData map[string]interface{}) error {
	cartData["userId"] = userID
	cartData["updatedAt"] = time.Now()

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"userId": userID},
		bson.M{"$set": cartData},
		&options.UpdateOptions{Upsert: func(b bool) *bool { return &b }(true)},
	)
	return err
}

func (r *CartRepository) AddItem(ctx context.Context, userID string, item map[string]interface{}) (map[string]interface{}, error) {
	cart, err := r.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if cart == nil {
		cart = map[string]interface{}{
			"userId":              userID,
			"items":               []interface{}{},
			"specialInstructions": map[string]string{},
		}
	}

	items, _ := cart["items"].([]interface{})
	itemID := item["id"].(string)

	// Check if item already exists
	found := false
	for i, existingItem := range items {
		if existingItemMap, ok := existingItem.(map[string]interface{}); ok {
			if existingItemMap["id"] == itemID {
				// Increment quantity
				qty, _ := existingItemMap["quantity"].(float64)
				existingItemMap["quantity"] = qty + 1
				items[i] = existingItemMap
				found = true
				break
			}
		}
	}

	if !found {
		if _, ok := item["quantity"]; !ok {
			item["quantity"] = 1
		}
		items = append(items, item)
	}

	cart["items"] = items
	cart["updatedAt"] = time.Now()

	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"userId": userID},
		bson.M{"$set": cart},
		&options.UpdateOptions{Upsert: func(b bool) *bool { return &b }(true)},
	)
	if err != nil {
		return nil, err
	}

	return cart, nil
}

func (r *CartRepository) UpdateItemQuantity(ctx context.Context, userID, itemID string, quantity int) (map[string]interface{}, error) {
	cart, err := r.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if cart == nil {
		return map[string]interface{}{"items": []interface{}{}}, nil
	}

	items, _ := cart["items"].([]interface{})
	for i, existingItem := range items {
		if itemMap, ok := existingItem.(map[string]interface{}); ok {
			if itemMap["id"] == itemID {
				if quantity <= 0 {
					// Remove item
					items = append(items[:i], items[i+1:]...)
				} else {
					itemMap["quantity"] = quantity
					items[i] = itemMap
				}
				break
			}
		}
	}

	cart["items"] = items
	cart["updatedAt"] = time.Now()

	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"userId": userID},
		bson.M{"$set": cart},
		&options.UpdateOptions{Upsert: true},
	)
	if err != nil {
		return nil, err
	}

	return cart, nil
}

func (r *CartRepository) RemoveItem(ctx context.Context, userID, itemID string) (map[string]interface{}, error) {
	return r.UpdateItemQuantity(ctx, userID, itemID, 0)
}

func (r *CartRepository) Clear(ctx context.Context, userID string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"userId": userID})
	return err
}