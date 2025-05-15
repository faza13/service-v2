package grpc

import (
	"context"
	"service/app/usecases"
)

type Grpc struct{}

func NewGrpc(ctx context.Context, usecase *usecases.Usecase) *Grpc {
	return &Grpc{}
}
