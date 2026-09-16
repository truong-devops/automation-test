package rpc

import (
	"context"
	"errors"
	"time"

	"automation-test/payments/internal/domain"
	"automation-test/payments/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const paymentMethod = "/payments.v1.PaymentService/GetPaymentByOrderID"

type PaymentRPCServer interface {
	GetPaymentByOrderID(context.Context, *wrapperspb.StringValue) (*structpb.Struct, error)
}

type Server struct {
	payments *usecase.PaymentUsecase
}

func RegisterPaymentServer(registrar grpc.ServiceRegistrar, payments *usecase.PaymentUsecase) {
	registrar.RegisterService(&grpc.ServiceDesc{
		ServiceName: "payments.v1.PaymentService",
		HandlerType: (*PaymentRPCServer)(nil),
		Methods:     []grpc.MethodDesc{{MethodName: "GetPaymentByOrderID", Handler: getPaymentByOrderIDHandler}},
	}, &Server{payments: payments})
}

func (s *Server) GetPaymentByOrderID(ctx context.Context, request *wrapperspb.StringValue) (*structpb.Struct, error) {
	payment, err := s.payments.GetByOrderID(ctx, request.Value)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	result, err := structpb.NewStruct(map[string]any{
		"id": payment.ID, "order_id": payment.OrderID, "amount": payment.Amount,
		"method": payment.Method, "status": payment.Status,
		"created_at": payment.CreatedAt.Format(time.RFC3339Nano),
		"updated_at": payment.UpdatedAt.Format(time.RFC3339Nano),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return result, nil
}

func getPaymentByOrderIDHandler(server any, ctx context.Context, decode func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	request := new(wrapperspb.StringValue)
	if err := decode(request); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return server.(PaymentRPCServer).GetPaymentByOrderID(ctx, request)
	}
	info := &grpc.UnaryServerInfo{Server: server, FullMethod: paymentMethod}
	handler := func(ctx context.Context, request any) (any, error) {
		return server.(PaymentRPCServer).GetPaymentByOrderID(ctx, request.(*wrapperspb.StringValue))
	}
	return interceptor(ctx, request, info, handler)
}
