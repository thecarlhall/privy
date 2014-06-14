package main

type Config struct {
	Organization string   `json:"organization"`
	OauthToken   string   `json:"oauth_token"`
	Repositories []string `json:"repositories"`
}
