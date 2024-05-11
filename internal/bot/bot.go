package bot

import (
	"fmt"

	"github.com/Crampustallin/discord_bot/internal/models"
	"github.com/bwmarrin/discordgo"
)

type Storage interface {
	Upload(string) error
	GetObjUrl(string) (*models.Url, error)
	GetObjList() ([]string, error)
}

type commandFun func(s *discordgo.Session, i *discordgo.MessageCreate)

type Bot struct {
	session   *discordgo.Session
	vc        *discordgo.VoiceConnection
	userCount int
	connected bool
	storage   Storage
	commands  map[string]commandFun
	users     map[string]bool
	Timer     int
}

func NewBot(token string) *Bot {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}
	return &Bot{
		session:  s,
		commands: make(map[string]commandFun),
		users:    make(map[string]bool),
		Timer:    10,
	}
}

func NewBotWithStorage(token string, storage Storage) *Bot {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}
	return &Bot{
		session:  s,
		storage:  storage,
		commands: make(map[string]commandFun),
		users:    make(map[string]bool),
		Timer:    10,
	}
}

func (b *Bot) Start() {
	b.commands["url"] = b.UrlCommandFun
	b.commands["list"] = b.ListCommandFun

	b.session.AddHandler(b.commandsHandler)
	b.session.AddHandler(b.channelCreateHandler)
	b.session.AddHandler(b.voiceStateUpdateHandler)

	fmt.Println("Bot started...")

	if err := b.session.Open(); err != nil {
		panic(err)
	}
}

func (b *Bot) Close() {
	b.session.Close()
}
