package main

import "sync"

// 储存room和server的信息
type repository interface {
	Store(key string, value string) error
	Delete(key string) error
}

type LocalRepository struct {
	m sync.Map
}

func NewRepository() repository {
	return &LocalRepository{}
}

func (l *LocalRepository) Store(key string, value string) error {
	l.m.Store(key, value)
	return nil
}

func (l *LocalRepository) Delete(key string) error {
	l.m.Delete(key)
	return nil
}
