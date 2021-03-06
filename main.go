package main

import (
	"encoding/gob"
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

	go checklowprio(discord)

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
		fmt.Println("  -----  WARNING Session is being shut down! (This has been done intentionally)")
		os.Exit(3)
	}

	// Causes bot to ignore itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// colonArray := ColonSplit(s, m)
	argumentArray, command := easySplit(s, m)

	// Launch commands here. // if statement here checks if the channel is botchannel
	if m.ChannelID == channel.botchannel {

		go joinme(s, m, command)
		go lifecheck(s, m)
		go help(s, m)
		go spam(s, m, argumentArray, command)
		go music(s, m, argumentArray, command)
		go stop(s, m, argumentArray, command)
		go lineup(s, m, argumentArray, command)
		go setupGame2(s, m, argumentArray, command)
		go setupGame3(s, m, argumentArray, command)
		go addToLeage(s, m, argumentArray, command)
		go lowprio(s, m, argumentArray, command)
	}

	callingUser := m.Author

	if m.ChannelID == "415295523927621632" {

		argumentArray, command := easySplit(s, m)
		if callingUser.ID == "179945962490429451" || callingUser.ID == "220120722939445248" || callingUser.ID == "242788102928596992" {
			go addToLeage(s, m, argumentArray, command)
		}
	}
	callingMember, _ := s.GuildMember(channel.guild, callingUser.ID)

	for _, v := range callingMember.Roles {
		if v == "281077340845506562" || callingUser.ID == "220120722939445248" || callingUser.ID == "242788102928596992" {
			go viewgame(s, m)
			go viewTeam(s, m, argumentArray, command)
			return
		}
	}
	return
}

//
// ------------------------------------------------------- //
//

type LowPrioPlayerInfo struct {
	ID          string
	TimeCreated int64
	Duration    int64
}

func checklowprio(s *discordgo.Session) {

	for {

		allFiles, _ := ioutil.ReadDir("./lowprio/")

		for _, v := range allFiles {

			fileName := v.Name()
			InFile, err := os.Open("./lowprio/" + fileName)
			if err != nil {
				fmt.Println(err)
				return
			}

			dec := gob.NewDecoder(InFile)

			var Res LowPrioPlayerInfo
			err2 := dec.Decode(&Res)
			if err2 != nil {
				s.ChannelMessageSend(channel.botchannel, "Unable to decode gob data from ./lowprio/"+fileName)
				fmt.Println(err2)
				return
			}

			currentTime := time.Now().Unix()

			if Res.TimeCreated <= currentTime-Res.Duration {
				s.ChannelMessageSend(channel.botchannel, "Removing player "+Res.ID+" from low priority")
				s.GuildMemberRoleRemove(channel.guild, Res.ID, "286629932996493313")
				InFile.Close()
				os.Remove("./lowprio/" + fileName)
			}

			InFile.Close()
		}
		time.Sleep(1 * time.Minute)
	}
}

func lowprio(s *discordgo.Session, m *discordgo.MessageCreate, commandArray []string, command string) {

	if command == "lowprio" {

		callingMember, _ := s.GuildMember(channel.guild, m.Author.ID)

		for _, v := range callingMember.Roles {
			if v == "277034907803582464" {

				fmt.Println(commandArray[0])

				_, err := strconv.Atoi(commandArray[0])
				if err != nil {
					s.ChannelMessageSend(channel.botchannel, "That doesn't look like a persons ID to me. I'm looking for something like !lowprio *37285847386918486* 40")
					return
				}
				duration, err := strconv.Atoi(commandArray[1])
				if err != nil {
					s.ChannelMessageSend(channel.botchannel, "That doesn't look like a duration in minutes to me. I'm looking for something like !lowprio *54375982347598374* 40")
				}

				err2 := s.GuildMemberRoleAdd(channel.guild, commandArray[0], "286629932996493313")
				if err2 != nil {
					s.ChannelMessageSend(channel.botchannel, "Something went wrong ):, Check bot console for details")
					fmt.Println(err2)
					return
				} else {
					s.ChannelMessageSend(channel.botchannel, "User sucessfully gained role for"+commandArray[1]+"minutes")

					// Timer save stuff here
					if _, err := os.Stat("./lowprio/" + commandArray[0]); os.IsNotExist(err) {
						os.Create("./lowprio/" + commandArray[0])
					} else {
						os.OpenFile("./lowprio"+commandArray[0], os.O_WRONLY, 0777)
					}
					var enc *gob.Encoder
					OutFile, err := os.OpenFile("./lowprio/"+commandArray[0], os.O_WRONLY, 0777)
					if err == nil {
						enc = gob.NewEncoder(OutFile)
					}
					defer OutFile.Close()

					timeMade := time.Now().Unix()

					OutData := LowPrioPlayerInfo{
						ID:          commandArray[0],
						Duration:    int64(duration * 60),
						TimeCreated: timeMade,
					}

					err3 := enc.Encode(OutData)
					if err3 != nil {
						s.ChannelMessageSend(channel.botchannel, "I was unable to encode the data for user "+OutData.ID+" this is bad")
						return
					}

					return
				}
			}
		}
	}
}

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
