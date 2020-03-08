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
	g.reader, err = geoip2.Open("misc/GeoLite2-City.mmdb")
	return
}

func (g *Geo) unload() (err error) {
	if g.reader != nil {
		err = g.reader.Close()
	}
	return
}

func (g *Geo) info(addr string) error {
	i := strings.LastIndex(addr, ":")
	if i > 0 {
		addr = addr[:i]
	}
	ip := net.ParseIP("182.61.200.6")
	if ip != nil {
		if city, err := g.reader.City(ip); err == nil {
			log.Println(city.Country.Names["en"], city.Country.Names["zh-CN"])
		} else {
			return err
		}
	}
	return nil
}

var (
	geo Geo
)

func prepareGeo() error {
	err := geo.load("misc/GeoLite2-City.mmdb")
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
