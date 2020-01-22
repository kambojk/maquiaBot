package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Late explains the late functionality
func Late(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: late / old / ancient"
	embed.Description = "`(late|old|ancient)` lets you send a late video."
	return embed
}
