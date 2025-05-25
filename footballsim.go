package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/teams", getTeams)

	fmt.Println("Starting server at http://localhost:8079")
	router.Run("localhost:8079")
}

type team struct {
	ID       string `json:"id"`
	Name     string `json:"name"`      //team name
	Points   int    `json:"points"`    //PTS
	Played   int    `json:"played"`    //P
	Win      int    `json:"win"`       //W
	Drawn    int    `json:"drawn"`     //D
	Lost     int    `json:"lost"`      //L
	GoalDiff int    `json:"goal_diff"` //GD
}

var teams = []team{
	{ID: "1", Name: "Chelsea", Points: 0, Played: 0, Win: 0, Drawn: 0, Lost: 0, GoalDiff: 0},
	{ID: "2", Name: "Arsenal", Points: 0, Played: 0, Win: 0, Drawn: 0, Lost: 0, GoalDiff: 0},
	{ID: "3", Name: "Man City", Points: 0, Played: 0, Win: 0, Drawn: 0, Lost: 0, GoalDiff: 0},
	{ID: "4", Name: "Liverpool", Points: 0, Played: 0, Win: 0, Drawn: 0, Lost: 0, GoalDiff: 0},
}

func getTeams(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, teams)
}
