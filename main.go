package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Model for course
type Course struct {
	CourseId    string  `json:"course"`
	CourseName  string  `json:"coursename"`
	CoursePrice int     `json:"price"`
	Author      *Author `json:"author"` //this will have a custom type which we already declared below so we will define the type as pointers such that the data will not be use as a copy
}

// 2nd model
type Author struct {
	FullName string `json:"fullname"`
	Website  string `json:"website"`
}

// fake database
var courses []Course

// MiddleWare or Helper file
// this will help to accomplish certain tasks
func (c *Course) isEmpty() bool {
	return c.CourseName == "" //we didn't mention CourseId because we will generate that
}

// Controllers
// serve home route - we are creating it so that our page will not look empty and have some message
// and also when the request will come from any method than that will be governed by this function

func serveHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Welcome to Golang</h1>"))
}

func getAllCourses(w http.ResponseWriter, r *http.Request) { //this is created such that we can send all our database into it
	fmt.Println("Get all courses")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses) //this is how we throw all of the thing in our fake DB to the json
}

// Want to get one course based on the id requested

func getOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get one course")
	w.Header().Set("Content-Type", "application/json")

	//grab the id of course which we are getting from request
	params := mux.Vars(r)

	//1. loop through courses
	//2. find the matching id
	//3. and return the response

	for _, course := range courses {
		if course.CourseId == params["id"] { // params contains key value pair
			json.NewEncoder(w).Encode(course)
			return
		}
	}
	json.NewEncoder(w).Encode("No courses found with given id")

}

func createOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create one course")
	w.Header().Set("Content-Type", "application/json")

	// what if - body is empty
	if r.Body == nil {
		json.NewEncoder(w).Encode("Please send some data")
	}

	// what if - body is {}
	var course Course
	_ = json.NewDecoder(r.Body).Decode(&course) // here we are getting the json so we will decode it and dont want any return value that's why using _
	if course.isEmpty() {
		json.NewEncoder(w).Encode("No data inside JSON")
		return
	}
	//generate unique id , string
	// append course into courses

	rand.Seed(time.Now().UnixNano())
	course.CourseId = strconv.Itoa(rand.Intn(100))
	courses = append(courses, course)
	json.NewEncoder(w).Encode(course)
	return

}

func updateOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update one course")
	w.Header().Set("Content-Type", "application/json")

	//step-1: Grab the id from request
	params := mux.Vars(r)

	//step-2: loop over the DB
	//step-3: Match the id
	//step-4: Remove the data of that index
	//step-5: now add the data with id

	for index, course := range courses {
		if course.CourseId == params["id"] {
			courses = append(courses[:index], courses[index+1:]...)
			var course Course // Here we are creating an object of type struct so that we can pass its reference when we decode the JSON
			json.NewDecoder(r.Body).Decode(&course)
			course.CourseId = params["id"] //Overwriting the ID
			courses = append(courses, course)
			json.NewEncoder(w).Encode(course)
			return
		}
	}
}

func deleteOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete one course")
	w.Header().Set("Content-Type", "application/json")

	//grab id
	params := mux.Vars(r)

	// loop, find id, remove the data with id

	for index, course := range courses {
		if course.CourseId == params["id"] {
			courses = append(courses[:index], courses[index+1:]...)
			json.NewEncoder(w).Encode("Deleting is Successful")
			break
		}
	}
}

func main() {

	fmt.Println("Building API")
	r := mux.NewRouter()

	//seeding

	courses = append(courses, Course{CourseId: "1", CourseName: "ReactJS", CoursePrice: 299, Author: &Author{FullName: "Nishchal Gupta", Website: "lco.dev"}})
	courses = append(courses, Course{CourseId: "2", CourseName: "MERN Stack", CoursePrice: 599, Author: &Author{FullName: "Nishchal Gupta", Website: "go.dev"}})

	//routing
	r.HandleFunc("/", serveHome).Methods("GET")
	r.HandleFunc("/courses", getAllCourses).Methods("GET")
	r.HandleFunc("/course/{id}", getOneCourse).Methods("GET") // Here {x} should be same as we use in params["x"]
	r.HandleFunc("/course", createOneCourse).Methods("POST")
	r.HandleFunc("/course/{id}", updateOneCourse).Methods("PUT") // Here PUT will bring both id and data
	r.HandleFunc("/course/{id}", deleteOneCourse).Methods("DELETE")

	//listen to port

	log.Fatal(http.ListenAndServe(":8000", r))
}
