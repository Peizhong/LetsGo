package main

import (
	"encoding/json"
	"log"

	"go.uber.org/dig"
)

type Config struct {
	Prefix string
}

func NewReadOnlyConnection() (*Config, error) {
	var cfg Config
	err := json.Unmarshal([]byte(`{"prefix": "[read only] "}`), &cfg)
	return &cfg, err
}

func NewReadWriteConnection() (*Config, error) {
	var cfg Config
	err := json.Unmarshal([]byte(`{"prefix": "[read write] "}`), &cfg)
	return &cfg, err

}

type GatewayParams struct {
	dig.In

	WriteToConn   *Config `name:"rw"`
	ReadFromConn  *Config `name:"ro"`
	ReadFromConnx *Config `name:"rox"`
}

type GatewayParamsOut struct {
	dig.Out

	WriteToConn   *Config `name:"rw"`
	ReadFromConn  *Config `name:"ro"`
	ReadFromConnx *Config `name:"rox"`
}

func NewReadWriteOut() (GatewayParamsOut, error) {
	var cfg Config
	err := json.Unmarshal([]byte(`{"prefix": "[read write] "}`), &cfg)
	return GatewayParamsOut{WriteToConn: &cfg, ReadFromConn: &cfg}, err
}

func main() {
	type mm struct {
		cfg *Config
	}
	c := dig.New()
	err := c.Provide(NewReadOnlyConnection, dig.Name("roxx"))
	//err = c.Provide(NewReadWriteConnection, dig.Name("rw"))
	err = c.Provide(NewReadWriteOut)
	err = c.Invoke(func(rw GatewayParams) {
		log.Println(rw.WriteToConn.Prefix, rw.ReadFromConn.Prefix, rw.ReadFromConnx.Prefix)
	})
	if err != nil {
		panic(err)
	}
}
