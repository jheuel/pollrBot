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

	if strings.Contains(update.Message.Text, locAboutCommand) {
		msg := tgbotapi.NewMessage(int64(update.Message.From.ID), locAboutMessage)
		_, err = bot.Send(&msg)
		if err != nil {
			return fmt.Errorf("could not send message: %v", err)
		}
		return err
	}

	if strings.Contains(update.Message.Text, locDonateCommand) {
		fmt.Println("found donate command")
		msg := tgbotapi.NewMessage(int64(update.Message.From.ID), locDonateMessage)
		msg.ParseMode = tgbotapi.ModeHTML
		_, err = bot.Send(&msg)
		if err != nil {
			return fmt.Errorf("could not send message: %v", err)
		}
		return err
	}

	state, pollid, err = st.GetState(update.Message.From.ID)
	if err != nil {
		// could not retrieve state -> state is zero
		state = ohHi
		log.Printf("could not get state from database: %v\n", err)
	}

	if strings.Contains(update.Message.Text, locEditCommand) {
		polls, err := st.GetPollsByUser(update.Message.From.ID)
		if err != nil || len(polls) == 0 {
			log.Printf("could not get polls of user with userid %d: %v", update.Message.From.ID, err)
			msg := tgbotapi.NewMessage(int64(update.Message.From.ID), locNoMessageToEdit)
			_, err = bot.Send(&msg)
			if err != nil {
				return fmt.Errorf("could not send message: %v", err)
			}
			return fmt.Errorf("could not find message to edit: %v", err)
		}

		var p *poll
		for _, p = range polls {
			if p.ID == pollid {
				break
			}
		}

		_, err = sendEditMessage(bot, update, p)
		if err != nil {
			return fmt.Errorf("could not send edit message: %v", err)
		}
		return nil
	}

	if strings.Contains(update.Message.Text, "/start") || pollid < 0 && state != waitingForQuestion {
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

	if state == editQuestion {
		p, err := st.GetPoll(pollid)
		if err != nil {
			return fmt.Errorf("could not get poll: %v", err)
		}

		p.Question = update.Message.Text

		_, err = st.SavePoll(p)
		if err != nil {
			return fmt.Errorf("could not save poll: %v", err)
		}

		msg := tgbotapi.NewMessage(
			int64(update.Message.From.ID),
			fmt.Sprintf(locGotEditQuestion, p.Question))
		_, err = bot.Send(&msg)
		if err != nil {
			return fmt.Errorf("could not send message: %v", err)
		}

		state = editPoll
		err = st.SaveState(update.Message.From.ID, pollid, state)
		if err != nil {
			return fmt.Errorf("could not save state: %v", err)
		}
		//return nil
	}

	if state == editPoll {
		p, err := st.GetPoll(pollid)
		if err != nil {
			return fmt.Errorf("could not get poll: %v", err)
		}

		body := "This is the poll currently selected:\n<pre>\n"
		body += p.Question + "\n"
		for i, o := range p.Options {
			body += fmt.Sprintf("%d. %s", i+1, o.Text) + "\n"
		}
		body += "</pre>\n\n"

		msg := tgbotapi.NewMessage(
			update.Message.Chat.ID,
			body)
		msg.ParseMode = tgbotapi.ModeHTML
		msg.ReplyMarkup = buildEditMarkup(p, false, false)

		_, err = bot.Send(msg)
		if err != nil {
			return fmt.Errorf("could not send message: %v", err)
		}
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
