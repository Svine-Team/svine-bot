package main

import (
    "fmt"
    "os"
    "log"
    "os/signal"
    "syscall"

    "github.com/joho/godotenv"
    "github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
    Token string
)

// Use godot package to load/read the .env file and
//   return the value of the key.
func getEnvVariable(key string) string {

    // load .env file
    err := godotenv.Load()

    if err != nil {
        log.Fatalf("Error loading .env file")
    }

    return os.Getenv(key)
}

func main() {

    Token := getEnvVariable("BOT_TOKEN")

    // Create a new Discord session using the provided bot token.
    dg, err := discordgo.New("Bot " + Token)
    if err != nil {
        fmt.Println("error creating Discord session,", err)
        return
    }

    // Register the messageCreate func as a callback for MessageCreate events.
    dg.AddHandler(messageCreate)

    // In this example, we only care about receiving message events.
    dg.Identify.Intents = discordgo.IntentsGuildMessages

    // Open a websocket connection to Discord and begin listening.
    err = dg.Open()
    if err != nil {
        fmt.Println("error opening connection,", err)
        return
    }

    // Wait here until CTRL-C or other term signal is received.
    fmt.Println("Bot is now running.  Press CTRL-C to exit.")
    sc := make(chan os.Signal, 1)
    signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
    <-sc

    // Cleanly close down the Discord session.
    dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
    // fmt.Println(m.Message)
    // Ignore all messages created by the bot itself
    // This isn't required in this specific example but it's a good practice.
    if m.Author.ID == s.State.User.ID {
        return
    }

    // If the message is "ping" reply with "Pong!"
    if m.Content == "ping" {
        s.ChannelMessageSend(m.ChannelID, "Pong!")
    }

    // If the message is "pong" reply with "Ping!"
    if m.Content == "pong" {
        s.ChannelMessageSend(m.ChannelID, "Ping!")
    }
}
