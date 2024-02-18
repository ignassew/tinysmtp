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

func NewClient(conn net.Conn) *Client {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	rw := bufio.NewReadWriter(reader, writer)
	return &Client{Conn{*rw}}
}

// spec: https://datatracker.ietf.org/doc/html/rfc5321#section-4.2
type ReplyCode int

const (
	// System status, or system help reply
	SystemStatus ReplyCode = 211

	// Help message (Information on how to use the receiver or the
	// meaning of a particular non-standard command; this reply is useful
	// only to the human user)
	HelpMessage ReplyCode = 214

	// Service ready
	ServiceReady ReplyCode = 220

	// Service closing transmission channel
	ServiceClosing ReplyCode = 221

	// Requested mail action okay, completed
	MailActionOK ReplyCode = 250

	// User not local; will forward to <forward-path> (See Section 3.4)
	ForwardingUser ReplyCode = 251

	// Cannot VRFY user, but will accept message and attempt delivery
	// (See Section 3.5.3)
	CannotVRFY ReplyCode = 252

	// Start mail input; end with <CRLF>.<CRLF>
	StartMailInput ReplyCode = 354

	// Service not available, closing transmission channel
	// (This may be a reply to any command if the service knows it must
	// shut down)
	ServiceNotAvailable ReplyCode = 421

	// Requested mail action not taken: mailbox unavailable (e.g.,
	// mailbox busy or temporarily blocked for policy reasons)
	MailboxUnavailable ReplyCode = 450

	// Requested action aborted: local error in processing
	ActionAborted ReplyCode = 451

	// Requested action not taken: insufficient system storage
	ActionNotTaken ReplyCode = 452

	// Server unable to accommodate parameters
	ServerUnable ReplyCode = 455

	// Syntax error, command unrecognized (This may include errors such
	// as command line too long)
	SyntaxError ReplyCode = 500

	// Syntax error in parameters or arguments
	ParameterSyntaxError ReplyCode = 501

	// Command not implemented (see Section 4.2.4)
	CommandNotImplemented ReplyCode = 502

	// Bad sequence of commands
	BadSequence ReplyCode = 503

	// Command parameter not implemented
	ParameterNotImplemented ReplyCode = 504

	// Requested action not taken: mailbox unavailable (e.g., mailbox
	// not found, no access, or command rejected for policy reasons)
	MailboxNotAvailable ReplyCode = 550

	// User not local; please try <forward-path>
	UserNotLocal ReplyCode = 551

	// Requested mail action aborted: exceeded storage allocation
	ActionAbortedStorage ReplyCode = 552

	// Requested action not taken: mailbox name not allowed (e.g.,
	// mailbox syntax incorrect)
	MailboxNameNotAllowed ReplyCode = 553

	// Transaction failed (Or, in the case of a connection-opening
	// response, "No SMTP service here")
	TransactionFailed ReplyCode = 554

	// MAIL FROM/RCPT TO parameters not recognized or not implemented
	ParameterNotRecognized ReplyCode = 555
)

// spec: https://datatracker.ietf.org/doc/html/rfc5321#section-4.2
//
// returns: ReplyCode, message, isMultiline
func (c *Client) recv_command() (ReplyCode, string, bool) {
	s, err := c.conn.ReadString('\n')
	if err != nil {
		panic("fuck")
	}

	codeInt, _ := strconv.Atoi(s[0:3])
	code := ReplyCode(codeInt)

	isMultiline := s[3] == '-'

	var message string

	if isMultiline {
		message = s[4:]
		nextReplyCode, nextMessage, nextIsMultiline := c.recv_command()
		isMultiline = nextIsMultiline

		if nextReplyCode != code {
			panic("reply code is not the same")
		}
		message += nextMessage
	} else {
		message = s[3:]
	}

	return code, strings.TrimSpace(message), isMultiline
}

// spec: https://datatracker.ietf.org/doc/html/rfc5321#section-4.1.1.1
func (c *Client) ehlo(domain string) {
	c.conn.writeLine("EHLO %s", domain)
	r, m, _ := c.recv_command()
	fmt.Printf("%v: %s\n", r, m)
}

// spec: https://datatracker.ietf.org/doc/html/rfc5321#section-4.1.1.1
func (c *Client) helo(domain string) {
	c.conn.writeLine("HELO %s", domain)
	r, m, _ := c.recv_command()
	fmt.Printf("%v: %s\n", r, m)
}

// spec: https://datatracker.ietf.org/doc/html/rfc5321#section-4.1.1.2
func (c *Client) mail(reversePath string) {
	c.conn.writeLine("MAIL FROM: <%s>", reversePath)
	r, m, _ := c.recv_command()
	fmt.Printf("%v: %s\n", r, m)
}

// spec: https://datatracker.ietf.org/doc/html/rfc5321#section-4.1.1.3
func (c *Client) mailRecipient(forwardPath string) {
	c.conn.writeLine("RCPT TO: <%s>", forwardPath)
	r, m, _ := c.recv_command()
	fmt.Printf("%v: %s\n", r, m)
}

// spec: https://datatracker.ietf.org/doc/html/rfc5321#section-4.1.1.4
func (c *Client) mailData(data string) {
	c.conn.writeLine("DATA")
	{
		r, m, _ := c.recv_command()
		fmt.Printf("%v: %s\n", r, m)
	}

	c.conn.WriteString(data)
	c.conn.WriteString("\r\n.\r\n")
	c.conn.Flush()
	{
		r, m, _ := c.recv_command()
		fmt.Printf("%v: %s\n", r, m)
	}
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

func (c *Client) SendMail(from string, to string, mail *mail) {
	c.mail(from)
	c.mailRecipient(to)
	c.mailData(fmt.Sprint(mail))
}
