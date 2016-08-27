package main

import (
	"encoding/json"
	"expvar"
	"flag"
	"fmt"
	"github.com/facebookgo/grace/gracehttp"
	expand "github.com/openvenues/gopostal/expand"
	parser "github.com/openvenues/gopostal/parser"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"syscall"
	"time"
)

var (
	readErrors      *expvar.Int
	marshalErrors   *expvar.Int
	unmarshalErrors *expvar.Int

	parseRequests *expvar.Int
	parseSuccess  *expvar.Int
	parseAvg      *expvar.Float

	expandRequests *expvar.Int
	expandSuccess  *expvar.Int
	expandAvg      *expvar.Float

	timeToParse  int64
	timeToExpand int64
)

func init() {

	readErrors = expvar.NewInt("ErrRead")
	marshalErrors = expvar.NewInt("ErrMarshal")
	unmarshalErrors = expvar.NewInt("ErrUnmarshal")

	parseRequests = expvar.NewInt("ReqParse")
	parseSuccess = expvar.NewInt("SuccessParse")
	parseAvg = expvar.NewFloat("AvgParse")

	expandRequests = expvar.NewInt("ReqExpand")
	expandSuccess = expvar.NewInt("SuccessExpand")
	expandAvg = expvar.NewFloat("AvgExpand")

	timeToParse = 0
	timeToExpand = 0
}

type Request struct {
	Query string `json:"query"`
}

// because this: https://github.com/golang/go/issues/15030
// but anyway we want to restict access to localhost so oh well...
// (20160826/thisisaaronland)

func ExpvarHandlerFunc(host string) http.HandlerFunc {

	f := func(w http.ResponseWriter, r *http.Request) {

		remote := strings.Split(r.RemoteAddr, ":")

		if remote[0] != host {

			http.Error(w, "No soup for you!", http.StatusForbidden)
			return
		}

		// This is copied wholesale from
		// https://golang.org/src/expvar/expvar.go

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		fmt.Fprintf(w, "{\n")

		first := true

		expvar.Do(func(kv expvar.KeyValue) {
			if !first {
				fmt.Fprintf(w, ",\n")
			}

			first = false
			fmt.Fprintf(w, "%q: %s", kv.Key, kv.Value)
		})

		fmt.Fprintf(w, "\n}\n")
	}

	return http.HandlerFunc(f)
}

func ExpandHandler(w http.ResponseWriter, r *http.Request) {

	req, err := ParseRequest(r, expandRequests)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t1 := time.Now()
	expansions := expand.ExpandAddress(req.Query)
	t2 := time.Since(t1)

	go func(t time.Duration) {

		ns := t.Nanoseconds()
		ms := ns / (int64(time.Millisecond) / int64(time.Nanosecond))

		tte := atomic.AddInt64(&timeToExpand, ms)

		req, err := strconv.ParseFloat(expandRequests.String(), 64)

		if err != nil {
			return
		}

		avg := float64(tte) / req
		expandAvg.Set(avg)

	}(t2)

	WriteResponse(w, expansions, expandSuccess)
}

func ParserHandler(w http.ResponseWriter, r *http.Request) {

	req, err := ParseRequest(r, parseRequests)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t1 := time.Now()
	parsed := parser.ParseAddress(req.Query)
	t2 := time.Since(t1)

	go func(t time.Duration) {

		ns := t.Nanoseconds()
		ms := ns / (int64(time.Millisecond) / int64(time.Nanosecond))

		ttp := atomic.AddInt64(&timeToParse, ms)

		req, err := strconv.ParseFloat(parseRequests.String(), 64)

		if err != nil {
			return
		}

		avg := float64(ttp) / req
		parseAvg.Set(avg)

	}(t2)

	WriteResponse(w, parsed, parseSuccess)
}

func ParseRequest(r *http.Request, requests *expvar.Int) (*Request, error) {

	requests.Add(1)

	q, err := ioutil.ReadAll(r.Body)

	if err != nil {
		readErrors.Add(1)
		return nil, err
	}

	var req Request
	err = json.Unmarshal(q, &req)

	if err != nil {
		unmarshalErrors.Add(1)
		return nil, err
	}

	return &req, nil
}

func WriteResponse(w http.ResponseWriter, rsp interface{}, successes *expvar.Int) {

	rsp_encoded, err := json.Marshal(rsp)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Write(rsp_encoded)
	successes.Add(1)
}

func main() {

	var host = flag.String("host", "localhost", "The hostname to listen for requests on")
	var port = flag.Int("port", 8080, "The port number to listen for requests on")
	var pidfile = flag.String("pidfile", "", "Where to write a PID file for wof-libpostal-server. If empty the PID file will be written to wof-libpostal-server.pid in the current directory")

	flag.Parse()

	go func() {

		if *pidfile == "" {

			cwd, err := os.Getwd()

			if err != nil {
				panic(err)
			}

			fname := fmt.Sprintf("%s.pid", os.Args[0])

			*pidfile = filepath.Join(cwd, fname)
		}

		fh, err := os.Create(*pidfile)

		if err != nil {
			panic(err)
		}

		defer fh.Close()

		pid := os.Getpid()
		strpid := strconv.Itoa(pid)

		fh.Write([]byte(strpid))

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			<-sigs

			os.Remove(*pidfile)
			os.Exit(0)
		}()

	}()

	endpoint := fmt.Sprintf("%s:%d", *host, *port)

	mux := http.NewServeMux()
	mux.HandleFunc("/parse", ParserHandler)
	mux.HandleFunc("/parser", ParserHandler) // this is legacy from testing and should be considered DEPRICATED
	mux.HandleFunc("/expand", ExpandHandler)

	// see comments above
	expvarHandler := ExpvarHandlerFunc(*host)
	mux.HandleFunc("/debug/vars", expvarHandler)

	err := gracehttp.Serve(&http.Server{Addr: endpoint, Handler: mux})

	if err != nil {
		log.Fatal(err)
	}

	os.Remove(*pidfile)
	os.Exit(0)
}
