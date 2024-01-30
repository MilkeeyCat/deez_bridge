package main

import (
	"sync"

	"github.com/MilkeeyCat/deez_bridge/internal/bridge"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".dev.env", ".env")
	var wg sync.WaitGroup

	bridge := bridge.NewBridge()
	bridge.Run(&wg)

	wg.Wait()
}
