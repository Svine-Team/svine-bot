package interactions

import (
	"github.com/bwmarrin/discordgo"
)

type InteractionMap struct {
	User discordgo.User
	Role discordgo.Role
}

func InteractionMapper(s *discordgo.Session, i *discordgo.InteractionCreate) InteractionMap {
	options := i.ApplicationCommandData().Options

	var optionsMap InteractionMap

	for _, option := range options {
		switch option.Type {
		case discordgo.ApplicationCommandOptionUser:
			optionsMap.User = *option.UserValue(s)

		case discordgo.ApplicationCommandOptionRole:
			optionsMap.Role = *option.RoleValue(s, i.GuildID)
		}
	}

	return optionsMap
}
