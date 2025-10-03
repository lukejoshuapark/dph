package posthog

type CreateFlagRequestModel struct {
	Key         string `json:"key"`
	Description string `json:"name"`
	Active      bool   `json:"active"`
}
