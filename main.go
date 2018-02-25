package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const (
	debug  bool   = false
	prefix string = "!"
)

var (
	stopMusicChannel = make(chan bool)
	stopSpamChannel  = make(chan bool)

	channel channelInfoType

	outputChannel string = "0"
)

type channelInfoType struct {
	token          string
	guild          string
	botchannel     string
	generalchannel string
	teamdaedelus   string
	teamdesolator  string
	teamorchid     string
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
			guild:          "378651898519093248",
			botchannel:     "410195120135077889",
			generalchannel: "378651898519093250",
		}
	} else {
		channel = channelInfoType{
			token:          string(digTok),
			guild:          "242808761641730048",
			botchannel:     "285511576700846081",
			generalchannel: "285922016534462464",
			teamdaedelus:   "379528439415177217",
			teamdesolator:  "379528393064054785",
			teamorchid:     "408643128346673153",
		}
	}

	// Log into discord
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

		// colonArray := ColonSplit(s, m)
		argumentArray, command := easySplit(s, m)

		go joinme(s, m, command)
		go lifecheck(s, m)
		go help(s, m)
		go spam(s, m, argumentArray, command)
		go music(s, m, argumentArray, command)
		go stop(s, m, argumentArray, command)
		go lineup(s, m, argumentArray, command)
		go viewTeam(s, m, argumentArray, command)
		go setupGame(s, m, argumentArray, command)
		go addToLeage(s, m, argumentArray, command)
	}
	go viewgame(s, m)
}

//
// ------------------------------------------------------- //
//

func stop(s *discordgo.Session, m *discordgo.MessageCreate, commandArray []string, command string) {
	if command == "stop" {
		if commandArray[0] == "music" {
			stopMusicChannel <- true
		}
		if commandArray[0] == "spam" {
			s.ChannelMessageSend(m.ChannelID, "Spam stopped")
			stopSpamChannel <- true
		}
	}
}

func help(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Content == "!help" || m.Content == "!help " {
		helpMessage, err := ioutil.ReadFile("./helpMessage.txt")
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Floppy has forgotten to add the helpMessage.txt file so please go spam him so he fixes.")
			return
		}
		s.ChannelMessageSend(m.ChannelID, string(helpMessage))
	}
}

func lifecheck(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Content == "areyoualive?" {
		s.ChannelMessageSend(m.ChannelID, "Yes, I'm alive and working")

	}
	return
}

func spam(s *discordgo.Session, m *discordgo.MessageCreate, commandArray []string, command string) {

	if command == "spam" {

		if len(commandArray) == 4 {
			if commandArray[3] == "botchannel" {
				outputChannel = channel.botchannel
			} else if commandArray[3] == "generalchat" {
				outputChannel = channel.generalchannel
			} else {
				s.ChannelMessageSend(m.ChannelID, "Channel could not be found. Make sure you spelled it correctly")
				return
			}

			wait_sec, err := strconv.ParseFloat(commandArray[2], 64)

			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Interval couldn't be read correctly. Check your command ")
				return
			}

			for x, err := strconv.Atoi(commandArray[1]); x > 0; x -= 1 {
				select {
				case <-stopSpamChannel:
					return
				default:
					if err != nil {
						s.ChannelMessageSend(m.ChannelID, "amount could not be converted into array, Double check your command.")
						return
					}
					s.ChannelMessageSend(outputChannel, commandArray[0])
					time.Sleep(time.Duration(wait_sec) * time.Second * 60)
					if x < 1 {
						s.ChannelMessageSend(channel.botchannel, "I finished spamming "+outputChannel)
						return
					}

				}
			}

		} else {
			s.ChannelMessageSend(m.ChannelID, "Wrong amount of arguments, use syntax:  ```!spam \"Message here\" 2 60 generalchat``` where 2 is amount of messages and 60 is minutes inbetween")
		}

	}
}

func joinme(s *discordgo.Session, m *discordgo.MessageCreate, command string) {
	if command == "joinme" {
		callinguser := m.Author

		fmt.Println(callinguser)
		fmt.Println(callinguser.ID)

		_, err := joinUserVoiceChannel(s, callinguser.ID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Could not join your voice channel, Maybe i don't have permission to?")
			return
		}
	}
}
