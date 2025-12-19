package appid

import (
	"context"

	"github.com/fulmenhq/gofulmen/appidentity"

	appidentityassets "github.com/fulmenhq/forge-workhorse-groningen/internal/assets/appidentity"
)

func init() {
	// Best-effort registration.
	//
	// Explicit identity overrides remain authoritative (Options.ExplicitPath and
	// FULMEN_APP_IDENTITY_PATH). Embedded identity provides standalone-binary
	// behavior when no external `.fulmen/app.yaml` can be found.
	_ = appidentity.RegisterEmbeddedIdentityYAML(appidentityassets.YAML)
}

func Get(ctx context.Context) (*appidentity.Identity, error) {
	return appidentity.Get(ctx)
}
