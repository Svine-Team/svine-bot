package commands

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func pivoHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	msgArgs := make([]interface{}, 0, len(options))
	msgFormat := "Successfully added role to user!\n"

	var (
		user *discordgo.User
		role *discordgo.Role
	)

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

	err := s.GuildMemberRoleAdd(i.GuildID, user.ID, role.ID)
	if err != nil {
		log.Panicf("Couldn't add the role '%v' to user '%v'", role.ID, user.ID)
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf(
				msgFormat,
				msgArgs...,
			),
		},
	})

	if err != nil {
		log.Panicf("Couldn't interact with user '%v'", user.ID)
	}
}

func lohHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	msgArgs := make([]interface{}, 0, len(options))
	msgFormat := "Successfully added role to user!\n"

	var (
		user *discordgo.User
	)

	if opt, ok := optionMap["user-option"]; !ok {
		log.Panicf("Couldn't get user from options")
	} else {
		user = opt.UserValue(nil)
		msgArgs = append(msgArgs, user.ID)
		msgFormat += "> user-option: <@%s>\n"
	}

	err := s.GuildMemberRoleAdd(i.GuildID, user.ID, "1144009767895957524") // хихихаха
	if err != nil {
		log.Panicf("Couldn't add the role 'loh' to user '%v'", user.ID)
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf(
				msgFormat,
				msgArgs...,
			),
		},
	})

	if err != nil {
		log.Panicf("Couldn't interact with user '%v'", user.ID)
	}
}

func Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	command := i.ApplicationCommandData().Name

	switch command {
	case string(EPivo):
		pivoHandler(s, i)
	case string(ELoh):
		lohHandler(s, i)
	}
}
