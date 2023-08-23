package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Svine-Team/svine-bot/pkg/operations"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
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

// All commands and options must have a description
// Commands/options without description will fail the registration
// of the command.
var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name: "basic-command",
			Description: "Basic command",
		},
		{
			Name: "cool-basic-command",
			Description: "Basic command",
		},
        {
            Name: "pivo",
            Description: "Ranking by the times we met",
            Options: []*discordgo.ApplicationCommandOption{
                {
                    Type:        discordgo.ApplicationCommandOptionUser,
                    Name:        "user-option",
                    Description: "User option",
                    Required:    true,
                },
                {
                    Type:        discordgo.ApplicationCommandOptionRole,
                    Name:        "role-option",
                    Description: "Role option",
                    Required:    true,
                },
            },
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

        "pivo": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
            options := i.ApplicationCommandData().Options
            optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
            for _, opt := range options {
                optionMap[opt.Name] = opt
            }

            msgArgs := make([]interface{}, 0, len(options))
            msgFormat := "Successfully added role to user!\n"
            var user *discordgo.User
            var role *discordgo.Role

            if opt, ok := optionMap["user-option"]; !ok {
                log.Panicf("Couldn't get user from options")
            } else {
                user = opt.UserValue(nil)
                msgArgs = append(msgArgs, user.ID)
                msgFormat += "> user-option: <@%s>\n"
            }

            if opt, ok := optionMap["role-option"]; !ok {
                log.Panicf("Couldn't get role from options")
            } else {
                role = opt.RoleValue(nil, "")
                msgArgs = append(msgArgs, role.ID)
                msgFormat += "> role-option: <@%s>\n"
            }

            err := session.GuildMemberRoleAdd(i.GuildID, user.ID, role.ID)
            if err != nil {
                log.Panicf("Couldn't add the role '%v' to user '%v'", role.ID, user.ID)
            }

            s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
                Type: discordgo.InteractionResponseChannelMessageWithSource,
                Data: &discordgo.InteractionResponseData{
                    Content: fmt.Sprintf(
                        msgFormat,
                        msgArgs...,
                    ),
                },
            })
        },
    }
)

var session *discordgo.Session

func createRoles(s *discordgo.Session, ready *discordgo.Ready) {
    user := s.State.User
    log.Printf("Logged in as %v#%v", user.Username, user.Discriminator)

    // TODO: Make an adapter to fetch list from different sources (json,
    //   server...)
    roleNamesToCreate := []string{"created-test-role", "new-role"}

    for _, guild := range ready.Guilds {
        createdRoles, err := operations.CreateRolesForGuild(session, guild, roleNamesToCreate)

        if err != nil {
            log.Panicf("Couldn't create role: %v", err)
        }

        if len(createdRoles) == 0 {
            log.Printf("All roles already exist for guild '%v'", guild.ID)
            continue
        }

        log.Printf("Created new roles for guild '%v': %v", guild.ID, createdRoles)
    }

}

func initHandlers() {
    session.AddHandler(createRoles)

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
