package main

import (
       "flag"
	"fmt"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/whosonfirst/go-whosonfirst-libpostal/http"
	"github.com/whosonfirst/go-whosonfirst-log"
	gohttp "net/http"
	"os"
)

func main() {

	var host = flag.String("host", "localhost", "The hostname to listen for requests on")
	var port = flag.Int("port", 8080, "The port number to listen for requests on")

	flag.Parse()

	logger := log.SimpleWOFLogger()

	parser_handler, err := http.ParserHandler()

	if err != nil {
		logger.Fatal("failed to create parser handler, because %v", err)
	}

	expand_handler, err := http.ExpandHandler()

	if err != nil {
		logger.Fatal("failed to create expand handler, because %v", err)
	}

	ping_handler, err := http.PingHandler()

	if err != nil {
		logger.Fatal("failed to create ping handler, because %v", err)
	}

	endpoint := fmt.Sprintf("%s:%d", *host, *port)

	mux := gohttp.NewServeMux()
	
	mux.Handle("/parse", parser_handler)
	mux.Handle("/expand", expand_handler)
	mux.Handle("/ping", ping_handler)	

	err = gracehttp.Serve(&gohttp.Server{Addr: endpoint, Handler: mux})

	if err != nil {
		logger.Fatal("failed to start HTTP server, because %v", err)
	}

	os.Exit(0)
}
