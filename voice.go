package main

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	"github.com/rylio/ytdl"
	"io"
	"strconv"
	"time"
)

func music(s *discordgo.Session, m *discordgo.MessageCreate, commandArray []string) {

	defer recover()

	// Get amount of args in command
	comLength := len(commandArray)

	// Selects if it should start music in "Multiple times" or "Single time" mode
	if comLength == 2 {
		PlayMusic(s, m, commandArray)
	} else if comLength == 3 {

		// Get int value from discord amount input
		amount, err := strconv.Atoi(commandArray[2])
		if err != nil {
			panic(err)
		}

		for count := 1; count <= amount; count++ {
			PlayMusic(s, m, commandArray)
		}
	}

}

// The actual command
func PlayMusic(s *discordgo.Session, m *discordgo.MessageCreate, commandArray []string) {

	if commandArray[0] == "play" {

		defer recover()

		callinguser := m.Author

		// Youtube stream options
		options := dca.StdEncodeOptions
		options.RawOutput = true
		options.Bitrate = 96
		options.Application = "lowdelay"
		options.Volume = 15

		videoInfo, err := ytdl.GetVideoInfo(commandArray[1])
		if err != nil {
			panic(err)
			s.ChannelMessageSend(m.ChannelID, "Could not get required youtube video info")
		}

		format := videoInfo.Formats.Extremes(ytdl.FormatAudioBitrateKey, true)[0]
		downloadURL, err := videoInfo.GetDownloadURL(format)
		if err != nil {
			panic(err)
			s.ChannelMessageSend(m.ChannelID, "Could not parse video URL")
		}

		encodingSession, err := dca.EncodeFile(downloadURL.String(), options)
		if err != nil {
			panic(err)
		}
		defer encodingSession.Cleanup()

		listeningchannel, err := joinUserVoiceChannel(s, callinguser.ID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Could not join your voice channel, Maybe i don't have permission to? Are you in a voice channel so i can join you?")
			panic(err)
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
					panic(err)
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
