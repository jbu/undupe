package main

import (
	"crypto/md5"
	//"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"regexp"

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
			fmt.Println("=1=", er)
                        continue
		}
		fi, er := fd.Stat()
		if er != nil {
			fmt.Println("=1=", er)
		}
		flen := fi.Size()
		//r := bufio.NewReader(fd)
		f := fileHash{path: s, sample: make([]byte, buflen)}
		_, er2 := fd.Read(f.sample)
		if er2 != nil && flen != 0{
			fmt.Printf("=2= %d, %d, %v, %s\n", flen, len(f.sample), er2, s)
		}
		fd.Close()
		readOut <- &f
	}
	readQuits <- true
}

func HashWorker(readOut <-chan *fileHash, hashOut chan<- *fileHash, hashQuits chan<- bool) {
	h := md5.New() // h is a hash.
	for r := range readOut {
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
        	if err != nil {
            		return err
        	}
		if info.IsDir() {
			return nil
		}
		p <- pth
        	return err
	}
	go func() {
    		err := filepath.Walk(path, f)
        	if err != nil {
            		fmt.Println("=3= ",err)
        	}
		close(p)
	}()
	return p
}

func main() {
	runtime.GOMAXPROCS(2)
	var nthreads int
	var nworkers int
	var targDirS string
	flag.IntVar(&nthreads, "nthreads", 4, "number of threads")
	flag.IntVar(&nworkers, "nworkers", 2, "number of workers")
	flag.StringVar(&targDirS, "target", "", "target directory")
	flag.Parse()

	if *nworkers == 0 {
		nworkers = nthreads
	}

 	targDirRE := regexp.MustCompile(targDirS + ".*")

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
		p := path.Clean(fh.path);
		switch {
		case len(m[fh.hash]) == 0:
			m[fh.hash] = []string{p}
		case targDirRE.MatchString(p):
			m[fh.hash] = []string{p}
		case targDirRE.MatchString(m[fh.hash][0]):
			// do nothing
		case true:
			m[fh.hash] = append(m[fh.hash], p)
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
