package main

type Config struct {
	Organization string   `json:"organization"`
	OauthToken   string   `json:"oauth_token"`
	Projects     []string `json:"projects"`
}
