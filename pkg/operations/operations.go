// Package operations provides a centralized gateway for processing requests
// on Discord API operations. It is mostly some small CRUD wrappers.
package operations

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

const (
	roleHoist             = true
	roleMentionable           = true
)

// Couldn't pass simple constants as params.
func pointerTo[T any](v T) *T {
	return &v
}

func CreateRole(
	session *discordgo.Session,
	guild *discordgo.Guild,
	roleName string,
	roleColor int,
) (*discordgo.Role, error) {
	role, err := session.GuildRoleCreate(guild.ID, &discordgo.RoleParams{
		Name:        roleName,
		Color:       &roleColor,
		Hoist:       pointerTo(roleHoist),
		Mentionable: pointerTo(roleMentionable),
	})
	if err != nil {
		return nil, fmt.Errorf("Unable to create role: %w", err)
	}

	err = session.State.RoleAdd(guild.ID, role)
	if err != nil {
		return nil, fmt.Errorf("Unable to add role to state cache: %w", err)
	}

	return role, nil
}
