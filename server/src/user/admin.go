package user

import "qoj/server/config"

func FindAdminByUsername(username string) (User, error) {
	var user User
	err := config.DB.
		QueryRow("SELECT RTRIM(username), password, RTRIM(fullname), profile_picture FROM users WHERE username = $1 AND is_admin = TRUE", username).
		Scan(&user.Username, &user.Password, &user.Fullname, &user.ProfilePicture)
	if err != nil {
		return User{}, err
	}
	return user, nil
}
