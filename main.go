package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func init() {
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()
}

var token string
var command int
var isPlaying bool = false

func main() {

	if token == "" {
		fmt.Println("Pass Token as Param")
		return
	}

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}
	dg.AddHandler(ready)
	dg.AddHandler(messageCreate)
	dg.AddHandler(guildCreate)

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

		// Look for the message sender in that guild's current voice states.
		for _, vs := range g.VoiceStates {
			if vs.UserID == m.Author.ID {
				voiceLine := strings.Split(m.Content, " ")
				if len(voiceLine) <= 1 {
					return
				}
				switch voiceLine[1] {
				case "emo":
					command = 1
				case "theri":
					command = 2
				case "aiyo":
					command = 3
				case "iladi":
					command = 4
				case "wtf":
					command = 5
				case "davara":
					command = 6
				case "daedalus":
					command = 7
				default:
					return
				}

				if isPlaying == false {
					isPlaying = true
					fmt.Printf("%s: %s -> %s", g.Name, m.Author, voiceLine[1])
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
				switch voiceLine[1] {
				case "ratata":
					command = 51
				case "die":
					command = 52
				case "bock":
					command = 54
				case "sad":
					command = 53
				default:
					return
				}

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

// loadSound attempts to load an encoded sound file from disk.
func loadSound() (buffer2 [][]uint8, err error) {

	var filename string
	var buffer = make([][]byte, 0)

	switch command {
	case 1:
		filename = "emodmg.dca"
	case 2:
		filename = "theri.dca"
	case 3:
		filename = "aiyayo.dca"
	case 4:
		filename = "siruthai.dca"
	case 5:
		filename = "wtf.dca"
	case 6:
		filename = "davara.dca"
	case 7:
		filename = "daedalus3.dca"

	case 51:
		filename = "ratata.dca"
	case 52:
		filename = "imightdie.dca"
	case 53:
		filename = "sad.dca"
	case 54:
		filename = "bock.dca"
	}

	file, err := os.Open(filename)
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
		return err
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
