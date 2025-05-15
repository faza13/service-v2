package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	grpcController "service/app/controllers/grpc"
	"service/config"
)

func RunGrpcServer(ctx context.Context, cfg *config.Config, grpcController *grpcController.Grpc) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Grpc.Port))
	if err != nil {
		log.Fatalf("grpc failed to listen: %v", err)
	}

	s := grpc.NewServer()

	if err := s.Serve(lis); err != nil {
		log.Fatalf("grpc failed to serve: %v", err)
	}

}
