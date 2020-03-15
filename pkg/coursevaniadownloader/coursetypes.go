package coursevaniadownloader

// Course struct is used for storing the course related information
type Course struct {
	Name     string `json:"name"`
	MimeType string `json:"mimeType"`
}

// CourseFolder struct for storing the course related info
type CourseFolder struct {
	Files []File `json:"files"`
}

// File struct for storing the individual file info
type File struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	MimeType     string `json:"mimeType"`
	ModifiedTime string `json:"modifiedTime"`
}

// DownloadStatus states if the file is downloaded successfully or not
type DownloadStatus map[string]bool
