package tinysmtp

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

type Conn struct {
	bufio.ReadWriter
}

func (c *Conn) writeLine(format string, args ...any) {
	c.WriteString(fmt.Sprintf(format+"\r\n", args...))
	c.Flush()
}

type Client struct {
	conn Conn
}

func NewClient(conn net.Conn) (*Client, error) {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	rw := bufio.NewReadWriter(reader, writer)
    c := &Client{Conn{*rw}}
	if r, _, _ := c.recvReply(); r.err != nil {
		return nil, r.err
	}
    return c, nil
}

// spec: https://datatracker.ietf.org/doc/html/rfc5321#section-4.2
func (c *Client) recvReply() (replyCode *ReplyCode, message string, isMultiline bool) {
	s, err := c.conn.ReadString('\n')
	if err != nil {
		panic("fuck")
	}

	codeInt, _ := strconv.Atoi(s[0:3])
	replyCode = NewReplyCode(codeInt, message)

	isMultiline = s[3] == '-'

	if isMultiline {
		message = s[4:]
		nextReplyCode, nextMessage, nextIsMultiline := c.recvReply()
		isMultiline = nextIsMultiline

		if nextReplyCode != replyCode {
			panic("reply code is not the same")
		}
		message += nextMessage
	} else {
		message = s[3:]
	}

	return replyCode, strings.TrimSpace(message), isMultiline
}

// spec: https://datatracker.ietf.org/doc/html/rfc5321#section-4.1.1.1
func (c *Client) ehlo(domain string) error {
	c.conn.writeLine("EHLO %s", domain)
	if r, _, _ := c.recvReply(); r.err != nil {
		return r.err
	}
	return nil
}

// spec: https://datatracker.ietf.org/doc/html/rfc5321#section-4.1.1.1
func (c *Client) helo(domain string) error {
	c.conn.writeLine("HELO %s", domain)
	if r, _, _ := c.recvReply(); r.err != nil {
		return r.err
	}
	return nil
}

// spec: https://datatracker.ietf.org/doc/html/rfc5321#section-4.1.1.2
func (c *Client) mail(reversePath string) error {
	c.conn.writeLine("MAIL FROM: <%s>", reversePath)
	if r, _, _ := c.recvReply(); r.err != nil {
		return r.err
	}
	return nil
}

// spec: https://datatracker.ietf.org/doc/html/rfc5321#section-4.1.1.3
func (c *Client) mailRecipient(forwardPath string) error {
	c.conn.writeLine("RCPT TO: <%s>", forwardPath)
	if r, _, _ := c.recvReply(); r.err != nil {
		return r.err
	}
	return nil
}

// spec: https://datatracker.ietf.org/doc/html/rfc5321#section-4.1.1.4
func (c *Client) mailData(data string) error {
	c.conn.writeLine("DATA")
	if r, _, _ := c.recvReply(); r.err != nil {
		return r.err
	}
	c.conn.WriteString(data)
	c.conn.WriteString("\r\n.\r\n")
	c.conn.Flush()
	if r, _, _ := c.recvReply(); r.err != nil {
		return r.err
	}
	return nil
}

// spec: https://datatracker.ietf.org/doc/html/rfc5321#section-4.1.1.5
func (c *Client) reset() error {
	c.conn.writeLine("RSET")
	if r, _, _ := c.recvReply(); r.err != nil {
		return r.err
	}
	return nil
}

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

func (c *Client) SendMail(from string, to string, mail *mail) error {
	if err := c.mail(from); err != nil {
		return err
	}
	if err := c.mailRecipient(to); err != nil {
		return err
	}
	if err := c.mailData(fmt.Sprint(mail)); err != nil {
		return err
	}
	return nil
}
