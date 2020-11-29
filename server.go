package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/rpc"
	"strconv"
	"strings"
)

const ADD_STUDENT = "1"
const STUDENT_AVERAGE = "2"
const GENERAL_AVERAGE = "3"
const SUBJECT_AVERAGE = "4"

var client *rpc.Client

func loadHtml(a string) string {
	html, _ := ioutil.ReadFile(a)
	return string(html)
}

func load(res http.ResponseWriter, html string) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	fmt.Fprintf(
		res,
		loadHtml(html),
	)
}

func redirect(res http.ResponseWriter, option string) {
	switch option {
	case ADD_STUDENT:
		load(res, "./html/student.html")
	}
}


func root(res http.ResponseWriter, req *http.Request) {
	load(res, "index.html")
}


func AddStudentGradeBySubject(res http.ResponseWriter, req *http.Request) {
	load(res, "./html/student.html")
	switch req.Method {
	case "POST":
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(res, "ParseForm() error %v", err)
			return
		}

		var subjects = make(map[string]map[string]float64)
		var subject = req.FormValue("subject")
		var studentName = req.FormValue("student")
		var grade = req.FormValue("grade")
		student := make(map[string]float64)
		gradef, err := strconv.ParseFloat(grade, 64)
		if err != nil {
			gradef = 0
		}
		student[studentName] = gradef
		subjects[subject] = student

		var result string
		client.Call("Server.AddStudentGradeBySubject", subjects, &result)
	}
}

func GetStudentsAverage(res http.ResponseWriter, req *http.Request) {
	var result string
	client.Call("Server.GetStudentsAverage", "", &result)
	s := strings.Split(result, ";")
	var content = ""
	for i := 1; i < len(s); i++ {
		// name
		if i%2 != 0 {
			content += "<tr><td>" + s[i] + "</td>"
		} else {
			// grade
			content += "<td>" + s[i] + "</td></tr>"
		}
	}
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	fmt.Fprintf(
		res,
		loadHtml("./html/students-average.html"),
		content,
	)
}

func GetGeneralAverage(res http.ResponseWriter, req *http.Request) {
	var result string
	client.Call("Server.GetStudentsGeneralAverage", "", &result)
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	fmt.Fprintf(
		res,
		loadHtml("./html/general-average.html"),
		result,
	)
	
}

func GetSubjectsAverage(res http.ResponseWriter, req *http.Request) {
	var result string
	client.Call("Server.GetSubjectsAverage", "", &result)
	s := strings.Split(result, ";")
	var content = ""
	for i := 1; i < len(s); i++ {
		// name
		if i%2 != 0 {
			content += "<tr><td>" + s[i] + "</td>"
		} else {
			// grade
			content += "<td>" + s[i] + "</td></tr>"
		}
	}
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	fmt.Fprintf(
		res,
		loadHtml("./html/subject-average.html"),
		content,
	)
}



func main() {
	c, err := rpc.Dial("tcp", "127.0.0.1:9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	client = c
	// client = c
	http.HandleFunc("/", root)
	http.HandleFunc("/student", AddStudentGradeBySubject)
	http.HandleFunc("/students-average", GetStudentsAverage)
	http.HandleFunc("/general-average", GetGeneralAverage)
	http.HandleFunc("/subjects-average", GetSubjectsAverage)
	
	fmt.Println("Running server in PORT 9000...")
	http.ListenAndServe(":9000", nil)
}
