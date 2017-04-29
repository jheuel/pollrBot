package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(2)
	}
}

func run() error {
	databaseFileName := os.Getenv("DB")
	if databaseFileName == "" {
		return fmt.Errorf("Did not find database file name $DB\n")
	}

	APIToken := os.Getenv("APITOKEN")
	if APIToken == "" {
		return fmt.Errorf("Did not find telegram API token $APITOKEN\n")
	}

	var st Store = NewSQLStore(databaseFileName)
	defer st.Close()

	bot, err := tgbotapi.NewBotAPI(APIToken)
	if err != nil {
		return fmt.Errorf("could not start bot: %v\n", err)
	}

	// bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		return fmt.Errorf("could not prepare update channel: %v", err)
	}

	for update := range updates {
		// INLINE QUERIES
		if update.InlineQuery != nil {
			log.Printf("InlineQuery from [%s]: %s", update.InlineQuery.From.UserName, update.InlineQuery.Query)

			err = st.SaveUser(update.InlineQuery.From)
			if err != nil {
				log.Printf("could not save user: %v", err)
			}

			err = handleInlineQuery(bot, update, st)
			if err != nil {
				log.Printf("could not handle inline query: %v", err)
			}

			continue
		}

		// CALLBACK QUERIES
		if update.CallbackQuery != nil {
			log.Printf("CallbackQuery from [%s]: %s", update.CallbackQuery.From.UserName, update.CallbackQuery.Data)

			err = st.SaveUser(update.CallbackQuery.From)
			if err != nil {
				log.Printf("could not save user: %v", err)
			}

			err = handleCallbackQuery(bot, update, st)
			if err != nil {
				log.Printf("could not handle callback query: %v", err)
			}

			continue
		}

		if update.Message == nil {
			continue
		}

		err = st.SaveUser(update.Message.From)
		if err != nil {
			log.Printf("could not save user: %v", err)
		}

		// Messages
		log.Printf("Message from [%s] %s", update.Message.From.UserName, update.Message.Text)

		// Conversations
		err = handleDialog(bot, update, st)
		if err != nil {
			log.Printf("could not handle dialog: %v", err)
		}
	}
	return nil
}
