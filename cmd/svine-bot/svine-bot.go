package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Svine-Team/svine-bot/internal/commands"
	"github.com/Svine-Team/svine-bot/internal/messages"
	"github.com/Svine-Team/svine-bot/internal/roles"
	"github.com/Svine-Team/svine-bot/pkg/env"
	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	token   string
	session *discordgo.Session
)

// Create a new Discord session using the provided bot token.
func init() {
	// TODO:https://github.com/cosmtrek/air/issues?q=is%3Aissue+log+is%3Aopen
	// Unfortunately, when launching app with `air`, everything that goes
	//   into stderr is held back until terminate. So you don't see any `log`
	//   messages in console and even redirecting them via `*>` or `2>` don't
	//   help.
	// Therefore, I set log to stdout: now it functions almost the same as
	// `fmt` but it's still important to differentiate between there two
	// packages. Use `log` when you want to log some messages for future
	// debugging. It's especially important in production environment. Almost
	// all analytic tools watch stderr so it has to be toggled to `os.Stderr`
	// in production.
	log.SetOutput(os.Stdout)
	token = env.GetEnvVariableByKey("BOT_TOKEN")

	var err error

	session, err = discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Cannot create the session: %v", err)
	}

	session.AddHandler(roles.InitialAdd)
	session.AddHandler(commands.Handler)
	session.AddHandler(messages.Handler)

	// In this example, we only care about receiving message events.
	session.Identify.Intents = discordgo.IntentsGuildMessages

}

func main() {
	// Open a websocket connection to Discord and begin listening.
	err := session.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
		return
	}

	log.Println("Adding commands...")
	registeredCommands := commands.Register(session)

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	log.Println("Removing commands...")
	commands.ClearRegistered(session, registeredCommands)

	fmt.Println("Closing session...")
	// Cleanly close down the Discord session.
	session.Close()

	fmt.Println("Gracefully shutting down...")
}
