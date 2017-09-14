package crawl

import (
	"fmt"
	walk "github.com/whosonfirst/walk"
	"os"
)

type CrawlFunc func(path string, info os.FileInfo) error

type Crawler struct {
	Root      string
	NFSKludge bool // https://github.com/whosonfirst/walk/tree/master#walkwalkwithnfskludge
}

func NewCrawler(path string) *Crawler {
	return &Crawler{
		Root:      path,
		NFSKludge: false,
	}
}

func (c Crawler) Crawl(cb CrawlFunc) error {

	walker := func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		_ = cb(path, info)
		return nil
	}

	var err error

	// See above

	if c.NFSKludge {
		err = walk.WalkWithNFSKludge(c.Root, walker)
	} else {
		err = walk.Walk(c.Root, walker)
	}

	if err != nil {
		fmt.Printf("error: %s\n", err)
	}

	return nil
}
