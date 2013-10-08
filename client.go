package main

import (
	"encoding/base64"
	"log"
	"net/http"

	"code.google.com/p/goauth2/oauth"
	"code.google.com/p/goauth2/oauth/jwt"

	calendar "code.google.com/p/google-api-go-client/calendar/v3"
)

const scope = "https://www.googleapis.com/auth/calendar.readonly"
const authURL = "https://accounts.google.com/o/oauth2/auth"
const tokenURL = "https://accounts.google.com/o/oauth2/token"

type ApiClient struct {
	ClientId   string
	EncodedKey string
	Token      *oauth.Token
}

func (c ApiClient) GetToken() *oauth.Token {

	keyBytes, err := base64.StdEncoding.DecodeString(c.EncodedKey)
	if err != nil {
		log.Fatal("Error decoding private key:", err)
	}

	t := jwt.NewToken(c.ClientId, scope, keyBytes)
	t.ClaimSet.Aud = tokenURL

	log.Print("Requesting new access token.\n")
	httpClient := &http.Client{}
	token, err := t.Assert(httpClient)
	if err != nil {
		log.Fatal("assertion error:", err)
	}

	log.Printf("New access token acquired.\n")
	return token
}

func (c ApiClient) Client() *http.Client {
	config := &oauth.Config{
		ClientId:     c.ClientId,
		ClientSecret: "notasecret",
		Scope:        scope,
		AuthURL:      authURL,
		TokenURL:     tokenURL,
	}

	transport := &oauth.Transport{
		Token:     c.Token,
		Config:    config,
		Transport: http.DefaultTransport,
	}

	return transport.Client()
}

func (c ApiClient) Api() *calendar.Service {
	client := c.Client()

	svc, _ := calendar.New(client)
	return svc
}
