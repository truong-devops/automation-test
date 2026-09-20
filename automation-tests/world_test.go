package automationtests

import (
	"net/http"
	"time"
)

type World struct {
	baseURL    string
	httpClient *http.Client

	aliases map[string]string

	lastStatus int
	lastBody   []byte
}

func newWorld() *World {
	return &World{
		baseURL: "http://localhost:8081",
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		aliases: make(map[string]string),
	}
}
