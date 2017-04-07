package main

/*
	"do this with that"
*/

import (
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-crawl"
	// "github.com/whosonfirst/go-whosonfirst-geojson"
	"os"
	"os/exec"
	"runtime"
	"time"
)

func main() {

	dothis := flag.String("do-this", "", "...")
	fromthere := flag.String("from-there", "", "...")
	procs := flag.Int("procs", 200, "...")
	verbose := flag.Bool("verbose", false, "...")

	flag.Parse()

	runtime.GOMAXPROCS(*procs)

	cb := func(abs_path string, info os.FileInfo) error {

		t1 := time.Now()

		cmd := exec.Command(*dothis, abs_path)
		out, err := cmd.Output()

		t2 := time.Since(t1)

		if *verbose {
			fmt.Printf("time to do this with %s: %v\n", abs_path, t2)
		}

		if err != nil {
			fmt.Printf("failed to do this with %s, because %v (%s)\n", abs_path, err, out)
			return err
		}

		if *verbose {
			fmt.Printf("%s", out)
		}

		return nil
	}

	c := crawl.NewCrawler(*fromthere)
	_ = c.Crawl(cb)

}
