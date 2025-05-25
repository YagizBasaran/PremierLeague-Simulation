package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"math"
	"math/rand"
	"time"
)

func main() {
	router := gin.Default()

	router.GET("/teams", getTeams)
	router.POST("/teams/:id/win", winnerTeamAPI)
	router.POST("/teams/:id/draw", drawTeamAPI)
	router.POST("/teams/:id/loss", loserTeamAPI)

	simulateMatch(&teams[0], &teams[1])
	simulateMatch(&teams[0], &teams[2])
	simulateMatch(&teams[0], &teams[3])
	simulateMatch(&teams[1], &teams[0])
	simulateMatch(&teams[2], &teams[0])
	simulateMatch(&teams[3], &teams[0])

	simulateMatch(&teams[1], &teams[2])
	simulateMatch(&teams[1], &teams[3])
	simulateMatch(&teams[2], &teams[1])
	simulateMatch(&teams[3], &teams[1])

	simulateMatch(&teams[2], &teams[3])
	simulateMatch(&teams[3], &teams[2])

	fmt.Println("Starting server at http://localhost:8079")
	router.Run("localhost:8079")
}

type team struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`          //team name
	Points       int     `json:"points"`        //PTS
	Played       int     `json:"played"`        //P
	Win          int     `json:"win"`           //W
	Drawn        int     `json:"drawn"`         //D
	Lost         int     `json:"lost"`          //L
	GoalsFor     int     `json:"goals_for"`     //GF
	GoalsAgainst int     `json:"goals_against"` //GA
	GoalDiff     int     `json:"goal_diff"`     //GD
	Elo          float64 `json:"elo"`           //ELO rating
	Tilt         float64 `json:"tilt"`          //TILT rating for teams play-style: high = offensive, low = defensive
}

var teams = []team{
	{ID: "1", Name: "Chelsea", Points: 0, Played: 0, Win: 0, Drawn: 0, Lost: 0, GoalsFor: 0, GoalsAgainst: 0, GoalDiff: 0, Elo: 1200, Tilt: 1.0},
	{ID: "2", Name: "Arsenal", Points: 0, Played: 0, Win: 0, Drawn: 0, Lost: 0, GoalsFor: 0, GoalsAgainst: 0, GoalDiff: 0, Elo: 1200, Tilt: 1.0},
	{ID: "3", Name: "Man City", Points: 0, Played: 0, Win: 0, Drawn: 0, Lost: 0, GoalsFor: 0, GoalsAgainst: 0, GoalDiff: 0, Elo: 1200, Tilt: 1.0},
	{ID: "4", Name: "Liverpool", Points: 0, Played: 0, Win: 0, Drawn: 0, Lost: 0, GoalsFor: 0, GoalsAgainst: 0, GoalDiff: 0, Elo: 1200, Tilt: 1.0},
}

func getTeams(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, teams)
}

func winnerTeamAPI(c *gin.Context) {
	id := c.Param("id")

	for i, t := range teams {
		if t.ID == id {
			registerWin(&teams[i], t.GoalsFor, t.GoalsAgainst)
			c.JSON(http.StatusOK, teams[i])
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
}

func loserTeamAPI(c *gin.Context) {
	id := c.Param("id")

	for i, t := range teams {
		if t.ID == id {
			registerLoss(&teams[i], t.GoalsFor, t.GoalsAgainst)
			c.JSON(http.StatusOK, teams[i])
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
}

func drawTeamAPI(c *gin.Context) {
	id := c.Param("id")

	for i, t := range teams {
		if t.ID == id {
			registerDraw(&teams[i], t.GoalsFor, t.GoalsAgainst)
			c.JSON(http.StatusOK, teams[i])
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
}

const K = 20

func simulateMatch(a *team, b *team) {
	rand.Seed(time.Now().UnixNano())

	EloA := 1.0 / (math.Pow(10, -(math.Abs(a.Elo-b.Elo))/400) + 1.0)
	EloB := 1.0 - EloA

	// 2.7 is the global average in football changes with teams styles
	expectedGoals := 2.7 * (a.Tilt + b.Tilt) / 2.0

	goalsA := int(math.Round(rand.Float64() * expectedGoals * a.Tilt / (a.Tilt + b.Tilt)))
	goalsB := int(math.Round(rand.Float64() * expectedGoals * b.Tilt / (a.Tilt + b.Tilt)))

	margin := math.Abs(float64(goalsA - goalsB))

	if goalsB < goalsA {
		baseChangeA := K * (1.0 - EloA) // A wins
		baseChangeB := K * (0.0 - EloB) // B loses
		if margin > 0 {
			scale := math.Sqrt(margin)
			baseChangeA *= scale
			baseChangeB *= scale
		}
		a.Elo += baseChangeA
		b.Elo += baseChangeB
		registerWin(a, goalsA, goalsB)
		registerLoss(b, goalsB, goalsA)
	} else if goalsA == goalsB {
		//draw
		a.Elo += K * (0.5 - EloA)
		b.Elo += K * (0.5 - EloB)
		registerDraw(a, goalsA, goalsB)
		registerDraw(b, goalsB, goalsA)
	} else {
		baseChangeA := K * (0.0 - EloA) // A loses
		baseChangeB := K * (1.0 - EloB) // B wins
		if margin > 0 {
			scale := math.Sqrt(margin)
			baseChangeA *= scale
			baseChangeB *= scale
		}
		a.Elo += baseChangeA
		b.Elo += baseChangeB
		registerWin(b, goalsB, goalsA)
		registerLoss(a, goalsA, goalsB)
	}

	// updating the play-styles of the teams based on goals scored
	a.Tilt = 0.98*a.Tilt + 0.02*(float64(goalsA+goalsB)/b.Tilt/expectedGoals)
	b.Tilt = 0.98*b.Tilt + 0.02*(float64(goalsA+goalsB)/a.Tilt/expectedGoals)
}

func registerWin(t *team, gf int, ga int) {
	t.Played++
	t.Win++
	t.Points += 3
	t.GoalsFor += gf
	t.GoalsAgainst += ga
	t.GoalDiff = t.GoalsFor - t.GoalsAgainst
}

func registerDraw(t *team, gf int, ga int) {
	t.Played++
	t.Drawn++
	t.Points++
	t.GoalsFor += gf
	t.GoalsAgainst += ga
	t.GoalDiff = t.GoalsFor - t.GoalsAgainst
}

func registerLoss(t *team, gf int, ga int) {
	t.Played++
	t.Lost++
	t.GoalsFor += gf
	t.GoalsAgainst += ga
	t.GoalDiff = t.GoalsFor - t.GoalsAgainst
}
