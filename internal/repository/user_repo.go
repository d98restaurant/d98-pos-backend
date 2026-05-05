package repository

import (
	"context"
	"time"

	"pos-backend/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(client *mongo.Client, dbName string) *UserRepository {
	return &UserRepository{
		collection: client.Database(dbName).Collection("users"),
	}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.Active = true

	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	user.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	filter := bson.M{"username": username}
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &user, err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	filter := bson.M{"email": email}
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &user, err
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user models.User
	filter := bson.M{"_id": objID}
	err = r.collection.FindOne(ctx, filter).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &user, err
}

func (r *UserRepository) FindAll(ctx context.Context) ([]models.User, error) {
	var users []models.User
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) Update(ctx context.Context, id string, updates bson.M) error {
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

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	_, err = r.collection.DeleteOne(ctx, filter)
	return err
}

func (r *UserRepository) UpdateLastLogin(ctx context.Context, id string) error {
	return r.Update(ctx, id, bson.M{"lastLogin": time.Now()})
}
