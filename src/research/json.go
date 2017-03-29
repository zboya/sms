package main

import (
	"encoding/json"
	"fmt"
)

type Conf struct {
	Rtmp struct {
		ListenAddr   string `json:"listenAddr"`
		ReadTimeout  int    `json:"readTimeoutSec"`
		WriteTimeout int    `json:"writeTimeoutSec"`
		ConnBuffSize int    `json:"connBufferSize"`
		GOPNum       int    `json:"gopNum"`
	} `json:"rtmp"`
	Flv struct {
		ListenAddr string `json:"listenAddr"`
	} `json:"flv"`
	Hls struct {
		ListenAddr    string `json:"listenAddr"`
		MaxBufferSize int    `json:"maxBufferSizeKB"`
	} `json:"hls"`

	MaxSubGroupNum     int `json:"maxSubGroupNum"`
	MaxSubGroupWorkers int `json:"maxSubGroupWorkers"`
	MaxBufferSize      int `json:"maxBufferSizeKB"`
}

func main() {
	conf := &Conf{}
	conf.Rtmp.ListenAddr = ":1935"
	conf.Rtmp.ReadTimeout = 30
	conf.Rtmp.WriteTimeout = 30
	conf.Rtmp.ConnBuffSize = 1024
	conf.Rtmp.GOPNum = 1

	b, err := json.Marshal(conf)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(b))
}
