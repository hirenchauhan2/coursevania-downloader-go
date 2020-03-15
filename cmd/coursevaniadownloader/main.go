package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/hirenchauhan2/coursevania-downloader-go/internals/logger"
	"github.com/hirenchauhan2/coursevania-downloader-go/pkg/coursevaniadownloader"
)

func init() {
	config := logger.Configuration{
		EnableConsole:     true,
		ConsoleLevel:      logger.Info,
		ConsoleJSONFormat: false,
		EnableFile:        true,
		FileLevel:         logger.Info,
		FileJSONFormat:    false,
		FileLocation:      "coursevaniadownloader.log",
	}
	err := logger.NewLogger(config, logger.InstanceZapLogger)
	if err != nil {
		log.Fatalf("Could not instantiate log %s", err.Error())
	}
}

func main() {
	contextLogger := logger.WithFields(logger.Fields{"location": "main"})

	courseName := flag.String("course", "", "Name of the course required. (should not include the string [coursevania.com] in name)")

	flag.Parse()

	// courseName := "The Complete 2020 Flutter Development Bootcamp with Dart"
	// courseName	:= "Complete React Developer in 2020 (w Redux, Hooks, GraphQL)"
	if str := strings.TrimSpace(*courseName); str == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	contextLogger.Infof("Downloading the course: %s", *courseName)

	_, err := coursevaniadownloader.DownloadCourse("[coursevania.com] " + *courseName)

	if err != nil {
		contextLogger.Errorf("Error while downloading the course: %v", err)
	}
}
