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
	path_ = path.Join(b.root, path_)
	err := os.MkdirAll(path.Dir(path_), 0o0755)
	if err != nil {
		return err
	}
	return os.WriteFile(path_, data, 0o0644)
}
