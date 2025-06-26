package helper

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/lucas11776-golang/http/utils/env"
	"github.com/spf13/cast"
)

type Cast struct{}

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

// Comment
func Truncate(str string, limit int, suffix string) string {
	if len(str) > limit {
		return fmt.Sprintf("%s%s", str[:limit], suffix)
	}

	return str
}

// Comment
func (ctx *Cast) ToString(value interface{}) string {
	return cast.ToString(value)
}

// Comment
func (ctx *Cast) ToInt8(value interface{}) int8 {
	return cast.ToInt8(value)
}

// Comment
func (ctx *Cast) ToInt16(value interface{}) int16 {
	return cast.ToInt16(value)
}

// Comment
func (ctx *Cast) ToInt32(value interface{}) int32 {
	return cast.ToInt32(value)
}

// Comment
func (ctx *Cast) ToInt(value interface{}) int {
	return cast.ToInt(value)
}

// Comment
func (ctx *Cast) ToInt64(value interface{}) int64 {
	return cast.ToInt64(value)
}

// Comment
func (ctx *Cast) ToFloat32(value interface{}) float32 {
	return cast.ToFloat32(value)
}

// Comment
func (ctx *Cast) ToFloat64(value interface{}) float64 {
	return cast.ToFloat64(value)
}
