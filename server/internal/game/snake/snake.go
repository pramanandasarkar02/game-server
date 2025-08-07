package snake

import (
	// "crypto/rand"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"math/rand"

	"github.com/gorilla/websocket"
)


type PlayerSnakeState struct {
    ID string `json:"id"`
    X int   `json:"x"`
    Y int   `json:"y"`
    Direction Point `json:"direction"`
    Tail []Point    `json:"tail"`
    Length int     `json:"length"`
    Score int      `json:"score"`
    Alive bool     `json:"alive"`
}

type Point struct {
    X int
    Y int
}

type GameState struct {
    PlayerSnakerStates map[string]*PlayerSnakeState `json:"player_snake_states"`
    Food []Point    `json:"food"`
    Poison []Point   `json:"poison"`
    Width int   `json:"width"`
    Height int  `json:"height"`
    FoodValue int `json:"food_value"`
    mu sync.Mutex
}

type Queue struct {
    PlayerSnakerStates []PlayerSnakeState
    mu sync.Mutex
}


var (
    upgrader = websocket.Upgrader{
        ReadBufferSize:  1024,
        WriteBufferSize: 1024,
        CheckOrigin: func(r *http.Request) bool {
            return true
        },
    }
    
    gameState = &GameState{
        PlayerSnakerStates: map[string]*PlayerSnakeState{},
        Food: []Point{},
        Poison: []Point{},
        Width: 80,
        Height: 60,
        FoodValue: 10,
    }

    conns     = make(map[string]*websocket.Conn)
	connsMu   sync.Mutex
	ticker    = time.NewTicker(200 * time.Millisecond) // Game updates every 200ms
	// rand.Seed(time.Now().UnixNano())
)

func broadcastGameState() {
    gameState.mu.Lock()
    data, err := json.Marshal(gameState)
    gameState.mu.Unlock()
    if err != nil {
        log.Println("Failed to marshal game state:", err)
        return
    }
    
    connsMu.Lock()
    defer connsMu.Unlock()
    for id, conn := range conns {
        if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
            log.Println("Failed to broadcast game state:", err)

            conn.Close()
            delete(conns, id)
            gameState.mu.Lock()
            delete(gameState.PlayerSnakerStates, id)
            gameState.mu.Unlock()
        }
    }
    
}

func spawnFoodOrPoison(isFood bool) {
    gameState.mu.Lock()
    defer gameState.mu.Lock()
    var pos Point 
    for {
		pos = Point{X: rand.Intn(gameState.Width), Y: rand.Intn(gameState.Height)}
		occupied := false
		for _, p := range gameState.Food {
			if p == pos {
				occupied = true
				break
			}
		}
		for _, p := range gameState.Poison {
			if p == pos {
				occupied = true
				break
			}
		}
		for _, snake := range gameState.PlayerSnakerStates {
			if (snake.X == pos.X && snake.Y == pos.Y) || containsPoint(snake.Tail, pos) {
				occupied = true
				break
			}
		}
		if !occupied {
			break
		}
	}
	if isFood {
		gameState.Food = append(gameState.Food, pos)
	} else {
		gameState.Poison = append(gameState.Poison, pos)
	}
}

func containsPoint(points []Point, p Point) bool {
    for _, point := range points {
        if point.X == p.X && point.Y == p.Y {
            return true
        }
    }
    return false
}


func gameLoop() {
    for range ticker.C {
        gameState.mu.Lock()

        for _, snake := range gameState.PlayerSnakerStates {
            if !snake.Alive {
                continue
            }
            newX, newY := snake.X + snake.Direction.X, snake.Y + snake.Direction.Y
            // out ranged
            if newX < 0 || newX >= gameState.Width || newY < 0 || newY >= gameState.Height {
                snake.Alive = false
                continue
            }
            // self collision
            if containsPoint(snake.Tail, Point{newX, newY}) {
                snake.Alive = false
                continue
            }

            // update position
            snake.Tail = append([]Point{{
                X: snake.X,
                Y: snake.Y,
            }}, snake.Tail[:snake.Length]...)
            snake.X, snake.Y = newX, newY

            // check food
            for i, food := range gameState.Food {
                if snake.X == food.X && snake.Y == food.Y {
                    snake.Length++
                    snake.Score += gameState.FoodValue
                    gameState.Food = append(gameState.Food[:i], gameState.Food[i+1:]...)
                    spawnFoodOrPoison(true)
                    break
                }
            }

            // check poison position
            for i, poison := range gameState.Poison {
                if snake.X == poison.X && snake.Y == poison.Y {
                    snake.Alive = false
                    gameState.Poison = append(gameState.Poison[:i], gameState.Poison[i+1:]...)
                    spawnFoodOrPoison(false)
                    break
                }
            }

        }
        if len(gameState.Food) == 0 {
            spawnFoodOrPoison(true)
        }
        if len(gameState.Poison) == 0 {
            spawnFoodOrPoison(false)
        }
        gameState.mu.Unlock()
        broadcastGameState()
    }
}




func handleConnecttions(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        http.Error(w, "Something went wrong", http.StatusBadRequest)
        return 
    }
    id := r.URL.Query().Get("id")
    if id == "" {
        conn.Close()
        return 
    }
    gameState.mu.Lock()
    snake := &PlayerSnakeState{
        ID: id,
        X: rand.Intn(gameState.Width),
        Y: rand.Intn(gameState.Height),
        Direction: Point{X: 1, Y: 0},
        Tail: []Point{},
        Length: 1,
        Score: 0,
        Alive: true,
    }
    gameState.PlayerSnakerStates[id] = snake
    gameState.mu.Unlock()

    connsMu.Lock()
    conns[id] = conn
    connsMu.Unlock()

    for {
        _, msg, err := conn.ReadMessage()
        if err != nil {
            log.Printf("Error reading from %s: %v", id, err)
            conn.Close()
            connsMu.Lock()
            delete(conns, id)
            connsMu.Unlock()
            gameState.mu.Lock()
            delete(gameState.PlayerSnakerStates, id)
            gameState.mu.Unlock()
            broadcastGameState()
            return 
        }
        var input struct{
            Type string `json:"type"`
            Direction Point `json:"direction"`
        }
        if err:= json.Unmarshal(msg, &input); err == nil && input.Type == "setDirection" {
            gameState.mu.Lock()
            if snake, exists := gameState.PlayerSnakerStates[id]; exists && snake.Alive {
                if (snake.Direction.X != input.Direction.X || snake.Direction.Y != input.Direction.Y) && (input.Direction.X == 0 || input.Direction.Y == 0)  && (snake.Direction.X != 0 || snake.Direction.Y != 0) {
                    snake.Direction = input.Direction
                }
            }
            gameState.mu.Unlock()
        }
    }
}

func StartGame() {
    spawnFoodOrPoison(true)
    spawnFoodOrPoison(false)
    go gameLoop()
}

func SetupRoutes() {
    http.HandleFunc("/ws", handleConnecttions)
    go StartGame()
    log.Fatal(http.ListenAndServe(":8080", nil))
}