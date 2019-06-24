module github.com/peizhong/letsgo

go 1.12

require (
	github.com/coreos/etcd v3.3.13+incompatible
	github.com/davecgh/go-spew v1.1.1
	github.com/emirpasic/gods v1.12.0
	github.com/go-ini/ini v1.42.0
	github.com/golang/mock v1.3.1
	github.com/golang/protobuf v1.3.1
	github.com/google/uuid v1.1.1
	github.com/gorilla/mux v1.7.2
	github.com/jinzhu/gorm v1.9.9
	github.com/jinzhu/now v1.0.1
	github.com/micro/go-micro v1.7.0
	github.com/micro/mdns v0.1.0
	github.com/micro/micro v1.7.0
	github.com/mitchellh/hashstructure v1.0.0
	github.com/satori/go.uuid v1.2.0
	github.com/stretchr/testify v1.3.0
	github.com/tidwall/evio v1.0.2
	github.com/tidwall/gjson v1.2.1
	github.com/tidwall/match v1.0.1 // indirect
	github.com/tidwall/pretty v1.0.0 // indirect
	golang.org/x/crypto v0.0.0-20190621222207-cc06ce4a13d4 // indirect
	golang.org/x/sys v0.0.0-20190624142023-c5567b49c5d0 // indirect
	golang.org/x/tools v0.0.0-20190621195816-6e04913cbbac // indirect
	google.golang.org/grpc v1.21.1

)

replace github.com/docker/docker v0.7.3-0.20190309235953-33c3200e0d16 => github.com/docker/docker v1.13.1
