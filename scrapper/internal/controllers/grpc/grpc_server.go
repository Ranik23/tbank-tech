package grpc

import (
	"context"
	"errors"
	"time"

	"github.com/Ranik23/tbank-tech/scrapper/api/proto/gen"
	"github.com/Ranik23/tbank-tech/scrapper/internal/metrics"
	"github.com/Ranik23/tbank-tech/scrapper/internal/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ScrapperServer struct {
	gen.UnimplementedScrapperServer
	usecase service.Service
}

func NewScrapperServer(usecase service.Service) *ScrapperServer {
	return &ScrapperServer{
		usecase: usecase,
	}
}

func (s *ScrapperServer) RegisterUser(ctx context.Context, req *gen.RegisterUserRequest) (*gen.RegisterUserResponse, error) {
	start := time.Now()
	metrics.TotalRequests.Inc()
	if err := s.usecase.RegisterUser(ctx, uint(req.GetTgUserId()), req.GetName(), req.GetToken()); err != nil {
		metrics.ErrorRequests.Inc()
		switch {
		case errors.Is(err, service.ErrUserAlreadyExists):
			return nil, status.Errorf(codes.AlreadyExists, "Already Exists")
		default:
			return nil, status.Errorf(codes.Internal, "Internal Server Error")
		}
	}
	
	duration := time.Since(start)
	metrics.RequestDuration.Observe(duration.Seconds())

	return &gen.RegisterUserResponse{Message: "Пользователь зарегистрирован!"}, nil
}

func (s *ScrapperServer) DeleteUser(ctx context.Context, req *gen.DeleteUserRequest) (*gen.DeleteUserResponse, error) {
	start := time.Now()
	metrics.TotalRequests.Inc()
	if err := s.usecase.DeleteUser(ctx, uint(req.GetTgUserId())); err != nil {
		metrics.ErrorRequests.Inc()
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			return nil, status.Errorf(codes.NotFound, "User Not Found")
		default:
			return nil, status.Errorf(codes.Internal, "Internal Server Error")
		}
	}
	duration := time.Since(start)
	metrics.RequestDuration.Observe(duration.Seconds())

	return &gen.DeleteUserResponse{Message: "Пользователь удален"}, nil
}

func (s *ScrapperServer) GetLinks(ctx context.Context, req *gen.GetLinksRequest) (*gen.ListLinksResponse, error) {
	start := time.Now()
	metrics.TotalRequests.Inc()
	links, err := s.usecase.GetLinks(ctx, uint(req.GetTgUserId()))
	if err != nil {
		metrics.ErrorRequests.Inc()
		return nil, status.Errorf(codes.Internal, "Internal Server Error")
	}

	var linksResponse []string
	for _, link := range links {
		linksResponse = append(linksResponse, link.Url)
	}
	duration := time.Since(start)
	metrics.RequestDuration.Observe(duration.Seconds())
	return &gen.ListLinksResponse{Links: linksResponse}, nil
}

func (s *ScrapperServer) AddLink(ctx context.Context, req *gen.AddLinkRequest) (*gen.AddLinkResponse, error) {
	start := time.Now()
	metrics.TotalRequests.Inc()
	if err := s.usecase.AddLink(ctx, req.GetUrl(), uint(req.GetTgUserId())); err != nil {
		metrics.ErrorRequests.Inc()
		return nil, status.Errorf(codes.Internal, "Internal Server Error")
	}
	duration := time.Since(start)
	metrics.RequestDuration.Observe(duration.Seconds())
	return &gen.AddLinkResponse{Message: "Успешно добавили ссылку!"}, nil
}

func (s *ScrapperServer) RemoveLink(ctx context.Context, req *gen.RemoveLinkRequest) (*gen.RemoveLinkResponse, error) {
	start := time.Now()
	metrics.TotalRequests.Inc()
	if err := s.usecase.RemoveLink(ctx, req.GetUrl(), uint(req.GetTgUserId())); err != nil {
		metrics.ErrorRequests.Inc()
		return nil, status.Errorf(codes.Internal, "Internal Server Error")
	}
	duration := time.Since(start)
	metrics.RequestDuration.Observe(duration.Seconds())
	return &gen.RemoveLinkResponse{Message: "Успешно удалили ссылку"}, nil
}
