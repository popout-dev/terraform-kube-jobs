package main

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-exec/tfexec"
	cp "github.com/otiai10/copy"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
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

type Logger struct{}

func NewLogger() Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	return Logger{}
}

func (logger *Logger) Info(message string) {
	log.Info().Msg(message)
}

func (logger *Logger) Error(err error) {
	log.Error().Stack().Err(err).Msg("")
}

func (logger *Logger) Fatal(err error) {
	log.Fatal().Stack().Err(err).Msg("")
}

func (logger *Logger) StackTrace(err error) error {
	return errors.Wrap(err, "")
}

type Key int

const (
	LoggerKey Key = iota
)

func AddLoggerToContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, LoggerKey, logger)
}

func LoggerFromContext(ctx context.Context) (Logger, error) {
	loggerValue := ctx.Value(LoggerKey)

	logger, ok := loggerValue.(Logger)
	if !ok {
		return Logger{}, errors.New("unexpected logger type")
	}

	return logger, nil
}

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

type EnvironmentVariableNotSetError struct {
	Variable string
}

func (e *EnvironmentVariableNotSetError) Error() string {
	return fmt.Sprintf("environment variable %s not set", e.Variable)
}

func GetEnvVar(name string) (string, error) {
	env := os.Getenv(name)
	if env == "" {
		return "", &EnvironmentVariableNotSetError{
			Variable: name,
		}
	}

	return env, nil
}

const (
	srcDirEnvVar    = "TF_SRC_DIR"
	destDirEnvVar   = "TF_DEST_DIR"
	tfVersionEnvVar = "TF_INSTALL_VERSION"
	tfActionEnvVar  = "TF_ACTION"
)

type Config struct {
	SourceDir        string
	DestinationDir   string
	TerraformVersion string
	TerraformAction  string
}

func GetConfig(ctx context.Context) (Config, error) {
	logger, err := LoggerFromContext(ctx)
	if err != nil {
		return Config{}, err
	}

	srcDir, err := GetEnvVar(srcDirEnvVar)
	if err != nil {
		return Config{}, logger.StackTrace(err)
	}

	destDir, err := GetEnvVar(destDirEnvVar)
	if err != nil {
		return Config{}, logger.StackTrace(err)
	}

	tfv, err := GetEnvVar(tfVersionEnvVar)
	if err != nil {
		return Config{}, logger.StackTrace(err)
	}

	tfa, err := GetEnvVar(tfActionEnvVar)
	if err != nil {
		return Config{}, logger.StackTrace(err)
	}

	conf := Config{
		SourceDir:        srcDir,
		DestinationDir:   destDir,
		TerraformVersion: tfv,
		TerraformAction:  tfa,
	}

	return conf, nil
}
