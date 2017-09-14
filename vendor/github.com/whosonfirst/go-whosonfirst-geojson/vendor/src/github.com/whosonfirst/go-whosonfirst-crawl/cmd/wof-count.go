package main

import (
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-crawl"
	"os"
	"runtime"
	"time"
)

func main() {

	procs := flag.Int("processes", runtime.NumCPU()*2, "The number of concurrent processes to use")
	nfs_kludge := flag.Bool("nfs-kludge", false, "Enable the (walk.go) NFS kludge to ignore 'readdirent: errno' 523 errors")

	flag.Parse()
	args := flag.Args()

	runtime.GOMAXPROCS(*procs)

	root := args[0]
	fmt.Println("count files and directories in ", root)

	var files int64
	var dirs int64

	callback := func(path string, info os.FileInfo) error {

		if info.IsDir() {
			dirs++
			return nil
		}

		files++
		return nil
	}

	t0 := time.Now()

	c := crawl.NewCrawler(root)
	c.NFSKludge = *nfs_kludge

	_ = c.Crawl(callback)

	t1 := float64(time.Since(t0)) / 1e9
	fmt.Printf("walked %d files (and %d dirs) in %.3f seconds\n", files, dirs, t1)
}
