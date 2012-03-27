package main

import (
	"fmt"
	"os"
	"path/filepath"
	"crypto/md5"
	"runtime"
	"flag"
//	"bufio"
//	"hash"
)

const buflen = 1<<16

type fileHash struct {
	path, hash string
	sample []byte
}

type vis struct {
	paths chan string
}

func (v *vis) VisitDir(path string, f *os.FileInfo) bool {
	return true
}

func ReadWorker(walkOut <-chan string, readOut chan<- *fileHash, readQuits chan<- bool) {
	for s := range walkOut {
		fd,_ := os.Open(s) 
		//r := bufio.NewReader(fd)
		f := fileHash{path: s, sample:make([]byte, buflen)}
		fd.Read(f.sample)
		fmt.Printf("read %s\n",s)
		fd.Close()
		readOut <- &f
	}
	readQuits <- true
}

func HashWorker(readOut <-chan *fileHash, hashOut chan<- *fileHash, hashQuits chan<- bool) {
	h := md5.New() // h is a hash.
	for r := range readOut {
		//fmt.Println("hash ",r.path)
		h.Write(r.sample)
		r.sample = nil //save space
		r.hash = fmt.Sprintf("%x", h.Sum())
		hashOut <- r
		h.Reset()
	}
	hashQuits <- true
}

func (v *vis) VisitFile(path string, f *os.FileInfo) {
	v.paths <- path
}

func AsyncWalk(path string) (<-chan string) {
	v := vis{paths: make(chan string)}
	go func () {
		filepath.Walk(path, &v, nil)
		close(v.paths)
	}()
	return v.paths
}

func main() {
	runtime.GOMAXPROCS(2)
	var nthreads = flag.Int("nthreads", 2, "number of threads")
	var nworkers = flag.Int("nworkers", 0, "number of workers")
	flag.Parse()

	if *nworkers == 0 {nworkers = nthreads}

	path := flag.Arg(0)
	walkOut := AsyncWalk(path)
	readOut := make(chan *fileHash)
	hashOut := make(chan *fileHash)
	readQuits := make(chan bool)
	hashQuits := make(chan bool)

	for i := 0; i <= *nworkers; i++ {
		go ReadWorker(walkOut, readOut, readQuits)
	}
	go func() {
		q := *nworkers
		for _ = range readQuits {
			if q == 0 {
				close(readOut)
				close(readQuits)
			}
			q--
		}
	}()

	for i := 0; i <= *nworkers; i++ {
		go HashWorker(readOut, hashOut, hashQuits)
	}
	go func() {
		q := *nworkers
		for _ = range hashQuits {
			if q == 0 {
				close(hashOut)
				close(hashQuits)
			}
			q--
		}
	}()


	m := make(map[string] []string)

	for fh := range hashOut {
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

