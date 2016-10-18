package storage

import (
	"log"
	"testing"
)

func TestCreateAlert(t *testing.T) {
	x := &Alert{Name: "emre"}
	s := NewStorage("mongodb://robocop:6Hi3QhgfWfmM@ds013162.mlab.com:13162/tmpmail-dev")
	id := s.CreateAlert(x)
	log.Println(id)

	a := s.GetAlert("57c1f71289c75aed0825c405")
	log.Println(a)
}

func TestCheckToken(t *testing.T) {
	s := NewStorage("mongodb://robocop:6Hi3QhgfWfmM@ds013162.mlab.com:13162/tmpmail-dev")
	res, err := s.CheckToken("gJ292HobqSLA2BVSX4aj-hCM0reGKPSqg-gmM4w-a0M=")

	if err != nil {
		log.Fatal("User not found")
	}

	log.Println(res)
}

func TestCreateUser(t *testing.T) {
	s := NewStorage("mongodb://robocop:6Hi3QhgfWfmM@ds013162.mlab.com:13162/tmpmail-dev")
	token, err := s.CreateUser(&User{Name: "emre", Email: "emrekayan@gmail.com"})

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Ok. Token: %s", token)
}
