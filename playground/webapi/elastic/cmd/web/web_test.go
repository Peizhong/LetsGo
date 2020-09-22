package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestBinding(t *testing.T) {
	engine := buildHandler()
	s := httptest.NewServer(engine)
	defer s.Close()

	resp, err := s.Client().Get(s.URL + "/api/app/123?key=apple")
	if !assert.NoError(t, err) {
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	t.Log(string(body))

	data := PostDemoData{
		Id:    "key",
		Value: "value",
	}
	byts, _ := json.Marshal(data)
	resp, err = s.Client().Post(s.URL+"/api/app/123?key=apple", "application/json", bytes.NewReader(byts))
	if !assert.NoError(t, err) {
		return
	}
	body, _ = ioutil.ReadAll(resp.Body)
	t.Log(string(body))
}
