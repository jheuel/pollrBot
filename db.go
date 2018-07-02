package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type sqlStore struct {
	db *sql.DB
}

func (st *sqlStore) Close() {
	err := st.db.Close()
	if err != nil {
		log.Printf("could not close database properly: %v\n", err)
	}
}

type closable interface {
	Close() error
}

func close(c closable) {
	err := c.Close()
	if err != nil {
		log.Printf("could not close stmt or rows properly: %v\n", err)
	}
}

func newSQLStore(databaseFile string) *sqlStore {
	st := &sqlStore{}
	var err error
	st.db, err = sql.Open("sqlite3", databaseFile)
	if err != nil {
		log.Fatalf("could not open database: %v", err)
	}
	if err := st.db.Ping(); err != nil {
		log.Fatalf("could not open database: %v", err)
	}

	schema := `
	CREATE TABLE IF NOT EXISTS poll(
		ID INTEGER PRIMARY KEY ASC,
		UserID INTEGER,
		LastSaved INTEGER,
		CreatedAt INTEGER,
		Type INTEGER,
		Private INTEGER,
		Inactive INTEGER,
		Question TEXT);
	CREATE TABLE IF NOT EXISTS pollmsg(
		MessageID INTEGER PRIMARY KEY,
		ChatId INTEGER,
		PollID INTEGER);
	CREATE TABLE IF NOT EXISTS pollinlinemsg(
		InlineMessageID TEXT PRIMARY KEY,
		PollID INTEGER);
	CREATE TABLE IF NOT EXISTS answer(
		ID INTEGER PRIMARY KEY ASC,
		PollID INTEGER,
		OptionID INTEGER,
		LastSaved INTEGER,
		CreatedAt INTEGER,
		UserID INTEGER);
	CREATE TABLE IF NOT EXISTS option(
		ID INTEGER PRIMARY KEY ASC,
		PollID INTEGER,
		Ctr INTEGER,
		Text TEXT);
	CREATE TABLE IF NOT EXISTS dialog(
		UserID INTEGER PRIMARY KEY,
		PollId INTEGER,
		state INTEGER);
	CREATE TABLE IF NOT EXISTS user(
		ID INTEGER PRIMARY KEY,
		FirstName TEXT,
		LastName Text,
		LastSaved INTEGER,
		CreatedAt INTEGER,
		UserName TEXT);

	CREATE INDEX IF NOT EXISTS poll_index ON poll(ID);
	CREATE INDEX IF NOT EXISTS pollmsg_index ON pollmsg(MessageID);
	CREATE INDEX IF NOT EXISTS answer_index ON answer(PollID);
	CREATE INDEX IF NOT EXISTS option_index ON option(PollID);
	`

	if _, err := st.db.Exec(schema); err != nil {
		log.Fatalf("could not create schema: %v", err)
	}

	return st
}

func (st *sqlStore) GetUser(userid int) (*tgbotapi.User, error) {
	u := &tgbotapi.User{ID: userid}

	row := st.db.QueryRow("SELECT FirstName, LastName, UserName FROM user WHERE ID = ?", userid)
	if err := row.Scan(&u.FirstName, &u.LastName, &u.UserName); err != nil {
		return u, fmt.Errorf(`could not scan user "%d": %v`, u.ID, err)
	}

	return u, nil
}

func (st *sqlStore) GetPoll(pollid int) (*poll, error) {
	p := &poll{ID: pollid}
	var err error
	row := st.db.QueryRow("SELECT UserID, Question, Inactive, Type FROM poll WHERE ID = ?", pollid)
	if err := row.Scan(&p.UserID, &p.Question, &p.Inactive, &p.Type); err != nil {
		return p, fmt.Errorf("could not scan poll #%d: %v", p.ID, err)
	}

	p.Options, err = st.GetOptions(p.ID)
	if err != nil {
		return p, fmt.Errorf("could not query options: %v", err)
	}

	p.Answers, err = st.GetAnswers(p.ID)
	if err != nil {
		return p, fmt.Errorf("could not query answers: %v", err)
	}

	return p, nil
}

