package main

type Config struct {
	OauthToken   string   `json:"oauth_token"`
	Organization string   `json:"organization"`
	Repositories []string `json:"repositories"`
}
