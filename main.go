package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	debug bool = true
)

var (
	stopMusicChannel = make(chan bool)
	stopSpamChannel  = make(chan bool)

	channel channelInfoType

	outputChannel string = "0"
)

type channelInfoType struct {
	token          string
	botchannel     string
	generalchannel string
}

func main() {

	debTok, err := ioutil.ReadFile("./debugToken.txt")
	digTok, err := ioutil.ReadFile("./digitalToken.txt")

	if err != nil {
		panic("Could not open token.txt file")
	}

	// Dirty hack to counter a strange bug with the ReadFile method
	if debug {
		debTok = debTok[0 : len(debTok)-1]
	} else {
		digTok = digTok[0 : len(debTok)-1]
	}

	if debug {
		channel = channelInfoType{
			token:          string(debTok),
			botchannel:     "410195120135077889",
			generalchannel: "410195120135077889",
		}
	} else {
		channel = channelInfoType{
			token:          string(digTok),
			botchannel:     "285511576700846081",
			generalchannel: "285922016534462464",
		}
	}

	discord, err := discordgo.New(channel.token)
	defer discord.Close()

	if err != nil {
		fmt.Println("The discord session returned an error. Maybe incorrect auth token?")
		return
	}

	// Runs function whenMessage() at every discord message.
	discord.AddHandler(whenMessage)

	fmt.Println("Loading...")

	err = discord.Open()
	if err != nil {
		fmt.Println("I was unable to connect to discord. ,", err)
		return
	} else {
		fmt.Println("Session opened. Bot should now apear as a user")
	}

	// Checks for task kill instruction (like ^C)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

// See "discord.AddHandler(messageCreate)"
func whenMessage(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Content == "restart" {
		s.ChannelMessageSend(m.ChannelID, "*Restarting bot, Should only take a few seconds.*")
		s.Close()
		fmt.Println("Session is being shut down! (This has been done intentionally)")
		os.Exit(3)
	}

	// Causes bot to ignore itself
	if m.Author.ID == s.State.User.ID {
		return
	}
	// Launch commands here. // if statement here checks if the channel is botchannel
	if m.ChannelID == channel.botchannel {

		// Splits the words into string array on ' : '
		commandArray := strings.Split(m.Content, " : ")

		go joinme(s, m, commandArray)
		go lifecheck(s, m, commandArray)
		go help(s, m, commandArray)
		go spam(s, m, commandArray)
		go music(s, m, commandArray)
		go stop(s, m, commandArray)
	}
}

//
// ------------------------------------------------------- //
//

func stop(s *discordgo.Session, m *discordgo.MessageCreate, commandArray []string) {
	if commandArray[0] == "stop" {
		if commandArray[1] == "music" {
			stopMusicChannel <- true
		}
		if commandArray[1] == "spam" {
			s.ChannelMessageSend(m.ChannelID, "Spam stopped")
			stopSpamChannel <- true
		}
	}
}

func help(s *discordgo.Session, m *discordgo.MessageCreate, commandArray []string) {

	if commandArray[0] == "help" {
		helpMessage, err := ioutil.ReadFile("./helpMessage.txt")
		if err != nil {
			panic("helpMessage.txt could not be read")
		}
		s.ChannelMessageSend(m.ChannelID, string(helpMessage))
	}
}

func lifecheck(s *discordgo.Session, m *discordgo.MessageCreate, commandArray []string) {
	if m.Content == "areyoualive?" {
		s.ChannelMessageSend(m.ChannelID, "Yes, I'm alive and working")

	}
}

func spam(s *discordgo.Session, m *discordgo.MessageCreate, commandArray []string) {

	defer recover()

	if commandArray[0] == "spam" {

		if commandArray[4] == "botchannel" {
			outputChannel = channel.botchannel
		} else if commandArray[4] == "generalchat" {
			outputChannel = channel.generalchannel
		} else {
			s.ChannelMessageSend(m.ChannelID, "Channel could not be found. Make sure you spelled it correctly")
		}

		wait_sec, err := strconv.ParseFloat(commandArray[3], 64)

		if err != nil {
			panic(err)
			s.ChannelMessageSend(m.ChannelID, "Interval couldn't be read correctly. Check your command ")
		}

		for x, err := strconv.Atoi(commandArray[2]); x > 0; x -= 1 {
			select {
			case <-stopSpamChannel:
				return
			default:
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "amount could not be converted into array, Double check your command.")
					panic(err)
					return
				}
				s.ChannelMessageSend(outputChannel, commandArray[1])
				time.Sleep(time.Duration(wait_sec) * time.Second * 60)
				if x < 1 {
					return
				}

			}
		}
	}
}

func joinme(s *discordgo.Session, m *discordgo.MessageCreate, commandArray []string) {
	if m.Content == "joinme" {
		callinguser := m.Author

		fmt.Println(callinguser)
		fmt.Println(callinguser.ID)

		_, err := joinUserVoiceChannel(s, callinguser.ID)
		if err != nil {
			panic(err)
			s.ChannelMessageSend(m.ChannelID, "Could not join your voice channel, Maybe i don't have permission to?")
		}
	}
}
