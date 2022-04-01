package db_models

import "time"

// Represents a connection from one user to another.
type UserFollowing struct {
	Id          string    `db:"id"`
	FollowerId  string    `db:"follower_id"`
	FollowingId string    `db:"following_id"`
	CreatedOn   time.Time `db:"created_on"`
}

func (t UserFollowing) GetId() string { return t.Id }
