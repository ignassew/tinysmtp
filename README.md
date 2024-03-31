# tinysmtp

SMTP implementation as per [RFC5321](https://datatracker.ietf.org/doc/html/rfc5321)

## Testing

### Maildev

You can run a test SMTP server with

```sh
docker run -p 1080:1080 -p 1025:1025 maildev/maildev
```

## Client

### TODO

Minimum implementation:

- [x] EHLO
- [x] HELO
- [x] MAIL
- [x] RCPT
- [x] DATA
- [x] RSET
- [x] NOOP
- [x] QUIT
- [ ] VRFY
