package http

import (
	"fmt"
	"io/fs"
	"math/rand"
	"strconv"
	"strings"
	"testing"

	"github.com/open2b/scriggo"
	"github.com/open2b/scriggo/native"
)

func TestView(t *testing.T) {
	t.Run("TestReader", func(t *testing.T) {
		world := int(rand.Float64() * 10000)
		view := NewView(&viewReaderTest{
			Files: scriggo.Files{
				"simple.html": []byte(strings.Join([]string{
					`<h1>Hello World: {{ world }}</h1>`,
				}, "\r\n")),
			},
		}, "html")

		data, err := view.Read("simple", ViewData{
			"world": world,
		}, nil)

		if err != nil {
			t.Fatalf("Failed to parse view: %s", err.Error())
		}

		expected := strings.Join([]string{"<h1>Hello World: ", strconv.Itoa(world), "</h1>"}, "")

		if expected != string(data) {
			t.Fatalf("Expected view to be (%s) but got (%s)", expected, string(data))
		}
	})

	t.Run("TestIfElse", func(t *testing.T) {
		view := NewView(&viewReaderTest{
			Files: scriggo.Files{
				"if.html": []byte(strings.Join([]string{
					`<h1>{% if age < 18 %}You can not drive{% else if age >= 21 %}You can drive code 12 or 14{% else %}You can drive code 10{% end %}</h1>`,
				}, "\r\n")),
			},
		}, "html")

		data, err := view.Read("if", ViewData{
			"age": 17,
		}, nil)

		if err != nil {
			t.Fatalf("Failed to parse view: %s", err.Error())
		}

		expected := "<h1>You can not drive</h1>"

		if expected != string(data) {
			t.Fatalf("Expected view to be (%s) but got (%s)", expected, string(data))
		}

		data, err = view.Read("if", ViewData{
			"age": 18,
		}, nil)

		if err != nil {
			t.Fatalf("Failed to parse view: %s", err.Error())
		}

		expected = "<h1>You can drive code 10</h1>"

		if expected != string(data) {
			t.Fatalf("Expected view to be (%s) but got (%s)", expected, string(data))
		}

		data, err = view.Read("if", ViewData{
			"age": 21,
		}, nil)

		if err != nil {
			t.Fatalf("Failed to parse view: %s", err.Error())
		}

		expected = "<h1>You can drive code 12 or 14</h1>"

		if expected != string(data) {
			t.Fatalf("Expected view to be (%s) but got (%s)", expected, string(data))
		}
	})

	t.Run("TestFor", func(t *testing.T) {
		view := NewView(&viewReaderTest{
			Files: scriggo.Files{
				"for.html": []byte(strings.Join([]string{
					`<ul>{% for city in cities %}<li>{{ city }}</li>{% end %}</ul>`,
				}, "\r\n")),
			},
		}, "html")

		cities := []string{"Pretoria", "New York", "Cape Town"}

		data, err := view.Read("for", ViewData{
			"cities": &cities,
		}, nil)

		if err != nil {
			t.Fatalf("Failed to parse view: %s", err.Error())
		}

		expected := strings.Join([]string{
			"<ul>",
			"<li>Pretoria</li>",
			"<li>New York</li>",
			"<li>Cape Town</li>",
			"</ul>",
		}, "")

		if expected != string(data) {
			t.Fatalf("Expected view to be (%s) but got (%s)", expected, string(data))
		}
	})

	t.Run("TestSwitch", func(t *testing.T) {
		view := NewView(&viewReaderTest{
			Files: scriggo.Files{
				"switch.html": []byte(strings.Join([]string{
					`<h1>{% switch role %}{% case "user" %}You are a user{% default %}You are a guest{% end %}</h1>`,
				}, "\r\n")),
			},
		}, "html")

		data, err := view.Read("switch", ViewData{
			"role": "guest",
		}, nil)

		if err != nil {
			t.Fatalf("Failed to parse view: %s", err.Error())
		}

		expected := "<h1>You are a guest</h1>"

		if expected != string(data) {
			t.Fatalf("Expected view to be (%s) but got (%s)", expected, string(data))
		}

		data, err = view.Read("switch", ViewData{
			"role": "user",
		}, nil)

		if err != nil {
			t.Fatalf("Failed to parse view: %s", err.Error())
		}

		expected = "<h1>You are a user</h1>"

		if expected != string(data) {
			t.Fatalf("Expected view to be (%s) but got (%s)", expected, string(data))
		}
	})

	t.Run("TestSubDirectoryView", func(t *testing.T) {
		view := NewView(&viewReaderTest{
			Files: scriggo.Files{
				"authentication/login.html": []byte(strings.Join([]string{
					"<h1>Login page</h1>",
				}, "\r\n")),
			},
		}, "html")

		data, err := view.Read("authentication.login", ViewData{}, nil)

		if err != nil {
			t.Fatalf("Failed to parse view: %s", err.Error())
		}

		expected := "<h1>Login page</h1>"

		if expected != string(data) {
			t.Fatalf("Expected view to be (%s) but got (%s)", expected, string(data))
		}
	})

	t.Run("TestViewDeclarations", func(t *testing.T) {
		email := "jeo@doe.com"

		type User struct {
			Email string
		}

		newUser := func(email string) *User {
			return &User{Email: email}
		}

		view := NewView(&viewReaderTest{
			Files: scriggo.Files{
				"profile.html": []byte(strings.Join([]string{
					fmt.Sprintf(`<h1>{{ user("%s").Email }}</h1>`, email),
				}, "\r\n")),
			},
		}, "html", native.Declarations{
			"user": newUser,
		})

		data, err := view.Read("profile", ViewData{}, nil)

		if err != nil {
			t.Fatalf("Failed to parse view: %s", err.Error())
		}

		expected := fmt.Sprintf("<h1>%s</h1>", email)

		if expected != string(data) {
			t.Fatalf("Expected view to be (%s) but got (%s)", expected, string(data))
		}
	})
}

type viewReaderTest struct {
	Files scriggo.Files
}

// Comment
func (ctx *viewReaderTest) Open(name string) (fs.File, error) {
	return ctx.Files.Open(name)
}
