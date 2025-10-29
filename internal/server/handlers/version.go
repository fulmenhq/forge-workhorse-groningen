package handlers

import (
	"encoding/json"
	"net/http"
	"runtime"

	"github.com/fulmenhq/gofulmen/crucible"
)

// AppVersion is injected from main via SetVersionInfo
var (
	AppVersion   = "dev"
	AppCommit    = "unknown"
	AppBuildDate = "unknown"
)

// SetVersionInfo sets the version information for the handler
func SetVersionInfo(version, commit, buildDate string) {
	AppVersion = version
	AppCommit = commit
	AppBuildDate = buildDate
}

// VersionResponse represents the version information response
type VersionResponse struct {
	App     AppInfo     `json:"app"`
	SSOT    SSOTInfo    `json:"ssot"`
	Runtime RuntimeInfo `json:"runtime"`
}

// AppInfo contains application version details
type AppInfo struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildDate string `json:"buildDate"`
}

// SSOTInfo contains SSOT version information
type SSOTInfo struct {
	Gofulmen string `json:"gofulmen"`
	Crucible string `json:"crucible"`
}

// RuntimeInfo contains runtime environment information
type RuntimeInfo struct {
	Go string `json:"go"`
}

// VersionHandler handles version information requests
func VersionHandler(w http.ResponseWriter, r *http.Request) {
	version := crucible.GetVersion()

	response := VersionResponse{
		App: AppInfo{
			Name:      "groningen",
			Version:   AppVersion,
			Commit:    AppCommit,
			BuildDate: AppBuildDate,
		},
		SSOT: SSOTInfo{
			Gofulmen: version.Gofulmen,
			Crucible: version.Crucible,
		},
		Runtime: RuntimeInfo{
			Go: runtime.Version(),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
