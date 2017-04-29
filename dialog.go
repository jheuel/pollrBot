package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func handleDialog(bot *tgbotapi.BotAPI, update tgbotapi.Update, st Store) error {
	state := ohHi
	pollid := -1
	var err error

	if !strings.Contains(update.Message.Text, locStartCommand) {
		state, pollid, err = st.GetState(update.Message.From.ID)
		if err != nil {
			// could not retrieve state -> state is zero
			state = ohHi
			log.Printf("could not get state from database: %v\n", err)
		}
	}

	if pollid < 0 && state != waitingForQuestion {
		state = ohHi
		err = st.SaveState(update.Message.From.ID, pollid, state)
		if err != nil {
			return fmt.Errorf("could not save state: %v", err)
		}
	}

	if state == ohHi {
		_, err = sendMainMenuMessage(bot, update)
		if err != nil {
			return fmt.Errorf("could not send main menu message: %v", err)
		}
		return nil
	}

	if state == waitingForQuestion {
		p := &poll{
			Question: update.Message.Text,
			UserID:   update.Message.From.ID,
		}

		pollid, err = st.SavePoll(p)
		if err != nil {
			return fmt.Errorf("could not save poll: %v", err)
		}

		msg := tgbotapi.NewMessage(int64(update.Message.From.ID), locGotQuestion)
		_, err = bot.Send(&msg)
		if err != nil {
			return fmt.Errorf("could not send message: %v", err)
		}

		state = waitingForOption
		err = st.SaveState(update.Message.From.ID, pollid, state)
		if err != nil {
			return fmt.Errorf("could not save state: %v", err)
		}

		return nil
	}

	if state == waitingForOption {
		opts := []option{
			option{
				PollID: pollid,
				Text:   update.Message.Text,
			}}

		err = st.SaveOptions(opts)
		if err != nil {
			return fmt.Errorf("could not save option: %v", err)
		}
		p, err := st.GetPoll(pollid)
		if err != nil {
			return fmt.Errorf("could not get poll: %v", err)
		}

		_, err = sendInterMessage(bot, update, p)
		if err != nil {
			return fmt.Errorf("could not send inter message: %v", err)
		}
		return nil
	}

	if state == pollDone {
		p, err := st.GetPoll(pollid)
		if err != nil {
			return fmt.Errorf("could not get poll: %v", err)
		}

		_, err = postPoll(bot, p, int64(update.Message.From.ID))
		if err != nil {
			return fmt.Errorf("could not post poll: %v", err)
		}
		return nil
	}

	return nil
}
