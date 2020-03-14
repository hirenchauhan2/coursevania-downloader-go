package utils

// IsFolder function checks if the given mime type is
// folder or not
func IsFolder(mimeType string) bool {
	return "application/vnd.google-apps.folder" == mimeType
}
