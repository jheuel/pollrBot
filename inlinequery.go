package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
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
		article := tgbotapi.NewInlineQueryResultArticleHTML(strconv.Itoa(p.ID), p.Question, buildPollListing(p, st))
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

func handleInlineQueryAdmin(bot *tgbotapi.BotAPI, update tgbotapi.Update, st Store) error {
	splits := strings.Split(update.InlineQuery.Query, ":")
	if len(splits) < 1 {
		return fmt.Errorf("Could not convert query to pollid")
	}
	active := false
	if len(splits) == 2 {
		if splits[1] == "p" {
			active = true
		}
	}

	pollid, err := strconv.Atoi(splits[0])
	if err != nil {
		return fmt.Errorf("Could not convert query to pollid: %v", err)
	}
	p, err := st.GetPoll(pollid)
	if err != nil {
		return fmt.Errorf("could not get polls for user: %v", err)
	}
	polls := []*poll{p}

	if len(polls) > maxPollsInlineQuery {
		polls = polls[0 : maxPollsInlineQuery-1]
	}
	results := make([]interface{}, len(polls))
	for i, p := range polls {
		log.Println(p)
		article := tgbotapi.NewInlineQueryResultArticleHTML(strconv.Itoa(p.ID), p.Question, buildPollListing(p, st))
		if len(p.Options) > 0 && active {
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
