package main

import (
	"github.com/TicketsBot/GoServerCounter/database"
	"github.com/TicketsBot/GoServerCounter/http"
	"time"
)

func main() {
	db := database.NewDatabase()
	go pollDatabase(db)

	http.StartServer()
}

func pollDatabase(db *database.Database) {
	for {
		http.Lock.Lock()
		http.Count = db.GetServerCount()
		http.Lock.Unlock()

		time.Sleep(time.Second * 5)
	}
}
