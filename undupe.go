package main

import (
	"crypto/md5"
	//"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

//	"bufio"
//	"hash"
)

const buflen = 1 << 16

type fileHash struct {
	path, hash string
	sample     []byte
}

func ReadWorker(walkOut <-chan string, readOut chan<- *fileHash, readQuits chan<- bool) {
	for s := range walkOut {
		fd, er := os.Open(s)
		if er != nil {
			fmt.Println("===", er)
		}
		//r := bufio.NewReader(fd)
		f := fileHash{path: s, sample: make([]byte, buflen)}
		_, er = fd.Read(f.sample)
		//fmt.Printf("read %s\n",s)
		if er != nil {
			fmt.Println("===", er)
		}
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
		r.hash = fmt.Sprintf("%x", h.Sum(nil))
		hashOut <- r
		h.Reset()
	}
	hashQuits <- true
}

func AsyncWalk(path string) <-chan string {
	p := make(chan string)
	f := func(pth string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
        	if err != nil {
            		return err
        	}
		p <- pth
        	return nil
	}
	go func() {
    		err := filepath.Walk(path, f)
        	if err != nil {
            		fmt.Println("=== ",err)
        	}
		close(p)
	}()
	return p
}

func main() {
	runtime.GOMAXPROCS(2)
	var nthreads = flag.Int("nthreads", 2, "number of threads")
	var nworkers = flag.Int("nworkers", 0, "number of workers")
	flag.Parse()

	if *nworkers == 0 {
		nworkers = nthreads
	}

	path := flag.Arg(0)
	walkOut := AsyncWalk(path)
	readOut := make(chan *fileHash, 10)
	hashOut := make(chan *fileHash, 10)
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

	m := make(map[string][]string)

	for fh := range hashOut {
		if _, ok := m[fh.hash]; ok {
			m[fh.hash] = append(m[fh.hash], fh.path)
		} else {
			m[fh.hash] = []string{fh.path}
		}
	}

	for _, val := range m {
		if len(val) > 1 {
			//fmt.Println(key)
			for i := range val {
				fmt.Print(val[i], "\t")
			}
			fmt.Println()
		}
	}
}
