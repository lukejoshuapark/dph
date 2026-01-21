package posthog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/lukejoshuapark/dph/util"
)

type Service interface {
	ListFlags(ctx context.Context, projectId string) ([]*FlagResponseModel, error)
	GetFlagByKey(ctx context.Context, projectId string, key string) (*FlagResponseModel, error)
	CreateFlag(ctx context.Context, projectId string, requestModel *CreateFlagRequestModel) error
	PatchFlag(ctx context.Context, projectId string, id int, requestModel *PatchFlagRequestModel) error
}

func NewService(baseUrl string, apiKey string) Service {
	return &impl{
		baseUrl: baseUrl,
		apiKey:  apiKey,
	}
}

type impl struct {
	baseUrl string
	apiKey  string
}

var _ Service = (*impl)(nil)

func (s *impl) ListFlags(ctx context.Context, projectId string) ([]*FlagResponseModel, error) {
	flags := make([]*FlagResponseModel, 0)
	next := (*string)(nil)

	for {
		page, err := s.listFlagPage(ctx, projectId, next)
		if err != nil {
			return nil, fmt.Errorf("failed to list flag page: %w", err)
		}

		flags = append(flags, page.Results...)
		if page.Next == nil {
			break
		}

		next = page.Next
	}

	return flags, nil
}

func (s *impl) listFlagPage(ctx context.Context, projectId string, next *string) (*PaginatedResponseModel[FlagResponseModel], error) {
	url := fmt.Sprintf("%s/api/projects/%s/feature_flags", s.baseUrl, projectId)
	if next != nil {
		url = *next
	}

	req, err := prepareRequest(ctx, http.MethodGet, url, s.apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare request: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to complete http request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errorForFailureResponse(res)
	}

	return decodeResponseModel[PaginatedResponseModel[FlagResponseModel]](res)
}

func (s *impl) GetFlagByKey(ctx context.Context, projectId string, key string) (*FlagResponseModel, error) {
	url := fmt.Sprintf("%s/api/projects/%s/feature_flags?search=%s", s.baseUrl, projectId, key)
	req, err := prepareRequest(ctx, http.MethodGet, url, s.apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare request: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to complete http request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errorForFailureResponse(res)
	}

	page, err := decodeResponseModel[PaginatedResponseModel[FlagResponseModel]](res)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	result, err := util.Single(page.Results, func(model *FlagResponseModel) bool {
		return model.Key == key
	})

	if err == util.ErrTooManyElements {
		return nil, fmt.Errorf("multiple flags found with key %s", key)
	}

	if err == util.ErrNoElementFound {
		return nil, nil
	}

	return result, nil
}

func (s *impl) CreateFlag(ctx context.Context, projectId string, requestModel *CreateFlagRequestModel) error {
	url := fmt.Sprintf("%s/api/projects/%s/feature_flags", s.baseUrl, projectId)
	req, err := prepareRequestWithRequestModel(ctx, http.MethodPost, url, s.apiKey, requestModel)
	if err != nil {
		return fmt.Errorf("failed to prepare request: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to complete http request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		return errorForFailureResponse(res)
	}

	return nil
}

func (s *impl) PatchFlag(ctx context.Context, projectId string, id int, requestModel *PatchFlagRequestModel) error {
	url := fmt.Sprintf("%s/api/projects/%s/feature_flags/%d", s.baseUrl, projectId, id)
	req, err := prepareRequestWithRequestModel(ctx, http.MethodPatch, url, s.apiKey, requestModel)
	if err != nil {
		return fmt.Errorf("failed to prepare request: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to complete http request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errorForFailureResponse(res)
	}

	return nil
}

func prepareRequest(ctx context.Context, method string, url string, apiKey string) (*http.Request, error) {
	return prepareRequestWithRequestModel[any](ctx, method, url, apiKey, nil)
}

func prepareRequestWithRequestModel[T any](ctx context.Context, method string, url string, apiKey string, requestModel *T) (*http.Request, error) {
	body, contentLength, err := prepareRequestModel(requestModel)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare request body: %w", err)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Content-Length", fmt.Sprintf("%d", contentLength))
	}

	req = req.WithContext(ctx)
	return req, nil
}

func prepareRequestModel[T any](requestModel *T) (io.Reader, int, error) {
	if requestModel == nil {
		return nil, 0, nil
	}

	data, err := json.Marshal(requestModel)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to marshal request body: %w", err)
	}

	return bytes.NewBuffer(data), len(data), nil
}

func decodeResponseModel[T any](res *http.Response) (*T, error) {
	var responseModel T

	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&responseModel); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &responseModel, nil
}

func errorForFailureResponse(res *http.Response) error {
	text, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("request failed with status %d and failed to read body: %w", res.StatusCode, err)
	}

	return fmt.Errorf("request failed with status %d: %s", res.StatusCode, string(text))
}
