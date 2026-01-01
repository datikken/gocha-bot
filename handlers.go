package main

import (
	"log"
	"strings"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
		"inviteclient":  func(bot *tgbotapi.BotAPI, chatID int64, user *User, fullText string) {
			inviteClient(bot, chatID, user, fullText)
		},
		"clients": func(bot *tgbotapi.BotAPI, chatID int64, user *User, fullText string) {
			listClients(bot, chatID, user)
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