func (st *sqlStore) GetPollNewer(pollid int, userid int) (*poll, error) {
	p := &poll{}
	var err error
	row := st.db.QueryRow("SELECT UserID, ID, Question, Inactive, Type FROM poll WHERE ID > ? AND UserID = ? ORDER BY ID ASC LIMIT 1", pollid, userid)
	if err := row.Scan(&p.UserID, &p.ID, &p.Question, &p.Inactive, &p.Type); err != nil {
		return p, fmt.Errorf("could not scan poll #%d: %v", p.ID, err)
	}

	p.Options, err = st.GetOptions(p.ID)
	if err != nil {
		return p, fmt.Errorf("could not query options: %v", err)
	}

	p.Answers, err = st.GetAnswers(p.ID)
	if err != nil {
		return p, fmt.Errorf("could not query answers: %v", err)
	}

	return p, nil
}

func (st *sqlStore) GetPollOlder(pollid int, userid int) (*poll, error) {
	p := &poll{}
	var err error
	row := st.db.QueryRow("SELECT UserID, ID, Question, Inactive, Type FROM poll WHERE ID < ? AND UserID = ? ORDER BY ID DESC LIMIT 1", pollid, userid)
	if err := row.Scan(&p.UserID, &p.ID, &p.Question, &p.Inactive, &p.Type); err != nil {
		return p, fmt.Errorf("could not scan poll #%d: %v", p.ID, err)
	}

	p.Options, err = st.GetOptions(p.ID)
	if err != nil {
		return p, fmt.Errorf("could not query options: %v", err)
	}

	p.Answers, err = st.GetAnswers(p.ID)
	if err != nil {
		return p, fmt.Errorf("could not query answers: %v", err)
	}

	return p, nil
}

func (st *sqlStore) GetState(userid int) (state int, pollid int, err error) {
	row := st.db.QueryRow("SELECT state, pollid FROM dialog WHERE UserID = ?", userid)
	if err := row.Scan(&state, &pollid); err != nil {
		return state, pollid, fmt.Errorf("could not scan state from row: %v", err)
	}
	return state, pollid, nil
}

func (st *sqlStore) SaveState(userid int, pollid int, state int) (err error) {
	res, err := st.db.Exec("UPDATE dialog SET state = ? WHERE UserID = ?", userid, state)
	if err != nil {
		return fmt.Errorf("could not update state in database: %v", err)
	}

	if aff, err := res.RowsAffected(); aff == 0 || err != nil {
		_, err = st.db.Exec("INSERT OR REPLACE INTO dialog(UserID, pollid, state) values(?, ?, ?)", userid, pollid, state)
		if err != nil {
			return fmt.Errorf("could not insert or replace state database entry: %v", err)
		}
	}

	return nil
}

func (st *sqlStore) GetPollsByUser(userid int) ([]*poll, error) {
	polls := make([]*poll, 0)
	var err error

	row, err := st.db.Query("SELECT ID, UserID, Question, Inactive, Type FROM poll WHERE UserID = ? ORDER BY ID DESC LIMIT 3", userid)
	if err != nil {
		return polls, fmt.Errorf("could not query polls for userid #%d: %v", userid, err)
	}

	for row.Next() {
		p := &poll{UserID: userid}
		if err := row.Scan(&p.ID, &p.UserID, &p.Question, &p.Inactive, &p.Type); err != nil {
			return polls, fmt.Errorf("could not scan poll for userid #%d: %v", userid, err)
		}
		p.Options, err = st.GetOptions(p.ID)
		if err != nil {
			return polls, fmt.Errorf("could not query options: %v", err)
		}

		p.Answers, err = st.GetAnswers(p.ID)
		if err != nil {
			return polls, fmt.Errorf("could not query answers: %v", err)
		}

		polls = append(polls, p)
	}
	return polls, nil
}

