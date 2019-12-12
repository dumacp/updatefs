package main

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/dumacp/keycloak"
	"golang.org/x/oauth2"
)

const (
	keycloakurl = "https://fleet.nebulae.com.co/auth"
	clientid    = "devices"
	redirecturl = "https://fleet.nebulae.com.co/updatefs/*"
	realm       = "DEVICES"
)

var ctx context.Context
var serverkey keycloak.Keycloak

func keycloakinit() error {

	ctx = context.Background()
	config := &keycloak.ServerConfig{
		Url:         keycloakurl,
		ClientID:    clientid,
		RedirectUrl: redirecturl,
		Realm:       realm,
	}

	transport := loadLocalCert()
	client := &http.Client{
		Transport: transport,
	}
	ctx = keycloak.NewClientContext(ctx, client)

	var err error
	serverkey, err = keycloak.NewConfig(ctx, config)
	if err != nil {
		return err
	}
	return nil
}

func keycloakNewToken() (*oauth2.Token, error) {

	name, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	names := strings.Split(name, ".")

	token, err := serverkey.TokenRequest(ctx, names[0], names[0])
	if err != nil {
		return nil, err
	}

	return token, nil
}

func keycloakTokenSource(token *oauth2.Token) oauth2.TokenSource {
	tokenSource := serverkey.TokenSource(ctx, token)
	return tokenSource
}

func keycloakinfo(tokensource oauth2.TokenSource) (map[string]interface{}, error) {
	return serverkey.UserInfo(ctx, tokensource)
}
