package main

import (
	"context"

	cp "github.com/otiai10/copy"
)

func main() {
	ctx := context.Background()
	logger := NewLogger()

	ctx = AddLoggerToContext(ctx, logger)

	config, err := GetConfig(ctx)
	if err != nil {
		logger.Fatal(err)
	}

	err = cp.Copy(config.SourceDir, config.DestinationDir)
	if err != nil {
		logger.Fatal(logger.StackTrace(err))
	}

	err = TerraformAction(ctx, config)
	if err != nil {
		logger.Fatal(logger.StackTrace(err))
	}
}
