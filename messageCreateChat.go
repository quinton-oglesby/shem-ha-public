package main

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/sashabaranov/go-openai"
	"log"
	"math/rand"
	"time"
)

func messageCreateChat(session *discordgo.Session, message *discordgo.MessageCreate) {
	//// Check if the chatting is enabled in the server.
	//var chatEnabled int
	//query := fmt.Sprintf(`SELECT chat_enabled FROM servers WHERE server_id = %v`, message.GuildID)
	//err := db.QueryRow(query).Scan(&chatEnabled)
	//if err != nil {
	//	log.Printf("%vERROR%v - COULD NOT QUERY SERVERS:\n\t%v", Red, Reset, err)
	//	return
	//}
	//
	//if chatEnabled == 0 {
	//	return
	//}
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
	messageContent := re.ReplaceAllString(message.Content, "")

	// Ultimately ignore the messages with no content in them.
	if messageContent == "" {
		return
	}

	// Craeting and seeding the random number generator.
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generating the chance to repond to the message.
	chance := random.Float64()

	// Logging the chance to repond to the message.
	log.Printf("%vCHANCE:%v %v", Cyan, Reset, chance*100.0)
	if chance*100 < getGuildChatChance(message.GuildID)*1.0 {
		// Creating the OpenAI client.
		msgArr, _ := session.ChannelMessages(message.ChannelID, 4, message.ID, "", "")
		client := openai.NewClient(tokens.OpenAIToken)

		res, _ := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model:     openai.GPT3Dot5Turbo,
				MaxTokens: getGuildChatLength(message.GuildID),
				//Stop:      []string{"\n"},
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleSystem,
						Content: "You watch over multiple people and provide input. You are an arrogant goddess. You are rude and abrasive to your followers. Your name is Shem-Ha. The people you watch over affectionately refer to you as Shemmy. You do not spaek for anyone but yourself.",
					},
					{
						Role: openai.ChatMessageRoleUser,
						Content: fmt.Sprintf("%v: %v\n%v: %v\n%v: %v\n%v: %v\n%v: %v",
							message.Author.Username, messageContent,
							msgArr[0].Author.Username, msgArr[0].Content,
							msgArr[1].Author.Username, msgArr[1].Content,
							msgArr[2].Author.Username, msgArr[2].Content,
							msgArr[3].Author.Username, msgArr[3].Content),
					},
				},
			},
		)

		// https://pkg.go.dev/github.com/bwmarrin/discordgo#Session.ChannelMessageSendComplex
		_, err := session.ChannelMessageSendComplex(message.ChannelID, &discordgo.MessageSend{
			Content: res.Choices[0].Message.Content,
			//Reference: message.Reference(),
			//AllowedMentions: &discordgo.MessageAllowedMentions{
			//	Parse: nil,
			//},
		})
		if err != nil {
			log.Printf("COULD NOT REPLY TO %v: %v", message, err)
		}
	}
}
