package grpc

import (
	"context"

	"github.com/Ranik23/tbank-tech/scrapper/api/proto/gen"
	"github.com/Ranik23/tbank-tech/scrapper/internal/service"
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
	if err := s.usecase.RegisterUser(ctx, uint(req.GetTgUserId()), req.GetName(), req.GetToken()); err != nil {
		return nil, err
	}
	return &gen.RegisterUserResponse{
		Message: "User Created",
	}, nil
}

func (s *ScrapperServer) DeleteUser(ctx context.Context, req *gen.DeleteUserRequest) (*gen.DeleteUserResponse, error) {
	if err := s.usecase.DeleteUser(ctx, uint(req.GetTgUserId())); err != nil {
		return nil, err
	}

	return &gen.DeleteUserResponse{
		Message: "User Deleted",
	}, nil
}

func (s *ScrapperServer) GetLinks(ctx context.Context, req *gen.GetLinksRequest) (*gen.ListLinksResponse, error) {
	links, err := s.usecase.GetLinks(ctx, uint(req.GetTgUserId()))
	if err != nil {
		return nil, err
	}

	var linksResponse []*gen.LinkResponse

	for _, link := range links {
		linksResponse = append(linksResponse, &gen.LinkResponse{
			Url: link.Url,
			Id:  int64(link.ID),
		})
	}

	return &gen.ListLinksResponse{
		Links: linksResponse,
		Size:  int32(len(linksResponse)),
	}, nil
}

func (s *ScrapperServer) AddLink(ctx context.Context, req *gen.AddLinkRequest) (*gen.LinkResponse, error) {
	if err := s.usecase.AddLink(ctx, req.GetUrl(), uint(req.GetTgUserId())); err != nil {
		return nil, err
	}
	return &gen.LinkResponse{
		Id:  1,
		Url: req.GetUrl(),
	}, nil
}

func (s *ScrapperServer) RemoveLink(ctx context.Context, req *gen.RemoveLinkRequest) (*gen.LinkResponse, error) {
	if err := s.usecase.RemoveLink(ctx, req.GetUrl(), uint(req.GetTgUserId())); err != nil {
		return nil, err
	}

	return &gen.LinkResponse{
		Id:  1,
		Url: req.GetUrl(),
	}, nil
}
