package main

import (
	"encoding/json"
	"expvar"
	"flag"
	"fmt"
	"github.com/facebookgo/grace/gracehttp"
	expand "github.com/openvenues/gopostal/expand"
	parser "github.com/openvenues/gopostal/parser"
	sanitize "github.com/whosonfirst/go-sanitize"
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
	readErrors    *expvar.Int
	marshalErrors *expvar.Int
	inputErrors   *expvar.Int

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
	inputErrors = expvar.NewInt("ErrInput")

	parseRequests = expvar.NewInt("ReqParse")
	parseSuccess = expvar.NewInt("SuccessParse")
	parseAvg = expvar.NewFloat("AvgParse")

	expandRequests = expvar.NewInt("ReqExpand")
	expandSuccess = expvar.NewInt("SuccessExpand")
	expandAvg = expvar.NewFloat("AvgExpand")

	timeToParse = 0
	timeToExpand = 0
}

type Query struct {
	Address string
}

type HTTPError struct {
	error
	Code    int
	Message string
}

func (e HTTPError) Error() string {
	return e.Message
}

func NewHTTPError(code int, message string) *HTTPError {

	e := HTTPError{
		Code:    code,
		Message: message,
	}

	return &e
}

// because this: https://github.com/golang/go/issues/15030
// but anyway we want to restict access to localhost so oh well...
// (20160826/thisisaaronland)

func ExpvarHandlerFunc(host string) http.HandlerFunc {

	f := func(w http.ResponseWriter, r *http.Request) {

		remote := strings.Split(r.RemoteAddr, ":")

		if remote[0] != "127.0.0.1" && remote[0] != host {

			http.Error(w, "No soup for you!", http.StatusForbidden)
			return
		}

		_, err := IsValidMethod(r, []string{"GET"})

		if err != nil {
			http.Error(w, err.Error(), err.Code)
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

	address, err := GetAddress(r)

	if err != nil {
		http.Error(w, err.Error(), err.Code)
		return
	}

	t1 := time.Now()
	expansions := expand.ExpandAddress(address)
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

	address, err := GetAddress(r)

	if err != nil {
		http.Error(w, err.Error(), err.Code)
		return
	}

	t1 := time.Now()
	parsed := parser.ParseAddress(address)
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

	query := r.URL.Query()

	if query.Get("format") == "keys" {
		WriteResponse(w, FormatParsed(parsed), parseSuccess)
	} else {
		WriteResponse(w, parsed, parseSuccess)
	}
}

func FormatParsed(parsed []parser.ParsedComponent) map[string][]string {

	rsp := make(map[string][]string)

	for _, component := range parsed {

		key := component.Label
		value := component.Value

		possible, ok := rsp[key]

		if !ok {
			possible = make([]string, 0)
		}

		possible = append(possible, value)
		rsp[key] = possible
	}

	return rsp
}

func GetAddress(r *http.Request) (string, *HTTPError) {

	_, err := IsValidMethod(r, []string{"GET"})

	if err != nil {
		return "", err
	}

	query, err := ParseRequest(r, parseRequests)

	if err != nil {
		// log.Printf("parse error %s (%d)\n", err.Error(), err.Code)
		return "", err
	}

	q, err := ParseQuery(query)

	if err != nil {
		// log.Printf("query error %s (%d)\n", err.Error(), err.Code)
		return "", err
	}

	return q, nil
}

func IsValidMethod(r *http.Request, allowed []string) (bool, *HTTPError) {

	method := r.Method
	ok := false

	for _, this := range allowed {

		if this == method {
			ok = true
			break
		}
	}

	if !ok {
		readErrors.Add(1)
		return false, NewHTTPError(http.StatusMethodNotAllowed, "")
	}

	return true, nil
}

func ParseRequest(r *http.Request, requests *expvar.Int) (*Query, *HTTPError) {

	requests.Add(1)

	query := r.URL.Query()

	address := strings.Trim(query.Get("address"), " ")

	if address == "" {
		inputErrors.Add(1)
		return nil, NewHTTPError(http.StatusBadRequest, "E_INSUFFICIENT_QUERY")
	}

	q := Query{
		Address: address,
	}

	return &q, nil
}

func ParseQuery(query *Query) (string, *HTTPError) {

	opts := sanitize.DefaultOptions()
	q, err := sanitize.SanitizeString(query.Address, opts)

	if err != nil {
		inputErrors.Add(1)
		return "", NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if q == "" {
		inputErrors.Add(1)
		return "", NewHTTPError(http.StatusBadRequest, "E_INVALID_QUERY")
	}

	return q, nil
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

		if *pidfile == "-" {
			return
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
