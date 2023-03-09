package main

import "github.com/bwmarrin/discordgo"

// Decalaring default member permission.
var defaultMemberPermissions int64 = discordgo.PermissionAdministrator

// Declaring min and max values of the chance command option.
var minChanceValue float64 = 0
var maxChanceValue float64 = 100

// Declaring the max value allowed for a response.
var minLengthValue float64 = 64
var maxLengthValue float64 = 1024

var commandMap = []*discordgo.ApplicationCommand{
	{
		Name:                     "echo",
		Description:              "This echoes your text to the specified channel as Shem-Ha.",
		DefaultMemberPermissions: &defaultMemberPermissions,

		// Registering the option available for this command.
		// https://pkg.go.dev/github.com/bwmarrin/discordgo#ApplicationCommandOption
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionChannel,
				Name:        "channel",
				Description: "This is the specified channel that you want to echo your message to.",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "text",
				Description: "This is the text that you want Shem-Ha to echo.",
				Required:    true,
			},
		},
	},
	{
		Name:                     "setup",
		Description:              "This command prepares Shem-Ha for use in your server. This must be run first!",
		DefaultMemberPermissions: &defaultMemberPermissions,
	},
	{
		Name:                     "chat_list_channels",
		Description:              "This command lists the channels Shem-Ha is allowed to respond in.",
		DefaultMemberPermissions: &defaultMemberPermissions,
	},
	{
		Name:                     "chat_add_channel",
		Description:              "This command adds a channel to the list that Shem-Ha is allowed to chat in.",
		DefaultMemberPermissions: &defaultMemberPermissions,

		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionChannel,
				Name:        "channel",
				Description: "The channel that you want to add to the list of approved channels.",
				Required:    true,
			},
		},
	},
	{
		Name:                     "chat_remove_channel",
		Description:              "This command removes a channel to the list that Shem-Ha is allowed to chat in.",
		DefaultMemberPermissions: &defaultMemberPermissions,

		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionChannel,
				Name:        "channel",
				Description: "The channel that you want to remove from the list of approved channels.",
				Required:    true,
			},
		},
	},
	{
		Name:                     "chat_get_chance",
		Description:              "This command checks the chat response chance.",
		DefaultMemberPermissions: &defaultMemberPermissions,
	},
	{
		Name:                     "chat_set_chance",
		Description:              "This command sets the chat response chance.",
		DefaultMemberPermissions: &defaultMemberPermissions,

		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionNumber,
				Name:        "percentage",
				Description: "This value is the chance that Shem-Ha will respond to a message, must be between 0 and 100.",
				Required:    true,
				MinValue:    &minChanceValue,
				MaxValue:    maxChanceValue,
			},
		},
	},

	//{
	//	Name:                     "list_channels",
	//	Description:              "This command shows all the channels that Shem-Ha is allowed to chat in!",
	//	DefaultMemberPermissions: &defaultMemberPermissions,
	//},
	//{
	//	Name:                     "get_chance",
	//	Description:              "This returns the value of the chance that Shem-Ha will respond to a message.",
	//	DefaultMemberPermissions: &defaultMemberPermissions,
	//},
	//{
	//	Name:                     "get_length",
	//	Description:              "This returns the maximum length of a response from Shem-Ha in tokens. A token is about 4 characters.",
	//	DefaultMemberPermissions: &defaultMemberPermissions,
	//},
	//{
	//	Name:                     "set_chance",
	//	Description:              "This sets the value of the chance that Shem-Ha will respond to a message.",
	//	DefaultMemberPermissions: &defaultMemberPermissions,
	//	// Registering the option available for this command.
	//	// https://pkg.go.dev/github.com/bwmarrin/discordgo#ApplicationCommandOption
	//	Options: []*discordgo.ApplicationCommandOption{
	//		{
	//			Type:        discordgo.ApplicationCommandOptionNumber,
	//			Name:        "percentage",
	//			Description: "This value is the chance that Shem-Ha will respond to a message, must be between 0 and 100.",
	//			Required:    true,
	//			MinValue:    &minChanceValue,
	//			MaxValue:    maxChanceValue,
	//		},
	//	},
	//},
	//{
	//	Name:                     "set_length",
	//	Description:              "This sets the maximum length of a response from Shem-Ha in tokens. A token is about 4 characters.",
	//	DefaultMemberPermissions: &defaultMemberPermissions,
	//	// Registering the option available for this command.
	//	// https://pkg.go.dev/github.com/bwmarrin/discordgo#ApplicationCommandOption
	//	Options: []*discordgo.ApplicationCommandOption{
	//		{
	//			Type:        discordgo.ApplicationCommandOptionInteger,
	//			Name:        "tokens",
	//			Description: "This is the maximum response length in tokens. A token is about 4 characters.",
	//			Required:    true,
	//			MinValue:    &minLengthValue,
	//			MaxValue:    maxLengthValue,
	//		},
	//	},
	//},
	//{
	//	Name:                     "add_channel",
	//	Description:              "This adds a channel to the list of channels Shem-Ha is allowed to reply in.",
	//	DefaultMemberPermissions: &defaultMemberPermissions,
	//	// Registering the option available for this command.
	//	// https://pkg.go.dev/github.com/bwmarrin/discordgo#ApplicationCommandOption
	//	Options: []*discordgo.ApplicationCommandOption{
	//		{
	//			Type:        discordgo.ApplicationCommandOptionChannel,
	//			Name:        "channel",
	//			Description: "The channel that you want to add to the list of approved channels.",
	//			Required:    true,
	//		},
	//	},
	//},
	//{
	//	Name:                     "remove_channel",
	//	Description:              "This removes a channel from the list of channels Shem-Ha is allowed to reply in.",
	//	DefaultMemberPermissions: &defaultMemberPermissions,
	//	// Registering the option available for this command.
	//	// https://pkg.go.dev/github.com/bwmarrin/discordgo#ApplicationCommandOption
	//	Options: []*discordgo.ApplicationCommandOption{
	//		{
	//			Type:        discordgo.ApplicationCommandOptionChannel,
	//			Name:        "channel",
	//			Description: "The channel that you want to remove from the list of approved channels.",
	//			Required:    true,
	//		},
	//	},
	//},
}
