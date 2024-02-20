package tinysmtp

import (
	"errors"
	"fmt"
)

// spec: https://datatracker.ietf.org/doc/html/rfc5321#section-4.2
type ReplyCode struct {
	code    int
	message string
	err     error
}

var (
	// System status, or system help reply
	SystemStatus ReplyCode = ReplyCode{211, "", nil}

	// Help message (Information on how to use the receiver or the
	// meaning of a particular non-standard command; this reply is useful
	// only to the human user)
	HelpMessage ReplyCode = ReplyCode{214, "", nil}

	// Service ready
	ServiceReady ReplyCode = ReplyCode{220, "", nil}

	// Service closing transmission channel
	ServiceClosing ReplyCode = ReplyCode{221, "", nil}

	// Requested mail action okay, completed
	MailActionOK ReplyCode = ReplyCode{250, "", nil}

	// User not local; will forward to <forward-path> (See Section 3.4)
	ForwardingUser ReplyCode = ReplyCode{251, "", nil}

	// Cannot VRFY user, but will accept message and attempt delivery
	// (See Section 3.5.3)
	CannotVRFY ReplyCode = ReplyCode{252, "", nil}

	StartMailInput ReplyCode = ReplyCode{354, "", nil}

	// Service not available, closing transmission channel
	// (This may be a reply to any command if the service knows it must
	// shut down)
	ServiceNotAvailable ReplyCode = ReplyCode{
		421,
		"",
		nil,
	}

	// Requested mail action not taken: mailbox unavailable (e.g.,
	// mailbox busy or temporarily blocked for policy reasons)
	MailboxUnavailable ReplyCode = ReplyCode{
		450,
		"",
		nil,
	}

	ActionAborted ReplyCode = ReplyCode{
		451,
		"",
		nil,
	}

	ActionNotTaken ReplyCode = ReplyCode{
		452,
		"",
		nil,
	}

	ServerUnable ReplyCode = ReplyCode{
		455,
		"",
		nil,
	}

	// Syntax error, command unrecognized (This may include errors such
	// as command line too long)
	SyntaxErrorCode ReplyCode = ReplyCode{
		500,
		"",
		nil,
	}

	ParameterSyntaxErrorCode ReplyCode = ReplyCode{
		501,
		"",
		nil,
	}

	CommandNotImplemented ReplyCode = ReplyCode{
		502,
		"",
		nil,
	}

	BadSequence ReplyCode = ReplyCode{
		503,
		"",
		nil,
	}

	ParameterNotImplemented ReplyCode = ReplyCode{
		504,
		"",
		nil,
	}

	// Requested action not taken: mailbox unavailable (e.g., mailbox
	// not found, no access, or command rejected for policy reasons)
	MailboxNotAvailable ReplyCode = ReplyCode{
		550,
		"",
		nil,
	}

	UserNotLocal ReplyCode = ReplyCode{
		551,
		"",
		nil,
	}

	ActionAbortedStorage ReplyCode = ReplyCode{
		552,
		"",
		nil,
	}

	// Requested action not taken: mailbox name not allowed (e.g.,
	// mailbox syntax incorrect)
	MailboxNameNotAllowed ReplyCode = ReplyCode{
		553,
		"",
		nil,
	}

	// Transaction failed (Or, in the case of a connection-opening
	// response, "No SMTP service here")
	TransactionFailed ReplyCode = ReplyCode{
		554,
		"",
		nil,
	}

	ParameterNotRecognized ReplyCode = ReplyCode{
		555,
		"",
		nil,
	}
)

