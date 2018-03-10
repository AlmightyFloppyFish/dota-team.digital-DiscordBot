package main

import (
	"encoding/gob"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type houseSetup struct {
	P5v5 []string
	P1v1 []string
	H1v1 []string
}

func viewTeam(s *discordgo.Session, m *discordgo.MessageCreate, commandArray []string, command string) {
	if command == "viewteam" {

		if commandArray[0] == "" {
			s.ChannelMessageSend(m.ChannelID, "`!viewteam daedalus/desolator/orchid/all`")
			return
		}

		var (
			daedalus  houseSetup
			desolator houseSetup
			orchid    houseSetup
		)

		switch commandArray[0] {
		case "all":
			// !viewteam all

			var err error

			err = decodeGob("teamdaedalus", &daedalus)
			err = decodeGob("teamdesolator", &desolator)
			err = decodeGob("teamorchid", &orchid)
			if err == nil {
				var (
					orchidText    string = "**__Team orchid__**\n**5V5**\n```\n" + orchid.P5v5[0] + "\n" + orchid.P5v5[1] + "\n" + orchid.P5v5[2] + "\n" + orchid.P5v5[3] + "\n" + orchid.P5v5[4] + "\n```\n**1V1:**\n```\n" + orchid.P1v1[0] + "\n" + orchid.P1v1[1] + "\n```\n**1V1 Heroes:**\n```\n" + orchid.H1v1[0] + "\n" + orchid.H1v1[1] + "```"
					daedalusText  string = "**__Team Daedalus__**\n**5V5**\n```\n" + daedalus.P5v5[0] + "\n" + daedalus.P5v5[1] + "\n" + daedalus.P5v5[2] + "\n" + daedalus.P5v5[3] + "\n" + daedalus.P5v5[4] + "\n```\n**1V1:**\n```\n" + daedalus.P1v1[0] + "\n" + daedalus.P1v1[1] + "\n```\n**1V1 Heroes:**\n```\n" + daedalus.H1v1[0] + "\n" + daedalus.H1v1[1] + "```"
					desolatorText string = "**__Team desolator__**\n**5V5**\n```\n" + desolator.P5v5[0] + "\n" + desolator.P5v5[1] + "\n" + desolator.P5v5[2] + "\n" + desolator.P5v5[3] + "\n" + desolator.P5v5[4] + "\n```\n**1V1:**\n```\n" + desolator.P1v1[0] + "\n" + desolator.P1v1[1] + "\n```\n**1V1 Heroes:**\n```\n" + desolator.H1v1[0] + "\n" + desolator.H1v1[1] + "```"
				)
				s.ChannelMessageSend(m.ChannelID, daedalusText+"\n\n"+desolatorText+"\n\n"+orchidText)
			}
		case "daedalus":
			// !viewteams deadalus
			err := decodeGob("teamdaedalus", &daedalus)
			if err == nil {
				var daedalusText string = "**__Team Daedalus__**\n**5V5**\n```\n" + daedalus.P5v5[0] + "\n" + daedalus.P5v5[1] + "\n" + daedalus.P5v5[2] + "\n" + daedalus.P5v5[3] + "\n" + daedalus.P5v5[4] + "\n```\n**1V1:**\n```\n" + daedalus.P1v1[0] + "\n" + daedalus.P1v1[1] + "\n```\n**1V1 Heroes:**\n```\n" + daedalus.H1v1[0] + "\n" + daedalus.H1v1[1] + "```"
				s.ChannelMessageSend(m.ChannelID, daedalusText)
			} else {
				fmt.Println(err)
			}
		case "desolator":
			// !viewteams desolator
			err := decodeGob("teamdesolator", &desolator)
			if err == nil {
				var desolatorText string = "**__Team desolator__**\n**5V5**\n```\n" + desolator.P5v5[0] + "\n" + desolator.P5v5[1] + "\n" + desolator.P5v5[2] + "\n" + desolator.P5v5[3] + "\n" + desolator.P5v5[4] + "\n```\n**1V1:**\n```\n" + desolator.P1v1[0] + "\n" + desolator.P1v1[1] + "\n```\n**1V1 Heroes:**\n```\n" + desolator.H1v1[0] + "\n" + desolator.H1v1[1] + "```"
				s.ChannelMessageSend(m.ChannelID, desolatorText)
			}
		case "orchid":
			// !viewteams orchid
			err := decodeGob("teamorchid", &orchid)
			if err == nil {
				var orchidText string = "**__Team orchid__**\n**5V5**\n```\n" + orchid.P5v5[0] + "\n" + orchid.P5v5[1] + "\n" + orchid.P5v5[2] + "\n" + orchid.P5v5[3] + "\n" + orchid.P5v5[4] + "\n```\n**1V1:**\n```\n" + orchid.P1v1[0] + "\n" + orchid.P1v1[1] + "\n```\n**1V1 Heroes:**\n```\n" + orchid.H1v1[0] + "\n" + orchid.H1v1[1] + "```"
				s.ChannelMessageSend(m.ChannelID, orchidText)

			}
		default:
			s.ChannelMessageSend(m.ChannelID, "`!viewteams daedalus/desolator/orchid/all`")
		}
	}
}

type random1v1setup struct {
	Game1 [3]string
	Game2 [3]string
	Game3 [3]string
}

func addToLeage(s *discordgo.Session, m *discordgo.MessageCreate, commandArray []string, command string) {
	if command == "addLeague" {

		userIDbytes := commandArray[0]
		userID := strings.TrimLeft(strings.TrimRight(userIDbytes, "> "), " <@!")

		s.ChannelMessageSend(channel.botchannel, "Attempting to add ID: __"+string(userID)+"__ To group `League Players`")

		err := s.GuildMemberRoleAdd(channel.guild, string(userID), "416133412013998080")
		if err != nil {
			s.ChannelMessageSend(channel.botchannel, "Something went wrong ):, Check bot console for details")
			fmt.Println(err)
			return
		} else {
			s.ChannelMessageSend(m.ChannelID, "User sucessfully gained role")
			return
		}
	}
	if command == "removeLeague" {

		userIDbytes := commandArray[0]
		userID := strings.TrimLeft(strings.TrimRight(userIDbytes, "> "), " <@!")

		s.ChannelMessageSend(channel.botchannel, "Attempting to remove ID: __"+string(userID)+"__ From group `League Players`")

		err := s.GuildMemberRoleRemove(channel.guild, string(userID), "416133412013998080")
		if err != nil {
			s.ChannelMessageSend(channel.botchannel, "Something went wrong ):, Check bot console for details")
			fmt.Println(err)
			return
		} else {
			s.ChannelMessageSend(m.ChannelID, "User sucessfully lost role")
		}
	}
	return
}

// I am well aware that copying large chunks like this is a bad habbit, got other projects i want to work on so this is a ghetto quick solution
func setupGame2(s *discordgo.Session, m *discordgo.MessageCreate, commandArray []string, command string) {

	if m.Content == "!setupgame2" || m.Content == "!setupgame2 " {

		s.ChannelMessageSend(m.ChannelID, "Atempting to open, "+"./teams/"+makeFolderName()+"/"+"teamdaedalus | teamdesolator.gob")

		var err error
		var (
			daed houseSetup
			orch houseSetup
		)

		// This is horrible code and I'll make a helper method in the future. Just lazy today
		err = decodeGob("teamdaedalus", &daed)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Could not read info from gob files, are both teams complete?; Aborting (check console for debug information)")
			fmt.Println(err)
			return
		}

		err = decodeGob("teamorchid", &orch)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Could not read info from gob files, are all both complete?; Aborting (check console for debug information)")
			fmt.Println(err)
			return
		}

		var randomNumber int
		var completeSetup random1v1setup

		// Randomization init
		rand.Seed(time.Now().UTC().UnixNano())
		randomNumber = rand.Intn(2)

		fmt.Println(completeSetup)

		// New pattern for second selection
		rand.Seed(time.Now().UTC().UnixNano())
		randomNumber = rand.Intn(2)

		// Decides setup for second player in deso depending on random number
		if randomNumber == 0 {
			fmt.Println(daed.P1v1[0])
			completeSetup.Game2[0] = daed.P1v1[1]
			completeSetup.Game2[1] = orch.P1v1[0]
			orch.P1v1 = append(orch.P1v1[:0], orch.P1v1[0+1:]...)
		} else if randomNumber == 1 {
			completeSetup.Game2[0] = daed.P1v1[1]
			completeSetup.Game2[1] = orch.P1v1[1]
			orch.P1v1 = append(orch.P1v1[:1], orch.P1v1[1+1:]...)
		}

		completeSetup.Game1[0] = daed.P1v1[0]
		completeSetup.Game1[1] = orch.P1v1[0]

		/* All players randomized for 1v1 by now */

		// All preffered heroes in one array
		var allHeroes []string
		allHeroes = append(allHeroes, daed.H1v1[0])
		allHeroes = append(allHeroes, daed.H1v1[1])
		allHeroes = append(allHeroes, orch.H1v1[0])
		allHeroes = append(allHeroes, orch.H1v1[1])

		// New random number for Hero randomization
		rand.Seed(time.Now().UTC().UnixNano())
		randomNumber = rand.Intn(4) // 0 -> 3

		completeSetup.Game1[2] = allHeroes[randomNumber]
		allHeroes = append(allHeroes[:randomNumber], allHeroes[randomNumber+1:]...) // Remove choosen hero from array of possibles

		rand.Seed(time.Now().UTC().UnixNano())
		randomNumber = rand.Intn(4) // 0 -> 3

		completeSetup.Game2[2] = allHeroes[randomNumber]
		allHeroes = append(allHeroes[:randomNumber], allHeroes[randomNumber+1:]...) // Remove choosen hero from array of possibles

		fmt.Println(completeSetup)

		// Save to gob file

		var fileerr error
		var gobfile *os.File

		// If gobfile doesn't exist, create it. Else open it
		if _, err := os.Stat("./teams/" + makeFolderName() + "/1v1lineup.gob"); os.IsNotExist(err) {
			gobfile, fileerr = os.Create("./teams/" + makeFolderName() + "/1v1lineup.gob")
			gobfile.Close()
			gobfile, fileerr = os.OpenFile("./teams/"+makeFolderName()+"/1v1lineup.gob", os.O_WRONLY, 0777)
		} else {
			gobfile, fileerr = os.OpenFile("./teams/"+makeFolderName()+"/1v1lineup.gob", os.O_WRONLY, 0777)
		}
		defer gobfile.Close()
		if fileerr == nil {
			gobencoder := gob.NewEncoder(gobfile)
			encerr := gobencoder.Encode(completeSetup)
			if encerr != nil {
				fmt.Println(encerr)
				return
			}
		} else {
			s.ChannelMessageSend(m.ChannelID, "1v1lineup.gob could not be made or read, aborting...")
			fmt.Println(fileerr)
			return
		}

		// Fancy output here to discord chat here
		s.ChannelMessageSend(channel.botchannel, "__**1 Versus 1 Lineups:**__\n```\n"+completeSetup.Game1[0]+" VS "+completeSetup.Game1[1]+"   |   "+completeSetup.Game1[2]+"\n"+completeSetup.Game2[0]+" VS "+completeSetup.Game2[1]+"   |   "+completeSetup.Game2[2]+"```")
	}
	return
}

