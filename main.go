package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

var token string
var fileName string
var isPlaying bool = false

type CommandData struct {
	FileName string
	Message  string
	Type     string
}

var payload map[string]CommandData

var helpCommand string

const infoCommand = `
I'm the intellectual brainchild of  Frooster. My code can be found at https://github.com/vathzen/EmoBot.
For any issues and recommendations please contact my author at https://vathzen.in/discord; though he will probably say Vaaila Vechuko`

func getDataFromJson() {
	payload = nil
	fileContent, err := ioutil.ReadFile("./commands.json")
	if err != nil {
		return
	}
	err = json.Unmarshal(fileContent, &payload)
	if err != nil {
		return
	}
}

func main() {

	token = os.Getenv("EMO_TOKEN")

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	updateHelpCommand()

	dg.AddHandler(ready)
	dg.AddHandler(messageCreate)
	dg.AddHandler(guildCreate)
	dg.AddHandler(onKTVJoin)

	dg.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
	}

	fmt.Println("Emobot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	s.UpdateGameStatus(0, "!emo")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	currentTime := time.Now()

	if strings.HasPrefix(m.Content, "!emo") {

		// Find the channel that the message came from.
		c, err := s.State.Channel(m.ChannelID)
		if err != nil {
			return
		}

		// Find the guild for that channel.
		g, err := s.State.Guild(c.GuildID)
		if err != nil {
			return
		}

		voiceLine := strings.Split(m.Content, " ")

		if voiceLine[1] == "help" {
			s.ChannelMessageSend(m.ChannelID, helpCommand)
			s.MessageReactionAdd(c.ID, m.ID, "\U0001F44D")
			return
		} else if voiceLine[1] == "info" {
			s.ChannelMessageSend(m.ChannelID, infoCommand)
			s.MessageReactionAdd(c.ID, m.ID, "\U0001F44D")
			return
		} else {
			// Look for the message sender in that guild's current voice states.
			for _, vs := range g.VoiceStates {
				if vs.UserID == m.Author.ID {

					if len(voiceLine) <= 1 {
						return
					}

					getDataFromJson()

					fileName = payload[voiceLine[1]].FileName
					if fileName == "" {
						fmt.Printf("%s %s: %s -> %s Not Found\n", currentTime.Format("2006.01.02 15:04:05"), g.Name, m.Author, voiceLine[1])
						s.MessageReactionAdd(c.ID, m.ID, "\U00002639")
						return
					}

					if isPlaying == false {
						isPlaying = true
						fmt.Printf("%s %s: %s -> %s\n", currentTime.Format("2006.01.02 15:04:05"), g.Name, m.Author, voiceLine[1])
						s.MessageReactionAdd(c.ID, m.ID, "\U0001F602")
						err = playSound(s, g.ID, vs.ChannelID)
					}
					isPlaying = false

					if err != nil {
						fmt.Println("Error Playing sound:", err)
					}
					return
				}
			}
		}
	}

	if strings.HasPrefix(m.Content, "!d2") {

		c, err := s.State.Channel(m.ChannelID)
		if err != nil {
			return
		}
		g, err := s.State.Guild(c.GuildID)
		if err != nil {
			return
		}
		for _, vs := range g.VoiceStates {
			if vs.UserID == m.Author.ID {
				voiceLine := strings.Split(m.Content, " ")
				if len(voiceLine) <= 1 {
					return
				}

				getDataFromJson()

				fileName = payload[voiceLine[1]].FileName

				if isPlaying == false {
					isPlaying = true
					fmt.Printf("%s: %s -> %s\n", g.Name, m.Author, voiceLine[1])
					err = playSound(s, g.ID, vs.ChannelID)
				}
				isPlaying = false

				if err != nil {
					fmt.Println("Error Playing sound:", err)
				}
				return
			}
		}

	}
}

func guildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {

	if event.Guild.Unavailable {
		return
	}

	for _, channel := range event.Guild.Channels {
		if channel.ID == event.Guild.ID {
			_, _ = s.ChannelMessageSend(channel.ID, "Emobot is ready! Type !emo while in a voice channel to play a sound.")
			return
		}
	}
}

func onKTVJoin(s *discordgo.Session, event *discordgo.VoiceStateUpdate) {
	//KTV ID = 963324557995962398
	// My ID = 300626364418097162
	//Vas ID = 214451937322467330

	getDataFromJson()

	fileName = payload["padida"].FileName
	currentTime := time.Now()

	g, err := s.State.Guild(event.GuildID)
	if err != nil {
		return
	}

	if event.VoiceState.UserID == "214451937322467330" && event.BeforeUpdate == nil {

		fmt.Printf("Vaseey Joined : %s at %s \n", g.Name, currentTime.Format("2006.01.02 15:04:05"))
		fileName = payload["vaseey"].FileName

		if isPlaying == false {
			isPlaying = true
			err := playSound(s, event.GuildID, event.VoiceState.ChannelID)
			if err != nil {
				fmt.Print(err)
			}
		}
		isPlaying = false
	}

	// if event.VoiceState.UserID == "963324557995962398" && event.BeforeUpdate == nil {

	// 	fmt.Printf("KTV Joined : %s at %s \n", g.Name, currentTime.Format("2006.01.02 15:04:05"))
	// 	fileName = payload["padida"].FileName

	// 	if isPlaying == false {
	// 		isPlaying = true
	// 		err := playSound(s, event.GuildID, event.VoiceState.ChannelID)
	// 		if err != nil {
	// 			fmt.Print(err)
	// 		}
	// 	}
	// 	isPlaying = false
	// }
}

// loadSound attempts to load an encoded sound file from disk.
func loadSound() (buffer2 [][]uint8, err error) {

	var buffer = make([][]byte, 0)

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error opening dca file :", err)
		return nil, err
	}

	var opuslen int16

	for {
		// Read opus frame length from dca file.
		err = binary.Read(file, binary.LittleEndian, &opuslen)

		// If this is the end of the file, just return.
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err := file.Close()
			if err != nil {
				return nil, err
			}
			return buffer, nil
		}

		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return nil, err
		}

		// Read encoded pcm from dca file.
		InBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &InBuf)

		// Should not be any end of file errors
		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return nil, err
		}
		// Append encoded pcm data to the buffer.
		buffer = append(buffer, InBuf)
	}
}

// playSound plays the current buffer to the provided channel.
func playSound(s *discordgo.Session, guildID, channelID string) (err error) {

	buffer, err2 := loadSound()

	if err2 != nil {
		return err
	}

	// Join the provided voice channel.
	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		if _, ok := s.VoiceConnections[guildID]; ok {
			vc = s.VoiceConnections[guildID]
		} else {
			return nil
		}
	}

	// Start speaking.
	vc.Speaking(true)

	// Send the buffer data.
	for _, buff := range buffer {
		vc.OpusSend <- buff
	}

	buffer = nil

	// Stop speaking
	vc.Speaking(false)

	// Disconnect from the provided voice channel.
	vc.Disconnect()

	return nil
}

func updateHelpCommand() {
	var emoCommands = `
!emo commands
-------------
info: Learn more about the bot
`

	var d2Commands = `
!d2 commands
-------------
`

	getDataFromJson()

	for key, element := range payload {
		if element.Type == "emo" {
			emoCommands += key + ": " + element.Message + "\n"
		} else if element.Type == "d2" {
			d2Commands += key + ": " + element.Message + "\n"
		}

	}

	helpCommand = emoCommands + d2Commands
}
