package main

import (
	"github.com/HistoryLabs/events-api/server"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	server.Init()
}
