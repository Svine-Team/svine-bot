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

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name: "basic-command",
			// All commands and options must have a description
			// Commands/options without description will fail the registration
			// of the command.
			Description: "Basic command",
		},
		{
			Name: "cool-basic-command",
			// All commands and options must have a description
			// Commands/options without description will fail the registration
			// of the command.
			Description: "Basic command",
		},
    }

    commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
        "basic-command": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
            s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
                Type: discordgo.InteractionResponseChannelMessageWithSource,
                Data: &discordgo.InteractionResponseData{
                    Content: "Response!",
                },
            })
        },

        "cool-basic-command": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
            s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
                Type: discordgo.InteractionResponseChannelMessageWithSource,
                Data: &discordgo.InteractionResponseData{
                    Content: "COOL Response!",
                },
            })
        },
    }
)

var session *discordgo.Session

func initHandlers() {
    session.AddHandler(func(s *discordgo.Session, ready *discordgo.Ready) {
        user := s.State.User
        log.Printf("Logged in as %v#%v", user.Username, user.Discriminator)
    })

    session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
        if commandHandler, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
            commandHandler(s, i)
        }
    })

    // Register the messageCreate func as a callback for MessageCreate events.
    session.AddHandler(messageCreate)

    // In this example, we only care about receiving message events.
    session.Identify.Intents = discordgo.IntentsGuildMessages

}

// Create a new Discord session using the provided bot token.
func init() {
    Token := getEnvVariable("BOT_TOKEN")

    var err error

    session, err = discordgo.New("Bot " + Token)
    if err != nil {
        log.Fatalf("Cannot create the session: %v", err)
    }

    initHandlers()

}

type Commands = []*discordgo.ApplicationCommand

func registerCommands() Commands {
    registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
    for i, command := range commands {
        registeredCommand, err := session.ApplicationCommandCreate(session.State.User.ID, "", command)
        if err != nil {
            log.Panicf("Cannot create '%v' command: %v", command.Name, err)
        }
        registeredCommands[i] = registeredCommand
    }

    return registeredCommands
}

func removeRegisteredCommands(registeredCommands Commands) {
    for _, command := range registeredCommands {
        err := session.ApplicationCommandDelete(session.State.User.ID, "", command.ID)
        if err != nil {
            log.Panicf("Cannot delete '%v' command: %v", command.Name, err)
        }
    }
}

func main() {
    // Open a websocket connection to Discord and begin listening.
    err := session.Open()
    if err != nil {
        log.Fatalf("Cannot open the session: %v", err)
        return
    }

    log.Println("Adding commands...")
    registeredCommands := registerCommands()

    // Wait here until CTRL-C or other term signal is received.
    fmt.Println("Bot is now running.  Press CTRL-C to exit.")
    sc := make(chan os.Signal, 1)
    signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
    <-sc

    fmt.Println("Closing session...")
    // Cleanly close down the Discord session.
    session.Close()

    log.Println("Removing commands...")
    removeRegisteredCommands(registeredCommands)
    fmt.Println("Gracefully shutting down...")
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
