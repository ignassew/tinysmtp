package tinysmtp

import (
	"fmt"
	"time"
)

type mail struct {
	from      string
	to        string
	subject   string
	date      time.Time
	messageID string
	body      string
}

func NewMail(
	from string,
	to string,
	subject string,
	date time.Time,
	messageID string,
	body string,
) *mail {
	return &mail{
		from,
		to,
		subject,
		date,
		messageID,
		body,
	}

}

func (m mail) String() string {
	var message string

	message += fmt.Sprintf("From: <%s>\r\n", m.from)
	message += fmt.Sprintf("To: <%s>\r\n", m.to)
	message += fmt.Sprintf("Subject: %s\r\n", m.subject)
	message += fmt.Sprintf("Date: %s\r\n", m.date.Format(time.RFC1123Z))
	message += "\r\n"
	message += m.body

	return message
}
