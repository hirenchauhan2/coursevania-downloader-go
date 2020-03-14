package main

import (
	"fmt"

	"github.com/hirenchauhan2/download-coursevania-courses/pkg/coursevaniadownloader"
)

func main() {

	courseName := "The Complete 2020 Flutter Development Bootcamp with Dart"
	// courseName	:= "Complete React Developer in 2020 (w Redux, Hooks, GraphQL)"

	course, _ := coursevaniadownloader.GetCourseLinks("[coursevania.com] " + courseName)

	fmt.Println(course[:1])
}
