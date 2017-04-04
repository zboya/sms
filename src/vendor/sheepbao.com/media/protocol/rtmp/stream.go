package rtmp

import (
	"errors"
	"time"

	"sheepbao.com/glog"
	"sheepbao.com/media/av"
	"sheepbao.com/media/protocol/rtmp/cache"
	"sheepbao.com/media/utils/cmap"
)

var (
	EmptyID = ""
)

type RtmpStream struct {
	streams cmap.ConcurrentMap
}

func NewRtmpStream() *RtmpStream {
	ret := &RtmpStream{
		streams: cmap.New(),
	}
	go ret.CheckAlive()
	return ret
}

func (rs *RtmpStream) HandleReader(r av.ReadCloser) {
	info := r.Info()
	var s *Stream
	item, ok := rs.streams.Get(info.Key)
	if !ok {
		s = NewStream()
		rs.streams.Set(info.Key, s)
	} else {
		s = item.(*Stream)
		if s != nil {
			s.TransStop()
			id := s.ID()
			if id != EmptyID && id != info.UID {
				ns := NewStream()
				s.Copy(ns)
				s = ns
				rs.streams.Set(info.Key, ns)
			}
		}
	}
	s.AddReader(r)
}

func (rs *RtmpStream) HandleWriter(w av.WriteCloser) {
	info := w.Info()
	var s *Stream
	ok := rs.streams.Has(info.Key)
	if !ok {
		s = NewStream()
		rs.streams.Set(info.Key, s)
	} else {
		item, ok := rs.streams.Get(info.Key)
		if ok {
			s = item.(*Stream)
			s.AddWriter(w)
		}
	}

}

func (rs *RtmpStream) GetStreams() cmap.ConcurrentMap {
	return rs.streams
}

func (rs *RtmpStream) CheckAlive() {
	for {
		<-time.After(5 * time.Second)
		for item := range rs.streams.IterBuffered() {
			v := item.Val.(*Stream)
			if v.CheckAlive() == 0 {
				rs.streams.Remove(item.Key)
				if v.r != nil {
					glog.Infof("check alive and remove: %v", v.r.Info())
				}
			}
		}
	}
}

type Stream struct {
	isStart bool
	cache   *cache.Cache
	r       av.ReadCloser
	ws      cmap.ConcurrentMap
}

type PackWriterCloser struct {
	init bool
	w    av.WriteCloser
}

func (p *PackWriterCloser) GetWriter() av.WriteCloser {
	return p.w
}

func NewStream() *Stream {
	return &Stream{
		cache: cache.NewCache(),
		ws:    cmap.New(),
	}
}

func (s *Stream) ID() string {
	if s.r != nil {
		return s.r.Info().UID
	}
	return EmptyID
}

func (s *Stream) GetReader() av.ReadCloser {
	return s.r
}

func (s *Stream) GetWs() cmap.ConcurrentMap {
	return s.ws
}

func (s *Stream) Copy(dst *Stream) {
	for item := range s.ws.IterBuffered() {
		v := item.Val.(*PackWriterCloser)
		s.ws.Remove(item.Key)
		v.w.CalcBaseTimestamp()
		dst.AddWriter(v.w)
	}
}

func (s *Stream) AddReader(r av.ReadCloser) {
	s.r = r
	go s.TransStart()
}

func (s *Stream) AddWriter(w av.WriteCloser) {
	info := w.Info()
	pw := &PackWriterCloser{w: w}
	s.ws.Set(info.UID, pw)
}

func (s *Stream) TransStart() {
	defer func() {
		if s.r != nil {
			glog.Infof("[%v] Transport stop", s.r.Info())
		}
		// debug mode don't use it
		// 	if r := recover(); r != nil {
		// 		glog.Errorln("rtmp TransStart panic: ", r)
		// 	}
	}()

	s.isStart = true
	var p av.Packet
	for {
		if !s.isStart {
			s.isStart = false
			s.closeInter()
			return
		}
		err := s.r.Read(&p)
		if err != nil {
			s.isStart = false
			s.closeInter()
			return
		}
		s.cache.Write(p)

		for item := range s.ws.IterBuffered() {
			v := item.Val.(*PackWriterCloser)
			if !v.init {
				if err = s.cache.Send(v.w); err != nil {
					glog.Errorf("[%s] send cache packet error: %v, remove", v.w.Info(), err)
					s.ws.Remove(item.Key)
					continue
				}
				v.init = true
			} else {
				if err = v.w.Write(p); err != nil {
					glog.Errorf("[%s] write packet error: %v, remove", v.w.Info(), err)
					s.ws.Remove(item.Key)
				}
			}
		}
	}
}

func (s *Stream) TransStop() {
	if s.isStart && s.r != nil {
		s.r.Close(errors.New("stop old"))
	}
	s.isStart = false
}

func (s *Stream) CheckAlive() (n int) {
	if s.r != nil && s.isStart {
		if s.r.Alive() {
			n++
		} else {
			s.r.Close(errors.New("read timeout"))
		}
	}
	for item := range s.ws.IterBuffered() {
		v := item.Val.(*PackWriterCloser)
		if v.w != nil {
			if !v.w.Alive() {
				glog.Infof("[%v] player closed and remove\n", v.w.Info())
				s.ws.Remove(item.Key)
				v.w.Close(errors.New("write timeout"))
				continue
			}
			n++
		}
	}
	return
}

func (s *Stream) closeInter() {
	for item := range s.ws.IterBuffered() {
		v := item.Val.(*PackWriterCloser)
		if v.w != nil {
			if v.w.Info().IsInterval() {
				v.w.Close(errors.New("closed"))
				glog.Infof("[%v] player closed and remove\n", v.w.Info())
				s.ws.Remove(item.Key)
			}
		}
	}
}
