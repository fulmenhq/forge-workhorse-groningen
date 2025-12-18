package appid

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fulmenhq/gofulmen/appidentity"
)

var executableDir = func() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	if exe == "" {
		return "", fmt.Errorf("executable path is empty")
	}
	return filepath.Dir(exe), nil
}

// Get loads Fulmen app identity with a portable fallback.
//
// Behavior:
// - If FULMEN_APP_IDENTITY_PATH is set, it is used (gofulmen behavior).
// - Otherwise, it searches upward from the current working directory.
// - If not found, it searches upward from the executable directory.
//
// The returned identity is cached by gofulmen via GetWithOptions(ExplicitPath=...).
func Get(ctx context.Context) (*appidentity.Identity, error) {
	identityPath, err := findIdentityPath()
	if err != nil {
		return nil, err
	}

	return appidentity.GetWithOptions(ctx, appidentity.Options{ExplicitPath: identityPath})
}

func findIdentityPath() (string, error) {
	// Preserve gofulmen’s env var override semantics.
	if envPath := os.Getenv(appidentity.EnvIdentityPath); envPath != "" {
		absPath, err := filepath.Abs(envPath)
		if err != nil {
			return "", fmt.Errorf("invalid %s path: %w", appidentity.EnvIdentityPath, err)
		}
		if _, err := os.Stat(absPath); err == nil {
			return absPath, nil
		}
		return "", &appidentity.NotFoundError{SearchedPaths: []string{absPath + " (from " + appidentity.EnvIdentityPath + ")"}}
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}
	absCWD, err := filepath.Abs(cwd)
	if err != nil {
		return "", fmt.Errorf("invalid current directory: %w", err)
	}

	searched := make([]string, 0, appidentity.MaxSearchDepth*2)
	if path, paths := searchUp(absCWD); path != "" {
		return path, nil
	} else {
		searched = append(searched, paths...)
	}

	exeDir, err := executableDir()
	if err == nil {
		absExeDir, err := filepath.Abs(exeDir)
		if err == nil && absExeDir != "" && absExeDir != absCWD {
			if path, paths := searchUp(absExeDir); path != "" {
				return path, nil
			} else {
				for _, p := range paths {
					searched = append(searched, p+" (fallback: executable dir)")
				}
			}
		}
	}

	return "", &appidentity.NotFoundError{StartDir: absCWD, SearchedPaths: searched}
}

func searchUp(startDir string) (found string, searchedPaths []string) {
	searchedPaths = make([]string, 0, appidentity.MaxSearchDepth)
	currentDir := startDir

	for depth := 0; depth < appidentity.MaxSearchDepth; depth++ {
		candidate := filepath.Join(currentDir, appidentity.DefaultIdentityPath)
		searchedPaths = append(searchedPaths, candidate)

		if _, err := os.Stat(candidate); err == nil {
			return candidate, searchedPaths
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			break
		}
		currentDir = parentDir
	}

	return "", searchedPaths
}
