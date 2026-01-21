package application

import (
	"context"
	"fmt"

	"github.com/lukejoshuapark/dph/config"
	"github.com/lukejoshuapark/dph/posthog"
	"github.com/lukejoshuapark/dph/schema"
)

func PersistFlagsToFile(ctx context.Context, cfg *config.Config, filePath string) error {
	posthogService := posthog.NewService(cfg.ApiBaseUrl, cfg.PersonalApiKey)

	flags, err := posthogService.ListFlags(ctx, cfg.ProjectId)
	if err != nil {
		return fmt.Errorf("failed to list flags: %w", err)
	}

	definition := &schema.Definition{
		Flags: map[string]*schema.FlagDefinition{},
	}

	for _, flag := range flags {
		definition.Flags[flag.Key] = &schema.FlagDefinition{
			Description: flag.Description,
		}
	}

	if err := schema.PersistToFile(filePath, definition); err != nil {
		return fmt.Errorf("failed to write flag definitions to file: %w", err)
	}

	return nil
}
