package repository

import (
	"context"
	"time"

	"pos-backend/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrderRepository struct {
	collection *mongo.Collection
}

func NewOrderRepository(client *mongo.Client, dbName string) *OrderRepository {
	return &OrderRepository{
		collection: client.Database(dbName).Collection("orders"),
	}
}

func (r *OrderRepository) Create(ctx context.Context, order *models.Order) error {
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, order)
	if err != nil {
		return err
	}

	order.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *OrderRepository) FindByID(ctx context.Context, id string) (*models.Order, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var order models.Order
	filter := bson.M{"_id": objID}
	err = r.collection.FindOne(ctx, filter).Decode(&order)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &order, err
}

func (r *OrderRepository) FindAll(ctx context.Context) ([]models.Order, error) {
	var orders []models.Order
	opts := options.Find().SetSort(bson.M{"createdAt": -1})
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &orders); err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) FindByStatus(ctx context.Context, status models.OrderStatus) ([]models.Order, error) {
	var orders []models.Order
	filter := bson.M{"status": status}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &orders); err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) FindByTable(ctx context.Context, tableNumber int, activeOnly bool) ([]models.Order, error) {
	var orders []models.Order
	filter := bson.M{"tableNumber": tableNumber}
	if activeOnly {
		filter["status"] = bson.M{"$in": []models.OrderStatus{
			models.OrderStatusPending,
			models.OrderStatusAccepted,
			models.OrderStatusPreparing,
		}}
	}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &orders); err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) Update(ctx context.Context, id string, updates bson.M) error {
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

func (r *OrderRepository) UpdateStatus(ctx context.Context, id string, status models.OrderStatus) error {
	return r.Update(ctx, id, bson.M{"status": status})
}

func (r *OrderRepository) GetNextOrderNumber(ctx context.Context) (int, error) {
	// This would typically use a counter collection
	// For simplicity, we'll return a timestamp-based number
	return int(time.Now().UnixNano() % 1000000), nil
}

func (r *OrderRepository) FindCreditCustomers(ctx context.Context) ([]map[string]interface{}, error) {
	pipeline := mongo.Pipeline{
		{{"$match", bson.M{
			"payment.method": "credit",
			"payment.status": "credit_due",
			"status":         bson.M{"$ne": "cancelled"},
		}}},
		{{"$group", bson.M{
			"_id": bson.M{
				"customerId":   "$payment.customerId",
				"customerName": "$payment.customerName",
				"customerPhone":"$payment.customerPhone",
			},
			"totalDue": bson.M{"$sum": "$total"},
			"orders":   bson.M{"$push": "$$ROOT"},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []map[string]interface{}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}