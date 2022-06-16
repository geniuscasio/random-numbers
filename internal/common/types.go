package common

import (
	"github.com/google/uuid"
	"time"
)

type GenerateRequest struct {
	From  int64 `json:"from_number" binding:"required,gte=-1000000000,lte=1000000000"`
	To    int64 `json:"to_number" binding:"required,gte=-1000000000,lte=1000000000,gtfield=From"`
	Total int64 `json:"total_numbers" binding:"required,gte=1,lte=10000"`
}

type Order struct {
	Order string `form:"sort"`
}

type GenerateResponse []Row

type Row struct {
	Number int64 `json:"number"`
	Count  int64 `json:"count"`
}

type DetailsResponse struct {
	Name                string `json:"name"`
	CreatedAt           string `json:"date_created"`
	NumberOfGenerations int64  `json:"number_of_generations"`
}

type Session struct {
	ID string `json:"session_id" binding:"required"`
}

type UserCredentials struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"pwd" binding:"required"`
}

type User struct {
	ID uuid.UUID
	UserCredentials
	Name           string `json:"name" binding:"required"`
	PasswordHash   string
	CallsStatistic map[string]int64
	CreatedAt      time.Time
}
