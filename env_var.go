package main

import (
	"fmt"
	"os"
)

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
