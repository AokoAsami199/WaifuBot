package main

import (
	"os"

	"github.com/Karitham/WaifuBot/discord"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.NewConsoleWriter())
	log.Logger = log.Output(os.Stderr)
	log.Logger = log.Level(zerolog.TraceLevel)

	//nolint:errcheck
	godotenv.Load()
	appID := os.Getenv("APP_ID")
	token := os.Getenv("BOT_TOKEN")

	close := discord.LS(appID, token)
	defer close()

	select {}
}
