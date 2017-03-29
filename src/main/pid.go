package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"sheepbao.com/glog"
)

// save pid
var CurDir string

func init() {
	CurDir = getParentDirectory(getCurrentDirectory())
}
func getParentDirectory(dirctory string) string {
	return substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}
func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		glog.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}
func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}
func SavePid() error {
	pidFilename := CurDir + "/pid/" + filepath.Base(os.Args[0]) + ".pid"
	pid := os.Getpid()
	return ioutil.WriteFile(pidFilename, []byte(strconv.Itoa(pid)), 0755)
}
