package main

import (
	"fmt"
	"os"
	"path/filepath"
	"crypto/md5"
	"runtime"
	"flag"
	"bufio"
	"hash"
)

type fileHash struct {
	path, hash string
}

type vis struct {
	req, resp chan *fileHash
	quits chan bool
}

func (v *vis) VisitDir(path string, f *os.FileInfo) bool {
	return true
}

const buflen = 1024000

func doHashTick(w *hash.Hash, bufs *[2][buflen]byte, ticks <-chan int) {
	for idx := range(ticks) {
		//fmt.Printf("hbuf %v\n",bufs[idx])
		//fmt.Print(".")
		(*w).Write(bufs[idx][:])
	}
}

func HashWorker(req <-chan *fileHash, resp chan<- *fileHash, quits chan<- bool) {
	h := md5.New() // h is a hash.
	var bufs [2][buflen]byte
	//b := &buf
	idx := 0
	for r := range req {
		//fmt.Println("<-req ",r.path)
		fmt.Print(",")
		fd, er := os.Open(r.path) 
		ticks := make(chan int)
		go doHashTick(&h, &bufs, ticks)
		if er == nil {
			r := bufio.NewReader(fd)
		 	_,ok := r.Read(bufs[idx][:])
			for ok == nil {
				//fmt.Printf("buf %v, idx %v\n",bufs, idx)
				//h.Write(bufs[idx][:])
				ticks <- idx
				if idx == 0 {idx = 1} else {idx = 0}
				_,ok = r.Read(bufs[idx][:])
			}
		}
		close(ticks)
		fd.Close()
		r.hash = fmt.Sprintf("%x", h.Sum())
		fmt.Println(*r)
		resp <- r
		h.Reset()
	}
	quits <- true
}

func (v *vis) VisitFile(path string, f *os.FileInfo) {
//	fmt.Println(".",path)
	//fmt.Print(".")
	v.req <- &fileHash{path:path}
}

func NewVis(path string, nworkers int) *vis {
	v := vis{}
	v.req = make(chan *fileHash, 5)
	v.resp = make(chan *fileHash, 5)
	v.quits = make(chan bool)
	for i := 0 ; i < nworkers ; i++ {
		fmt.Println("starting workers: ", i)
		go HashWorker(v.req, v.resp, v.quits)
	}
	go func () {
		filepath.Walk(path, &v, nil)
		close(v.req)
	}()
	go func () {
		i := nworkers
		for _ = range v.quits {
			i--
			if i == 0 {
				close(v.quits)
				close(v.resp)
			}
		}
	}()
	return &v
}

func main() {
	runtime.GOMAXPROCS(2)
	var nthreads = flag.Int("nthreads", 2, "number of threads")
	var nworkers = flag.Int("nworkers", 0, "number of workers")
	flag.Parse()

	if *nworkers == 0 {nworkers = nthreads}

	path := flag.Arg(0)
	v := NewVis(path, *nworkers)


	m := make(map[string] []string)

	for fh := range v.resp {
		if _, ok := m[fh.hash]; ok {
			m[fh.hash] = append(m[fh.hash], fh.path)
		} else {
			m[fh.hash] = []string{fh.path}
		}
	}

	fmt.Println("---")
	for key, val := range(m) {
		if len(val) > 1 {
			fmt.Println(key)
			for i := range(val) {
				fmt.Println("\t",val[i])
			}
		}
	}
}

