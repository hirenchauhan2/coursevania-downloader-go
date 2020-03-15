package coursevaniadownloader

import (
	"errors"
	"net/http"
	"os"

	"github.com/hirenchauhan2/coursevania-downloader-go/internals/logger"
	"github.com/hirenchauhan2/coursevania-downloader-go/internals/utils"
)

// DownloadFolder will download the course folder and
// return the Course with links
func DownloadFolder(client *http.Client, folderURL string) (DownloadStatus, error) {
	var downloadStatus DownloadStatus = make(DownloadStatus, 0)

	contextLogger := logger.WithFields(logger.Fields{"location": "DownloadFolder"})

	contextLogger.Infof("Inside DownloadFolder() to download: %s", folderURL)

	contextLogger.Infof(" DownloadFolder() calling getCourseFolder: %s", folderURL)

	// last statement before exiting the function
	defer contextLogger.Infof("Exiting the DownloadFolder() for: %s", folderURL)

	courseFolder, errFolder := getCourseFolder(client, folderURL)
	if errFolder != nil {
		contextLogger.Errorf("Unable to download the course folder: %s", folderURL)
		downloadStatus[folderURL] = false
		return downloadStatus, errFolder
	}

	contextLogger.Infof(" DownloadFolder() successfully got the result from getCourseFolder()")

	// get the count of files inside a folder
	fileCount := len(courseFolder.Files)
	contextLogger.Infof("Folders/Files count: %d", fileCount)

	for _, file := range courseFolder.Files {
		fileName := file.Name
		link := folderURL + fileName
		contextLogger.Infof("Current Folder/file: %s", fileName)

		// if the file is actually an folder
		if utils.IsCourseFolder(file.MimeType) {
			folderLink := link + "/"
			contextLogger.Infof("Folder Link: %s", folderLink)
			// Create directory for this folder
			contextLogger.Infof("Creating the directory for folder: %s", file.Name)
			dirCreated, dirExists := utils.CreateDirectory(fileName)
			if !dirCreated && !dirExists {
				contextLogger.Errorf("Error in creating the directory")
				downloadStatus[folderLink] = false
				return downloadStatus, errors.New("Error in creating the directory")
			}
			contextLogger.Infof("Created the directory for folder, changing the directory: %s", file.Name)

			chErr := os.Chdir(fileName)

			if chErr != nil {
				contextLogger.Errorf("There was an error while changing the directory. Error %v", chErr)
				downloadStatus[folderLink] = false
				return downloadStatus, chErr
			}

			subFolderDownloaded, dlErr := DownloadFolder(client, folderLink)

			if dlErr != nil {
				contextLogger.Errorf("Error while downloading the folder: %s, %v ", fileName, dlErr)
			}

			downloadStatus = mergeDownloadStatuses(downloadStatus, subFolderDownloaded)

			chErr = os.Chdir("..")
			if chErr != nil {
				contextLogger.Errorf("There was an error while changing the directory %v", chErr)
			}
		} else {
			contextLogger.Infof("File Link: %s\n", link)

			fileDownloaded, dlErr := downloadFile(client, link)

			if dlErr != nil {
				contextLogger.Errorf("Downloading failed for file %s", link)
				downloadStatus[link] = false
			}
			downloadStatus[link] = fileDownloaded

			contextLogger.Infof("File downloaded successfully")
		}
	}
	return downloadStatus, nil
}
