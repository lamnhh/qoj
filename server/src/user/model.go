package user

import (
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"qoj/server/config"
	"strings"
)

type LoginAuth struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password,omitempty" binding:"required"`
}

type User struct {
	LoginAuth
	Fullname       string `json:"fullname" binding:"required"`
	ProfilePicture string `json:"profilePicture"`
}

type PatchUser struct {
	Fullname string `json:"fullname"`
}

type PutPasswordUser struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

func Login(username string, password string) (User, int, error) {
	user, err := FindUserByUsername(username)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, http.StatusNotFound, errors.New(fmt.Sprintf("User `%s` does not exist", username))
		}
		return User{}, http.StatusInternalServerError, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return User{}, http.StatusBadRequest, errors.New("Wrong password")
	}
	user.Password = ""
	return user, http.StatusOK, nil
}

func FindUserByUsername(username string) (User, error) {
	var user User
	err := config.DB.
		QueryRow("SELECT RTRIM(username), password, RTRIM(fullname), profile_picture FROM users WHERE username = $1", username).
		Scan(&user.Username, &user.Password, &user.Fullname, &user.ProfilePicture)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func CreateNewUser(user User) error {
	hashedPassword := HashPassword(user.Password)
	_, err := config.DB.Exec("SELECT create_user($1, $2, $3)", user.Username, hashedPassword, user.Fullname)
	return err
}

func UpdateUserProfile(username string, modifier map[string]interface{}) (User, error) {
	keyList := make([]string, 0)
	valList := []interface{}{username}
	count := 1
	for k, v := range modifier {
		count++
		keyList = append(keyList, fmt.Sprintf("%s = $%d", k, count))
		valList = append(valList, v)
	}

	// No modifier given, return current profile
	if len(keyList) == 0 {
		return FindUserByUsername(username)
	}

	var user User
	cmd := fmt.Sprintf(
		"UPDATE users SET %s WHERE username = $1 RETURNING RTRIM(username), password, RTRIM(fullname), profile_picture",
		strings.Join(keyList, ", "),
	)

	err := config.DB.
		QueryRow(cmd, valList...).
		Scan(&user.Username, &user.Password, &user.Fullname, &user.ProfilePicture)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func updatePassword(username string, newPassword string) error {
	hashedPassword := HashPassword(newPassword)
	_, err := config.DB.Exec("UPDATE users SET password = $1 WHERE username = $2", hashedPassword, username)
	return err
}