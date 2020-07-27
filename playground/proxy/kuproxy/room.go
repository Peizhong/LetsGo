package main

//go:generate mockgen -destination ./mock/room.go -package mock_main -source room.go

type room interface {
	Join(room string, endpoint string)
	Leave(room string, endpoint string)
}
