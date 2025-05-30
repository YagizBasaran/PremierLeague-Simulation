package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"math"
	"math/rand"
	"time"

	"os"
)

func main() {
	router := gin.Default()

	router.Static("/web", "./public")

	//API Endpoints
	router.GET("/teams", getTeams)
	router.GET("/results", getMatchResults)
	router.GET("/team", getTeam)
	router.POST("/finish-week", weeklyScheduleHandler)
	router.POST("/reset", reset)
	router.POST("/finish-season", playWholeLeague)

	router.GET("/championship-ratio", func(c *gin.Context) {
		ratios := calculateChampionshipRatio()
		c.JSON(http.StatusOK, ratios)
	})

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

	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/web/")
	})

	// fmt.Println("Starting server at http://localhost:8079")
	// router.Run("localhost:8079")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8079"
	}
	fmt.Println("Starting server on port", port)
	router.Run(":" + port)

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

func getTeam(c *gin.Context) {
	name := c.Query("name")
	id := c.Query("id")

	for _, t := range teams {
		if t.Name == name || t.ID == id {
			c.JSON(http.StatusOK, t)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
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

// Functions
// Points based calculation (can be elo and 1000 games simulation based too)
func calculateChampionshipRatio() map[string]float64 {
	gamesRemaining := 6 - currentWeek
	maxPointsPerTeam := gamesRemaining * 3

	ratios := make(map[string]float64)
	totalProbability := 0.0

	for i, team := range teams {
		maxPoints := team.Points + maxPointsPerTeam

		canWin := true
		for j, otherTeam := range teams {
			if i != j {
				if otherTeam.Points > maxPoints {
					canWin = false
					break
				}
			}
		}

		if !canWin {
			ratios[team.Name] = 0.0
		} else {
			score := float64(team.Points) + float64(gamesRemaining)*0.5
			ratios[team.Name] = score
			totalProbability += score
		}
	}

	for name, score := range ratios {
		if totalProbability > 0 {
			ratios[name] = (score / totalProbability) * 100
		}
	}

	return ratios
}

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
	//increased to 5 to see more goals
	expectedGoals := 5.0 * (a.Tilt + b.Tilt) / 2.0

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
