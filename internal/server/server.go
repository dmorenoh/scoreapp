package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"scoreapp/internal/score"
)

type App struct {
	scoreService score.Service
}

func NewApp(scoreService score.Service) *App {
	return &App{
		scoreService: scoreService,
	}
}

func SetupRouter(app *App) *gin.Engine {
	r := gin.Default()
	r.POST("/user/:user_id/score", app.submit)
	r.GET("/ranking", app.getRanking)
	return r
}

func (a *App) submit(c *gin.Context) {

	var req SubmitScoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid json"})
		return
	}

	if req.Total != 0 && req.Score != "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid json"})
		return
	}

	if req.Total != 0 {
		if err := a.scoreService.SubmitAbsolute(c, req.User, req.Total); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "service error"})
			return
		}
	}

	if req.Score != "" {
		val, err := ScoreVariationValue(req.Score)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid json"})
			return
		}
		if err := a.scoreService.SubmitRelative(c, req.User, val); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "service error"})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "accepted"})
	return

}

func (a *App) getRanking(c *gin.Context) {
	value := c.Query("type")

	filter, err := BuildFilter(value)
	if err != nil {
		c.JSON(400, gin.H{"message": "invalid json"})
		return
	}

	result, err := a.scoreService.Find(c, filter)
	if err != nil {
		c.JSON(400, gin.H{"message": "invalid json"})
		return
	}

	c.JSON(200, result)
}

type SubmitScoreRequest struct {
	User  uint   `json:"user" binding:"required"`
	Total int    `json:"total"`
	Score string `json:"score"`
}
