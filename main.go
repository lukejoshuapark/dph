package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/lukejoshuapark/dph/application"
	"github.com/lukejoshuapark/dph/config"
	"github.com/lukejoshuapark/dph/schema"
	"github.com/lukejoshuapark/dph/util"
)

var filePath string
var dryRun bool
var reverse bool

func init() {
	flag.StringVar(&filePath, "f", "flags.yml", "Path to the flag definition file")
	flag.BoolVar(&dryRun, "d", false, "If set, will not make any changes to PostHog")
	flag.BoolVar(&reverse, "r", false, "If set, will operate in reverse - the flags currently in PostHog will be populated in the specified flag definition file")
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	flag.Parse()

	cfg, err := config.LoadFromEnvironment()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if reverse {
		if err := application.PersistFlagsToFile(ctx, cfg, filePath); err != nil {
			return fmt.Errorf("failed to persist flags to file: %w", err)
		}

		return nil
	}

	definition, err := schema.LoadFromFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to load definition: %w", err)
	}

	flags := util.MapMap(definition.Flags, func(key string, value *schema.FlagDefinition) *schema.Flag {
		return schema.FlagFromDefinition(key, value)
	})

	return application.ProcessFlags(ctx, cfg, flags, dryRun)
}
