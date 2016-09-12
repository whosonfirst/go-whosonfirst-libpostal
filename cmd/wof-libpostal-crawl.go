package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-crawl"
	"github.com/whosonfirst/go-whosonfirst-csv"
	"github.com/whosonfirst/go-whosonfirst-geojson"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type LibpostalResponse []LibpostalElement

type LibpostalElement struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

func parse(endpoint string, f *geojson.WOFFeature, properties []string) (time.Duration, string, string) {

	addr := make([]string, 0)

	for _, k := range properties {

		k = fmt.Sprintf("properties.%s", k)

		v, ok := f.StringValue(k)

		if ok {
			addr = append(addr, v)
		}
	}

	str_addr := strings.Join(addr, " ")

	url := fmt.Sprintf("%s/parse", endpoint)

	t1 := time.Now()

	req, err := http.NewRequest("GET", url, nil)

	q := req.URL.Query()
	q.Add("address", str_addr)

	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	t2 := time.Since(t1)

	// fmt.Printf("time to query %s, %v\n", str_addr, t2)

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	r := LibpostalResponse{}
	err = json.Unmarshal(body, &r)

	if err != nil {
		panic(err)
	}

	return t2, str_addr, string(body)
}

func main() {

	var host = flag.String("libpostal-host", "", "The host for the libpostal endpoint")
	var port = flag.Int("libpostal-port", 8080, "The host for the libpostal port")
	var out = flag.String("output", "libpostal.csv", "Where to write output data")
	var props = flag.String("properties", "sg:address,sg:city,sg:province,sg:postcode", "A comma-separated list of GeoJSON properties used to construct an address")
	var processes = flag.Int("processes", (runtime.NumCPU() * 2), "The number of concurrent processes to clone data with")

	flag.Parse()
	args := flag.Args()

	runtime.GOMAXPROCS(*processes)

	var ttq int64
	var files int64

	endpoint := fmt.Sprintf("http://%s:%d", *host, *port)

	properties := make([]string, 0)

	for _, k := range strings.Split(*props, ",") {

		k = strings.Trim(k, " ")

		if k != "" {
			properties = append(properties, k)
		}
	}

	if len(properties) == 0 {
		log.Fatal("Invalid properties list")
	}

	fieldnames := []string{"id", "address", "results"}

	writer, err := csv.NewDictWriterFromPath(*out, fieldnames)

	if err != nil {
		log.Fatal(err)
	}

	writer.WriteHeader()

	mu := new(sync.Mutex)
	wmu := new(sync.Mutex)

	callback := func(path string, info os.FileInfo) error {

		if info.IsDir() {
			return nil
		}

		feature, err := geojson.UnmarshalFile(path)

		if err != nil {
			fmt.Println(path)
			panic(err)
		}

		t, address, results := parse(endpoint, feature, properties)

		id := feature.Id()

		go func(id int, address string, results string) {

			str_id := strconv.Itoa(id)

			row := make(map[string]string)
			row["id"] = str_id
			row["address"] = address
			row["results"] = results

			wmu.Lock()
			writer.WriteRow(row)
			wmu.Unlock()

		}(id, address, results)

		ns := t.Nanoseconds()
		ms := ns / int64(time.Millisecond)

		mu.Lock()
		files += 1
		ttq += ms
		mu.Unlock()

		return nil
	}

	for _, root := range args {

		t1 := time.Now()
		c := crawl.NewCrawler(root)

		_ = c.Crawl(callback)
		t2 := float64(time.Since(t1)) / 1e9

		avg := float64(ttq) / float64(files)

		fmt.Printf("parsed %d files in %s in %.3f seconds avg ttq: %.f6 ms\n", files, root, t2, avg)
	}
}
