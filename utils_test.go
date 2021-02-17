package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetEnvVar(t *testing.T) {
	os.Setenv("FOO", "BAR")
	assert.Equal(t, GetEnvVar("FOO", "BAZ"), "BAR")
	assert.Equal(t, GetEnvVar("QUX", "BAZ"), "BAZ")
}
