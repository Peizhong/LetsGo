package main

type room interface {
	Join(string)
	Leave(string)
}
