package tinysmtp

import (
	"fmt"
	"testing"
	"time"
    "crypto/tls"
)

func TestX(t *testing.T) {
    conn, err := tls.Dial("tcp", "smtp.gmail.com:465", nil)
	// conn, err := net.Dial("tcp", "localhost:1025")
	if err != nil {
		panic(fmt.Sprintf("failed to connect to the smtp server: %s", err.Error()))
	}
	client, err := NewClient(conn)
	client.ehlo("google.com")
	err = client.SendMail("szpiren@google.com",
		"szpiren@google.com",
		NewMail(
			"szpiren@google.com",
			"szpiren@google.com",
			"Canonical Test Message",
			time.Now(),
            "",
			"Hi!\nHope you're doing fine.\nOkay, bye!",
		),
	)
	if err != nil {
		panic(fmt.Sprintf("failed send an email: %s", err.Error()))
	}
}
