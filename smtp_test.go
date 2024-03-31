package tinysmtp

import (
	// "crypto/tls"
	"fmt"
	"net"
	"testing"
	"time"
)

const SERVER_ADDRESS = "localhost:1025"

func prepareClient(t *testing.T, address string) *Client {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		t.Fatalf("failed to connect to the smtp server: %s", err.Error())
	}

	client, err := NewClient(conn)
	if err != nil {
		t.Fatalf("failed to create the client: %s", err.Error())
	}
	return client
}

func TestSendMail(t *testing.T) {
	client := prepareClient(t, SERVER_ADDRESS)

	client.ehlo("example.com")

	err := func() error {
		var mail *mail = NewMail("iggy@example.com", "iggy@example.com", "Canonical Test Message", time.Now(), "", "Hi!\nHope you're doing fine.\nOkay, bye!")
		if err := client.mail("iggy@example.com"); err != nil {
			return err
		}
		if err := client.mailRecipient("iggy@example.com"); err != nil {
			return err
		}
		if err := client.mailData(fmt.Sprint(mail)); err != nil {
			return err
		}
		return nil
	}()
	if err != nil {
		t.Fatalf("failed send an email: %s", err.Error())
	}
}

func TestVRFY(t *testing.T) {
	client := prepareClient(t, SERVER_ADDRESS)

	client.ehlo("example.com")

	_, err := client.vrfy("iggy@example.com")
	if err != nil {
		t.Fatalf("failed to verify the email address: %s", err.Error())
	}
}
