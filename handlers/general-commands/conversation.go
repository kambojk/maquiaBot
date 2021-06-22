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

	// Get number of quotes and if links should be removed
	num := 2
	excludeLinks := true
	if convoRegex.MatchString(m.Content) {
		num, err = strconv.Atoi(convoRegex.FindStringSubmatch(m.Content)[2])
		if err != nil || num < 2 {
			num = 2
		}
	}

	// Get the username
	user := &discordgo.User{}
	username := ""
	userQuotes := serverData.Quotes
	//number := 0
	if convoRegex.MatchString(m.Content) {
		username = convoRegex.FindStringSubmatch(m.Content)[2]
		if discordUser, err := s.User(username); err == nil {
			username = discordUser.Username
			tools.ErrRead(s, err)
		} else if len(m.Mentions) > 0 {
			username = m.Mentions[0].Username
		}
		if len(strings.Split(username, " ")) > 1 {
			if _, err = strconv.Atoi(strings.Split(username, " ")[1]); err == nil {
				username = strings.Split(username, " ")[0]
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
			s.ChannelMessageSend(m.ChannelID, "the hoes: " + member.User.Username)
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
		} else if len(userQuotes[j].Attachments) > 0 {
			convo = append(convo, "**"+userQuotes[j].Author.Username+"**: "+userQuotes[j].ContentWithMentionsReplaced()+" "+userQuotes[j].Attachments[0].URL)
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
