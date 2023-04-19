package main

import (
	"CCP_backend/backend"
)

func main() {
	server := backend.NewGin()
	server.Init()
	server.Start(":9090")
}
