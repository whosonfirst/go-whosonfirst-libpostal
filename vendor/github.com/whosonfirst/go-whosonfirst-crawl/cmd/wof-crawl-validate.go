package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-crawl"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"sync/atomic"
	"time"
)

func main() {

	procs := flag.Int("processes", runtime.NumCPU()*2, "The number of concurrent processes to use")
	// nfs_kludge := flag.Bool("nfs-kludge", false, "Enable the (walk.go) NFS kludge to ignore 'readdirent: errno' 523 errors")

	flag.Parse()

	runtime.GOMAXPROCS(*procs)

	t0 := time.Now()

	var files int64
	var dirs int64

	var okay int64
	var errors int64

	for _, root := range flag.Args() {

		fmt.Println("validate JSON files in", root)

		callback := func(path string, info os.FileInfo) error {

			if info.IsDir() {
				atomic.AddInt64(&dirs, 1)
				return nil
			}

			atomic.AddInt64(&files, 1)

			fh, err := os.Open(path)

			if err != nil {
				log.Printf("failed to open %s, because %s\n", path, err)
				atomic.AddInt64(&errors, 1)
				return nil
			}

			defer fh.Close()

			body, err := ioutil.ReadAll(fh)

			if err != nil {
				log.Printf("failed to read %s, because %s\n", path, err)
				atomic.AddInt64(&errors, 1)
				return nil
			}

			var stub interface{}

			err = json.Unmarshal(body, &stub)

			if err != nil {
				log.Printf("failed to parse %s, because %s\n", path, err)
				atomic.AddInt64(&errors, 1)
				return nil
			}

			atomic.AddInt64(&okay, 1)
			return nil
		}

		c := crawl.NewCrawler(root)
		// c.NFSKludge = *nfs_kludge

		c.Crawl(callback)

	}

	t1 := float64(time.Since(t0)) / 1e9

	log.Printf("walked %d files (and %d dirs) in %.3f seconds\n", files, dirs, t1)
	log.Printf("okay %d errors %d\n", okay, errors)

}
