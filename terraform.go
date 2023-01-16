package main

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-exec/tfexec"
)

func TerraformAction(ctx context.Context, config Config) error {
	logger, err := LoggerFromContext(ctx)
	if err != nil {
		return err
	}

	installer := &releases.ExactVersion{
		Product: product.Terraform,
		Version: version.Must(version.NewVersion(config.TerraformVersion)),
	}

	execPath, err := installer.Install(ctx)
	if err != nil {
		return logger.StackTrace(err)
	}

	tf, err := tfexec.NewTerraform(config.DestinationDir, execPath)
	if err != nil {
		return logger.StackTrace(err)
	}

	err = tf.Init(ctx, tfexec.Upgrade(true))
	if err != nil {
		return logger.StackTrace(err)
	}

	switch config.TerraformAction {
	case "APPLY":
		err = tf.Apply(ctx)
		if err != nil {
			return logger.StackTrace(err)
		}
	case "DESTROY":
		err = tf.Destroy(ctx)
		if err != nil {
			return logger.StackTrace(err)
		}
	default:
		return fmt.Errorf("unknown terraform command %s", config.TerraformAction)
	}

	return nil
}
