package coursevaniadownloader

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/hirenchauhan2/coursevania-downloader-go/internals/logger"
)

// httpCall a helper function to call the API to get the data
func httpCall(client *http.Client, url string) (*http.Response, error) {
	contextLogger := logger.WithFields(logger.Fields{"location": "httpCall"})

	contextLogger.Infof("inside the httpCall() for: %s", url)

	defer contextLogger.Infof("Exiting the httpCall() for: %s", url)

	requestBody, errReqBody := json.Marshal(map[string]string{
		"password": "null",
	})

	if errReqBody != nil {
		contextLogger.Errorf("There was some error while creating request body")
		return nil, errReqBody
	}

	request, errReq := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))

	if errReq != nil {
		contextLogger.Errorf("There was some error while creating request body")
		return nil, errReq
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// send the request to server
	resp, getErr := client.Do(request)

	if getErr != nil {
		contextLogger.Errorf("There was some error fetching the course details.")
		return nil, getErr
	}

	statusCode := resp.StatusCode

	if statusCode >= 500 {
		return nil, errors.New("There was an internal error at server")
	} else if statusCode == 404 {
		contextLogger.Errorf("404 - The course not found!")
		return nil, errors.New("404 - Not Found")
	}

	return resp, nil
}
