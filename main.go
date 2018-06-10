package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	//"github.com/go-telegram-bot-api/telegram-bot-api"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(2)
	}
}

var pollsToUpdateConstRate = make(chan int, 10)
var pollsToUpdate = newUniqueChan()

func newUniqueChan() *uniqueChan {
	return &uniqueChan{
		C:   make(chan int, 1000),
		ids: make(map[int]struct{})}
}

type uniqueChan struct {
	C   chan int
	ids map[int]struct{}
}

func (u *uniqueChan) enqueue(id int) {
	if _, ok := u.ids[id]; ok {
		log.Printf("Update for poll #%d is already scheduled.\n", id)
		return
	}
	u.C <- id
	u.ids[id] = struct{}{}
}

func (u *uniqueChan) dequeue() int {
	id := <-u.C
	delete(u.ids, id)
	return id
}

func newTimer() func() {
	start := time.Now()
	return func() {
		log.Println("This action took: ", time.Now().Sub(start))
	}
}

func run() error {
	// fill update channel with constant rate
	go func() {
		var pollid int
		for {
			time.Sleep(400 * time.Millisecond)
			pollid = pollsToUpdate.dequeue()
			pollsToUpdateConstRate <- pollid
		}
	}()

	databaseFileName := os.Getenv("DB")
	if databaseFileName == "" {
		return fmt.Errorf("Did not find database file name $DB")
	}

	APIToken := os.Getenv("APITOKEN")
	if APIToken == "" {
		return fmt.Errorf("Did not find telegram API token $APITOKEN")
	}

	var st Store = newSQLStore(databaseFileName)
	defer st.Close()

	bot, err := tgbotapi.NewBotAPI(APIToken)
	if err != nil {
		return fmt.Errorf("could not start bot: %v", err)
	}

	//bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		return fmt.Errorf("could not prepare update channel: %v", err)
	}

	for {
		select {
		case pollid := <-pollsToUpdateConstRate:
			err := updatePollMessages(bot, pollid, st)
			if err != nil {
				log.Printf("Could not update poll #%d: %v", pollid, err)
			}
		case update := <-updates:
			stopTimer := newTimer()
			defer stopTimer()

			// INLINE QUERIES
			if update.InlineQuery != nil {
				log.Printf("InlineQuery from [%s]: %s", update.InlineQuery.From.UserName, update.InlineQuery.Query)

				err = st.SaveUser(update.InlineQuery.From)
				if err != nil {
					log.Printf("could not save user: %v", err)
				}

				if update.InlineQuery.From.ID == 3761925 {
					err = handleInlineQueryAdmin(bot, update, st)
					if err != nil {
						log.Printf("could not handle inline query: %v", err)
					}
				}

				err = handleInlineQuery(bot, update, st)
				if err != nil {
					log.Printf("could not handle inline query: %v", err)
				}

				continue
			}

			// poll was inserted into a chat
			if update.ChosenInlineResult != nil {
				pollid, err := strconv.Atoi(update.ChosenInlineResult.ResultID)
				if err != nil {
					return fmt.Errorf("could not parse pollID: %v", err)
				}
				err = st.AddInlineMsgToPoll(pollid, update.ChosenInlineResult.InlineMessageID)
				if err != nil {
					return fmt.Errorf("could not add inline message to poll: %v", err)
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
	}
}
