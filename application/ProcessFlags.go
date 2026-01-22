package application

import (
	"context"
	"fmt"

	"github.com/lukejoshuapark/dph/config"
	"github.com/lukejoshuapark/dph/notification"
	"github.com/lukejoshuapark/dph/posthog"
	"github.com/lukejoshuapark/dph/schema"
	"github.com/lukejoshuapark/dph/util"
)

func ProcessFlags(ctx context.Context, cfg *config.Config, flags []*schema.Flag, dryRun bool) error {
	var notificationService notification.Service
	if cfg.GoogleChatWebhookUrl != "" {
		notificationService = notification.NewGoogleChatService(cfg.GoogleChatWebhookUrl)
	} else {
		notificationService = notification.NewEmptyService()
	}

	posthogService := posthog.NewService(cfg.ApiBaseUrl, cfg.PersonalApiKey)
	if err := DeleteSurplusFlags(ctx, notificationService, posthogService, cfg, flags, dryRun); err != nil {
		return fmt.Errorf("failed to delete surplus flags: %w", err)
	}

	return util.Parallel(flags, func(flag *schema.Flag) error {
		return ProcessFlag(ctx, notificationService, posthogService, cfg, flag, dryRun)
	})
}
