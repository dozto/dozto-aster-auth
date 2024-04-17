package helper

import (
	"path/filepath"
	"runtime"
)

var ()

func RootPath() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(b), "../..")
}
