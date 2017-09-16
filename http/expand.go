package http

import (
	"github.com/openvenues/gopostal/expand"
	gohttp "net/http"
)

func ExpandHandler() (gohttp.Handler, error) {

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		address, err := GetAddress(req)

		if err != nil {
			gohttp.Error(rsp, err.Error(), err.Code)
			return
		}

		expansions := postal.ExpandAddress(address)
		WriteResponse(rsp, expansions)
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
