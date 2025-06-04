package env

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Comment
func Load(path string) {
	if err := godotenv.Load(path); err != nil {
		panic(err)
	}
}

// Comment
func Env(key string) string {
	return os.Getenv(key)
}

// Comment
func EnvInt(key string) int {
	n, _ := strconv.Atoi(Env(key))

	return n
}
