package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/msteinert/pam"
)

func PAMAuth(username, password string) bool {
	t, err := pam.StartFunc("", "", func(s pam.Style, msg string) (string, error) {

		switch s {
		case pam.PromptEchoOff:
			return password, nil
		case pam.PromptEchoOn:
			return username, nil
		case pam.ErrorMsg:
			log.Print(msg)
			return "", nil
		case pam.TextInfo:
			fmt.Println(msg)
			return "", nil
		}
		return "", errors.New("unrecognized message style")
	})
	if err != nil {
		log.Fatalf("Start: %s", err.Error())
		return false
	}
	err = t.Authenticate(0)
	if err != nil {
		log.Fatalf("authenticate: %s", err.Error())
		return false
	}
	log.Println("authentication succeeded")
	return true
}
