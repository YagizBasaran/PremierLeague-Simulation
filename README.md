# 🏆 Premier League Simulator

A web-based football league simulation app where teams play weekly fixtures, scores and stats are updated dynamically, and real-time predictions for the championship are calculated.
Built using **Go (Gin)** for the backend and **HTML/CSS/JavaScript** for the frontend.

## 1 Live Demo

Deployed via Render: [Project Link](https://premierleague-simulation.onrender.com)
<br><br>
https://premierleague-simulation.onrender.com/web/

## 2 Setup & Installation

### a. Clone the repo

```bash
git clone https://github.com/YagizBasaran/PremierLeague-Simulation.git
cd PremierLeague-Simulation
```

### b. Install dependencies

```bash
go mod tidy
```

### c. Run the app

```bash
go run footballsim.go
```

By default, the app runs at:  
📍 `http://localhost:8079`  
Frontend is served from:  
📍 `http://localhost:8079/web/`

> To deploy, make sure your environment sets the correct `PORT` variable (e.g. Render).

---

## 🌐 API Endpoints

| Method | Endpoint             | Description                          |
|--------|----------------------|--------------------------------------|
| GET    | `/teams`             | Returns current teams ALL data       |
| GET    | `/results`           | Returns all match results so far     |
| GET    | `/championship-ratio`| Returns title prediction percentages |
||||
| POST   | `/finish-week`       | Simulates the next week              |
| POST   | `/finish-season`     | Simulates all remaining weeks        |
| POST   | `/reset`             | Resets the league                    |
| POST   | `/teams/:id/win`     | Manually register a win              |
| POST   | `/teams/:id/draw`    | Manually register a draw             |
| POST   | `/teams/:id/loss`    | Manually register a loss             |

---

## Some Concepts I want to highlight

### 1) ELO & Tilt System
Algorithms obtained from http://clubelo.com/System
- **ELO** rating reflects a team's strength based on match outcomes and goal margins. Just like chess.
  - `E = 1 / (10(-dr/400) + 1)` where "dr" is the Elo point difference of the 2 clubs.

- **ELO Exchange**: When clubs play each other and win or lose, they exchange points. The number of points exchanged must be determined so that a certain win rate between two clubs makes the Elo difference between both clubs converge towards the Elo difference that corresponds to this win rate.<br><br>The following equation satisfies this constraint: `ΔElo_1X2 = (R - E) * k` <br>where R is the result (1 for win, 0.5 for draw, 0 for loss).
There is one degree of freedom in this equation which in the weight index k that has to be chosen. A higher k will have the ratings converge quicker to their true values but will suffer from more variation. A smaller k provides more stable values that take longer to converge.
<br><br>ClubElo uses a weight index of k = 20.

- **Tilt** reflects offensive style (more tilt = more goals, less tilt = defensive team).
  - `New_tilt = 0.98 * Old_tilt + 0.02 *Game_total_goals/Opposition_tilt/Expected_Game_total_goals`

 
- Expected goals are calculated using both ELO and Tilt.
- Match outcomes adjust both ELO and Tilt dynamically.

### 2) Championship Prediction
Could've used Monte Carlo simulation to run 1000 matches every week and calculate the champion https://algolritmo.com/index.php/2017/02/14/predicting-the-champions-league/. However, using a basic points-based model worked just as fine and take less time to implement:

- Each team's potential score = `current points + max possible points / 2`
- Probability is normalized over all "still-alive" teams.
---

### References
- https://go.dev/doc/tutorial/web-service-gin
- https://www.youtube.com/watch?v=9bRMLKBbFMQ
- http://clubelo.com/
