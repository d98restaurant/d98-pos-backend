package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SettingsRepository struct {
	collection *mongo.Collection
}

func NewSettingsRepository(client *mongo.Client, dbName string) *SettingsRepository {
	return &SettingsRepository{
		collection: client.Database(dbName).Collection("settings"),
	}
}

func (r *SettingsRepository) Find(ctx context.Context) (map[string]interface{}, error) {
	var settings map[string]interface{}
	err := r.collection.FindOne(ctx, bson.M{}).Decode(&settings)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return settings, err
}

func (r *SettingsRepository) Update(ctx context.Context, settings map[string]interface{}) error {
	settings["updatedAt"] = time.Now()
	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{},
		bson.M{"$set": settings},
		opts,
	)
	return err
}
