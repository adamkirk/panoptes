package api

import (
	"log/slog"
	"net/http"

	"github.com/adamkirk/panoptes/internal/domain/users"
	"github.com/danielgtaylor/huma/v2"
)

type AuthConfig interface {
	GetMasterToken() string
}

type AuthMiddleware struct {
	api huma.API
}

type AuthRepo interface {
	ByID(id string) (*users.AccessToken, error)
}

type TokenVerifier interface {
	HashMatches(hash string, val string) (bool)
}

func NewAuthMiddleware(api huma.API, repo AuthRepo, verifier TokenVerifier) func(ctx huma.Context, next func(huma.Context)) {
	return func (ctx huma.Context, next func(huma.Context)) {
		authRequired := false

		var neededScopes []string
		for _, opScheme := range ctx.Operation().Security {
			var ok bool

			if neededScopes, ok = opScheme["scopes"]; ok {
				authRequired = true
				break
			}
		}

		if ! authRequired {
			next(ctx)
			return
		}

		key := ctx.Header("X-Access-Key-ID")
		token := ctx.Header("X-Access-Key-Token")

		if key != "" && token != "" {
			accessToken, err := repo.ByID(key)

			if err != nil {
				slog.Error("failed to get auth token", "error", err)
				huma.WriteErr(api, ctx, http.StatusInternalServerError, "failed to verify access token")
				return
			}

			if ! verifier.HashMatches(accessToken.SecretHash, token) {
				huma.WriteErr(api, ctx, http.StatusUnauthorized, "Not authorized to perform this action.")
				return
			}

			if accessToken.User.Can(neededScopes) {
				next(ctx)
			}
			
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Not authorized to perform this action.")
			
			return
		}

		if len(neededScopes) > 0 {
			// TOOD: check the scopes, this way is a bit bonkers
		}

		huma.WriteErr(api, ctx, http.StatusUnauthorized, "Not auth mechanism found")
	}
}