func setupGame3(s *discordgo.Session, m *discordgo.MessageCreate, commandArray []string, command string) {

	if m.Content == "!setupgame3" || m.Content == "!setupgame3 " {

		s.ChannelMessageSend(m.ChannelID, "Atempting to open, "+"./teams/"+makeFolderName()+"/"+"teamorchid/teamdesolator/teamdaedalus.gob")

		var err error
		var (
			daed houseSetup
			deso houseSetup
			orch houseSetup
		)

		// This is horrible code and I'll make a helper method in the future. Just lazy today
		err = decodeGob("teamdaedalus", &daed)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Could not read info from gob files, are all 3 teams complete?; Aborting (check console for debug information)")
			fmt.Println(err)
			return
		}

		err = decodeGob("teamdesolator", &deso)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Could not read info from gob files, are all 3 teams complete?; Aborting (check console for debug information)")
			fmt.Println(err)
			return
		}

		err = decodeGob("teamorchid", &orch)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Could not read info from gob files, are all 3 teams complete?; Aborting (check console for debug information)")
			fmt.Println(err)
			return
		}

		var randomNumber int
		var completeSetup random1v1setup

		// Randomization init
		rand.Seed(time.Now().UTC().UnixNano())
		randomNumber = rand.Intn(2)

		fmt.Println(completeSetup)
		fmt.Println(deso.P1v1[1])

		// Decides setup for first player in deso depending on random number
		if randomNumber == 0 {
			completeSetup.Game1[0] = deso.P1v1[0]
			completeSetup.Game1[1] = daed.P1v1[0]
			daed.P1v1 = append(daed.P1v1[:0], daed.P1v1[0+1:]...) // Removes value from slice so the two remaining values can be assumed as last team
		} else if randomNumber == 1 {
			completeSetup.Game1[0] = deso.P1v1[0]
			completeSetup.Game1[1] = daed.P1v1[1]
			daed.P1v1 = append(daed.P1v1[:1], daed.P1v1[1+1:]...)
		}

		// New pattern for second selection
		rand.Seed(time.Now().UTC().UnixNano())
		randomNumber = rand.Intn(2)

		// Decides setup for second player in deso depending on random number
		if randomNumber == 0 {
			fmt.Println(deso.P1v1[0])
			completeSetup.Game2[0] = deso.P1v1[1]
			completeSetup.Game2[1] = orch.P1v1[0]
			orch.P1v1 = append(orch.P1v1[:0], orch.P1v1[0+1:]...)
		} else if randomNumber == 1 {
			completeSetup.Game2[0] = deso.P1v1[1]
			completeSetup.Game2[1] = orch.P1v1[1]
			orch.P1v1 = append(orch.P1v1[:1], orch.P1v1[1+1:]...)
		}

		// Matches the remaining player from orchid and daedalus
		completeSetup.Game3[0] = daed.P1v1[0]
		completeSetup.Game3[1] = orch.P1v1[0]

		/* All players randomized for 1v1 by now */

		// All preffered heroes in one array
		var allHeroes []string
		allHeroes = append(allHeroes, deso.H1v1[0])
		allHeroes = append(allHeroes, deso.H1v1[1])
		allHeroes = append(allHeroes, daed.H1v1[0])
		allHeroes = append(allHeroes, daed.H1v1[1])
		allHeroes = append(allHeroes, orch.H1v1[0])
		allHeroes = append(allHeroes, orch.H1v1[1])

		// New random number for Hero randomization
		rand.Seed(time.Now().UTC().UnixNano())
		randomNumber = rand.Intn(6) // 0 -> 5

		completeSetup.Game1[2] = allHeroes[randomNumber]
		allHeroes = append(allHeroes[:randomNumber], allHeroes[randomNumber+1:]...) // Remove choosen hero from array of possibles

		rand.Seed(time.Now().UTC().UnixNano())
		randomNumber = rand.Intn(5) // 0 -> 5

		completeSetup.Game2[2] = allHeroes[randomNumber]
		allHeroes = append(allHeroes[:randomNumber], allHeroes[randomNumber+1:]...) // Remove choosen hero from array of possibles

		rand.Seed(time.Now().UTC().UnixNano())
		randomNumber = rand.Intn(4) // 0 -> 5

		completeSetup.Game3[2] = allHeroes[randomNumber]
		allHeroes = append(allHeroes[:randomNumber], allHeroes[randomNumber+1:]...) // Remove choosen hero from array of possibles

		fmt.Println(completeSetup)

		// Save to gob file

		var fileerr error
		var gobfile *os.File

		// If gobfile doesn't exist, create it. Else open it
		if _, err := os.Stat("./teams/" + makeFolderName() + "/1v1lineup.gob"); os.IsNotExist(err) {
			gobfile, fileerr = os.Create("./teams/" + makeFolderName() + "/1v1lineup.gob")
			gobfile.Close()
			gobfile, fileerr = os.OpenFile("./teams/"+makeFolderName()+"/1v1lineup.gob", os.O_WRONLY, 0777)
		} else {
			gobfile, fileerr = os.OpenFile("./teams/"+makeFolderName()+"/1v1lineup.gob", os.O_WRONLY, 0777)
		}
		defer gobfile.Close()
		if fileerr == nil {
			gobencoder := gob.NewEncoder(gobfile)
			encerr := gobencoder.Encode(completeSetup)
			if encerr != nil {
				fmt.Println(encerr)
				return
			}
		} else {
			fmt.Println(fileerr)
			return
		}

		// Fancy output here to discord chat here
		s.ChannelMessageSend(channel.botchannel, "__**1 Versus 1 Lineups:**__\n```\n"+completeSetup.Game1[0]+" VS "+completeSetup.Game1[1]+"   |   "+completeSetup.Game1[2]+"\n"+completeSetup.Game2[0]+" VS "+completeSetup.Game2[1]+"   |   "+completeSetup.Game2[2]+"\n"+completeSetup.Game3[0]+" VS "+completeSetup.Game3[1]+"   |   "+completeSetup.Game3[2]+"\n```")
	}
	return
}

