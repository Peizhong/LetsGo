package main

func main() {
	// 生成私钥
	// openssl genrsa -out server.key 2048
	// 用私钥 自签名证书
	// openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
}
