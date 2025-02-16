package usecase

import (
	"net/http"
	"tbank/internal/config"
	"tbank/internal/models/requests"
	"tbank/internal/models/responses"
	"tbank/internal/storage"
)



type UseCase interface {
	RegisterChat(id int) 	error
	DeleteChat(id int) 		error 
	Help()					error 
	AddLink(chatID int, addLinkRequest requests.AddLinkRequest) 			(responses.LinkResponse, error)
	RemoveLink(chatID int, removeLinkRequest requests.RemoveLinkRequest) 	(responses.LinkResponse, error)
	ListLinks(chatID int) 													(responses.ListLinksResponse, error) 
}


type UseCaseImpl struct {
	сfg config.Config
	client http.Client
	storage storage.Storage
}

func (uc *UseCaseImpl) RegisterChat(id int) error {
	return nil
}

func (uc *UseCaseImpl) DeleteChat(id int) error {
	return nil
}

func (uc *UseCaseImpl) Help() error {
	return nil
}

func (uc *UseCaseImpl) AddLink(chatID int, addLinkRequest requests.AddLinkRequest) (responses.LinkResponse, error) {

	return responses.LinkResponse{}, nil
}

func (uc *UseCaseImpl) RemoveLink(chatID int, removeLinkRequest requests.RemoveLinkRequest) (responses.LinkResponse, error) {
	return responses.LinkResponse{}, nil
}

func (uc *UseCaseImpl) ListLinks(chatID int) (responses.ListLinksResponse, error) {
	return responses.ListLinksResponse{}, nil
}