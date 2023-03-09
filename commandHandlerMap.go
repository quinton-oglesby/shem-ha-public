package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
)

// Creating a map of event handlers to respond to application commands.
// https://pkg.go.dev/github.com/bwmarrin/discordgo#EventHandler
var commandHandlers = map[string]func(session *discordgo.Session, interaction *discordgo.InteractionCreate){
	"echo": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		// Grabbing the channel ID and the content of the message to echo.
		channel := interaction.ApplicationCommandData().Options[0].ChannelValue(session)
		content := interaction.ApplicationCommandData().Options[1].StringValue()
		msg, err := session.ChannelMessageSend(channel.ID, content)
		if err != nil {
			log.Printf("COULD NOT SEND MESSAGE '%v': %v", msg, err)

			//https://pkg.go.dev/github.com/bwmarrin/discordgo#Session.InteractionRespond
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("COULD NOT SEND MESSAGE '%v': %v", msg, err),
				},
			})
			return

		}

		// Responding to the interaction.
		//https://pkg.go.dev/github.com/bwmarrin/discordgo#Session.InteractionRespond
		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Successfully sent '%v' to channel '%v'", content, channel.Name),
			},
		})
	},
	//"get_chance": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	//	// Responding to the interaction.
	//	//https://pkg.go.dev/github.com/bwmarrin/discordgo#Session.InteractionRespond
	//	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
	//		Type: discordgo.InteractionResponseChannelMessageWithSource,
	//		Data: &discordgo.InteractionResponseData{
	//			Content: fmt.Sprintf("The current response chance is %v percent.", gpt3Parameters.Chance),
	//		},
	//	})
	//},
	//"get_length": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	//	// Responding to the interaction.
	//	//https://pkg.go.dev/github.com/bwmarrin/discordgo#Session.InteractionRespond
	//	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
	//		Type: discordgo.InteractionResponseChannelMessageWithSource,
	//		Data: &discordgo.InteractionResponseData{
	//			Content: fmt.Sprintf("The current response length is %v tokens.", gpt3Parameters.Length),
	//		},
	//	})
	//},
	//"set_chance": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	//	gpt3Parameters.Chance = interaction.ApplicationCommandData().Options[0].FloatValue()
	//
	//	// Marshall the new parameters to save.
	//	jsonBytes, err := json.Marshal(gpt3Parameters)
	//	if err != nil {
	//		log.Println("ERROR MARSHALING JSON: ", err)
	//
	//		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
	//			Type: discordgo.InteractionResponseChannelMessageWithSource,
	//			Data: &discordgo.InteractionResponseData{
	//				Content: fmt.Sprintf("FAILED TO UPDATE CHANCE: %v", err),
	//			},
	//		})
	//
	//		return
	//	}
	//
	//	// Save updated parameters into parameters.json.
	//	err = os.WriteFile("parameters.json", jsonBytes, 0644)
	//	if err != nil {
	//		log.Panicln("ERROR SAVING JSON: ", err)
	//
	//		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
	//			Type: discordgo.InteractionResponseChannelMessageWithSource,
	//			Data: &discordgo.InteractionResponseData{
	//				Content: fmt.Sprintf("FAILED TO UPDATE CHANCE: %v", err),
	//			},
	//		})
	//
	//		return
	//	}
	//
	//	// Responding to the interaction.
	//	//https://pkg.go.dev/github.com/bwmarrin/discordgo#Session.InteractionRespond
	//	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
	//		Type: discordgo.InteractionResponseChannelMessageWithSource,
	//		Data: &discordgo.InteractionResponseData{
	//			Content: fmt.Sprintf("Successfully updated the response chance. The reponse chance is now %v percent.", gpt3Parameters.Chance),
	//		},
	//	})
	//},
	//"set_length": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	//	gpt3Parameters.Length = interaction.ApplicationCommandData().Options[0].IntValue()
	//
	//	// Marshall the new parameters to save.
	//	jsonBytes, err := json.Marshal(gpt3Parameters)
	//	if err != nil {
	//		log.Println("ERROR MARSHALING JSON: ", err)
	//
	//		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
	//			Type: discordgo.InteractionResponseChannelMessageWithSource,
	//			Data: &discordgo.InteractionResponseData{
	//				Content: fmt.Sprintf("FAILED TO UPDATE LENGTH: %v", err),
	//			},
	//		})
	//
	//		return
	//	}
	//
	//	// Save updated parameters into parameters.json.
	//	err = os.WriteFile("parameters.json", jsonBytes, 0644)
	//	if err != nil {
	//		log.Println("ERROR SAVING JSON: ", err)
	//
	//		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
	//			Type: discordgo.InteractionResponseChannelMessageWithSource,
	//			Data: &discordgo.InteractionResponseData{
	//				Content: fmt.Sprintf("FAILED TO UPDATE LENGTH: %v", err),
	//			},
	//		})
	//
	//		return
	//	}
	//
	//	// Responding to the interaction.
	//	//https://pkg.go.dev/github.com/bwmarrin/discordgo#Session.InteractionRespond
	//	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
	//		Type: discordgo.InteractionResponseChannelMessageWithSource,
	//		Data: &discordgo.InteractionResponseData{
	//			Content: fmt.Sprintf("Successfully updated the response length. The reponse length is now %v tokens.", gpt3Parameters.Length),
	//		},
	//	})
	//},
	//"pop_channel": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	//	// Snagging the target channel ID.
	//	targetChannelName := interaction.ApplicationCommandData().Options[0].ChannelValue(session).Name
	//	targetChannelID := interaction.ApplicationCommandData().Options[0].ChannelValue(session).ID
	//
	//	// Checking if channel is already in the list of approved channels.
	//	if !stringInArray(targetChannelID, channels.Channels) {
	//		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
	//			Type: discordgo.InteractionResponseChannelMessageWithSource,
	//			Data: &discordgo.InteractionResponseData{
	//				Content: fmt.Sprintf("Channel '%v' is not in the list of approved channels.", targetChannelName),
	//			},
	//		})
	//
	//		return
	//	}
	//
	//	// Remove channel from the list of channels allowed.
	//	channels.Channels = removeStringFromArray(targetChannelID, channels.Channels)
	//
	//	// Marshall the new channels to save.
	//	jsonBytes, err := json.Marshal(channels)
	//	if err != nil {
	//		log.Println("ERROR MARSHALING JSON: ", err)
	//
	//		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
	//			Type: discordgo.InteractionResponseChannelMessageWithSource,
	//			Data: &discordgo.InteractionResponseData{
	//				Content: fmt.Sprintf("FAILED TO UPDATE CHANNELS: %v", err),
	//			},
	//		})
	//
	//		return
	//	}
	//
	//	// Save updated parameters into parameters.json.
	//	err = os.WriteFile("channels.json", jsonBytes, 0644)
	//	if err != nil {
	//		log.Println("ERROR SAVING JSON: ", err)
	//
	//		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
	//			Type: discordgo.InteractionResponseChannelMessageWithSource,
	//			Data: &discordgo.InteractionResponseData{
	//				Content: fmt.Sprintf("FAILED TO UPDATE CHANNELS: %v", err),
	//			},
	//		})
	//
	//		return
	//	}
	//
	//	// Responding to the interaction.
	//	//https://pkg.go.dev/github.com/bwmarrin/discordgo#Session.InteractionRespond
	//	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
	//		Type: discordgo.InteractionResponseChannelMessageWithSource,
	//		Data: &discordgo.InteractionResponseData{
	//			Content: fmt.Sprintf("Successfully removed '%v' from the list of channels I am allowed to respond in.", targetChannelName),
	//		},
	//	})
	//},
	//"append_channel": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	//
	//	// Snagging the target channel ID.
	//	targetChannelName := interaction.ApplicationCommandData().Options[0].ChannelValue(session).Name
	//	targetChannelID := interaction.ApplicationCommandData().Options[0].ChannelValue(session).ID
	//
	//	// Checking if channel is already in the list of approved channels.
	//	if stringInArray(targetChannelID, channels.Channels) {
	//		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
	//			Type: discordgo.InteractionResponseChannelMessageWithSource,
	//			Data: &discordgo.InteractionResponseData{
	//				Content: fmt.Sprintf("Channel '%v' is already in the list of approved channels.", targetChannelName),
	//			},
	//		})
	//
	//		return
	//	}
	//
	//	// Add channel to the list of channels allowed.
	//	channels.Channels = append(channels.Channels, targetChannelID)
	//
	//	// Marshall the new channels to save.
	//	jsonBytes, err := json.Marshal(channels)
	//	if err != nil {
	//		log.Println("ERROR MARSHALING JSON: ", err)
	//
	//		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
	//			Type: discordgo.InteractionResponseChannelMessageWithSource,
	//			Data: &discordgo.InteractionResponseData{
	//				Content: fmt.Sprintf("FAILED TO UPDATE CHANNELS: %v", err),
	//			},
	//		})
	//
	//		return
	//	}
	//
	//	// Save updated parameters into parameters.json.
	//	err = os.WriteFile("channels.json", jsonBytes, 0644)
	//	if err != nil {
	//		log.Println("ERROR SAVING JSON: ", err)
	//
	//		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
	//			Type: discordgo.InteractionResponseChannelMessageWithSource,
	//			Data: &discordgo.InteractionResponseData{
	//				Content: fmt.Sprintf("FAILED TO UPDATE CHANNELS: %v", err),
	//			},
	//		})
	//
	//		return
	//	}
	//
	//	// Responding to the interaction.
	//	//https://pkg.go.dev/github.com/bwmarrin/discordgo#Session.InteractionRespond
	//	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
	//		Type: discordgo.InteractionResponseChannelMessageWithSource,
	//		Data: &discordgo.InteractionResponseData{
	//			Content: fmt.Sprintf("Successfully added '%v' to the list of channels I am allowed to respond in.", targetChannelName),
	//		},
	//	})
	//},
	//"list_channels": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	//
	//	if len(channels.Channels) > 0 {
	//		chnls := ""
	//		for _, channelID := range channels.Channels {
	//			channel, err := session.Channel(channelID)
	//			if err != nil {
	//				log.Println("ERROR RETREIVING CHANNELS: ", err)
	//
	//				session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
	//					Type: discordgo.InteractionResponseChannelMessageWithSource,
	//					Data: &discordgo.InteractionResponseData{
	//						Content: fmt.Sprintf("FAILED TO GET CHANNELS: %v", err),
	//					},
	//				})
	//				return
	//			}
	//
	//			chnls += channel.Name + "\n"
	//		}
	//
	//		// Responding to the interaction.
	//		//https://pkg.go.dev/github.com/bwmarrin/discordgo#Session.InteractionRespond
	//		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
	//			Type: discordgo.InteractionResponseChannelMessageWithSource,
	//			Data: &discordgo.InteractionResponseData{
	//				Content: chnls,
	//			},
	//		})
	//	} else {
	//		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
	//			Type: discordgo.InteractionResponseChannelMessageWithSource,
	//			Data: &discordgo.InteractionResponseData{
	//				Content: "CHANNEL LIST IS EMPTY",
	//			},
	//		})
	//	}
	//
	//},
}
