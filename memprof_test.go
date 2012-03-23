package mprof

import (
	"testing"
)


func TestCreate(t *testing.T) {
	m := MakeMemProf("tmprof")
}

func TestLabel(t *testing.T) {
	m := MakeMemProf("tmprof")
}

type MemProf struct {
	dir string
	cnt int
	stop chan bool
}

func (m *MemProf) startDumps(interval int64) {
	go func() {
		i := 0
		for {
			s:=fmt.Sprintf("%s/%d.mprof",m.dir,i)
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
	s:=fmt.Sprintf("%s/%s.mprof",m.dir, label)
	f, _ := os.Create(s)
	pprof.WriteHeapProfile(f)
	f.Close()
}

