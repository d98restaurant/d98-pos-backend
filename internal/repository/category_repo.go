package repository

import (
	"go.mongodb.org/mongo-driver/mongo/options"
	"context"
	"time"

	"pos-backend/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CategoryRepository struct {
	collection *mongo.Collection
}

func NewCategoryRepository(client *mongo.Client, dbName string) *CategoryRepository {
	return &CategoryRepository{
		collection: client.Database(dbName).Collection("categories"),
	}
}

func (r *CategoryRepository) Create(ctx context.Context, category *models.Category) error {
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()

	// Get max sort order if not set
	if category.SortOrder == 0 {
		var lastCategory models.Category
		opts := options.FindOne().SetSort(bson.M{"sortOrder": -1})
		err := r.collection.FindOne(ctx, bson.M{}, opts).Decode(&lastCategory)
		if err == nil {
			category.SortOrder = lastCategory.SortOrder + 1
		}
	}

	result, err := r.collection.InsertOne(ctx, category)
	if err != nil {
		return err
	}
	category.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *CategoryRepository) FindAll(ctx context.Context) ([]models.Category, error) {
	var categories []models.Category
	opts := options.Find().SetSort(bson.M{"sortOrder": 1})
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &categories); err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *CategoryRepository) FindByID(ctx context.Context, id string) (*models.Category, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var category models.Category
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&category)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &category, err
}

func (r *CategoryRepository) Update(ctx context.Context, id string, updates bson.M) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	updates["updatedAt"] = time.Now()
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": updates})
	return err
}

func (r *CategoryRepository) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}