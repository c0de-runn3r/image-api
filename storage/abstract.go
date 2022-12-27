package storage

type StorageInterface interface {
	AddPair(id string, name string)
	GetValue(id string) string
}

type StorageClient struct {
	StorageInterface
}

func NewStorageClient(storage StorageInterface) *StorageClient {
	return &StorageClient{storage}
}
