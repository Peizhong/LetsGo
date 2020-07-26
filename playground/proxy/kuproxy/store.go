package main

type store interface {
	Set(string, string)
	Remove(string)
}