func lineup(s *discordgo.Session, m *discordgo.MessageCreate, commandArray []string, command string) {
	if command == "teamdaedalus" || command == "teamdesolator" || command == "teamorchid" {

		var fullTeam houseSetup

		if len(commandArray) == 3 {
			fullTeam.P5v5 = strings.Split(commandArray[0], ", ")
			fullTeam.P1v1 = strings.Split(commandArray[1], ", ")
			fullTeam.H1v1 = strings.Split(commandArray[2], ", ")

			// Check if slices correct length
			if len(fullTeam.P5v5) == 5 && len(fullTeam.P1v1) == 2 && len(fullTeam.H1v1) == 2 {

				// Makes folder for entire housematch
				if _, err := os.Stat("./teams/" + makeFolderName()); os.IsNotExist(err) {
					os.MkdirAll("./teams/"+makeFolderName(), 0777)
				}

				var err error

				switch command {
				case "teamdesolator":
					// Encoding for deso
					err = encodeToGob("teamdesolator", fullTeam)

				case "teamdaedalus":
					// Encoding for daedalus
					err = encodeToGob("teamdaedalus", fullTeam)
				case "teamorchid":
					// Encoding for Orchid
					err = encodeToGob("teamorchid", fullTeam)
				}

				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "Could not encode to gob file. Your team layout could not be saved")
					fmt.Println(err)
					return
				} else {
					s.ChannelMessageSend(m.ChannelID, "Team layout saved")
				}

			} else {
				s.ChannelMessageSend(m.ChannelID, "Incorrect number of players or heroes. Expected 5, 2, 2, got: "+strconv.Itoa(len(fullTeam.P5v5))+", "+strconv.Itoa(len(fullTeam.P1v1))+", "+strconv.Itoa(len(fullTeam.H1v1)))
				return
			}
		}
	}
	return
}

