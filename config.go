package main

import "context"

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
