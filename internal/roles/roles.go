package roles

import (
	"log"

	"github.com/Svine-Team/svine-bot/pkg/operations"
	"github.com/bwmarrin/discordgo"
)

func InitialAdd(s *discordgo.Session, ready *discordgo.Ready) {
	user := s.State.User
	log.Printf("Logged in as %v#%v", user.Username, user.Discriminator)

	// TODO: Make an adapter to fetch list from different sources (json,
	//   server...)
	roleNamesToCreate := []string{"created-test-role", "new-role"}

	for _, guild := range ready.Guilds {
		createdRoles, err := operations.CreateRolesForGuild(s, guild, roleNamesToCreate)

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
