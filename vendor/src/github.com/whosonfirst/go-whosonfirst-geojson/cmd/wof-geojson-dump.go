package main

import (
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-geojson"
)

func main() {

	flag.Parse()
	args := flag.Args()

	for _, path := range args {

		f, parse_err := geojson.UnmarshalFile(path)

		if parse_err != nil {
			panic(parse_err)
		}

		fmt.Printf("# %s\n", path)

		fmt.Printf("ID is %d\n", f.Id())
		fmt.Printf("Name is %s\n", f.Name())
		fmt.Printf("Placetype is %s\n", f.Placetype())

	}

}
