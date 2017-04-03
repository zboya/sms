package main

import (
	"fmt"
	"path"
	"strings"
)

func main() {
	var pathstr string = "live/03.m3u8"
	p := path.Ext(pathstr)
	fmt.Println(p)
	p = strings.TrimRight("live/03.m3u8", p)
	fmt.Println(p)
}
