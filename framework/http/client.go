package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	_http "net/http"
	"strings"

	"github.com/peizhong/letsgo/framework/log"
)

type Header struct {
	K, V string
}

type HTTPResponse struct {
	Body    []byte
	Headers []Header
}

func (r HTTPResponse) String() string {
	return string(r.Body)
}

func (r HTTPResponse) Fill(obj interface{}) {
	json.Unmarshal(r.Body, obj)
}

func Do(method, url string, headers []Header, body string, query ...interface{}) (*HTTPResponse, error) {
	req := strings.Builder{}
	req.WriteString(url)
	ql := len(query)
	qc := ql / 2
	if qc > 0 {
		qc := qc * 2
		for i := 0; i < qc; i += 2 {
			if i == 0 {
				req.WriteString(fmt.Sprintf("?%v=%v", query[i], query[i+1]))
			} else {
				req.WriteString(fmt.Sprintf("&%v=%v", query[i], query[i+1]))
			}
		}
	}
	reqURL := req.String()
	log.Info("requset url: %v", reqURL)
	client := _http.Client{}
	r, err := _http.NewRequest(method, reqURL, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}
	for _, h := range headers {
		r.Header.Set(h.K, h.V)
	}
	response, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	rHeaders := make([]Header, len(response.Header))
	var i int
	for k, v := range response.Header {
		rHeaders[i] = Header{k, v[0]}
		i++
	}
	return &HTTPResponse{
		Body:    bytes,
		Headers: rHeaders,
	}, nil
}
