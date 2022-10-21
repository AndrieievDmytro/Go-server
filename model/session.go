package model

import "database/sql"

type Session struct {
	Id    string
	Uid   int
	Start string
	End   sql.NullString
}
