package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
	"time"
)

func handleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	user, err := getOrCreateUser(message.From.ID, message.From.UserName)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		return
	}

	text := strings.TrimPrefix(message.Text, "/")

	commands := map[string]func(*tgbotapi.BotAPI, int64, *User, string){
		"start": func(bot *tgbotapi.BotAPI, chatID int64, user *User, fullText string) {
			sendStartMessage(bot, chatID, user)
		},
		"inviteclient": func(bot *tgbotapi.BotAPI, chatID int64, user *User, fullText string) {
			inviteClient(bot, chatID, user, fullText)
		},
		"clients": func(bot *tgbotapi.BotAPI, chatID int64, user *User, fullText string) {
			listClients(bot, chatID, user)
		},
		"visit": func(bot *tgbotapi.BotAPI, chatID int64, user *User, fullText string) {
			handleVisit(bot, chatID, user, fullText)
		},
	}

	textFields := strings.Fields(text)
	commandPart := textFields[0]

	if handler, exists := commands[commandPart]; exists {
		handler(bot, message.Chat.ID, user, message.Text)
	}
	if _, exists := commands[commandPart]; !exists {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Unknown command. Use /start or start for help.")
		bot.Send(msg)
	}
}

func sendStartMessage(bot *tgbotapi.BotAPI, chatID int64, user *User) {
	text := `
Welcome to the Barbershop Bot!
Use /inviteclient to save customer.
Use /clients to get list of customers
User /visit 2026-01-01 10:10 @username to create a visit for a client.
`

	msg := tgbotapi.NewMessage(chatID, text)
	bot.Send(msg)
}

func inviteClient(bot *tgbotapi.BotAPI, chatID int64, user *User, fullText string) {
	// Parse the command: /inviteclient @username
	parts := strings.Fields(fullText)
	if len(parts) < 2 {
		msg := tgbotapi.NewMessage(chatID, "Usage: /inviteclient @username")
		bot.Send(msg)
		return
	}

	username := strings.TrimPrefix(parts[1], "@")
	if username == "" {
		msg := tgbotapi.NewMessage(chatID, "Invalid username")
		bot.Send(msg)
		return
	}

	// Create invited user
	err := createInvitedUser(username, user.ID)
	if err != nil {
		log.Printf("Error creating invited user: %v", err)
		msg := tgbotapi.NewMessage(chatID, "Error inviting client")
		bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(chatID, "Client "+username+" invited successfully")
	bot.Send(msg)
}

func listClients(bot *tgbotapi.BotAPI, chatID int64, user *User) {
	clients, err := getMyClients(user.ID)
	if err != nil {
		log.Printf("Error getting clients: %v", err)
		msg := tgbotapi.NewMessage(chatID, "Error retrieving clients")
		bot.Send(msg)
		return
	}

	if len(clients) == 0 {
		msg := tgbotapi.NewMessage(chatID, "No clients found")
		bot.Send(msg)
		return
	}

	text := "Your clients:\n"
	for _, client := range clients {
		text += "- @" + client.Username + "\n"
	}

	msg := tgbotapi.NewMessage(chatID, text)
	bot.Send(msg)
}

func handleVisit(bot *tgbotapi.BotAPI, chatID int64, user *User, fullText string) {
	// Parse the command: /visit 2026-01-01 10:10 @username
	parts := strings.Fields(fullText)
	if len(parts) < 4 {
		msg := tgbotapi.NewMessage(chatID, "Usage: /visit YYYY-MM-DD HH:MM @username")
		bot.Send(msg)
		return
	}

	dateStr := parts[1]
	timeStr := parts[2]
	username := strings.TrimPrefix(parts[3], "@")

	if username == "" {
		msg := tgbotapi.NewMessage(chatID, "Invalid username")
		bot.Send(msg)
		return
	}

	// Parse date and time
	dateTimeStr := dateStr + " " + timeStr
	visitDate, err := time.Parse("2006-01-02 15:04", dateTimeStr)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Invalid date/time format. Use YYYY-MM-DD HH:MM")
		bot.Send(msg)
		return
	}

	// Check if the user is a client of the current user
	client, err := getUserByUsername(username)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Client not found")
		bot.Send(msg)
		return
	}

	// Verify the client was invited by this user
	if client.InvitedBy != user.ID {
		msg := tgbotapi.NewMessage(chatID, "You can only schedule visits for your own clients")
		bot.Send(msg)
		return
	}

	// Create the visit
	err = createVisit(client.ID, visitDate)
	if err != nil {
		log.Printf("Error creating visit: %v", err)
		msg := tgbotapi.NewMessage(chatID, "Error creating visit")
		bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(chatID, "Visit scheduled for "+username+" on "+dateStr+" at "+timeStr)
	bot.Send(msg)
}
