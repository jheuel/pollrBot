package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// Store is an interface for the persistent storage
// should allow easier swapping of databases
type Store interface {
	Close()
	AddMsgToPoll(pollid int, messageid int, chatid int64) error
	AddInlineMsgToPoll(pollid int, inlinemessageid string) error
	RemoveInlineMsg(inlinemessageid string) error
	GetPoll(pollid int) (*poll, error)
	GetUser(userid int) (*tgbotapi.User, error)
	GetPollsByUser(userid int) ([]*poll, error)
	GetPollID(messageid int) (int, error)
	GetPollNewer(pollid int, userid int) (*poll, error)
	GetPollOlder(pollid int, userid int) (*poll, error)
	GetAllPollMsg(pollid int) ([]pollident, error)
	GetAllPollInlineMsg(pollid int) ([]pollident, error)
	GetState(userid int) (state int, pollid int, err error)
	SaveState(userid int, pollid int, state int) error
	SaveUser(*tgbotapi.User) error
	SavePoll(*poll) (int, error)
	SaveOptions([]option) error
	SaveAnswer(*poll, answer) (unvoted bool, err error)
}

type answer struct {
	ID       int
	PollID   int
	UserID   int
	OptionID int
}

type option struct {
	ID     int
	PollID int
	Text   string
	Ctr    int
}

type poll struct {
	ID        int
	MessageID int
	UserID    int
	Question  string
	Inactive  int
	Private   int
	Type      int
	Options   []option
	Answers   []answer
}

func isInactive(poll *poll) bool {
	return poll.Inactive == inactive
}

func isMultipleChoice(poll *poll) bool {
	return poll.Type == multipleChoice
}
