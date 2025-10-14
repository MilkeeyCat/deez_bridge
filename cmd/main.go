package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/MilkeeyCat/deez_bridge/internal/bridge"
	"github.com/MilkeeyCat/deez_bridge/internal/logger"
	"github.com/joho/godotenv"
)

func main() {
	logger.SetupLogger("./mumbo_jumbo.logs")
	if err := godotenv.Load(".dev.env", ".env"); err != nil {
		logger.Fatal(fmt.Sprintf("failed to load env vars: %v", err.Error()))
	}

	bridge := bridge.NewBridge()
	bridge.Open()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	logger.Logger.Info("exiting...")
	bridge.Close()
}
