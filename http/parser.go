package http

import (
	"github.com/openvenues/gopostal/parser"
	gohttp "net/http"
)

func ParserHandler() (gohttp.Handler, error) {

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		address, err := GetAddress(req)

		if err != nil {
			gohttp.Error(rsp, err.Error(), err.Code)
			return
		}

		parsed := postal.ParseAddress(address)

		query := req.URL.Query()

		if query.Get("format") == "keys" {
			WriteResponse(rsp, FormatParsed(parsed))
		} else {
			WriteResponse(rsp, parsed)
		}
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}

func FormatParsed(parsed []postal.ParsedComponent) map[string][]string {

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
