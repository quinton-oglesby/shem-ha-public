package main

import (
	"database/sql"
	"fmt"
	"log"
)

var (
	Reset  = "\033[0m"
	Bold   = "\033[1m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Gray   = "\033[37m"
	White  = "\033[97m"
)

// Function to check whether a guild is already set up.
func guildIsSetup(guildID string) bool {
	var id int64
	// Perform a single row query to check if the guild is setup.
	query := fmt.Sprintf(`SELECT server_id FROM servers WHERE server_id= %v`, guildID)
	err := db.QueryRow(query).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("%vERROR%v - COULD NOT QUERY SERVERS:\n\t%v", Red, Reset, err)
		return false
	}

	if err == sql.ErrNoRows {
		return false
	} else {
		return true
	}
}

func chatGetChannels(guildID string) []string {
	query := fmt.Sprintf(`SELECT * FROM chat_channels WHERE server_id = %v;`, guildID)

	rows, err := db.Query(query)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("%vERROR%v - COULD NOT QUERY CHANNELS:\n\t%v", Red, Reset, err)
		return []string{}
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

		channelIDs = append(channelIDs, channelID)
	}

	return channelIDs
}

func guidIsExempt(guildID string) bool {
	var exempt int

	// Performa a single row query to check if the guild is exempt from paying.
	query := fmt.Sprintf(`SELECT exempt from servers WHERE server_id = %v`, guildID)
	err := db.QueryRow(query).Scan(&exempt)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("%vERROR%v - COULD NOT QUERY SERVERS:\n\t%v", Red, Reset, err)
		return false
	}

	if err == sql.ErrNoRows {
		return false
	} else if exempt == 0 {
		return false
	} else {
		return true
	}
}

func guildIsPaid(guildID string) bool {
	var paid int

	// Performa a single row query to check if the guild is exempt from paying.
	query := fmt.Sprintf(`SELECT paid from servers WHERE server_id = %v`, guildID)
	err := db.QueryRow(query).Scan(&paid)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("%vERROR%v - COULD NOT QUERY SERVERS:\n\t%v", Red, Reset, err)
		return false
	}

	if err == sql.ErrNoRows {
		return false
	} else if paid == 0 {
		return false
	} else {
		return true
	}
}

func getGuildChatChance(guildID string) float64 {
	var chance float64

	// Performa a single row query to check if the guild chat respone chance.
	query := fmt.Sprintf(`SELECT chat_chance from servers WHERE server_id = %v`, guildID)
	err := db.QueryRow(query).Scan(&chance)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("%vERROR%v - COULD NOT QUERY SERVERS:\n\t%v", Red, Reset, err)
		return 0
	}

	if err == sql.ErrNoRows {
		return 0
	} else {
		return chance
	}
}

func getGuildChatLength(guildID string) int {
	var length int

	// Performa a single row query to check if the guild chat respone chance.
	query := fmt.Sprintf(`SELECT chat_length from servers WHERE server_id = %v`, guildID)
	err := db.QueryRow(query).Scan(&length)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("%vERROR%v - COULD NOT QUERY SERVERS:\n\t%v", Red, Reset, err)
		return 0
	}

	if err == sql.ErrNoRows {
		return 0
	} else {
		return length
	}
}

// Function to check if a string is in an array, returns true or false.
func inArray(str string, list []string) bool {
	for _, i := range list {
		if i == str {
			return true
		}
	}

	return false
}

// Function to remove a string from an array, returning the newly updated array.
func removeFromArray(str string, list []string) []string {
	for i, j := range list {
		if j == str {
			return append(list[:i], list[i+1:]...)
		}
	}

	return list
}
