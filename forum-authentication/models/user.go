package models

import (
	"database/sql"
	"errors"
	"regexp"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

// CheckEmailExists verifies if an email is already registered in the database.
func CheckEmailExists(email string) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", email).Scan(&exists)
	return exists, err
}

// CheckUsernameExists verifies if an email is already registered in the database.
func CheckUsernameExists(username string) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", username).Scan(&exists)
	return exists, err
}

// RegisterUser creates a new user with the given email, username, and hashed password.
func RegisterUser(email, username, password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	userID, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	sessionToken, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	_, err = db.Exec("INSERT INTO users (id, email, username, password, session_token) VALUES (?, ?, ?, ?, ?)",
		userID.String(), email, username, hashedPassword, sessionToken.String())
	return sessionToken.String(), err
}

// AuthenticateUser checks the user's email and password, returning their ID if valid.
func AuthenticateUser(email, password string) (string, error) {
	var userID, hashedPassword string

	err := db.QueryRow("SELECT id, password FROM users WHERE email = ?", email).Scan(&userID, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("invalid credentials")
		}
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	sessionToken, _ := uuid.NewV4()
	_, err = db.Exec("UPDATE users SET session_token = ? WHERE id = ?", sessionToken.String(), userID)
	if err != nil {
		return "", err
	}

	return sessionToken.String(), nil
}

// GetUsernameByID retrieves a username based on the user ID.
func GetIDBySessionToken(sessionToken string) (string, string, error) {
	var username string
	var userID string
	err := db.QueryRow("SELECT id, username FROM users WHERE session_token = ?", sessionToken).Scan(&userID, &username)
	if err != nil {
		return "", "", err
	}
	return userID, username, nil
}

func GoogleGithubUser(id string, email string, provider string) (string, error) {
	var username string
	err := db.QueryRow("SELECT username FROM users WHERE id = ? AND oauth_provider = ?", id, provider).Scan(&username)
	if err == sql.ErrNoRows {
		re := regexp.MustCompile(`^[^@]+`)
		username = re.FindString(email)

		// User does not exist, create a new one
		_, err = db.Exec("INSERT INTO users (id, email, username, oauth_provider) VALUES (?, ?, ?, ?)", id, email, username, provider)
		if err != nil && provider == "google" {
			return "", errors.New("You already registered as GitHub user with this email. Please log in with GitHub account")
		}
		if err != nil && provider == "github" {
			return "", errors.New("You already registered as Google user with this email. Please log in with Google account")
		}
		if err != nil  {
			return "", err
		}
	} else if err != nil {
		return "", err
	}

	sessionToken, _ := uuid.NewV4()
	_, err = db.Exec("UPDATE users SET session_token = ? WHERE id = ?", sessionToken.String(), id)
	if err != nil {
		return "", err
	}

	return sessionToken.String(), nil
}
