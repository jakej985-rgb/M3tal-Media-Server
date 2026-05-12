package auth

import (
	"encoding/json"
	"os"

	"golang.org/x/crypto/bcrypt"
)

// User represents a dashboard user
type User struct {
	Username string `json:"username"`
	Hash     string `json:"token_hash"`
	Role     string `json:"role"`
}

// UpdateUser updates or creates a user in the specified JSON file
func UpdateUser(filePath, username, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	var users []User
	data, err := os.ReadFile(filePath)
	if err == nil {
		json.Unmarshal(data, &users)
	}

	found := false
	for i, u := range users {
		if u.Username == username {
			users[i].Hash = string(hash)
			found = true
			break
		}
	}
	if !found {
		users = append(users, User{Username: username, Hash: string(hash), Role: "admin"})
	}

	newData, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, newData, 0644)
}
