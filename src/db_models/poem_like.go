package db_models

import "time"

// Represents a like on a poem.
type PoemLike struct {
	Id        string    `db:"id"`
	UserId    string    `db:"user_id"`
	PoemId    string    `db:"poem_id"`
	CreatedOn time.Time `db:"created_on"`
}

func (t PoemLike) GetId() string { return t.Id }
