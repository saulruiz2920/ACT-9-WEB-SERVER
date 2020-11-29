package main

import (
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"strconv"
)

type Server struct{
	subjects 	map[string]map[string]float64
	students  	map[string]map[string]float64
}

func (this *Server) AddStudentGradeBySubject(data map[string]map[string]float64, reply *string) error {
	fmt.Print("METHOD: AddStudentGradeBySubject, PARAMS: map subjects: ")
	fmt.Println(data)
	for subject, student := range data {
		for studentName, grade := range student {
			// Create a subject with a student ELSE add student with grade to existing subject
			if this.subjects[subject] == nil {
				this.subjects[subject] = student
			} else {
				this.subjects[subject][studentName] = grade
			}
			// Create a student with a subject ELSE add subject with grade to existing student
			if this.students[studentName] == nil {
				subj := make(map[string]float64)
				subj[subject] = grade
				this.students[studentName] = subj
			} else {
				this.students[studentName][subject] = grade
			}
		}
	}
	*reply = "Success"
	return nil
}

func (this *Server) GetStudentsGeneralAverage(noData string, reply *string) error {
	fmt.Println("METHOD: GetStudentsGeneralAverage")
	var avrg float64
	var elements float64 = 0
	for _, student := range this.subjects {
		for _, grade := range student {
			avrg += grade
			elements += 1
		}
	}
	if elements <= 0 {
		err := errors.New("No information")
		fmt.Println(err)
        return err
	}
	*reply = FloatToString(avrg/elements)
	return nil
}

func (this *Server) GetSubjectAverage(subject string, reply *string) error {
	fmt.Println("METHOD: GetSubjectAverage,         PARAMS: subject: " + subject)
	_, ok := this.subjects[subject]
    if ok {
		var avrg float64
		for _, grade := range this.subjects[subject] {
			avrg += grade
		}
		var totalcount = float64(len(this.subjects[subject]))
		avrg = avrg / totalcount
		*reply = FloatToString(avrg)
    } else {
		err := errors.New("Subject not found. Value: " + subject)
		fmt.Println(err)
        return err
    }
	return nil
}

func GetSubjectAverage(this *Server, subject string) string {
	var avrg float64
	for _, grade := range this.subjects[subject] {
		avrg += grade
	}
	var totalcount = float64(len(this.subjects[subject]))
	avrg = avrg / totalcount
	return FloatToString(avrg)
}

func (this *Server) GetSubjectsAverage(subject string, reply *string) error {
	fmt.Println("METHOD: GetSubjectsAverage")
	var subjects string
	for subject, _ := range this.subjects {
		subjects += ";" + subject + ";" + GetSubjectAverage(this, subject)
	}
	*reply = subjects
	return nil
}

func GetStudentAverage(this *Server, studentName string) string {
	var avrg float64
	for _, grade := range this.students[studentName] {
		avrg += grade
	}
	var totalcount = float64(len(this.students[studentName]))
	avrg = avrg / totalcount
	return FloatToString(avrg)
}

func (this *Server) GetStudentsAverage(noData string, reply *string) error {
	fmt.Println("METHOD: GetStudentsAverage")
	var students string
	for student, _ := range this.students {
		students += ";" + student + ";" + GetStudentAverage(this, student)
	}
	*reply = students
	return nil
}


func FloatToString(float float64) string {
	// f float
   	// 2 bit precision
   	// 64 bit float
	return strconv.FormatFloat(float, 'f', 2, 64)
}

/**
* Creates a Server struct and initialize its maps
*/
func InitServer() *Server {
	var server = new(Server)
	server.subjects = make(map[string]map[string]float64)
	server.students = make(map[string]map[string]float64)
	return server
}

func server() {
	rpc.Register(InitServer())
	ln, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Server running")
	for {
		c, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go rpc.ServeConn(c)
	}
}

func main() {
	go server()

	var input string
	fmt.Scanln(&input)
}