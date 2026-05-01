package launch

import (
	"os"
	"testing"

	"github.com/ollama/ollama/envconfig"
)

func TestSkipIntegrations(t *testing.T) {
	// Save original state
	origSpecs := integrationSpecs
	origOrder := launcherIntegrationOrder
	defer func() {
		integrationSpecs = origSpecs
		launcherIntegrationOrder = origOrder
		rebuildIntegrationSpecIndexes()
	}()

	// Set skip env var
	os.Setenv("OLLAMA_SKIP_INTEGRATIONS", "openclaw,claude")
	defer os.Unsetenv("OLLAMA_SKIP_INTEGRATIONS")

	// Trigger filtering (manually since init already ran)
	filterIntegrations()
	rebuildIntegrationSpecIndexes()

	// Check if openclaw is gone
	_, err := LookupIntegrationSpec("openclaw")
	if err == nil {
		t.Errorf("expected openclaw to be skipped, but it was found")
	}

	// Check if claude is gone
	_, err = LookupIntegrationSpec("claude")
	if err == nil {
		t.Errorf("expected claude to be skipped, but it was found")
	}

	// Check if others are still there
	_, err = LookupIntegrationSpec("copilot")
	if err != nil {
		t.Errorf("expected copilot to be found, but got error: %v", err)
	}

	// Check visible specs
	visible := ListVisibleIntegrationSpecs()
	for _, spec := range visible {
		if spec.Name == "openclaw" || spec.Name == "claude" {
			t.Errorf("skipped integration %s found in visible list", spec.Name)
		}
	}

	// Check launcher order
	for _, name := range launcherIntegrationOrder {
		if name == "openclaw" || name == "claude" {
			t.Errorf("skipped integration %s found in launcher order", name)
		}
	}
}

func TestSkipIntegrationsByAlias(t *testing.T) {
	// Save original state
	origSpecs := integrationSpecs
	origOrder := launcherIntegrationOrder
	defer func() {
		integrationSpecs = origSpecs
		launcherIntegrationOrder = origOrder
		rebuildIntegrationSpecIndexes()
	}()

	// Set skip env var using an alias
	os.Setenv("OLLAMA_SKIP_INTEGRATIONS", "clawdbot")
	defer os.Unsetenv("OLLAMA_SKIP_INTEGRATIONS")

	// Trigger filtering
	filterIntegrations()
	rebuildIntegrationSpecIndexes()

	// Check if openclaw is gone (since clawdbot is its alias)
	_, err := LookupIntegrationSpec("openclaw")
	if err == nil {
		t.Errorf("expected openclaw to be skipped by alias, but it was found")
	}
}
