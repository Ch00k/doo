package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnvVar(t *testing.T) {
	os.Setenv("FOO", "BAR")
	assert.Equal(t, GetEnvVar("FOO", "BAZ"), "BAR")
	assert.Equal(t, GetEnvVar("QUX", "BAZ"), "BAZ")
}
