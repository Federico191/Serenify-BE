package model

import "github.com/google/uuid"

type SeminarResponse struct {
    ID uuid.UUID `json:"id"`
    Title string `json:"title"`
    Place string `json:"place"`
    Time string `json:"time"`
    Price int `json:"price"`
    PhotoLink string `json:"photo_link"`
}

type SeminarDetailResponse struct {
    ID uuid.UUID `json:"id"`
    Title string `json:"title"`
    Place string `json:"place"`
    Time string `json:"time"`
    Price int `json:"price"`
    PhotoLink string `json:"photo_link"`
    Description string `json:"description"`
}