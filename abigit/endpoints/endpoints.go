package endpoints

import (
	"context"
	"encoding/hex"
	goalone "github.com/bwmarrin/go-alone"
	"github.com/bwmarrin/snowflake"
	"github.com/codemicro/abigit/abigit/config"
	"github.com/codemicro/abigit/abigit/db"
	"github.com/codemicro/abigit/abigit/static"
	"github.com/codemicro/abigit/abigit/urls"
	"github.com/codemicro/abigit/abigit/util"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"time"
)

const (
	sessionCookieKey = "burntpocket-session"
	sessionDuration  = time.Hour * 24 * 15
)

type Endpoints struct {
	db            *db.DB
	sessionSigner *goalone.Sword

	oidcProvider *oidc.Provider
	oidcVerifier *oidc.IDTokenVerifier
	oauth2Config *oauth2.Config

	authStates *authStateManager
}

func New(dbi *db.DB) (*Endpoints, error) {
	key, err := dbi.FetchSigningKey()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	provider, err := oidc.NewProvider(ctx, config.OIDC.Issuer)
	if err != nil {
		return nil, errors.WithMessage(err, "cannot load OIDC information")
	}

	oauth2Config := &oauth2.Config{
		ClientID:     config.OIDC.ClientID,
		ClientSecret: config.OIDC.ClientSecret,
		RedirectURL:  config.HTTP.ExternalURL + urls.AuthOIDCInbound,

		// Discovery returns the OAuth2 endpoints.
		Endpoint: provider.Endpoint(),

		// "openid" is a required scope for OpenID Connect flows.
		Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
	}

	return &Endpoints{
		db:            dbi,
		sessionSigner: goalone.New(key.Key, goalone.Timestamp),
		oidcProvider:  provider,
		oidcVerifier:  provider.Verifier(&oidc.Config{ClientID: config.OIDC.ClientID}),
		oauth2Config:  oauth2Config,
		authStates:    newAuthStateManager(),
	}, nil
}

func (e *Endpoints) SetupApp() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler:          util.JSONErrorHandler,
		DisableStartupMessage: !config.Debug.Enabled,
	})

	app.Get(urls.Index, e.Index)

	app.Get(urls.AuthOIDCOutbound, e.authOIDCOutbound)
	app.Get(urls.AuthOIDCInbound, e.authOIDCInbound)

	app.Use("/", static.NewHandler())

	return app
}

func (e *Endpoints) generateSessionToken(userID snowflake.ID) string {
	return hex.EncodeToString(e.sessionSigner.Sign(userID.Bytes()))
}

type sessionInfo struct {
	ID        snowflake.ID
	CreatedAt time.Time
}

func (e *Endpoints) getSessionInformation(ctx *fiber.Ctx) *sessionInfo {
	debugLog := log.With().Str("location", "getSessionInformation").Logger()

	cookieContent := ctx.Cookies(sessionCookieKey)
	if cookieContent == "" {
		debugLog.Debug().Msg("session cookie not set")
		return nil
	}

	decodedCookie, err := hex.DecodeString(cookieContent)
	if err != nil {
		debugLog.Debug().Err(err).Msg("could not hex decode session cookie")
		return nil
	}

	if _, err := e.sessionSigner.Unsign(decodedCookie); err != nil {
		debugLog.Debug().Err(err).Msg("invalid signature")
		return nil
	}

	parsedToken := e.sessionSigner.Parse(decodedCookie)

	parsedUserID, err := snowflake.ParseBytes(parsedToken.Payload)
	if err != nil {
		log.Error().Err(err).Msg("signed session token doesn't contain a valid user ID")
		return nil
	}

	si := &sessionInfo{
		ID:        parsedUserID,
		CreatedAt: parsedToken.Timestamp,
	}

	if time.Now().UTC().Sub(si.CreatedAt) < 0 {
		debugLog.Debug().Msg("session token expired")
		return nil
	}

	return si
}
