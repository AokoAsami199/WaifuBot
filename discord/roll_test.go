package discord

import (
	"context"
	"os"
	"strconv"
	"testing"

	"github.com/Karitham/WaifuBot/service/anilist"
	"github.com/Karitham/WaifuBot/service/store"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/session"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Test_bot_Roller(t *testing.T) {
	d, _ := t.Deadline()
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()

	_ = godotenv.Load("../.env")
	log.Level(zerolog.TraceLevel)

	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		t.Skip("BOT_TOKEN not set")
	}
	guildID := os.Getenv("GUILD_ID")
	if guildID == "" {
		t.Skip("GUILD_ID not set")
	}
	appID := os.Getenv("APP_ID")
	if appID == "" {
		t.Skip("APP_ID not set")
	}

	randomer := NewMockRandomer()
	randomer.RandomFunc.SetDefaultReturn(anilist.Character{
		Name: anilist.Name{
			Full: "Midori Asakusa",
		},
		Description: `Asakusa is a freshman at Shibahama Highschool who likes anime.
		She is full of curiosity and imagination, yet not that sociable.
		Her policy is “Anime is all about its setting” and she has kept various settings inspired by her daily life in her sketchbook.
		\n\nSayaka Kanamori, Asakusa’s classmate, used to be the only person who knows about Asakusa’s passion,
		but when they visited Anime Club together they got along with Tsubame Mizusaki,
		which led them to start up Eizouken. Her goal is to explore her own “Incredible World”.`,
		Siteurl: "https://anilist.co/character/149676",
		Image:   anilist.Image{Large: "https://s4.anilist.co/file/anilistcdn/character/large/b149676-TxktIkOe5Xl1.jpg"},
		ID:      149676,
	}, nil)
	storer := NewMockStorer()
	storer.GetFunc.SetDefaultReturn([]store.Character{}, nil)
	storer.PutFunc.SetDefaultHook(func(i int, c store.Character) error {
		t.Log(c)
		return nil
	})
	storer.PutFunc.SetDefaultReturn(nil)

	b := Setup(t, appID, token)
	if b == nil {
		t.Fatal("Setup failed")
		return
	}

	b.RegisterGuild(discord.GuildID(MustAtoi(guildID)), map[string]Commander{
		"roll": Command{fn: b.Roller(randomer, storer), cmd: api.CreateCommandData{
			Name:        "roll",
			Description: "test a roll",
		}},
	})

	if err := b.Open(ctx); err != nil {
		t.Fatal(err)
		return
	}
	defer b.s.Close()

	<-ctx.Done()
}

func Setup(t *testing.T, appID, token string, commands ...Command) *bot {
	s, err := session.New("Bot " + token)
	if err != nil {
		t.Fatal(err)
		return nil
	}

	b := New(s, appIDFromString(appID))
	if err != nil {
		t.Fatal(err)
		return nil
	}

	return b
}

func MustAtoi(s string) uint64 {
	i, _ := strconv.ParseUint(s, 10, 64)
	return i
}
