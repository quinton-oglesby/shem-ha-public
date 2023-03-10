package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
)

func messageCreateMarkov(session *discordgo.Session, message *discordgo.MessageCreate) {
	// Check if bot is allowed to respond in the channel.
	if !inArray(message.ChannelID, chatGetChannels(message.GuildID)) {
		return
	}

	// Ignore all messages that were created by the bot itself.
	if message.Author.ID == session.State.User.ID {
		return
	}

	// Ignore all messages with the discriminator #0000 (Webhooks).
	if message.WebhookID != "" {
		return
	}

	// Filter out all URLs in the message.
	message_content := re.ReplaceAllString(message.Content, "")

	// Ultimately ignore all messages with no content in them.
	if message_content == "" {
		return
	}

	// Store the message data into the table.
	query := fmt.Sprintf(`INSERT INTO messages(server_id, channel_id, user_id, content) VALUES(%v, %v, %v, "%v");`,
		message.GuildID, message.ChannelID, message.Author.ID, message_content)
	_, err := db.Exec(query)
	if err != nil {
		log.Println(err)
		return
	}
}
