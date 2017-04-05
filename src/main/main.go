package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"time"

	_ "net/http/pprof"

	"sheepbao.com/glog"
	"sheepbao.com/media/av"
	"sheepbao.com/media/container/flv"
	"sheepbao.com/media/protocol/hls"
	"sheepbao.com/media/protocol/httpflv"
	"sheepbao.com/media/protocol/httpopera"
	"sheepbao.com/media/protocol/rtmp"
)

const (
	programName = "SMS"
	VERSION     = "1.1.1"
)

var (
	buildTime string
	prof      = flag.String("pprofAddr", "", "golang pprof debug address.")
	rtmpAddr  = flag.String("rtmpAddr", ":1935", "The rtmp server address to bind.")
	flvAddr   = flag.String("flvAddr", ":8081", "the http-flv server address to bind.")
	hlsAddr   = flag.String("hlsAddr", ":8080", "the hls server address to bind.")
	operaAddr = flag.String("operaAddr", "", "the http operation or config address to bind: 8082.")
	flvDvr    = flag.Bool("flvDvr", false, "enable flv dvr")
)

var (
	Getters []av.GetWriter
)

func BuildTime() string {
	return buildTime
}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s Version[%s]\r\nUsage: %s [OPTIONS]\r\n", programName, VERSION, os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
}

func catchSignal() {
	// windows unsupport syscall.SIGSTOP
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGTERM)
	<-sig
	glog.Println("recieved signal!")
	// select {}
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			glog.Errorln("main panic: ", r)
			time.Sleep(1 * time.Second)
		}
	}()

	stream := rtmp.NewRtmpStream()
	// hls
	startHls()
	// flv dvr
	startFlvDvr()
	// rtmp
	startRtmp(stream, Getters)
	// http-flv
	startHTTPFlv(stream)
	// http-opera
	startHTTPOpera(stream)
	// pprof
	startPprof()
	// my log
	mylog()
	// block
	catchSignal()
}

func startHls() *hls.Server {
	hlsListen, err := net.Listen("tcp", *hlsAddr)
	if err != nil {
		glog.Fatal(err)
	}

	hlsServer := hls.NewServer()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				glog.Errorln("hls server panic: ", r)
			}
		}()
		hlsServer.Serve(hlsListen)
	}()
	Getters = append(Getters, hlsServer)
	return hlsServer
}

func startRtmp(stream *rtmp.RtmpStream, getters []av.GetWriter) {
	rtmplisten, err := net.Listen("tcp", *rtmpAddr)
	if err != nil {
		glog.Fatal(err)
	}

	rtmpServer := rtmp.NewRtmpServer(stream, getters)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				glog.Errorln("hls server panic: ", r)
			}
		}()
		rtmpServer.Serve(rtmplisten)
	}()
}

func startHTTPFlv(stream *rtmp.RtmpStream) {
	flvListen, err := net.Listen("tcp", *flvAddr)
	if err != nil {
		glog.Fatal(err)
	}

	hdlServer := httpflv.NewServer(stream)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				glog.Errorln("hls server panic: ", r)
			}
		}()
		hdlServer.Serve(flvListen)
	}()
}

func startHTTPOpera(stream *rtmp.RtmpStream) {
	if *operaAddr != "" {
		opListen, err := net.Listen("tcp", *operaAddr)
		if err != nil {
			glog.Fatal(err)
		}
		opServer := httpopera.NewServer(stream)
		go func() {
			defer func() {
				if r := recover(); r != nil {
					glog.Errorln("hls server panic: ", r)
				}
			}()
			opServer.Serve(opListen)
		}()
	}
}

func startFlvDvr() {
	if *flvDvr {
		fd := new(flv.FlvDvr)
		Getters = append(Getters, fd)
	}
}

func startPprof() {
	if *prof != "" {
		go func() {
			if err := http.ListenAndServe(*prof, nil); err != nil {
				glog.Fatal("enable pprof failed: ", err)
			}
		}()
	}
}

func mylog() {
	fmt.Println("")
	glog.Printf("SMS Version:  %s\tBuildTime:  %s\n", VERSION, BuildTime())
	glog.Printf("SMS Start, Rtmp Listen On %s\n", *rtmpAddr)
	glog.Printf("SMS Start, Hls Listen On %s\n", *hlsAddr)
	glog.Printf("SMS Start, HTTP-flv Listen On %s\n", *flvAddr)
	if *operaAddr != "" {
		glog.Printf("SMS Start, HTTP-Operation Listen On %s\n", *operaAddr)
	}
	if *prof != "" {
		glog.Printf("SMS Start, Pprof Server Listen On %s\n", *prof)
	}
	if *flvDvr {
		glog.Printf("SMS Start, Flv Dvr Save On [%s]", "app/streamName.flv")
	}
	SavePid()
}
