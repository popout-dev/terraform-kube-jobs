package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"testing"

	cp "github.com/otiai10/copy"
	"github.com/stretchr/testify/assert"
)

func TestTerraformAction(t *testing.T) {
	outDir := filepath.Join(CurrentDir(), "out")
	testDir := filepath.Join(CurrentDir(), "test")

	tests := map[string]struct {
		Config   Config
		FileName string
	}{
		"runs an apply": {
			Config: Config{
				SourceDir:        testDir,
				DestinationDir:   outDir,
				TerraformVersion: "1.3.7",
				TerraformAction:  "APPLY",
			},
			FileName: "test.txt",
		},
		"runs a destroy": {
			Config: Config{
				SourceDir:        testDir,
				DestinationDir:   outDir,
				TerraformVersion: "1.3.7",
				TerraformAction:  "DESTROY",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := AddLoggerToContext(context.Background(), NewLogger())

			if test.Config.TerraformAction == "APPLY" {
				err := cp.Copy(test.Config.SourceDir, test.Config.DestinationDir)
				if err != nil {
					t.Fatal(err)
				}
			}

			err := TerraformAction(ctx, test.Config)
			if err != nil {
				t.Fatal(err)
			}

			if test.Config.TerraformAction == "APPLY" {
				assert.FileExists(t, filepath.Join(test.Config.DestinationDir, test.FileName))
			}
		})
	}

	os.RemoveAll(outDir)
}

func CurrentDir() string {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return currentDir
}
