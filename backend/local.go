package backend

import (
	"os"
	"path"
)

type LocalBackend struct {
	root string
}

func (b LocalBackend) Read(path_ string) ([]byte, error) {
	return os.ReadFile(path.Join(b.root, path_))
}

func (b LocalBackend) Write(path_ string, data []byte) error {
	return os.WriteFile(path.Join(b.root, path_), data, 0o0644)
}
