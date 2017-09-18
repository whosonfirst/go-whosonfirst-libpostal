package http

import (
	"github.com/openvenues/gopostal/expand"
	"github.com/whosonfirst/go-whosonfirst-log"
	gohttp "net/http"
	"time"
)

func ExpandHandler(logger *log.WOFLogger) (gohttp.Handler, error) {

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		address, err := GetAddress(req)

		if err != nil {
			gohttp.Error(rsp, err.Error(), err.Code)
			return
		}

		t1 := time.Now()

		defer func() {
			t2 := time.Since(t1)
			logger.Status("expand '%s' %v", address, t2)
		}()

		expansions := postal.ExpandAddress(address)
		WriteResponse(rsp, expansions)
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
