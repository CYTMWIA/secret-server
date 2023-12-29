package backend

type StorageBackend interface {
	Read(path string) ([]byte, error)
	Write(path string, data []byte) error
}

var DefaultStorageBackend StorageBackend = &LocalBackend{config: LocalBackendConfig{Root: "data"}}
