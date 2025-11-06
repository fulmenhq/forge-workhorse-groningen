package cmd

import (
	"context"
	"testing"

	"github.com/fulmenhq/gofulmen/appidentity"
)

func TestAppIdentityLoading(t *testing.T) {
	t.Run("load app identity from .fulmen/app.yaml", func(t *testing.T) {
		// Load app identity the same way the application does
		ctx := context.Background()
		identity, err := appidentity.Get(ctx)

		// Should load successfully
		if err != nil {
			t.Fatalf("Failed to load app identity: %v", err)
		}

		if identity == nil {
			t.Fatal("App identity is nil")
		}

		// Log the full identity for debugging
		t.Logf("Loaded identity: %+v", identity)

		// Check all expected fields are populated
		expectedFields := map[string]string{
			"Vendor":     identity.Vendor,
			"BinaryName": identity.BinaryName,
			"EnvPrefix":  identity.EnvPrefix,
			"ConfigName": identity.ConfigName,
		}

		for fieldName, value := range expectedFields {
			if value == "" {
				t.Errorf("App identity field %s is empty (expected: non-empty)", fieldName)
			} else {
				t.Logf("✅ %s = '%s'", fieldName, value)
			}
		}

		// Specific assertions for template values
		if identity.Vendor != "fulmen" {
			t.Errorf("Expected vendor 'fulmen', got '%s'", identity.Vendor)
		}

		if identity.BinaryName != "groningen" {
			t.Errorf("Expected binary_name 'groningen', got '%s'", identity.BinaryName)
		}

		if identity.EnvPrefix != "GRONINGEN_" {
			t.Errorf("Expected env_prefix 'GRONINGEN_', got '%s'", identity.EnvPrefix)
		}

		if identity.ConfigName != "groningen" {
			t.Errorf("Expected config_name 'groningen', got '%s'", identity.ConfigName)
		}
	})
}
