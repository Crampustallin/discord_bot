package tools

import (
	"errors"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/pion/rtp"
	"github.com/pion/webrtc/v3/pkg/media/oggwriter"
)

func CreatePionRtpPacket(p *discordgo.Packet) *rtp.Packet {
	return &rtp.Packet{
		Header: rtp.Header{
			Version:        2,
			PayloadType:    0x78,
			SequenceNumber: p.Sequence,
			Timestamp:      p.Timestamp,
			SSRC:           p.SSRC,
		},
		Payload: p.Opus,
	}
}

func HandleConversation(chanId string, c chan *discordgo.Packet) (string, error) {
	fileName := chanId + time.Now().Format("20060102150405") + ".ogg"
	file, err := oggwriter.New(fmt.Sprintf(fileName), 48000, 2)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Failed to create file %s.ogg, giving up on recording: %v\n", fileName, err.Error()))
	}

	defer file.Close()

	for p := range c {
		rPacket := CreatePionRtpPacket(p)
		if err := file.WriteRTP(rPacket); err != nil {
			return "", errors.New(fmt.Sprintf("Failed to create file %s.ogg, giving up on recording: %v\n", fileName, err.Error()))
		}
	}
	return fileName, nil
}
