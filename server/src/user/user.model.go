package user

import (
	"qoj/server/config"
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
	hashedPassword := hashPassword(user.Password)
	_, err := config.DB.Exec("SELECT create_user($1, $2, $3)", user.Username, hashedPassword, user.Fullname)
	return err
}