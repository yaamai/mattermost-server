// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package web

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc"
	"github.com/mattermost/mattermost-server/v5/einterfaces"
	"github.com/mattermost/mattermost-server/v5/mlog"
	oauthoidc "github.com/mattermost/mattermost-server/v5/model/oidc"
	"golang.org/x/oauth2"
)

var oauth2Config oauth2.Config
var verifier *oidc.IDTokenVerifier

func testOidcCompleteHandler(c *Context, w http.ResponseWriter, r *http.Request) {
	mlog.Info("OK")

	// Verify state and errors.
	ctx := context.TODO()
	oauth2Token, err := oauth2Config.Exchange(ctx, r.URL.Query().Get("code"))
	if err != nil {
		// handle error
		mlog.Error(err.Error())
	}

	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		// handle missing token
		mlog.Error("failed to extract id token")
	}

	// Parse and verify ID Token payload.
	ctx = context.TODO()
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		mlog.Error(err.Error())
	}

	// Extract custom claims
	var claims struct {
		Email string `json:"email"`
	}
	if err := idToken.Claims(&claims); err != nil {
		mlog.Error(err.Error())
	}

	mlog.Info("Claims:")
	mlog.Info(claims.Email)

	provider := &oauthoidc.OidcProvider{}
	einterfaces.RegisterOauthProvider("oidc", provider)
	jsonBytes, err := json.Marshal(map[string]string{"sub": idToken.Subject, "email": claims.Email})
	mlog.Info(string(jsonBytes))
	user, appErr := c.App.LoginByOAuth("oidc", bytes.NewReader(jsonBytes), "")
	if appErr != nil {

		mlog.Error("some error")
		appErr.Translate(c.App.T)
		mlog.Error(appErr.Error())
	}
	// mlog.Error(err)
	appErr = c.App.DoLogin(w, r, user, "")
	if appErr != nil {
		appErr.Translate(c.App.T)
		c.Err = appErr
		return
	}

	c.App.AttachSessionCookies(w, r)

	redirectUrl := c.GetSiteURLHeader()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.Redirect(w, r, redirectUrl, http.StatusTemporaryRedirect)

	ReturnStatusOK(w)
}

func testOidcLoginHandler(c *Context, w http.ResponseWriter, r *http.Request) {

	issuerUrl := *c.App.Config().OidcSettings.IssuerUrl
	clientId := *c.App.Config().OidcSettings.ClientId
	clientSecret := *c.App.Config().OidcSettings.ClientSecret
	scopes := *c.App.Config().OidcSettings.Scopes

	ctx := context.TODO()
	provider, err := oidc.NewProvider(ctx, issuerUrl)
	if err != nil {
		mlog.Error(err.Error())
	}

	verifier = provider.Verifier(&oidc.Config{ClientID: clientId})

	oauth2Config = oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  c.GetSiteURLHeader() + "/oidc/complete",
		Endpoint:     provider.Endpoint(),
		// "openid" scope set by config model
		Scopes: strings.Split(scopes, ","),
	}

	http.Redirect(w, r, oauth2Config.AuthCodeURL("asdfasdfowiearwuo"), http.StatusFound)
}
