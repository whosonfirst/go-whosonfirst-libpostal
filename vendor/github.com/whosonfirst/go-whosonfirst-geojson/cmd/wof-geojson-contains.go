package main

import (
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-geojson"
	"log"
	"strconv"
	"strings"
)

func main() {

	var lat = flag.Float64("latitude", 0.0, "")
	var lon = flag.Float64("longitude", 0.0, "")
	var point = flag.String("point", "", "")

	flag.Parse()
	args := flag.Args()

	if *point != "" {

		parts := strings.Split(*point, ",")

		if len(parts) != 2 {
			log.Fatal("Can not parse point")
		}

		str_lat := strings.Trim(parts[0], " ")
		str_lon := strings.Trim(parts[1], " ")

		fl_lat, err := strconv.ParseFloat(str_lat, 64)

		if err != nil {
			log.Fatal(err)
		}

		fl_lon, err := strconv.ParseFloat(str_lon, 64)

		if err != nil {
			log.Fatal(err)
		}

		*lat = fl_lat
		*lon = fl_lon
	}

	for _, path := range args {

		f, parse_err := geojson.UnmarshalFile(path)

		if parse_err != nil {
			panic(parse_err)
		}

		polygons := f.GeomToPolygons()
		contains := false

		fmt.Printf("%s has this many polygons: %d\n", path, len(polygons))

		for i, poly := range polygons {

			fmt.Printf("%s #%d has %d interior rings and a total of %d points\n", path, (i + 1), len(poly.InteriorRings), poly.CountPoints())

			c := poly.Contains(*lat, *lon)

			fmt.Printf("%s #%d contains point %t\n", path, (i + 1), c)

			if c {
				contains = true
			}
		}

		fmt.Printf("%s contains point: %t\n", path, contains)

		fmt.Printf("%s f.Contains() point: %t\n", path, f.Contains(*lat, *lon))
		fmt.Println("---")
	}

}
