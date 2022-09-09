package endpoints

import (
	"context"
	"github.com/codemicro/abigit/abigit/config"
	"github.com/codemicro/abigit/abigit/db"
	"github.com/codemicro/abigit/abigit/db/models"
	"github.com/codemicro/abigit/abigit/urls"
	"github.com/codemicro/abigit/abigit/util"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"time"
)

func (e *Endpoints) authOIDCOutbound(ctx *fiber.Ctx) error {
	if e.getSessionInformation(ctx) != nil {
		// they're already logged in, chief!
		return ctx.Redirect(urls.Index)
	}

	authState := e.authStates.New("")

	return ctx.Redirect(e.oauth2Config.AuthCodeURL(authState.Key))
}

func (e *Endpoints) authOIDCInbound(ctx *fiber.Ctx) error {
	providedState := ctx.Query("state")

	authState := e.authStates.Get(providedState)
	if authState == nil {
		return util.NewRichError(fiber.StatusBadRequest, "Invalid state", nil)
	} else {
		e.authStates.Delete(providedState)
	}

	c, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	oauth2Token, err := e.oauth2Config.Exchange(c, ctx.Query("code"))
	if err != nil {
		return errors.WithStack(err)
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return util.NewRichError(fiber.StatusBadRequest, "missing ID token", nil)
	}

	idToken, err := e.oidcVerifier.Verify(context.Background(), rawIDToken)
	if err != nil {
		return errors.WithStack(err)
	}

	var claims struct {
		Email   string `json:"email"`
		Subject string `json:"sub"`
	}
	if err := idToken.Claims(&claims); err != nil {
		return errors.WithStack(err)
	}

	user, err := e.db.GetUserByExternalID(claims.Subject)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			user = &models.User{
				EmailAddress: claims.Email,
				ExternalID:   claims.Subject,
			}
			if err := e.db.RegisterUser(user); err != nil {
				return errors.WithStack(err)
			}
		} else {
			return errors.WithStack(err)
		}
	}

	token := e.generateSessionToken(user.ID)

	ctx.Cookie(&fiber.Cookie{
		Name:     sessionCookieKey,
		Value:    token,
		Expires:  time.Now().UTC().Add(sessionDuration),
		Secure:   config.HTTP.SecureCookies(),
		HTTPOnly: true,
	})

	return ctx.Redirect(urls.Index)
}