func (st *sqlStore) GetPollID(messageid int) (int, error) {
	var pollid int

	rows, err := st.db.Query("SELECT PollID FROM pollmsg WHERE MessageID = ?", messageid)
	if err != nil {
		return pollid, fmt.Errorf("could not query pollid: %v", err)
	}
	defer close(rows)
	for rows.Next() {
		err = rows.Scan(&pollid)
		if err != nil {
			return pollid, fmt.Errorf("could not scan pollid: %v", err)
		}
	}
	return pollid, nil
}

type pollident struct {
	MessageID       int
	InlineMessageID string
	ChatID          int64
}

func (st *sqlStore) GetAllPollMsg(pollid int) ([]pollident, error) {
	msgs := make([]pollident, 0)

	rows, err := st.db.Query("SELECT MessageID, ChatID FROM pollmsg WHERE PollID = ?", pollid)
	if err != nil {
		return msgs, fmt.Errorf("could not query pollmsgs: %v", err)
	}
	defer close(rows)
	var msg pollident
	for rows.Next() {
		err = rows.Scan(&msg.MessageID, &msg.ChatID)
		if err != nil {
			return msgs, fmt.Errorf("could not scan pollmsgs: %v", err)
		}
		msgs = append(msgs, msg)
	}
	return msgs, nil
}

func (st *sqlStore) GetAllPollInlineMsg(pollid int) ([]pollident, error) {
	msgs := make([]pollident, 0)

	rows, err := st.db.Query("SELECT InlineMessageID FROM pollinlinemsg WHERE PollID = ?", pollid)
	if err != nil {
		return msgs, fmt.Errorf("could not query pollinlinemsg: %v", err)
	}
	defer close(rows)
	var msg pollident
	for rows.Next() {
		err = rows.Scan(&msg.InlineMessageID)
		if err != nil {
			return msgs, fmt.Errorf("could not scan pollid: %v", err)
		}
		msgs = append(msgs, msg)
	}
	return msgs, nil
}

func (st *sqlStore) GetOptions(pollid int) ([]option, error) {

	options := make([]option, 0)

	rows, err := st.db.Query("SELECT Ctr, PollID, ID, Text FROM option WHERE PollID = ?", pollid)
	if err != nil {
		return options, fmt.Errorf("could not query options: %v", err)
	}
	defer close(rows)
	var o option
	for rows.Next() {
		err = rows.Scan(&o.Ctr, &o.PollID, &o.ID, &o.Text)
		if err != nil {
			return options, fmt.Errorf("could not scan option: %v", err)
		}
		options = append(options, o)
	}
	return options, nil
}

func (st *sqlStore) GetAnswers(pollid int) ([]answer, error) {
	answers := make([]answer, 0)

	rows, err := st.db.Query("SELECT ID, PollID, OptionID, UserID FROM answer WHERE PollID = ?", pollid)
	if err != nil {
		return answers, fmt.Errorf("could not query answers: %v", err)
	}
	defer close(rows)
	var a answer
	for rows.Next() {
		err = rows.Scan(&a.ID, &a.PollID, &a.OptionID, &a.UserID)
		if err != nil {
			return answers, fmt.Errorf("could not scan answer: %v", err)
		}
		answers = append(answers, a)
	}
	return answers, nil
}

