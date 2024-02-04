package main

import (
	"fmt"
	"os"
	"os/signal"

	bridge "github.com/MilkeeyCat/deez_bridge/internal"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".dev.env", ".env")

	bridge := bridge.NewBridge()
	bridge.Open()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	fmt.Println("leaving ...")
	bridge.Close()
}
