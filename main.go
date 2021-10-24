package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	//"github.com/go-telegram-bot-api/telegram-bot-api"

	_ "github.com/mattn/go-sqlite3"
)

var c crypt
var adminID int64 = 3761925

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(2)
	}
}

var pollsToUpdateConstRate = make(chan int, 200)
var pollsToUpdate = newUniqueChan()

func newUniqueChan() *uniqueChan {
	return &uniqueChan{
		C:   make(chan int, 10000),
		ids: make(map[int]struct{})}
}

type uniqueChan struct {
	C   chan int
	ids map[int]struct{}
}

func (u *uniqueChan) enqueue(id int) {
	if _, ok := u.ids[id]; ok {
		log.Printf("Update for poll #%d is already scheduled (%d scheduled updates).\n", id, len(u.ids))
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

func newTimer(name string) func() {
	start := time.Now()
	return func() {
		log.Printf("%s took: %s", name, time.Since(start))
	}
}

func run() error {
	// fill update channel with constant rate
	go func() {
		var pollid int
		for {
			time.Sleep(900 * time.Millisecond)
			pollid = pollsToUpdate.dequeue()
			pollsToUpdateConstRate <- pollid
		}
	}()

	webhookURL := os.Getenv("URL")
	if webhookURL == "" {
		return fmt.Errorf("Did not find webhook URL $URL")
	}

	key := os.Getenv("SECRET")
	if key == "" {
		return fmt.Errorf("Did not find secret key $SECRET")
	}
	c = newCrypt(key)

	databaseFileName := os.Getenv("DB")
	if databaseFileName == "" {
		return fmt.Errorf("Did not find database file name $DB")
	}

	APIToken := os.Getenv("APITOKEN")
	if APIToken == "" {
		return fmt.Errorf("Did not find telegram API token $APITOKEN")
	}

	ADMINID := os.Getenv("ADMINID")
	if ADMINID != "" {
		id, err := strconv.Atoi(ADMINID)
		if err != nil {
			return fmt.Errorf("Could not set adminID to %s: %v", ADMINID, err)
		}
		adminID = int64(id)
	}

	var st Store = newSQLStore(databaseFileName)
	defer st.Close()

	client := &http.Client{
		Timeout: time.Second * 1,
	}

	bot, err := tgbotapi.NewBotAPIWithClient(APIToken, client)
	if err != nil {
		return fmt.Errorf("could not start bot: %v", err)
	}

	//bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(webhookURL + bot.Token))
	if err != nil {
		log.Fatal(err)
	}
	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}
	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}
	updates := bot.ListenForWebhook("/" + bot.Token)

	// static files
	staticFiles := []string{
		"favicon.ico",
		"manifest.json",
		"asset-manifest.json",
		"robots.txt",
		"service-worker.js",
		"logo192.png",
		"logo512.png",
		"index.html",
	}

	for _, staticFile := range staticFiles {
		f := fmt.Sprintf("/%s", staticFile)
		http.Handle(f, serveStatic("/home/jheuel/services/pollrBot/static/"+f))
	}
	http.Handle("/", serveStatic("/home/jheuel/services/pollrBot/static/index.html"))

	http.HandleFunc("/poll/", pollHandler(st))
	http.Handle("/messageme/", contactHandler(bot))

	fs := http.FileServer(http.Dir("/home/jheuel/services/pollrBot/static/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	go http.ListenAndServe("127.0.0.1:8765", nil)

	reloadTimer := time.NewTimer(24 * time.Hour)
	for {
		select {
		case <-reloadTimer.C:
			return fmt.Errorf("Reload")
		case pollid := <-pollsToUpdateConstRate:
			stopTimer := newTimer("update poll")
			log.Printf("Updating poll %d\n", pollid)
			err := updatePollMessages(bot, pollid, st)
			if err != nil {
				log.Printf("Could not update poll #%d: %v", pollid, err)
			}
			stopTimer()
		case update := <-updates:
			stopTimer := newTimer("handle updates pdate")
			defer stopTimer()

			// INLINE QUERIES
			if update.InlineQuery != nil {
				log.Printf("InlineQuery from [%s]: %s", update.InlineQuery.From.UserName, update.InlineQuery.Query)

				err = st.SaveUser(update.InlineQuery.From)
				if err != nil {
					log.Printf("could not save user: %v", err)
				}

				if update.InlineQuery.From.ID == int(adminID) {
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
