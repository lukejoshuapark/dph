package application

import (
	"context"
	"fmt"

	"github.com/lukejoshuapark/dph/config"
	"github.com/lukejoshuapark/dph/posthog"
	"github.com/lukejoshuapark/dph/schema"
	"github.com/lukejoshuapark/dph/util"
)

func DeleteSurplusFlags(ctx context.Context, posthogService posthog.Service, cfg *config.Config, flags []*schema.Flag, dryRun bool) error {
	currentFlags, err := posthogService.ListFlags(ctx, cfg.ProjectId)
	if err != nil {
		return fmt.Errorf("failed to list flags: %w", err)
	}

	currentFlagKeys := util.Map(currentFlags, func(f *posthog.FlagResponseModel) string { return f.Key })
	definedFlagKeys := util.Map(flags, func(f *schema.Flag) string { return f.Key })
	surplusFlagKeys := util.Difference(currentFlagKeys, definedFlagKeys)

	if len(surplusFlagKeys) < 1 {
		return nil
	}

	return util.Parallel(surplusFlagKeys, func(key string) error {
		flag, err := posthogService.GetFlagByKey(ctx, cfg.ProjectId, key)
		if err != nil {
			return fmt.Errorf("failed to get flag by key '%s': %w", key, err)
		}

		if flag == nil {
			return fmt.Errorf("could not find flag with key '%s'", key)
		}

		if dryRun {
			fmt.Printf("☠️ [DRY RUN] Would have deleted flag with key '%s' (ID: %d)\n", key, flag.Id)
			return nil
		}

		requestModel := &posthog.PatchFlagRequestModel{
			Deleted: util.Ptr(true),
		}

		if err := posthogService.PatchFlag(ctx, cfg.ProjectId, flag.Id, requestModel); err != nil {
			return fmt.Errorf("failed to delete flag with key '%s' (ID: %d): %w", key, flag.Id, err)
		}

		fmt.Printf("☠️ Deleted flag with key '%s' (ID: %d)\n", key, flag.Id)
		return nil
	})
}
