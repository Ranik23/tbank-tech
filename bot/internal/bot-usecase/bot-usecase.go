package botusecase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"tbank/bot/config"
	"tbank/bot/internal/models/requests"
	"tbank/bot/internal/models/responses"
	"tbank/bot/internal/storage"
)


var (
	ErrCodeUnknown = fmt.Errorf("response code unknown")
)


type UseCase interface {
	RegisterChat(сtx context.Context, id int64) 												error
	DeleteChat(ctx context.Context, id int64) 													error 
	Help(ctx context.Context)																	error 
	AddLink(ctx context.Context, chatID int64, link string, tags []string, filters []string) 	(*responses.LinkResponse, error)
	RemoveLink(ctx context.Context, chatID int64, link string) 									(*responses.LinkResponse, error)
	ListLinks(ctx context.Context, chatID int64) 												(*responses.ListLinksResponse, error) 
}


type UseCaseImpl struct {
	config 		*config.Config
	client 		http.Client
	logger 		*slog.Logger
	storage storage.Storage
}

func NewUseCaseImpl(config *config.Config, storage storage.Storage, logger *slog.Logger) *UseCaseImpl {
	return &UseCaseImpl{
		config: config,
		storage: storage,
		logger: logger,
	}
}

func (uc *UseCaseImpl) RegisterChat(ctx context.Context, chatID int64) error {

	url := fmt.Sprintf("%s:%s/tg-chat/%d", uc.config.ScrapperService.Host, uc.config.ScrapperService.Port, chatID)

	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return err
	}

	resp, err := uc.client.Do(req)
	if err != nil {
		return err
	}
	
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusBadRequest:

		var errorResponse responses.ApiErrorResponse

		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return err
		}
		return fmt.Errorf("error - %s", errorResponse.Description)
	default:
		return ErrCodeUnknown
	}
}

func (uc *UseCaseImpl) DeleteChat(ctx context.Context, chatID int64) error {
	url := fmt.Sprintf("%s:%s/tg-chat/%d", uc.config.ScrapperService.Host, uc.config.ScrapperService.Port, chatID)

	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := uc.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusBadRequest:
		var errorResponse responses.ApiErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return fmt.Errorf("failed to decode error response: %w", err)
		}
		return fmt.Errorf("bad request: %s", errorResponse.Description)

	case http.StatusNotFound:
		var errorResponse responses.ApiErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return fmt.Errorf("failed to decode error response: %w", err)
		}
		return fmt.Errorf("chat not found: %s", errorResponse.Description)

	default:
		return ErrCodeUnknown
	}
}


func (uc *UseCaseImpl) Help(ctx context.Context) error {
	return nil
}

func (uc *UseCaseImpl) AddLink(ctx context.Context, chatID int64, link string, tags []string, filters []string) (*responses.LinkResponse, error) {
	url := fmt.Sprintf("%s:%s/links", uc.config.ScrapperService.Host, uc.config.ScrapperService.Port)

	addLinkRequest := requests.AddLinkRequest{
		Link: link,
		Tags: tags,
		Filters: filters,
	}

	body, err := json.Marshal(addLinkRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Tg-Chat-Id", strconv.FormatInt(chatID, 10))

	resp, err := uc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {

	case http.StatusOK:
		var response responses.LinkResponse

		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return nil, err
		}

		return &response, nil

	case http.StatusBadRequest:
		var errorResponse responses.ApiErrorResponse

		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("error - %s", errorResponse.Description)

	default:
		return nil, ErrCodeUnknown
	}
}

func (uc *UseCaseImpl) RemoveLink(ctx context.Context, chatID int64, link string) (*responses.LinkResponse, error) {
	url := fmt.Sprintf("%s:%s/links", uc.config.ScrapperService.Host, uc.config.ScrapperService.Port)

	removeLinkRequest := requests.RemoveLinkRequest{
		Link: link,
	}

	body, err := json.Marshal(removeLinkRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "DELETE", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Tg-Chat-Id", strconv.FormatInt(chatID, 10))

	resp, err := uc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var response responses.LinkResponse
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return nil, err
		}
		return &response, nil

	case http.StatusBadRequest:
		var errorResponse responses.ApiErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("error - %s", errorResponse.Description)

	case http.StatusNotFound:
		var errorResponse responses.ApiErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("link not found - %s", errorResponse.Description)

	default:
		return nil, ErrCodeUnknown
	}
}

func (uc *UseCaseImpl) ListLinks(ctx context.Context, chatID int64) (*responses.ListLinksResponse, error) {

	url := fmt.Sprintf("%s:%s/links", uc.config.ScrapperService.Host, uc.config.ScrapperService.Port)


	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Tg-Chat-Id", strconv.FormatInt(chatID, 10))

	resp, err := uc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:

		var response responses.ListLinksResponse

		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return nil, err
		
		}
		return &response, nil

	case http.StatusBadRequest:

		var errorResponse responses.ApiErrorResponse

		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("error - %s", errorResponse.Description)
	default:

		return nil, ErrCodeUnknown
	}
}