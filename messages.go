package main

import (
	"fmt"
	"html"
	"log"
	"strconv"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/kyokomi/emoji"
)

func postPoll(bot *tgbotapi.BotAPI, p *poll, chatid int64) (tgbotapi.Message, error) {
	share := tgbotapi.InlineKeyboardButton{
		Text:              locSharePoll,
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
	//shareButton := tgbotapi.InlineKeyboardButton{
	//Text:              locSharePoll,
	//SwitchInlineQuery: &p.Question,
	//}
	pollDoneButton := tgbotapi.NewInlineKeyboardButtonData(
		locPollDoneButton, fmt.Sprintf("%s:%d", pollDoneQuery, p.ID))

	buttons := make([]tgbotapi.InlineKeyboardButton, 0)
	buttons = append(buttons, pollDoneButton)
	//buttons = append(buttons, shareButton)

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

func sendEditMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update, p *poll) (tgbotapi.Message, error) {
	body := "This is the poll currently selected:\n<pre>\n"
	body += p.Question + "\n"
	for i, o := range p.Options {
		body += fmt.Sprintf("%d. %s", i+1, o.Text) + "\n"
	}
	body += "</pre>\n\n"
	msg := tgbotapi.NewMessage(int64(update.Message.From.ID), body)
	msg.ParseMode = tgbotapi.ModeHTML

	msg.ReplyMarkup = buildEditMarkup(p, false, false)

	return bot.Send(&msg)
}

func buildPollMarkup(p *poll) *tgbotapi.InlineKeyboardMarkup {
	buttonrows := make([][]tgbotapi.InlineKeyboardButton, 0) //len(p.Options), len(p.Options))
	row := -1

	votesForOption := make(map[int]int)
	for _, o := range p.Options {
		for _, a := range p.Answers {
			if a.OptionID == o.ID {
				votesForOption[o.ID]++
			}
		}
	}

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
		label := fmt.Sprintf("%s (%d)", o.Text, votesForOption[o.ID])
		callback := fmt.Sprintf("%d:%d", p.ID, o.ID)
		button := tgbotapi.NewInlineKeyboardButtonData(label, callback)
		buttonrows[row] = append(buttonrows[row], button)
	}
	markup := tgbotapi.NewInlineKeyboardMarkup(buttonrows...)

	return &markup
}

func buildPollListing(p *poll, st Store) (listing string) {
	polledUsers := make(map[int]struct{})
	listOfUsers := make([][]*tgbotapi.User, len(p.Options))
	votesForOption := make(map[int]int)
	for i, o := range p.Options {
		for _, a := range p.Answers {
			if a.OptionID == o.ID {
				votesForOption[o.ID]++
				u, err := st.GetUser(a.UserID)
				if err != nil {
					log.Printf("could not get user: %v", err)
					listOfUsers[i] = append(listOfUsers[i], &tgbotapi.User{ID: a.UserID})
					continue
				}
				polledUsers[u.ID] = struct{}{}
				listOfUsers[i] = append(listOfUsers[i], u)
			}
		}
	}

	listing += emoji.Sprintf(":bar_chart:<b>%s</b>\n", html.EscapeString(p.Question))
	//log.Printf("Create listing for question: %s\n", p.Question)

	for i, o := range p.Options {
		var part string
		if len(p.Answers) > 0 {
			part = fmt.Sprintf(" (%.0f%%)", 100.*float64(votesForOption[o.ID])/float64(len(polledUsers)))
			if votesForOption[o.ID] != o.Ctr {
				log.Printf("Counter for option #%d is off: %d stored vs. %d counted", o.ID, o.Ctr, votesForOption[o.ID])
			}
		}
		listing += fmt.Sprintf("\n<b>%s</b>%s", html.EscapeString(o.Text), part)

		usersOnAnswer := len(listOfUsers[i])
		if len(p.Answers) < maxNumberOfUsersListed && usersOnAnswer > 0 {
			for j := 0; j+1 < usersOnAnswer; j++ {
				listing += "\n\u251C " + html.EscapeString(getDisplayUserName(listOfUsers[i][j]))
			}
			listing += "\n\u2514 " + html.EscapeString(getDisplayUserName(listOfUsers[i][usersOnAnswer-1]))
		}
		listing += "\n"
	}
	listing += emoji.Sprint(fmt.Sprintf("\n%d :busts_in_silhouette:\n", len(polledUsers)))
	return listing
}

func buildEditMarkup(p *poll, noOlder, noNewer bool) *tgbotapi.InlineKeyboardMarkup {
	query := fmt.Sprintf("e:%d", p.ID)

	buttonrows := make([][]tgbotapi.InlineKeyboardButton, 0)
	buttonrows = append(buttonrows, make([]tgbotapi.InlineKeyboardButton, 0))
	buttonrows = append(buttonrows, make([]tgbotapi.InlineKeyboardButton, 0))
	buttonrows = append(buttonrows, make([]tgbotapi.InlineKeyboardButton, 0))

	buttonLast := tgbotapi.NewInlineKeyboardButtonData("\u2B05", query+":-")
	buttonNext := tgbotapi.NewInlineKeyboardButtonData("\u27A1", query+":+")
	if noOlder {
		buttonLast = tgbotapi.NewInlineKeyboardButtonData("\u2B05", "dummy")
	}
	if noNewer {
		buttonNext = tgbotapi.NewInlineKeyboardButtonData("\u27A1", "dummy")
	}
	buttonrows[0] = append(buttonrows[0], buttonLast, buttonNext)
	buttonInactive := tgbotapi.NewInlineKeyboardButtonData(locToggleOpen, query+":c")
	if isInactive(p) {
		buttonInactive = tgbotapi.NewInlineKeyboardButtonData(locToggleInactive, query+":c")
	}
	buttonMultipleChoice := tgbotapi.NewInlineKeyboardButtonData(locToggleSingleChoice, query+":m")
	// if isMultipleChoice(p) {
	// 	buttonMultipleChoice = tgbotapi.NewInlineKeyboardButtonData(locToggleMultipleChoice, query+":m")
	// }
	buttonrows[1] = append(buttonrows[1], buttonInactive)
	if !isMultipleChoice(p) {
		buttonrows[1] = append(buttonrows[1], buttonMultipleChoice)
	}
	buttonEditQuestion := tgbotapi.NewInlineKeyboardButtonData(locEditQuestionButton, query+":q")
	buttonAddOptions := tgbotapi.NewInlineKeyboardButtonData(locAddOptionButton, query+":o")

	buttonrows[2] = append(buttonrows[2], buttonEditQuestion, buttonAddOptions)
	markup := tgbotapi.NewInlineKeyboardMarkup(buttonrows...)

	return &markup
}

func getDisplayUserName(u *tgbotapi.User) string {
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
