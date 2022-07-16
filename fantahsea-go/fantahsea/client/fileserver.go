package client

import "fantahsea/util"

// Download file from file-server
func DownloadFile(fileKey string, user *util.User, absPath string) error {
	return nil // todo
}

// Validate the file key, return true if it's valid else false
func ValidateFileKey(fileKey string, user *util.User) (bool, error) {
	return true, nil // todo
}
