package db

import (
	"Go-server/model"
	"database/sql"
	"fmt"
	"time"
)

type PstgDB struct {
	pdb *sql.DB
}

func NewPstgrDB(cfg *model.Config) (*PstgDB, error) {
	conn, err := sql.Open("postgres", cfg.ConnectString)
	if err != nil {
		return nil, err
	}

	pstgDB := &PstgDB{pdb: conn}
	if err := conn.Ping(); err != nil {
		return nil, err
	}

	fmt.Println("Database connected")
	return pstgDB, nil
}

func (o *PstgDB) CreateAccount(login string, psw string) error {
	_, err := o.pdb.Exec(
		`INSERT INTO account(username, password) 
			VALUES($1, crypt($2, gen_salt('bf',8)));`, login, psw)

	if err != nil {
		return fmt.Errorf(
			"Couldn't insert value: %s, %s in account", login, psw)
	}

	return nil
}

func (o *PstgDB) SelectAccount(login string, psw string) (*model.Account, error) {
	var account model.Account
	err := o.pdb.QueryRow(
		`SELECT user_id, username, password 
		FROM account 
		WHERE 
			username = $1 AND 
			password = crypt($2, password);`, login, psw).Scan(
		&account.Id, &account.Login, &account.Pass)

	if err != nil {
		return nil, fmt.Errorf("Couldn't select account: %v", err.Error())
	}

	return &account, nil
}

func (o *PstgDB) SelectAccountLogin(login string) (*model.Account, error) {
	var account model.Account
	err := o.pdb.QueryRow(
		`SELECT user_id, username, password 
		FROM account 
		WHERE 
			username = $1 ;`, login).Scan(
		&account.Id, &account.Login, &account.Pass)

	if err != nil {
		return nil, fmt.Errorf("Couldn't select account: %v", err.Error())
	}

	return &account, nil
}

func (o *PstgDB) InsertSession(sid string, uid string, session_start time.Time) error {
	_, err := o.pdb.Exec(
		`INSERT INTO session VALUES($1, $2, $3, NULL);`, sid, uid, session_start)
	if err != nil {
		return fmt.Errorf(
			"Couldn't insert session. SID: %s, Error: %v", sid, err.Error())
	}

	return nil
}

func (o *PstgDB) SelectSessionSID(sid string) (*model.Session, error) {
	session := model.Session{}
	err := o.pdb.QueryRow(
		`SELECT session_id, user_id, session_start, session_end
		FROM session 
		WHERE session_id = $1;`, sid).Scan(
		&session.Id, &session.Uid, &session.Start, &session.End)

	if err != nil {
		return nil, fmt.Errorf(
			"Couldn't select session. SID: %s, Error: %v", sid, err.Error())
	}

	return &session, nil
}

func (o *PstgDB) DeleteSessionSID(sid string) error {
	_, err := o.pdb.Exec(
		`DELETE FROM session
		WHERE session_id = $1;`, sid)
	if err != nil {
		return fmt.Errorf(
			"Couldn't delete session. SID: %s; Error: %v ", sid, err.Error())
	}

	return nil
}
