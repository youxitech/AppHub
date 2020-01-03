package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func panic400(str string, args ...interface{}) {
	msg := fmt.Sprintf(str, args...)
	panic(&Err{msg, 0, 400})
}

func panicErr(err error) {
	panic400(err.Error())
}

// A-Z, a-z
var chars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

// lower/upper case alphabets
func randomStr(n int) string {
	b := make([]byte, n)

	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}

	return string(b)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
