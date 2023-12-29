package backend

import (
	"os"
	"path"
)

type LocalBackendConfig struct {
	Root string
}

type LocalBackend struct {
	config LocalBackendConfig
}

func (b *LocalBackend) Read(path_ string) ([]byte, error) {
	return os.ReadFile(path.Join(b.config.Root, path_))
}

func (b *LocalBackend) Write(path_ string, data []byte) error {
	path_ = path.Join(b.config.Root, path_)
	err := os.MkdirAll(path.Dir(path_), 0o0755)
	if err != nil {
		return err
	}
	return os.WriteFile(path_, data, 0o0644)
}
