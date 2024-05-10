package model

type AnswerRequest struct {
    Answer []string `json:"answer" binding:"required,min=20"`
}

type AnswerResponse struct {
    Score int `json:"score"`
    Description string `json:"description"`
}