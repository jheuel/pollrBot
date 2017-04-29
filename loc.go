package main

import "fmt"

var (
	locGotQuestion          = "OK now that we got a question, please send answer options to your poll."
	locStartCommand         = "/start"
	locCreateNewPoll        = "create new poll"
	locInlineInsertPoll     = "insert poll into chat"
	locNewQuestion          = "Great! Send a question for your new poll, please."
	locFinishedCreatingPoll = "Finished creating a new poll\n\nPreview:\n"
	locMainMenu             = "I can help you create, send and manage polls.\n\nWhat do you want to do?"
	locPollDoneButton       = "finalize"
	locAddedOption          = fmt.Sprintf(
		"You can add more options by sending messages each containing one option. If you are done, please push the %s button.\n\nPreview:\n",
		locPollDoneButton)
)
