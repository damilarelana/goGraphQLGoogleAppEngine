package models

import "time"

// User fields declared
type User struct {
	ID   string `json:"id" datastore:"-"`
	Name string `json:"name"`
}

// Post fields declared
type Post struct {
	ID        string    `json:"id" datastore:"-"`
	UserID    string    `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
	Content   string    `json:"content"`
}
