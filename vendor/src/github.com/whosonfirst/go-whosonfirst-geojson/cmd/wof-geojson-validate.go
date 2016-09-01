package main

import (
	"flag"
	"fmt"
	crawl "github.com/whosonfirst/go-whosonfirst-crawl"
	geojson "github.com/whosonfirst/go-whosonfirst-geojson"
	"os"
	"runtime"
	"sync"
	"time"
)

func main() {

	var source = flag.String("source", "", "Where to look for files")
	var procs = flag.Int("processes", runtime.NumCPU()*2, "Number of concurrent processes to use")

	flag.Parse()

	runtime.GOMAXPROCS(*procs)

	// TO DO: ensure *source exists and is a directory

	t1 := time.Now()
	count := 0

	wg := new(sync.WaitGroup)

	callback := func(source string, info os.FileInfo) error {

		wg.Add(1)
		defer wg.Done()

		if info.IsDir() {
			return nil
		}

		count += 1

		_, err := geojson.UnmarshalFile(source)

		if err != nil {
			fmt.Println(source, err)
		}

		return nil
	}

	c := crawl.NewCrawler(*source)
	_ = c.Crawl(callback)

	wg.Wait()

	t2 := time.Since(t1)

	fmt.Printf("time to validate %d files: %v\n", count, t2)
}
