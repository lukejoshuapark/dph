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

func ProcessFlag(ctx context.Context, notificationService notification.Service, posthogService posthog.Service, cfg *config.Config, flag *schema.Flag, dryRun bool) error {
	existingFlag, err := posthogService.GetFlagByKey(ctx, cfg.ProjectId, flag.Key)
	if err != nil {
		return fmt.Errorf("failed to get flag by key '%s': %w", flag.Key, err)
	}

	if existingFlag == nil {
		return createNewFlag(ctx, notificationService, posthogService, cfg, flag, dryRun)
	} else {
		return updateExistingFlag(ctx, posthogService, cfg, flag, existingFlag, dryRun)
	}
}

func createNewFlag(ctx context.Context, notificationService notification.Service, posthogService posthog.Service, cfg *config.Config, flag *schema.Flag, dryRun bool) error {
	if dryRun {
		fmt.Printf("✨ [DRY RUN] Would have created flag '%s'\n", flag.Key)
		return nil
	}

	requestModel := &posthog.CreateFlagRequestModel{
		Key:         flag.Key,
		Description: flag.Description,
		Active:      false,
	}

	flagId, err := posthogService.CreateFlag(ctx, cfg.ProjectId, requestModel)
	if err != nil {
		return fmt.Errorf("failed to create flag '%s': %w", flag.Key, err)
	}

	flagUrl := fmt.Sprintf("%s/project/%s/feature_flags/%d", cfg.ApiBaseUrl, cfg.ProjectId, flagId)

	if err := notificationService.PushCreateNotification(flag.Key, flag.Description, flagUrl); err != nil {
		fmt.Printf("⚠️ Failed to send create notification for flag '%s': %v\n", flag.Key, err)
	}

	fmt.Printf("✨ Created flag '%s'\n", flag.Key)

	return nil
}

func updateExistingFlag(ctx context.Context, posthogService posthog.Service, cfg *config.Config, flag *schema.Flag, existingFlag *posthog.FlagResponseModel, dryRun bool) error {
	if flag.Description == existingFlag.Description {
		return nil
	}

	if dryRun {
		fmt.Printf("👆 [DRY RUN] Would have updated flag '%s'\n", flag.Key)
		return nil
	}

	requestModel := &posthog.PatchFlagRequestModel{
		Description: util.Ptr(flag.Description),
	}

	if err := posthogService.PatchFlag(ctx, cfg.ProjectId, existingFlag.Id, requestModel); err != nil {
		return fmt.Errorf("failed to update flag '%s' (ID: %d): %w", flag.Key, existingFlag.Id, err)
	}

	fmt.Printf("👆 Updated flag '%s' (ID: %d)\n", flag.Key, existingFlag.Id)
	return nil
}
