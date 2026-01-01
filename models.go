package main

import "time"

type User struct {
	ID         int       `json:"id"`
	TelegramID int64     `json:"telegram_id"`
	Username   string    `json:"username"`
	Role       string    `json:"role"` // "customer", "barber", or "admin"
	InvitedBy  int       `json:"invited_by"`
	CreatedAt  time.Time `json:"created_at"`
}

type Visit struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	VisitDate time.Time `json:"visit_date"`
	Status    string    `json:"status"` // "scheduled", "completed", "cancelled"
	CreatedAt time.Time `json:"created_at"`
}
