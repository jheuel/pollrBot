package main

const createNewPollQuery = "createNewPoll"
const createPollQuery = "createpoll"
const pollDoneQuery = "polldone"

const (
	ohHi = iota
	waitingForQuestion
	waitingForOption
	pollDone
)

var maxNumberOfUsersListed = 100
var maxPollsInlineQuery = 5
