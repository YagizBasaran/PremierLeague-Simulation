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

	//API Endpoints
	router.GET("/teams", getTeams)
	router.GET("/results", getMatchResults)
	router.POST("/weekly-schedule", weeklyScheduleHandler)
	router.POST("/reset", reset)
	router.POST("/all-league-schedule", playWholeLeague)

	router.POST("/teams/:id/win", winnerTeamAPI)
	router.POST("/teams/:id/draw", drawTeamAPI)
	router.POST("/teams/:id/loss", loserTeamAPI)

	// simulateMatch(&teams[0], &teams[1])
	// simulateMatch(&teams[0], &teams[2])
	// simulateMatch(&teams[0], &teams[3])
	// simulateMatch(&teams[1], &teams[0])
	// simulateMatch(&teams[2], &teams[0])
	// simulateMatch(&teams[3], &teams[0])
	// simulateMatch(&teams[1], &teams[2])
	// simulateMatch(&teams[1], &teams[3])
	// simulateMatch(&teams[2], &teams[1])
	// simulateMatch(&teams[3], &teams[1])
	// simulateMatch(&teams[2], &teams[3])
	// simulateMatch(&teams[3], &teams[2])

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

type match struct {
	Week     int    `json:"week"`
	HomeTeam string `json:"home"`
	AwayTeam string `json:"away"`
	Score    string `json:"score"`
}

var teams = []team{
	{ID: "1", Name: "Chelsea", Points: 0, Played: 0, Win: 0, Drawn: 0, Lost: 0, GoalsFor: 0, GoalsAgainst: 0, GoalDiff: 0, Elo: 1200, Tilt: 1.0},
	{ID: "2", Name: "Arsenal", Points: 0, Played: 0, Win: 0, Drawn: 0, Lost: 0, GoalsFor: 0, GoalsAgainst: 0, GoalDiff: 0, Elo: 1200, Tilt: 1.0},
	{ID: "3", Name: "Man City", Points: 0, Played: 0, Win: 0, Drawn: 0, Lost: 0, GoalsFor: 0, GoalsAgainst: 0, GoalDiff: 0, Elo: 1200, Tilt: 1.0},
	{ID: "4", Name: "Liverpool", Points: 0, Played: 0, Win: 0, Drawn: 0, Lost: 0, GoalsFor: 0, GoalsAgainst: 0, GoalDiff: 0, Elo: 1200, Tilt: 1.0},
}

// randomly assigned league matches
var schedule = [][][2]int{
	{{0, 1}, {2, 3}},
	{{0, 2}, {1, 3}},
	{{0, 3}, {1, 2}},
	{{1, 0}, {3, 2}},
	{{2, 0}, {3, 1}},
	{{3, 0}, {2, 1}},
}

var currentWeek = 0
var results []match

//API Functions

func getTeams(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, teams)
}

func getMatchResults(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, results)
}

func weeklyScheduleHandler(c *gin.Context) {

	if currentWeek == len(schedule) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Season completed"})
		return
	}
	simulateCurrentWeek()
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Week %d simulated", currentWeek)})
}

func reset(c *gin.Context) {
	currentWeek = 0
	teams = []team{
		{ID: "1", Name: "Chelsea", Points: 0, Played: 0, Win: 0, Drawn: 0, Lost: 0, GoalsFor: 0, GoalsAgainst: 0, GoalDiff: 0, Elo: 1200, Tilt: 1.0},
		{ID: "2", Name: "Arsenal", Points: 0, Played: 0, Win: 0, Drawn: 0, Lost: 0, GoalsFor: 0, GoalsAgainst: 0, GoalDiff: 0, Elo: 1200, Tilt: 1.0},
		{ID: "3", Name: "Man City", Points: 0, Played: 0, Win: 0, Drawn: 0, Lost: 0, GoalsFor: 0, GoalsAgainst: 0, GoalDiff: 0, Elo: 1200, Tilt: 1.0},
		{ID: "4", Name: "Liverpool", Points: 0, Played: 0, Win: 0, Drawn: 0, Lost: 0, GoalsFor: 0, GoalsAgainst: 0, GoalDiff: 0, Elo: 1200, Tilt: 1.0},
	}
	results = nil
}

func playWholeLeague(c *gin.Context) {

	for i := currentWeek; i < len(schedule); i++ {
		simulateCurrentWeek()
	}
	c.JSON(http.StatusOK, gin.H{"message": "Whole season simulated"})
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

//Functions

func simulateCurrentWeek() {
	if currentWeek == len(schedule) {
		return
	}

	thisWeeksMatchs := schedule[currentWeek]

	for _, game := range thisWeeksMatchs {
		homeTeam := &teams[game[0]]
		awayTeam := &teams[game[1]]

		goalsHome, goalsAway := simulateMatch(homeTeam, awayTeam)
		score := fmt.Sprintf("%d - %d", goalsHome, goalsAway)

		match := match{Week: currentWeek + 1, HomeTeam: homeTeam.Name, AwayTeam: awayTeam.Name, Score: score}
		results = append(results, match)
	}
	currentWeek++
}

const K = 20

func simulateMatch(a *team, b *team) (int, int) {
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

	return goalsA, goalsB
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