func viewgame(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Content == "!viewgame" {

		gobfile, err := os.OpenFile("./teams/"+makeFolderName()+"/1v1lineup.gob", os.O_RDONLY, 0777)
		defer gobfile.Close()

		var completeSetup random1v1setup

		if err == nil {
			gobdecoder := gob.NewDecoder(gobfile)
			err := gobdecoder.Decode(&completeSetup)
			if err != nil {
				fmt.Println(err)
				return
			}
		} else {
			fmt.Println(err)
			return
		}

		s.ChannelMessageSend(m.ChannelID, "__**1 Versus 1 Lineups:**__\n```\n"+completeSetup.Game1[0]+" VS "+completeSetup.Game1[1]+"   |   "+completeSetup.Game1[2]+"\n"+completeSetup.Game2[0]+" VS "+completeSetup.Game2[1]+"   |   "+completeSetup.Game2[2])

	}
	return
}

func makeFolderName() string {
	timenow := time.Now()

	year := timenow.Year()
	day := timenow.Day()
	month := timenow.Month()

	out := strconv.Itoa(year) + month.String() + strconv.Itoa(day)
	return out
}

func encodeToGob(filename string, obj interface{}) error {

	var fileerr error
	var gobfile *os.File

	if _, err := os.Stat("./teams/" + makeFolderName() + "/" + filename + ".gob"); os.IsNotExist(err) {
		gobfile, fileerr = os.Create("./teams/" + makeFolderName() + "/" + filename + ".gob")
		gobfile.Close()
		gobfile, fileerr = os.OpenFile("./teams/"+makeFolderName()+"/"+filename+".gob", os.O_WRONLY, 0777)
	} else {
		gobfile, fileerr = os.OpenFile("./teams/"+makeFolderName()+"/"+filename+".gob", os.O_WRONLY, 0777)
	}

	defer gobfile.Close()
	if fileerr == nil {
		fmt.Println("Was gonna encode: ")
		fmt.Println(obj)
		gobencoder := gob.NewEncoder(gobfile)
		encerr := gobencoder.Encode(obj)
		if encerr != nil {
			return encerr
		}
	} else {
		return fileerr
	}
	return fileerr
}

// Write function to decode
func decodeGob(filename string, madeData *houseSetup) error {

	gobfile, err := os.OpenFile("./teams/"+makeFolderName()+"/"+filename+".gob", os.O_RDONLY, 0777)
	defer gobfile.Close()

	if err == nil {
		gobdecoder := gob.NewDecoder(gobfile)
		err := gobdecoder.Decode(&madeData)
		return err
	} else {
		return err
	}
	return err
}
