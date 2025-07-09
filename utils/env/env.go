package env

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

// Comment
func Load(path string) {
	if err := godotenv.Load(path); err != nil {
		panic(err)
	}
}

// Comment
func Set(key string, value interface{}) {
	os.Setenv(key, cast.ToString(value))
}

// Comment
func Env(key string) string {
	return os.Getenv(key)
}

// Comment
func EnvInt(key string) int {
	return cast.ToInt(Env(key))
}

// Comment
func EnvInt64(key string) int64 {
	return cast.ToInt64(Env(key))
}

// Comment
func EnvBool(key string) bool {
	return cast.ToBool(Env(key))
}
