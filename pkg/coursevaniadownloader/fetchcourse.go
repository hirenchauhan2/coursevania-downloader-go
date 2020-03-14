package coursevaniadownloader

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/hirenchauhan2/download-coursevania-courses/internals/utils"
)

var baseURL string = "https://coursevania.courses.workers.dev/"

// DownloadCourse function is used for downloading the course that is passed
// as argument to this function
func DownloadCourse(courseName string) error {
	courseURL := baseURL + courseName

	requestBody, errReqBody := json.Marshal(map[string]string{
		"password": "null",
	})
	if errReqBody != nil {
		log.Fatalln("There was some error while creating request body")
	}

	// 2o seconds of timeout for a request
	timeout := time.Duration(20 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	request, errReq := http.NewRequest("POST", courseURL, bytes.NewBuffer(requestBody))
	if errReq != nil {
		log.Fatalln("There was some error while creating request body")
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, getErr := client.Do(request)

	if getErr != nil {
		log.Fatalln("There was some error fetching the course details.")
	}

	if resp.StatusCode == 404 {
		log.Fatalln("Oops course not found. Please enter correct course name")
	}

	defer resp.Body.Close()

	respBody, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatalln("Unable to read the response body")
	}

	courseInfo := Course{}

	// Convert the JSON object to CourseFolder struct
	jsonErr := json.Unmarshal(respBody, &courseInfo)

	if jsonErr != nil {
		log.Fatalln("Unable to unmarshalling JSON result to course type")
	}
	courseName1 := courseInfo.Name
	fmt.Println("Course Name: ", courseName1)

	mimeType := courseInfo.MimeType
	// application/vnd.google-apps.folder
	fmt.Println("Mime Type: ", mimeType)

	// did we got the main course folder?
	if utils.IsFolder(mimeType) {
		fmt.Println("We have the Course Folder")

		// Yes, now we need to create the folder in our system's pwd(current working directory!)
		// with same name as course name
		// get the current working directory
		pwd, wdErr := os.Getwd()
		if wdErr != nil {
			log.Println("Could not get the working directory.", wdErr)
		}
		fmt.Println("pwd: ", pwd)
		chErr := os.Chdir(pwd)
		if chErr != nil {
			log.Fatalln("Could not change directory", chErr)
		}
		// TODO: create directory in the working directory
		mkErr := os.Mkdir(courseName1, 0755)
		if mkErr != nil {
			log.Fatalln("Unable to create the directory", mkErr)
		}
		chErr = os.Chdir(courseName1)
		if chErr != nil {
			log.Fatalln("Could not change directory", chErr)
		}
		// TODO: need to call get links recursively for getting the folder and sub folder's details(Struct needs to be created)

	}

	return nil
}

// GetCourseLinks will return the Course with links
func GetCourseLinks(courseName string) ([]string, error) {
	//baseURL := "https://coursevania.courses.workers.dev/"
	courseURL := baseURL + courseName + "/"
	var links []string

	e := DownloadCourse(courseName)

	if e != nil {
		log.Fatalln("Course Not Found or internal server error")
	}

	fmt.Println("Course Link:", courseURL)

	timeout := time.Duration(20 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	requestBody, errReqBody := json.Marshal(map[string]string{
		"password": "null",
	})

	if errReqBody != nil {
		log.Fatalln("There was some error while creating request body")
	}

	request, errReq := http.NewRequest("POST", courseURL, bytes.NewBuffer(requestBody))
	if errReq != nil {
		log.Fatalln("There was some error while creating request body")
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, getErr := client.Do(request)

	if getErr != nil {
		log.Fatalln("There was some error fetching the course details.")
	}

	statusCode := resp.StatusCode
	fmt.Println(statusCode)

	if statusCode >= 500 {
		log.Fatalf("There was an internal error at server. Could not process further.")
	}

	fmt.Println(resp.Header.Get("Content-Length"))

	defer resp.Body.Close()

	respBody, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatalln("Unable to read the response body")
	}

	courseFolder := CourseFolder{}

	// Convert the JSON object to CourseFolder struct
	jsonErr := json.Unmarshal(respBody, &courseFolder)

	if jsonErr != nil {
		log.Fatalln("Unable to unmarshalling JSON result to course type")
	}

	// create list of folders with the files inside
	folders := make(map[string]CourseFolder)

	// create links for each folder and visit each folder
	folderCount := len(courseFolder.Files)

	log.Println("Folders count: ", folderCount)

	for i := 0; i < folderCount; i++ {
		link := courseURL + courseFolder.Files[i].Name
		link1 := link + "/"

		// log.Println("Visiting: ", link1)

		subRequest, subReqErr := http.NewRequest("POST", link1, bytes.NewBuffer(requestBody))
		if subReqErr != nil {
			log.Fatalln("There was some error while creating request body")
		}
		subRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		subReqResp, subReqGetErr := client.Do(subRequest)

		if subReqGetErr != nil {
			log.Fatalln("There was some error fetching the sob folder details.", subReqErr)
		}

		defer subReqResp.Body.Close()

		subReqRespBody, readErr := ioutil.ReadAll(subReqResp.Body)
		if readErr != nil {
			log.Fatalln("Unable to read the response body")
		}

		subFolder := CourseFolder{}

		// Convert the JSON object to CourseFolder struct
		jsonErr = json.Unmarshal(subReqRespBody, &subFolder)

		if jsonErr != nil {
			log.Fatalln("Unable to unmarshalling JSON result to course type")
		}

		log.Println("Files in sub folder:", len(subFolder.Files))

		folders[link] = subFolder
	}

	// fmt.Println(folders)

	// Create a final list of all course files in one list
	for folder, subFolder := range folders {
		for _, file := range subFolder.Files {
			fileLink := folder + "/" + file.Name
			links = append(links, fileLink)
		}
	}

	return links, nil
}

// func httpCall (client *http.Client, url string)
