package main

import (
	"github.com/TicketsBot/GoServerCounter/config"
	"github.com/TicketsBot/GoServerCounter/http"
)

func main() {
	config.LoadConfig()
	http.StartServer()
}
