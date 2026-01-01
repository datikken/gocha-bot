package main

import (
	"database/sql"
	"log"
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
	role := "customer"
	if username == "datikken" {
		role = "admin"
		// Update to admin
		_, err = db.Exec("UPDATE users SET role = 'admin' WHERE id = $1", id)
		if err != nil {
			log.Printf("Error updating user to admin: %v", err)
		}
	} else if id == 1 {
		role = "barber"
		// Update to barber
		_, err = db.Exec("UPDATE users SET role = 'barber' WHERE id = $1", id)
		if err != nil {
			log.Printf("Error updating first user to barber: %v", err)
		}
	}
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

func createVisit(userID int, visitDate time.Time) error {
	_, err := db.Exec("INSERT INTO visits (user_id, visit_date) VALUES ($1, $2)", userID, visitDate)
	return err
}

func getUpcomingVisits(userID int) ([]Visit, error) {
	rows, err := db.Query("SELECT id, user_id, visit_date, status, created_at FROM visits WHERE user_id = $1 AND visit_date > $2 AND status = 'scheduled' ORDER BY visit_date", userID, time.Now())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var visits []Visit
	for rows.Next() {
		var v Visit
		err := rows.Scan(&v.ID, &v.UserID, &v.VisitDate, &v.Status, &v.CreatedAt)
		if err != nil {
			return nil, err
		}
		visits = append(visits, v)
	}
	return visits, nil
}

func getAllUsers() ([]User, error) {
	rows, err := db.Query("SELECT id, telegram_id, username, role, invited_by, created_at FROM users")
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

func getCustomers() ([]User, error) {
	rows, err := db.Query("SELECT id, telegram_id, username, role, invited_by, created_at FROM users WHERE role = 'customer'")
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
