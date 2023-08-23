package operations

// Package operations provides a centralized gateway for processing requests
// on Discord API operations. It is mostly some small CRUD wrappers.

import (
	"fmt"
	"log"

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

// Checks if the role exists already before creating.
func CreateRolesForGuild(
	session *discordgo.Session,
	guild *discordgo.Guild,
    roleNamesToCreate []string,
) ([]*discordgo.Role, error) {
    createdRoles := []*discordgo.Role{}

    for _, roleNameToCreate := range roleNamesToCreate {
        log.Printf("Creating ranks (roles) for guild '%v'", guild.ID)
        existingRoles, err := session.GuildRoles(guild.ID)
        if err != nil {
            log.Panicf("Couldn't fetch guild roles for guild '%v'", guild.ID)
        }
        existingRoleNames := make(map[string]bool)
        for _, role := range existingRoles {
            existingRoleNames[role.Name] = true
        }

        if _, roleAlreadyExists := existingRoleNames[roleNameToCreate]; roleAlreadyExists {
            continue
        }

        createdRole, err := CreateRole(session, guild, roleNameToCreate, 100)
        if err != nil {
            return createdRoles, fmt.Errorf("During creation of role '%v' on guild '%v' an error has occured %v", createdRole.ID, guild.ID, err)
        }

        createdRoles = append(createdRoles, createdRole)
    }

    return createdRoles, nil
}
