package gateway

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"

	"github.com/Ranik23/tbank-tech/bot/api/proto/gen"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func RunGateway(ctx context.Context, grpcAddr string, httpAddr string, logger *slog.Logger) error {
	const op = "Gateway.RunGateway"
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := gen.RegisterBotHandlerFromEndpoint(ctx, mux, grpcAddr, opts); err != nil {
		logger.Error(op, slog.String("message", "Failed to register gRPC handler"), slog.String("error", err.Error()))
		return err
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error(op, slog.String("message", "Failed to read request body"), slog.String("error", err.Error()))
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		r.Body = io.NopCloser(bytes.NewReader(body))
		logger.Info(op, slog.String("incoming_request", string(body)))
		mux.ServeHTTP(w, r)
	}) 

	logger.Info(op, slog.String("message", "Starting Bot Proxy Server"), slog.String("httpAddr", httpAddr))
	return http.ListenAndServe(httpAddr, handler)
}
