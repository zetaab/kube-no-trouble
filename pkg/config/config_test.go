package config

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"k8s.io/client-go/util/homedir"
)

func TestValidLogLevelFromFlags(t *testing.T) {
	oldArgs := os.Args[1]
	defer func() { os.Args[1] = oldArgs }()

	var validLevels = []string{"trace", "debug", "info", "warn", "error", "fatal", "panic", "", "disabled"}
	for i, level := range validLevels {
		// reset for testing
		pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

		os.Args[1] = "--log-level=" + level
		config, err := NewFromFlags()

		if err != nil {
			t.Errorf("Flags parsing failed %s", err)
		}

		expected := ZeroLogLevel(i - 1)
		actual := config.LogLevel

		if actual != expected {
			t.Errorf("Config not parsed correctly: %s \nactual %d, expected %d", level, actual, expected)
		}
	}
}

func TestInvalidLogLevelFromFlags(t *testing.T) {
	var testLevel ZeroLogLevel

	if err := testLevel.Set("bad"); err == nil {
		t.Errorf("Should not parse invalid flag")
	}
}

func TestNewFromFlags(t *testing.T) {
	// reset for testing
	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

	config, err := NewFromFlags()

	if err != nil {
		t.Errorf("Flags parsing failed %s", err)
	}

	if !config.Cluster && config.Output != "text" {
		t.Errorf("Config not parsed correctly")
	}
}

func TestNewFromFlagsKubeconfigEnv(t *testing.T) {
	// reset for testing
	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

	testVal := "my-file.conf"

	err := os.Setenv("KUBECONFIG", testVal)
	if err != nil {
		t.Errorf("failed to set env variable: %s", err)
	}

	config, err := NewFromFlags()

	if err != nil {
		t.Errorf("Formatting flags failed %s", err)
	}

	if config.Kubeconfig != testVal {
		t.Errorf("kubeconfig option not loaded correctly from ebv variable, expected: %s, got: %s", testVal, config.Kubeconfig)
	}
}

func TestNewFromFlagsKubeconfigHome(t *testing.T) {
	// reset for testing
	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

	expected := homedir.HomeDir() + "/.kube/config"
	err := os.Unsetenv("KUBECONFIG")
	if err != nil {
		t.Errorf("failed to unset KUBECONFIG env variable: %s", err)
	}

	config, err := NewFromFlags()

	if err != nil {
		t.Errorf("Formatting flags failed %s", err)
	}

	if config.Kubeconfig != expected {
		t.Errorf("kubeconfig option not set to expected default, expected: %s, got: %s", expected, config.Kubeconfig)
	}
}

func TestEnvOrStringVariable(t *testing.T) {
	err := os.Setenv("FOO", "1")
	if err != nil {
		t.Errorf("failed to set env variable: %e", err)
	}

	i := envOrString("FOO", "default")
	if i != "1" {
		t.Errorf("Expected to get env variable, got %s instead", i)
	}
}

func TestEnvOrStringDefault(t *testing.T) {
	err := os.Unsetenv("FOO")
	if err != nil {
		t.Errorf("failed to unset env variable: %e", err)
	}

	i := envOrString("FOO", "default")
	if i != "default" {
		t.Errorf("Expected to get default string, got %s instead", i)
	}
}

func TestValidateAdditionalResources(t *testing.T) {
	resources := []string{
		"Test.v1.example.com",
		"ManagedCertificates.v1.networking.gke.io",
		"ManagedCertificates.networking.gke.io",
	}

	err := validateAdditionalResources(resources)

	if err != nil {
		t.Errorf("expected resources %s to pass validation: %s", resources, err)
	}
}

func TestValidateAdditionalResourcesFail(t *testing.T) {
	testCases := [][]string{
		[]string{"abcdef"},
		[]string{""},
		[]string{"test.v1.com"},
	}

	for _, tc := range testCases {
		err := validateAdditionalResources(tc)

		if err == nil {
			t.Errorf("expected resources %s to fail validation: %s", tc, err)
		}
	}
}

func TestValidateTargetVersion(t *testing.T) {
	validVersions := []string{
		"", "1.16", "1.16.3", "1.2.3",
	}
	invalidVersions := []string{
		"1", "v12.3.4", "v1.2", "1.blah", "nope",
	}

	for _, v := range validVersions {
		err := validateTargetVersion(v)

		if err != nil {
			t.Errorf("Expected %s to fail validation, it failed instead with: %s", v, err)
		}
	}

	for _, v := range invalidVersions {
		err := validateTargetVersion(v)

		if err == nil {
			t.Errorf("Expected %s to fail validation, it succeeded instead", v)
		}
	}
}
