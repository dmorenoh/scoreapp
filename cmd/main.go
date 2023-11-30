package main

import (
	"log"
	"scoreapp/internal/score"
	"scoreapp/internal/server"
)

func main() {
	scores := make(map[uint]*score.Entity)
	repository := score.NewInMemoryRepository(scores)
	scoreService := score.NewService(repository)
	app := server.NewApp(scoreService)
	router := server.SetupRouter(app)

	if err := router.Run(":8080"); err != nil {
		log.Fatal("error starting server")
	}
}
