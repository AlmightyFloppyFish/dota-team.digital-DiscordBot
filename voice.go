package main

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"github.com/jeffreymkabot/ytdl"
	"github.com/jonas747/dca"
	"io"
	"strconv"
	"time"
)

func music(s *discordgo.Session, m *discordgo.MessageCreate, commandArray []string, command string) {

	// Get amount of args in command
	comLength := len(commandArray)

	// Selects if it should start music in "Multiple times" or "Single time" mode
	if comLength == 1 {
		PlayMusic(s, m, commandArray, command)
	} else if comLength == 2 {

		// Get int value from discord amount input
		amount, err := strconv.Atoi(commandArray[1])
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Could not understand how many times i was supposed to play that, You said to play it: "+commandArray[1]+"times")
			return
		}

		for count := 1; count <= amount; count++ {
			PlayMusic(s, m, commandArray, command)
		}
	}
}

// The actual command
func PlayMusic(s *discordgo.Session, m *discordgo.MessageCreate, commandArray []string, command string) {

	if command == "play" {

		callinguser := m.Author

		// Youtube stream options
		options := dca.StdEncodeOptions
		options.RawOutput = true
		options.Bitrate = 96
		options.Application = "lowdelay"
		options.Volume = 15

		videoInfo, err := ytdl.GetVideoInfo(commandArray[0])
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Could not get required youtube video info, Link is probably broken. Maybe the video is blocked in Sweden?")
			return
		}

		format := videoInfo.Formats.Extremes(ytdl.FormatAudioBitrateKey, true)[0]
		downloadURL, err := videoInfo.GetDownloadURL(format)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Could not parse video URL: parse.go | GetDownloadURL()")
			return
		}

		encodingSession, err := dca.EncodeFile(downloadURL.String(), options)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Could not create audio stream")
		}
		defer encodingSession.Cleanup()

		listeningchannel, err := joinUserVoiceChannel(s, callinguser.ID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Could not join your voice channel, Maybe i don't have permission to? Are you in a voice channel so i can join you?")
			return
		}

		done := make(chan error)
		dca.NewStream(encodingSession, listeningchannel, done)

		for {
			select {
			case <-stopMusicChannel:
				s.ChannelMessageSend(m.ChannelID, "Stopping song...")
				encodingSession.Stop()
				return
			case err2 := <-done:
				if err2 != nil && err2 != io.EOF {
					return
				}
				return
			default:
				// Prevent consumption of CPU
				time.Sleep(1 * time.Second)
			}
		}
	}
	return
}

func findUserVoiceState(session *discordgo.Session, userid string) (*discordgo.VoiceState, error) {
	for _, guild := range session.State.Guilds {
		for _, vs := range guild.VoiceStates {
			if vs.UserID == userid {
				return vs, nil

			}

		}

	}
	return nil, errors.New("Could not find user's voice state")

}
func joinUserVoiceChannel(session *discordgo.Session, userID string) (*discordgo.VoiceConnection, error) {
	// Find a user's current voice channel
	vs, err := findUserVoiceState(session, userID)
	if err != nil {
		return nil, err

	}

	// Join the user's channel and start unmuted and deafened.
	return session.ChannelVoiceJoin(vs.GuildID, vs.ChannelID, false, true)
}
