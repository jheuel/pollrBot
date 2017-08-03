package main

import (
	"fmt"
	"log"
	"strconv"

	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func postPoll(bot *tgbotapi.BotAPI, p *poll, chatid int64) (tgbotapi.Message, error) {
	share := tgbotapi.InlineKeyboardButton{
		Text:              "share poll",
		SwitchInlineQuery: &p.Question,
	}
	new := tgbotapi.NewInlineKeyboardButtonData(locCreateNewPoll, createPollQuery)

	buttons := tgbotapi.NewInlineKeyboardRow(share, new)
	markup := tgbotapi.NewInlineKeyboardMarkup(buttons)
	messageTxt := locFinishedCreatingPoll
	messageTxt += p.Question + "\n\n"

	for i, o := range p.Options {
		messageTxt += strconv.Itoa(i+1) + ") " + o.Text + "\n"
	}
	msg := tgbotapi.NewMessage(chatid, messageTxt)
	msg.ReplyMarkup = markup

	return bot.Send(msg)
}

func sendMainMenuMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) (tgbotapi.Message, error) {
	buttons := make([]tgbotapi.InlineKeyboardButton, 0)
	buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData("create poll", createPollQuery))
	markup := tgbotapi.NewInlineKeyboardMarkup(buttons)
	messageTxt := locMainMenu
	msg := tgbotapi.NewMessage(int64(update.Message.From.ID), messageTxt)
	msg.ReplyMarkup = markup

	return bot.Send(msg)
}

func sendInterMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update, p *poll) (tgbotapi.Message, error) {
	buttons := make([]tgbotapi.InlineKeyboardButton, 0)
	buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(locPollDoneButton, fmt.Sprintf("%s:%d", pollDoneQuery, p.ID)))
	markup := tgbotapi.NewInlineKeyboardMarkup(buttons)
	messageTxt := locAddedOption
	messageTxt += p.Question + "\n\n"

	for i, o := range p.Options {
		messageTxt += strconv.Itoa(i+1) + ") " + o.Text + "\n"
	}
	msg := tgbotapi.NewMessage(int64(update.Message.From.ID), messageTxt)
	msg.ReplyMarkup = markup

	return bot.Send(msg)
}

func sendNewQuestionMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update, st Store) error {
	msg := tgbotapi.NewMessage(int64(update.CallbackQuery.From.ID), locNewQuestion)
	_, err := bot.Send(&msg)
	if err != nil {
		return fmt.Errorf("could not send message: %v", err)
	}

	err = st.SaveState(update.CallbackQuery.From.ID, -1, waitingForQuestion)
	if err != nil {
		return fmt.Errorf("could not change state to waiting for questions: %v", err)
	}
	return nil
}

func emojify(number int) string {
	str := strconv.Itoa(number)
	for number, emoji := range emojinumbers {
		str = strings.Replace(str, number, emoji, -1)
	}
	return str
}

func buildPollMarkup(p *poll) *tgbotapi.InlineKeyboardMarkup {
	buttonrows := make([][]tgbotapi.InlineKeyboardButton, 0) //len(p.Options), len(p.Options))
	row := -1

	for _, o := range p.Options {
		textWidth := 0
		if row != -1 {
			for _, b := range buttonrows[row] {
				textWidth += len(b.Text)
			}
		}
		textWidth += len(o.Text)
		if row == -1 || textWidth > 30 {
			row++
			buttonrows = append(buttonrows, make([]tgbotapi.InlineKeyboardButton, 0))
		}
		label := fmt.Sprintf("%s (%d)", o.Text, o.Ctr)
		callback := fmt.Sprintf("%d:%d", p.ID, o.ID)
		button := tgbotapi.NewInlineKeyboardButtonData(label, callback)
		buttonrows[row] = append(buttonrows[row], button)
	}
	markup := tgbotapi.NewInlineKeyboardMarkup(buttonrows...)

	return &markup
}

func buildPollListing(p *poll, st Store) (listing string) {
	listOfUsers := make([][]*tgbotapi.User, len(p.Options))
	for i, o := range p.Options {
		for _, a := range p.Answers {
			if a.OptionID == o.ID {
				u, err := st.GetUser(a.UserID)
				if err != nil {
					log.Printf("could not get user: %v", err)
					listOfUsers[i] = append(listOfUsers[i], &tgbotapi.User{ID: a.UserID})
					continue
				}
				listOfUsers[i] = append(listOfUsers[i], u)
			}
		}
	}

	listing += fmt.Sprintf("%s\n\n", p.Question)
	for i, o := range p.Options {
		var nr string
		if len(p.Options) < 10 {
			nr = emojify(i + 1)
		} else {
			nr = fmt.Sprintf("%d)", i+1)
		}

		var part string
		if len(p.Answers) > 0 {
			part = fmt.Sprintf(" (%.0f%%)", 100.*float64(o.Ctr)/float64(len(p.Answers)))
			// part = fmt.Sprintf("\n%s *%s* (%d/%d):\n ", emojify(i+1), o.Text, o.Ctr, len(p.Answers))
		}

		listing += fmt.Sprintf("\n%s *%s*%s", nr, o.Text, part)
		if len(p.Answers) < maxNumberOfUsersListed {
			listing += ":\n "
			for _, u := range listOfUsers[i] {
				listing += " " + getDisplayUserName(u) + ","
			}
			listing = listing[:len(listing)-1]
		}
		listing += "\n"

	}
	return listing
}

func getDisplayUserName(u *tgbotapi.User) string {
	//if u.UserName != "" {
	////return fmt.Sprintf(" @%s", u.UserName)
	//return fmt.Sprintf(" %s", u.UserName)
	//}

	if u.FirstName == "" && u.LastName == "" {
		return strconv.Itoa(u.ID)
	} else if u.FirstName != "" {
		name := u.FirstName
		if u.LastName != "" {
			name += " " + u.LastName
		}
		return name
	}
	return u.LastName
}
