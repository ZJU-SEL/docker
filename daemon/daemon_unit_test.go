package daemon

import (
	"testing"

	"github.com/docker/docker/runconfig"
)

func TestParseSecurityOpt(t *testing.T) {
	container := &Container{}
	config := &runconfig.HostConfig{}

	// test apparmor
	config.SecurityOpt = []string{"apparmor:test_profile"}
	if err := parseSecurityOpt(container, config); err != nil {
		t.Fatalf("Unexpected parseSecurityOpt error: %v", err)
	}
	if container.AppArmorProfile != "test_profile" {
		t.Fatalf("Unexpected AppArmorProfile, expected: \"test_profile\", got %q", container.AppArmorProfile)
	}

	// test valid label
	config.SecurityOpt = []string{"label:user:USER"}
	if err := parseSecurityOpt(container, config); err != nil {
		t.Fatalf("Unexpected parseSecurityOpt error: %v", err)
	}

	// test invalid label
	config.SecurityOpt = []string{"label"}
	if err := parseSecurityOpt(container, config); err == nil {
		t.Fatal("Expected parseSecurityOpt error, got nil")
	}

	// test invalid opt
	config.SecurityOpt = []string{"test"}
	if err := parseSecurityOpt(container, config); err == nil {
		t.Fatal("Expected parseSecurityOpt error, got nil")
	}
}

//only test whether path is absolute
func TestParseVolumeMountConfig(t *testing.T) {
	container := &Container{}
	container.hostConfig = &runconfig.HostConfig{}
	container.Config = &runconfig.Config{Volumes: make(map[string]struct{})}

	container.hostConfig.Binds = []string{"data:/data"}
	if _, err := container.parseVolumeMountConfig(); err == nil {
		t.Fatal("Expected parseVolumeMountConfig error, got nil")
	}

	container.hostConfig.Binds = []string{"/data:data"}
	if _, err := container.parseVolumeMountConfig(); err == nil {
		t.Fatal("Expected parseVolumeMountConfig error, got nil")
	}

	container.Config.Volumes["data"] = struct{}{}
	if _, err := container.parseVolumeMountConfig(); err == nil {
		t.Fatal("Expected parseVolumeMountConfig error, got nil")
	}
	t.Log("parseVolumeMountConfig test passed")
}
