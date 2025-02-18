package grpcserver

import (
	"context"
	"tbank/scrapper/internal/storage"
	"tbank/scrapper/api/proto/gen"
	"tbank/scrapper/internal/usecase"
)

type ScrapperServer struct {
	gen.UnimplementedScrapperServer
	usecase usecase.UseCase
	storage storage.Storage
}

func NewScrapperServer() *ScrapperServer {
	return &ScrapperServer{}
}

func (s *ScrapperServer) RegisterChat(ctx context.Context, req *gen.RegisterChatRequest) (*gen.RegisterChatResponse, error) {
	return &gen.RegisterChatResponse{Message: "Чат зарегистрирован"}, nil
}

func (s *ScrapperServer) DeleteChat(ctx context.Context, req *gen.DeleteChatRequest) (*gen.DeleteChatResponse, error) {
	return &gen.DeleteChatResponse{Message: "Чат удалён"}, nil
}

func (s *ScrapperServer) GetLinks(ctx context.Context, req *gen.GetLinksRequest) (*gen.ListLinksResponse, error) {
	return &gen.ListLinksResponse{
		Links: []*gen.LinkResponse{
			{Id: 1, Url: "https://example.com", Tags: []string{"новости"}, Filters: []string{"акции"}},
		},
		Size: 1,
	}, nil
}

func (s *ScrapperServer) AddLink(ctx context.Context, req *gen.AddLinkRequest) (*gen.LinkResponse, error) {
	return &gen.LinkResponse{
		Id:      123,
		Url:     req.Link,
		Tags:    req.Tags,
		Filters: req.Filters,
	}, nil
}

func (s *ScrapperServer) RemoveLink(ctx context.Context, req *gen.RemoveLinkRequest) (*gen.LinkResponse, error) {
	return &gen.LinkResponse{
		Id: 1, Url: req.Link,
	}, nil
}
