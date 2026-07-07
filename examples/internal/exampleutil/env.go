// Package exampleutil holds tiny helpers shared by the runnable examples. It is
// internal to the examples tree and is not part of the SDK's public API.
package exampleutil

import (
	"log"
	"os"
	"strconv"
)

// MustEnv returns the value of an environment variable or exits if it is unset.
func MustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("environment variable %s is required", key)
	}
	return v
}

// MustEnvInt returns an integer environment variable or exits on error.
func MustEnvInt(key string) int64 {
	v := MustEnv(key)
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		log.Fatalf("environment variable %s must be an integer: %v", key, err)
	}
	return n
}
