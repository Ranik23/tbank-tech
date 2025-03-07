package grpcserver

import (
	"context"
	"tbank/scrapper/api/proto/gen"
	"tbank/scrapper/internal/models"
	"tbank/scrapper/internal/usecase"
)

type ScrapperServer struct {
	gen.UnimplementedScrapperServer
	usecase usecase.UseCase
}

func NewScrapperServer(usecase usecase.UseCase) *ScrapperServer {
	return &ScrapperServer{
		usecase: usecase,
	}
}

func (s *ScrapperServer) RegisterUser(ctx context.Context, req *gen.RegisterUserRequest) (*gen.RegisterUserResponse, error) {
	if err := s.usecase.RegisterUser(ctx, uint(req.GetTgUserId()), req.GetName()); err != nil {
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
			Id: int64(link.ID),
		})
	}

	return &gen.ListLinksResponse{
		Links: linksResponse,
		Size: int32(len(linksResponse)),
	}, nil
}

func (s *ScrapperServer) AddLink(ctx context.Context, req *gen.AddLinkRequest) (*gen.LinkResponse, error) {

	link := models.Link{
		Url: req.GetUrl(),
	}

	newLink, err := s.usecase.AddLink(ctx, link, uint(req.GetTgUserId()))
	if err != nil {
		return nil, err
	}

	return &gen.LinkResponse{
		Id: int64(newLink.ID),
		Url: newLink.Url,
	}, nil
}

func (s *ScrapperServer) RemoveLink(ctx context.Context, req *gen.RemoveLinkRequest) (*gen.LinkResponse, error) {
	
	link := models.Link{
		Url: req.GetUrl(),
	}

	if err := s.usecase.RemoveLink(ctx, link, uint(req.GetTgUserId())); err != nil {
		return nil, err
	}

	return &gen.LinkResponse{
		Id: int64(link.ID),
		Url: link.Url,
	}, nil
}
