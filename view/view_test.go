package view

import (
	"embed"
	"io/fs"
	"math/rand"
	"strconv"
	"strings"
	"testing"
)

//go:embed views/*
var views embed.FS

type WriterTest struct {
}

// Comment
func (ctx *WriterTest) Open(name string) (fs.File, error) {
	return views.Open(strings.Join([]string{"views", name}, "/"))
}

func TestView(t *testing.T) {
	view := Init(&WriterTest{}, "html")

	t.Run("TestReader", func(t *testing.T) {
		world := int(rand.Float64() * 10000)

		data, err := view.Read("simple", Data{
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
		data, err := view.Read("if", Data{
			"age": 17,
		})

		if err != nil {
			t.Fatalf("Failed to parse view: %s", err.Error())
		}

		expected := "<h1>You can not drive</h1>"

		if expected != string(data) {
			t.Fatalf("Expected view to be (%s) but got (%s)", expected, string(data))
		}

		data, err = view.Read("if", Data{
			"age": 18,
		})

		if err != nil {
			t.Fatalf("Failed to parse view: %s", err.Error())
		}

		expected = "<h1>You can drive code 10</h1>"

		if expected != string(data) {
			t.Fatalf("Expected view to be (%s) but got (%s)", expected, string(data))
		}

		data, err = view.Read("if", Data{
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

		data, err := view.Read("for", Data{
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
		data, err := view.Read("switch", Data{
			"role": "guest",
		})

		if err != nil {
			t.Fatalf("Failed to parse view: %s", err.Error())
		}

		expected := "<h1>You are a guest</h1>"

		if expected != string(data) {
			t.Fatalf("Expected view to be (%s) but got (%s)", expected, string(data))
		}

		data, err = view.Read("switch", Data{
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
}