func (st *sqlStore) SaveAnswer(p *poll, a answer) (unvoted bool, err error) {
	tx, err := st.db.Begin()
	if err != nil {
		return false, fmt.Errorf("could not begin database transaction: %v", err)
	}
	defer func() {
		if err != nil {
			if err := tx.Rollback(); err != nil {
				log.Printf("could not rollback database change: %v", err)
			}
			return
		}
		err = tx.Commit()
	}()

	// find previous votes in this poll
	stmt, err := tx.Prepare("SELECT OptionID FROM answer WHERE PollID = ? AND UserID = ?")
	if err != nil {
		return false, fmt.Errorf("could not prepare sql statement: %v", err)
	}
	defer close(stmt)

	rows, err := stmt.Query(a.PollID, a.UserID)
	if err != nil {
		return false, fmt.Errorf("could not query option rows: %v", err)
	}
	defer close(rows)
	var optionid int
	prevOpts := make([]int, 0) // len should be 0 or 1
	for rows.Next() {
		err = rows.Scan(&optionid)
		if err != nil {
			return false, fmt.Errorf("could not scan optionid: %v", err)
		}
		prevOpts = append(prevOpts, optionid)
	}

	if len(prevOpts) > 0 { // user voted before

		// user clicked the same answer again
		if contains(prevOpts, a.OptionID) {
			stmt, err = tx.Prepare("DELETE FROM answer where PollID = ? AND UserID = ? AND OptionID = ?")
			if err != nil {
				return false, fmt.Errorf("could not prepare sql statement: %v", err)
			}
			_, err = stmt.Exec(a.PollID, a.UserID, a.OptionID)
			if err != nil {
				return false, fmt.Errorf("could not delete previous answer: %v", err)
			}

			// decrement previously selected option(s)
			stmt, err = tx.Prepare("UPDATE option SET Ctr = Ctr - 1 WHERE Ctr > 0 AND ID = ?")
			if err != nil {
				return false, fmt.Errorf("could not prepare sql statement: %v", err)
			}
			if _, err = stmt.Exec(a.OptionID); err != nil {
				return false, fmt.Errorf("could not decrement option: %v", err)
			}
			return true, nil
		}

		if !isMultipleChoice(p) {

			// decrement previously selected option(s)
			stmt, err = tx.Prepare("UPDATE option SET Ctr = Ctr - 1 WHERE ID = ? AND Ctr > 0")
			if err != nil {
				return false, fmt.Errorf("could not prepare sql statement: %v", err)
			}
			for _, o := range prevOpts {
				if _, err = stmt.Exec(o); err != nil {
					return false, fmt.Errorf("could not decrement option: %v", err)
				}
			}

			// update answer
			stmt, err = tx.Prepare("UPDATE answer SET OptionID = ?, LastSaved = ? WHERE UserID = ? AND PollID = ?")
			if err != nil {
				return false, fmt.Errorf("could not prepare sql statement: %v", err)
			}
			_, err = stmt.Exec(a.OptionID, time.Now().UTC().Unix(), a.UserID, a.PollID)
			if err != nil {
				return false, fmt.Errorf("could not update vote: %v", err)
			}
		} else {
			// new vote
			stmt, err = tx.Prepare("INSERT INTO answer(PollID, OptionID, UserID, LastSaved, CreatedAt) values(?, ?, ?, ?, ?)")
			if err != nil {
				return false, fmt.Errorf("could not prepare sql statement: %v", err)
			}
			_, err = stmt.Exec(a.PollID, a.OptionID, a.UserID, time.Now().UTC().Unix(), time.Now().UTC().Unix())
			if err != nil {
				return false, fmt.Errorf("could not insert answer: %v", err)
			}
		}
	} else {
		// new vote
		stmt, err = tx.Prepare("INSERT INTO answer(PollID, OptionID, UserID, LastSaved, CreatedAt) values(?, ?, ?, ?, ?)")
		if err != nil {
			return false, fmt.Errorf("could not prepare sql statement: %v", err)
		}
		_, err = stmt.Exec(a.PollID, a.OptionID, a.UserID, time.Now().UTC().Unix(), time.Now().UTC().Unix())
		if err != nil {
			return false, fmt.Errorf("could not insert answer: %v", err)
		}
	}

	// increment counter for option
	stmt, err = tx.Prepare("UPDATE option set Ctr = Ctr + 1 WHERE ID = ?")
	if err != nil {
		return false, fmt.Errorf("could not prepare sql statement: %v", err)
	}

	_, err = stmt.Exec(a.OptionID)
	if err != nil {
		return false, fmt.Errorf("could not increment option counter: %v", err)
	}

	return false, nil
}

