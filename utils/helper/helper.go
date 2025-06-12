package helper

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/lucas11776-golang/http/utils/env"
	"github.com/spf13/cast"
)

// comment
func Uri(path ...interface{}) string {
	uri := []string{}

	for _, p := range path {
		uri = append(uri, strings.Trim(cast.ToString(p), "/"))
	}

	return strings.Join(uri, "/")
}

// Comment
func Url(path ...interface{}) string {
	return strings.Join([]string{strings.TrimRight(env.Env("APP_URL"), "/"), Uri(path...)}, "/")
}

// Comment
func Subdomain(domain string, path ...interface{}) string {

	u, err := url.Parse(env.Env("APP_URL"))

	if err != nil {
		return ""
	}

	return fmt.Sprintf("%s://%s.%s/%s", u.Scheme, domain, u.Host, Uri(path...))
}

// Comment
func Format(time time.Time, layout string) string {
	return time.Format(layout)
}
