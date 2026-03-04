package posthog

type CreateFlagRequestModel struct {
	Key         string                        `json:"key"`
	Description string                        `json:"name"`
	Active      bool                          `json:"active"`
	Filters     CreateFlagRequestModelFilters `json:"filters"`
}

type CreateFlagRequestModelFilters struct {
	Groups       []CreateFlagRequestModelFiltersGroup `json:"groups"`
	Payloads     map[string]any                       `json:"payloads"`
	Multivariant *struct{}                            `json:"multivariant"`
}

type CreateFlagRequestModelFiltersGroup struct {
	Variant           *struct{} `json:"variant"`
	Properties        []any     `json:"properties"`
	RolloutPercentage int       `json:"rollout_percentage"`
}

func NewCreateFlagRequestModel(key string, description string) *CreateFlagRequestModel {
	return &CreateFlagRequestModel{
		Key:         key,
		Description: description,
		Filters: CreateFlagRequestModelFilters{
			Groups: []CreateFlagRequestModelFiltersGroup{
				{
					Properties: []any{},
				},
			},
			Payloads: map[string]any{},
		},
	}
}
