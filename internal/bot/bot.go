package bot

import (
	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	session *discordgo.Session
}

func NewBot(s *discordgo.Session) *Bot {
	return &Bot{session: s}
}

func (b *Bot) Start() {
	b.session.AddHandler(VoiceStateUpdateHandler)

	err := b.session.Open()
	if err != nil {
		panic(err)
	}
}

func (b *Bot) Close() {
	b.session.Close()
}
