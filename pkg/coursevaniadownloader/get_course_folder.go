package coursevaniadownloader

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/hirenchauhan2/coursevania-downloader-go/internals/logger"
)

func getCourseFolder(client *http.Client, folderURL string) (CourseFolder, error) {
	var courseFolder = CourseFolder{}
	contextLogger := logger.WithFields(logger.Fields{"location": "getCourseFolder"})
	contextLogger.Infof("Inside the getCourseFolder() for: %s", folderURL)

	// last statement before exiting the function
	defer contextLogger.Infof("Exiting the getCourseFolder() for: %s", folderURL)

	// call the api to get course folder
	resp, getErr := httpCall(client, folderURL)

	if getErr != nil {
		contextLogger.Errorf("There was some error fetching the course details", getErr)
		return courseFolder, getErr
	}

	defer resp.Body.Close()

	respBody, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		contextLogger.Errorf("Unable to read the response body", readErr)
		return courseFolder, readErr
	}

	// Convert the JSON object to CourseFolder struct
	jsonErr := json.Unmarshal(respBody, &courseFolder)

	if jsonErr != nil {
		contextLogger.Errorf("Unable to unmarshalling JSON result to course type", jsonErr)
		return courseFolder, jsonErr
	}

	return courseFolder, nil
}
