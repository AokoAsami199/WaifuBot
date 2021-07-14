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

	//nolint:errcheck
	godotenv.Load()
	appID := os.Getenv("APP_ID")
	token := os.Getenv("BOT_TOKEN")

	rm := discord.LS(appID, token)
	defer rm()

	select {}
}
