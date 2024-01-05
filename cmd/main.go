package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// Player represents the player's data structure
type Player struct {
	PlayerName string  `json:"playerName"`
	PositionX  float32 `json:"positionX"`
	PositionY  float32 `json:"positionY"`
}

func main() {
	http.HandleFunc("/", handleRequest)
	http.ListenAndServe(":8080", nil)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {

	// Set environment variables
	os.Setenv("DB_USERNAME", "secret_username")
	os.Setenv("DB_PASSWORD", "secret_password")
	os.Setenv("DB_NAME", "secret_database_name")

	if r.Method == "POST" {
		playerName := r.FormValue("playerName")

		// Connect to the database (replace with your database credentials)
		db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s",
			os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME")))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		// Check if the player exists
		var player Player
		err = db.QueryRow("SELECT playerName, positionX, positionY FROM players WHERE playerName=?", playerName).
			Scan(&player.PlayerName, &player.PositionX, &player.PositionY)

		if err == sql.ErrNoRows {
			// Player doesn't exist, create a new player
			_, err := db.Exec("INSERT INTO players (playerName, positionX, positionY) VALUES (?, 0, 0)", playerName)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			player = Player{PlayerName: playerName, PositionX: 0, PositionY: 0}
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Convert player data to JSON
		playerJSON, err := json.Marshal(player)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Send JSON response
		w.Header().Set("Content-Type", "application/json")
		w.Write(playerJSON)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
