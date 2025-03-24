package grpc

import (
	"context"
	"errors"
	"log"

	"github.com/Ranik23/tbank-tech/scrapper/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrorHandlingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	resp, err := handler(ctx, req)
	if err != nil {
		log.Printf("Error in method %s: %v", info.FullMethod, err)

		switch {
		case errors.Is(err, service.ErrUserAlreadyExists):
			return nil, status.Errorf(codes.AlreadyExists, "User already exists")
		case errors.Is(err, service.ErrUserNotFound):
			return nil, status.Errorf(codes.NotFound, "User not found")
		default:
			return nil, status.Errorf(codes.Internal, "Internal Server Error")
		}
	}
	return resp, nil
}




