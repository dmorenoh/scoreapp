package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"scoreapp/internal/score"
	"scoreapp/internal/server"
	"sync"
	"testing"
	"time"
)

func Test_Client(t *testing.T) {
	scores := make(map[uint]*score.Entity)
	repository := score.NewInMemoryRepository(scores)
	scoreService := score.NewService(repository)
	app := server.NewApp(scoreService)
	router := server.SetupRouter(app)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		if err := router.Run(":8081"); err != nil {
			log.Fatal("error starting server")
		}
	}()

	client := http.Client{Timeout: 30 * time.Second}
	buffer := bytes.NewBuffer([]byte(`{"user": 123, "score": "+100"}`))
	requestURL := "http://localhost:8081/user/123/score"
	_, err := client.Post(requestURL, "application/json", buffer)
	assert.NoError(t, err)

	wg.Done()
	wg.Wait()
}
