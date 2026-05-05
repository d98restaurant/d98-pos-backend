package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BusinessRepository struct {
	collection *mongo.Collection
}

func NewBusinessRepository(client *mongo.Client, dbName string) *BusinessRepository {
	return &BusinessRepository{
		collection: client.Database(dbName).Collection("business"),
	}
}

func (r *BusinessRepository) Find(ctx context.Context) (map[string]interface{}, error) {
	var business map[string]interface{}
	err := r.collection.FindOne(ctx, bson.M{}).Decode(&business)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return business, err
}

func (r *BusinessRepository) Update(ctx context.Context, business map[string]interface{}) error {
	business["updatedAt"] = time.Now()
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{},
		bson.M{"$set": business},
		&options.UpdateOptions{Upsert: func(b bool) *bool { return &b }(true)},
	)
	return err
}