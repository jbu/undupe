package main

import (
	"fmt"
	"time"
	"os"
	"path/filepath"
	"crypto/md5"
	"runtime"
	"runtime/pprof"
	"log"
	"flag"
	"bufio"
)

type fileHash struct {
	path, hash string
}

type vis struct {
	c chan fileHash
	running int
	visitedAll bool
}

func (v *vis) VisitDir(path string, f *os.FileInfo) bool {
	return true
}

func doHash2(path string) (hString string) {
	h := md5.New() // h is a hash.
	fd, er := os.Open(path) 
	defer fd.Close()
	if er == nil {
		r := bufio.NewReader(fd)
		buf := make([]byte, 1024000)
		 _,ok := r.Read(buf)
		for ok == nil {
			h.Write(buf)
			_,ok = r.Read(buf)
		}
	}
	hString = fmt.Sprintf("%x", h.Sum())
	return
}

func doHash(path string, v* vis) {
	v.running++
	hString := doHash2(path)
	v.c <- fileHash{path, hString}
	v.running--
	if v.visitedAll && (v.running == 0) {
		close(v.c)
	}
}

func (v *vis) VisitFile(path string, f *os.FileInfo) {
	fmt.Print(v.running, " ")
	go doHash(path, v) 
}

func NewVis() *vis {
	v := vis{}
	v.c = make(chan fileHash)
	v.visitedAll = false
	return &v
}


func doWalk(path string, v *vis) <-chan bool {
	e := make(chan bool)
	go func () {
		filepath.Walk(path, v, nil)
		v.visitedAll = true
		e <- true
	}()
	return e
}

func main() {
	runtime.GOMAXPROCS(2)
	var cpuprofile = flag.String("cpuprof", "", "write cpu profile to file")
	var memprofile = flag.String("memprof", "", "write memory profile to this file")
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
	    		log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
    	}
	if *memprofile != "" {
		s:=fmt.Sprintf("%s/start.mprof",*memprofile)
		f, _ := os.Create(s)
		pprof.WriteHeapProfile(f)
		f.Close()
		go func() {
			i := 0
			for {
				s:=fmt.Sprintf("%s/%d.mprof",*memprofile,i)
				f, _ := os.Create(s)
				pprof.WriteHeapProfile(f)
				f.Close()
				i++
				time.Sleep(10e7)
			}
		}()
	}
	path := flag.Arg(0)
	v := NewVis()
	m := make(map[string] []string)
	e := doWalk(path, v)
	colchan := make(chan bool)
	go func () {
		for fh := range v.c {
			fmt.Println("consuming ", fh.path)
			if _, ok := m[fh.hash]; ok {
				m[fh.hash] = append(m[fh.hash], fh.path)
			} else {
				m[fh.hash] = []string{fh.path}
			}
		}
		colchan <- true
	}()
	<-e
	<-colchan
	fmt.Println("---")
	if *memprofile != "" {
		s:=fmt.Sprintf("%s/preend.mprof",*memprofile)
		f, _ := os.Create(s)
		pprof.WriteHeapProfile(f)
		f.Close()
	}
	for key, val := range(m) {
		if len(val) > 1 {
			fmt.Println(key)
			for i := range(val) {
				fmt.Println("\t",val[i])
			}
		}
	}
	if *memprofile != "" {
		s:=fmt.Sprintf("%s/end.mprof",*memprofile)
		f, _ := os.Create(s)
		pprof.WriteHeapProfile(f)
		f.Close()
	}
}

