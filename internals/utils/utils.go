package utils

import (
	"os"

	"github.com/hirenchauhan2/coursevania-downloader-go/internals/logger"
)

// IsCourseFolder function checks if the given mime type is
// folder or not
func IsCourseFolder(mimeType string) bool {
	return "application/vnd.google-apps.folder" == mimeType
}

// CreateDirectory creates a directory if not exists.
func CreateDirectory(dirName string) (bool, bool) {
	src, err := os.Stat(dirName)
	contextLogger := logger.WithFields(logger.Fields{"location": "DownloCreateDirectoryadFolder"})

	if os.IsNotExist(err) {
		errDir := os.MkdirAll(dirName, 0755)
		if errDir != nil {
			contextLogger.Panicf("%v", err)
			return false, false
		}
		// created, exists
		return true, false
	}

	if src.Mode().IsRegular() {
		contextLogger.Errorf("%s already exist as a file!", dirName)
		// created, exists
		return false, true
	}

	return false, false
}
