package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"automation-test/payments/configs"
	httpdelivery "automation-test/payments/internal/delivery/http"
	"automation-test/payments/internal/delivery/rpc"
	"automation-test/payments/internal/repository"
	"automation-test/payments/internal/usecase"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func main() {
	config := configs.Load()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		log.Fatal(err)
	}
	pingContext, cancelPing := context.WithTimeout(ctx, 10*time.Second)
	err = mongoClient.Ping(pingContext, nil)
	cancelPing()
	if err != nil {
		log.Fatalf("connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(context.Background())

	repository := repository.NewMongoPaymentRepository(mongoClient.Database(config.MongoDatabase))
	if err := repository.EnsureIndexes(ctx); err != nil {
		log.Fatalf("create MongoDB indexes: %v", err)
	}
	payments := usecase.NewPaymentUsecase(repository)
	orderClient, err := rpc.NewOrderClient(config.OrdersGRPCAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer orderClient.Close()

	grpcListener, err := net.Listen("tcp", config.GRPCAddress)
	if err != nil {
		log.Fatal(err)
	}
	grpcServer := grpc.NewServer()
	rpc.RegisterPaymentServer(grpcServer, payments)
	healthServer := health.NewServer()
	healthServer.SetServingStatus("payments.v1.PaymentService", healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	reflection.Register(grpcServer)
	go func() {
		log.Printf("payments gRPC listening on %s", config.GRPCAddress)
		if err := grpcServer.Serve(grpcListener); err != nil {
			log.Printf("gRPC server stopped: %v", err)
		}
	}()

	httpServer := &http.Server{
		Addr: config.HTTPAddress, Handler: httpdelivery.NewHandler(payments, orderClient).Routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}
	go func() {
		log.Printf("payments HTTP listening on %s", config.HTTPAddress)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("HTTP server stopped: %v", err)
			stop()
		}
	}()

	<-ctx.Done()
	shutdownContext, cancelShutdown := context.WithTimeout(context.Background(), config.ShutdownTimeout)
	defer cancelShutdown()
	_ = httpServer.Shutdown(shutdownContext)
	grpcServer.GracefulStop()
}
