package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type frontendPoll struct {
	Question string
	Options  []frontendOption
	Cnt      int
}

type frontendOption struct {
	Text    string
	Answers []frontendAnswer
	Cnt     int
}

type frontendAnswer struct {
	UserName string
}

func buildFrontendPoll(p *poll, st Store) (f frontendPoll) {
	polledUsers := make(map[int]struct{})
	votesForOption := make(map[int]int)
	f.Question = p.Question
	for _, o := range p.Options {
		var opt frontendOption
		opt.Text = o.Text
		for _, a := range p.Answers {
			if !isMultipleChoice(p) {
				if _, ok := polledUsers[a.UserID]; ok {
					continue
				}
			}
			if a.OptionID == o.ID {
				votesForOption[o.ID]++
				u, err := st.GetUser(a.UserID)
				if err != nil {
					log.Printf("could not get user: %v", err)
					continue
				}
				polledUsers[u.ID] = struct{}{}
				opt.Answers = append(opt.Answers, frontendAnswer{
					UserName: getDisplayUserName(u),
				})
			}
		}
		opt.Cnt = len(opt.Answers)
		f.Options = append(f.Options, opt)
	}
	f.Cnt = len(polledUsers)
	return f
}

func serveStatic(filename string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filename)
	})
}

func contactHandler(bot *tgbotapi.BotAPI) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		type Response struct {
			Success bool
		}

		var m struct {
			FirstName string
			LastName  string
			Email     string
			Message   string
		}
		err := json.NewDecoder(r.Body).Decode(&m)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Printf("Someone tried to contact you: %+v\n", m)
		message := fmt.Sprintf("%s %s (%s): %s\n", m.FirstName, m.LastName, m.Email, m.Message)
		msg := tgbotapi.NewMessage(adminID, message)
		_, err = bot.Send(msg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(Response{true})
	})
}

func pollHandler(st Store) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		p := strings.TrimPrefix(r.URL.Path, "/poll/")

		s, err := c.Decrypt(p)
		if err != nil {
			log.Printf("could not decrypt %s: %v", p, err)
		}
		id, err := strconv.Atoi(s)
		if err != nil {
			log.Printf("could not parse %s to int: %v", p, err)
			return
		}
		poll, err := st.GetPoll(id)
		if err != nil {
			log.Printf("could not get poll %d: %v", id, err)
			return
		}
		log.Printf("Prepare frontendPoll for poll %d\n", id)
		frontendPoll := buildFrontendPoll(poll, st)
		w.Header().Set("Content-Type", "application/json")

		// for local development
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
		json.NewEncoder(w).Encode(frontendPoll)
	}
}
