package main

type User struct {
	Id    string
	Name  string
	Phone string
}

var users = map[string]*User{
	"1": {
		Id:    "1",
		Name:  "68",
		Phone: "13800000000",
	},
	"2": {
		Id:    "2",
		Name:  "B3",
		Phone: "13800000001",
	},
}
