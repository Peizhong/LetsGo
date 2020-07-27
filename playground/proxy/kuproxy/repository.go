package main

type Repository interface {
	Create(key string, value string) error
	Delete(key string) error
}
