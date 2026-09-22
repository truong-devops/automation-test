package automationtests

import (
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

type World struct {
	baseURL     string
	httpClient  *http.Client
	mongoClient *mongo.Client

	aliases         map[string]string
	requestFixtures map[string][]byte

	lastStatus int
	lastBody   []byte
}

func newWorld() *World {
	return &World{
		baseURL:    "http://localhost:8081",
		httpClient: &http.Client{Timeout: 5 * time.Second},

		aliases:         make(map[string]string),
		requestFixtures: make(map[string][]byte),
	}
}

func (w *World) resetScenarioState() {
	w.aliases = make(map[string]string)
	w.requestFixtures = make(map[string][]byte)

	w.lastStatus = 0
	w.lastBody = nil
}
