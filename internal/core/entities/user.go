package entities

import "time"

const UserTableName = "users"

type User struct {
	Id            string
	Name          string
	Email         string
	Password      string
	Verified      bool
	Code          string
	SendExpiresAt time.Time
	Token         string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
