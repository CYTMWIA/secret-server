package main

import (
	"os"
	"path"
)

type StorageBackend interface {
	Read(path string) ([]byte, error)
	Write(path string, data []byte) error
}

type LocalBackend struct {
	root string
}

var DefaultStorageBackend StorageBackend = LocalBackend{root: "data"}

func (b LocalBackend) Read(path_ string) ([]byte, error) {
	return os.ReadFile(path.Join(b.root, path_))
}

func (b LocalBackend) Write(path_ string, data []byte) error {
	return os.WriteFile(path.Join(b.root, path_), data, 0o0644)
}
