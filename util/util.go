package util

import (
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/viper"
)

// GetSecret fetches a mounted secret defined in the config
func GetSecret(secret string) (string, error) {
	buffer, err := ioutil.ReadFile(viper.GetString(secret))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(buffer)), err
}

// StringInSlice returns if string in slice
// I wish native go included this, but not too big of a deal
// to define it once in a common library...
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// FuncName returns the function name of caller. Used to provide more context
// on traces
// https://stackoverflow.com/questions/10742749/get-name-of-function-using-reflection-in-golang
func FuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	nameFull := runtime.FuncForPC(pc).Name() // main.foo
	nameEnd := filepath.Ext(nameFull)        // .foo
	name := strings.TrimPrefix(nameEnd, ".") // foo
	return name
}
