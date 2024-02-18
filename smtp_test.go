package tinysmtp

import (
	"net"
	"testing"
	"time"
)

func TestX(t *testing.T) {
    conn, err := net.Dial("tcp", "localhost:1025")
	if err != nil {
		panic("fuck")
	}
	client := *NewClient(conn)
	client.recv_command()
	client.ehlo("google.com")
    client.SendMail("tx@google.com",
        "rx@google.com",
        NewMail(
            "tx@google.com",
            "rx@google.com",
            "Canonical Test Message",
            time.Now(),
            "randomID@google.com",
            `Hi!
Hope you're doing fine.
Okay, bye!`,
        ),
    )
}
