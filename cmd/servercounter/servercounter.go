package main

import (
	"github.com/TicketsBot/GoServerCounter/config"
	"github.com/TicketsBot/GoServerCounter/database"
	"github.com/TicketsBot/GoServerCounter/http"
	"time"
)

func main() {
	config.LoadConfig()

	db := database.NewDatabase()
	go pollDatabase(db)

	http.StartServer()
}

func pollDatabase(db *database.Database) {
	for {
		http.Count = db.GetServerCount()
		time.Sleep(time.Second * 5)
	}
}
