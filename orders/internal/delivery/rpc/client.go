package rpc

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const paymentMethod = "/payments.v1.PaymentService/GetPaymentByOrderID"

type PaymentClient struct {
	connection *grpc.ClientConn
}

func NewPaymentClient(address string) (*PaymentClient, error) {
	connection, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &PaymentClient{connection: connection}, nil
}

func (c *PaymentClient) GetByOrderID(ctx context.Context, orderID string) (map[string]any, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	response := new(structpb.Struct)
	if err := c.connection.Invoke(ctx, paymentMethod, wrapperspb.String(orderID), response); err != nil {
		return nil, err
	}
	return response.AsMap(), nil
}

func (c *PaymentClient) Close() error { return c.connection.Close() }
