package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func handleInlineQuery(bot *tgbotapi.BotAPI, update tgbotapi.Update, st Store) error {
	polls, err := st.GetPollsByUser(update.InlineQuery.From.ID)
	if err != nil {
		return fmt.Errorf("could not get polls for user: %v", err)
	}

	if len(polls) > maxPollsInlineQuery {
		polls = polls[0 : maxPollsInlineQuery-1]
	}
	results := make([]interface{}, len(polls))
	for i, p := range polls {
		log.Println(p)
		article := tgbotapi.NewInlineQueryResultArticleMarkdown(strconv.Itoa(p.ID), p.Question, buildPollListing(p, st))
		if len(p.Options) > 0 {
			article.ReplyMarkup = buildPollMarkup(p)
		}
		article.Description = locInlineInsertPoll

		results[i] = article

	}
	inlineConfig := tgbotapi.InlineConfig{
		InlineQueryID:     update.InlineQuery.ID,
		Results:           results,
		IsPersonal:        true,
		CacheTime:         0,
		SwitchPMText:      locCreateNewPoll,
		SwitchPMParameter: createNewPollQuery,
	}

	_, err = bot.AnswerInlineQuery(inlineConfig)
	if err != nil {
		return fmt.Errorf("could not answer inline query: %v", err)
	}

	return nil
}
