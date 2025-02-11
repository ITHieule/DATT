package types

import "time"

type User struct {
    ID          int       `json:"id"`
    Name        string    `json:"name"`
    Email       string    `json:"email"`
    Password    string    `json:"-"`
    Phone       string    `json:"phone"`
    AvatarURL   string    `json:"avatar_url"`
    CreatedAt   time.Time `json:"created_at"`
}