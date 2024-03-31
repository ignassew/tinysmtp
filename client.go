package tinysmtp

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
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
	if r, _, _ := c.recvReply(); r.GetError() != nil {
		return nil, r.GetError()
	}
	return c, nil
}

// spec: https://datatracker.ietf.org/doc/html/rfc5321#section-4.2
func (c *Client) recvReply() (reply Reply, message string, isMultiline bool) {
	s, err := c.conn.ReadString('\n')
	if err != nil {
		panic("fuck")
	}

	codeInt, _ := strconv.Atoi(s[0:3])
	replyCode := ReplyCode(codeInt)
	reply = NewReply(replyCode, message)

	isMultiline = s[3] == '-'

	if isMultiline {
		message = s[4:]
		nextReplyCode, nextMessage, nextIsMultiline := c.recvReply()
		isMultiline = nextIsMultiline

		if nextReplyCode != reply {
			panic("reply code is not the same")
		}
		message += nextMessage
	} else {
		message = s[3:]
	}

	return reply, strings.TrimSpace(message), isMultiline
}

// spec: https://datatracker.ietf.org/doc/html/rfc5321#section-4.1.1.1
func (c *Client) ehlo(domain string) error {
	c.conn.writeLine("EHLO %s", domain)
	if r, _, _ := c.recvReply(); r.GetError() != nil {
		return r.GetError()
	}
	return nil
}

// spec: https://datatracker.ietf.org/doc/html/rfc5321#section-4.1.1.1
func (c *Client) helo(domain string) error {
	c.conn.writeLine("HELO %s", domain)
	if r, _, _ := c.recvReply(); r.GetError() != nil {
		return r.GetError()
	}
	return nil
}

// spec: https://datatracker.ietf.org/doc/html/rfc5321#section-4.1.1.2
func (c *Client) mail(reversePath string) error {
	c.conn.writeLine("MAIL FROM: <%s>", reversePath)
	if r, _, _ := c.recvReply(); r.GetError() != nil {
		return r.GetError()
	}
	return nil
}

// spec: https://datatracker.ietf.org/doc/html/rfc5321#section-4.1.1.3
func (c *Client) mailRecipient(forwardPath string) error {
	c.conn.writeLine("RCPT TO: <%s>", forwardPath)
	if r, _, _ := c.recvReply(); r.GetError() != nil {
		return r.GetError()
	}
	return nil
}

// spec: https://datatracker.ietf.org/doc/html/rfc5321#section-4.1.1.4
func (c *Client) mailData(data string) error {
	c.conn.writeLine("DATA")
	if r, _, _ := c.recvReply(); r.GetError() != nil {
		return r.GetError()
	}
	c.conn.WriteString(data)
	c.conn.WriteString("\r\n.\r\n")
	c.conn.Flush()
	if r, _, _ := c.recvReply(); r.GetError() != nil {
		return r.GetError()
	}
	return nil
}

// spec: https://datatracker.ietf.org/doc/html/rfc5321#section-4.1.1.5
func (c *Client) reset() error {
	c.conn.writeLine("RSET")
	if r, _, _ := c.recvReply(); r.GetError() != nil {
		return r.GetError()
	}
	return nil
}

// spec: https://datatracker.ietf.org/doc/html/rfc5321#section-4.1.1.9
func (c *Client) noop() error {
	c.conn.writeLine("NOOP")
	r, _, _ := c.recvReply()

	if r.GetError() != nil {
		return r.GetError()
	}

	if r.code != MailActionOK {

	}
	return nil
}
