package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/oschwald/geoip2-golang"
)

type Geo struct {
	reader *geoip2.Reader
}

func (g *Geo) load(file string) (err error) {
	g.reader, err = geoip2.Open(file)
	return
}

func (g *Geo) unload() (err error) {
	if g.reader != nil {
		err = g.reader.Close()
	}
	return
}

func getIP(addr string) (r string) {
	i := strings.LastIndex(addr, ":")
	if i > 0 {
		r = addr[:i]
	}
	return
}

func (g *Geo) info(addr string) (err error) {
	ip := net.ParseIP(getIP(addr))
	if ip != nil {
		if city, err := g.reader.City(ip); err == nil {
			if en, ok := city.Country.Names["en"]; ok {
				log.Println(en, city.Country.Names["zh-CN"])
				return nil
			}
		}
	}
	log.Println("unknow ip", addr)
	return err
}

var (
	geo Geo
)

func prepareGeo() error {
	err := geo.load("playground/webapi/tcloud/misc/GeoLite2-City.mmdb")
	return err
}

func unloadGeo() error {
	err := geo.unload()
	return err
}

func recordGeo(ip string) {
	geo.info(ip)
}

func VisitorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "visitor are")
}
