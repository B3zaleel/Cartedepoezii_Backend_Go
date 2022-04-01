package db_models

import "time"

// Represents a comment on a poem or reply to a comment.
type Comment struct {
	Id        string    `db:"id"`
	UserId    string    `db:"user_id"`
	PoemId    string    `db:"poem_id"`
	CommentId string    `db:"comment_id"`
	Text      string    `db:"text"`
	CreatedOn time.Time `db:"created_on"`
}

func (t Comment) GetId() string { return t.Id }
