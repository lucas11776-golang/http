package mail

import (
	"io"

	"github.com/go-gomail/gomail"
	"github.com/lucas11776-golang/http/utils/env"
)

type Mail struct {
	host     string
	port     int
	username string
	password string
	message  *gomail.Message
}

type Address struct {
	Name  string
	Email string
}

// Comment
func NewMail() *Mail {
	return &Mail{
		host:     env.Env("MAIL_HOST"),
		port:     env.EnvInt("MAIL_PORT"),
		username: env.Env("MAIL_USERNAME"),
		password: env.Env("MAIL_PASSWORD"),
		message:  gomail.NewMessage(),
	}
}

// Comment
func (ctx *Mail) Host(host string) *Mail {
	ctx.host = host

	return ctx
}

// Comment
func (ctx *Mail) Port(port int) *Mail {
	ctx.port = port

	return ctx
}

// Comment
func (ctx *Mail) Username(username string) *Mail {
	ctx.username = username

	return ctx
}

// Comment
func (ctx *Mail) Password(password string) *Mail {
	ctx.password = password

	return ctx
}

// Comment
func (ctx *Mail) Attachment(path ...string) *Mail {
	for _, p := range path {
		ctx.message.Attach(p)
	}

	return ctx
}

// Comment
func (ctx *Mail) Attach(name string, content []byte) *Mail {
	ctx.message.Attach(name, gomail.SetCopyFunc(func(w io.Writer) error {
		_, err := w.Write(content)
		return err
	}))

	return ctx
}

// Comment
func (ctx *Mail) From(email string) *Mail {
	ctx.message.SetHeader("From", email)

	return ctx
}

// Comment
func (ctx *Mail) FromAddress(address *Address) *Mail {
	ctx.message.SetAddressHeader("From", address.Email, address.Name)

	return ctx
}

// Comment
func (ctx *Mail) To(email ...string) *Mail {
	ctx.message.SetHeader("To", email...)

	return ctx
}

// Comment
func (ctx *Mail) ToAddress(address ...*Address) *Mail {
	for _, a := range address {
		ctx.message.SetAddressHeader("To", a.Email, a.Name)
	}

	return ctx
}

// Comment
func (ctx *Mail) Cc(email ...string) *Mail {
	ctx.message.SetHeader("Cc", email...)

	return ctx
}

// Comment
func (ctx *Mail) CcAddress(address ...Address) *Mail {
	for _, a := range address {
		ctx.message.SetAddressHeader("Cc", a.Email, a.Name)
	}

	return ctx
}

// Comment
func (ctx *Mail) Subject(subject string) *Mail {
	ctx.message.SetHeader("Subject", subject)
	return ctx
}

// Comment
func (ctx *Mail) SendText(body string) error {
	ctx.message.SetBody("text/plain", body)

	return ctx.send()
}

// Comment
func (ctx *Mail) SendHtml(body string) error {
	ctx.message.SetBody("text/html", body)

	return ctx.send()
}

// comment
func (ctx *Mail) send() error {
	return gomail.NewDialer(ctx.host, ctx.port, ctx.username, ctx.password).DialAndSend(ctx.message)
}
