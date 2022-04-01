package db_models

import "time"

// Represents a poem.
type Poem struct {
	Id        string    `db:"id"`
	CreatedOn time.Time `db:"created_on"`
	UpdatedOn time.Time `db:"updated_on"`
	UserId    string    `db:"user_id"`
	Title     string    `db:"title"`
	Text      string    `db:"text"`
}

func (t Poem) GetId() string { return t.Id }
