package config

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	os.Setenv("FOO_TEST", "BAR")
	defer os.Unsetenv("FOO_TEST")
	if v := GetEnv("FOO_TEST"); v != "BAR" {
		t.Fatalf("expected BAR got %q", v)
	}
}
