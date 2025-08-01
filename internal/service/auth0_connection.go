package service

import (
	"context"
	"github.com/go-resty/resty/v2"
)

type Auth0Credentials struct {
	Domain       string
	ClientID     string
	ClientSecret string
	Audience     string
}

func GetManagementToken(ctx context.Context, creds Auth0Credentials) (string, error) {
	client := resty.New()

	var resp struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
	}

	_, err := client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]string{
			"grant_type":    "client_credentials",
			"client_id":     creds.ClientID,
			"client_secret": creds.ClientSecret,
			"audience":      creds.Audience,
		}).
		SetResult(&resp).
		Post("https://" + creds.Domain + "/oauth/token")

	if err != nil {
		return "", err
	}

	return resp.AccessToken, nil
}
