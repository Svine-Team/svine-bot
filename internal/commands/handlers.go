package commands

import (
	"fmt"
	"log"

	"github.com/Svine-Team/svine-bot/internal/interactions"
	"github.com/bwmarrin/discordgo"
)

func pivoHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := interactions.InteractionMapper(s, i)

	user, role := options.User, options.Role

	err := s.GuildMemberRoleAdd(i.GuildID, user.ID, role.ID)
	if err != nil {
		log.Panicf("Couldn't add the role '%v' to user '%v'", role.Name, user.Username)
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Successfully added role '%v' to user '%v'", role.Name, user.Username),
		},
	})

	if err != nil {
		log.Panicf("Couldn't interact with user '%v'", user.Username)
	}
}

func Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	command := i.ApplicationCommandData().Name

	switch command {
	case string(EPivo):
		pivoHandler(s, i)
	}
}
