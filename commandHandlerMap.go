package main

import (
	"database/sql"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
)

// Creating a map of event handlers to respond to application commands.
// https://pkg.go.dev/github.com/bwmarrin/discordgo#EventHandler
var commandHandlers = map[string]func(session *discordgo.Session, interaction *discordgo.InteractionCreate){
	"echo": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		if !guildIsSetup(interaction.GuildID) {
			//https://pkg.go.dev/github.com/bwmarrin/discordgo#Session.InteractionRespond
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "You need to run /setup first!",
				},
			})
			return
		}

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
	"setup": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		var serverID int

		// Performing a single row query to check if the guild is already setup.
		query := fmt.Sprintf(`SELECT server_id FROM servers WHERE server_id = %v;`, interaction.GuildID)
		err := db.QueryRow(query).Scan(&serverID)
		if err != nil && err != sql.ErrNoRows {
			log.Printf("%vERROR%v - COULD NOT QUERY SERVERS: %v", Red, Reset, err)
		}

		if err == sql.ErrNoRows {
			// Creating and executing a SQL query to set up.
			query = fmt.Sprintf(
				`INSERT INTO servers(server_id, exempt, paid, chat_enabled, frame_enabled, markov_enabled, chat_chance, chat_length, chat_limit, frame_frequency, markov_frequency, markov_owo)
					  VALUES(%v, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)`,
				interaction.GuildID)

			_, err = db.Exec(query)
			if err != nil {
				// Reporting error to user.
				session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf("COULD NOT RUN SETUP: %v", err),
					},
				})

				log.Println("COULD NOT RUN SETUP: %v", err)
				return
			}

		} else {
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "You're already set up!",
				},
			})
		}

		// Responding to the interaction.
		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "All set up!",
			},
		})
	},
	"chat_list_channels": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		if !guildIsSetup(interaction.GuildID) {
			//https://pkg.go.dev/github.com/bwmarrin/discordgo#Session.InteractionRespond
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "You need to run /setup first!",
				},
			})
			return
		}

		// Creating and executing the query to list the channels.
		var count int64
		query := fmt.Sprintf(`SELECT COUNT(*) FROM channels WHERE server_id = %v;`, interaction.GuildID)
		_ = db.QueryRow(query).Scan(&count)
		switch {
		case count == 0:
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "I am not allowed to respond in any channels currently!",
				},
			})
		default:
			query = fmt.Sprintf(`SELECT * FROM channels WHERE server_id = %v;`, interaction.GuildID)

			rows, err := db.Query(query)
			if err != nil && err != sql.ErrNoRows {
				log.Printf("%vERROR%v - COULD NOT QUERY CHANNELS:\n\t%v", Red, Reset, err)
				return
			}

			// Getting the list of channel IDs that she is allowed to respond in.
			var channelIDs []string
			var channelID string
			var serverID string
			for rows.Next() {
				err := rows.Scan(&serverID, &channelID)
				if err != nil {
					log.Printf("%vERROR%v - COULD NOT RETRIEVE CHANNEL FROM ROW:\n\t%v", Red, Reset, err)
					continue
				}
				log.Println(channelID)

				channelIDs = append(channelIDs, channelID)
			}

			// Constructing the string of channels that she is allowed to respond in.
			var response string
			log.Println(channelIDs)
			for _, chid := range channelIDs {
				response += fmt.Sprintln("\t<#" + chid + ">")
			}

			response = fmt.Sprintln("I am allowed to respond in the following channel(s):\n", response)
			log.Println(response)

			// Finally, responding to the interaction.
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: response,
				},
			})
		}
	},
	"chat_add_channel": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		if !guildIsSetup(interaction.GuildID) {
			//https://pkg.go.dev/github.com/bwmarrin/discordgo#Session.InteractionRespond
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "You need to run /setup first!",
				},
			})
			return
		}

		var channelID string

		// Snagging the server ID.
		serverID := interaction.GuildID

		// Performing a single row query to check if the user already has the card in their collection.
		query := fmt.Sprintf(`SELECT channel_id FROM channels WHERE server_id = %v AND channel_id = %v;`,
			interaction.GuildID, interaction.ChannelID)
		err := db.QueryRow(query).Scan(&channelID)
		if err != nil && err != sql.ErrNoRows {
			log.Printf("%vERROR%v - COULD NOT QUERY CHANNELS: %v", Red, Reset, err)
		}

		if err == sql.ErrNoRows {
			// Snagging the target channel ID.
			channelID = interaction.ApplicationCommandData().Options[0].ChannelValue(session).ID

			// Creating and executing a SQL query to add the channel to the list of approved channels.
			query = fmt.Sprintf(`INSERT INTO channels(server_id, channel_id) VALUES (%v, %v)`, serverID, channelID)

			_, err = db.Exec(query)
			if err != nil {
				// Reporting error to user.
				session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf("COULD NOT ADD CHANNEL: %v", err),
					},
				})

				log.Println("COULD NOT ADD CHANNEL: ", err)
				return
			}

			// Responding to the interaction.
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Successfully added <#%v> to the list of channels I am allowed to respond in.", channelID),
				},
			})
		} else {
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Channel is already approved!",
				},
			})
		}
	},
	"chat_remove_channel": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		if !guildIsSetup(interaction.GuildID) {
			//https://pkg.go.dev/github.com/bwmarrin/discordgo#Session.InteractionRespond
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "You need to run /setup first!",
				},
			})
			return
		}

		var serverID string

		// Snagging the target channel ID.
		channelID := interaction.ApplicationCommandData().Options[0].ChannelValue(session).ID

		// Performing a single row query to check if the channel is already disallowed..
		query := fmt.Sprintf(`SELECT server_id FROM channels WHERE channel_id = %v;`,
			channelID)
		err := db.QueryRow(query).Scan(&serverID)
		if err != nil && err != sql.ErrNoRows {
			log.Printf("%vERROR%v - COULD NOT QUERY CHANNELS: %v", Red, Reset, err)
		}

		if err == sql.ErrNoRows {
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Channel is already disallowed!",
				},
			})
		} else {
			// Snagging the server ID.
			serverID := interaction.GuildID

			// Creating and executing a SQL query to remove the channel to the list of approved channels.
			query = fmt.Sprintf(`DELETE FROM channels WHERE server_id = %v AND channel_id = %v`, serverID, channelID)

			_, err = db.Exec(query)
			if err != nil {
				// Reporting error to user.
				session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf("COULD NOT ADD CHANNEL: %v", err),
					},
				})

				log.Println("COULD NOT REMOVE CHANNEL: ", err)
				return
			}

			// Responding to the interaction.
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Successfully removed <#%v> from the list of channels I am allowed to respond in.", channelID),
				},
			})
		}
	},
	"chat_get_chance": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		// Responding to the interaction.
		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("The current response chance is %v percent.",
					getGuildChatChance(interaction.GuildID)),
			},
		})
	},
	"chat_set_chance": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		chance := interaction.ApplicationCommandData().Options[0].FloatValue()

		// Creating a single row query to set the chance of the chat response chance.
		query := fmt.Sprintf(`UPDATE servers SET chat_chance = %v WHERE server_id = %v`, chance, interaction.GuildID)
		_, err := db.Query(query)
		if err != nil {
			log.Printf("%vERROR%v - COULD NOT QUERY SERVERS: %v", Red, Reset, err)

			// Responding to the interaction.
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Ran into an error, please try again later!"),
				},
			})
		}

		// Responding to the interaction.
		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Successfully updated the response chance. The reponse chance is now %v percent.", chance),
			},
		})
	},
	"chat_get_length": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		// Responding to the interaction.
		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("The current response chance is %v tokens.",
					getGuildChatLength(interaction.GuildID)),
			},
		})
	},
	"chat_set_length": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		tokens := interaction.ApplicationCommandData().Options[0].IntValue()

		// Creating a single row query to set the chance of the chat response chance.
		query := fmt.Sprintf(`UPDATE servers SET chat_length = %v WHERE server_id = %v`, tokens, interaction.GuildID)
		_, err := db.Query(query)
		if err != nil {
			log.Printf("%vERROR%v - COULD NOT QUERY SERVERS: %v", Red, Reset, err)

			// Responding to the interaction.
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Ran into an error, please try again later!"),
				},
			})
		}

		// Responding to the interaction.
		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Successfully updated the response length. The reponse length is now %v tokens.", tokens),
			},
		})
	},

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
