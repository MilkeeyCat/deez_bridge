package main

import (
	"fmt"

	bridge "github.com/MilkeeyCat/deez_bridge/internal"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".dev.env", ".env")

	bridge := bridge.NewBridge()
	bridge.Run()

    fmt.Println("leaving...")
}
