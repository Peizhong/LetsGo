package http

import (
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
	Body []byte
}

func (r HTTPResponse) String() string {
	return string(r.Body)
}

func (r HTTPResponse) Fill(obj interface{}) {
	json.Unmarshal(r.Body, obj)
}

func Get(url string, headers []Header, query ...interface{}) (*HTTPResponse, error) {
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
	r, err := _http.NewRequest("GET", reqURL, nil)
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
	return &HTTPResponse{
		Body: bytes,
	}, nil
}
