package http

import (
	"github.com/openvenues/gopostal/parser"
	gohttp "net/http"
)

func ExpandHandler() (gohttp.Handler, error) {

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		address, err := GetAddress(req)

		if err != nil {
			gohttp.Error(w, err.Error(), err.Code)
			return
		}

		parsed := parser.ParseAddress(address)

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
