package repository

import (
	"context"
	"errors"

	"automation-test/orders/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrderRepository interface {
	Create(context.Context, *domain.Order) error
	List(context.Context) ([]domain.Order, error)
	GetByID(context.Context, string) (*domain.Order, error)
	Update(context.Context, *domain.Order) error
	Delete(context.Context, string) error
}

type MongoOrderRepository struct {
	collection *mongo.Collection
}

func NewMongoOrderRepository(db *mongo.Database) *MongoOrderRepository {
	return &MongoOrderRepository{collection: db.Collection("orders")}
}

func (r *MongoOrderRepository) Create(ctx context.Context, order *domain.Order) error {
	doc := bson.M{
		"customer_name": order.CustomerName, "item": order.Item,
		"quantity": order.Quantity, "total_amount": order.TotalAmount,
		"status": order.Status, "created_at": order.CreatedAt, "updated_at": order.UpdatedAt,
	}
	result, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return err
	}
	order.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func (r *MongoOrderRepository) List(ctx context.Context) ([]domain.Order, error) {
	cursor, err := r.collection.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var docs []struct {
		ID           primitive.ObjectID `bson:"_id"`
		domain.Order `bson:",inline"`
	}
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}
	orders := make([]domain.Order, 0, len(docs))
	for _, doc := range docs {
		doc.Order.ID = doc.ID.Hex()
		orders = append(orders, doc.Order)
	}
	return orders, nil
}

func (r *MongoOrderRepository) GetByID(ctx context.Context, id string) (*domain.Order, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrInvalidID
	}
	var doc struct {
		ID           primitive.ObjectID `bson:"_id"`
		domain.Order `bson:",inline"`
	}
	if err := r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&doc); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	doc.Order.ID = doc.ID.Hex()
	return &doc.Order, nil
}

func (r *MongoOrderRepository) Update(ctx context.Context, order *domain.Order) error {
	objectID, err := primitive.ObjectIDFromHex(order.ID)
	if err != nil {
		return domain.ErrInvalidID
	}
	update := bson.M{"$set": bson.M{
		"customer_name": order.CustomerName, "item": order.Item,
		"quantity": order.Quantity, "total_amount": order.TotalAmount,
		"status": order.Status, "updated_at": order.UpdatedAt,
	}}
	result, err := r.collection.UpdateByID(ctx, objectID, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *MongoOrderRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrInvalidID
	}
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return domain.ErrNotFound
	}
	return nil
}
