package payments

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"payments-app/internal/config"
	"payments-app/internal/database"
	entities "payments-app/internal/entities/payments"
	"payments-app/rpc/proto"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PaymentServiceServerImpl struct {
	proto.UnimplementedPaymentServiceServer
	Config config.Config
}

type Server struct {
	config config.Config
}

func NewServer() (*Server, error) {
	fmt.Println("Hello we are here")
	var appConfig config.Config
	if err := config.LoadConfig(&appConfig); err != nil {
		return nil, err
	}

	return &Server{
		config: appConfig,
	}, nil
}

func (s *Server) setupGRPCServer() (*grpc.Server, net.Listener, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.Server.Port))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create listener: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterPaymentServiceServer(grpcServer, &PaymentServiceServerImpl{Config: s.config})

	return grpcServer, listener, nil
}

func (s *Server) setupGatewayHandler(ctx context.Context) (http.Handler, error) {
	gwmux := runtime.NewServeMux(runtime.WithErrorHandler(runtime.DefaultHTTPErrorHandler))

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	endpoint := fmt.Sprintf("localhost:%d", s.config.Server.Port)

	if err := proto.RegisterPaymentServiceHandlerFromEndpoint(ctx, gwmux, endpoint, opts); err != nil {
		return nil, fmt.Errorf("failed to register gateway: %v", err)
	}

	return gwmux, nil
}

func (s *Server) startGRPCServer(grpcServer *grpc.Server, listener net.Listener) {
	log.Printf("Starting gRPC server on port %d...", s.config.Server.Port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
		return
	}
	log.Printf("grpc started")
}

func (s *Server) startHTTPServer(handler http.Handler) error {
	log.Printf("Starting HTTP server on port %d...", s.config.Server.HTTPPort)
	return http.ListenAndServe(fmt.Sprintf(":%d", s.config.Server.HTTPPort), handler)
}

func (s *Server) initDatabase() error {
	return database.InitDB(s.config)
}

func (s *Server) autoMigrate() error {
	// AutoMigrate payment table
	log.Println("Running migration..")
	err := database.DB.AutoMigrate(&entities.Payment{})
	if err != nil {
		return err
	}
	log.Println("Database migration completed.")
	return nil
}

func (s *Server) Start() error {
	if err := s.initDatabase(); err != nil {
		return fmt.Errorf("failed to initialize database: %v", err)
	}

	if err := s.autoMigrate(); err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	grpcServer, listener, err := s.setupGRPCServer()
	if err != nil {
		return fmt.Errorf("failed to setup gRPC server: %v", err)
	}

	ctx := context.Background()
	handler, err := s.setupGatewayHandler(ctx)
	if err != nil {
		return fmt.Errorf("failed to setup gateway handler: %v", err)
	}

	go s.startGRPCServer(grpcServer, listener)
	return s.startHTTPServer(handler)
}

func StartServer() error {
	server, err := NewServer()
	if err != nil {
		return fmt.Errorf("failed to create server: %v", err)
	}

	return server.Start()
}
