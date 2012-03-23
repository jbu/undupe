package mprof

import (
	"fmt"
	"time"
	"os"
	"runtime/pprof"
)

type Mprof struct {
	path string
	cnt int
	stop chan bool
}

func NewMprof(path string) (m *Mprof) {
	if path == "" {return nil}
	_ = os.MkdirAll(path, 0666)
	return &Mprof{path, 0, make(chan bool)}
}

func (m *Mprof) Dump(label string) {
	if m == nil {return}
	s:=fmt.Sprintf("%s/%s.mprof",m.path, label)
	f, _ := os.Create(s)
	pprof.WriteHeapProfile(f)
	f.Close()
}

func (m *Mprof) StartDumps(interval int64) {
	if m == nil {return}
	go func() {
		i := 0
		for {
			s:=fmt.Sprintf("%d",i)
			m.Dump(s)
			i++
			select {
				case <- m.stop: break
				case <- time.After(interval):
			}
		}
	}()
}

func (m *Mprof) StopDumps() {
	if m == nil {return}
	m.stop <- true;
}

