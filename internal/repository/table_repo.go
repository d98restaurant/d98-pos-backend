package repository

import (
	"context"
	"time"

	"pos-backend/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TableRepository struct {
	collection *mongo.Collection
}

func NewTableRepository(client *mongo.Client, dbName string) *TableRepository {
	return &TableRepository{
		collection: client.Database(dbName).Collection("tables"),
	}
}

func (r *TableRepository) Create(ctx context.Context, table *models.Table) error {
	table.CreatedAt = time.Now()
	table.UpdatedAt = time.Now()
	table.Status = models.TableStatusAvailable
	table.RunningOrderCount = 0
	table.TotalRunningAmount = 0

	result, err := r.collection.InsertOne(ctx, table)
	if err != nil {
		return err
	}
	table.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *TableRepository) FindAll(ctx context.Context) ([]models.Table, error) {
	var tables []models.Table
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &tables); err != nil {
		return nil, err
	}
	return tables, nil
}

func (r *TableRepository) FindByNumber(ctx context.Context, tableNumber int) (*models.Table, error) {
	var table models.Table
	err := r.collection.FindOne(ctx, bson.M{"tableNumber": tableNumber}).Decode(&table)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &table, err
}

func (r *TableRepository) UpdateByNumber(ctx context.Context, tableNumber int, updates bson.M) error {
	updates["updatedAt"] = time.Now()
	_, err := r.collection.UpdateOne(ctx, bson.M{"tableNumber": tableNumber}, bson.M{"$set": updates})
	return err
}

func (r *TableRepository) DeleteByNumber(ctx context.Context, tableNumber int) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"tableNumber": tableNumber})
	return err
}

func (r *TableRepository) IncrementRunningCount(ctx context.Context, tableNumber int, amount float64) error {
	_, err := r.collection.UpdateOne(ctx,
		bson.M{"tableNumber": tableNumber},
		bson.M{"$inc": bson.M{"runningOrderCount": 1, "totalRunningAmount": amount}},
	)
	return err
}

func (r *TableRepository) DecrementRunningCount(ctx context.Context, tableNumber int, amount float64) error {
	_, err := r.collection.UpdateOne(ctx,
		bson.M{"tableNumber": tableNumber},
		bson.M{"$inc": bson.M{"runningOrderCount": -1, "totalRunningAmount": -amount}},
	)
	return err
}

func (r *TableRepository) ResetRunningOrders(ctx context.Context, tableNumber int) error {
	_, err := r.collection.UpdateOne(ctx,
		bson.M{"tableNumber": tableNumber},
		bson.M{"$set": bson.M{"runningOrderCount": 0, "totalRunningAmount": 0, "currentSessionId": ""}},
	)
	return err
}
