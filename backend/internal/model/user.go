package model

import (
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func initDefaultData() error {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		_, err = DB.Exec(
			"INSERT INTO users (username, password, role) VALUES (?, ?, ?)",
			"admin", string(hashedPassword), "admin",
		)
		if err != nil {
			return err
		}
		log.Println("Default admin user created: admin/admin123")
	}

	return nil
}

func GetUserByUsername(username string) (*User, error) {
	user := &User{}
	err := DB.QueryRow(
		"SELECT id, username, password, role, created_at, updated_at FROM users WHERE username = ?",
		username,
	).Scan(&user.ID, &user.Username, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserByID(id int) (*User, error) {
	user := &User{}
	err := DB.QueryRow(
		"SELECT id, username, password, role, created_at, updated_at FROM users WHERE id = ?",
		id,
	).Scan(&user.ID, &user.Username, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetAllUsers() ([]User, error) {
	rows, err := DB.Query("SELECT id, username, role, created_at, updated_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Username, &user.Role, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func CreateUser(username, password, role string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = DB.Exec(
		"INSERT INTO users (username, password, role) VALUES (?, ?, ?)",
		username, string(hashedPassword), role,
	)
	return err
}

func UpdateUserPassword(id int, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = DB.Exec(
		"UPDATE users SET password = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		string(hashedPassword), id,
	)
	return err
}

func DeleteUser(id int) error {
	_, err := DB.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}

func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
