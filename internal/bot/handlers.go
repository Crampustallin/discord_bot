package bot

import (
	"fmt"
	"time"

	"github.com/Crampustallin/discord_bot/internal/bot/tools"
	"github.com/bwmarrin/discordgo"
)

func (b *Bot) channelCreateHandler(s *discordgo.Session, cc *discordgo.ChannelCreate) {
	fmt.Println("=======here+++======")
	if b.connected {
		return
	}
	b.connected = true
	b.vc = b.join(cc.GuildID, cc.ID)
}

func (b *Bot) voiceStateUpdateHandler(s *discordgo.Session, c *discordgo.VoiceStateUpdate) {
	if b.connected {
		return
	}
	if c.VoiceState.ChannelID == "" {
		fmt.Println("The user: " + c.UserID + " left the channel.")
		return
	}
	b.connected = true
	b.vc = b.join(c.GuildID, c.VoiceState.ChannelID)
}

func (b *Bot) join(guildId, channelId string) *discordgo.VoiceConnection {
	go func() {
		time.Sleep(10 * time.Second)
		b.disconnect()
	}()

	v, err := b.session.ChannelVoiceJoin(guildId, channelId, true, false)
	if err != nil {
		return nil
	}

	fmt.Println("Joined the channel: " + channelId)
	go func() {
		fileName, err := tools.HandleConversation(channelId, v.OpusRecv)
		if err != nil {
			fmt.Println(err.Error())
		}
		b.FileNameSend <- fileName
	}()
	return v
}

func (b *Bot) disconnect() {
	fmt.Println("Disconnecting from " + b.vc.ChannelID)
	close(b.vc.OpusRecv)
	b.vc.Disconnect()
	b.connected = false
}
