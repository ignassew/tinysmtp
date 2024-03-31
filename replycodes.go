package tinysmtp

import (
	"errors"
	"fmt"
)

type ReplyCode int

const (
	// System status, or system help reply
	SystemStatus = 211

	// Help message (Information on how to use the receiver or the
	// meaning of a particular non-standard command; this reply is useful
	// only to the human user
	HelpMessage = 214

	// Service ready
	ServiceReady = 220

	// Service closing transmission channel
	ServiceClosing = 221

	// Requested mail action okay, completed
	MailActionOK = 250

	// User not local; will forward to <forward-path> (See Section 3.4)
	ForwardingUser = 251

	// Cannot VRFY user, but will accept message and attempt delivery
	// (See Section 3.5.3)
	CannotVRFY = 252

	StartMailInput = 354

	// Service not available, closing transmission channel
	// (This may be a reply to any command if the service knows it must
	// shut down)
	ServiceNotAvailable = 421

	// Requested mail action not taken: mailbox unavailable (e.g.,
	// mailbox busy or temporarily blocked for policy reasons)
	MailboxUnavailable = 450

	ActionAborted  = 451
	ActionNotTaken = 452
	ServerUnable   = 455

	// Syntax error, command unrecognized (This may include errors such
	// as command line too long)
	SyntaxErrorCode = 500

	ParameterSyntaxErrorCode = 501
	CommandNotImplemented    = 502
	BadSequence              = 503
	ParameterNotImplemented  = 504

	// Requested action not taken: mailbox unavailable (e.g., mailbox
	// not found, no access, or command rejected for policy reasons)
	MailboxNotAvailable = 550

	UserNotLocal         = 551
	ActionAbortedStorage = 552

	// Requested action not taken: mailbox name not allowed (e.g.,
	// mailbox syntax incorrect)
	MailboxNameNotAllowed = 553

	// Transaction failed (Or, in the case of a connection-opening
	// response, "No SMTP service here")
	TransactionFailed = 554

	ParameterNotRecognized = 555
)

// spec: https://datatracker.ietf.org/doc/html/rfc5321#section-4.2
type Reply struct {
	code    ReplyCode
	message string
}

func NewReply(code ReplyCode, message string) Reply {
	return Reply{code, message}
}

func (r *Reply) GetError() error {
	switch r.code {
	case ServiceNotAvailable:
		return errors.New(fmt.Sprintf("%d: service not available: %s", r.code, r.message))
	case MailboxUnavailable:
		return errors.New(fmt.Sprintf("%d: requested mail action not taken: mailbox unavailable: %s", r.code, r.message))
	case ActionAborted:
		return errors.New(fmt.Sprintf("%d: requested action aborted: local error in processing: %s", r.code, r.message))
	case ActionNotTaken:
		return errors.New(fmt.Sprintf("%d: requested action not taken: insufficient system storage: %s", r.code, r.message))
	case ServerUnable:
		return errors.New(fmt.Sprintf("%d: server unable to accommodate parameters: %s", r.code, r.message))
	case SyntaxErrorCode:
		return errors.New(fmt.Sprintf("%d: syntax error, command unrecognized: %s", r.code, r.message))
	case ParameterSyntaxErrorCode:
		return errors.New(fmt.Sprintf("%d: syntax error in parameters or arguments: %s", r.code, r.message))
	case CommandNotImplemented:
		return errors.New(fmt.Sprintf("%d: command not implemented: %s", r.code, r.message))
	case BadSequence:
		return errors.New(fmt.Sprintf("%d: bad sequence of commands: %s", r.code, r.message))
	case ParameterNotImplemented:
		return errors.New(fmt.Sprintf("%d: command parameter not implemented: %s", r.code, r.message))
	case MailboxNotAvailable:
		return errors.New(fmt.Sprintf("%d: requested action not taken: mailbox unavailable: %s", r.code, r.message))
	case UserNotLocal:
		return errors.New(fmt.Sprintf("%d: user not local; please try <forward-path>: %s", r.code, r.message))
	case ActionAbortedStorage:
		return errors.New(fmt.Sprintf("%d: requested mail action aborted: exceeded storage allocation: %s", r.code, r.message))
	case MailboxNameNotAllowed:
		return errors.New(fmt.Sprintf("%d: requested action not taken: mailbox name not allowed: %s", r.code, r.message))
	case TransactionFailed:
		return errors.New(fmt.Sprintf("%d: transaction failed: %s", r.code, r.message))
	case ParameterNotRecognized:
		return errors.New(fmt.Sprintf("%d: MAIL FROM/RCPT TO parameters not recognized or not implemented: %s", r.code, r.message))
	default:
		if r.code/100 == 5 || r.code/100 == 4 {
			return errors.New(fmt.Sprintf("%d: %s", r.code, r.message))
		}
	}

	return nil
}
