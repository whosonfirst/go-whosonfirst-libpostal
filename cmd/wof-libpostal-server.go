package main

import (
	"fmt"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/whosonfirst/go-whosonfirst-libpostal/http"
	"github.com/whosonfirst/go-whosonfirst-log"
	"os"
)

func main() {

	var host = flag.String("host", "localhost", "The hostname to listen for requests on")
	var port = flag.Int("port", 8080, "The port number to listen for requests on")

	flag.Parse()

	parser_hander, err := http.ParserHandler()

	if err != nil {
		logger.Fatal("failed to create parser handler, because %v", err)
	}

	expand_hander, err := http.ExpandHandler()

	if err != nil {
		logger.Fatal("failed to create expand handler, because %v", err)
	}

	endpoint := fmt.Sprintf("%s:%d", *host, *port)

	mux := http.NewServeMux()
	mux.HandleFunc("/parse", parser_handler)
	mux.HandleFunc("/expand", expand_handler)

	err := gracehttp.Serve(&http.Server{Addr: endpoint, Handler: mux})

	if err != nil {
		logger.Fatal("failed to start HTTP server, because %v", err)
	}

	os.Exit(0)
}
