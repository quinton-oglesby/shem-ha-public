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
					  VALUES(%v, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0)`,
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
				Content: "You're all set up!\n\nPlease consider setting up a monthly donation of $5 or more as Shem-Ha's chat function does cost money out of my own pocket to run. \n\nhttps://ko-fi.com/variableformation",
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
					Content: "I am not allowed to save from or respond in any channels currently!",
				},
			})
		default:
			channelIDs := getGuildChatChannels(interaction.GuildID)

			// Constructing the string of channels that she is allowed to respond in.
			var response string
			log.Println(channelIDs)
			for _, chid := range channelIDs {
				response += fmt.Sprintln("\t<#" + chid + ">")
			}

			response = fmt.Sprintln("I am allowed to respond in the following channel(s):\n", response)

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
		query := fmt.Sprintf(`SELECT channel_id FROM chat_channels WHERE server_id = %v AND channel_id = %v;`,
			interaction.GuildID, interaction.ChannelID)
		err := db.QueryRow(query).Scan(&channelID)
		if err != nil && err != sql.ErrNoRows {
			log.Printf("%vERROR%v - COULD NOT QUERY CHANNELS: %v", Red, Reset, err)
		}

		if err == sql.ErrNoRows {
			// Snagging the target channel ID.
			channelID = interaction.ApplicationCommandData().Options[0].ChannelValue(session).ID

			// Creating and executing a SQL query to add the channel to the list of approved channels.
			query = fmt.Sprintf(`INSERT INTO chat_channels(server_id, channel_id) VALUES (%v, %v)`, serverID, channelID)

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
					Content: fmt.Sprintf("Successfully added <#%v> to the list of channels I am allowed to gather data from and respond in.", channelID),
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
		query := fmt.Sprintf(`SELECT server_id FROM chat_channels WHERE channel_id = %v;`,
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
			query = fmt.Sprintf(`DELETE FROM chat_channels WHERE server_id = %v AND channel_id = %v`, serverID, channelID)

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
					Content: fmt.Sprintf("Successfully removed <#%v> from the list of channels I am allowed to save data from and respond in.", channelID),
				},
			})
		}
	},
	"chat_get_chance": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
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
	"markov_list_channels": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
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
		query := fmt.Sprintf(`SELECT COUNT(*) FROM markov_channels WHERE server_id = %v;`, interaction.GuildID)
		_ = db.QueryRow(query).Scan(&count)
		switch {
		case count == 0:
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "I am not allowed to generate a chain in any channels currently!",
				},
			})
		default:
			channelIDs := getGuildMarkovChannels(interaction.GuildID)

			// Constructing the string of channels that she is allowed to respond in.
			var response string
			log.Println(channelIDs)
			for _, chid := range channelIDs {
				response += fmt.Sprintln("\t<#" + chid + ">")
			}

			response = fmt.Sprintln("I am allowed to generate a chain in the following channel(s):\n", response)

			// Finally, responding to the interaction.
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: response,
				},
			})
		}
	},
	"markov_add_channel": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
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
		query := fmt.Sprintf(`SELECT channel_id FROM markov_channels WHERE server_id = %v AND channel_id = %v;`,
			interaction.GuildID, interaction.ChannelID)
		err := db.QueryRow(query).Scan(&channelID)
		if err != nil && err != sql.ErrNoRows {
			log.Printf("%vERROR%v - COULD NOT QUERY CHANNELS: %v", Red, Reset, err)
		}

		if err == sql.ErrNoRows {
			// Snagging the target channel ID.
			channelID = interaction.ApplicationCommandData().Options[0].ChannelValue(session).ID

			// Creating and executing a SQL query to add the channel to the list of approved channels.
			query = fmt.Sprintf(`INSERT INTO markov_channels(server_id, channel_id) VALUES (%v, %v)`, serverID, channelID)

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
					Content: fmt.Sprintf("Successfully added <#%v> to the list of channels I am allowed to generate a chain in.", channelID),
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
	"markov_remove_channel": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
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
		query := fmt.Sprintf(`SELECT server_id FROM markov_channels WHERE channel_id = %v;`,
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
			query = fmt.Sprintf(`DELETE FROM markov_channels WHERE server_id = %v AND channel_id = %v`, serverID, channelID)

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
					Content: fmt.Sprintf("Successfully removed <#%v> from the list of channels I am allowed to generate a chain in.", channelID),
				},
			})
		}
	},
	"markov_get_chance": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
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

		// Responding to the interaction.
		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("The current response chance is %v percent.",
					getGuildMarkovChance(interaction.GuildID)),
			},
		})
	},
	"markov_set_chance": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
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

		chance := interaction.ApplicationCommandData().Options[0].FloatValue()

		// Creating a single row query to set the chance of the chat response chance.
		query := fmt.Sprintf(`UPDATE servers SET markov_frequency = %v WHERE server_id = %v`, chance, interaction.GuildID)
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
				Content: fmt.Sprintf("Successfully updated the markov chance. The reponse chance is now %v percent.", chance),
			},
		})
	},
	"owo": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		owo := getOwO(interaction.GuildID)

		if owo {
			query := fmt.Sprintf(`UPDATE servers SET markov_owo = 0 WHERE server_id = %v`, interaction.GuildID)
			_, err := db.Query(query)
			if err != nil {
				log.Printf("%vERROR%v - COULD NOT QUERY SERVERS: %v", Red, Reset, err)
			}

			// Responding to the interaction.
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(":("),
				},
			})
		} else {
			query := fmt.Sprintf(`UPDATE servers SET markov_owo = 1 WHERE server_id = %v`, interaction.GuildID)
			_, err := db.Query(query)
			if err != nil {
				log.Printf("%vERROR%v - COULD NOT QUERY SERVERS: %v", Red, Reset, err)
			}

			// Responding to the interaction.
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(":)"),
				},
			})
		}
	},
}
