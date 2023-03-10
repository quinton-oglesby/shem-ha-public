package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/varsapphire/OwO"
	"log"
	"math/rand"
	"strings"
	"time"
)

func messageCreateMarkov(session *discordgo.Session, message *discordgo.MessageCreate) {
	// Check if bot is allowed to respond in the channel.
	if !inArray(message.ChannelID, getGuildMarkovChannels(message.GuildID)) {
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

	// Craeting and seeding the random number generator.
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generating the chance to repond to the message.
	chance := random.Float64()

	// Logging the chance to repond to the message.
	log.Printf("%vCHANCE:%v %v", Yellow, Reset, chance*100.0)
	if chance*100 < getGuildMarkovChance(message.GuildID)*1.0 {
		// Grab all the user's messages from the database.
		query = fmt.Sprintf(`SELECT content FROM messages WHERE user_id = %v;`, message.Author.ID)
		rows, err := db.Query(query)
		if err != nil {
			log.Println(err)
		}

		// Add all the snagged messages to one giant string.
		content := ""
		for rows.Next() {
			var message string
			rows.Scan(&message)
			content = content + message + " "
		}

		// Create the Markov chain with an order of 2 to mimic the user.
		chain := NewChain(2)

		// Feed in the giant string of messages for training.
		chain.Build(strings.NewReader(content))

		// Generate the chain.
		content_chain := chain.Generate(32)
		if err != nil {
			log.Println(err)
		}

		// Creating a webhook to chat in the channel, will be deleting it afterwards.
		webhook, err := session.WebhookCreate(message.ChannelID, "mimic", "")
		if err != nil {
			fmt.Println("Webhook Creation Error ", err)
		}

		if getOwO(message.GuildID) {
			content_chain = OwO.WhatsThis(content_chain)
		}

		// Setting the parameters for the webhook that will mimic the user.
		params := discordgo.WebhookParams{}
		params.Content = content_chain
		params.Username = message.Author.Username
		params.AvatarURL = message.Author.AvatarURL(message.Author.Avatar)

		// Executing the webhook that will mimic the user.
		_, err = session.WebhookExecute(webhook.ID, webhook.Token, true, &params)
		if err != nil {
			log.Println(err)
			return
		}

		err = session.WebhookDelete(webhook.ID)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
