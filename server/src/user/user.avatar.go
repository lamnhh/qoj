package user

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"qoj/server/config"
	"sync"
)

var locks map[string]*sync.Mutex
var locksMutex sync.Mutex

func updateProfilePicture(username string, header *multipart.FileHeader) (string, error) {
	locksMutex.Lock()
	if locks[username] == nil {
		locks[username] = &sync.Mutex{}
	}
	locksMutex.Unlock()

	locks[username].Lock()
	defer locks[username].Unlock()

	uploadedFile, err := header.Open()
	if err != nil {
		return "", err
	}

	targetPath := filepath.Join(".", "server", "profile-picture", username)
	targetFile, err := os.Create(targetPath)
	if err != nil {
		return "", err
	}

	// Copy uploadedFile to targetFile
	if _, err = io.Copy(targetFile, uploadedFile); err != nil {
		return "", err
	}

	// Update profile_picture of `username` to '/profile-picture/username'
	_, err = config.DB.Exec(
		"UPDATE users SET profile_picture = $1 WHERE username = $2",
		"/profile-picture/" + username,
		username,
	)
	if err != nil {
		return "", err
	}
	return "/profile-picture/" + username, err
}

func InitialiseAvatarLocks() {
	locks = make(map[string]*sync.Mutex)
}