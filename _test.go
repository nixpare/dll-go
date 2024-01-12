//go:generate dll-go -output hello-world-dll.go main.go
//go:generate go build -o hello-world.dll -buildmode=c-shared .
package main

import (
	"errors"
	"fmt"
	"time"
)

type message struct {
	sender string
	receiver string
	date time.Time
	msg string
	err error
}

func main() {
	msg := &message{
		sender: "Ale",
		receiver: "Eli",
		date: time.Now(),
		msg: "Ciao eli",
	}
	fmt.Println(Print(msg))

	msg.sender = "Eli"
	msg.receiver = "Ale"
	msg.msg = "Ciao, come stai?"
	msg.date = time.Now()
	msg.err = errors.New("ERRORE")
	fmt.Println(PrintDLL(msg))
}

//dll Print(msg *message) (n int, b error) = ./hello-world.dll
func Print(msg *message) (int, error) {
	n, _ := fmt.Printf(
		"%s, you received a message from %s at %v\n%s\n",
		msg.receiver, msg.sender, msg.date, msg.msg,
	)
	return n, msg.err
}