func (st *sqlStore) AddMsgToPoll(pollid int, messageid int, chatid int64) error {
	tx, err := st.db.Begin()
	if err != nil {
		return fmt.Errorf("could not begin database transaction: %v", err)
	}
	defer func() {
		if err != nil {
			if err := tx.Rollback(); err != nil {
				log.Printf("could not rollback database change: %v", err)
			}
			return
		}
		err = tx.Commit()
	}()

	stmt, err := tx.Prepare("INSERT OR REPLACE INTO pollmsg(PollID, MessageID, ChatID) values(?, ?, ?)")
	if err != nil {
		return fmt.Errorf("could not build sql insert statement: %v", err)
	}
	defer close(stmt)

	_, err = stmt.Exec(pollid, messageid, chatid)
	if err != nil {
		return fmt.Errorf("could not add message to poll: %v", err)
	}

	return nil
}

func (st *sqlStore) AddInlineMsgToPoll(pollid int, inlinemessageid string) error {
	tx, err := st.db.Begin()
	if err != nil {
		return fmt.Errorf("could not begin database transaction: %v", err)
	}
	defer func() {
		if err != nil {
			if err := tx.Rollback(); err != nil {
				log.Printf("could not rollback database change: %v", err)
			}
			return
		}
		err = tx.Commit()
	}()

	// InlineMessageId is the primary key
	stmt, err := tx.Prepare("INSERT OR REPLACE INTO pollinlinemsg(PollID, InlineMessageID) values(?, ?)")
	if err != nil {
		return fmt.Errorf("could not build sql insert statement: %v", err)
	}
	defer close(stmt)

	_, err = stmt.Exec(pollid, inlinemessageid)
	if err != nil {
		return fmt.Errorf("could not add message to poll: %v", err)
	}

	return nil
}

func (st *sqlStore) RemoveInlineMsg(inlinemessageid string) error {
	tx, err := st.db.Begin()
	if err != nil {
		return fmt.Errorf("could not begin database transaction: %v", err)
	}
	defer func() {
		if err != nil {
			if err := tx.Rollback(); err != nil {
				log.Printf("could not rollback database change: %v", err)
			}
			return
		}
		err = tx.Commit()
	}()

	stmt, err := tx.Prepare("DELETE FROM pollinlinemsg WHERE InlineMessageID = ?")
	if err != nil {
		return fmt.Errorf("could not build sql insert statement: %v", err)
	}
	defer close(stmt)

	_, err = stmt.Exec(inlinemessageid)
	if err != nil {
		return fmt.Errorf("could not remove inline message: %v", err)
	}

	return nil

}

func (st *sqlStore) SaveOptions(options []option) error {
	// option gets passed by value as we only change id numbers
	// and do not append new elements this should be fine

	tx, err := st.db.Begin()
	if err != nil {
		return fmt.Errorf("could not begin database transaction: %v", err)
	}
	defer func() {
		if err != nil {
			if err := tx.Rollback(); err != nil {
				log.Printf("could not rollback database change: %v", err)
			}
			return
		}
		err = tx.Commit()
	}()
	stmt, err := tx.Prepare("INSERT OR REPLACE INTO option(PollID, Ctr, Text) values(?, ?, ?)")
	if err != nil {
		return fmt.Errorf("could not prepare insert sql statement for options: %v", err)
	}
	defer close(stmt)

	var res sql.Result
	var id64 int64
	for i := 0; i < len(options); i++ {
		res, err = stmt.Exec(options[i].PollID, options[i].Ctr, options[i].Text)
		if err != nil {
			return fmt.Errorf("could not insert option into sql database: %v", err)
		}
		id64, err = res.LastInsertId()
		if err != nil {
			return fmt.Errorf("could not get id of last insert: %v", err)
		}
		options[i].ID = int(id64)
	}
	return nil
}

