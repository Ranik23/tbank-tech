package httpserver

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"tbank/scrapper/api/proto/gen"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func RunGateway(ctx context.Context, grpcAddr string, httpAddr string) error {
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithInsecure()}

	if err := gen.RegisterScrapperHandlerFromEndpoint(ctx, mux, grpcAddr, opts); err != nil {
		return err
	}

	fmt.Println("адрес - ", httpAddr)
	slog.Info("Запуск прокси-сервера", slog.String("httpAddr", httpAddr))
	return http.ListenAndServe(httpAddr, mux)
}