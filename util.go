package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"cjting.me/apphub/parser"
)

// work around sqlite3 error: can't scan string into time.Time
type MyTime time.Time

var emptyArray = make([]int, 0)

func (t *MyTime) Scan(v interface{}) error {
	var vt time.Time

	if t, ok := v.(time.Time); ok {
		vt = t
	}

	if t, ok := v.(string); ok {
		var err error
		vt, err = time.Parse("2006-01-02 15:04:05", t)
		if err != nil {
			return err
		}
	}

	*t = MyTime(vt)
	return nil
}

func (t MyTime) MarshalJSON() ([]byte, error) {
	return time.Time(t).MarshalJSON()
}

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

// for ios: this is `$short_version.$bundle_version`, e.g. "1.2.3.4"
// for android: this is $VersionName
func getFullVersion(a *parser.AppInfo) string {
	if a.Platform == "ios" {
		return fmt.Sprintf("%s.%s", a.IOSShortVersion, a.IOSBundleVersion)
	}

	return a.AndroidVersionName
}
