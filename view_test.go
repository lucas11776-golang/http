package http

import (
	"io/fs"
	"math/rand"
	"strconv"
	"strings"
	"testing"

	"github.com/lucas11776-golang/http/utils/reader"
	"github.com/open2b/scriggo"
)

func TestView(t *testing.T) {
	view := InitView(&viewReaderTest{
		cache: make(scriggo.Files),
	}, "html")

	t.Run("TestReader", func(t *testing.T) {
		world := int(rand.Float64() * 10000)

		data, err := view.Read("simple", ViewData{
			"world": world,
		})

		if err != nil {
			t.Fatalf("Failed to parse view: %s", err.Error())
		}

		expected := strings.Join([]string{"<h1>Hello World: ", strconv.Itoa(world), "</h1>"}, "")

		if expected != string(data) {
			t.Fatalf("Expected view to be (%s) but got (%s)", expected, string(data))
		}
	})

	t.Run("TestIfElse", func(t *testing.T) {
		data, err := view.Read("if", ViewData{
			"age": 17,
		})

		if err != nil {
			t.Fatalf("Failed to parse view: %s", err.Error())
		}

		expected := "<h1>You can not drive</h1>"

		if expected != string(data) {
			t.Fatalf("Expected view to be (%s) but got (%s)", expected, string(data))
		}

		data, err = view.Read("if", ViewData{
			"age": 18,
		})

		if err != nil {
			t.Fatalf("Failed to parse view: %s", err.Error())
		}

		expected = "<h1>You can drive code 10</h1>"

		if expected != string(data) {
			t.Fatalf("Expected view to be (%s) but got (%s)", expected, string(data))
		}

		data, err = view.Read("if", ViewData{
			"age": 21,
		})

		if err != nil {
			t.Fatalf("Failed to parse view: %s", err.Error())
		}

		expected = "<h1>You can drive code 12 or 14</h1>"

		if expected != string(data) {
			t.Fatalf("Expected view to be (%s) but got (%s)", expected, string(data))
		}
	})

	t.Run("TestFor", func(t *testing.T) {
		cities := []string{"Pretoria", "New York", "Cape Town"}

		data, err := view.Read("for", ViewData{
			"cities": &cities,
		})

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
		data, err := view.Read("switch", ViewData{
			"role": "guest",
		})

		if err != nil {
			t.Fatalf("Failed to parse view: %s", err.Error())
		}

		expected := "<h1>You are a guest</h1>"

		if expected != string(data) {
			t.Fatalf("Expected view to be (%s) but got (%s)", expected, string(data))
		}

		data, err = view.Read("switch", ViewData{
			"role": "user",
		})

		if err != nil {
			t.Fatalf("Failed to parse view: %s", err.Error())
		}

		expected = "<h1>You are a user</h1>"

		if expected != string(data) {
			t.Fatalf("Expected view to be (%s) but got (%s)", expected, string(data))
		}
	})

	t.Run("TestSubDirectoryView", func(t *testing.T) {
		data, err := view.Read("authentication.login", ViewData{})

		if err != nil {
			t.Fatalf("Failed to parse view: %s", err.Error())
		}

		expected := "<h1>Login page</h1>"

		if expected != string(data) {
			t.Fatalf("Expected view to be (%s) but got (%s)", expected, string(data))
		}
	})
}

var viewReaderTestFS = scriggo.Files{
	"simple.html": []byte(strings.Join([]string{
		`<h1>Hello World: {{ world }}</h1>`,
	}, "\r\n")),
	"for.html": []byte(strings.Join([]string{
		`<ul>{% for city in cities %}<li>{{ city }}</li>{% end %}</ul>`,
	}, "\r\n")),
	"if.html": []byte(strings.Join([]string{
		`<h1>{% if age < 18 %}You can not drive{% else if age >= 21 %}You can drive code 12 or 14{% else %}You can drive code 10{% end %}</h1>`,
	}, "\r\n")),
	"switch.html": []byte(strings.Join([]string{
		`<h1>{% switch role %}{% case "user" %}You are a user{% default %}You are a guest{% end %}</h1>`,
	}, "\r\n")),
	"authentication/login.html": []byte(strings.Join([]string{
		"<h1>Login page</h1>",
	}, "\r\n")),
}

type viewReaderTest struct {
	cache scriggo.Files
}

// Comment
func (ctx *viewReaderTest) Open(name string) (fs.File, error) {
	return viewReaderTestFS.Open(name)
}

// Comment
func (ctx *viewReaderTest) Cache(name string) (scriggo.Files, error) {
	return reader.ReadCache(ctx, ctx.cache, name)
}
