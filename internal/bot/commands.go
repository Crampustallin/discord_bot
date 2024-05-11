package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) UrlCommandFun(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Content == "" {
		s.ChannelMessageSend(m.ChannelID, "The url command usage: !url [key]")
		return
	}
	url, err := b.storage.GetObjUrl(m.Content)
	if err != nil {
		panic(err)
	}
	message := fmt.Sprintf("%s\nExpires in %s", url.Link, url.Expires)
	_, err = s.ChannelMessageSend(m.ChannelID, message)
	if err != nil {
		fmt.Println(err)
	}
}

func (b *Bot) ListCommandFun(s *discordgo.Session, m *discordgo.MessageCreate) {
	list, err := b.storage.GetObjList()
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "The list couldn't be obtain")
		return
	}
	_, err = s.ChannelMessageSend(m.ChannelID, strings.Join(list, "\n"))
	if err != nil {
		fmt.Println(err)
	}
}
