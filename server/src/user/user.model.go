package user

import "qoj/server/config"

type UserLogin struct {
	Username string `json:"username" binding:"required"`
	Fullname string `json:"fullname" binding:"required"`
}

type User struct {
	UserLogin
	Password string `json:"password" binding:"required"`
}

func FindUserByUsername(username string) (User, error){
	var user User
	err := config.DB.
		QueryRow("SELECT RTRIM(username), password, RTRIM(fullname) FROM users WHERE username = $1", username).
		Scan(&user.Username, &user.Password, &user.Fullname)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func CreateNewUser(user User) error {
	hashedPassword := HashPassword(user.Password)
	_, err := config.DB.Query("SELECT create_user($1, $2, $3)", user.Username, hashedPassword, user.Fullname)
	return err
}