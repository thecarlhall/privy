package main

type Config struct {
	Debug        bool         `json:"debug,omitempty"`
	OauthToken   string       `json:"oauth_token"`
	Repositories []Repository `json:"repositories"`
}

type Repository struct {
	Organization string   `json:"organization"`
	Projects     []string `json:"projects"`
}
