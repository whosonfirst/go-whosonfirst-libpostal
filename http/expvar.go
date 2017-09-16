package http

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

// because this: https://github.com/golang/go/issues/15030
// but anyway we want to restict access to localhost so oh well...
// (20160826/thisisaaronland)

func ExpvarHandler(host string) (gohttp.Handler, error) {

	f := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		remote := strings.Split(req.RemoteAddr, ":")

		if remote[0] != "127.0.0.1" && remote[0] != host {

			gohttp.Error(rsp, "No soup for you!", gohttp.StatusForbidden)
			return
		}

		/*
			_, err := IsValidMethod(req, []string{"GET"})

			if err != nil {
				gohttp.Error(rsp, err.Error(), err.Code)
			}
		*/

		// This is copied wholesale from
		// https://golang.org/src/expvar/expvar.go

		rsp.Header().Set("Content-Type", "application/json; charset=utf-8")
		fmt.Fprintf(rsp, "{\n")

		first := true

		expvar.Do(func(kv expvar.KeyValue) {

			if !first {
				fmt.Fprintf(w, ",\n")
			}

			first = false
			fmt.Fprintf(rsp, "%q: %s", kv.Key, kv.Value)
		})

		fmt.Fprintf(rsp, "\n}\n")
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
