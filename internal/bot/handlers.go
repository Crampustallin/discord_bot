package bot

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/pion/rtp"
	"github.com/pion/webrtc/v3/pkg/media/oggwriter"
)

func createPionRtpPacket(p *discordgo.Packet) *rtp.Packet {
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

func handleConversation(chanId int, c chan *discordgo.Packet) (string, error) {
	file, err := oggwriter.New(fmt.Sprintf("%d.ogg", chanId), 48000, 2)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Failed to create file %d.ogg, giving up on recording: %v\n", chanId, err.Error()))
	}

	defer file.Close()

	for p := range c {
		rPacket := createPionRtpPacket(p)
		if err := file.WriteRTP(rPacket); err != nil {
			return "", errors.New(fmt.Sprintf("Failed to create file %d.ogg, giving up on recording: %v\n", chanId, err.Error()))
		}
	}
	return strconv.Itoa(chanId), nil
}

func VoiceStateUpdateHandler(s *discordgo.Session, vsu *discordgo.VoiceStateUpdate) {
	// TODO: make user joined or created vc event handler
}
