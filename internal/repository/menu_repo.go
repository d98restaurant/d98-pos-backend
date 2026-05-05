package repository

import (
	"go.mongodb.org/mongo-driver/mongo/options"
	"context"
	"time"

	"pos-backend/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MenuRepository struct {
	collection *mongo.Collection
}

func NewMenuRepository(client *mongo.Client, dbName string) *MenuRepository {
	return &MenuRepository{
		collection: client.Database(dbName).Collection("menu"),
	}
}

func (r *MenuRepository) Create(ctx context.Context, item *models.MenuItem) error {
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, item)
	if err != nil {
		return err
	}

	item.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *MenuRepository) FindAll(ctx context.Context) ([]models.MenuItem, error) {
	var items []models.MenuItem
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *MenuRepository) FindByID(ctx context.Context, id string) (*models.MenuItem, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var item models.MenuItem
	filter := bson.M{"_id": objID}
	err = r.collection.FindOne(ctx, filter).Decode(&item)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &item, err
}

func (r *MenuRepository) Update(ctx context.Context, id string, updates bson.M) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	updates["updatedAt"] = time.Now()
	filter := bson.M{"_id": objID}
	update := bson.M{"$set": updates}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *MenuRepository) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	_, err = r.collection.DeleteOne(ctx, filter)
	return err
}