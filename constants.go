package main

const createNewPollQuery = "createNewPoll"
const createPollQuery = "createpoll"
const pollDoneQuery = "polldone"

const (
	ohHi = iota
	waitingForQuestion
	waitingForOption
	pollDone
	editPoll
	editQuestion
	addOption
)

const (
	open = iota
	inactive
)

var maxNumberOfUsersListed = 100
var maxPollsInlineQuery = 10
