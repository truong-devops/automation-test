package rpc

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const orderMethod = "/orders.v1.OrderService/GetOrder"

type OrderClient struct {
	connection *grpc.ClientConn
}

func NewOrderClient(address string) (*OrderClient, error) {
	connection, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &OrderClient{connection: connection}, nil
}

func (c *OrderClient) GetByID(ctx context.Context, orderID string) (map[string]any, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	response := new(structpb.Struct)
	if err := c.connection.Invoke(ctx, orderMethod, wrapperspb.String(orderID), response); err != nil {
		return nil, err
	}
	return response.AsMap(), nil
}

func (c *OrderClient) Close() error { return c.connection.Close() }
