package models

import "time"

const UserTableName = "users"

type User struct {
	Id            string    `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Password      string    `json:"password"`
	Verified      bool      `json:"verified"`
	Code          string    `json:"code"`
	SendExpiresAt time.Time `json:"send_expires_at"`
	Token         string    `json:"token"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
