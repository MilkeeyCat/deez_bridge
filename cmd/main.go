package main

import (
	"os"
	"os/signal"

	"github.com/MilkeeyCat/deez_bridge/internal/bridge"
	"github.com/MilkeeyCat/deez_bridge/internal/logger"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".dev.env", ".env")
	logger.SetupLogger("./mumbo_jumbo.logs")

	bridge := bridge.NewBridge()
	bridge.Open()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	logger.Logger.Info("exiting...")
	bridge.Close()
}
