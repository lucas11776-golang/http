package mail

import (
	"testing"

	"github.com/lucas11776-golang/http/utils/env"
)

func TestMail(t *testing.T) {
	t.Run("TestMailCredentials", func(t *testing.T) {
		env.Set("MAIL_HOST", "smtp.gmail.com")
		env.Set("MAIL_PORT", "587")
		env.Set("MAIL_USERNAME", "joe")
		env.Set("MAIL_PASSWORD", "test@123")

		mail := NewMail()

		if mail.host != env.Env("MAIL_HOST") {
			t.Fatalf("Expected host to be (%s) but got (%s)", mail.host, env.Env("MAIL_HOST"))
		}

		if mail.port != env.EnvInt("MAIL_PORT") {
			t.Fatalf("Expected port to be (%v) but got (%s)", mail.port, env.Env("MAIL_PORT"))
		}

		if mail.username != env.Env("MAIL_USERNAME") {
			t.Fatalf("Expected username to be (%v) but got (%s)", mail.username, env.Env("MAIL_USERNAME"))
		}

		if mail.password != env.Env("MAIL_PASSWORD") {
			t.Fatalf("Expected password to be (%v) but got (%s)", mail.password, env.Env("MAIL_PASSWORD"))
		}
	})

	// t.Run("TestMailSend", func(t *testing.T) {
	// 	mail := NewMail()

	// 	err := mail.Host("localhost").
	// 		Port(1025).
	// 		FromAddress(&Address{Email: "jane@doe.com", Name: "Jane"}).
	// 		To("jeo@doe.com").
	// 		Subject("Testing with hello message").
	// 		SendHtml("Hello World!!!")

	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}
	// })
}
