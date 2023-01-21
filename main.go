package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// declare our secret token string
var (
	Token string
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
}

func main() {
	// start a new discord go bot session using the token we declared
	dg, err := discordgo.New("Bot " + "authentication token")
	if err != nil {
		fmt.Println("error creating Discord session: %w", err)
	}

	// Register the messageCreate func as a callback for MessageCreate events
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection: %w", err)
	}

	// Wait unti term signal such as CTRL-C or other as specified is pressed
	fmt.Println("Bot is now running. Presss CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session
	dg.Close()
}

// This function will be called everytime a message is created on any channel that the authenticated bot has access to
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore all messages created by the bot
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!agents" {

		// call the valo api to return a list of playable agents //TODO Separate all APIs to new file, this will get messy otherwise
		req, err := http.NewRequest("GET", "https://valorant-api.com/v1/agents?isPlayableCharacter=true", nil)
		if err != nil {
			fmt.Println(err)
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalln(err)
		}

		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatalln(err)
			}

			// _, err = s.ChannelMessageSend(m.ChannelID, string(b))
			fmt.Println("response body: ", string(b))

			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("error getting agent data with status code %w", resp.StatusCode)
		}
	}

}
