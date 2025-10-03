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

func init() {
	flag.StringVar(&filePath, "f", "flags.yml", "Path to the flag definition file")
	flag.BoolVar(&dryRun, "d", false, "If set, will not make any changes to PostHog")
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

	definition, err := schema.LoadFromFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to load definition: %w", err)
	}

	flags := util.MapMap(definition.Flags, func(key string, value *schema.FlagDefinition) *schema.Flag {
		return schema.FlagFromDefinition(key, value)
	})

	return application.ProcessFlags(ctx, cfg, flags, dryRun)
}
