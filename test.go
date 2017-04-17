package main

import (
	"fmt"
)

// SimpleChaincode example simple Chaincode implementation

type Error struct {
	Err string
}
type EmployeeWorkDetails interface {
	HoursPerWeek()
	HoursPerDay()
}
type Calendar struct {
	Temp string
}
type OutputTransaction struct {
	To    string
	Value string
}
type Transaction struct {
	Id      string
	Type    string
	From    string
	Inputs  []string
	Outputs []OutputTransaction
}

type TransactionNode struct {
	TransactionId string
	Prev          *TransactionNode
	Next          *TransactionNode
}
type TransacationList struct {
	Current *TransactionNode
}
type LeaveList struct {
	Temp string
}
type Employee struct {
	Id              string
	Name            string
	Contract        string
	Projects        []string
	Location        string
	MaxHoursPerWeek string
	C               Calendar
	L               LeaveList
	Tl              TransacationList
}
type DeliveryManager struct {
	Id       string
	Name     string
	Projects []string
	Tl       TransacationList
}
type Project struct {
	Id        string
	Name      string
	StartDate string
	EndDate   string
	Dm        string
	Tl        TransacationList
}

func (e Employee) HoursPerWeek() string {
	if e.Location == "France" {
		return "32"
	} else if e.Location == "US" {
		return "40"
	} else {
		return "47.5"
	}
}

func (e Employee) HoursPerDay() string {
	if e.Location == "France" {
		return "6.4"
	} else if e.Location == "US" {
		return "8"
	} else {
		return "9.5"
	}
}

func main() {
	// c:=Calendar{"temp Calendar"}
	// l:=LeaveList{"temp LeaveList"}
	id := "akshay_id"
	name := "akshay meher"
	contract := "permanment"
	location := "banglore"
	maxHoursPerWeek := "168"
	calender := Calendar{"temporary_calender"}
	leaveList := LeaveList{"temporary_leavelist"}
	project := []string{"pr1", "pr2"}
	var Tl TransacationList
	Tl.Current = nil
	// Id string
	// Name string
	// Contract string
	// Projects []string
	// Location string
	// MaxHoursPerWeek string
	// C Calendar
	// L LeaveList
	// Tl TransacationList
	var e Employee
	e = Employee{id, name, contract, project, location, maxHoursPerWeek, calender, leaveList, Tl}
	fmt.Println(e)
	// fmt.Println(e.HoursPerDay())
	// fmt.Println(e.HoursPerWeek())

}
