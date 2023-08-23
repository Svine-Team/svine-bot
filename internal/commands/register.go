package commands

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

type Commands = []*discordgo.ApplicationCommand

func Register(s *discordgo.Session) Commands {
	registeredCommands := make([]*discordgo.ApplicationCommand, len(List))

	for i, command := range List {
		registeredCommand, err := s.ApplicationCommandCreate(s.State.User.ID, "", command)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", command.Name, err)
		}

		registeredCommands[i] = registeredCommand
	}

	return registeredCommands
}

func ClearRegistered(s *discordgo.Session, commands Commands) {
	for _, command := range commands {
		err := s.ApplicationCommandDelete(s.State.User.ID, "", command.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", command.Name, err)
		}
	}
}
