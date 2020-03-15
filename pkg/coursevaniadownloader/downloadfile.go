package coursevaniadownloader

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/hirenchauhan2/coursevania-downloader-go/internals/logger"
)

func downloadFile(client *http.Client, fileURL string) (bool, error) {
	contextLogger := logger.WithFields(logger.Fields{"location": "downloadFile"})

	contextLogger.Infof("Inside the downloadFile() for %s", fileURL)
	// last statement before exiting the function
	defer contextLogger.Infof("Exiting the downloadFile() for: %s", fileURL)

	resp, err := client.Get(fileURL)

	if err != nil {
		contextLogger.Errorf("Request to download the file failed. %v", err)
		return false, err
	}

	defer resp.Body.Close()

	workingDir, wdErr := os.Getwd()

	if wdErr != nil {
		return false, wdErr
	}

	fileName := filepath.Join(workingDir, filepath.Base(fileURL))

	contextLogger.Infof("Creating File: %s", fileName)

	file, fErr := os.Create(fileName)

	if fErr != nil {
		contextLogger.Errorf("Unable to create the file: %v", fErr)
		return false, fErr
	}

	contextLogger.Infof("Created File: %s", fileName)

	defer file.Close()

	contextLogger.Infof("Copying the response data to %s", fileName)
	_, err = io.Copy(file, resp.Body)

	if err != nil {
		contextLogger.Errorf("Unable to copy the downloaded data into file: %v", err)
		return false, err
	}

	contextLogger.Infof("Successfully copied the response data to %s", fileName)

	return true, nil
}
