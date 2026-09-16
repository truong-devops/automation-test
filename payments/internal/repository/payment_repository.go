package repository

import (
	"context"
	"errors"

	"automation-test/payments/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PaymentRepository interface {
	Create(context.Context, *domain.Payment) error
	List(context.Context) ([]domain.Payment, error)
	GetByID(context.Context, string) (*domain.Payment, error)
	GetByOrderID(context.Context, string) (*domain.Payment, error)
	Update(context.Context, *domain.Payment) error
	Delete(context.Context, string) error
}

type MongoPaymentRepository struct {
	collection *mongo.Collection
}

func NewMongoPaymentRepository(db *mongo.Database) *MongoPaymentRepository {
	return &MongoPaymentRepository{collection: db.Collection("payments")}
}

func (r *MongoPaymentRepository) EnsureIndexes(ctx context.Context) error {
	_, err := r.collection.Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.D{{Key: "order_id", Value: 1}}})
	return err
}

func (r *MongoPaymentRepository) Create(ctx context.Context, payment *domain.Payment) error {
	doc := bson.M{
		"order_id": payment.OrderID, "amount": payment.Amount, "method": payment.Method,
		"status": payment.Status, "created_at": payment.CreatedAt, "updated_at": payment.UpdatedAt,
	}
	result, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return err
	}
	payment.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func (r *MongoPaymentRepository) List(ctx context.Context) ([]domain.Payment, error) {
	cursor, err := r.collection.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var docs []struct {
		ID             primitive.ObjectID `bson:"_id"`
		domain.Payment `bson:",inline"`
	}
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}
	payments := make([]domain.Payment, 0, len(docs))
	for _, doc := range docs {
		doc.Payment.ID = doc.ID.Hex()
		payments = append(payments, doc.Payment)
	}
	return payments, nil
}

func (r *MongoPaymentRepository) GetByID(ctx context.Context, id string) (*domain.Payment, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrInvalidID
	}
	return r.findOne(ctx, bson.M{"_id": objectID}, nil)
}

func (r *MongoPaymentRepository) GetByOrderID(ctx context.Context, orderID string) (*domain.Payment, error) {
	return r.findOne(ctx, bson.M{"order_id": orderID}, options.FindOne().SetSort(bson.D{{Key: "created_at", Value: -1}}))
}

func (r *MongoPaymentRepository) findOne(ctx context.Context, filter any, opts *options.FindOneOptions) (*domain.Payment, error) {
	var doc struct {
		ID             primitive.ObjectID `bson:"_id"`
		domain.Payment `bson:",inline"`
	}
	var result *mongo.SingleResult
	if opts == nil {
		result = r.collection.FindOne(ctx, filter)
	} else {
		result = r.collection.FindOne(ctx, filter, opts)
	}
	if err := result.Decode(&doc); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	doc.Payment.ID = doc.ID.Hex()
	return &doc.Payment, nil
}

func (r *MongoPaymentRepository) Update(ctx context.Context, payment *domain.Payment) error {
	objectID, err := primitive.ObjectIDFromHex(payment.ID)
	if err != nil {
		return domain.ErrInvalidID
	}
	update := bson.M{"$set": bson.M{
		"order_id": payment.OrderID, "amount": payment.Amount, "method": payment.Method,
		"status": payment.Status, "updated_at": payment.UpdatedAt,
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

func (r *MongoPaymentRepository) Delete(ctx context.Context, id string) error {
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
