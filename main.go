package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/go-sql-driver/mysql"
)

// Creating a struct to hold the Discord and OpenAI Token held within tokens.json.
type Tokens struct {
	DiscordToken string
	OpenAIToken  string
}

// Making a struct to hold the MySQL server logon parameters.
type MySQLParameters struct {
	Username string
	Password string
	Database string
}

// Globalizing the Tokens struct to hold the data.
var tokens Tokens

// Globalizing the MySQLParameters struct to hold the data.
var mySQLParameters MySQLParameters

// Global variable to hold database connection.
var db *sql.DB

// Global variable to hold the regex string
var re *regexp.Regexp

func main() {

	// Retrieve the tokens from the tokens.json file.
	tokensFile, err := os.ReadFile("tokens.json")
	if err != nil {
		log.Fatal("COULD NOT READ 'tokens.json' FILE: ", err)
	}

	// Retrieve the parameters from db_data.json file.
	dbParameters, err := os.ReadFile("db_data.json")
	if err != nil {
		log.Println("Could not open sql_data file.")
		log.Fatal(err)
	}

	// Unmarshal the tokens and database parameters.
	json.Unmarshal(tokensFile, &tokens)
	json.Unmarshal(dbParameters, &mySQLParameters)

	// Compile regex string.
	re, _ = regexp.Compile(`([\w+]+\:\/\/)?([\w\d-]+\.)*[\w-]+[\.\:]\w+([\/\?\=\&\#\.]?[\w-]+)*\/?`)

	// Set up the parameters for the database connection.
	configuration := mysql.Config{
		User:   mySQLParameters.Username,
		Passwd: mySQLParameters.Password,
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: mySQLParameters.Database,
	}

	// Open a connection to the database.
	db, err = sql.Open("mysql", configuration.FormatDSN())
	if err != nil {
		log.Fatal("ERROR OPENING DATABASE CONNECTION: ", err)
	}

	// Create a new Discord session using the provided bot token.
	session, err := discordgo.New("Bot " + tokens.DiscordToken)
	if err != nil {
		log.Fatal("ERROR CREATING DISCORD SESSION: ", err)
	}

	// Identify we want all intents.
	session.Identify.Intents = discordgo.IntentsAll

	// Now we open a websocket connection to Discord and begin listening.
	err = session.Open()
	if err != nil {
		log.Fatal("ERROR OPENING WEBSOCKET: ", err)
	}

	// Making a map of registered commands.
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commandMap))

	// Looping through the commands array and registering them.
	// https://pkg.go.dev/github.com/bwmarrin/discordgo#Session.ApplicationCommandCreate
	for i, command := range commandMap {
		registeredCommand, err := session.ApplicationCommandCreate(session.State.User.ID, "", command)
		if err != nil {
			log.Printf("CANNOT CREATE '%v' COMMAND: %v", command.Name, err)
		}

		registeredCommands[i] = registeredCommand
	}

	// Looping through the array of interaction handlers and adding them to the session.
	session.AddHandler(func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		if handler, ok := commandHandlers[interaction.ApplicationCommandData().Name]; ok {
			handler(session, interaction)
		}
	})

	session.AddHandler(messageCreateChat)
	session.AddHandler(messageCreateMarkov)

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	//// // Lopping through the registeredCommands array and deleting all the commands.
	//for _, v := range registeredCommands {
	//	err := session.ApplicationCommandDelete(session.State.User.ID, "", v.ID)
	//	if err != nil {
	//		log.Printf("CANNOT DELETE '%v' COMMAND: %v", v.Name, err)
	//	}
	//}

	// Cleanly close the Discord session.
	session.Close()
}
