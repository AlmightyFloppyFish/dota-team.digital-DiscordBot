package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/mattn/go-shellwords"
	"strings"
)

// Splits the words into string array on ' : '
func ColonSplit(s *discordgo.Session, m *discordgo.MessageCreate) []string {

	byteArr := []byte(m.Content)
	prefixByte := []byte(prefix)

	if byteArr[0] == prefixByte[0] {

		commandArray := strings.Split(m.Content, " : ")

		return commandArray
	}
	return []string{}
}

func easySplit(s *discordgo.Session, m *discordgo.MessageCreate) ([]string, string) {

	byteArr := []byte(m.Content)
	prefixByte := []byte(prefix)

	if byteArr[0] == prefixByte[0] {

		// Iterate untill space
		for index, val := range byteArr {

			noneArr := []byte(" ")

			// Create string from remaining
			if val == noneArr[0] {

				command := string(byteArr[1:index])
				arguments := byteArr[index:]

				strArgs, err := shellwords.Parse(string(arguments))
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "Could not parse command. Make sure you typed it correctly and filled all required parameters")
					return []string{}, "InvalidCommand"
				}

				return strArgs, command
			}
		}

	}
	return []string{}, "InvalidCommand"
}
