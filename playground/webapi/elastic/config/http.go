package config

type HTTPConfig struct {
	LitenAddress string
}

var http = &HTTPConfig{
	LitenAddress: "localhost:8000",
}

func HTTP() *HTTPConfig {
	return http
}
