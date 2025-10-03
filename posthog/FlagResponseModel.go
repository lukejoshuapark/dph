package posthog

type FlagResponseModel struct {
	Id          int    `json:"id"`
	Key         string `json:"key"`
	Description string `json:"name"`
	Deleted     bool   `json:"deleted"`
}
