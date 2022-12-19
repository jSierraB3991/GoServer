package models

import (
    "time"
)

type Post struct {
    Id          string      `json:"id"`
    Content     string      `json:"content"`
    CreateAt    time.Time   `json:"create_at"`
    UserId      string      `json:"user_id"`
}
