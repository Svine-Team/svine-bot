package commands

import "github.com/bwmarrin/discordgo"

var List = []*discordgo.ApplicationCommand{
	{
		Name:        string(EPivo),
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
