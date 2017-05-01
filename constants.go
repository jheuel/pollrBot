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

var maxNumberOfUsersListed = 10
var maxPollsInlineQuery = 10

var emojinumbers = map[string]string{
	"0": "️0️⃣",
	"1": "1️⃣",
	"2": "2️⃣",
	"3": "3️⃣",
	"4": "4️⃣",
	"5": "5️⃣",
	"6": "6️⃣",
	"7": "7️⃣",
	"8": "8️⃣",
	"9": "9️⃣",
}