func NewReplyCode(code int, message string) *ReplyCode {
	var r ReplyCode
	switch code {
	case SystemStatus.code:
		r = SystemStatus
		r.message = message
	case HelpMessage.code:
		r = HelpMessage
		r.message = message
	case ServiceReady.code:
		r = ServiceReady
		r.message = message
	case ServiceClosing.code:
		r = ServiceClosing
		r.message = message
	case MailActionOK.code:
		r = MailActionOK
		r.message = message
	case ForwardingUser.code:
		r = ForwardingUser
		r.message = message
	case CannotVRFY.code:
		r = CannotVRFY
		r.message = message
	case StartMailInput.code:
		r = StartMailInput
		r.message = message
	case ServiceNotAvailable.code:
		r = ServiceNotAvailable
		r.message = message
		r.err = errors.New(fmt.Sprintf("%d: service not available: %s", r.code, r.message))
	case MailboxUnavailable.code:
		r = MailboxUnavailable
		r.message = message
		r.err = errors.New(fmt.Sprintf("%d: requested mail action not taken: mailbox unavailable: %s", r.code, r.message))
	case ActionAborted.code:
		r = ActionAborted
		r.message = message
		r.err = errors.New(fmt.Sprintf("%d: requested action aborted: local error in processing: %s", r.code, r.message))
	case ActionNotTaken.code:
		r = ActionNotTaken
		r.message = message
		r.err = errors.New(fmt.Sprintf("%d: requested action not taken: insufficient system storage: %s", r.code, r.message))
	case ServerUnable.code:
		r = ServerUnable
		r.message = message
		r.err = errors.New(fmt.Sprintf("%d: server unable to accommodate parameters: %s", r.code, r.message))
	case SyntaxErrorCode.code:
		r = SyntaxErrorCode
		r.message = message
		r.err = errors.New(fmt.Sprintf("%d: syntax error, command unrecognized: %s", r.code, r.message))
	case ParameterSyntaxErrorCode.code:
		r = ParameterSyntaxErrorCode
		r.message = message
		r.err = errors.New(fmt.Sprintf("%d: syntax error in parameters or arguments: %s", r.code, r.message))
	case CommandNotImplemented.code:
		r = CommandNotImplemented
		r.message = message
		r.err = errors.New(fmt.Sprintf("%d: command not implemented: %s", r.code, r.message))
	case BadSequence.code:
		r = BadSequence
		r.message = message
		r.err = errors.New(fmt.Sprintf("%d: bad sequence of commands: %s", r.code, r.message))
	case ParameterNotImplemented.code:
		r = ParameterNotImplemented
		r.message = message
		r.err = errors.New(fmt.Sprintf("%d: command parameter not implemented: %s", r.code, r.message))
	case MailboxNotAvailable.code:
		r = MailboxNotAvailable
		r.message = message
		r.err = errors.New(fmt.Sprintf("%d: requested action not taken: mailbox unavailable: %s", r.code, r.message))
	case UserNotLocal.code:
		r = UserNotLocal
		r.message = message
		r.err = errors.New(fmt.Sprintf("%d: user not local; please try <forward-path>: %s", r.code, r.message))
	case ActionAbortedStorage.code:
		r = ActionAbortedStorage
		r.message = message
		r.err = errors.New(fmt.Sprintf("%d: requested mail action aborted: exceeded storage allocation: %s", r.code, r.message))
	case MailboxNameNotAllowed.code:
		r = MailboxNameNotAllowed
		r.message = message
		r.err = errors.New(fmt.Sprintf("%d: requested action not taken: mailbox name not allowed: %s", r.code, r.message))
	case TransactionFailed.code:
		r = TransactionFailed
		r.message = message
		r.err = errors.New(fmt.Sprintf("%d: transaction failed: %s", r.code, r.message))
	case ParameterNotRecognized.code:
		r = ParameterNotRecognized
		r.message = message
		r.err = errors.New(fmt.Sprintf("%d: MAIL FROM/RCPT TO parameters not recognized or not implemented: %s", r.code, r.message))
	default:
		if code/100 == 5 || code/100 == 4 {
			r.code = code
			r.message = message
			r.err = errors.New(fmt.Sprintf("%d: %s", r.code, r.message))
		} else {
			r = ReplyCode{code, message, nil}
		}
	}

	return &r
}
