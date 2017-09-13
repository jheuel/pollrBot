package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func handleCallbackQuery(bot *tgbotapi.BotAPI, update tgbotapi.Update, st Store) error {
	if update.CallbackQuery.Data == createPollQuery {
		return sendNewQuestionMessage(bot, update, st)
	}

	if strings.Contains(update.CallbackQuery.Data, pollDoneQuery) {
		return handlePollDoneQuery(bot, update, st)
	}

	pollid, optionid, err := parseQueryPayload(update)
	if err != nil {
		return fmt.Errorf("could not parse query payload: %v", err)
	}

	if update.CallbackQuery.InlineMessageID != "" {
		if err := st.AddInlineMsgToPoll(pollid, update.CallbackQuery.InlineMessageID); err != nil {
			return fmt.Errorf("could not add inline message to poll: %v", err)
		}
	}

	newAnswer := answer{
		UserID:   update.CallbackQuery.From.ID,
		PollID:   pollid,
		OptionID: optionid,
	}
	unvoted, err := st.SaveAnswer(newAnswer)
	if err != nil {
		return fmt.Errorf("could not save answers: %v", err)
	}

	p, err := st.GetPoll(pollid)
	if err != nil {
		return fmt.Errorf("could not get poll: %v", err)
	}

	var choice option
	for _, o := range p.Options {
		if o.ID == newAnswer.OptionID {
			choice = o
		}
	}

	msgs, err := st.GetAllPollMsg(pollid)
	if err != nil {
		return fmt.Errorf("could not get all pollmsgs: %v", err)
	}

	var ed tgbotapi.EditMessageTextConfig
	ed.Text = buildPollListing(p, st)
	ed.ParseMode = tgbotapi.ModeMarkdown

	ed.ReplyMarkup = buildPollMarkup(p)

	for _, msg := range msgs {
		ed.ChatID = msg.ChatID
		ed.MessageID = msg.MessageID

		_, err = bot.Send(ed)
		if err != nil {
			return fmt.Errorf("could not update message: %v", err)
		}
	}
	// reset
	ed.ChatID = 0
	ed.MessageID = 0

	msgs, err = st.GetAllPollInlineMsg(p.ID)
	if err != nil {
		return fmt.Errorf("could not get all poll inline messages: %v", err)
	}

	for _, msg := range msgs {
		ed.InlineMessageID = msg.InlineMessageID
		_, err = bot.Send(ed)
		if err != nil {
			log.Println(fmt.Errorf("could not update inline message: %v", err))
		}
	}
	popupText := fmt.Sprintf(`You voted for "%s"`, choice.Text)
	if unvoted {
		popupText = fmt.Sprintf("Seems like you deleted your vote.")
	}

	callbackConfig := tgbotapi.NewCallback(
		update.CallbackQuery.ID,
		popupText)
	_, err = bot.AnswerCallbackQuery(callbackConfig)
	if err != nil {
		return fmt.Errorf("could not send answer to callback query: %v", err)
	}

	return nil
}

func handlePollDoneQuery(bot *tgbotapi.BotAPI, update tgbotapi.Update, st Store) error {
	splits := strings.Split(update.CallbackQuery.Data, ":")
	if len(splits) < 2 {
		return fmt.Errorf("query did not contain the pollid")
	}
	pollid, err := strconv.Atoi(splits[1])
	if err != nil {
		return fmt.Errorf("could not convert string payload to int: %v", err)
	}

	p, err := st.GetPoll(pollid)
	if err != nil {
		return fmt.Errorf("could not get poll: %v", err)
	}
	_, err = postPoll(bot, p, int64(update.CallbackQuery.From.ID))
	if err != nil {
		return fmt.Errorf("could not post finished poll: %v", err)
	}
	err = st.SaveState(update.CallbackQuery.From.ID, p.ID, pollDone)
	if err != nil {
		return fmt.Errorf("could not change state to poll done: %v", err)
	}
	return nil
}

func parseQueryPayload(update tgbotapi.Update) (pollid int, optionid int, err error) {
	dataSplit := strings.Split(update.CallbackQuery.Data, ":")
	if len(dataSplit) != 2 {
		return pollid, optionid, fmt.Errorf("could not parse response")
	}
	pollid, err = strconv.Atoi(dataSplit[0])
	if err != nil {
		return pollid, optionid, fmt.Errorf("could not convert CallbackQuery data pollid to int: %v", err)
	}

	optionid, err = strconv.Atoi(dataSplit[1])
	if err != nil {
		return pollid, optionid, fmt.Errorf("could not convert CallbackQuery data OptionID to int: %v", err)
	}
	return pollid, optionid, nil

}
