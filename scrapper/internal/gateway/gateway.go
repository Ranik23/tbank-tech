package gateway

import (
	"context"
	"log/slog"
	"net/http"
	"tbank/scrapper/api/proto/gen"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func RunGateway(ctx context.Context, grpcAddr string, httpAddr string, logger *slog.Logger) error {
	const op = "Gateway.RunGateway"
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithInsecure()}

	if err := gen.RegisterScrapperHandlerFromEndpoint(ctx, mux, grpcAddr, opts); err != nil {
		logger.Error(op, slog.String("message", "Failed to register gRPC handler"), slog.String("error", err.Error()))
		return err
	}

	logger.Info(op, slog.String("message", "Starting proxy server"), slog.String("httpAddr", httpAddr))
	return http.ListenAndServe(httpAddr, mux)
}
