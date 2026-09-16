package rpc

import (
	"context"
	"errors"
	"time"

	"automation-test/orders/internal/domain"
	"automation-test/orders/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const orderMethod = "/orders.v1.OrderService/GetOrder"

type OrderRPCServer interface {
	GetOrder(context.Context, *wrapperspb.StringValue) (*structpb.Struct, error)
}

type Server struct {
	orders *usecase.OrderUsecase
}

func RegisterOrderServer(registrar grpc.ServiceRegistrar, orders *usecase.OrderUsecase) {
	registrar.RegisterService(&grpc.ServiceDesc{
		ServiceName: "orders.v1.OrderService",
		HandlerType: (*OrderRPCServer)(nil),
		Methods:     []grpc.MethodDesc{{MethodName: "GetOrder", Handler: getOrderHandler}},
	}, &Server{orders: orders})
}

func (s *Server) GetOrder(ctx context.Context, request *wrapperspb.StringValue) (*structpb.Struct, error) {
	order, err := s.orders.GetByID(ctx, request.Value)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) || errors.Is(err, domain.ErrInvalidID) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	result, err := structpb.NewStruct(map[string]any{
		"id": order.ID, "customer_name": order.CustomerName, "item": order.Item,
		"quantity": order.Quantity, "total_amount": order.TotalAmount,
		"status": order.Status, "created_at": order.CreatedAt.Format(time.RFC3339Nano),
		"updated_at": order.UpdatedAt.Format(time.RFC3339Nano),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return result, nil
}

func getOrderHandler(server any, ctx context.Context, decode func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	request := new(wrapperspb.StringValue)
	if err := decode(request); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return server.(OrderRPCServer).GetOrder(ctx, request)
	}
	info := &grpc.UnaryServerInfo{Server: server, FullMethod: orderMethod}
	handler := func(ctx context.Context, request any) (any, error) {
		return server.(OrderRPCServer).GetOrder(ctx, request.(*wrapperspb.StringValue))
	}
	return interceptor(ctx, request, info, handler)
}
