package http

import (
	"encoding/json"
	"io/ioutil"
	gohttp "net/http"
	"strings"

	"github.com/whosonfirst/go-sanitize"
)

type address struct {
	Address string `json:"address"`
}

func GetAddress(r *gohttp.Request) (string, *Error) {

	method := r.Method

	if method == "GET" {
		return GetAddressForGetRequest(r)
	}

	if method == "POST" {
		return GetAddressForPostRequest(r)
	}

	return "", NewError(gohttp.StatusMethodNotAllowed, "")
}

func GetAddressForPostRequest(r *gohttp.Request) (string, *Error) {

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return "", NewError(gohttp.StatusInternalServerError, "SOMETHING_WENT_WRONG")
	}

	var address address
	err = json.Unmarshal(body, &address)

	if err != nil {
		return "", NewError(gohttp.StatusBadRequest, "E_INVALID_QUERY")
	}

	if address.Address == "" {
		return "", NewError(gohttp.StatusBadRequest, "E_INSUFFICIENT_QUERY")
	}

	return address.Address, nil
}

func GetAddressForGetRequest(r *gohttp.Request) (string, *Error) {

	query, err := ParseRequest(r)

	if err != nil {
		return "", err
	}

	q, err := ParseQuery(query)

	if err != nil {
		return "", err
	}

	return q, nil
}

func ParseRequest(r *gohttp.Request) (*Query, *Error) {

	query := r.URL.Query()

	address := strings.Trim(query.Get("address"), " ")

	if address == "" {
		return nil, NewError(gohttp.StatusBadRequest, "E_INSUFFICIENT_QUERY")
	}

	q := Query{
		Address: address,
	}

	return &q, nil
}

func ParseQuery(query *Query) (string, *Error) {

	opts := sanitize.DefaultOptions()
	q, err := sanitize.SanitizeString(query.Address, opts)

	if err != nil {
		return "", NewError(gohttp.StatusBadRequest, err.Error())
	}

	if q == "" {
		return "", NewError(gohttp.StatusBadRequest, "E_INVALID_QUERY")
	}

	return q, nil
}

func WriteResponse(w gohttp.ResponseWriter, rsp interface{}) {

	rsp_encoded, err := json.Marshal(rsp)

	if err != nil {
		gohttp.Error(w, err.Error(), gohttp.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Write(rsp_encoded)
}
