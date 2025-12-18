package appid

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/fulmenhq/gofulmen/appidentity"
)

const testIdentityYAML = "" +
	"app:\n" +
	"  vendor: acme\n" +
	"  binary_name: testapp\n" +
	"  env_prefix: TESTAPP_\n" +
	"  config_name: testapp\n"

func TestGet_UsesCWDSearchFirst(t *testing.T) {
	appidentity.Reset()
	defer appidentity.Reset()

	oldExecutableDir := executableDir
	executableDir = func() (string, error) {
		return "", os.ErrNotExist
	}
	defer func() { executableDir = oldExecutableDir }()

	root := t.TempDir()
	identityDir := filepath.Join(root, appidentity.DefaultIdentityDir)
	if err := os.MkdirAll(identityDir, 0o755); err != nil {
		t.Fatalf("mkdir identity dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(identityDir, appidentity.DefaultIdentityFilename), []byte(testIdentityYAML), 0o644); err != nil {
		t.Fatalf("write identity file: %v", err)
	}

	nested := filepath.Join(root, "a", "b")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatalf("mkdir nested: %v", err)
	}

	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	defer func() { _ = os.Chdir(oldWD) }()
	if err := os.Chdir(nested); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	identity, err := Get(context.Background())
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if identity.BinaryName != "testapp" {
		t.Fatalf("unexpected BinaryName: %q", identity.BinaryName)
	}
}

func TestGet_FallsBackToExecutableDir(t *testing.T) {
	appidentity.Reset()
	defer appidentity.Reset()

	root := t.TempDir()
	identityDir := filepath.Join(root, appidentity.DefaultIdentityDir)
	if err := os.MkdirAll(identityDir, 0o755); err != nil {
		t.Fatalf("mkdir identity dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(identityDir, appidentity.DefaultIdentityFilename), []byte(testIdentityYAML), 0o644); err != nil {
		t.Fatalf("write identity file: %v", err)
	}

	fakeExeDir := filepath.Join(root, "bin")
	if err := os.MkdirAll(fakeExeDir, 0o755); err != nil {
		t.Fatalf("mkdir fake exe dir: %v", err)
	}

	oldExecutableDir := executableDir
	executableDir = func() (string, error) {
		return fakeExeDir, nil
	}
	defer func() { executableDir = oldExecutableDir }()

	oldEnv := os.Getenv(appidentity.EnvIdentityPath)
	_ = os.Unsetenv(appidentity.EnvIdentityPath)
	defer func() {
		if oldEnv == "" {
			_ = os.Unsetenv(appidentity.EnvIdentityPath)
			return
		}
		_ = os.Setenv(appidentity.EnvIdentityPath, oldEnv)
	}()

	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	defer func() { _ = os.Chdir(oldWD) }()

	outside := t.TempDir()
	if err := os.Chdir(outside); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	identity, err := Get(context.Background())
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if identity.BinaryName != "testapp" {
		t.Fatalf("unexpected BinaryName: %q", identity.BinaryName)
	}
}
