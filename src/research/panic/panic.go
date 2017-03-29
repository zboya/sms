package main

import (
	"time"

	"flag"

	"sheepbao.com/glog"
)

func main() {
	flag.Parse()
	glog.Infoln("start")
	time.Sleep(5 * time.Second)
	p()
	time.Sleep(1)
	glog.Infoln("22222222")
}

func p() {
	panic("test")
}
