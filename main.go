package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"strings"
)

// Структура для хранения состояния пользователя
type UserState struct {
	Stage    string
	Choice   string
	UserData map[string]string
}

var userStates = make(map[int64]*UserState)

func getUserState(chatID int64) *UserState {
	if state, exists := userStates[chatID]; exists {
		return state
	}
	// Инициализация нового состояния
	state := &UserState{
		Stage:    "CHOOSING",
		UserData: make(map[string]string),
	}
	userStates[chatID] = state
	return state
}

func factsToStr(userData map[string]string) string {
	var facts []string
	for key, value := range userData {
		facts = append(facts, fmt.Sprintf("%s - %s", key, value))
	}
	return strings.Join(facts, "\n")
}

func handleStart(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	chatID := update.Message.Chat.ID
	state := getUserState(chatID)
	state.Stage = "CHOOSING"

	replyText := "Hi! My name is Doctor Botter."
	if len(state.UserData) > 0 {
		keys := make([]string, 0, len(state.UserData))
		for key := range state.UserData {
			keys = append(keys, key)
		}
		replyText += fmt.Sprintf(" You already told me your %s. Why don't you tell me something more?", strings.Join(keys, ", "))
	} else {
		replyText += " I will hold a more complex conversation with you. Why don't you tell me something about yourself?"
	}

	msg := tgbotapi.NewMessage(chatID, replyText)
	msg.ReplyMarkup = replyKeyboardMarkup()
	bot.Send(msg)
}

func handleMessage(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	chatID := update.Message.Chat.ID
	state := getUserState(chatID)

	switch state.Stage {
	case "CHOOSING":
		text := update.Message.Text
		if text == "Done" {
			handleDone(update, bot)
			return
		} else if text == "Something else..." {
			state.Stage = "TYPING_CHOICE"
			msg := tgbotapi.NewMessage(chatID, "Alright, please send me the category first, for example 'Most impressive skill'")
			bot.Send(msg)
		} else {
			state.Choice = text
			state.Stage = "TYPING_REPLY"
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Your %s? Yes, I would love to hear about that!", text))
			bot.Send(msg)
		}

	case "TYPING_REPLY":
		state.UserData[state.Choice] = update.Message.Text
		state.Choice = ""
		state.Stage = "CHOOSING"
		replyText := fmt.Sprintf("Neat! Just so you know, this is what you already told me:\n%s\nYou can tell me more, or change your opinion on something.", factsToStr(state.UserData))
		msg := tgbotapi.NewMessage(chatID, replyText)
		msg.ReplyMarkup = replyKeyboardMarkup()
		bot.Send(msg)

	case "TYPING_CHOICE":
		state.Choice = update.Message.Text
		state.Stage = "TYPING_REPLY"
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Your %s? Yes, I would love to hear about that!", state.Choice))
		bot.Send(msg)
	}
}

func handleDone(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	chatID := update.Message.Chat.ID
	state := getUserState(chatID)
	replyText := fmt.Sprintf("I learned these facts about you: \n%s\nUntil next time!", factsToStr(state.UserData))
	msg := tgbotapi.NewMessage(chatID, replyText)
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
	bot.Send(msg)

	// Удаление состояния пользователя
	delete(userStates, chatID)
}

func replyKeyboardMarkup() tgbotapi.ReplyKeyboardMarkup {
	keys := [][]tgbotapi.KeyboardButton{
		{tgbotapi.NewKeyboardButton("Age"), tgbotapi.NewKeyboardButton("Favourite colour")},
		{tgbotapi.NewKeyboardButton("Number of siblings"), tgbotapi.NewKeyboardButton("Something else...")},
		{tgbotapi.NewKeyboardButton("Done")},
	}
	return tgbotapi.NewReplyKeyboard(keys...)
}

func main() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not set")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				handleStart(update, bot)
			}
		} else {
			handleMessage(update, bot)
		}
	}
}
