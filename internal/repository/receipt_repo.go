package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReceiptRepository struct {
	collection *mongo.Collection
}

func NewReceiptRepository(client *mongo.Client, dbName string) *ReceiptRepository {
	return &ReceiptRepository{
		collection: client.Database(dbName).Collection("receipts"),
	}
}

func (r *ReceiptRepository) FindByID(ctx context.Context, id string) (map[string]interface{}, error) {
	var receipt map[string]interface{}
	err := r.collection.FindOne(ctx, bson.M{"receiptId": id}).Decode(&receipt)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return receipt, err
}

func (r *ReceiptRepository) Create(ctx context.Context, receipt map[string]interface{}) error {
	receipt["createdAt"] = time.Now()
	_, err := r.collection.InsertOne(ctx, receipt)
	return err
}