package config

import "github.com/lukejoshuapark/environment"

type Config struct {
	ProjectId            string `environment:"DPH_PROJECT_ID"`
	PersonalApiKey       string `environment:"DPH_PERSONAL_API_KEY"`
	ApiBaseUrl           string `environment:"DPH_API_BASE_URL,https://us.posthog.com"`
	GoogleChatWebhookUrl string `environment:"DPH_GOOGLE_CHAT_WEBHOOK_URL,"`
}

func LoadFromEnvironment() (*Config, error) {
	cfg := &Config{}
	if err := environment.Populate(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
