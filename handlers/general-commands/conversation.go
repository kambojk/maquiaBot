package gencommands

import (
	"crypto/rand"
	tools "maquiaBot/tools"
	"math/big"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Conversation creates a conversation based off of quotes
func Conversation(s *discordgo.Session, m *discordgo.MessageCreate) {
	convoRegex, _ := regexp.Compile(`(?i)(convo?|conversation)?\s+(.+)?(\d+)?`)
	linkRegex, _ := regexp.Compile(`(?i)(https://www|https://|www)\S+`)

	// Get server
	server, err := s.Guild(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "This is not a server!")
		return
	}
	serverData, _ := tools.GetServer(*server, s)
	if len(serverData.Quotes) == 0 {
		s.ChannelMessageSend(m.ChannelID, "No quotes saved for this server! Please see `help quoteadd` to see how to add quotes!")
		return
	}

	// Get number of quotes and if links should be removed and username
	user := &discordgo.User{}
	username := ""
	userQuotes := serverData.Quotes
	num := 2
	excludeLinks := true
	if convoRegex.MatchString(m.Content) {

		username = convoRegex.FindStringSubmatch(m.Content)[2]
		if discordUser, err := s.User(username); err == nil {
			username = discordUser.Username
			tools.ErrRead(s, err)
		} else if len(m.Mentions) > 0 {
			username = m.Mentions[0].Username
		}

		if len(strings.Split(username, " ")) > 1 {
			if num, err = strconv.Atoi(strings.Split(username, " ")[1]); err == nil {
				username = strings.Split(username, " ")[0]
			} else {
				num = 2
			}
		}
		// Get user
		
		members, _ := s.GuildMembers(m.GuildID, "", 1000)
		sort.Slice(members, func(i, j int) bool {
			time1, _ := members[i].JoinedAt.Parse()
			time2, _ := members[j].JoinedAt.Parse()
			return time1.Unix() < time2.Unix()
		})
		for _, member := range members {
			if strings.HasPrefix(strings.ToLower(member.User.Username), strings.ToLower(username)) || strings.HasPrefix(strings.ToLower(member.Nick), strings.ToLower(username)) {
				user, _ = s.User(member.User.ID)
				break
			}
		}

		if user.ID == "" {
			s.ChannelMessageSend(m.ChannelID, "No user with the name **"+username+"** found!")
			return
		}

		userQuotes = []discordgo.Message{}
		for _, quote := range serverData.Quotes {
			if quote.Author.ID == user.ID {
				userQuotes = append(userQuotes, quote)
			}
		}
		if len(userQuotes) == 0 {
			s.ChannelMessageSend(m.ChannelID, "No quotes saved for **"+user.Username+"**! Please see `help quoteadd` to see how to add quotes!")
			return
		}
	}
	// Create the Convo .
	convo := []string{}
	for i := 0; i < num; i++ {
		if len(userQuotes) == 0 {
			break
		}
		roll, _ := rand.Int(rand.Reader, big.NewInt(int64(len(userQuotes))))
		j := roll.Int64()
		if len(userQuotes[j].Attachments) <= 0 || excludeLinks {
			text := userQuotes[j].ContentWithMentionsReplaced()
			if linkRegex.MatchString(text) {
				text = strings.TrimSpace(linkRegex.ReplaceAllString(text, ""))
			}
			if text == "" {
				i--
				continue
			}
			convo = append(convo, "**"+userQuotes[j].Author.Username+"**: "+userQuotes[j].ContentWithMentionsReplaced())
			userQuotes = append(userQuotes[:j],userQuotes[j+1:]...)
		} else if len(userQuotes[j].Attachments) > 0 {
			convo = append(convo, "**"+userQuotes[j].Author.Username+"**: "+userQuotes[j].ContentWithMentionsReplaced()+" "+userQuotes[j].Attachments[0].URL)
			userQuotes = append(userQuotes[:j],userQuotes[j+1:]...)
		} else {
			i--
			continue
		}

		if len(strings.Join(convo, "\n")) > 2000 {
			convo = convo[:len(convo)-1]
			break
		}
	}
	if excludeLinks {
		s.ChannelMessageSend(m.ChannelID, "```md\n"+strings.Join(convo, "\n")+"```")
	} else {
		s.ChannelMessageSend(m.ChannelID, strings.Join(convo, "\n"))
	}
}
