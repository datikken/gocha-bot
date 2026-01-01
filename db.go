package main

import (
	"database/sql"
	"time"
)

func getOrCreateUser(telegramID int64, username string) (*User, error) {
	var user User
	err := db.QueryRow("SELECT id, telegram_id, username, role, invited_by, created_at FROM users WHERE telegram_id = $1", telegramID).Scan(&user.ID, &user.TelegramID, &user.Username, &user.Role, &user.InvitedBy, &user.CreatedAt)
	if err != sql.ErrNoRows {
		return &user, nil
	}

	// Create new user
	result, err := db.Exec("INSERT INTO users (telegram_id, username) VALUES ($1, $2)", telegramID, username)
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	role := "client"
	user = User{
		ID:         int(id),
		TelegramID: telegramID,
		Username:   username,
		Role:       role,
		CreatedAt:  time.Now(),
	}
	return &user, nil
}

func updateUserRole(telegramID int64, role string) error {
	_, err := db.Exec("UPDATE users SET role = $1 WHERE telegram_id = $2", role, telegramID)
	return err
}

func createInvitedUser(username string, invitedBy int) error {
	_, err := db.Exec("INSERT INTO users (telegram_id, username, role, invited_by) VALUES ($1, $2, $3, $4)", 0, username, "client", invitedBy)
	return err
}

func getMyClients(invitedBy int) ([]User, error) {
	rows, err := db.Query("SELECT id, telegram_id, username, role, invited_by, created_at FROM users WHERE role = 'client' AND invited_by = $1", invitedBy)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.TelegramID, &u.Username, &u.Role, &u.InvitedBy, &u.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func getAdmins() ([]User, error) {
	rows, err := db.Query("SELECT id, telegram_id, username, role, invited_by, created_at FROM users WHERE role = 'admin'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.TelegramID, &u.Username, &u.Role, &u.InvitedBy, &u.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func getLastVisit(userID int) (*Visit, error) {
	var v Visit
	err := db.QueryRow("SELECT id, user_id, visit_date, status, created_at FROM visits WHERE user_id = $1 AND status = 'scheduled' ORDER BY visit_date DESC LIMIT 1", userID).Scan(&v.ID, &v.UserID, &v.VisitDate, &v.Status, &v.CreatedAt)
	if err != sql.ErrNoRows {
		return &v, nil
	}
	return nil, nil
}

func createVisit(userID int, visitDate time.Time) error {
	_, err := db.Exec("INSERT INTO visits (user_id, visit_date) VALUES ($1, $2)", userID, visitDate)
	return err
}

func getUserByUsername(username string) (*User, error) {
	var user User
	err := db.QueryRow("SELECT id, telegram_id, username, role, invited_by, created_at FROM users WHERE username = $1", username).Scan(&user.ID, &user.TelegramID, &user.Username, &user.Role, &user.InvitedBy, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
