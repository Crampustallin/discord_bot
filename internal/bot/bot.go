package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	session      *discordgo.Session
	vc           *discordgo.VoiceConnection
	userCount    int
	connected    bool
	FileNameSend chan string
}

func NewBot(token string) *Bot {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}
	fileNameSend := make(chan string)
	return &Bot{session: s, FileNameSend: fileNameSend}
}

func (b *Bot) Start() {
	b.session.AddHandler(b.channelCreateHandler)
	b.session.AddHandler(b.voiceStateUpdateHandler)

	fmt.Println("Bot started...")

	err := b.session.Open()
	if err != nil {
		panic(err)
	}
}

func (b *Bot) Close() {
	b.session.Close()
	close(b.FileNameSend)
}
