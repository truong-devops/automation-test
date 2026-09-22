package automationtests

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	mongoURI         = "mongodb://localhost:27018"
	ordersDatabase   = "orders_db"
	ordersCollection = "orders"
)

func (w *World) connectMongo(ctx context.Context) error {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}

	w.mongoClient = client

	fmt.Println("MongoDB connected!")

	return nil
}

func (w *World) resetMongo(ctx context.Context) error {
	if w.mongoClient == nil {
		return fmt.Errorf("mongo client is not connected")
	}

	result, err := w.mongoClient.Database(ordersDatabase).Collection(ordersCollection).DeleteMany(ctx, bson.M{})
	if err != nil {
		return err
	}

	fmt.Printf("MongoDB reset, deleted %d orders\n", result.DeletedCount)

	return nil
}

func (w *World) closeMongo(ctx context.Context) error {
	if w.mongoClient == nil {
		return nil
	}

	err := w.mongoClient.Disconnect(ctx)

	w.mongoClient = nil

	return err
}
