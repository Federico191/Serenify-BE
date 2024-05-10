package model

import "github.com/google/uuid"

type ArticleResponse struct {
    ID uuid.UUID `json:"id"`
    Title string `json:"title"`
    Content string `json:"content"`
    PhotoLink string `json:"photo_link"`
}