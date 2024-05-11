package bot

import (
	"fmt"
	"strings"
	"time"

	"github.com/Crampustallin/discord_bot/internal/bot/tools"
	"github.com/bwmarrin/discordgo"
)

func (b *Bot) commandsHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot || m.Author.ID == s.State.User.ID {
		return
	}
	content := m.Content
	if !strings.HasPrefix(content, "!") {
		return
	}

	splited := strings.Split(content, " ")
	command := strings.Replace(splited[0], "!", "", 1)
	if c, ok := b.commands[command]; ok {
		m.Content = strings.Join(splited[1:], " ")
		c(s, m)
	} else {
		_, err := s.ChannelMessageSend(m.ChannelID, "No command found "+command)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (b *Bot) channelCreateHandler(s *discordgo.Session, cc *discordgo.ChannelCreate) {
	if b.connected {
		return
	}
	b.connected = true
	b.vc = b.join(cc.GuildID, cc.ID)
}

func (b *Bot) voiceStateUpdateHandler(s *discordgo.Session, c *discordgo.VoiceStateUpdate) {
	if c.VoiceState.ChannelID == "" {
		if b.connected {
			b.userLeft(c.UserID)
		}
		return
	}
	if b.connected {
		if c.UserID != s.State.User.ID && c.VoiceState.ChannelID == b.vc.ChannelID {
			fmt.Println(b.users)
			b.users[c.UserID] = true
			fmt.Println(b.users)
		}
		if c.UserID != s.State.User.ID && c.VoiceState.ChannelID != b.vc.ChannelID {
			b.userLeft(c.UserID)
		}
		return
	}
	b.connected = true
	b.vc = b.join(c.GuildID, c.VoiceState.ChannelID)
	b.users[c.UserID] = true
}

func (b *Bot) join(guildId, channelId string) *discordgo.VoiceConnection {
	v, err := b.session.ChannelVoiceJoin(guildId, channelId, true, false)
	if err != nil {
		b.connected = false
		fmt.Println("Failed to join the voice chat")
		return nil
	}

	go func() {
		time.Sleep(time.Duration(b.Timer) * time.Second)
		if b.connected {
			b.disconnect()
		}
	}()

	fmt.Println("Joined the channel: " + channelId)
	go func() {
		fileName, err := tools.HandleConversation(channelId, v.OpusRecv)
		if err != nil {
			fmt.Println(err.Error())
		}
		if b.storage != nil {
			if err = b.storage.Upload(fileName); err != nil {
				b.session.ChannelMessageSend(channelId, "Failed to upload file")
			} else {
				b.session.ChannelMessageSend(channelId, fileName)
			}
		}
	}()
	return v
}

func (b *Bot) disconnect() {
	fmt.Println("Disconnecting from " + b.vc.ChannelID)
	close(b.vc.OpusRecv)
	b.vc.Disconnect()
	b.connected = false
}

func (b *Bot) userLeft(userId string) {
	fmt.Println("The user: " + userId + " left the channel.")
	delete(b.users, userId)
	if len(b.users) <= 0 && b.connected {
		b.disconnect()
	}
}