func (st *sqlStore) SaveUser(u *tgbotapi.User) error {
	tx, err := st.db.Begin()
	if err != nil {
		return fmt.Errorf("could not begin database transaction: %v", err)
	}
	defer func() {
		if err != nil {
			if err := tx.Rollback(); err != nil {
				log.Printf("could not rollback database change: %v", err)
			}
			return
		}
		err = tx.Commit()
	}()

	stmt, err := tx.Prepare("SELECT count(1) FROM user WHERE ID = ?")
	if err != nil {
		return fmt.Errorf("could not prepare sql statement: %v", err)
	}
	defer close(stmt)

	var cnt int
	err = stmt.QueryRow(u.ID).Scan(&cnt)
	if err != nil {
		return fmt.Errorf("could not check if user '%s' exists: %v", u.UserName, err)
	}
	if cnt != 0 {
		stmt, err = tx.Prepare("UPDATE user SET FirstName = ?, LastName = ?, UserName = ?, LastSaved = ? WHERE ID = ?")
		if err != nil {
			return fmt.Errorf("could not prepare sql statement: %v", err)
		}
		_, err = stmt.Exec(u.FirstName, u.LastName, u.UserName, time.Now().UTC().Unix(), u.ID)
		if err != nil {
			return fmt.Errorf("could not update user entry: %v", err)
		}
		return nil
	}

	stmt, err = tx.Prepare("INSERT INTO user(ID, FirstName, LastName, UserName, LastSaved, CreatedAt) values(?, ?, ?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("could not prepare sql insert statement: %v", err)
	}
	defer close(stmt)

	_, err = stmt.Exec(u.ID, u.FirstName, u.LastName, u.UserName, time.Now().UTC().Unix(), time.Now().UTC().Unix())
	if err != nil {
		return fmt.Errorf("could not execute sql insert statement: %v", err)
	}

	return nil
}

func (st *sqlStore) SavePoll(p *poll) (id int, err error) {
	tx, err := st.db.Begin()
	if err != nil {
		return int(id), fmt.Errorf("could not begin database transaction: %v", err)
	}
	defer func() {
		if err != nil {
			if err := tx.Rollback(); err != nil {
				log.Printf("could not rollback database change: %v", err)
			}
			return
		}
		err = tx.Commit()
	}()

	if p.ID != 0 {
		var stmt *sql.Stmt
		stmt, err = tx.Prepare("UPDATE poll SET UserID = ?, Question = ?, Inactive = ?, Private = ?, Type = ?, LastSaved = ?, CreatedAt = ? WHERE ID = ?")
		if err != nil {
			return id, fmt.Errorf("could not prepare sql statement: %v", err)
		}
		defer close(stmt)
		_, err = stmt.Exec(p.UserID, p.Question, p.Inactive, p.Private, p.Type, time.Now().UTC().Unix(), time.Now().UTC().Unix(), p.ID)
		if err != nil {
			return id, fmt.Errorf("could not update user entry: %v", err)
		}
		return id, nil
	}

	stmt, err := tx.Prepare("INSERT INTO poll(UserID, Question, Inactive, Private, Type, LastSaved, CreatedAt) values(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return id, fmt.Errorf("could not prepare sql insert statement: %v", err)
	}
	defer close(stmt)

	var res sql.Result
	res, err = stmt.Exec(p.UserID, p.Question, p.Inactive, p.Private, p.Type, time.Now().UTC().Unix(), time.Now().UTC().Unix())
	if err != nil {
		return id, fmt.Errorf("could not execute sql insert statement: %v", err)
	}

	id64, err := res.LastInsertId()
	if err != nil {
		return id, fmt.Errorf("could not get id of last insert: %v", err)
	}
	id = int(id64)

	return id, nil
}

func contains(slice []int, n int) bool {
	for _, i := range slice {
		if i == n {
			return true
		}
	}
	return false
}
