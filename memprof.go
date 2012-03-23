package mprof

import (
	"fmt"
	"time"
	"os"
	"runtime/pprof"
)

type MemProf struct {
	path string
	cnt int
	stop chan bool
}

func MakeMemProf(path string) (m *MemProf) {
	_ = os.MkdirAll(path)
	return &MemProf{dir, 0, make(chan bool)}
}

func (m *MemProf) startDumps(interval int64) {
	go func() {
		i := 0
		for {
			s:=fmt.Sprintf("%s/%d.mprof",m.path,i)
			f, _ := os.Create(s)
			pprof.WriteHeapProfile(f)
			f.Close()
			i++
			select {
				case <- m.stop: break
				case <- time.After(interval):
			}
		}
	}()
}

func (m *MemProf) stopDumps() {
	m.stop <- true;
}

func (m *MemProf) dump(label string) {
	s:=fmt.Sprintf("%s/%s.mprof",m.path, label)
	f, _ := os.Create(s)
	pprof.WriteHeapProfile(f)
	f.Close()
}

