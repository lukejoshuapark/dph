package posthog

type PaginatedResponseModel[T any] struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []*T    `json:"results"`
}
