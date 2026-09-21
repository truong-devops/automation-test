package automationtests

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type World struct {
	baseURL     string
	httpClient  *http.Client
	mongoClient *mongo.Client

	aliases map[string]string

	lastStatus int
	lastBody   []byte
}

func newWorld() *World {
	return &World{
		baseURL: "http://localhost:8081",
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		aliases: make(map[string]string),
	}
}

func (w *World) connectMongo(ctx context.Context) error {
	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI("mongodb://localhost:27018"),
	)
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

func (w *World) closeMongo(ctx context.Context) error {
	if w.mongoClient == nil {
		return nil
	}

	return w.mongoClient.Disconnect(ctx)
}

func (w *World) resetMongo(ctx context.Context) error {
	result, err := w.mongoClient.Database("orders_db").Collection("orders").DeleteMany(ctx, bson.M{})

	if err != nil {
		return err
	}

	fmt.Printf("MongoDB reset, deleted %d orders\n", result.DeletedCount)

	return nil
}
