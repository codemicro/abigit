package config

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"os"
	"strings"
)

func InitLogging() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	ll := zerolog.InfoLevel
	if Debug.Enabled {
		ll = zerolog.DebugLevel
	}

	log.Logger = zerolog.New(os.Stdout).Level(ll).With().Timestamp().Stack().Logger()

	if Debug.Enabled {
		log.Warn().Msg("debug mode enabled")
	}
}

var Debug = struct {
	Enabled bool
}{
	Enabled: asBool(fetchFromFile("debug.enabled")),
}

type httpConfig struct {
	Host          string
	Port          int
	ExternalURL   string
	secureCookies bool
}

var HTTP = httpConfig{
	Host:          asString(withDefault("http.host", "0.0.0.0")),
	Port:          asInt(withDefault("http.port", 8080)),
	ExternalURL:   strings.TrimSuffix(asString(required("http.externalURL")), "/"),
	secureCookies: asBool(withDefault("http.secureCookies", true)),
}

func (h httpConfig) SecureCookies() bool {
	if Debug.Enabled {
		return false
	}
	return h.secureCookies
}

var Database = struct {
	Filename string
}{
	Filename: asString(withDefault("database.filename", "abigit.sqlite3.db")),
}

var OIDC = struct {
	ClientID     string
	ClientSecret string
	Issuer       string
}{
	ClientID:     asString(required("oidc.clientID")),
	ClientSecret: asString(required("oidc.clientSecret")),
	Issuer:       asString(required("oidc.issuer")),
}

var Git = struct {
	RepositoriesPath string
}{
	RepositoriesPath: asString(required("git.repositoriesPath")),
}
