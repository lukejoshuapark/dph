package posthog

type PatchFlagRequestModel struct {
	Description *string `json:"name,omitempty"`
	Deleted     *bool   `json:"deleted,omitempty"`
}
