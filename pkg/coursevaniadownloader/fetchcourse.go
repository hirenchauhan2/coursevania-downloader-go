package coursevaniadownloader

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/hirenchauhan2/coursevania-downloader-go/internals/logger"

	"github.com/hirenchauhan2/coursevania-downloader-go/internals/utils"
)

const baseURL = "https://coursevania.courses.workers.dev/"

// DownloadCourse function is used for downloading the course that is passed
// as argument to this function
func DownloadCourse(courseName string) (bool, error) {
	courseURL := baseURL + courseName

	contextLogger := logger.WithFields(logger.Fields{"location": "DownloadCourse"})
	contextLogger.Infof("Inside the DownloadCourse() for %s", courseName)

	defer contextLogger.Infof("Exiting the DownloadCourse() for %s", courseName)

	client := http.Client{}

	resp, getErr := httpCall(&client, courseURL)

	if getErr != nil {
		contextLogger.Errorf("There was some error fetching the course details. %s", getErr)
		return false, getErr
	}

	defer resp.Body.Close()

	respBody, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		contextLogger.Errorf("Unable to read the response body. Error: %v", readErr)
		return false, readErr
	}

	courseInfo := Course{}

	// Convert the JSON object to CourseFolder struct
	jsonErr := json.Unmarshal(respBody, &courseInfo)

	if jsonErr != nil {
		contextLogger.Errorf("Unable to convert the JSON result to course type. %s", jsonErr)
		return false, jsonErr
	}
	courseName1 := courseInfo.Name
	contextLogger.Infof("Course Name: %s", courseName1)

	mimeType := courseInfo.MimeType

	// did we got the main course folder?
	if utils.IsCourseFolder(mimeType) {
		// Yes, now we need to create the folder in our system's pwd(current working directory!)
		// with same name as course name
		contextLogger.Infof("Creating the course directory")
		dirCreated, dirExists := utils.CreateDirectory(courseName1)
		if !dirCreated && !dirExists {
			contextLogger.Errorf("Unable to create the directory: %s", courseName1)
			return false, errors.New("Unknown error, could not find or create the directory" + courseName1)
		}
		contextLogger.Infof("Created the course directory")

		chErr := os.Chdir(courseName1)
		if chErr != nil {
			contextLogger.Errorf("Could not change the directory to %s, %v", courseName1, chErr)
		}

		courseURL = courseURL + "/"
		contextLogger.Infof("The course url: %s", courseURL)
		downloadStatus, dlErr := DownloadFolder(&client, courseURL)

		if dlErr != nil {
			contextLogger.Errorf("There was an error while downloading the course files. %v", dlErr)
			return false, dlErr
		}

		contextLogger.Infof("=========================================================")
		contextLogger.Infof("Course Downloader Status")
		contextLogger.Infof("Course Name: %s", courseName1)
		for file, status := range downloadStatus {
			contextLogger.Infof("File: %s, Downloaded: %v \n", file, status)
		}
		contextLogger.Infof("=========================================================")

		return true, nil
	}

	return false, errors.New("Unknown error, kindly check the course name properly")
}

func mergeDownloadStatuses(mainStatus DownloadStatus, subFolderStatus DownloadStatus) DownloadStatus {
	for file, status := range subFolderStatus {
		mainStatus[file] = status
	}
	return mainStatus
}
